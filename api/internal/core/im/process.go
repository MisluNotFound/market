package im

import (
	"log"
	"math"
	"time"

	"github.com/mislu/market-api/internal/service"
	"github.com/mislu/market-api/internal/types/request"
)

const (
	send = iota + 1
	ack
	fail
)

func (c *Client) handleMessage(message *request.Message) {
	log.Printf("read message from %s\n", c.userID)
	switch message.Type {
	case send:
		err := sendMessage(message)
		if err != nil {
			handleFail(c, message)
			return
		}
	case ack:
		// ackMessage(message)
	}
}

func handleRetry(client *Client, msgID string) {
	client.pendingMu.Lock()
	defer client.pendingMu.Unlock()

	pm, exists := client.pending[msgID]
	if !exists {
		return
	}

	// 超出最大重试（3次）
	if pm.retries >= 3 {
		log.Printf("消息 %s 达到最大重试次数", msgID)
		delete(client.pending, msgID)
		return
	}

	pm.retries++
	nextRetry := time.Duration(math.Pow(2, float64(pm.retries))) * 500 * time.Millisecond
	pm.timer.Reset(nextRetry)

	go func() {
		if err := client.conn.WriteJSON(pm.Message); err != nil {
			log.Printf("重传失败: %v", err)
		}
	}()
}

// sendMessage handle message type 'send'(1). Save message to database.
func sendMessage(message *request.Message) error {
	err := service.SaveMessage(message)
	if err != nil {
		return err
	}

	toClient, online := getClient(message.To)
	if !online {
		// handle offline
		service.RecordLastReadMessage(message.From, message.To, message.ID)
		log.Println("用户下线", toClient)
		return nil
	}

	toClient.send <- *message
	log.Println("写入用户channel", toClient)
	// toClient.pendingMu.Lock()
	// defer toClient.pendingMu.Unlock()
	// pm := &PendingMessage{
	// 	Message:   *message,
	// 	expiresAt: time.Now().Add(500 * time.Millisecond),
	// }

	// pm.timer = time.AfterFunc(500*time.Millisecond, func() {
	// 	handleRetry(toClient, message.ID)
	// })

	// toClient.pending[message.ID] = pm

	return nil
}

// ackMessage handle message type 'ack'(2). Remove message from pending queue.
func ackMessage(message *request.Message) {
	toClient, online := getClient(message.To)
	if !online {
		// TODO handle offline message
		return
	}
	toClient.pendingMu.Lock()
	defer toClient.pendingMu.Unlock()

	if pm, exists := toClient.pending[message.ID]; exists {
		pm.timer.Stop()
		delete(toClient.pending, message.ID)
	}
}

func handleFail(client *Client, message *request.Message) {
	failedMessage := *message
	failedMessage.ID = failedMessage.TempID
	failedMessage.Type = fail

	client.send <- failedMessage
}
