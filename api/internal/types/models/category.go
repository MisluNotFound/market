package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// 整型数组类型
type JSONIntArray []int

func (j *JSONIntArray) Scan(value interface{}) error {
	return jsonScan(value, j)
}

func (j JSONIntArray) Value() (driver.Value, error) {
	return jsonValue(j)
}

// 字符串数组类型
type JSONStrArray []string

func (j *JSONStrArray) Scan(value interface{}) error {
	return jsonScan(value, j)
}

func (j JSONStrArray) Value() (driver.Value, error) {
	return jsonValue(j)
}

// 规则映射类型
type JSONRuleMap map[string]JSONIntArray

func (j *JSONRuleMap) Scan(value interface{}) error {
	return jsonScan(value, j)
}

func (j JSONRuleMap) Value() (driver.Value, error) {
	return jsonValue(j)
}

// 公共处理方法
func jsonScan(value interface{}, target interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("invalid JSON type")
	}

	if len(bytes) == 0 {
		return nil
	}

	return json.Unmarshal(bytes, target)
}

func jsonValue(value interface{}) (driver.Value, error) {
	if value == nil {
		return nil, nil
	}

	return json.Marshal(value)
}

type DataType string

const (
	DataTypeString  DataType = "STRING"
	DataTypeNumber  DataType = "NUMBER"
	DataTypeEnum    DataType = "ENUM"
	DataTypeBoolean DataType = "BOOLEAN"
	DataTypeDate    DataType = "DATE"
)

// 商品分类表
type Category struct {
	ID         uint         `gorm:"primaryKey;comment:分类ID"`
	TypeName   string       `gorm:"size:50;not null;index;comment:分类名称"`
	ParentID   uint         `gorm:"default:0;index;comment:父分类ID" json:"-"`
	Level      int8         `gorm:"check:level BETWEEN 1 AND 4;comment:层级"`
	IsLeaf     bool         `gorm:"default:false;comment:是否叶子节点"`
	Attributes JSONIntArray `gorm:"type:json;comment:绑定属性ID列表" json:"-"`
}

func (category Category) Exists() bool {
	return category.ID > 0
}

func (Category) TableNmae() string {
	return "category"
}

type ProductCategory struct {
	ProductID  string `gorm:"column:product_id;primaryKey;type:varchar(36)"`
	CategoryID uint   `gorm:"column:category_id;primaryKey"`
}

func (ProductCategory) TableName() string {
	return "product_category"
}

// 属性模板表
type AttributeTemplate struct {
	ID       uint         `gorm:"primaryKey;comment:属性ID"`
	Name     string       `gorm:"size:40;not null;index;comment:属性名称"`
	DataType DataType     `gorm:"type:varchar(20);check:data_type IN ('STRING','NUMBER','ENUM','BOOLEAN','DATE');comment:数据类型"`
	Required bool         `gorm:"default:false;comment:是否必填"`
	Options  JSONStrArray `gorm:"type:json;comment:枚举选项(仅ENUM类型需要)"`
	Unit     string       `gorm:"size:10;comment:计量单位(如cm、kg)"`
}

func (a AttributeTemplate) Exists() bool {
	return a.ID > 0
}

func (AttributeTemplate) TableName() string {
	return "attribute_template"
}

// 分类属性关系表
type CategoryAttribute struct {
	CategoryID  uint `gorm:"primaryKey;comment:分类ID"`
	AttributeID uint `gorm:"primaryKey;comment:属性ID"`
	IsRequired  bool `gorm:"comment:是否必填（覆盖模板规则）"`
	IsInherited bool `gorm:"comment:是否继承到子分类"`
}

func (CategoryAttribute) TableName() string {
	return "category_attribute"
}

type ProductAttribute struct {
	ProductID   string `gorm:"column:product_id;primaryKey"`
	AttributeID uint   `gorm:"column:attribute_id;primaryKey"`
	Value       string `gorm:"column:value;type:varchar(50)"`
}

func (ProductAttribute) TableName() string {
	return "product_attribute"
}
