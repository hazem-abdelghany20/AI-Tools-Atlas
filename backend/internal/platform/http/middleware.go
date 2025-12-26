package http

import (
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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
