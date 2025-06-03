package request

type CreateCategoryReq struct {
	CategoryName string         `form:"categoryName" json:"categoryName" binding:"required,max=50"`
	ParentID     uint           `form:"parentID" json:"parentID" binding:"omitempty"`
	Level        int            `form:"level" json:"level" binding:"omitempty,oneof=1 2 3"`
	Attributes   []AttributeReq `form:"attributes" json:"attributes" binding:"omitempty"`
}

type AttributeReq struct {
	Name     string   `form:"name" json:"name" binding:"required,max=50"`
	DataType string   `form:"dataType" json:"dataType" binding:"required,oneof=STRING NUMBER ENUM BOOLEAN DATE"`
	Required *bool    `form:"required" json:"required" binding:"required"`
	Options  []string `form:"options" json:"options" binding:"omitempty"`
	Unit     string   `form:"unit" json:"unit" binding:"omitempty,max=10"`
}

type UpdateCategoryReq struct {
	CategoryID   uint   `form:"categoryID" json:"categoryID" binding:"required"`
	CategoryName string `form:"categoryName" json:"categoryName" binding:"required,max=50"`
}

type DeleteCategoryReq struct {
	CategoryID uint `form:"categoryID" json:"categoryID" binding:"required"`
}

type CreateAttributeReq struct {
	CategoryID int `form:"categoryID" json:"categoryID" binding:"required"`
	AttributeReq
}

type DeleteAttributeReq struct {
	AttributeID int `form:"attributeID" json:"attributeID" binding:"required"`
}

type UpdateAttributeReq struct {
	AttributeID int `form:"attributeID" json:"attributeID" binding:"required"`
	AttributeReq
}

type CreateInterestTagReq struct {
	TagName    string `form:"tagName" json:"tagName" binding:"required,max=50"`
	CategoryID int    `form:"categoryID" json:"categoryID" binding:"required"`
}

type UpdateInterestTagReq struct {
	TagID      int    `form:"tagID" json:"tagID" binding:"required"`
	TagName    string `form:"tagName" json:"tagName" binding:"required,max=50"`
	CategoryID int    `form:"categoryID" json:"categoryID" binding:"required"`
}

type DeleteInterestTagReq struct {
	TagID int `form:"tagID" json:"tagID" binding:"required"`
}
