package entities

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data any `json:"data,omitempty"`
}


func NewSuccessResponse(data any) *Response {
	return &Response {
		Code: 200,
		Msg: "success",
		Data: data,
	}
}

func NewErrorResponse(code int, msg string) *Response {
	return &Response {
		Code: code,
		Msg: msg,
	}
}