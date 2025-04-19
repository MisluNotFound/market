package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/mislu/market-api/internal/types/request"
)

func GetAssert() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.GetAssertReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		// TODO implement
	}
}
