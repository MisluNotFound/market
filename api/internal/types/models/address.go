package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Address struct {
	Model
	Address      string  `gorm:"column:address;type:varchar(255);not null" json:"address"`
	City         string  `gorm:"column:city;type:varchar(36);not null" json:"city"`
	District     string  `gorm:"column:district;type:varchar(36);not null" json:"district"`
	Province     string  `gorm:"column:province;type:varchar(36);not null" json:"province"`
	Street       string  `gorm:"column:street;type:varchar(36);not null" json:"street"`
	StreetNumber string  `gorm:"column:street_number;type:varchar(36);not null" json:"streetNumber"`
	Latitude     float64 `gorm:"column:latitude;type:decimal(12,8);not null" json:"latitude"`
	Longitude    float64 `gorm:"column:longitude;type:decimal(12,8);not null" json:"longitude"`
}

type UserAddress struct {
	Model
	UserID    string `gorm:"column:user_id;type:varchar(36);not null" json:"userID"`
	IsDefault bool   `gorm:"column:is_default;type:boolean;not null" json:"isDefault"`
	AddressID string `gorm:"column:address_id;type:varchar(36);not null" json:"addressID"`
	Phone     string `gorm:"column:phone;type:varchar(36);not null" json:"phone"`
	Receiver  string `gorm:"column:receiver;type:varchar(36);not null" json:"receiver"`
	Detail    string `gorm:"column:detail;type:varchar(50);not null" json:"detail"`
}

func (Address) TableName() string {
	return "address"
}

func (a *Address) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	return nil
}

func (u *UserAddress) Exists() bool {
	return len(u.ID) > 0
}
