package models

type User struct {
	Model

	Username     string `gorm:"column:username;type:varchar(50);not null;" json:"username"`
	Password     string `gorm:"column:password;type:varchar(255);not null;" json:"-"`
	Phone        string `gorm:"column:phone;type:varchar(50);not null;" json:"phone"`
	Gender       string `gorm:"column:gender;type:varchar(10)" json:"gender"`
	Avatar       string `gorm:"column:avatar;type:varchar(100)" json:"avatar"`
	Salt         string `gorm:"column:salt;type:varchar(100)" json:"-"`
	SelectedTags bool   `gorm:"column:selected_tags;" json:"selected_tags"`
}

func (User) TableName() string {
	return "user"
}

func (u User) Exists() bool {
	return len(u.ID) > 0
}
