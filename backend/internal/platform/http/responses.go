package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SuccessResponse sends a successful response with data
func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

// ListResponse sends a list response with data and metadata
func ListResponse(c *gin.Context, data interface{}, meta map[string]interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"data": data,
		"meta": meta,
	})
}

// ErrorResponse sends an error response
func ErrorResponse(c *gin.Context, statusCode int, code string, message string, details interface{}) {
	c.JSON(statusCode, gin.H{
		"error": gin.H{
			"code":    code,
			"message": message,
			"details": details,
		},
	})
}
