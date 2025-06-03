package recommend

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/mislu/market-api/internal/core/mq"
	"github.com/mislu/market-api/internal/core/mq/memory"
	"github.com/mislu/market-api/internal/utils/app"
	"github.com/zhenghaoz/gorse/client"
)

var GlobalWorker *RecommendationWorker

// Feedback represents the structure of a recommendation feedback message
type Feedback struct {
	UserId       string `json:"userId"`
	ItemId       string `json:"itemId"`
	FeedbackType string `json:"feedbackType"`
	Timestamp    int64  `json:"timestamp"`
}

// RecommendationWorker receives messages and forwards them to Gorse
type RecommendationWorker struct {
	consumer    mq.Queue
	gorseClient *client.GorseClient
}

func NewRecommendationWorker(consumer mq.Queue) *RecommendationWorker {
	gorseConfig := app.GetConfig().Gorse
	gorseClient := client.NewGorseClient(gorseConfig.Endpoint, gorseConfig.ApiKey)
	return &RecommendationWorker{
		consumer:    consumer,
		gorseClient: gorseClient,
	}
}

func InitGlobalWorker() {
	gorseConfig := app.GetConfig().Gorse
	var (
		err   error
		queue mq.Queue
	)

	switch gorseConfig.MQ.Type {
	case "memory":
		queue = memory.NewInMemoryQueue(gorseConfig.MQ.Memory.Size)
		// TODO: support other queue
	}

	if err != nil {
		panic(err)
	}

	GlobalWorker = NewRecommendationWorker(queue)
	go GlobalWorker.Work(context.Background())
}

// Work starts consuming messages from the queue and processes them with Gorse
func (w *RecommendationWorker) Work(ctx context.Context) error {
	// Start consuming messages
	msgChan, err := w.consumer.Consume(ctx)
	if err != nil {
		return fmt.Errorf("failed to start consumer: %w", err)
	}

	// Process messages in a loop
	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				log.Println("Consumer channel closed")
				return nil
			}
			if err := w.processMessage(ctx, msg); err != nil {
				log.Printf("Failed to process message %s: %v", msg.ID, err)
				// Optionally implement retry logic or dead letter queue
				continue
			}
			log.Printf("Successfully processed message %s", msg.ID)

		case <-ctx.Done():
			log.Println("Context canceled, stopping worker")
			return ctx.Err()
		}
	}
}

// processMessage handles a single message by parsing it and sending to Gorse
func (w *RecommendationWorker) processMessage(ctx context.Context, msg mq.Message) error {
	// Parse message content
	var feedback Feedback
	if err := json.Unmarshal(msg.Content, &feedback); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	// Validate feedback
	if feedback.UserId == "" || feedback.ItemId == "" || feedback.FeedbackType == "" {
		return fmt.Errorf("invalid feedback data: userId=%s, itemId=%s, feedbackType=%s",
			feedback.UserId, feedback.ItemId, feedback.FeedbackType)
	}

	// Send feedback to Gorse with retry
	const maxRetries = 3
	for attempt := 1; attempt <= maxRetries; attempt++ {
		_, err := w.gorseClient.InsertFeedback(ctx, []client.Feedback{{
			FeedbackType: feedback.FeedbackType,
			UserId:       feedback.UserId,
			ItemId:       feedback.ItemId,
			Timestamp:    time.Now().Format(time.RFC3339),
		}})
		if err == nil {
			return nil
		}
		log.Printf("Attempt %d failed to send feedback to Gorse: %v", attempt, err)
		if attempt < maxRetries {
			select {
			case <-time.After(time.Second * time.Duration(attempt)):
				continue
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
	return fmt.Errorf("failed to send feedback to Gorse after %d attempts", maxRetries)
}

func (w *RecommendationWorker) InsertFeedback(ctx context.Context, feedbacks []Feedback) error {
	for _, feedback := range feedbacks {
		msg, err := json.Marshal(feedback)
		if err != nil {
			continue
		}

		w.consumer.Publish(ctx, mq.Message{
			ID:      fmt.Sprintf("%s-%s-%s", feedback.UserId, feedback.ItemId, feedback.FeedbackType),
			Content: msg,
		})
	}

	return nil
}
