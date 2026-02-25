package response

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type APIResponse struct {
	RequestID string    `json:"request_id"`
	Status    int       `json:"status"`
	Success   bool      `json:"success"`
	Message   string    `json:"message"`
	Data      any       `json:"data,omitempty"`
	Error     *APIError `json:"error,omitempty"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Success(c *gin.Context, status int, data any) {
	c.JSON(status, APIResponse{
		RequestID: uuid.New().String(),
		Status:    status,
		Success:   true,
		Message:   "success",
		Data:      data,
	})
}

func Error(c *gin.Context, status int, code string, message string) {
	c.JSON(status, APIResponse{
		RequestID: uuid.New().String(),
		Status:    status,
		Success:   false,
		Message:   message,
		Error: &APIError{
			Code:    code,
			Message: message,
		},
	})
}
