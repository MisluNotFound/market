package models

type User struct {
	Model

	Username        string `gorm:"column:username;type:varchar(50);not null;" json:"username"`
	Password        string `gorm:"column:password;type:varchar(255);not null;" json:"-"`
	School          string `gorm:"column:school;type:varchar(50);not null;" json:"school"`
	Phone           string `gorm:"column:phone;type:varchar(50);not null;" json:"phone"`
	Gender          string `gorm:"column:gender;type:varchar(10)" json:"gender"`
	Avatar          string `gorm:"column:avatar;type:varchar(100)" json:"avatar"`
	Salt            string `gorm:"column:salt;type:varchar(100)" json:"-"`
	SellerCredit    int    `gorm:"column:seller_credit;type:int;default:0;" json:"sellerCredit"`
	PurchaserCredit int    `gorm:"column:purchaser_credit;type:int;default:0;" json:"purchaseCredit"`

	// Is user certificated, in other words, is user a student
	IsCertificated bool `gorm:"column:is_certificated;type:bool;default:false;" json:"isCertificated"`

	// TODO Is it useful?
	Address string `gorm:"column:address;type:varchar(255);not null;" json:"address"`
}

func (User) TableName() string {
	return "user"
}

func (u User) Exists() bool {
	return len(u.ID) > 0
}
