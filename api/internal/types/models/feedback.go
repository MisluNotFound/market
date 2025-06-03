package models

import "time"

type Feedback struct {
	ID        string    `gorm:"column:id;type:varchar(36);primary_key" json:"id"`
	UserID    string    `gorm:"column:user_id;varchar(36)" json:"userID"`
	ItemID    string    `gorm:"column:item_id;varchar(36)" json:"itemID"`
	Timestamp time.Time `json:"timestamp"`
}

func (Feedback) TableName() string {
	return "feedback"
}
