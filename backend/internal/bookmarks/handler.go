package bookmarks

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/your-org/ai-tools-atlas-backend/internal/platform/responses"
)

// AddBookmarkRequest represents the request body for adding a bookmark
type AddBookmarkRequest struct {
	ToolID uint `json:"tool_id" binding:"required"`
}

// Handler handles HTTP requests for bookmarks
type Handler struct {
	service Service
}

// NewHandler creates a new bookmarks handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers bookmark routes on the given router group
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup, authMiddleware gin.HandlerFunc, optionalAuthMiddleware gin.HandlerFunc) {
	me := rg.Group("/me")
	{
		me.GET("/bookmarks", optionalAuthMiddleware, h.GetBookmarks)
		me.POST("/bookmarks", optionalAuthMiddleware, h.AddBookmark)
		me.DELETE("/bookmarks/:tool_id", optionalAuthMiddleware, h.RemoveBookmark)
	}
}

// GetBookmarks handles GET /api/v1/me/bookmarks
func (h *Handler) GetBookmarks(c *gin.Context) {
	userID, sessionID := h.getUserOrSession(c)

	bookmarks, err := h.service.GetBookmarks(userID, sessionID)
	if err != nil {
		responses.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch bookmarks", nil)
		return
	}

	// Extract tools from bookmarks for simpler response
	tools := make([]interface{}, len(bookmarks))
	for i, b := range bookmarks {
		tools[i] = b.Tool
	}

	responses.Success(c, tools)
}

// AddBookmark handles POST /api/v1/me/bookmarks
func (h *Handler) AddBookmark(c *gin.Context) {
	userID, sessionID := h.getUserOrSession(c)

	// Ensure we have either user or session
	if userID == 0 && sessionID == "" {
		sessionID = h.createSessionID(c)
	}

	var req AddBookmarkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "tool_id is required", nil)
		return
	}

	bookmark, err := h.service.AddBookmark(userID, sessionID, req.ToolID)
	if err != nil {
		if errors.Is(err, ErrAlreadyBookmarked) {
			responses.Error(c, http.StatusConflict, "ALREADY_BOOKMARKED", "Tool is already bookmarked", nil)
			return
		}
		responses.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to add bookmark", nil)
		return
	}

	responses.Created(c, bookmark)
}

// RemoveBookmark handles DELETE /api/v1/me/bookmarks/:tool_id
func (h *Handler) RemoveBookmark(c *gin.Context) {
	userID, sessionID := h.getUserOrSession(c)

	toolIDStr := c.Param("tool_id")
	toolID, err := strconv.ParseUint(toolIDStr, 10, 32)
	if err != nil {
		responses.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid tool_id", nil)
		return
	}

	err = h.service.RemoveBookmark(userID, sessionID, uint(toolID))
	if err != nil {
		if errors.Is(err, ErrBookmarkNotFound) {
			responses.Error(c, http.StatusNotFound, "NOT_FOUND", "Bookmark not found", nil)
			return
		}
		responses.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to remove bookmark", nil)
		return
	}

	responses.NoContent(c)
}

// getUserOrSession extracts user_id from auth context or session_id from cookie
func (h *Handler) getUserOrSession(c *gin.Context) (uint, string) {
	// Try to get authenticated user first
	if userIDVal, exists := c.Get("user_id"); exists {
		if userID, ok := userIDVal.(uint); ok && userID > 0 {
			return userID, ""
		}
	}

	// Fall back to session ID
	sessionID, _ := c.Cookie("session_id")
	return 0, sessionID
}

// createSessionID generates a new session ID and sets it as a cookie
func (h *Handler) createSessionID(c *gin.Context) string {
	sessionID := uuid.New().String()
	maxAge := 365 * 24 * 60 * 60 // 1 year

	// Use secure cookies when not in development mode
	secure := c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https"

	c.SetCookie(
		"session_id",
		sessionID,
		maxAge,
		"/",
		"",
		secure,
		true, // httpOnly - prevent XSS access
	)

	return sessionID
}
