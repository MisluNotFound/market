package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	Model
	ProductID   string    `gorm:"column:product_id;type:varchar(36);not null" json:"productID"`
	UserID      string    `gorm:"column:user_id;type:varchar(36);not null" json:"userID"`
	SellerID    string    `gorm:"column:seller_id;type:varchar(36);not null" json:"sellerID"`
	Status      int       `gorm:"column:status;type:int;not null" json:"status"`
	TotalAmount float64   `gorm:"column:total_amount;type:decimal(10,2);not null" json:"totalAmount"`
	ShipAmount  float64   `gorm:"column:ship_amount;type:decimal(10,2);not null" json:"shipAmount"`
	PayTime     time.Time `json:"payTime"`
	ShipTime    time.Time `json:"shipTime"`
	FinishTime  time.Time `json:"finishTime"`
	IsEvaluated bool      `gorm:"column:is_evaluated;type:bool;not null" json:"isEvaluated"`
	PayMethod   string    `gorm:"column:pay_method;type:varchar(36);not null" json:"payMethod"`
}

func (Order) TableName() string {
	return "order"
}

func (o Order) Exists() bool {
	return len(o.ID) > 0
}

func (o Order) IsOwner(userID string) bool {
	return o.UserID == userID
}

func (o Order) IsSeller(userID string) bool {
	return o.SellerID == userID
}

type OrderComment struct {
	ID        int       `gorm:"column:id;type:bigint;primary_key;auto_increment" json:"id"`
	OrderID   string    `gorm:"column:order_id;type:varchar(36);not null" json:"orderID"`
	UserID    string    `gorm:"column:user_id;type:varchar(36);not null" json:"userID"`
	ProductID string    `gorm:"column:product_id;type:varchar(36);not null" json:"productID"`
	Comment   string    `gorm:"column:comment;type:varchar(255);not null" json:"comment"`
	IsGood    bool      `gorm:"column:is_good;type:bool;" json:"isGood"`
	ParentID  int       `gorm:"column:parent_id;type:bigint;" json:"parentID"`
	ReplyTo   string    `gorm:"column:reply_to;type:varchar(36);" json:"replyTo"` // user's name
	IsTop     bool      `gorm:"column:is_top;type:bool;" json:"isTop"`
	Pics      string    `gorm:"column:pics;type:varchar(500);" json:"pics"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null" json:"updatedAt"`
}

func (OrderComment) TableName() string {
	return "order_comment"
}

func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.ID == "" {
		o.ID = uuid.New().String()
	}
	return nil
}
