package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/mislu/market-api/internal/types/exceptions"
)

func BindRequest[T any](ctx *gin.Context, ptr T) exceptions.APIError {
	ctx.ShouldBindUri(ptr)

	switch ctx.ContentType() {
	case gin.MIMEJSON:
		if err := ctx.ShouldBindWith(ptr, binding.JSON); err != nil {
			return exceptions.BadRequestError(err, exceptions.ParameterBindingError)
		}
	case gin.MIMEMultipartPOSTForm:
		if err := ctx.ShouldBindWith(ptr, binding.FormMultipart); err != nil {
			return exceptions.BadRequestError(err, exceptions.ParameterBindingError)
		}
	case gin.MIMEPOSTForm:
		if err := ctx.ShouldBindWith(ptr, binding.FormPost); err != nil {
			return exceptions.BadRequestError(err, exceptions.ParameterBindingError)
		}
	default:
		if err := ctx.ShouldBindQuery(ptr); err != nil {
			return exceptions.BadRequestError(err, exceptions.ParameterBindingError)
		}
	}

	return nil
}
