package models

type SearchHistory struct {
	ID         int    `gorm:"column:id;type:bigint;primary_key;auto_increment" json:"-"`
	UserID     string `gorm:"column:user_id;type:varchar(36);not null;index:idx_user_id"`
	Keyword    string `gorm:"column:keyword;type:varchar(255);not null"`
	SearchTime int64  `gorm:"column:search_time;type:bigint;not null"`
}

func (SearchHistory) TableName() string {
	return "search_history"
}
