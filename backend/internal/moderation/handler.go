package moderation

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/your-org/ai-tools-atlas-backend/internal/platform/responses"
	"github.com/your-org/ai-tools-atlas-backend/internal/tools"
)

// Handler handles HTTP requests for moderation
type Handler struct {
	service     Service
	toolService tools.Service
}

// NewHandler creates a new moderation handler
func NewHandler(service Service, toolService tools.Service) *Handler {
	return &Handler{
		service:     service,
		toolService: toolService,
	}
}

// RegisterRoutes registers public reporting routes
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup, optionalAuthMiddleware gin.HandlerFunc) {
	// Tool reports
	rg.POST("/tools/:slug/report", optionalAuthMiddleware, h.ReportTool)
	// Review reports
	rg.POST("/reviews/:id/report", optionalAuthMiddleware, h.ReportReview)
}

// RegisterAdminRoutes registers admin moderation routes
func (h *Handler) RegisterAdminRoutes(rg *gin.RouterGroup) {
	moderation := rg.Group("/moderation")
	{
		// Reports queue
		moderation.GET("/queue", h.ListModerationQueue)
		moderation.GET("/reports", h.ListPendingReports)
		moderation.GET("/reports/:id", h.GetReport)
		moderation.PATCH("/reports/:id", h.UpdateReportStatus)

		// Review moderation actions
		moderation.PATCH("/reviews/:id/approve", h.ApproveReview)
		moderation.PATCH("/reviews/:id/hide", h.HideReview)
		moderation.PATCH("/reviews/:id/remove", h.RemoveReview)

		// Moderation history
		moderation.GET("/history/:review_id", h.GetModerationHistory)
	}
}

// ReportTool handles POST /api/v1/tools/:slug/report
func (h *Handler) ReportTool(c *gin.Context) {
	slug := c.Param("slug")

	// Get tool by slug
	tool, err := h.toolService.GetToolBySlug(slug)
	if err != nil {
		if errors.Is(err, tools.ErrToolNotFound) {
			responses.Error(c, http.StatusNotFound, "NOT_FOUND", "Tool not found", nil)
			return
		}
		responses.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to find tool", nil)
		return
	}

	var input CreateReportInput
	if err := c.ShouldBindJSON(&input); err != nil {
		responses.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
		return
	}

	// Get user ID if authenticated
	var userID *uint
	if idVal, exists := c.Get("user_id"); exists {
		if uid, ok := idVal.(uint); ok {
			userID = &uid
		}
	}

	report, err := h.service.CreateToolReport(tool.ID, userID, input)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidReason):
			responses.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "Invalid reason. Must be: spam, abuse, misinformation, or other", nil)
		case errors.Is(err, ErrAlreadyReported):
			responses.Error(c, http.StatusConflict, "ALREADY_REPORTED", "You have already reported this tool today", nil)
		default:
			responses.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create report", nil)
		}
		return
	}

	responses.Created(c, report)
}

// ReportReview handles POST /api/v1/reviews/:id/report
func (h *Handler) ReportReview(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		responses.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid review ID", nil)
		return
	}

	var input CreateReportInput
	if err := c.ShouldBindJSON(&input); err != nil {
		responses.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
		return
	}

	// Get user ID if authenticated
	var userID *uint
	if uidVal, exists := c.Get("user_id"); exists {
		if uid, ok := uidVal.(uint); ok {
			userID = &uid
		}
	}

	report, err := h.service.CreateReviewReport(uint(id), userID, input)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidReason):
			responses.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "Invalid reason. Must be: spam, abuse, misinformation, or other", nil)
		case errors.Is(err, ErrAlreadyReported):
			responses.Error(c, http.StatusConflict, "ALREADY_REPORTED", "You have already reported this review today", nil)
		default:
			responses.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create report", nil)
		}
		return
	}

	responses.Created(c, report)
}

// ListPendingReports handles GET /api/v1/admin/moderation/reports
func (h *Handler) ListPendingReports(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	reports, total, err := h.service.ListPendingReports(page, pageSize)
	if err != nil {
		responses.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch reports", nil)
		return
	}

	responses.List(c, reports, map[string]interface{}{
		"page":      page,
		"page_size": pageSize,
		"total":     total,
	})
}

// GetReport handles GET /api/v1/admin/moderation/reports/:id
func (h *Handler) GetReport(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		responses.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid report ID", nil)
		return
	}

	report, err := h.service.GetReportByID(uint(id))
	if err != nil {
		if errors.Is(err, ErrReportNotFound) {
			responses.Error(c, http.StatusNotFound, "NOT_FOUND", "Report not found", nil)
			return
		}
		responses.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch report", nil)
		return
	}

	responses.Success(c, report)
}

// UpdateStatusInput represents input for updating report status
type UpdateStatusInput struct {
	Status string `json:"status" binding:"required"`
}

// UpdateReportStatus handles PATCH /api/v1/admin/moderation/reports/:id
func (h *Handler) UpdateReportStatus(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		responses.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid report ID", nil)
		return
	}

	var input UpdateStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		responses.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
		return
	}

	// Get reviewer's user ID
	reviewerIDVal, exists := c.Get("user_id")
	if !exists {
		responses.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}
	reviewerID, ok := reviewerIDVal.(uint)
	if !ok {
		responses.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Invalid user context", nil)
		return
	}

	err = h.service.UpdateReportStatus(uint(id), input.Status, reviewerID)
	if err != nil {
		switch {
		case errors.Is(err, ErrReportNotFound):
			responses.Error(c, http.StatusNotFound, "NOT_FOUND", "Report not found", nil)
		case errors.Is(err, ErrInvalidStatus):
			responses.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "Invalid status. Must be: pending, reviewed, or dismissed", nil)
		default:
			responses.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update report", nil)
		}
		return
	}

	responses.Success(c, map[string]string{"message": "Report status updated"})
}

// ListModerationQueue handles GET /api/v1/admin/moderation/queue
func (h *Handler) ListModerationQueue(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	filters := ReportFilters{
		Type:   c.Query("type"),
		Status: c.DefaultQuery("status", "pending"),
	}

	reports, total, err := h.service.ListReports(filters, page, pageSize)
	if err != nil {
		responses.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch moderation queue", nil)
		return
	}

	responses.List(c, reports, map[string]interface{}{
		"page":      page,
		"page_size": pageSize,
		"total":     total,
	})
}

// ApproveReview handles PATCH /api/v1/admin/moderation/reviews/:id/approve
func (h *Handler) ApproveReview(c *gin.Context) {
	h.handleReviewModeration(c, "approve")
}

// HideReview handles PATCH /api/v1/admin/moderation/reviews/:id/hide
func (h *Handler) HideReview(c *gin.Context) {
	h.handleReviewModeration(c, "hide")
}

// RemoveReview handles PATCH /api/v1/admin/moderation/reviews/:id/remove
func (h *Handler) RemoveReview(c *gin.Context) {
	h.handleReviewModeration(c, "remove")
}

func (h *Handler) handleReviewModeration(c *gin.Context, action string) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		responses.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid review ID", nil)
		return
	}

	var input ModerationActionInput
	// Input is optional, so we don't fail on binding error
	_ = c.ShouldBindJSON(&input)

	// Get moderator's user ID
	moderatorIDVal, exists := c.Get("user_id")
	if !exists {
		responses.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}
	moderatorID, ok := moderatorIDVal.(uint)
	if !ok {
		responses.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Invalid user context", nil)
		return
	}

	var review interface{}
	switch action {
	case "approve":
		review, err = h.service.ApproveReview(uint(id), moderatorID, input)
	case "hide":
		review, err = h.service.HideReview(uint(id), moderatorID, input)
	case "remove":
		review, err = h.service.RemoveReview(uint(id), moderatorID, input)
	}

	if err != nil {
		switch {
		case errors.Is(err, ErrReviewNotFound):
			responses.Error(c, http.StatusNotFound, "NOT_FOUND", "Review not found", nil)
		default:
			responses.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to "+action+" review", nil)
		}
		return
	}

	responses.Success(c, review)
}

// GetModerationHistory handles GET /api/v1/admin/moderation/history/:review_id
func (h *Handler) GetModerationHistory(c *gin.Context) {
	idParam := c.Param("review_id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		responses.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid review ID", nil)
		return
	}

	history, err := h.service.GetModerationHistory(uint(id))
	if err != nil {
		switch {
		case errors.Is(err, ErrReviewNotFound):
			responses.Error(c, http.StatusNotFound, "NOT_FOUND", "Review not found", nil)
		default:
			responses.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch moderation history", nil)
		}
		return
	}

	responses.Success(c, history)
}
