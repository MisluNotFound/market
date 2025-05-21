package rabbit

import (
	"context"
	"testing"
	"time"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
)

func TestRabbitMQIntegration(t *testing.T) {
	// 配置 RabbitMQ
	config := Config{
		URL:      "amqp://guest:guest@localhost:5672/",
		Exchange: "test_exchange",
		Queue: QueueConfig{
			Name:       "test_queue",
			Durable:    true,
			AutoDelete: false,
		},
		DLX: DLXConfig{
			Enabled:    true,
			Exchange:   "dlx_exchange",
			Queue:      "dlx_queue",
			RoutingKey: "dlx_routing_key",
			TTL:        time.Minute * 30,
		},
	}

	// 创建客户端
	client, err := NewClient(config)
	assert.NoError(t, err)
	defer client.Close()

	// 测试生产者
	t.Run("Producer", func(t *testing.T) {
		producer := NewProducer(client)

		// 测试普通消息
		err := producer.Publish(context.Background(), "test_routing_key", "test message")
		assert.NoError(t, err)

		// 测试带确认的消息
		err = producer.PublishWithConfirm(context.Background(), "test_routing_key", "confirmed message")
		assert.NoError(t, err)
	})

	// 测试消费者
	t.Run("Consumer", func(t *testing.T) {
		msgChan := make(chan string, 1)
		handler := func(ctx context.Context, delivery *amqp.Delivery) error {
			msgChan <- string(delivery.Body)
			return nil
		}

		consumer := NewConsumer(client, "test_queue", handler)
		go func() {
			err := consumer.Start(context.Background())
			assert.NoError(t, err)
		}()

		// 等待消息
		select {
		case msg := <-msgChan:
			assert.Equal(t, "test message", msg)
		case <-time.After(5 * time.Second):
			t.Fatal("timeout waiting for message")
		}
	})

	// 测试死信队列
	t.Run("DLX", func(t *testing.T) {
		dlxHandler, err := NewDLXHandler(client)
		assert.NoError(t, err)

		dlxMsgChan := make(chan string, 1)
		dlxHandlerFunc := func(ctx context.Context, delivery *amqp.Delivery) error {
			dlxMsgChan <- string(delivery.Body)
			return nil
		}

		go func() {
			err := dlxHandler.HandleDeadLetters(context.Background(), dlxHandlerFunc)
			assert.NoError(t, err)
		}()

		// 发送会变成死信的消息
		producer := NewProducer(client)
		err = producer.Publish(context.Background(), "invalid_routing_key", "dlx test message")
		assert.NoError(t, err)

		// 等待死信消息
		select {
		case msg := <-dlxMsgChan:
			assert.Equal(t, "dlx test message", msg)
		case <-time.After(5 * time.Second):
			t.Fatal("timeout waiting for DLX message")
		}
	})
}