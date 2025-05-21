package dbinitialize

import (
	"encoding/json"
	"os"

	"github.com/mislu/market-api/internal/types/models"
	"gorm.io/gorm"
)

// 初始化分类（电子产品->手机->智能手机）
func InitCategories(db *gorm.DB) error {
	type CategoryInit struct {
		Level      int8   `json:"level"`
		TypeName   string `json:"typeName"`
		ParentName string `json:"parentName"`
		IsLeaf     bool   `json:"isLeaf"`
	}

	jsonData, err := os.ReadFile("D:\\repository\\m-market\\api\\internal\\db\\db_initialize\\categories.json")
	if err != nil {
		panic(err)
	}

	var items []CategoryInit
	if err := json.Unmarshal(jsonData, &items); err != nil {
		panic(err)
	}

	// 初始化所有分类
	cs := make([]*models.Category, 0, len(items))
	for _, item := range items {
		c := &models.Category{
			TypeName: item.TypeName,
			Level:    item.Level,
			IsLeaf:   item.IsLeaf,
		}

		cs = append(cs, c)
		if err := db.Where("type_name = ?", c.TypeName).FirstOrCreate(c).Error; err != nil {
			return err
		}
	}

	for _, c := range items {
		var parentID int
		if err := db.Model(&models.Category{}).Select("id").Where("type_name = ?", c.ParentName).Find(&parentID).Error; err != nil {
			return err
		}

		if err := db.Model(&models.Category{}).Where("type_name = ?", c.TypeName).Update("parent_id", parentID).Error; err != nil {
			return err
		}
	}

	return nil
}

// 初始化属性模板
func InitAttributes(db *gorm.DB) error {
	type InitAttribute struct {
		Name     string   `json:"name"`
		DataType string   `json:"dataType"`
		Required bool     `json:"required"`
		Options  []string `json:"options"`
		Unit     string   `json:"unit"`
	}

	jsonData, err := os.ReadFile("D:\\repository\\m-market\\api\\internal\\db\\db_initialize\\attrubutes.json")
	if err != nil {
		return err
	}

	var items []InitAttribute
	if err := json.Unmarshal(jsonData, &items); err != nil {
		return err
	}

	for _, item := range items {
		a := models.AttributeTemplate{
			Name:     item.Name,
			DataType: models.DataType(item.DataType),
			Required: item.Required,
			Options:  item.Options,
			Unit:     item.Unit,
		}

		err := db.Where("name = ?", a.Name).FirstOrCreate(&a).Error
		if err != nil {
			return err
		}
	}

	return nil
}

// 建立分类与属性关系
func LinkCategoryAttributes(db *gorm.DB) error {
	type CateAttrInit struct {
		CategoryName string   `json:"categoryName"`
		Attributes   []string `json:"attributes"`
	}

	jsonData, err := os.ReadFile("D:\\repository\\m-market\\api\\internal\\db\\db_initialize\\category_attr.json")
	if err != nil {
		panic(err)
	}

	var cateAttrs []CateAttrInit
	if err := json.Unmarshal(jsonData, &cateAttrs); err != nil {
		panic(err)
	}

	for _, attrs := range cateAttrs {
		var attrIDs []uint
		if err := db.Model(&models.AttributeTemplate{}).Select("id").Where("name in ?", attrs.Attributes).Find(&attrIDs).Error; err != nil {
			return err
		}

		var cid int
		if err := db.Model(&models.Category{}).Select("id").Where("type_name = ?", attrs.CategoryName).Find(&cid).Error; err != nil {
			return err
		}

		links := make([]models.CategoryAttribute, 0, len(attrIDs))
		for _, attr := range attrIDs {
			links = append(links, models.CategoryAttribute{
				CategoryID:  uint(cid),
				AttributeID: attr,
			})
		}

		if err := db.Create(&links).Error; err != nil {
			return err
		}
	}

	return nil
}

func RunInit(db *gorm.DB) error {
	if err := InitAttributes(db); err != nil {
		return err
	}

	if err := InitCategories(db); err != nil {
		return err
	}

	return LinkCategoryAttributes(db)
}
