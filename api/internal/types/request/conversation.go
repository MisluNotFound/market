package request

import "mime/multipart"

type ConversationIDReq struct {
	ConversationID string `form:"conversationID" uri:"conversationID"`
}

type Message struct {
	TempID    string `json:"tempID"` // 发送方生成的临时唯一id，用于server向发送方发送ack消息
	ID        string `json:"id"`
	From      string `json:"from"`
	To        string `json:"to"`
	Content   string `json:"content"`
	MediaType string `json:"mediaType"` // text/image/link/video
	Type      uint   `json:"type"`      // message/ack
}

type UploadMediaFileReq struct {
	ConversationIDReq
	File *multipart.FileHeader `form:"file"`
}

type CreateConversationReq struct {
	FromUserID string `form:"fromUserID" binding:"required"`
	ToUserID   string `form:"toUserID" binding:"required"`
	ProductID  string `form:"productID" binding:"required"`
}

type GetConversationListReq struct {
	UserIDReq
}

type GetMessagesReq struct {
	FromUserID string `form:"fromUserID" binding:"required"`
	ToUserID   string `form:"toUserID" binding:"required"`
	PageReq
}
