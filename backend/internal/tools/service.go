package tools

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrToolNotFound      = errors.New("tool not found")
	ErrSlugRequired      = errors.New("slug is required")
	ErrNameRequired      = errors.New("name is required")
	ErrSlugExists        = errors.New("slug already exists")
	ErrCategoryRequired  = errors.New("primary_category_id is required")
	ErrInvalidCategoryID = errors.New("invalid category ID")
)

// CreateToolInput represents input for creating a new tool
type CreateToolInput struct {
	Slug              string  `json:"slug"`
	Name              string  `json:"name"`
	LogoURL           string  `json:"logo_url,omitempty"`
	Tagline           string  `json:"tagline,omitempty"`
	Description       string  `json:"description,omitempty"`
	BestFor           string  `json:"best_for,omitempty"`
	PrimaryUseCases   string  `json:"primary_use_cases,omitempty"`
	PricingSummary    string  `json:"pricing_summary,omitempty"`
	TargetRoles       string  `json:"target_roles,omitempty"`
	Platforms         string  `json:"platforms,omitempty"`
	HasFreeTier       bool    `json:"has_free_tier"`
	OfficialURL       string  `json:"official_url,omitempty"`
	PrimaryCategoryID uint    `json:"primary_category_id"`
}

// UpdateToolInput represents input for updating an existing tool
type UpdateToolInput struct {
	Name              *string `json:"name,omitempty"`
	LogoURL           *string `json:"logo_url,omitempty"`
	Tagline           *string `json:"tagline,omitempty"`
	Description       *string `json:"description,omitempty"`
	BestFor           *string `json:"best_for,omitempty"`
	PrimaryUseCases   *string `json:"primary_use_cases,omitempty"`
	PricingSummary    *string `json:"pricing_summary,omitempty"`
	TargetRoles       *string `json:"target_roles,omitempty"`
	Platforms         *string `json:"platforms,omitempty"`
	HasFreeTier       *bool   `json:"has_free_tier,omitempty"`
	OfficialURL       *string `json:"official_url,omitempty"`
	PrimaryCategoryID *uint   `json:"primary_category_id,omitempty"`
}

// Service defines the interface for tool business logic
type Service interface {
	ListTools(filters ToolFilters, page, pageSize int) ([]Tool, int64, error)
	SearchTools(query string, filters ToolFilters, page, pageSize int) ([]Tool, int64, error)
	GetToolBySlug(slug string) (*Tool, error)
	GetToolByID(id uint) (*Tool, error)
	GetToolAlternatives(slug string) (*AlternativesResult, error)
	// Admin methods
	ListToolsAdmin(search string, includeArchived bool, page, pageSize int) ([]Tool, int64, error)
	GetToolByIDAdmin(id uint) (*Tool, error)
	CreateTool(input CreateToolInput) (*Tool, error)
	UpdateTool(id uint, input UpdateToolInput) (*Tool, error)
	ArchiveTool(id uint) error
}

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new tool service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// ListTools returns paginated tools with filters
func (s *service) ListTools(filters ToolFilters, page, pageSize int) ([]Tool, int64, error) {
	// Validate and set defaults
	page, pageSize = s.validatePagination(page, pageSize)
	filters.Sort = ValidateSort(filters.Sort)
	filters.Price = ValidatePrice(filters.Price)

	return s.repo.ListTools(filters, page, pageSize)
}

// SearchTools performs a free-text search with filters
func (s *service) SearchTools(query string, filters ToolFilters, page, pageSize int) ([]Tool, int64, error) {
	// Validate and set defaults
	page, pageSize = s.validatePagination(page, pageSize)
	filters.Sort = ValidateSort(filters.Sort)
	filters.Price = ValidatePrice(filters.Price)

	return s.repo.SearchTools(query, filters, page, pageSize)
}

// GetToolBySlug finds a tool by its slug
func (s *service) GetToolBySlug(slug string) (*Tool, error) {
	tool, err := s.repo.GetToolBySlug(slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrToolNotFound
		}
		return nil, err
	}
	return tool, nil
}

// GetToolByID finds a tool by its ID
func (s *service) GetToolByID(id uint) (*Tool, error) {
	tool, err := s.repo.GetToolByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrToolNotFound
		}
		return nil, err
	}
	return tool, nil
}

// GetToolAlternatives returns similar and alternative tools for a given tool slug
func (s *service) GetToolAlternatives(slug string) (*AlternativesResult, error) {
	// First get the tool to find its ID
	tool, err := s.GetToolBySlug(slug)
	if err != nil {
		return nil, err
	}

	return s.repo.GetToolAlternatives(tool.ID, 6)
}

// validatePagination ensures page and pageSize have valid values
func (s *service) validatePagination(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

// ListToolsAdmin returns paginated tools for admin view
func (s *service) ListToolsAdmin(search string, includeArchived bool, page, pageSize int) ([]Tool, int64, error) {
	page, pageSize = s.validatePagination(page, pageSize)
	return s.repo.ListToolsAdmin(search, includeArchived, page, pageSize)
}

// GetToolByIDAdmin finds a tool by ID (including archived)
func (s *service) GetToolByIDAdmin(id uint) (*Tool, error) {
	tool, err := s.repo.GetToolByIDAdmin(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrToolNotFound
		}
		return nil, err
	}
	return tool, nil
}

// CreateTool creates a new tool
func (s *service) CreateTool(input CreateToolInput) (*Tool, error) {
	// Validate required fields
	if input.Slug == "" {
		return nil, ErrSlugRequired
	}
	if input.Name == "" {
		return nil, ErrNameRequired
	}
	if input.PrimaryCategoryID == 0 {
		return nil, ErrCategoryRequired
	}

	// Check if slug already exists
	exists, err := s.repo.SlugExists(input.Slug, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrSlugExists
	}

	tool := &Tool{
		Slug:              input.Slug,
		Name:              input.Name,
		LogoURL:           input.LogoURL,
		Tagline:           input.Tagline,
		Description:       input.Description,
		BestFor:           input.BestFor,
		PrimaryUseCases:   input.PrimaryUseCases,
		PricingSummary:    input.PricingSummary,
		TargetRoles:       input.TargetRoles,
		Platforms:         input.Platforms,
		HasFreeTier:       input.HasFreeTier,
		OfficialURL:       input.OfficialURL,
		PrimaryCategoryID: input.PrimaryCategoryID,
	}

	if err := s.repo.Create(tool); err != nil {
		return nil, err
	}

	// Fetch the complete tool with relations
	return s.repo.GetToolByIDAdmin(tool.ID)
}

// UpdateTool updates an existing tool
func (s *service) UpdateTool(id uint, input UpdateToolInput) (*Tool, error) {
	// Get existing tool
	tool, err := s.repo.GetToolByIDAdmin(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrToolNotFound
		}
		return nil, err
	}

	// Apply updates
	if input.Name != nil {
		if *input.Name == "" {
			return nil, ErrNameRequired
		}
		tool.Name = *input.Name
	}
	if input.LogoURL != nil {
		tool.LogoURL = *input.LogoURL
	}
	if input.Tagline != nil {
		tool.Tagline = *input.Tagline
	}
	if input.Description != nil {
		tool.Description = *input.Description
	}
	if input.BestFor != nil {
		tool.BestFor = *input.BestFor
	}
	if input.PrimaryUseCases != nil {
		tool.PrimaryUseCases = *input.PrimaryUseCases
	}
	if input.PricingSummary != nil {
		tool.PricingSummary = *input.PricingSummary
	}
	if input.TargetRoles != nil {
		tool.TargetRoles = *input.TargetRoles
	}
	if input.Platforms != nil {
		tool.Platforms = *input.Platforms
	}
	if input.HasFreeTier != nil {
		tool.HasFreeTier = *input.HasFreeTier
	}
	if input.OfficialURL != nil {
		tool.OfficialURL = *input.OfficialURL
	}
	if input.PrimaryCategoryID != nil {
		if *input.PrimaryCategoryID == 0 {
			return nil, ErrCategoryRequired
		}
		tool.PrimaryCategoryID = *input.PrimaryCategoryID
	}

	if err := s.repo.Update(tool); err != nil {
		return nil, err
	}

	// Fetch the complete tool with relations
	return s.repo.GetToolByIDAdmin(tool.ID)
}

// ArchiveTool soft-deletes a tool
func (s *service) ArchiveTool(id uint) error {
	// Verify tool exists
	_, err := s.repo.GetToolByIDAdmin(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrToolNotFound
		}
		return err
	}

	return s.repo.Archive(id)
}
