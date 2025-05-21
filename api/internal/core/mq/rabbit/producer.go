package rabbit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"time"
)

// Producer RabbitMQ 消息生产者
type Producer struct {
	client *Client
}

// NewProducer 创建新的生产者
func NewProducer(client *Client) *Producer {
	return &Producer{client: client}
}

// Publish 发布消息
func (p *Producer) Publish(ctx context.Context, routingKey string, message interface{}) error {
	channel, err := p.client.Channel()
	if err != nil {
		return fmt.Errorf("failed to get channel: %v", err)
	}

	// 序列化消息
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	// 配置发布选项
	publishing := amqp.Publishing{
		ContentType:  "application/json",
		Body:         body,
		DeliveryMode: amqp.Persistent, // 持久化消息
		Timestamp:    time.Now(),
	}

	// 如果配置了死信队列，添加相关参数
	if p.client.config.DLX.Enabled {
		publishing.Headers = amqp.Table{
			"x-dead-letter-exchange":    p.client.config.DLX.Exchange,
			"x-dead-letter-routing-key": p.client.config.DLX.RoutingKey,
		}
		if p.client.config.DLX.TTL > 0 {
			publishing.Expiration = fmt.Sprintf("%d", p.client.config.DLX.TTL.Milliseconds())
		}
	}

	// 发布消息
	err = channel.Publish(
		p.client.config.Exchange,
		routingKey,
		false, // 强制
		false, // 立即
		publishing,
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}

	return nil
}

// PublishWithConfirm 发布消息并等待确认
func (p *Producer) PublishWithConfirm(ctx context.Context, routingKey string, message interface{}) error {
	channel, err := p.client.Channel()
	if err != nil {
		return fmt.Errorf("failed to get channel: %v", err)
	}

	// 启用发布确认
	if err := channel.Confirm(false); err != nil {
		return fmt.Errorf("failed to put channel in confirm mode: %v", err)
	}

	confirms := channel.NotifyPublish(make(chan amqp.Confirmation, 1))

	err = p.Publish(ctx, routingKey, message)
	if err != nil {
		return err
	}

	select {
	case confirm := <-confirms:
		if !confirm.Ack {
			return errors.New("failed to deliver message to broker")
		}
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(5 * time.Second):
		return errors.New("timeout waiting for publish confirmation")
	}

	return nil
}