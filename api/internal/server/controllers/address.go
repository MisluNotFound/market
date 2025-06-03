package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/mislu/market-api/internal/service"
	"github.com/mislu/market-api/internal/types/request"
)

func CreateAddress() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.CreateAddressReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		err := service.CreateAddress(req)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, "ok")
	}
}

func UpdateAddress() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.UpdateAddressReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		err := service.UpdateAddress(req)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, "ok")
	}
}

func DeleteAddress() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.DeleteAddressReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		err := service.DeleteAddress(req)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, "ok")
	}
}

func GetAddress() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.GetAddressReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		resp, err := service.GetAddress(req)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, resp)
	}
}

func SetDefaultAddress() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.SetDefaultAddressReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		err := service.SetDefaultAddress(req)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, "ok")
	}
}
