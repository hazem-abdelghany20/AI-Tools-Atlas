package http

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// SetAuthCookie sets the JWT token as an HTTP-only cookie
func SetAuthCookie(c *gin.Context, token string) {
	secure := os.Getenv("GIN_MODE") == "release" // Secure only in production
	maxAge := 7 * 24 * 60 * 60                   // 7 days in seconds

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		"auth_token", // name
		token,        // value
		maxAge,       // max age in seconds
		"/",          // path
		"",           // domain (empty = current domain)
		secure,       // secure (HTTPS only in production)
		true,         // httpOnly (prevent JS access)
	)
}

// ClearAuthCookie removes the auth cookie (for logout)
func ClearAuthCookie(c *gin.Context) {
	c.SetCookie("auth_token", "", -1, "/", "", false, true)
}

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
