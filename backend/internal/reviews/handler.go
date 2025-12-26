package reviews

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/your-org/ai-tools-atlas-backend/internal/platform/responses"
)

// Handler handles HTTP requests for reviews
type Handler struct {
	service Service
}

// NewHandler creates a new reviews handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers review routes on the given router group
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	tools := rg.Group("/tools")
	{
		tools.GET("/:slug/reviews", h.ListReviews)
		tools.POST("/:slug/reviews", authMiddleware, h.CreateReview)
	}

	// User reviews (authenticated)
	rg.GET("/me/reviews", authMiddleware, h.GetUserReviews)
}

// ListReviews handles GET /api/v1/tools/:slug/reviews
func (h *Handler) ListReviews(c *gin.Context) {
	slug := c.Param("slug")
	page, pageSize := h.parsePagination(c)
	sort := c.DefaultQuery("sort", SortNewest)

	reviews, total, err := h.service.ListReviews(slug, sort, page, pageSize)
	if err != nil {
		if errors.Is(err, ErrToolNotFound) {
			responses.Error(c, http.StatusNotFound, "NOT_FOUND", "Tool not found", nil)
			return
		}
		responses.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch reviews", nil)
		return
	}

	responses.List(c, reviews, map[string]interface{}{
		"page":      page,
		"page_size": pageSize,
		"total":     total,
	})
}

// CreateReview handles POST /api/v1/tools/:slug/reviews
func (h *Handler) CreateReview(c *gin.Context) {
	slug := c.Param("slug")

	// Get user ID from auth context
	userIDVal, exists := c.Get("user_id")
	if !exists {
		responses.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required", nil)
		return
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		responses.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Invalid user context", nil)
		return
	}

	// Parse request body
	var input CreateReviewInput
	if err := c.ShouldBindJSON(&input); err != nil {
		responses.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
		return
	}

	// Create review
	review, err := h.service.CreateReview(slug, userID, input)
	if err != nil {
		switch {
		case errors.Is(err, ErrToolNotFound):
			responses.Error(c, http.StatusNotFound, "NOT_FOUND", "Tool not found", nil)
		case errors.Is(err, ErrAlreadyReviewed):
			responses.Error(c, http.StatusConflict, "ALREADY_REVIEWED", "You have already reviewed this tool", nil)
		case errors.Is(err, ErrRatingRequired):
			responses.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "Rating is required", map[string]string{"rating_overall": "required"})
		case errors.Is(err, ErrInvalidRating):
			responses.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "Rating must be between 1 and 5", map[string]string{"rating": "invalid"})
		case errors.Is(err, ErrProsRequired):
			responses.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "Pros field is required", map[string]string{"pros": "required"})
		case errors.Is(err, ErrConsRequired):
			responses.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "Cons field is required", map[string]string{"cons": "required"})
		case errors.Is(err, ErrProsTooLong):
			responses.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "Pros must be 500 characters or less", map[string]string{"pros": "too_long"})
		case errors.Is(err, ErrConsTooLong):
			responses.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "Cons must be 500 characters or less", map[string]string{"cons": "too_long"})
		default:
			responses.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create review", nil)
		}
		return
	}

	responses.Created(c, review)
}

// GetUserReviews handles GET /api/v1/me/reviews
func (h *Handler) GetUserReviews(c *gin.Context) {
	// Get user ID from auth context
	userIDVal, exists := c.Get("user_id")
	if !exists {
		responses.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required", nil)
		return
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		responses.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Invalid user context", nil)
		return
	}

	page, pageSize := h.parsePagination(c)

	reviews, total, err := h.service.ListUserReviews(userID, page, pageSize)
	if err != nil {
		responses.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch reviews", nil)
		return
	}

	responses.List(c, reviews, map[string]interface{}{
		"page":      page,
		"page_size": pageSize,
		"total":     total,
	})
}

// parsePagination extracts pagination parameters from the request
func (h *Handler) parsePagination(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	return page, pageSize
}
