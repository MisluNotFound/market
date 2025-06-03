package response

import (
	"github.com/mislu/market-api/internal/types/models"
)

type OrderCommentDetail struct {
	models.OrderComment
	Username string               `json:"username"`
	Avatar   string               `json:"avatar"`
	Replies  []OrderCommentDetail `json:"replies,omitempty"`
}

type GetOrderCommentsResp struct {
	Comments []OrderCommentDetail `json:"comments"`
	Total    int64                `json:"total"`
	PageResp
}

type GetSellerCommentsResp struct {
	Comments []OrderCommentDetail `json:"comments"`
	Total    int64                `json:"total"`
	PageResp
	GoodRate float64 `json:"goodRate"`
}
