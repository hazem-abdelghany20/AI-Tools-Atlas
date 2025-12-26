package reviews

import (
	"github.com/your-org/ai-tools-atlas-backend/internal/domain"
	"gorm.io/gorm"
)

// Sort constants for reviews
const (
	SortNewest      = "newest"
	SortMostHelpful = "most_helpful"
	SortHighest     = "highest"
	SortLowest      = "lowest"
)

// Repository defines the interface for review data operations
type Repository interface {
	ListReviewsByTool(toolID uint, sort string, page, pageSize int) ([]domain.Review, int64, error)
	ListReviewsByUser(userID uint, page, pageSize int) ([]domain.Review, int64, error)
	CreateReview(review *domain.Review) error
	HasUserReviewed(toolID, userID uint) (bool, error)
	UpdateToolRatingAggregates(toolID uint) error
	GetToolBySlug(slug string) (*domain.Tool, error)
}

// repository implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new reviews repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// ListReviewsByTool returns paginated reviews for a tool
func (r *repository) ListReviewsByTool(toolID uint, sort string, page, pageSize int) ([]domain.Review, int64, error) {
	var reviews []domain.Review
	var total int64

	// Base query: only approved reviews
	query := r.db.Model(&domain.Review{}).
		Where("tool_id = ? AND moderation_status = ?", toolID, "approved")

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	query = r.applySorting(query, sort)

	// Apply pagination
	offset := (page - 1) * pageSize
	query = query.
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "display_name")
		}).
		Limit(pageSize).
		Offset(offset)

	if err := query.Find(&reviews).Error; err != nil {
		return nil, 0, err
	}

	return reviews, total, nil
}

// ListReviewsByUser returns paginated reviews by a user with tool info
func (r *repository) ListReviewsByUser(userID uint, page, pageSize int) ([]domain.Review, int64, error) {
	var reviews []domain.Review
	var total int64

	// Base query: all reviews by user
	query := r.db.Model(&domain.Review{}).Where("user_id = ?", userID)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination with tool preload
	offset := (page - 1) * pageSize
	err := r.db.Where("user_id = ?", userID).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "display_name")
		}).
		Preload("Tool", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "slug", "name", "logo_url")
		}).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&reviews).Error

	if err != nil {
		return nil, 0, err
	}

	return reviews, total, nil
}

// CreateReview inserts a new review
func (r *repository) CreateReview(review *domain.Review) error {
	return r.db.Create(review).Error
}

// HasUserReviewed checks if a user has already reviewed a tool
func (r *repository) HasUserReviewed(toolID, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&domain.Review{}).
		Where("tool_id = ? AND user_id = ?", toolID, userID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// UpdateToolRatingAggregates recalculates avg_rating_overall and review_count for a tool
func (r *repository) UpdateToolRatingAggregates(toolID uint) error {
	// Calculate aggregates from approved reviews
	var result struct {
		AvgRating   float64
		ReviewCount int64
	}

	err := r.db.Model(&domain.Review{}).
		Select("COALESCE(AVG(rating_overall), 0) as avg_rating, COUNT(*) as review_count").
		Where("tool_id = ? AND moderation_status = ?", toolID, "approved").
		Scan(&result).Error
	if err != nil {
		return err
	}

	// Update tool
	return r.db.Model(&domain.Tool{}).
		Where("id = ?", toolID).
		Updates(map[string]interface{}{
			"avg_rating_overall": result.AvgRating,
			"review_count":       result.ReviewCount,
		}).Error
}

// GetToolBySlug finds a tool by its slug
func (r *repository) GetToolBySlug(slug string) (*domain.Tool, error) {
	var tool domain.Tool
	err := r.db.Where("slug = ? AND archived_at IS NULL", slug).First(&tool).Error
	if err != nil {
		return nil, err
	}
	return &tool, nil
}

// applySorting adds the appropriate ORDER BY clause
func (r *repository) applySorting(query *gorm.DB, sort string) *gorm.DB {
	switch sort {
	case SortMostHelpful:
		return query.Order("helpful_count DESC, created_at DESC")
	case SortHighest:
		return query.Order("rating_overall DESC, created_at DESC")
	case SortLowest:
		return query.Order("rating_overall ASC, created_at DESC")
	case SortNewest:
		fallthrough
	default:
		return query.Order("created_at DESC")
	}
}

// ValidateSort ensures sort value is valid
func ValidateSort(sort string) string {
	switch sort {
	case SortNewest, SortMostHelpful, SortHighest, SortLowest:
		return sort
	default:
		return SortNewest
	}
}
