package rabbit

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/mislu/market-api/internal/db"
	"github.com/mislu/market-api/internal/types/models"
	"github.com/mislu/market-api/internal/utils/app"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	OrderQueueName     = "order_queue"
	DeadLetterQueue    = "dead_letter_queue"
	ExchangeName       = "order_exchange"
	DeadLetterExchange = "dead_letter_exchange"
)

var rabbitCli *RabbitMQClient

type RabbitMQClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func InitGlobalRabbitMQ() error {
	var err error
	rabbitCli, err = NewRabbitMQClient(app.GetConfig().Rabbit.Url)
	if err != nil {
		return err
	}
	err = rabbitCli.SetupQueues()
	if err != nil {
		return err
	}
	err = rabbitCli.ConsumeDeadLetterQueue()
	if err != nil {
		return err
	}

	return rabbitCli.ConsumeDeadLetterQueue()
}

func NewRabbitMQClient(url string) (*RabbitMQClient, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %v", err)
	}

	return &RabbitMQClient{
		conn:    conn,
		channel: channel,
	}, nil
}

func (c *RabbitMQClient) SetupQueues() error {
	err := c.channel.ExchangeDeclare(
		DeadLetterExchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare dead letter exchange: %v", err)
	}

	_, err = c.channel.QueueDeclare(
		DeadLetterQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare dead letter queue: %v", err)
	}

	err = c.channel.QueueBind(
		DeadLetterQueue,
		"",
		DeadLetterExchange,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind dead letter queue: %v", err)
	}

	err = c.channel.ExchangeDeclare(
		ExchangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %v", err)
	}

	args := amqp.Table{
		"x-dead-letter-exchange": DeadLetterExchange,
		"x-message-ttl":          int32(15 * 60 * 1000), // TTL 15min
	}

	_, err = c.channel.QueueDeclare(
		OrderQueueName,
		true,
		false,
		false,
		false,
		args,
	)
	if err != nil {
		return fmt.Errorf("failed to declare order queue: %v", err)
	}

	err = c.channel.QueueBind(
		OrderQueueName,
		"",
		ExchangeName,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind order queue: %v", err)
	}

	return nil
}

func SendOrder(order models.Order) error {
	body, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to marshal order: %v", err)
	}

	err = rabbitCli.channel.Publish(
		ExchangeName,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}

	return nil
}

func (c *RabbitMQClient) ConsumeDeadLetterQueue() error {
	msgs, err := c.channel.Consume(
		DeadLetterQueue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %v", err)
	}

	go func() {
		for msg := range msgs {
			var order models.Order
			if err := json.Unmarshal(msg.Body, &order); err != nil {
				log.Printf("Failed to unmarshal order: %v", err)
				msg.Nack(false, false)
				continue
			}

			err = updateOrderStatus(order.ID, "timeout")
			if err != nil {
				log.Printf("Failed to update order status: %v", err)
				msg.Nack(false, true)
				continue
			}

			log.Printf("Processed timeout order: %s", order.ID)
			msg.Ack(false)
		}
	}()

	return nil
}

func (c *RabbitMQClient) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}

func updateOrderStatus(orderID, status string) error {
	order, err := db.GetOne[models.Order](
		db.Equal("id", orderID),
	)

	if err != nil {
		return err
	}
	if order.Status != 1 {
		return nil
	}

	order.Status = 8
	return db.Update(&order)
}
