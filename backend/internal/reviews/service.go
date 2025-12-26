package reviews

import (
	"errors"
	"time"
	"unicode/utf8"

	"github.com/your-org/ai-tools-atlas-backend/internal/domain"
	"gorm.io/gorm"
)

// Error constants
var (
	ErrToolNotFound      = errors.New("tool not found")
	ErrAlreadyReviewed   = errors.New("user has already reviewed this tool")
	ErrInvalidRating     = errors.New("rating must be between 1 and 5")
	ErrProsRequired      = errors.New("pros field is required")
	ErrConsRequired      = errors.New("cons field is required")
	ErrProsTooLong       = errors.New("pros must be 500 characters or less")
	ErrConsTooLong       = errors.New("cons must be 500 characters or less")
	ErrRatingRequired    = errors.New("rating_overall is required")
)

// CreateReviewInput represents the input for creating a review
type CreateReviewInput struct {
	RatingOverall   int    `json:"rating_overall"`
	RatingEaseOfUse *int   `json:"rating_ease_of_use,omitempty"`
	RatingValue     *int   `json:"rating_value,omitempty"`
	RatingAccuracy  *int   `json:"rating_accuracy,omitempty"`
	RatingSpeed     *int   `json:"rating_speed,omitempty"`
	RatingSupport   *int   `json:"rating_support,omitempty"`
	Pros            string `json:"pros"`
	Cons            string `json:"cons"`
	PrimaryUseCase  string `json:"primary_use_case,omitempty"`
	ReviewerRole    string `json:"reviewer_role,omitempty"`
	CompanySize     string `json:"company_size,omitempty"`
	UsageContext    string `json:"usage_context,omitempty"`
}

// ReviewResponse represents a review with user info for API response
type ReviewResponse struct {
	ID               uint       `json:"id"`
	RatingOverall    int        `json:"rating_overall"`
	RatingEaseOfUse  *int       `json:"rating_ease_of_use,omitempty"`
	RatingValue      *int       `json:"rating_value,omitempty"`
	RatingAccuracy   *int       `json:"rating_accuracy,omitempty"`
	RatingSpeed      *int       `json:"rating_speed,omitempty"`
	RatingSupport    *int       `json:"rating_support,omitempty"`
	Pros             string     `json:"pros,omitempty"`
	Cons             string     `json:"cons,omitempty"`
	PrimaryUseCase   string     `json:"primary_use_case,omitempty"`
	ReviewerRole     string     `json:"reviewer_role,omitempty"`
	CompanySize      string     `json:"company_size,omitempty"`
	UsageContext     string     `json:"usage_context,omitempty"`
	HelpfulCount     int        `json:"helpful_count"`
	CreatedAt        string     `json:"created_at"`
	User             UserBrief  `json:"user"`
}

// UserBrief represents minimal user info for review response
type UserBrief struct {
	ID          uint   `json:"id"`
	DisplayName string `json:"display_name"`
}

// UserReviewResponse represents a review with tool info for user profile
type UserReviewResponse struct {
	ID               uint   `json:"id"`
	RatingOverall    int    `json:"rating_overall"`
	Pros             string `json:"pros,omitempty"`
	Cons             string `json:"cons,omitempty"`
	HelpfulCount     int    `json:"helpful_count"`
	ModerationStatus string `json:"moderation_status"`
	CreatedAt        string `json:"created_at"`
	Tool             struct {
		Slug    string `json:"slug"`
		Name    string `json:"name"`
		LogoURL string `json:"logo_url,omitempty"`
	} `json:"tool"`
}

// Service defines the interface for review business logic
type Service interface {
	ListReviews(slug string, sort string, page, pageSize int) ([]ReviewResponse, int64, error)
	ListUserReviews(userID uint, page, pageSize int) ([]UserReviewResponse, int64, error)
	CreateReview(slug string, userID uint, input CreateReviewInput) (*ReviewResponse, error)
}

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new review service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// ListReviews returns paginated reviews for a tool
func (s *service) ListReviews(slug string, sort string, page, pageSize int) ([]ReviewResponse, int64, error) {
	// Validate and set defaults
	page, pageSize = s.validatePagination(page, pageSize)
	sort = ValidateSort(sort)

	// Get tool by slug
	tool, err := s.repo.GetToolBySlug(slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, ErrToolNotFound
		}
		return nil, 0, err
	}

	// Get reviews
	reviews, total, err := s.repo.ListReviewsByTool(tool.ID, sort, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// Convert to response DTOs
	responses := make([]ReviewResponse, len(reviews))
	for i, r := range reviews {
		responses[i] = s.toReviewResponse(r)
	}

	return responses, total, nil
}

// ListUserReviews returns paginated reviews by a user with tool info
func (s *service) ListUserReviews(userID uint, page, pageSize int) ([]UserReviewResponse, int64, error) {
	// Validate pagination
	page, pageSize = s.validatePagination(page, pageSize)

	// Get reviews with tool preload
	reviews, total, err := s.repo.ListReviewsByUser(userID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// Convert to response DTOs
	responses := make([]UserReviewResponse, len(reviews))
	for i, r := range reviews {
		responses[i] = s.toUserReviewResponse(r)
	}

	return responses, total, nil
}

// toUserReviewResponse converts a domain review to user review API response
func (s *service) toUserReviewResponse(r domain.Review) UserReviewResponse {
	resp := UserReviewResponse{
		ID:               r.ID,
		RatingOverall:    r.RatingOverall,
		Pros:             r.Pros,
		Cons:             r.Cons,
		HelpfulCount:     r.HelpfulCount,
		ModerationStatus: r.ModerationStatus,
		CreatedAt:        r.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	// Tool info is loaded via Preload
	resp.Tool.Slug = r.Tool.Slug
	resp.Tool.Name = r.Tool.Name
	resp.Tool.LogoURL = r.Tool.LogoURL

	return resp
}

// CreateReview creates a new review for a tool
func (s *service) CreateReview(slug string, userID uint, input CreateReviewInput) (*ReviewResponse, error) {
	// Validate input
	if err := s.validateInput(input); err != nil {
		return nil, err
	}

	// Get tool by slug
	tool, err := s.repo.GetToolBySlug(slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrToolNotFound
		}
		return nil, err
	}

	// Check if user already reviewed this tool
	hasReviewed, err := s.repo.HasUserReviewed(tool.ID, userID)
	if err != nil {
		return nil, err
	}
	if hasReviewed {
		return nil, ErrAlreadyReviewed
	}

	// Create review
	review := &domain.Review{
		ToolID:           tool.ID,
		UserID:           userID,
		RatingOverall:    input.RatingOverall,
		RatingEaseOfUse:  input.RatingEaseOfUse,
		RatingValue:      input.RatingValue,
		RatingAccuracy:   input.RatingAccuracy,
		RatingSpeed:      input.RatingSpeed,
		RatingSupport:    input.RatingSupport,
		Pros:             input.Pros,
		Cons:             input.Cons,
		PrimaryUseCase:   input.PrimaryUseCase,
		ReviewerRole:     input.ReviewerRole,
		CompanySize:      input.CompanySize,
		UsageContext:     input.UsageContext,
		ModerationStatus: "approved", // Auto-approve for now (can be configurable)
	}

	if err := s.repo.CreateReview(review); err != nil {
		return nil, err
	}

	// Update tool rating aggregates - best effort, don't fail if this errors
	// In production, this should be logged for monitoring
	_ = s.repo.UpdateToolRatingAggregates(tool.ID)

	// Use current time for response since GORM sets CreatedAt after Create
	createdAt := time.Now().UTC()
	if !review.CreatedAt.IsZero() {
		createdAt = review.CreatedAt
	}

	return &ReviewResponse{
		ID:              review.ID,
		RatingOverall:   review.RatingOverall,
		RatingEaseOfUse: review.RatingEaseOfUse,
		RatingValue:     review.RatingValue,
		RatingAccuracy:  review.RatingAccuracy,
		RatingSpeed:     review.RatingSpeed,
		RatingSupport:   review.RatingSupport,
		Pros:            review.Pros,
		Cons:            review.Cons,
		PrimaryUseCase:  review.PrimaryUseCase,
		ReviewerRole:    review.ReviewerRole,
		CompanySize:     review.CompanySize,
		UsageContext:    review.UsageContext,
		HelpfulCount:    review.HelpfulCount,
		CreatedAt:       createdAt.Format("2006-01-02T15:04:05Z"),
		User: UserBrief{
			ID: userID,
		},
	}, nil
}

// validateInput validates the review input
func (s *service) validateInput(input CreateReviewInput) error {
	// Rating overall is required
	if input.RatingOverall == 0 {
		return ErrRatingRequired
	}

	// Validate rating range
	if input.RatingOverall < 1 || input.RatingOverall > 5 {
		return ErrInvalidRating
	}

	// Validate optional ratings
	if input.RatingEaseOfUse != nil && (*input.RatingEaseOfUse < 1 || *input.RatingEaseOfUse > 5) {
		return ErrInvalidRating
	}
	if input.RatingValue != nil && (*input.RatingValue < 1 || *input.RatingValue > 5) {
		return ErrInvalidRating
	}
	if input.RatingAccuracy != nil && (*input.RatingAccuracy < 1 || *input.RatingAccuracy > 5) {
		return ErrInvalidRating
	}
	if input.RatingSpeed != nil && (*input.RatingSpeed < 1 || *input.RatingSpeed > 5) {
		return ErrInvalidRating
	}
	if input.RatingSupport != nil && (*input.RatingSupport < 1 || *input.RatingSupport > 5) {
		return ErrInvalidRating
	}

	// Pros required
	if input.Pros == "" {
		return ErrProsRequired
	}
	if utf8.RuneCountInString(input.Pros) > 500 {
		return ErrProsTooLong
	}

	// Cons required
	if input.Cons == "" {
		return ErrConsRequired
	}
	if utf8.RuneCountInString(input.Cons) > 500 {
		return ErrConsTooLong
	}

	return nil
}

// validatePagination ensures page and pageSize have valid values
func (s *service) validatePagination(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

// toReviewResponse converts a domain review to API response
func (s *service) toReviewResponse(r domain.Review) ReviewResponse {
	resp := ReviewResponse{
		ID:              r.ID,
		RatingOverall:   r.RatingOverall,
		RatingEaseOfUse: r.RatingEaseOfUse,
		RatingValue:     r.RatingValue,
		RatingAccuracy:  r.RatingAccuracy,
		RatingSpeed:     r.RatingSpeed,
		RatingSupport:   r.RatingSupport,
		Pros:            r.Pros,
		Cons:            r.Cons,
		PrimaryUseCase:  r.PrimaryUseCase,
		ReviewerRole:    r.ReviewerRole,
		CompanySize:     r.CompanySize,
		UsageContext:    r.UsageContext,
		HelpfulCount:    r.HelpfulCount,
		CreatedAt:       r.CreatedAt.Format("2006-01-02T15:04:05Z"),
		User: UserBrief{
			ID: r.UserID,
		},
	}

	// Include user display name if loaded
	if r.User.DisplayName != "" {
		resp.User.DisplayName = r.User.DisplayName
	}

	return resp
}
