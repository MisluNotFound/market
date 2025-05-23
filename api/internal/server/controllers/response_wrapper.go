package controllers

import (
	exceptions "github.com/mislu/market-api/internal/types/exceptions"
	"github.com/mislu/market-api/pkg/entities"

	"github.com/gin-gonic/gin"
)

const (
	abortError  = "__abortError"
	payload     = "__payload"
	_contentType = "__contentType"
)

const (
	ResponseTypeJSON   = "json"
	ResponseTypeFile   = "file"
	ResponseTypeStream = "stream"
)

func AbortWithError(ctx *gin.Context, err exceptions.APIError) {
	ctx.Set(abortError, err)
	ctx.AbortWithError(err.ToResponse().Code, err)
}

func GetAbortError(ctx *gin.Context) exceptions.APIError {
	if abortErr, exists := ctx.Get(abortError); exists {
		if err, ok := abortErr.(exceptions.APIError); ok {
			return err
		}
	}

	return nil
}

func Success(ctx *gin.Context, contentType string, data any) {
	ctx.Set(_contentType, contentType)
	ctx.Set(payload, entities.NewSuccessResponse(data))
}

func GetPayLoad(ctx *gin.Context) *entities.Response {
	if apiError, exists := ctx.Get(payload); exists {
		if resp, ok := apiError.(*entities.Response); ok {
			return resp
		}
	}
	return nil
}

func GetContentType(ctx *gin.Context) string {
	if contentType, exists := ctx.Get(_contentType); exists {
		if ct, ok := contentType.(string); ok {
			return ct
		}
	}
	return ResponseTypeJSON
}
