package response

import "net/http"

type CommonErrorData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func NewError(code int, message string, data interface{}) CommonErrorData {
	return CommonErrorData{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func NewBadRequest(data interface{}) CommonErrorData {
	return NewError(http.StatusBadRequest, "Bad Request", data)
}

func NewNotFound(data interface{}) CommonErrorData {
	return NewError(http.StatusNotFound, "Not Found", data)
}

func NewUnauthorized(data interface{}) CommonErrorData {
	return NewError(http.StatusUnauthorized, "Unauthorized", data)
}

func NewConflict(data interface{}) CommonErrorData {
	return NewError(http.StatusConflict, "Conflict", data)
}

func NewInternalServerError(data interface{}) CommonErrorData {
	return NewError(http.StatusInternalServerError, "Internal Server Error", data)
}

type ErrorString struct {
	code    int
	message string
}

func NewErrorString(code int, msg string) ErrorString {
	return ErrorString{
		code:    code,
		message: msg,
	}
}

func (e ErrorString) Code() int {
	return e.code
}

func (e ErrorString) Error() string {
	return e.message
}

func (e ErrorString) Message() string {
	return e.message
}

func BadRequest(msg string) error {
	return NewErrorString(http.StatusBadRequest, msg)
}

func NotFound(msg string) error {
	return NewErrorString(http.StatusNotFound, msg)
}

func Conflict(msg string) error {
	return NewErrorString(http.StatusConflict, msg)
}

func InternalServerError(msg string) error {
	return NewErrorString(http.StatusInternalServerError, msg)
}

func UnauthorizedMessage(msg string) error {
	return NewErrorString(http.StatusUnauthorized, msg)
}

func ForbiddenMessage(msg string) error {
	return NewErrorString(http.StatusForbidden, msg)
}
