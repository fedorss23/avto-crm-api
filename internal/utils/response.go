package utils

import (
	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool `json:"success"`
	Message string `json:"message"`
	Data interface{} `json:"data,omitempty"`
	Error interface{} `json:"error,omitempty"`
	Meta interface{} `json:"meta,omitempty"`
}

type PaginationMeta struct {
	CurrentPage int `json:"currentPage"`
	PerPage int `json:"perPage"`
	TotalItems int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
}

func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data: data,
	})
}

func SuccessResponseWithoutBody(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
	})
}

func SuccessResponseWithMeta(c *gin.Context, statusCode int, message string, data interface{}, meta interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data: data,
		Meta: meta,
	})
}

func ErrorResponse(c *gin.Context, statusCode int, message string, err error) {
	c.JSON(statusCode, APIResponse{
		Success: false,
		Message: message,
		Error: err.Error(),
	})
}

func ValidationErrorResponse(c *gin.Context, errors map[string]string) {
	c.JSON(400, APIResponse{
		Success: false,
		Message: "Ошибка валидации",
		Error: errors,
	})
}