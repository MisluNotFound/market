package rabbit

import "time"

// Config RabbitMQ 连接配置
type Config struct {
	URL      string        // AMQP 连接URL
	Exchange string        // 交换机名称
	Queue    QueueConfig   // 队列配置
	DLX      DLXConfig     // 死信队列配置
}

// QueueConfig 队列配置
type QueueConfig struct {
	Name       string        // 队列名称
	Durable    bool          // 是否持久化
	AutoDelete bool          // 是否自动删除
	Exclusive  bool          // 是否排他队列
	NoWait     bool          // 是否非阻塞
	Args       map[string]interface{} // 额外参数
}

// DLXConfig 死信队列配置
type DLXConfig struct {
	Enabled    bool          // 是否启用死信队列
	Exchange   string        // 死信交换机
	Queue      string        // 死信队列
	RoutingKey string        // 路由键
	TTL        time.Duration // 消息TTL
}