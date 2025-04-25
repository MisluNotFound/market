package request

import "mime/multipart"

type ProductIDReq struct {
	ProductID string `uri:"productID" binding:"required"`
}

type CreateProductReq struct {
	UserIDReq
	OriginalPrice float64                 `form:"originalPrice"`
	Price         float64                 `form:"price" binding:"required"`
	Describe      string                  `form:"describe" binding:"required"`
	Pics          []*multipart.FileHeader `form:"pics" binding:"required,min=1,max=5"`
	ShipMethod    string                  `form:"shipMethod" binding:"required,oneof=included fixed"`
	ShipPrice     float64                 `form:"shipPrice"`
	CanSelfPickup bool                    `form:"canSelfPickup"`
	// TODO location

	// TODO 分类
}

type GetProductReq struct {
	UserIDReq
	ProductIDReq
}

type UpdateProductReq struct {
	UserIDReq
	ProductIDReq
	OriginalPrice float64                 `form:"originalPrice"`
	Price         float64                 `form:"price" binding:"required"`
	Describe      string                  `form:"describe" binding:"required"`
	DeletedPics   []string                `form:"deletedPics"`
	AddedPics     []*multipart.FileHeader `form:"addedPics"`
	ShipMethod    string                  `form:"shipMethod" binding:"required,oneof=included fixed"`
	ShipPrice     int                     `form:"shipPrice" binding:"required"`
	CanSelfPickup bool                    `form:"canSelfPickup"`

	// TODO location

	// TODO 分类
}

type UpdateProductStatusReq struct {
	UserIDReq
	ProductIDReq
}

type GetUserProductsReq struct {
	UserIDReq
	PageReq
}

type PageReq struct {
	Page int `form:"page" json:"page"`
	Size int `form:"size" json:"size"`
}

func (p *PageReq) Fill() {
	if p.Page == 0 {
		p.Page = 1
	}

	if p.Size == 0 {
		p.Size = 10
	}
}

type GetProductListReq struct {
	PageReq
}

type SearchProductReq struct {
	PageReq
	Keyword string `form:"keyword" binding:"required"`
}
