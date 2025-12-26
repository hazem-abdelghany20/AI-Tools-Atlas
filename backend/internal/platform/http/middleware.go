package http

import (
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/your-org/ai-tools-atlas-backend/internal/auth"
)

// CORSMiddleware handles Cross-Origin Resource Sharing
func CORSMiddleware(allowedOrigins string) gin.HandlerFunc {
	// Parse allowed origins into a slice
	origins := strings.Split(allowedOrigins, ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}

	return func(c *gin.Context) {
		requestOrigin := c.Request.Header.Get("Origin")

		// Check if request origin is in allowed list
		allowed := false
		for _, origin := range origins {
			if origin == "*" || origin == requestOrigin {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				allowed = true
				break
			}
		}

		// If no match found and we have allowed origins, use the first one as fallback
		// This handles cases where Origin header is missing (like direct API calls)
		if !allowed && len(origins) > 0 {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origins[0])
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// RequestLoggerMiddleware logs incoming requests
func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		duration := time.Since(startTime)
		statusCode := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path

		log.Printf("[%s] %s - %d - %v", method, path, statusCode, duration)
	}
}

// AuthRequired middleware verifies JWT token from cookie
func AuthRequired(authService *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from cookie
		cookie, err := c.Cookie("auth_token")
		if err != nil {
			ErrorResponse(c, 401, "unauthorized", "Authentication required", nil)
			c.Abort()
			return
		}

		// Validate token
		claims, err := authService.ValidateToken(cookie)
		if err != nil {
			ErrorResponse(c, 401, "unauthorized", "Invalid token", nil)
			c.Abort()
			return
		}

		// Set user context
		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)
		c.Set("user_email", claims.Email)

		c.Next()
	}
}

// AdminRequired middleware checks if user has admin role
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists {
			ErrorResponse(c, 403, "forbidden", "Access denied: role not found", nil)
			c.Abort()
			return
		}

		if role != "admin" {
			ErrorResponse(c, 403, "forbidden", "Access denied: admin role required", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalAuth middleware attempts to authenticate but doesn't require it
func OptionalAuth(authService *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try to get token from cookie
		cookie, err := c.Cookie("auth_token")
		if err != nil {
			// No token, continue without authentication
			c.Next()
			return
		}

		// Try to validate token
		claims, err := authService.ValidateToken(cookie)
		if err != nil {
			// Invalid token, continue without authentication
			c.Next()
			return
		}

		// Set user context if valid
		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)
		c.Set("user_email", claims.Email)

		c.Next()
	}
}
