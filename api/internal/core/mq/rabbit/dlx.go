package rabbit

import (
	"context"
	"fmt"
	"github.com/streadway/amqp"
)

// DLXHandler 死信队列处理器
type DLXHandler struct {
	client *Client
}

// NewDLXHandler 创建死信队列处理器
func NewDLXHandler(client *Client) (*DLXHandler, error) {
	if !client.config.DLX.Enabled {
		return nil, fmt.Errorf("DLX is not enabled in config")
	}

	// 初始化死信交换机和队列
	if err := initDLXResources(client); err != nil {
		return nil, fmt.Errorf("failed to initialize DLX resources: %v", err)
	}

	return &DLXHandler{client: client}, nil
}

// initDLXResources 初始化死信队列资源
func initDLXResources(client *Client) error {
	channel, err := client.Channel()
	if err != nil {
		return fmt.Errorf("failed to get channel: %v", err)
	}

	// 声明死信交换机
	err = channel.ExchangeDeclare(
		client.config.DLX.Exchange,
		"direct", // 类型
		true,     // 持久化
		false,    // 自动删除
		false,    // 内部
		false,    // 非阻塞
		nil,      // 参数
	)
	if err != nil {
		return fmt.Errorf("failed to declare DLX exchange: %v", err)
	}

	// 声明死信队列
	_, err = channel.QueueDeclare(
		client.config.DLX.Queue,
		true,  // 持久化
		false, // 自动删除
		false, // 排他
		false, // 非阻塞
		amqp.Table{
			"x-message-ttl": client.config.DLX.TTL.Milliseconds(),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to declare DLX queue: %v", err)
	}

	// 绑定死信队列到交换机
	err = channel.QueueBind(
		client.config.DLX.Queue,
		client.config.DLX.RoutingKey,
		client.config.DLX.Exchange,
		false, // 非阻塞
		nil,   // 参数
	)
	if err != nil {
		return fmt.Errorf("failed to bind DLX queue: %v", err)
	}

	return nil
}

// HandleDeadLetters 处理死信消息
func (h *DLXHandler) HandleDeadLetters(ctx context.Context, handler func(context.Context, *amqp.Delivery) error) error {
	consumer := NewConsumer(h.client, h.client.config.DLX.Queue, handler, func(o *ConsumerOptions) {
		o.AutoAck = false // 死信队列需要手动确认
	})

	return consumer.Start(ctx)
}

// GetDLXConfig 获取死信队列配置
func (h *DLXHandler) GetDLXConfig() DLXConfig {
	return h.client.config.DLX
}