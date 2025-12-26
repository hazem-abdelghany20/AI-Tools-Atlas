package http

import (
	"github.com/gin-gonic/gin"
	"github.com/your-org/ai-tools-atlas-backend/internal/platform/config"
)

// SetupRouter initializes the Gin router with middleware
func SetupRouter(cfg *config.Config) *gin.Engine {
	// Create router
	r := gin.New()

	// Add middleware
	r.Use(gin.Recovery())
	r.Use(CORSMiddleware(cfg.AllowedOrigins))
	r.Use(RequestLoggerMiddleware())

	// Health check endpoint (outside versioned API)
	r.GET("/health", HealthCheck)

	// Create API v1 group
	v1 := r.Group("/api/v1")
	{
		// Future endpoints will be added here
		_ = v1
	}

	return r
}
