package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	Model
	UserID        string  `gorm:"column:user_id;type:varchar(50);not null;" json:"userID"`
	OriginalPrice float64 `gorm:"column:original_price;type:decimal(10,2);not null;" json:"originalPrice"`
	Price         float64 `gorm:"column:price;type:decimal(10,2);not null;" json:"price"`
	Describe      string  `gorm:"column:describe;type:varchar(255);not null;" json:"describe"`
	Pics          string  `gorm:"column:pics;type:varchar(500);not null;" json:"avatar"`
	Condition     string  `gorm:"column:condition;varchar(20)" json:"condition"`
	UsedTime      string  `gorm:"column:usedTime;varchar(20)" json:"usedTime"`

	ShippingMethod string    `gorm:"column:shipping_method;type:varchar(50);not null;" json:"shippingMethod"`
	ShippingPrise  float64   `gorm:"column:shipping_price;type:decimal(10,2);not null;" json:"shippingPrice"`
	CanSelfPickup  bool      `gorm:"column:can_self_pickup;type:bool;default:false;" json:"canSelfPickup"`
	Location       string    `gorm:"column:location;type:varchar(100);not null;" json:"location"`
	PublishAt      time.Time `json:"publishAt"`

	IsPublished bool `gorm:"column:is_published;type:bool;default:false;" json:"isPublished"` // 是否通过审核
	IsSold      bool `gorm:"column:is_sold;type:bool;default:false;" json:"isSold"`           // 是否售出
	IsSelling   bool `gorm:"column:is_selling;type:bool;default:true;" json:"isSelling"`      // 是否下架
}

func (Product) TableName() string {
	return "product"
}

func (p Product) Exists() bool {
	return len(p.ID) > 0
}

func (p Product) IsOwner(userID string) bool {
	return p.UserID == userID
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}
