package service

import (
	"github.com/mislu/market-api/internal/db"
	"github.com/mislu/market-api/internal/types/exceptions"
	"github.com/mislu/market-api/internal/types/models"
	"github.com/mislu/market-api/internal/types/request"
	"gorm.io/gorm"
)

func CreateCategory(req *request.CreateCategoryReq) exceptions.APIError {
	category, err := db.GetOne[models.Category](
		db.Equal("type_name", req.CategoryName),
		db.Equal("parent_id", req.ParentID),
	)

	if err != nil {
		return exceptions.InternalServerError(err)
	}

	if category.Exists() {
		return nil
	}

	newCategory := models.Category{
		TypeName: req.CategoryName,
		Level:    int8(req.Level),
	}

	if req.ParentID != 0 {
		newCategory.ParentID = req.ParentID
	}

	if req.Level == 3 {
		newCategory.IsLeaf = true
	}

	err = db.WithTransaction(func(tx *gorm.DB) error {
		if err := db.Create(&newCategory, tx); err != nil {
			return err
		}

		if !newCategory.IsLeaf {
			return nil
		}

		return nil
	})

	if err != nil {
		return exceptions.InternalServerError(err)
	}

	return nil
}

func UpdateCategory(req *request.UpdateCategoryReq) exceptions.APIError {
	category, err := db.GetOne[models.Category](
		db.Equal("id", req.CategoryID),
	)

	if err != nil {
		return exceptions.InternalServerError(err)
	}

	if !category.Exists() {
		return nil
	}

	category.TypeName = req.CategoryName

	err = db.Update(&category)

	if err != nil {
		return exceptions.InternalServerError(err)
	}

	return nil
}

func DeleteCategory(req *request.DeleteCategoryReq) exceptions.APIError {
	category, err := db.GetOne[models.Category](
		db.Equal("id", req.CategoryID),
	)

	if err != nil {
		return exceptions.InternalServerError(err)
	}

	if !category.Exists() {
		return nil
	}

	err = db.Delete(&category)
	if err != nil {
		return exceptions.InternalServerError(err)
	}

	return nil
}

func CreateAttribute(req *request.CreateAttributeReq) exceptions.APIError {
	required := false
	if req.Required != nil && *req.Required {
		required = true
	}

	attribute := models.AttributeTemplate{
		Name:     req.Name,
		DataType: models.DataType(req.DataType),
		Required: required,
		Options:  req.Options,
		Unit:     req.Unit,
	}

	err := db.WithTransaction(func(tx *gorm.DB) error {
		if err := db.Create(&attribute, tx); err != nil {
			return err
		}

		categoryAttribute := models.CategoryAttribute{
			CategoryID:  uint(req.CategoryID),
			AttributeID: attribute.ID,
		}
		if err := db.Create(&categoryAttribute, tx); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return exceptions.InternalServerError(err)
	}

	return nil
}

func DeleteAttribute(req *request.DeleteAttributeReq) exceptions.APIError {
	attribute, err := db.GetOne[models.AttributeTemplate](
		db.Equal("id", req.AttributeID),
	)

	if err != nil {
		return exceptions.InternalServerError(err)
	}

	if !attribute.Exists() {
		return nil
	}

	categoryAttributes, err := db.GetAll[models.CategoryAttribute](
		db.Equal("attribute_id", req.AttributeID),
	)

	if err != nil {
		return exceptions.InternalServerError(err)
	}

	err = db.WithTransaction(func(tx *gorm.DB) error {
		if err := db.Delete(&attribute, tx); err != nil {
			return err
		}

		if err := db.Delete(&categoryAttributes, tx); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return exceptions.InternalServerError(err)
	}

	return nil
}

func UpdateAttribute(req *request.UpdateAttributeReq) exceptions.APIError {
	attribute, err := db.GetOne[models.AttributeTemplate](
		db.Equal("id", req.AttributeID),
	)

	if err != nil {
		return exceptions.InternalServerError(err)
	}

	if !attribute.Exists() {
		return nil
	}

	required := false
	if req.Required != nil && *req.Required {
		required = true
	}

	attribute.Name = req.Name
	attribute.DataType = models.DataType(req.DataType)
	attribute.Required = required
	attribute.Options = req.Options
	attribute.Unit = req.Unit

	err = db.Update(&attribute)
	if err != nil {
		return exceptions.InternalServerError(err)
	}

	return nil
}

func CreateInterestTag(req *request.CreateInterestTagReq) exceptions.APIError {
	tag := models.InterestTag{
		TagName:    req.TagName,
		CategoryID: req.CategoryID,
	}

	err := db.Create(&tag)
	if err != nil {
		return exceptions.InternalServerError(err)
	}

	return nil
}

func UpdateInterestTag(req *request.UpdateInterestTagReq) exceptions.APIError {
	tag, err := db.GetOne[models.InterestTag](
		db.Equal("id", req.TagID),
	)

	if err != nil {
		return exceptions.InternalServerError(err)
	}

	if !tag.Exists() {
		return nil
	}

	tag.TagName = req.TagName
	tag.CategoryID = req.CategoryID

	err = db.Update(&tag)
	if err != nil {
		return exceptions.InternalServerError(err)
	}

	return nil
}

func DeleteInterestTag(req *request.DeleteInterestTagReq) exceptions.APIError {
	tag, err := db.GetOne[models.InterestTag](
		db.Equal("id", req.TagID),
	)

	if err != nil {
		return exceptions.InternalServerError(err)
	}

	if !tag.Exists() {
		return nil
	}

	err = db.Delete(&tag)
	if err != nil {
		return exceptions.InternalServerError(err)
	}

	return nil
}
