package tools

import (
	"strings"

	"github.com/your-org/ai-tools-atlas-backend/internal/domain"
	"gorm.io/gorm"
)

// AlternativesResult holds similar and alternative tools
type AlternativesResult struct {
	Similar      []Tool
	Alternatives []Tool
}

// Repository defines the interface for tool data operations
type Repository interface {
	ListTools(filters ToolFilters, page, pageSize int) ([]Tool, int64, error)
	SearchTools(query string, filters ToolFilters, page, pageSize int) ([]Tool, int64, error)
	GetToolBySlug(slug string) (*Tool, error)
	GetToolByID(id uint) (*Tool, error)
	GetToolAlternatives(toolID uint, limit int) (*AlternativesResult, error)
	// Admin methods
	ListToolsAdmin(search string, includeArchived bool, page, pageSize int) ([]Tool, int64, error)
	GetToolByIDAdmin(id uint) (*Tool, error)
	Create(tool *Tool) error
	Update(tool *Tool) error
	Archive(id uint) error
	SlugExists(slug string, excludeID uint) (bool, error)
}

// repository implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new tool repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// ListTools returns paginated tools with filters
func (r *repository) ListTools(filters ToolFilters, page, pageSize int) ([]Tool, int64, error) {
	var tools []Tool
	var total int64

	query := r.buildBaseQuery(filters)

	// Count total
	if err := query.Model(&domain.Tool{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	query = r.applySorting(query, filters.Sort)

	// Apply pagination
	offset := (page - 1) * pageSize
	query = query.
		Preload("PrimaryCategory").
		Preload("Tags").
		Preload("Badges").
		Limit(pageSize).
		Offset(offset)

	if err := query.Find(&tools).Error; err != nil {
		return nil, 0, err
	}

	return tools, total, nil
}

// SearchTools performs a free-text search with filters
func (r *repository) SearchTools(searchQuery string, filters ToolFilters, page, pageSize int) ([]Tool, int64, error) {
	var tools []Tool
	var total int64

	query := r.buildBaseQuery(filters)

	// Add search conditions if query is not empty
	if searchQuery != "" {
		searchTerm := "%" + strings.ToLower(searchQuery) + "%"
		query = query.Where(
			"LOWER(tools.name) LIKE ? OR LOWER(tools.tagline) LIKE ? OR LOWER(tools.description) LIKE ? OR LOWER(tools.best_for) LIKE ? OR LOWER(tools.primary_use_cases) LIKE ?",
			searchTerm, searchTerm, searchTerm, searchTerm, searchTerm,
		)
	}

	// Count total
	if err := query.Model(&domain.Tool{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	query = r.applySorting(query, filters.Sort)

	// Apply pagination
	offset := (page - 1) * pageSize
	query = query.
		Preload("PrimaryCategory").
		Preload("Tags").
		Preload("Badges").
		Limit(pageSize).
		Offset(offset)

	if err := query.Find(&tools).Error; err != nil {
		return nil, 0, err
	}

	return tools, total, nil
}

// GetToolBySlug finds a tool by its slug
func (r *repository) GetToolBySlug(slug string) (*Tool, error) {
	var tool Tool
	err := r.db.
		Preload("PrimaryCategory").
		Preload("Tags").
		Preload("Badges").
		Preload("Media").
		Where("slug = ? AND archived_at IS NULL", slug).
		First(&tool).Error
	if err != nil {
		return nil, err
	}
	return &tool, nil
}

// GetToolByID finds a tool by its ID
func (r *repository) GetToolByID(id uint) (*Tool, error) {
	var tool Tool
	err := r.db.
		Preload("PrimaryCategory").
		Preload("Tags").
		Preload("Badges").
		Preload("Media").
		Where("id = ? AND archived_at IS NULL", id).
		First(&tool).Error
	if err != nil {
		return nil, err
	}
	return &tool, nil
}

// buildBaseQuery creates the base query with common filters
func (r *repository) buildBaseQuery(filters ToolFilters) *gorm.DB {
	query := r.db.Model(&domain.Tool{}).Where("tools.archived_at IS NULL")

	// Filter by category
	if filters.Category != "" {
		query = query.
			Joins("JOIN categories ON categories.id = tools.primary_category_id").
			Where("categories.slug = ?", filters.Category)
	}

	// Filter by price
	if filters.Price != "" {
		switch filters.Price {
		case "free":
			query = query.Where("tools.has_free_tier = ? AND (tools.pricing_summary ILIKE ? OR tools.pricing_summary ILIKE ?)", true, "%free%", "%$0%")
		case "freemium":
			query = query.Where("tools.has_free_tier = ?", true)
		case "paid":
			query = query.Where("tools.has_free_tier = ?", false)
		}
	}

	// Filter by minimum rating
	if filters.MinRating > 0 {
		query = query.Where("tools.avg_rating_overall >= ?", filters.MinRating)
	}

	// Filter by platform
	if filters.Platform != "" {
		platformTerm := "%" + strings.ToLower(filters.Platform) + "%"
		query = query.Where("LOWER(tools.platforms) LIKE ?", platformTerm)
	}

	return query
}

// applySorting adds the appropriate ORDER BY clause
func (r *repository) applySorting(query *gorm.DB, sort string) *gorm.DB {
	switch sort {
	case SortMostBookmarked:
		return query.Order("tools.bookmark_count DESC")
	case SortTrending:
		return query.Order("tools.trending_score DESC")
	case SortNewest:
		return query.Order("tools.created_at DESC")
	case SortTopRated:
		fallthrough
	default:
		return query.Order("tools.avg_rating_overall DESC, tools.review_count DESC")
	}
}

// GetToolAlternatives returns similar and alternative tools for a given tool
func (r *repository) GetToolAlternatives(toolID uint, limit int) (*AlternativesResult, error) {
	result := &AlternativesResult{
		Similar:      []Tool{},
		Alternatives: []Tool{},
	}

	// Get tool alternatives with relationship type
	var altRecords []domain.ToolAlternative
	err := r.db.Where("tool_id = ?", toolID).Find(&altRecords).Error
	if err != nil {
		return nil, err
	}

	// Collect alternative tool IDs by type
	var similarIDs, altIDs []uint
	for _, record := range altRecords {
		if record.RelationshipType == "similar" {
			similarIDs = append(similarIDs, record.AlternativeToolID)
		} else if record.RelationshipType == "alternative" {
			altIDs = append(altIDs, record.AlternativeToolID)
		}
	}

	// Fetch similar tools
	if len(similarIDs) > 0 {
		var similar []Tool
		err := r.db.
			Preload("PrimaryCategory").
			Preload("Tags").
			Where("id IN ? AND archived_at IS NULL", similarIDs).
			Limit(limit).
			Find(&similar).Error
		if err != nil {
			return nil, err
		}
		result.Similar = similar
	}

	// Fetch alternative tools
	if len(altIDs) > 0 {
		var alternatives []Tool
		err := r.db.
			Preload("PrimaryCategory").
			Preload("Tags").
			Where("id IN ? AND archived_at IS NULL", altIDs).
			Limit(limit).
			Find(&alternatives).Error
		if err != nil {
			return nil, err
		}
		result.Alternatives = alternatives
	}

	return result, nil
}

// ListToolsAdmin returns paginated tools for admin view (optionally including archived)
func (r *repository) ListToolsAdmin(search string, includeArchived bool, page, pageSize int) ([]Tool, int64, error) {
	var tools []Tool
	var total int64

	query := r.db.Model(&domain.Tool{})

	// Exclude archived unless explicitly requested
	if !includeArchived {
		query = query.Where("archived_at IS NULL")
	}

	// Search filter
	if search != "" {
		searchTerm := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(slug) LIKE ?", searchTerm, searchTerm)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	err := query.
		Preload("PrimaryCategory").
		Preload("Tags").
		Preload("Badges").
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&tools).Error

	if err != nil {
		return nil, 0, err
	}

	return tools, total, nil
}

// GetToolByIDAdmin finds a tool by ID (including archived)
func (r *repository) GetToolByIDAdmin(id uint) (*Tool, error) {
	var tool Tool
	err := r.db.
		Preload("PrimaryCategory").
		Preload("Tags").
		Preload("Badges").
		Preload("Media").
		Where("id = ?", id).
		First(&tool).Error
	if err != nil {
		return nil, err
	}
	return &tool, nil
}

// Create inserts a new tool
func (r *repository) Create(tool *Tool) error {
	return r.db.Create(tool).Error
}

// Update updates an existing tool
func (r *repository) Update(tool *Tool) error {
	return r.db.Save(tool).Error
}

// Archive soft-deletes a tool by setting archived_at
func (r *repository) Archive(id uint) error {
	return r.db.Model(&domain.Tool{}).
		Where("id = ?", id).
		Update("archived_at", gorm.Expr("NOW()")).Error
}

// SlugExists checks if a slug is already in use (excluding the given ID for updates)
func (r *repository) SlugExists(slug string, excludeID uint) (bool, error) {
	var count int64
	query := r.db.Model(&domain.Tool{}).Where("slug = ?", slug)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
