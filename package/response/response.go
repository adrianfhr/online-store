package response

import (
	"github.com/gin-gonic/gin"
)

// Response structure for both success and error responses
type CommonResponse struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// RespondSuccess creates a formatted success response
func RespondSuccess(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, CommonResponse{
		Code:    code,
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

// RespondError creates a formatted error response
func RespondError(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, CommonResponse{
		Code:    code,
		Status:  "error",
		Message: message,
		Data:    data,
	})
}

// Success message helper function
func GetSuccess() string {
	return "Operation successful"
}

// Error message helper function
func GetError() string {
	return "Operation failed"
}
