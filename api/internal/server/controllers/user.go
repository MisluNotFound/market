package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/mislu/market-api/internal/service"
	"github.com/mislu/market-api/internal/types/request"
)

// POST /api/user/register
func CreateUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.CreateUserReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		if err := service.CreateUser(req); err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, "ok")
	}
}

// PUT /api/user/{userID}/basic
func UpdateBasic() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.UpdateBasicReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		if err := service.UpdateBasic(req); err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, "ok")
	}
}

// PUT /api/user/{userID}/avatar
func UploadAvatar() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.UploadAvatarReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		resp, err := service.UploadAvatar(req)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, resp)
	}
}

// PUT /api/user/{userID}/password
func UpdatePassword() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.UpdatePasswordReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		if err := service.UpdatePassword(req); err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, "ok")
	}
}

// POST /api/user/login
func Login() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.LoginReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		resp, err := service.Login(req)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, resp)
	}
}

// GET /api/user/{userID}
func GetUserInfo() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.GetUserInfoReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		resp, err := service.GetUserInfo(req)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, resp)
	}
}

func RefreshAccessToken() func(c *gin.Context) {
	return func(c *gin.Context) {

	}
}

func SelectInterestTags() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.SelectInterestTagsReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		err := service.SelectInterestTags(req)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, "ok")
	}
}
