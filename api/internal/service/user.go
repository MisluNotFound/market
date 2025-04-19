package service

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	resourcemanager "github.com/mislu/market-api/internal/core/resource_manager"
	"github.com/mislu/market-api/internal/db"
	"github.com/mislu/market-api/internal/types/exceptions"
	"github.com/mislu/market-api/internal/types/models"
	"github.com/mislu/market-api/internal/types/request"
	"github.com/mislu/market-api/internal/types/response"
	"github.com/mislu/market-api/internal/utils/app"
	"github.com/mislu/market-api/internal/utils/lib"
)

var mobileRegex = regexp.MustCompile(`^1[3-9][0-9]{9}$`)
var (
	errUserNotFound = errors.New("user does not found")
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

	user, err := db.GetOne[*models.User](
		db.Equal("id", req.UserID),
	)

	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	if !user.Exists() {
		return resp, exceptions.BadRequestError(errUserNotFound, exceptions.UserNotExistsError)
	}

	key := resourcemanager.GenerateObjectKey(resourcemanager.UserAvatarBucket, req.UserID, req.Avatar.Filename)

	avatarFile, err := req.Avatar.Open()
	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	data := make([]byte, 1024 * 1024 * 10)
	n, err := avatarFile.Read(data)
	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	data = data[:n]

	err = resourcemanager.UploadFile(resourcemanager.UserAvatarBucket, key, data)
	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	urlParts := strings.Split(user.Avatar, "/")
	oldAvatarKey := urlParts[len(urlParts)-1]

	// http://ip/assert/key
	user.Avatar = fmt.Sprintf("http://%s/assert/%s", app.GetConfig().Server.BaseIP, key)

	// 忽略删除错误，由定时清理任务清除多余的文件
	resourcemanager.DeleteFile(resourcemanager.UserAvatarBucket, resourcemanager.GetObjectPath(resourcemanager.UserAvatarBucket, req.UserID, oldAvatarKey))

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

	return resp, nil
}

func GetUserInfo(req *request.GetUserInfoReq) (response.GetUserInfoResp, exceptions.APIError) {
	var resp response.GetUserInfoResp

	user, err := db.GetOne[*models.User](
		db.Equal("id", req.UserID),
	)

	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	if !user.Exists() {
		return resp, exceptions.BadRequestError(errUserNotFound, exceptions.UserNotExistsError)
	}

	resp.User = *user
	// TODO get user products
	return resp, nil
}
