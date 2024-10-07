package response

import "net/http"

// CommonSuccessData defines the structure of the success response
type CommonSuccessData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// NewSuccess is used to generate a success response
func NewSuccess(code int, message string, data interface{}) CommonSuccessData {
	return CommonSuccessData{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// NewOK returns a standard 200 OK response
func NewOK(data interface{}) CommonSuccessData {
	return NewSuccess(http.StatusOK, "Operation successful", data)
}

// NewCreated returns a standard 201 Created response
func NewCreated(data interface{}) CommonSuccessData {
	return NewSuccess(http.StatusCreated, "Resource created successfully", data)
}

// NewNoContent returns a standard 204 No Content response
func NewNoContent() CommonSuccessData {
	return NewSuccess(http.StatusNoContent, "No content", nil)
}
