package moderation

import (
	"time"

	"github.com/your-org/ai-tools-atlas-backend/internal/domain"
	"gorm.io/gorm"
)

// ReportFilters defines filters for listing reports
type ReportFilters struct {
	Type   string // "tool" or "review"
	Status string // "pending", "reviewed", "dismissed"
}

// Repository defines the interface for report data operations
type Repository interface {
	CreateReport(report *Report) error
	GetReportByID(id uint) (*Report, error)
	ListPendingReports(page, pageSize int) ([]Report, int64, error)
	ListReports(filters ReportFilters, page, pageSize int) ([]Report, int64, error)
	UpdateReportStatus(id uint, status string, reviewedBy uint) error
	CountUserReportsToday(userID uint, reportableType string, reportableID uint) (int64, error)
	CountAnonReportsToday(reportableType string, reportableID uint) (int64, error)

	// Review moderation
	GetReviewByID(id uint) (*domain.Review, error)
	UpdateReviewModerationStatus(id uint, status string, moderatorID uint) error

	// Moderation actions (audit log)
	CreateModerationAction(action *domain.ModerationAction) error
	GetModerationHistory(reviewID uint) ([]domain.ModerationAction, error)

	// Get reportable objects
	GetToolByID(id uint) (*domain.Tool, error)
}

// repository implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new moderation repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// CreateReport creates a new report
func (r *repository) CreateReport(report *Report) error {
	return r.db.Create(report).Error
}

// GetReportByID finds a report by ID
func (r *repository) GetReportByID(id uint) (*Report, error) {
	var report Report
	err := r.db.Preload("ReporterUser").Where("id = ?", id).First(&report).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}

// ListPendingReports returns reports with pending status
func (r *repository) ListPendingReports(page, pageSize int) ([]Report, int64, error) {
	return r.ListReports(ReportFilters{Status: "pending"}, page, pageSize)
}

// ListReports returns reports with optional filters
func (r *repository) ListReports(filters ReportFilters, page, pageSize int) ([]Report, int64, error) {
	var reports []Report
	var total int64

	offset := (page - 1) * pageSize

	query := r.db.Model(&Report{})

	// Apply filters
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	if filters.Type != "" {
		query = query.Where("reportable_type = ?", filters.Type)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.
		Preload("ReporterUser").
		Where(r.buildWhereClause(filters)).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&reports).Error

	if err != nil {
		return nil, 0, err
	}

	return reports, total, nil
}

func (r *repository) buildWhereClause(filters ReportFilters) map[string]interface{} {
	where := make(map[string]interface{})
	if filters.Status != "" {
		where["status"] = filters.Status
	}
	if filters.Type != "" {
		where["reportable_type"] = filters.Type
	}
	return where
}

// UpdateReportStatus updates a report's status
func (r *repository) UpdateReportStatus(id uint, status string, reviewedBy uint) error {
	now := time.Now()
	return r.db.Model(&Report{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":      status,
		"reviewed_by": reviewedBy,
		"reviewed_at": now,
	}).Error
}

// CountUserReportsToday counts reports from a user for a specific item today
func (r *repository) CountUserReportsToday(userID uint, reportableType string, reportableID uint) (int64, error) {
	var count int64
	today := time.Now().Truncate(24 * time.Hour)
	err := r.db.Model(&Report{}).
		Where("reporter_user_id = ? AND reportable_type = ? AND reportable_id = ? AND created_at >= ?",
			userID, reportableType, reportableID, today).
		Count(&count).Error
	return count, err
}

// CountAnonReportsToday counts anonymous reports for a specific item today
func (r *repository) CountAnonReportsToday(reportableType string, reportableID uint) (int64, error) {
	var count int64
	today := time.Now().Truncate(24 * time.Hour)
	err := r.db.Model(&Report{}).
		Where("reporter_user_id IS NULL AND reportable_type = ? AND reportable_id = ? AND created_at >= ?",
			reportableType, reportableID, today).
		Count(&count).Error
	return count, err
}

// GetReviewByID finds a review by ID
func (r *repository) GetReviewByID(id uint) (*domain.Review, error) {
	var review domain.Review
	err := r.db.
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "display_name", "email")
		}).
		Preload("Tool", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "slug", "name", "logo_url")
		}).
		Where("id = ?", id).First(&review).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

// UpdateReviewModerationStatus updates a review's moderation status
func (r *repository) UpdateReviewModerationStatus(id uint, status string, moderatorID uint) error {
	now := time.Now()
	return r.db.Model(&domain.Review{}).Where("id = ?", id).Updates(map[string]interface{}{
		"moderation_status": status,
		"moderated_by":      moderatorID,
		"moderated_at":      now,
	}).Error
}

// CreateModerationAction creates an audit log entry
func (r *repository) CreateModerationAction(action *domain.ModerationAction) error {
	return r.db.Create(action).Error
}

// GetModerationHistory returns moderation actions for a review
func (r *repository) GetModerationHistory(reviewID uint) ([]domain.ModerationAction, error) {
	var actions []domain.ModerationAction
	err := r.db.
		Preload("Moderator", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "display_name", "email")
		}).
		Where("review_id = ?", reviewID).
		Order("created_at DESC").
		Find(&actions).Error
	if err != nil {
		return nil, err
	}
	return actions, nil
}

// GetToolByID finds a tool by ID
func (r *repository) GetToolByID(id uint) (*domain.Tool, error) {
	var tool domain.Tool
	err := r.db.
		Preload("PrimaryCategory").
		Where("id = ?", id).First(&tool).Error
	if err != nil {
		return nil, err
	}
	return &tool, nil
}
