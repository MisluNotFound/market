package response

import "github.com/mislu/market-api/internal/types/models"

type SearchProductResp struct {
	Products []UserProduct `json:"products"`
	PageResp
}

type SearchHistoryResp struct {
	History []models.SearchHistory `json:"history"`
}