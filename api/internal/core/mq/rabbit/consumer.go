package rabbit

import (
	"context"
	"fmt"
	"github.com/streadway/amqp"
	"sync"
	"time"
)

// Consumer RabbitMQ 消息消费者
type Consumer struct {
	client      *Client
	queueName   string
	consumerTag string
	handler     HandlerFunc
	options     ConsumerOptions
	wg          sync.WaitGroup
}

// HandlerFunc 消息处理函数
type HandlerFunc func(ctx context.Context, delivery *amqp.Delivery) error

// ConsumerOptions 消费者选项
type ConsumerOptions struct {
	AutoAck      bool          // 是否自动确认消息
	Exclusive    bool          // 是否排他消费者
	NoLocal      bool          // 是否排除本地消息
	PrefetchCount int          // 预取数量
	PrefetchSize  int          // 预取大小
	RetryPolicy  RetryPolicy   // 重试策略
}

// RetryPolicy 重试策略
type RetryPolicy struct {
	MaxAttempts int           // 最大重试次数
	Delay       time.Duration // 重试延迟
}

// NewConsumer 创建新的消费者
func NewConsumer(client *Client, queueName string, handler HandlerFunc, opts ...func(*ConsumerOptions)) *Consumer {
	options := ConsumerOptions{
		AutoAck:      false,
		Exclusive:    false,
		NoLocal:      false,
		PrefetchCount: 1,
		PrefetchSize:  0,
		RetryPolicy: RetryPolicy{
			MaxAttempts: 3,
			Delay:       5 * time.Second,
		},
	}

	for _, opt := range opts {
		opt(&options)
	}

	return &Consumer{
		client:    client,
		queueName: queueName,
		handler:   handler,
		options:   options,
	}
}

// Start 开始消费消息
func (c *Consumer) Start(ctx context.Context) error {
	channel, err := c.client.Channel()
	if err != nil {
		return fmt.Errorf("failed to get channel: %v", err)
	}

	// 设置QoS
	if err := channel.Qos(
		c.options.PrefetchCount,
		c.options.PrefetchSize,
		false,
	); err != nil {
		return fmt.Errorf("failed to set QoS: %v", err)
	}

	// 声明队列
	queue, err := channel.QueueDeclare(
		c.queueName,
		true,  // 持久化
		false, // 自动删除
		false, // 排他
		false, // 非阻塞
		nil,   // 参数
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %v", err)
	}

	// 绑定队列到交换机
	err = channel.QueueBind(
		queue.Name,
		c.queueName, // 使用队列名作为路由键
		c.client.config.Exchange,
		false, // 非阻塞
		nil,   // 参数
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue: %v", err)
	}

	// 开始消费
	deliveries, err := channel.Consume(
		queue.Name,
		c.consumerTag,
		c.options.AutoAck,
		c.options.Exclusive,
		c.options.NoLocal,
		false, // 非阻塞
		nil,   // 参数
	)
	if err != nil {
		return fmt.Errorf("failed to consume: %v", err)
	}

	c.wg.Add(1)
	go c.consumeMessages(ctx, channel, deliveries)

	return nil
}

// consumeMessages 消费消息
func (c *Consumer) consumeMessages(ctx context.Context, channel *amqp.Channel, deliveries <-chan amqp.Delivery) {
	defer c.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case delivery, ok := <-deliveries:
			if !ok {
				return
			}

			c.processMessage(ctx, channel, &delivery)
		}
	}
}

// processMessage 处理单个消息
func (c *Consumer) processMessage(ctx context.Context, channel *amqp.Channel, delivery *amqp.Delivery) {
	var attempt int
	var err error

	for attempt = 0; attempt < c.options.RetryPolicy.MaxAttempts; attempt++ {
		if attempt > 0 {
			time.Sleep(c.options.RetryPolicy.Delay)
		}

		err = c.handler(ctx, delivery)
		if err == nil {
			break
		}
	}

	if !c.options.AutoAck {
		if err != nil {
			// 处理失败，拒绝消息
			_ = delivery.Reject(false) // 不重新入队
		} else {
			// 处理成功，确认消息
			_ = delivery.Ack(false)
		}
	}
}

// Stop 停止消费者
func (c *Consumer) Stop() error {
	c.wg.Wait()
	return nil
}