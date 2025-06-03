package request

import "mime/multipart"

type UserIDReq struct {
	UserID string `uri:"userID"`
}

type CreateUserReq struct {
	Username        string `form:"username" binding:"required"`
	Password        string `form:"password" binding:"required"`
	ConfirmPassword string `form:"confirmPassword" binding:"required,eqfield=Password"`
	Phone           string `form:"phone" binding:"required"`

	// Gender          string `form:"gender" binding:"oneof=male female"`
	// School          string `form:"school" binding:"required"`
}

type UpdateBasicReq struct {
	UserIDReq
	Username string `form:"username"`
	Gender   string `form:"gender" binding:"oneof=male female"`
}

type UploadAvatarReq struct {
	UserIDReq
	Avatar *multipart.FileHeader `form:"avatar" binding:"required"`
}

type UpdatePasswordReq struct {
	UserIDReq
	Password        string `form:"password" binding:"required"`
	NewPassword     string `form:"newPassword" binding:"required"`
	ConfirmPassword string `form:"confirmPassword" binding:"required,eqfield=NewPassword"`
}

type LoginReq struct {
	Phone    string `form:"phone" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type GetUserInfoReq struct {
	UserIDReq
}

type SelectInterestTagsReq struct {
	UserIDReq
	Tags []int `form:"tags" json:"tags"`
}
