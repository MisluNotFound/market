package request

import (
	"mime/multipart"
	"time"
)

type ProductIDReq struct {
	ProductID string `uri:"productID" binding:"required"`
}

type CreateProductReq struct {
	UserIDReq
	OriginalPrice float64                 `json:"originalPrice" form:"originalPrice"`
	Price         float64                 `json:"price" form:"price" binding:"required,gt=0"`
	Describe      string                  `json:"describe" form:"describe" binding:"required,min=1,max=255"`
	Pics          []*multipart.FileHeader `json:"pics" form:"pics" binding:"required,min=1,max=5"`
	ShipMethod    string                  `json:"shipMethod" form:"shipMethod" binding:"required,oneof=included fixed"`
	ShipPrice     float64                 `json:"shipPrice" form:"shipPrice"`
	CanSelfPickup bool                    `json:"canSelfPickup" form:"canSelfPickup"`
	Condition     string                  `json:"condition" form:"condition" binding:"required,oneof=new excellent good used"`
	UsedTime      string                  `form:"usedTime"`
	// TODO location

	Categories     []uint `form:"categories" binding:"required"`
	AttributesJson string `form:"attributes" binding:"required"`
}

type ProductDocument struct {
	ID         string        `json:"id"`
	Describe   string        `json:"describe"`
	Category   []string      `json:"category"`
	CreatedAt  time.Time     `json:"created_at"`
	Attributes []AttributeES `json:"attributes"`
	Price      float64       `json:"price"`
}

// AttributeES 定义嵌套 attributes 字段
type AttributeES struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type GetProductReq struct {
	UserIDReq
	ProductIDReq
}

type UpdateProductReq struct {
	UserIDReq
	ProductIDReq
	OriginalPrice float64                 `form:"originalPrice" binding:"required,gt=0"`
	Price         float64                 `form:"price" binding:"required,gt=0"`
	Describe      string                  `form:"describe" binding:"required"`
	DeletedPics   []string                `form:"deletedPics"`
	AddedPics     []*multipart.FileHeader `form:"addedPics"`
	ShipMethod    string                  `form:"shipMethod" binding:"required,oneof=included fixed"`
	ShipPrice     int                     `form:"shipPrice" binding:"required,gt=0"`
	CanSelfPickup bool                    `form:"canSelfPickup"`
	// TODO location
	Categories     []uint `form:"categories" binding:"required"`
	AttributesJson string `form:"attributes" binding:"required"`
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

type UpdateProductPriceReq struct {
	UserIDReq
	ProductIDReq
	Price float64 `form:"price" binding:"required,gt=0"`
}

type LikeProductReq struct {
	UserIDReq
	ProductIDReq
}

type DislikeProductReq struct {
	UserIDReq
	ProductIDReq
}

type GetUserLikesReq struct {
	UserIDReq
	PageReq
}
