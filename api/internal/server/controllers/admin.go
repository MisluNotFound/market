package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/mislu/market-api/internal/service"
	"github.com/mislu/market-api/internal/types/request"
)

func CreateCategory() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.CreateCategoryReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		if err := service.CreateCategory(req); err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, "ok")
	}
}

func UpdateCategory() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.UpdateCategoryReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		if err := service.UpdateCategory(req); err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, "ok")
	}
}

func DeleteCategory() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.DeleteCategoryReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		if err := service.DeleteCategory(req); err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, "ok")
	}
}

func CreateInterestTag() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.CreateInterestTagReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		if err := service.CreateInterestTag(req); err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, "ok")
	}
}

func UpdateInterestTag() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.UpdateInterestTagReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		if err := service.UpdateInterestTag(req); err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, "ok")
	}
}

func DeleteInterestTag() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.DeleteInterestTagReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		if err := service.DeleteInterestTag(req); err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, "ok")
	}
}

func CreateAttribute() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.CreateAttributeReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		if err := service.CreateAttribute(req); err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, "ok")
	}
}

func DeleteAttribute() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.DeleteAttributeReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		if err := service.DeleteAttribute(req); err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, "ok")
	}
}

func UpdateAttribute() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.UpdateAttributeReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		if err := service.UpdateAttribute(req); err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, "ok")
	}
}
