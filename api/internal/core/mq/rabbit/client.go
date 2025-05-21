package rabbit

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"sync"
)

// Client RabbitMQ 客户端
type Client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	config  Config
	mu      sync.Mutex
}

// NewClient 创建新的 RabbitMQ 客户端
func NewClient(config Config) (*Client, error) {
	conn, err := amqp.Dial(config.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("failed to open channel: %v", err)
	}

	// 声明交换机
	err = channel.ExchangeDeclare(
		config.Exchange,
		"direct", // 类型
		true,     // 持久化
		false,    // 自动删除
		false,    // 内部
		false,    // 非阻塞
		nil,      // 参数
	)
	if err != nil {
		_ = channel.Close()
		_ = conn.Close()
		return nil, fmt.Errorf("failed to declare exchange: %v", err)
	}

	return &Client{
		conn:    conn,
		channel: channel,
		config:  config,
	}, nil
}

// Close 关闭连接
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var errs []error

	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			errs = append(errs, err)
		}
		c.channel = nil
	}

	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			errs = append(errs, err)
		}
		c.conn = nil
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

// Channel 获取 channel (线程安全)
func (c *Client) Channel() (*amqp.Channel, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.channel == nil {
		return nil, errors.New("channel is closed")
	}
	return c.channel, nil
}