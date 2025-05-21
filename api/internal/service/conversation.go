package service

import (
	"errors"
	"fmt"
	"mime"
	"strings"
	"time"

	resourcemanager "github.com/mislu/market-api/internal/core/resource_manager"
	"github.com/mislu/market-api/internal/db"
	"github.com/mislu/market-api/internal/types/exceptions"
	"github.com/mislu/market-api/internal/types/models"
	"github.com/mislu/market-api/internal/types/request"
	"github.com/mislu/market-api/internal/types/response"
	"github.com/mislu/market-api/internal/utils/lib"
	"github.com/mislu/market-api/internal/utils/snowflake"
	"gorm.io/gorm"
)

const (
	text  = "text"
	image = "image"
	video = "video"
	link  = "link"
)

func SaveMessage(raw *request.Message) error {
	id, err := snowflake.IdGenerator.NextID()
	if err != nil {
		return err
	}

	raw.ID = id
	conversationID := getConversationID(raw.From, raw.To)

	message := &models.Message{
		ID:             id,
		FromUserID:     raw.From,
		ToUserID:       raw.To,
		ConversationID: conversationID,
		MediaType:      raw.MediaType,
		Timestamp:      time.Now(),
	}

	switch raw.MediaType {
	case text:
		message.Content = raw.Content
	case image, video:
		// 媒体文件通过http上传，这里默认为url
		message.Content = raw.Content
	case link:
		// 商品链接
		// TODO implement
	default:
		return nil
	}

	return db.WithTransaction(func(tx *gorm.DB) error {
		// 消息落盘
		err := db.Create(message, tx)
		if err != nil {
			return err
		}

		conversation, err := db.GetOne[models.Conversation](
			db.Equal("to_user_id", message.ToUserID),
			db.Equal("from_user_id", message.FromUserID),
		)

		// TODO 更新接收方的对话状态

		if err != nil {
			return err
		}

		if !conversation.MarkDeleted {
			return nil
		}

		conversation.MarkDeleted = false
		conversation.LastMessageTime = message.Timestamp
		conversation.LastMessageContent = message.Content
		return db.Update(&conversation, tx)
	})
}

func getConversationID(from, to string) string {
	if from > to {
		return fmt.Sprintf("%s:%s", to, from)
	}
	return fmt.Sprintf("%s:%s", from, to)
}

func UploadMediaFile(req *request.UploadMediaFileReq) (response.UploadMediaFileResp, exceptions.APIError) {
	var resp response.UploadMediaFileResp

	mediaType := GetCategoryByExt(req.File.Filename)
	switch mediaType {
	case image:
		if req.File.Size > picMaxSize {
			return resp, exceptions.BadRequestError(errors.New("file size exceeds limit"), exceptions.ImageFileSizeExceedError)
		}
	case video:

	default:
		return resp, exceptions.BadRequestError(errors.New("unsupported file type"), exceptions.UnsupportedFileTypeError)
	}
	key := resourcemanager.GenerateObjectKey(req.File.Filename)
	resp.Url = lib.GetResourceURL(int(resourcemanager.ConversationBucket), req.ConversationID, key)
	return resp, nil
}

var extToMIME = map[string]string{
	// 图片类型
	"jpg":  "image/jpeg",
	"jpeg": "image/jpeg",
	"png":  "image/png",
	"gif":  "image/gif",
	"webp": "image/webp",
	"bmp":  "image/bmp",
	"tiff": "image/tiff",

	// 视频类型
	"mp4":  "video/mp4",
	"mov":  "video/quicktime",
	"avi":  "video/x-msvideo",
	"mkv":  "video/x-matroska",
	"webm": "video/webm",
	"flv":  "video/x-flv",
	"wmv":  "video/x-ms-wmv",
}

// MIME类型 -> 分类映射表
var mimeToCategory = map[string]string{
	// 图片
	"image/jpeg":    "image",
	"image/png":     "image",
	"image/gif":     "image",
	"image/webp":    "image",
	"image/bmp":     "image",
	"image/tiff":    "image",
	"image/svg+xml": "image",

	// 视频
	"video/mp4":             "video",
	"video/quicktime":       "video",
	"video/x-msvideo":       "video",
	"video/x-matroska":      "video",
	"video/webm":            "video",
	"video/x-flv":           "video",
	"video/x-ms-wmv":        "video",
	"application/x-mpegURL": "video", // HLS流
}

// 通过扩展名获取文件分类
func GetCategoryByExt(ext string) string {
	// 统一处理扩展名格式
	ext = strings.ToLower(strings.TrimPrefix(ext, "."))

	// 获取MIME类型
	mimeType, ok := extToMIME[ext]
	if !ok {
		// 尝试通过系统MIME类型检测
		if detected := mime.TypeByExtension("." + ext); detected != "" {
			mimeType = detected
		} else {
			return ""
		}
	}

	// 返回分类
	return mimeToCategory[mimeType]
}

func GetCategoryByMIME(mimeType string) string {
	baseMIME, _, _ := mime.ParseMediaType(mimeType)
	return mimeToCategory[baseMIME]
}

func CreateConversation(req *request.CreateConversationReq) exceptions.APIError {
	conversation, err := db.GetOne[models.Conversation](
		db.Equal("from_user_id", req.FromUserID),
		db.Equal("to_user_id", req.ToUserID),
	)
	if err != nil {
		return exceptions.InternalServerError(err)
	}

	if conversation.FromUserID == req.FromUserID {
		conversation.MarkDeleted = false
		conversation.CurrentProductID = req.ProductID
		err := db.Update(&conversation)
		if err != nil {
			return exceptions.InternalServerError(err)
		}
		return nil
	}

	conversations := []*models.Conversation{
		&models.Conversation{
			FromUserID:       req.FromUserID,
			ToUserID:         req.ToUserID,
			CurrentProductID: req.ProductID,
		},
		&models.Conversation{
			FromUserID:       req.ToUserID,
			ToUserID:         req.FromUserID,
			CurrentProductID: req.ProductID,
		},
	}

	err = db.Create(conversations)
	if err != nil {
		return exceptions.InternalServerError(err)
	}

	return nil
}

func RecordLastReadMessage(fromUserID, toUserID string, lastMessageID string) error {
	conversation, err := db.GetOne[models.Conversation](
		// 获取接收方的conversation
		db.Equal("from_user_id", toUserID),
		db.Equal("to_user_id", fromUserID),
	)

	if err != nil {
		return err
	}

	message, err := db.GetOne[models.Message](
		db.Fields("id"),
		db.LessThan("id", lastMessageID),
		db.Page(1, 1),
	)

	if err != nil {
		return err
	}

	conversation.LastReadMessageID = message.ID
	return db.Update(&conversation)
}

func GetConversationList(req *request.GetConversationListReq) (response.GetConversationListResp, exceptions.APIError) {
	conversations, err := db.GetAll[models.Conversation](
		db.Equal("from_user_id", req.UserID),
		db.Equal("mark_deleted", false),
	)

	if err != nil {
		return nil, exceptions.InternalServerError(err)
	}

	resp := make([]response.ConversationWithUnReadCount, 0, len(conversations))

	for _, conversation := range conversations {
		user, err := db.GetOne[models.User](
			db.Fields("id", "username", "avatar"),
			db.Equal("id", conversation.ToUserID),
		)

		if err != nil {
			return nil, exceptions.InternalServerError(err)
		}

		unreadCount, err := db.GetCount[models.Message](
			db.GreaterThan("id", conversation.LastReadMessageID),
			db.Equal("conversation_id", getConversationID(conversation.FromUserID, conversation.ToUserID)),
		)
		if err != nil {
			return nil, exceptions.InternalServerError(err)
		}

		resp = append(resp, response.ConversationWithUnReadCount{
			User:        user,
			UnreadCount: int(unreadCount),
			Conversation: conversation,
		})
	}

	return resp, nil
}

func GetMessages(req *request.GetMessagesReq) (response.GetMessagesResp, exceptions.APIError) {
	var resp response.GetMessagesResp
	resp.Page = req.Page
	resp.Size = req.Size

	messages, err := db.GetAll[models.Message](
		db.Equal("conversation_id", getConversationID(req.FromUserID, req.ToUserID)),
		db.Page(req.Page, req.Size),
		db.OrderBy("id", true),
	)

	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	if len(messages) == 0 {
		return resp, nil
	}

	conversation, err := db.GetOne[models.Conversation](
		db.Equal("from_user_id", req.FromUserID),
		db.Equal("to_user_id", req.ToUserID),
	)
	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	conversation.LastReadMessageID = messages[0].ID
	err = db.Update(conversation)
	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	resp.Messages = messages
	return resp, nil
}
