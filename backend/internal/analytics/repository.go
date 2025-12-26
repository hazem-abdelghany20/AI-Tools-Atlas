package analytics

import (
	"time"

	"gorm.io/gorm"
)

// OverviewStats contains overall system statistics
type OverviewStats struct {
	TotalTools      int64 `json:"total_tools"`
	TotalCategories int64 `json:"total_categories"`
	TotalReviews    int64 `json:"total_reviews"`
	TotalBookmarks  int64 `json:"total_bookmarks"`
	TotalUsers      int64 `json:"total_users"`
	NewToolsWeek    int64 `json:"new_tools_week"`
	NewToolsMonth   int64 `json:"new_tools_month"`
	NewReviewsWeek  int64 `json:"new_reviews_week"`
	NewUsersWeek    int64 `json:"new_users_week"`
}

// TopTool represents a tool with its ranking metrics
type TopTool struct {
	ID            uint    `json:"id"`
	Slug          string  `json:"slug"`
	Name          string  `json:"name"`
	LogoURL       string  `json:"logo_url,omitempty"`
	BookmarkCount int     `json:"bookmark_count"`
	ReviewCount   int     `json:"review_count"`
	AvgRating     float64 `json:"avg_rating"`
}

// TopCategory represents a category with its ranking metrics
type TopCategory struct {
	ID        uint   `json:"id"`
	Slug      string `json:"slug"`
	Name      string `json:"name"`
	ToolCount int64  `json:"tool_count"`
}

// Repository defines the interface for analytics data operations
type Repository interface {
	GetOverviewStats() (*OverviewStats, error)
	GetTopToolsByBookmarks(limit int) ([]TopTool, error)
	GetTopToolsByRating(limit int) ([]TopTool, error)
	GetTopToolsByReviews(limit int) ([]TopTool, error)
	GetTopCategories(limit int) ([]TopCategory, error)
}

// repository implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new analytics repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// GetOverviewStats returns aggregate statistics
func (r *repository) GetOverviewStats() (*OverviewStats, error) {
	stats := &OverviewStats{}

	// Total counts
	r.db.Table("tools").Where("archived_at IS NULL").Count(&stats.TotalTools)
	r.db.Table("categories").Count(&stats.TotalCategories)
	r.db.Table("reviews").Count(&stats.TotalReviews)
	r.db.Table("bookmarks").Count(&stats.TotalBookmarks)
	r.db.Table("users").Count(&stats.TotalUsers)

	// New this week
	weekAgo := time.Now().AddDate(0, 0, -7)
	monthAgo := time.Now().AddDate(0, -1, 0)

	r.db.Table("tools").Where("created_at >= ? AND archived_at IS NULL", weekAgo).Count(&stats.NewToolsWeek)
	r.db.Table("tools").Where("created_at >= ? AND archived_at IS NULL", monthAgo).Count(&stats.NewToolsMonth)
	r.db.Table("reviews").Where("created_at >= ?", weekAgo).Count(&stats.NewReviewsWeek)
	r.db.Table("users").Where("created_at >= ?", weekAgo).Count(&stats.NewUsersWeek)

	return stats, nil
}

// GetTopToolsByBookmarks returns top tools by bookmark count
func (r *repository) GetTopToolsByBookmarks(limit int) ([]TopTool, error) {
	var tools []TopTool
	err := r.db.Raw(`
		SELECT t.id, t.slug, t.name, t.logo_url, t.bookmark_count, t.review_count, t.avg_rating_overall as avg_rating
		FROM tools t
		WHERE t.archived_at IS NULL
		ORDER BY t.bookmark_count DESC
		LIMIT ?
	`, limit).Scan(&tools).Error
	return tools, err
}

// GetTopToolsByRating returns top tools by rating
func (r *repository) GetTopToolsByRating(limit int) ([]TopTool, error) {
	var tools []TopTool
	err := r.db.Raw(`
		SELECT t.id, t.slug, t.name, t.logo_url, t.bookmark_count, t.review_count, t.avg_rating_overall as avg_rating
		FROM tools t
		WHERE t.archived_at IS NULL AND t.review_count >= 1
		ORDER BY t.avg_rating_overall DESC, t.review_count DESC
		LIMIT ?
	`, limit).Scan(&tools).Error
	return tools, err
}

// GetTopToolsByReviews returns top tools by review count
func (r *repository) GetTopToolsByReviews(limit int) ([]TopTool, error) {
	var tools []TopTool
	err := r.db.Raw(`
		SELECT t.id, t.slug, t.name, t.logo_url, t.bookmark_count, t.review_count, t.avg_rating_overall as avg_rating
		FROM tools t
		WHERE t.archived_at IS NULL
		ORDER BY t.review_count DESC
		LIMIT ?
	`, limit).Scan(&tools).Error
	return tools, err
}

// GetTopCategories returns top categories by tool count
func (r *repository) GetTopCategories(limit int) ([]TopCategory, error) {
	var categories []TopCategory
	err := r.db.Raw(`
		SELECT c.id, c.slug, c.name, COUNT(t.id) as tool_count
		FROM categories c
		LEFT JOIN tools t ON t.primary_category_id = c.id AND t.archived_at IS NULL
		GROUP BY c.id, c.slug, c.name
		ORDER BY tool_count DESC
		LIMIT ?
	`, limit).Scan(&categories).Error
	return categories, err
}
