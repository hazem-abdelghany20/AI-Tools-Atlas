package analytics

// TopToolsResponse contains different top tool lists
type TopToolsResponse struct {
	ByBookmarks []TopTool `json:"by_bookmarks"`
	ByRating    []TopTool `json:"by_rating"`
	ByReviews   []TopTool `json:"by_reviews"`
}

// Service defines the interface for analytics business logic
type Service interface {
	GetOverviewStats() (*OverviewStats, error)
	GetTopTools(limit int) (*TopToolsResponse, error)
	GetTopCategories(limit int) ([]TopCategory, error)
}

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new analytics service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// GetOverviewStats returns aggregate statistics
func (s *service) GetOverviewStats() (*OverviewStats, error) {
	return s.repo.GetOverviewStats()
}

// GetTopTools returns top tools by different metrics
func (s *service) GetTopTools(limit int) (*TopToolsResponse, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	byBookmarks, err := s.repo.GetTopToolsByBookmarks(limit)
	if err != nil {
		return nil, err
	}

	byRating, err := s.repo.GetTopToolsByRating(limit)
	if err != nil {
		return nil, err
	}

	byReviews, err := s.repo.GetTopToolsByReviews(limit)
	if err != nil {
		return nil, err
	}

	return &TopToolsResponse{
		ByBookmarks: byBookmarks,
		ByRating:    byRating,
		ByReviews:   byReviews,
	}, nil
}

// GetTopCategories returns top categories by tool count
func (s *service) GetTopCategories(limit int) ([]TopCategory, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}
	return s.repo.GetTopCategories(limit)
}
