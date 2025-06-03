package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/mislu/market-api/internal/service"
	"github.com/mislu/market-api/internal/types/request"
)

// POST /api/search/products
func SearchProduct() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.SearchProductReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		req.PageReq.Fill()
		req.UserID, _ = GetContextUserID(c)
		resp, err := service.SearchProduct(req)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, resp)
	}
}

// GET /api/search/{userID}/history
func GetSearchHistory() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.GetSearchHistoryReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		resp, err := service.GetSearchHistory(req)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, resp)
	}
}
