package response

import "github.com/mislu/market-api/internal/types/models"

type GetAllOrderStatusResp struct {
	Bought      int `json:"bought"`
	Sold        int `json:"sold"`
	BeEvaluated int `json:"beEvaluated"`
}

type GetOrderListResp struct {
	PageResp
	Orders []UserOrder `json:"orders"`
}

type UserOrder struct {
	models.User    `json:"user"`
	models.Order   `json:"order"`
	models.Product `json:"product"`
}

type GetOrderResp struct {
	UserOrder
}

type PayOrderResp struct {
	PayURL string `json:"payURL"`
}

type GetUnCommentOrderResp struct {
	Orders []UserOrder `json:"orders"`
}

type PurchaseProductResp struct {
	OrderID string `json:"orderID"`
}
