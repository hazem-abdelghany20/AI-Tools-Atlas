package reviews_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/your-org/ai-tools-atlas-backend/internal/domain"
	"github.com/your-org/ai-tools-atlas-backend/internal/reviews"
	"gorm.io/gorm"
)

// MockRepository is a mock implementation of reviews.Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) ListReviewsByTool(toolID uint, sort string, page, pageSize int) ([]domain.Review, int64, error) {
	args := m.Called(toolID, sort, page, pageSize)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]domain.Review), args.Get(1).(int64), args.Error(2)
}

func (m *MockRepository) CreateReview(review *domain.Review) error {
	args := m.Called(review)
	return args.Error(0)
}

func (m *MockRepository) HasUserReviewed(toolID, userID uint) (bool, error) {
	args := m.Called(toolID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) UpdateToolRatingAggregates(toolID uint) error {
	args := m.Called(toolID)
	return args.Error(0)
}

func (m *MockRepository) GetToolBySlug(slug string) (*domain.Tool, error) {
	args := m.Called(slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Tool), args.Error(1)
}

func (m *MockRepository) ListReviewsByUser(userID uint, page, pageSize int) ([]domain.Review, int64, error) {
	args := m.Called(userID, page, pageSize)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]domain.Review), args.Get(1).(int64), args.Error(2)
}

func TestServiceListReviews(t *testing.T) {
	t.Run("returns reviews for valid tool slug", func(t *testing.T) {
		mockRepo := new(MockRepository)
		tool := &domain.Tool{ID: 1, Slug: "chatgpt", Name: "ChatGPT"}
		expectedReviews := []domain.Review{
			{
				ID:            1,
				ToolID:        1,
				UserID:        1,
				RatingOverall: 5,
				Pros:          "Great AI",
				Cons:          "Expensive",
				CreatedAt:     time.Now(),
				User:          domain.User{ID: 1, DisplayName: "John Doe"},
			},
		}

		mockRepo.On("GetToolBySlug", "chatgpt").Return(tool, nil)
		mockRepo.On("ListReviewsByTool", uint(1), reviews.SortNewest, 1, 10).Return(expectedReviews, int64(1), nil)

		service := reviews.NewService(mockRepo)
		result, total, err := service.ListReviews("chatgpt", reviews.SortNewest, 1, 10)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, result, 1)
		assert.Equal(t, uint(1), result[0].ID)
		assert.Equal(t, 5, result[0].RatingOverall)
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns ErrToolNotFound when tool not found", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockRepo.On("GetToolBySlug", "non-existent").Return(nil, gorm.ErrRecordNotFound)

		service := reviews.NewService(mockRepo)
		result, total, err := service.ListReviews("non-existent", reviews.SortNewest, 1, 10)

		assert.Nil(t, result)
		assert.Equal(t, int64(0), total)
		assert.ErrorIs(t, err, reviews.ErrToolNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("validates sort option to default", func(t *testing.T) {
		mockRepo := new(MockRepository)
		tool := &domain.Tool{ID: 1, Slug: "chatgpt", Name: "ChatGPT"}

		mockRepo.On("GetToolBySlug", "chatgpt").Return(tool, nil)
		mockRepo.On("ListReviewsByTool", uint(1), reviews.SortNewest, 1, 10).Return([]domain.Review{}, int64(0), nil)

		service := reviews.NewService(mockRepo)
		_, _, err := service.ListReviews("chatgpt", "invalid", 1, 10)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("validates pagination defaults", func(t *testing.T) {
		mockRepo := new(MockRepository)
		tool := &domain.Tool{ID: 1, Slug: "chatgpt", Name: "ChatGPT"}

		mockRepo.On("GetToolBySlug", "chatgpt").Return(tool, nil)
		mockRepo.On("ListReviewsByTool", uint(1), reviews.SortNewest, 1, 10).Return([]domain.Review{}, int64(0), nil)

		service := reviews.NewService(mockRepo)
		_, _, err := service.ListReviews("chatgpt", "", 0, 0)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("limits page size to 100", func(t *testing.T) {
		mockRepo := new(MockRepository)
		tool := &domain.Tool{ID: 1, Slug: "chatgpt", Name: "ChatGPT"}

		mockRepo.On("GetToolBySlug", "chatgpt").Return(tool, nil)
		mockRepo.On("ListReviewsByTool", uint(1), reviews.SortNewest, 1, 100).Return([]domain.Review{}, int64(0), nil)

		service := reviews.NewService(mockRepo)
		_, _, err := service.ListReviews("chatgpt", "", 1, 200)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestServiceCreateReview(t *testing.T) {
	t.Run("creates review successfully", func(t *testing.T) {
		mockRepo := new(MockRepository)
		tool := &domain.Tool{ID: 1, Slug: "chatgpt", Name: "ChatGPT"}

		mockRepo.On("GetToolBySlug", "chatgpt").Return(tool, nil)
		mockRepo.On("HasUserReviewed", uint(1), uint(1)).Return(false, nil)
		mockRepo.On("CreateReview", mock.AnythingOfType("*domain.Review")).Return(nil)
		mockRepo.On("UpdateToolRatingAggregates", uint(1)).Return(nil)

		service := reviews.NewService(mockRepo)
		input := reviews.CreateReviewInput{
			RatingOverall: 5,
			Pros:          "Great AI",
			Cons:          "Expensive",
		}

		result, err := service.CreateReview("chatgpt", 1, input)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 5, result.RatingOverall)
		assert.Equal(t, "Great AI", result.Pros)
		assert.Equal(t, "Expensive", result.Cons)
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns ErrToolNotFound when tool not found", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockRepo.On("GetToolBySlug", "non-existent").Return(nil, gorm.ErrRecordNotFound)

		service := reviews.NewService(mockRepo)
		input := reviews.CreateReviewInput{
			RatingOverall: 5,
			Pros:          "Great",
			Cons:          "None",
		}

		result, err := service.CreateReview("non-existent", 1, input)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, reviews.ErrToolNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns ErrAlreadyReviewed when user already reviewed", func(t *testing.T) {
		mockRepo := new(MockRepository)
		tool := &domain.Tool{ID: 1, Slug: "chatgpt", Name: "ChatGPT"}

		mockRepo.On("GetToolBySlug", "chatgpt").Return(tool, nil)
		mockRepo.On("HasUserReviewed", uint(1), uint(1)).Return(true, nil)

		service := reviews.NewService(mockRepo)
		input := reviews.CreateReviewInput{
			RatingOverall: 5,
			Pros:          "Great",
			Cons:          "None",
		}

		result, err := service.CreateReview("chatgpt", 1, input)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, reviews.ErrAlreadyReviewed)
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns ErrRatingRequired when rating is 0", func(t *testing.T) {
		mockRepo := new(MockRepository)

		service := reviews.NewService(mockRepo)
		input := reviews.CreateReviewInput{
			RatingOverall: 0,
			Pros:          "Great",
			Cons:          "None",
		}

		result, err := service.CreateReview("chatgpt", 1, input)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, reviews.ErrRatingRequired)
	})

	t.Run("returns ErrInvalidRating when rating out of range", func(t *testing.T) {
		mockRepo := new(MockRepository)

		service := reviews.NewService(mockRepo)
		input := reviews.CreateReviewInput{
			RatingOverall: 6,
			Pros:          "Great",
			Cons:          "None",
		}

		result, err := service.CreateReview("chatgpt", 1, input)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, reviews.ErrInvalidRating)
	})

	t.Run("returns ErrProsRequired when pros is empty", func(t *testing.T) {
		mockRepo := new(MockRepository)

		service := reviews.NewService(mockRepo)
		input := reviews.CreateReviewInput{
			RatingOverall: 5,
			Pros:          "",
			Cons:          "None",
		}

		result, err := service.CreateReview("chatgpt", 1, input)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, reviews.ErrProsRequired)
	})

	t.Run("returns ErrConsRequired when cons is empty", func(t *testing.T) {
		mockRepo := new(MockRepository)

		service := reviews.NewService(mockRepo)
		input := reviews.CreateReviewInput{
			RatingOverall: 5,
			Pros:          "Great",
			Cons:          "",
		}

		result, err := service.CreateReview("chatgpt", 1, input)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, reviews.ErrConsRequired)
	})

	t.Run("returns ErrProsTooLong when pros exceeds 500 chars", func(t *testing.T) {
		mockRepo := new(MockRepository)

		service := reviews.NewService(mockRepo)
		longPros := make([]byte, 501)
		for i := range longPros {
			longPros[i] = 'a'
		}
		input := reviews.CreateReviewInput{
			RatingOverall: 5,
			Pros:          string(longPros),
			Cons:          "None",
		}

		result, err := service.CreateReview("chatgpt", 1, input)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, reviews.ErrProsTooLong)
	})

	t.Run("returns ErrConsTooLong when cons exceeds 500 chars", func(t *testing.T) {
		mockRepo := new(MockRepository)

		service := reviews.NewService(mockRepo)
		longCons := make([]byte, 501)
		for i := range longCons {
			longCons[i] = 'a'
		}
		input := reviews.CreateReviewInput{
			RatingOverall: 5,
			Pros:          "Great",
			Cons:          string(longCons),
		}

		result, err := service.CreateReview("chatgpt", 1, input)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, reviews.ErrConsTooLong)
	})

	t.Run("validates optional rating dimensions", func(t *testing.T) {
		mockRepo := new(MockRepository)

		service := reviews.NewService(mockRepo)
		invalidRating := 6
		input := reviews.CreateReviewInput{
			RatingOverall:   5,
			RatingEaseOfUse: &invalidRating,
			Pros:            "Great",
			Cons:            "None",
		}

		result, err := service.CreateReview("chatgpt", 1, input)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, reviews.ErrInvalidRating)
	})
}
