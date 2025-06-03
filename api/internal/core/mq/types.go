package mq

import (
	"context"
	"errors"
)

var (
	ErrClosed = errors.New("mq: connection closed")
)

type Message struct {
	ID      string
	Content []byte
}

// Queue interface defines the message queue operations
type Queue interface {
	Publish(ctx context.Context, message Message) error
	Consume(ctx context.Context) (<-chan Message, error)
	Close() error
}
