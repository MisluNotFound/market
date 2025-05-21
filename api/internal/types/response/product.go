package response

import "github.com/mislu/market-api/internal/types/models"

type CreateProductResp struct {
	Failures []UploadFileFailure `json:"failures"`
}

type UploadFileFailure struct {
	FileName string `json:"fileName"`
	Error    string `json:"error"`
}

type GetProductResp struct {
	UserProduct
	// TODO comment
}

type UserProduct struct {
	models.User    `json:"user"`
	models.Product `json:"product"`
	Categories     []uint          `json:"categories"`
	Attributes     map[uint]string `json:"attributes"`
}

type GetUserProductsResp struct {
	Products []UserProduct `json:"products"`
	PageResp
}

type PageResp struct {
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Size  int   `json:"size"`

	// TODO use it
	HasMore bool `json:"hasMore"`
}

type GetProductListResp struct {
	Products []UserProduct `json:"products"`
	PageResp
}

type GetAllCategoryResp []*WrappedCategory

type WrappedCategory struct {
	models.Category
	Children   []*WrappedCategory         `json:"children"`
	Attributes []models.AttributeTemplate `json:"attributes"`
}
