package service

import (
	"errors"
	"regexp"

	resourcemanager "github.com/mislu/market-api/internal/core/resource_manager"
	"github.com/mislu/market-api/internal/db"
	"github.com/mislu/market-api/internal/types/exceptions"
	"github.com/mislu/market-api/internal/types/models"
	"github.com/mislu/market-api/internal/types/request"
	"github.com/mislu/market-api/internal/types/response"
	"github.com/mislu/market-api/internal/utils/lib"
	"gorm.io/gorm"
)

var mobileRegex = regexp.MustCompile(`^1[3-9][0-9]{9}$`)
var (
	errUserNotFound = errors.New("user does not found")
)

const (
	picMaxSize   = 10 << 20
	videoMaxSize = 10 << 20
)

func CreateUser(req *request.CreateUserReq) exceptions.APIError {
	if ok := mobileRegex.MatchString(req.Phone); !ok {
		return exceptions.BadRequestError(errors.New("invalid phone number"), exceptions.ParameterBindingError)
	}

	user, err := db.GetOne[*models.User](
		db.Equal("phone", req.Phone),
	)

	if user.Exists() {
		return exceptions.BadRequestError(errors.New("phone has been bound"), exceptions.PhoneBoundError)
	}

	if err != nil {
		return exceptions.InternalServerError(err)
	}

	salt, err := lib.GenerateSalt()
	if err != nil {
		return exceptions.InternalServerError(err)
	}

	hashedPassword, err := lib.EncryptPassword(req.Password, salt)
	if err != nil {
		return exceptions.InternalServerError(err)
	}

	// TODO verify phone number
	user = &models.User{
		Username: req.Username,
		Password: hashedPassword,
		Phone:    req.Phone,
		Salt:     salt,
	}

	if err := db.Create(user); err != nil {
		return exceptions.InternalServerError(err)
	}

	return nil
}

func UpdateBasic(req *request.UpdateBasicReq) exceptions.APIError {
	user, err := db.GetOne[*models.User](
		db.Equal("id", req.UserID),
	)

	if err != nil {
		return exceptions.InternalServerError(err)
	}

	if !user.Exists() {
		return exceptions.BadRequestError(errUserNotFound, exceptions.UserNotExistsError)
	}

	if len(req.Gender) > 0 {
		user.Gender = req.Gender
	}

	if len(req.Username) > 0 {
		user.Username = req.Username
	}

	if err := db.Update(user); err != nil {
		return exceptions.InternalServerError(err)
	}

	return nil
}

func UploadAvatar(req *request.UploadAvatarReq) (response.UploadAvatarResp, exceptions.APIError) {
	var resp response.UploadAvatarResp

	if req.Avatar.Size > picMaxSize {
		return resp, exceptions.BadRequestError(errors.New("avatar size exceeds limit"), exceptions.ParameterBindingError)
	}

	user, err := db.GetOne[*models.User](
		db.Equal("id", req.UserID),
	)

	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	if !user.Exists() {
		return resp, exceptions.BadRequestError(errUserNotFound, exceptions.UserNotExistsError)
	}

	avatarFile, err := req.Avatar.Open()
	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	data := make([]byte, req.Avatar.Size)
	_, err = avatarFile.Read(data)
	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	key := resourcemanager.GenerateObjectKey(req.Avatar.Filename)
	path := resourcemanager.GetObjectPath(resourcemanager.UserBucket, req.UserID, key)

	err = resourcemanager.UploadFile(resourcemanager.UserBucket, path, data)
	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	oldKey := lib.SplitResourceURL(user.Avatar)

	user.Avatar = lib.GetResourceURL(int(resourcemanager.UserBucket), req.UserID, key)

	// 忽略删除错误，由定时清理任务清除多余的文件
	resourcemanager.DeleteFile(resourcemanager.UserBucket, resourcemanager.GetObjectPath(resourcemanager.UserBucket, req.UserID, oldKey))
	err = db.Update(user)
	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	resp.Avatar = user.Avatar
	return resp, nil
}

func UpdatePassword(req *request.UpdatePasswordReq) exceptions.APIError {
	user, err := db.GetOne[*models.User](
		db.Equal("id", req.UserID),
	)

	if err != nil {
		return exceptions.InternalServerError(err)
	}

	if !user.Exists() {
		return exceptions.BadRequestError(errUserNotFound, exceptions.UserNotExistsError)
	}

	password, err := lib.EncryptPassword(req.Password, user.Salt)
	if err != nil {
		return exceptions.InternalServerError(err)
	}

	if password != user.Password {
		return exceptions.BadRequestError(errors.New("incorrect password"), exceptions.IncorrectPasswordError)
	}

	salt, err := lib.GenerateSalt()
	if err != nil {
		return exceptions.InternalServerError(err)
	}

	hashedPassword, err := lib.EncryptPassword(req.NewPassword, salt)
	if err != nil {
		return exceptions.InternalServerError(err)
	}

	user.Password = hashedPassword
	user.Salt = salt

	err = db.Update(user)
	if err != nil {
		return exceptions.InternalServerError(err)
	}

	return nil
}

func Login(req *request.LoginReq) (response.LoginResp, exceptions.APIError) {
	var resp response.LoginResp

	user, err := db.GetOne[*models.User](
		db.Equal("phone", req.Phone),
	)

	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	hashedPassword, err := lib.EncryptPassword(req.Password, user.Salt)
	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	if user.Password != hashedPassword {
		return resp, exceptions.BadRequestError(errors.New("incorrect password"), exceptions.IncorrectPasswordError)
	}

	accessToken, err := lib.GenerateAccessToken(user.ID)
	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	refreshToken, err := lib.GenerateRefreshToken(user.ID)
	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	resp.AccessToken = accessToken
	resp.RefreshToken = refreshToken
	resp.UserID = user.ID

	return resp, nil
}

func GetUserInfo(req *request.GetUserInfoReq) (response.GetUserInfoResp, exceptions.APIError) {
	var resp response.GetUserInfoResp

	user, err := db.GetOne[models.User](
		db.Equal("id", req.UserID),
	)

	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	if !user.Exists() {
		return resp, exceptions.BadRequestError(errUserNotFound, exceptions.UserNotExistsError)
	}

	resp.User = user

	credit, err := db.GetOne[models.Credit](
		db.Equal("user_id", req.UserID),
	)
	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	resp.Credit = credit

	return resp, nil
}

func SelectInterestTags(req *request.SelectInterestTagsReq) exceptions.APIError {
	userTags := make([]models.UserInterests, 0, len(req.Tags))

	for _, tag := range req.Tags {
		userTags = append(userTags, models.UserInterests{
			UserID:        req.UserID,
			InterestTagID: tag,
		})
	}

	err := db.WithTransaction(func(tx *gorm.DB) error {
		if len(userTags) == 0 {
			return nil
		}

		err := db.Create(userTags)
		if err != nil {
			return err
		}

		user, err := db.GetOne[models.User](
			db.Equal("id", req.UserID),
		)
		user.SelectedTags = true
		err = db.Update(user)
		return err
	})

	if err != nil {
		return exceptions.InternalServerError(err)
	}

	return nil
}
