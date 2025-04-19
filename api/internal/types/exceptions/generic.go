package exceptions

import "github.com/mislu/market-api/pkg/entities"

type genericError struct {
	Err  error  `json:"-"`
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}

func (e *genericError) Error() string {
	return e.Err.Error()
}

func (e *genericError) ToResponse() *entities.Response {
	return &entities.Response{
		Code: e.Code,
		Msg:  e.Msg,
	}
}

func (e *genericError) IsError() bool {
	return e.Err == nil
}

func NewGenericError(code int, msg string, err error) APIError {
	return &genericError{
		Err: err,
		Msg:  msg,
		Code: code,
	}
}
