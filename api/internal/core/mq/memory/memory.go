package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/mislu/market-api/internal/core/mq"
)

// InMemoryQueue implements an in-memory message queue
type InMemoryQueue struct {
	messages chan mq.Message
	wg       sync.WaitGroup
	closed   bool
	mu       sync.Mutex
}

func NewInMemoryQueue(bufferSize int) *InMemoryQueue {
	return &InMemoryQueue{
		messages: make(chan mq.Message, bufferSize),
	}
}

func (q *InMemoryQueue) Publish(ctx context.Context, message mq.Message) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.closed {
		return errors.New("queue is closed")
	}
	select {
	case q.messages <- message:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (q *InMemoryQueue) Consume(ctx context.Context) (<-chan mq.Message, error) {
	if q.closed {
		return nil, errors.New("queue is closed")
	}
	return q.messages, nil
}

func (q *InMemoryQueue) Close() error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if !q.closed {
		q.closed = true
		close(q.messages)
	}
	return nil
}
