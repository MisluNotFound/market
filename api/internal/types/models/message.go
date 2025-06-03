package models

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	ID             string    `gorm:"primaryKey;type:varchar(50)" json:"id"`
	ConversationID string    `gorm:"column:conversation_id;index;type:varchar(80)" json:"conversation_id"`
	FromUserID     string    `gorm:"column:from_userd_id;index;type:varchar(36)" json:"from_user_id"`
	ToUserID       string    `gorm:"column:to_userd_id;index;type:varchar(36)" json:"to_user_id"`
	Content        string    `gorm:"type:text" json:"content"`
	MediaType      string    `gorm:"type:varchar(20)" json:"media_type"` // text/image/link
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"-"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"-"`
	Timestamp      time.Time `gorm:"timestamp;index" json:"timestamp"`
}

func (Message) TableName() string {
	return "message"
}

type Conversation struct {
	gorm.Model         `json:"-"`
	FromUserID         string    `gorm:"type:varchar(36);index:from_to_user_id" json:"fromUserID"`
	ToUserID           string    `gorm:"type:varchar(36);index:from_to_user_id" json:"toUserID"`
	LastMessageContent string    `gorm:"column:last_message_content;type:varchar(200)" json:"lastMessageContent"`
	LastMessageTime    time.Time `gorm:"column:last_message_time;index" json:"lastMessageTime"`
	LastReadMessageID  string    `gorm:"column:last_read_message_id;type:varchar(36)" json:"lastReadMessageID"` // 记录离线前的最后一条消息
	MarkDeleted        bool      `gorm:"column:mark_deleted;default:false" json:"-"`
}

func (Conversation) TableName() string {
	return "conversation"
}
