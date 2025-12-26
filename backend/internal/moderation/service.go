package moderation

import (
	"errors"

	"github.com/your-org/ai-tools-atlas-backend/internal/domain"
	"gorm.io/gorm"
)

var (
	ErrReportNotFound        = errors.New("report not found")
	ErrReviewNotFound        = errors.New("review not found")
	ErrInvalidReason         = errors.New("invalid reason")
	ErrAlreadyReported       = errors.New("you have already reported this item today")
	ErrInvalidStatus         = errors.New("invalid status")
	ErrInvalidModerationStatus = errors.New("invalid moderation status")
)

// Valid reasons for reporting
var validReasons = map[string]bool{
	"spam":           true,
	"abuse":          true,
	"misinformation": true,
	"other":          true,
}

// Valid statuses for reports
var validStatuses = map[string]bool{
	"pending":   true,
	"reviewed":  true,
	"dismissed": true,
}

// Valid moderation statuses for reviews
var validModerationStatuses = map[string]bool{
	"pending":  true,
	"approved": true,
	"rejected": true,
	"hidden":   true,
	"removed":  true,
}

// CreateReportInput represents input for creating a report
type CreateReportInput struct {
	Reason  string `json:"reason" binding:"required"`
	Comment string `json:"comment,omitempty"`
}

// ModerationActionInput represents input for moderation actions
type ModerationActionInput struct {
	Notes string `json:"notes,omitempty"`
}

// ReportWithContext includes the reportable object
type ReportWithContext struct {
	Report
	Tool   *domain.Tool   `json:"tool,omitempty"`
	Review *domain.Review `json:"review,omitempty"`
}

// ModerationActionResponse represents an action in the history
type ModerationActionResponse struct {
	ID          uint   `json:"id"`
	ActionType  string `json:"action_type"`
	Notes       string `json:"notes,omitempty"`
	CreatedAt   string `json:"created_at"`
	Moderator   struct {
		ID          uint   `json:"id"`
		DisplayName string `json:"display_name"`
	} `json:"moderator"`
}

// Service defines the interface for moderation business logic
type Service interface {
	CreateToolReport(toolID uint, userID *uint, input CreateReportInput) (*Report, error)
	CreateReviewReport(reviewID uint, userID *uint, input CreateReportInput) (*Report, error)
	GetReportByID(id uint) (*Report, error)
	ListPendingReports(page, pageSize int) ([]Report, int64, error)
	ListReports(filters ReportFilters, page, pageSize int) ([]ReportWithContext, int64, error)
	UpdateReportStatus(id uint, status string, reviewedBy uint) error

	// Review moderation
	ApproveReview(reviewID uint, moderatorID uint, input ModerationActionInput) (*domain.Review, error)
	HideReview(reviewID uint, moderatorID uint, input ModerationActionInput) (*domain.Review, error)
	RemoveReview(reviewID uint, moderatorID uint, input ModerationActionInput) (*domain.Review, error)
	GetModerationHistory(reviewID uint) ([]ModerationActionResponse, error)
}

// ReviewsRepository is a subset of reviews.Repository needed for rating updates
type ReviewsRepository interface {
	UpdateToolRatingAggregates(toolID uint) error
}

// service implements the Service interface
type service struct {
	repo        Repository
	reviewsRepo ReviewsRepository
}

// NewService creates a new moderation service
func NewService(repo Repository, reviewsRepo ReviewsRepository) Service {
	return &service{repo: repo, reviewsRepo: reviewsRepo}
}

// CreateToolReport creates a report for a tool
func (s *service) CreateToolReport(toolID uint, userID *uint, input CreateReportInput) (*Report, error) {
	return s.createReport("tool", toolID, userID, input)
}

// CreateReviewReport creates a report for a review
func (s *service) CreateReviewReport(reviewID uint, userID *uint, input CreateReportInput) (*Report, error) {
	return s.createReport("review", reviewID, userID, input)
}

func (s *service) createReport(reportableType string, reportableID uint, userID *uint, input CreateReportInput) (*Report, error) {
	// Validate reason
	if !validReasons[input.Reason] {
		return nil, ErrInvalidReason
	}

	// Check for duplicate reports today
	if userID != nil {
		count, err := s.repo.CountUserReportsToday(*userID, reportableType, reportableID)
		if err != nil {
			return nil, err
		}
		if count > 0 {
			return nil, ErrAlreadyReported
		}
	}

	report := &Report{
		ReportableType: reportableType,
		ReportableID:   reportableID,
		ReporterUserID: userID,
		Reason:         input.Reason,
		Comment:        input.Comment,
		Status:         "pending",
	}

	if err := s.repo.CreateReport(report); err != nil {
		return nil, err
	}

	return report, nil
}

// GetReportByID finds a report by ID
func (s *service) GetReportByID(id uint) (*Report, error) {
	report, err := s.repo.GetReportByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrReportNotFound
		}
		return nil, err
	}
	return report, nil
}

// ListPendingReports returns pending reports with pagination
func (s *service) ListPendingReports(page, pageSize int) ([]Report, int64, error) {
	page, pageSize = s.validatePagination(page, pageSize)
	return s.repo.ListPendingReports(page, pageSize)
}

// ListReports returns reports with filters and includes reportable objects
func (s *service) ListReports(filters ReportFilters, page, pageSize int) ([]ReportWithContext, int64, error) {
	page, pageSize = s.validatePagination(page, pageSize)

	reports, total, err := s.repo.ListReports(filters, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// Enrich reports with reportable objects
	result := make([]ReportWithContext, len(reports))
	for i, report := range reports {
		result[i] = ReportWithContext{Report: report}

		switch report.ReportableType {
		case "tool":
			tool, err := s.repo.GetToolByID(report.ReportableID)
			if err == nil {
				result[i].Tool = tool
			}
		case "review":
			review, err := s.repo.GetReviewByID(report.ReportableID)
			if err == nil {
				result[i].Review = review
			}
		}
	}

	return result, total, nil
}

// UpdateReportStatus updates a report's status
func (s *service) UpdateReportStatus(id uint, status string, reviewedBy uint) error {
	if !validStatuses[status] {
		return ErrInvalidStatus
	}

	// Check report exists
	_, err := s.repo.GetReportByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrReportNotFound
		}
		return err
	}

	return s.repo.UpdateReportStatus(id, status, reviewedBy)
}

// ApproveReview approves a review for public display
func (s *service) ApproveReview(reviewID uint, moderatorID uint, input ModerationActionInput) (*domain.Review, error) {
	return s.moderateReview(reviewID, moderatorID, "approved", "approve", input)
}

// HideReview hides a review from public display
func (s *service) HideReview(reviewID uint, moderatorID uint, input ModerationActionInput) (*domain.Review, error) {
	return s.moderateReview(reviewID, moderatorID, "hidden", "hide", input)
}

// RemoveReview permanently hides a review and updates rating aggregates
func (s *service) RemoveReview(reviewID uint, moderatorID uint, input ModerationActionInput) (*domain.Review, error) {
	return s.moderateReview(reviewID, moderatorID, "removed", "remove", input)
}

func (s *service) moderateReview(reviewID uint, moderatorID uint, status, actionType string, input ModerationActionInput) (*domain.Review, error) {
	// Get review first
	review, err := s.repo.GetReviewByID(reviewID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrReviewNotFound
		}
		return nil, err
	}

	// Update moderation status
	err = s.repo.UpdateReviewModerationStatus(reviewID, status, moderatorID)
	if err != nil {
		return nil, err
	}

	// Log moderation action
	action := &domain.ModerationAction{
		ReviewID:    reviewID,
		ModeratorID: moderatorID,
		ActionType:  actionType,
		Notes:       input.Notes,
	}
	_ = s.repo.CreateModerationAction(action) // Log error but don't fail

	// Update tool rating aggregates if review is hidden/removed or restored to approved
	if s.reviewsRepo != nil && (status == "hidden" || status == "removed" || status == "approved") {
		_ = s.reviewsRepo.UpdateToolRatingAggregates(review.ToolID) // Log error but don't fail
	}

	// Refresh review data
	review, err = s.repo.GetReviewByID(reviewID)
	if err != nil {
		return nil, err
	}

	return review, nil
}

// GetModerationHistory returns audit log of moderation actions for a review
func (s *service) GetModerationHistory(reviewID uint) ([]ModerationActionResponse, error) {
	// Verify review exists
	_, err := s.repo.GetReviewByID(reviewID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrReviewNotFound
		}
		return nil, err
	}

	actions, err := s.repo.GetModerationHistory(reviewID)
	if err != nil {
		return nil, err
	}

	// Convert to response DTOs
	result := make([]ModerationActionResponse, len(actions))
	for i, action := range actions {
		result[i] = ModerationActionResponse{
			ID:         action.ID,
			ActionType: action.ActionType,
			Notes:      action.Notes,
			CreatedAt:  action.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
		result[i].Moderator.ID = action.ModeratorID
		if action.Moderator.DisplayName != "" {
			result[i].Moderator.DisplayName = action.Moderator.DisplayName
		}
	}

	return result, nil
}

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
