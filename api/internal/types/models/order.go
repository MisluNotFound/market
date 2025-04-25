package models

import "time"

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