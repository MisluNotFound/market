package response

import "github.com/mislu/market-api/internal/types/models"

type UploadMediaFileResp struct {
	Url string `json:"url"`
}

type GetConversationListResp []ConversationWithUnReadCount

type ConversationWithUnReadCount struct {
	models.Conversation
	UnreadCount int `json:"unreadCount"`
	models.User
}


type GetMessagesResp struct {
	Messages []models.Message
	PageResp
}