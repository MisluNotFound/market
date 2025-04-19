package controllers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/mislu/market-api/internal/types/exceptions"
	"github.com/mislu/market-api/internal/types/request"
)

func MockGet() func(c *gin.Context) {
	return func(c *gin.Context) {
		Success(c, "ok")
	}
}

func MockPost() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.MockRequest{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, exceptions.BadRequestError(err, exceptions.ParameterBindingError))
			return
		}
		
		Success(c, "ok")
	}
}

func MockError() func(c *gin.Context) {
	return func(c *gin.Context) {
		AbortWithError(c, exceptions.BadRequestError(errors.New("mock error"), exceptions.MockError))
	}
}
