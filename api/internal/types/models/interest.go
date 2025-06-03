package models

import "time"

type InterestTag struct {
	ID         int    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TagName    string `gorm:"column:tag_name;type:varchar(50);not null;unique" json:"tagName"`
	CategoryID int    `gorm:"column:category_id;type:int;not null" json:"categoryID"`
}

func (InterestTag) TableName() string {
	return "interest_tags"
}

func (tag InterestTag) Exists() bool {
	return tag.ID > 0
}

type UserInterests struct {
	ID            int       `gorm:"id;type:bigint;not null;primary_key" json:"-"`
	UserID        string    `gorm:"column:user_id;type:varchar(36);not null" json:"userID"`
	InterestTagID int       `gorm:"column:interest_tag_id;type:int;not null" json:"interestTagID"`
	SelectedAt    time.Time `gorm:"column:selected_at;type:datetime;not null" json:"selectedAt"`
	UpdatedAt     time.Time `gorm:"column:updated_at;type:datetime;not null" json:"updatedAt"`
}

func (UserInterests) TableName() string {
	return "user_interests"
}
