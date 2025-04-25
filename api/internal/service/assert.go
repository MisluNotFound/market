package service

import (
	"errors"
	"path/filepath"

	resourcemanager "github.com/mislu/market-api/internal/core/resource_manager"
	"github.com/mislu/market-api/internal/types/exceptions"
	"github.com/mislu/market-api/internal/types/request"
	"github.com/mislu/market-api/internal/utils/app"
)

func GetAssert(req *request.GetAssertReq) (string, exceptions.APIError) {
	path := resourcemanager.GetObjectPath(resourcemanager.BucketType(req.Type), req.Owner, req.Key)
	exists, err := resourcemanager.FileExists(path)
	if err != nil {
		return "", exceptions.InternalServerError(err)
	}

	if !exists {
		return "", exceptions.BadRequestError(errors.New("resource not found"), exceptions.ResourceNotFoundError)
	}

	path = filepath.Join(app.GetConfig().OSS.Root, path)
	return path, nil
}

func RedirectAssert(req *request.GetAssertReq) (string, exceptions.APIError) {
	// TODO implement

	return "", nil
}
