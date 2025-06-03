package kafka

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/mislu/market-api/internal/core/mq"
)

type KafkaQueue struct {
	producer sarama.SyncProducer
	consumer sarama.Consumer
	topic    string
}

func NewKafkaQueue(brokers []string, topic string) (*KafkaQueue, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka producer: %w", err)
	}

	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		producer.Close()
		return nil, fmt.Errorf("failed to create kafka consumer: %w", err)
	}

	return &KafkaQueue{
		producer: producer,
		consumer: consumer,
		topic:    topic,
	}, nil
}

func (q *KafkaQueue) Publish(ctx context.Context, message mq.Message) error {
	_, _, err := q.producer.SendMessage(&sarama.ProducerMessage{
		Topic: q.topic,
		Key:   sarama.StringEncoder(message.ID),
		Value: sarama.ByteEncoder(message.Content),
	})
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}
	return nil
}

func (q *KafkaQueue) Consume(ctx context.Context) (<-chan mq.Message, error) {
	partitionConsumer, err := q.consumer.ConsumePartition(q.topic, 0, sarama.OffsetNewest)
	if err != nil {
		return nil, fmt.Errorf("failed to start consumer: %w", err)
	}

	messageChan := make(chan mq.Message)
	go func() {
		defer partitionConsumer.Close()
		for {
			select {
			case msg, ok := <-partitionConsumer.Messages():
				if !ok {
					close(messageChan)
					return
				}
				messageChan <- mq.Message{
					ID:      string(msg.Key),
					Content: msg.Value,
				}
			case <-ctx.Done():
				close(messageChan)
				return
			}
		}
	}()
	return messageChan, nil
}

func (q *KafkaQueue) Close() error {
	var producerErr, consumerErr error
	if err := q.producer.Close(); err != nil {
		producerErr = fmt.Errorf("failed to close producer: %w", err)
	}
	if err := q.consumer.Close(); err != nil {
		consumerErr = fmt.Errorf("failed to close consumer: %w", err)
	}
	if producerErr != nil || consumerErr != nil {
		return fmt.Errorf("errors closing kafka queue: producer=%v, consumer=%v", producerErr, consumerErr)
	}
	return nil
}
