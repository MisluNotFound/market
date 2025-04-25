package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/mislu/market-api/internal/service"
	"github.com/mislu/market-api/internal/types/request"
	"github.com/mislu/market-api/internal/utils/app"
)

// api/assert/{type}/{owner}/{key}
func GetAssert() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.GetAssertReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		switch app.GetConfig().OSS.Type {
		case "local_storage":
			path, err := service.GetAssert(req)
			if err != nil {
				AbortWithError(c, err)
				return
			}

			Success(c, ResponseTypeFile, path)
		}
	}
}
