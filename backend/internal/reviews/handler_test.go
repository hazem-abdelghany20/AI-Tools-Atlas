package reviews_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/your-org/ai-tools-atlas-backend/internal/reviews"
)

// MockService is a mock implementation of reviews.Service
type MockService struct {
	mock.Mock
}

func (m *MockService) ListReviews(slug string, sort string, page, pageSize int) ([]reviews.ReviewResponse, int64, error) {
	args := m.Called(slug, sort, page, pageSize)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]reviews.ReviewResponse), args.Get(1).(int64), args.Error(2)
}

func (m *MockService) ListUserReviews(userID uint, page, pageSize int) ([]reviews.UserReviewResponse, int64, error) {
	args := m.Called(userID, page, pageSize)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]reviews.UserReviewResponse), args.Get(1).(int64), args.Error(2)
}

func (m *MockService) CreateReview(slug string, userID uint, input reviews.CreateReviewInput) (*reviews.ReviewResponse, error) {
	args := m.Called(slug, userID, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reviews.ReviewResponse), args.Error(1)
}

func setupTestRouter(handler *reviews.Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	v1 := r.Group("/api/v1")

	authMiddleware := func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	}

	handler.RegisterRoutes(v1, authMiddleware)
	return r
}

func TestListReviews(t *testing.T) {
	t.Run("returns reviews for valid tool slug", func(t *testing.T) {
		mockService := new(MockService)
		expectedReviews := []reviews.ReviewResponse{
			{
				ID:            1,
				RatingOverall: 5,
				Pros:          "Great AI",
				Cons:          "Expensive",
				User:          reviews.UserBrief{ID: 1, DisplayName: "John"},
			},
		}

		mockService.On("ListReviews", "chatgpt", "newest", 1, 10).Return(expectedReviews, int64(1), nil)

		handler := reviews.NewHandler(mockService)
		router := setupTestRouter(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/tools/chatgpt/reviews", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response["data"])
		mockService.AssertExpectations(t)
	})

	t.Run("returns 404 for non-existent tool", func(t *testing.T) {
		mockService := new(MockService)
		mockService.On("ListReviews", "non-existent", "newest", 1, 10).Return(nil, int64(0), reviews.ErrToolNotFound)

		handler := reviews.NewHandler(mockService)
		router := setupTestRouter(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/tools/non-existent/reviews", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("handles pagination parameters", func(t *testing.T) {
		mockService := new(MockService)
		mockService.On("ListReviews", "chatgpt", "highest", 2, 20).Return([]reviews.ReviewResponse{}, int64(0), nil)

		handler := reviews.NewHandler(mockService)
		router := setupTestRouter(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/tools/chatgpt/reviews?page=2&page_size=20&sort=highest", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestCreateReview(t *testing.T) {
	t.Run("creates review successfully", func(t *testing.T) {
		mockService := new(MockService)
		expectedReview := &reviews.ReviewResponse{
			ID:            1,
			RatingOverall: 5,
			Pros:          "Great AI",
			Cons:          "Expensive",
		}

		input := reviews.CreateReviewInput{
			RatingOverall: 5,
			Pros:          "Great AI",
			Cons:          "Expensive",
		}

		mockService.On("CreateReview", "chatgpt", uint(1), input).Return(expectedReview, nil)

		handler := reviews.NewHandler(mockService)
		router := setupTestRouter(handler)

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/tools/chatgpt/reviews", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("returns 404 for non-existent tool", func(t *testing.T) {
		mockService := new(MockService)
		input := reviews.CreateReviewInput{
			RatingOverall: 5,
			Pros:          "Great",
			Cons:          "None",
		}

		mockService.On("CreateReview", "non-existent", uint(1), input).Return(nil, reviews.ErrToolNotFound)

		handler := reviews.NewHandler(mockService)
		router := setupTestRouter(handler)

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/tools/non-existent/reviews", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("returns 409 for duplicate review", func(t *testing.T) {
		mockService := new(MockService)
		input := reviews.CreateReviewInput{
			RatingOverall: 5,
			Pros:          "Great",
			Cons:          "None",
		}

		mockService.On("CreateReview", "chatgpt", uint(1), input).Return(nil, reviews.ErrAlreadyReviewed)

		handler := reviews.NewHandler(mockService)
		router := setupTestRouter(handler)

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/tools/chatgpt/reviews", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("returns 422 for validation errors", func(t *testing.T) {
		mockService := new(MockService)
		input := reviews.CreateReviewInput{
			RatingOverall: 0, // Invalid
			Pros:          "Great",
			Cons:          "None",
		}

		mockService.On("CreateReview", "chatgpt", uint(1), input).Return(nil, reviews.ErrRatingRequired)

		handler := reviews.NewHandler(mockService)
		router := setupTestRouter(handler)

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/tools/chatgpt/reviews", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("returns 400 for invalid JSON", func(t *testing.T) {
		mockService := new(MockService)

		handler := reviews.NewHandler(mockService)
		router := setupTestRouter(handler)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/tools/chatgpt/reviews", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetUserReviews(t *testing.T) {
	t.Run("returns user reviews successfully", func(t *testing.T) {
		mockService := new(MockService)
		expectedReviews := []reviews.UserReviewResponse{
			{
				ID:            1,
				RatingOverall: 5,
				Pros:          "Great",
				Cons:          "None",
			},
		}

		mockService.On("ListUserReviews", uint(1), 1, 10).Return(expectedReviews, int64(1), nil)

		handler := reviews.NewHandler(mockService)
		router := setupTestRouter(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/me/reviews", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})
}
