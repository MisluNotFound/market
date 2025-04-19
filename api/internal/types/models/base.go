package models

import "time"

type Model struct {
	ID        string `gorm:"column:id;primaryKey;type:varchar(36);default:(uuid())" json:"id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
