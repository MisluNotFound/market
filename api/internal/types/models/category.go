package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Category struct {
	TypeID     int       `gorm:"primaryKey;column:type_id"`
	TypeName   string    `gorm:"size:50;not null"`
	ParentID   int       `gorm:"default:0"`
	Level      int8      `gorm:"check:level BETWEEN 1 AND 4"`
	IsLeaf     bool      `gorm:"default:false"`
	Attributes JSONArray `gorm:"type:json"`
}

func (Category) TableName() string {
	return "product_category"
}

type AttributeTemplate struct {
	AttributeID int       `gorm:"primaryKey;column:attribute_id"`
	Name        string    `gorm:"size:40;not null"`
	DataType    DataType  `gorm:"type:ENUM('STRING','NUMBER','ENUM','BOOLEAN','DATE');not null"`
	Required    bool      `gorm:"default:false"`
	Options     JSONArray `gorm:"type:json"`
	Unit        string    `gorm:"size:10"`
	BindRules   JSONMap   `gorm:"type:json"`
}

func (AttributeTemplate) TableName() string {
	return "attribute_template"
}

type CategoryAttribute struct {
	TypeID           int `gorm:"primaryKey"`
	AttributeID      int `gorm:"primaryKey"`
	RequiredOverride *bool
}

func (CategoryAttribute) TableName() string {
	return "category_attribute"
}

type ProductAttribute struct {
	ItemID      string    `gorm:"primaryKey;size:20"`
	AttributeID int       `gorm:"primaryKey"`
	StrValue    string    `gorm:"size:255"`
	NumValue    float64   `gorm:"type:decimal(15,3)"`
	DateValue   time.Time `gorm:"type:date"`
}

func (ProductAttribute) TableName() string {
	return "product_attribute"
}

type DataType string

const (
	DataTypeString  DataType = "STRING"
	DataTypeNumber  DataType = "NUMBER"
	DataTypeEnum    DataType = "ENUM"
	DataTypeBoolean DataType = "BOOLEAN"
	DataTypeDate    DataType = "DATE"
)

type JSONArray []int

func (j *JSONArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, j)
}

func (j JSONArray) Value() (driver.Value, error) {
	return json.Marshal(j)
}

type JSONMap map[string]interface{}

func (j *JSONMap) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, j)
}

func (j JSONMap) Value() (driver.Value, error) {
	return json.Marshal(j)
}
