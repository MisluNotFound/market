package rabbitmq

import (
	"context"
	"fmt"
	"time"

	"github.com/mislu/market-api/internal/core/mq"
	"github.com/streadway/amqp"
)

// RabbitMQQueue implements a RabbitMQ-based message queue
type RabbitMQQueue struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
}

func NewRabbitMQQueue(url, queue string) (*RabbitMQQueue, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	_, err = ch.QueueDeclare(
		queue,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // args
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	return &RabbitMQQueue{
		conn:    conn,
		channel: ch,
		queue:   queue,
	}, nil
}

func (q *RabbitMQQueue) Publish(ctx context.Context, message mq.Message) error {
	err := q.channel.Publish(
		"",      // exchange
		q.queue, // routing key
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType:  "application/octet-stream",
			DeliveryMode: amqp.Persistent,
			MessageId:    message.ID,
			Body:         message.Content,
			Timestamp:    time.Now(),
		})
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}
	return nil
}

func (q *RabbitMQQueue) Consume(ctx context.Context) (<-chan mq.Message, error) {
	msgs, err := q.channel.Consume(
		q.queue, // queue
		"",      // consumer
		true,    // autoAck
		false,   // exclusive
		false,   // noLocal
		false,   // noWait
		nil,     // args
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start consumer: %w", err)
	}

	messageChan := make(chan mq.Message)
	go func() {
		defer close(messageChan)
		for {
			select {
			case msg, ok := <-msgs:
				if !ok {
					return
				}
				messageChan <- mq.Message{
					ID:      msg.MessageId,
					Content: msg.Body,
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return messageChan, nil
}

func (q *RabbitMQQueue) Close() error {
	var channelErr, connErr error
	if err := q.channel.Close(); err != nil {
		channelErr = fmt.Errorf("failed to close channel: %w", err)
	}
	if err := q.conn.Close(); err != nil {
		connErr = fmt.Errorf("failed to close connection: %w", err)
	}
	if channelErr != nil || connErr != nil {
		return fmt.Errorf("errors closing rabbitmq queue: channel=%v, connection=%v", channelErr, connErr)
	}
	return nil
}
