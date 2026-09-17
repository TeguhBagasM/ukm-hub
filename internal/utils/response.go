package utils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, status int, data any, message string) {
	c.JSON(status, gin.H{
		"success": true,
		"data":    data,
		"message": message,
	})
}

func ErrorResponse(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"success": false,
		"message": message,
	})
}

func WriteError(c *gin.Context, err error) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		ErrorResponse(c, appErr.Status, appErr.Message)
		return
	}
	ErrorResponse(c, http.StatusInternalServerError, "internal server error")
}
