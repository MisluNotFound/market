package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/mislu/market-api/internal/service"
	"github.com/mislu/market-api/internal/types/request"
)

// post /api/conversation/userID/
func UploadMediaFile() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.UploadMediaFileReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		resp, err := service.UploadMediaFile(req)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, resp)
	}
}

// post /api/conversation/create
func CreateConversation() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.CreateConversationReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		err := service.CreateConversation(req)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, "ok")
	}
}

// get /api/conversation/{userID}/list
func GetConversationList() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.GetConversationListReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		resp, err := service.GetConversationList(req)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, resp)
	}
}

// get /api/conversation/messages
func GetMessages() func(c *gin.Context) {
	return func(c *gin.Context) {
		req := &request.GetMessagesReq{}
		if err := BindRequest(c, req); err != nil {
			AbortWithError(c, err)
			return
		}

		req.Fill()

		resp, err := service.GetMessages(req)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		Success(c, ResponseTypeJSON, resp)
	}
}
