package models

type Credit struct {
	UserID          string  `gorm:"column:user_id;type:varchar(36);not null;primary_key" json:"userID"`
	TotalComment    int     `gorm:"column:total_comment;type:int;not null" json:"totalComment"`
	PositiveComment int     `gorm:"column:positive_comment;type:int;not null" json:"positiveComment"`
	NegativeComment int     `gorm:"column:negative_comment;type:int;not null" json:"negativeComment"`
	Reputation      float64 `gorm:"column:reputation;type:decimal(10,2);not null" json:"reputation"`
}

func (Credit) TableName() string {
	return "credit"
}
