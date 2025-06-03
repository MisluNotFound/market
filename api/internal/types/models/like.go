package models

type Like struct {
	UserID    string `gorm:"column:user_id;type:varchar(36);not null;primary_key"`
	ProductID string `gorm:"column:product_id;type:varchar(36);not null;primary_key"`
}

func (Like) TableName() string {
	return "like"
}

func (l Like) Exists() bool {
	return len(l.UserID) > 0 && len(l.ProductID) > 0
}
