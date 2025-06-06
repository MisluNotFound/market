package response

import "github.com/mislu/market-api/internal/types/models"

type UploadAvatarResp struct {
	Avatar string `json:"avatar"`
}

type LoginResp struct {
	RefreshToken   string `json:"refreshToken"`
	AccessToken    string `json:"accessToken"`
	UserID         string `json:"userID"`
	NeedSelectTags bool   `json:"needSelectTags"`
}

type GetUserInfoResp struct {
	models.User
	models.Credit
	// TODO 添加商品信息
}
