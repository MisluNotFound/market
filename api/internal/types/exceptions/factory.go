package exceptions

import "github.com/mislu/market-api/pkg/entities"

type APIError interface {
	error

	ToResponse() *entities.Response
	IsError() bool
}

func InternalServerError(err error) APIError {
	return NewGenericError(500, "Internal server error", err)
}

func BadRequestError(err error, msg string) APIError {
	return NewGenericError(400, msg, err)
}
