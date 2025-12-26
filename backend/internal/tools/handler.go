package tools

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/your-org/ai-tools-atlas-backend/internal/platform/responses"
)

// Handler handles HTTP requests for tools
type Handler struct {
	service Service
}

// NewHandler creates a new tool handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers tool routes on the given router group
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	tools := rg.Group("/tools")
	{
		tools.GET("", h.ListTools)
		tools.GET("/:slug", h.GetTool)
		tools.GET("/:slug/alternatives", h.GetToolAlternatives)
	}

	search := rg.Group("/search")
	{
		search.GET("/tools", h.SearchTools)
	}
}

// RegisterAdminRoutes registers admin tool routes
func (h *Handler) RegisterAdminRoutes(rg *gin.RouterGroup) {
	tools := rg.Group("/tools")
	{
		tools.GET("", h.AdminListTools)
		tools.POST("", h.AdminCreateTool)
		tools.GET("/:id", h.AdminGetTool)
		tools.PATCH("/:id", h.AdminUpdateTool)
		tools.DELETE("/:id", h.AdminArchiveTool)
	}
}

// ListTools handles GET /api/v1/tools
func (h *Handler) ListTools(c *gin.Context) {
	filters := h.parseFilters(c)
	page, pageSize := h.parsePagination(c)

	tools, total, err := h.service.ListTools(filters, page, pageSize)
	if err != nil {
		responses.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch tools", nil)
		return
	}

	responses.List(c, tools, map[string]interface{}{
		"page":      page,
		"page_size": pageSize,
		"total":     total,
	})
}

// SearchTools handles GET /api/v1/search/tools
func (h *Handler) SearchTools(c *gin.Context) {
	query := c.Query("q")
	filters := h.parseFilters(c)
	page, pageSize := h.parsePagination(c)

	tools, total, err := h.service.SearchTools(query, filters, page, pageSize)
	if err != nil {
		responses.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to search tools", nil)
		return
	}

	responses.List(c, tools, map[string]interface{}{
		"page":      page,
		"page_size": pageSize,
		"total":     total,
		"query":     query,
	})
}

// GetTool handles GET /api/v1/tools/:slug
func (h *Handler) GetTool(c *gin.Context) {
	slug := c.Param("slug")

	tool, err := h.service.GetToolBySlug(slug)
	if err != nil {
		if errors.Is(err, ErrToolNotFound) {
			responses.Error(c, http.StatusNotFound, "NOT_FOUND", "Tool not found", nil)
			return
		}
		responses.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch tool", nil)
		return
	}

	responses.Success(c, tool)
}

// GetToolAlternatives handles GET /api/v1/tools/:slug/alternatives
func (h *Handler) GetToolAlternatives(c *gin.Context) {
	slug := c.Param("slug")

	result, err := h.service.GetToolAlternatives(slug)
	if err != nil {
		if errors.Is(err, ErrToolNotFound) {
			responses.Error(c, http.StatusNotFound, "NOT_FOUND", "Tool not found", nil)
			return
		}
		responses.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch alternatives", nil)
		return
	}

	responses.Success(c, map[string]interface{}{
		"similar":      result.Similar,
		"alternatives": result.Alternatives,
	})
}

// parseFilters extracts filter parameters from the request
func (h *Handler) parseFilters(c *gin.Context) ToolFilters {
	minRating, _ := strconv.ParseFloat(c.Query("min_rating"), 64)

	return ToolFilters{
		Category:  c.Query("category"),
		Price:     c.Query("price"),
		MinRating: minRating,
		Platform:  c.Query("platform"),
		Sort:      c.DefaultQuery("sort", SortTopRated),
	}
}

// parsePagination extracts pagination parameters from the request
func (h *Handler) parsePagination(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	return page, pageSize
}

// AdminListTools handles GET /api/v1/admin/tools
func (h *Handler) AdminListTools(c *gin.Context) {
	search := c.Query("search")
	includeArchived := c.Query("archived") == "true"
	page, pageSize := h.parsePagination(c)

	tools, total, err := h.service.ListToolsAdmin(search, includeArchived, page, pageSize)
	if err != nil {
		responses.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch tools", nil)
		return
	}

	responses.List(c, tools, map[string]interface{}{
		"page":      page,
		"page_size": pageSize,
		"total":     total,
	})
}

// AdminGetTool handles GET /api/v1/admin/tools/:id
func (h *Handler) AdminGetTool(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		responses.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid tool ID", nil)
		return
	}

	tool, err := h.service.GetToolByIDAdmin(uint(id))
	if err != nil {
		if errors.Is(err, ErrToolNotFound) {
			responses.Error(c, http.StatusNotFound, "NOT_FOUND", "Tool not found", nil)
			return
		}
		responses.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch tool", nil)
		return
	}

	responses.Success(c, tool)
}

// AdminCreateTool handles POST /api/v1/admin/tools
func (h *Handler) AdminCreateTool(c *gin.Context) {
	var input CreateToolInput
	if err := c.ShouldBindJSON(&input); err != nil {
		responses.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
		return
	}

	tool, err := h.service.CreateTool(input)
	if err != nil {
		switch {
		case errors.Is(err, ErrSlugRequired):
			responses.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "Slug is required", map[string]string{"slug": "required"})
		case errors.Is(err, ErrNameRequired):
			responses.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "Name is required", map[string]string{"name": "required"})
		case errors.Is(err, ErrSlugExists):
			responses.Error(c, http.StatusConflict, "SLUG_EXISTS", "A tool with this slug already exists", nil)
		case errors.Is(err, ErrCategoryRequired):
			responses.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "Category is required", map[string]string{"primary_category_id": "required"})
		default:
			responses.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create tool", nil)
		}
		return
	}

	responses.Created(c, tool)
}

// AdminUpdateTool handles PATCH /api/v1/admin/tools/:id
func (h *Handler) AdminUpdateTool(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		responses.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid tool ID", nil)
		return
	}

	var input UpdateToolInput
	if err := c.ShouldBindJSON(&input); err != nil {
		responses.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
		return
	}

	tool, err := h.service.UpdateTool(uint(id), input)
	if err != nil {
		switch {
		case errors.Is(err, ErrToolNotFound):
			responses.Error(c, http.StatusNotFound, "NOT_FOUND", "Tool not found", nil)
		case errors.Is(err, ErrNameRequired):
			responses.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "Name cannot be empty", map[string]string{"name": "required"})
		case errors.Is(err, ErrCategoryRequired):
			responses.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "Category cannot be empty", map[string]string{"primary_category_id": "required"})
		default:
			responses.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update tool", nil)
		}
		return
	}

	responses.Success(c, tool)
}

// AdminArchiveTool handles DELETE /api/v1/admin/tools/:id
func (h *Handler) AdminArchiveTool(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		responses.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid tool ID", nil)
		return
	}

	err = h.service.ArchiveTool(uint(id))
	if err != nil {
		if errors.Is(err, ErrToolNotFound) {
			responses.Error(c, http.StatusNotFound, "NOT_FOUND", "Tool not found", nil)
			return
		}
		responses.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to archive tool", nil)
		return
	}

	c.Status(http.StatusNoContent)
}
