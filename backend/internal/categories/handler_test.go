package categories_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/your-org/ai-tools-atlas-backend/internal/categories"
	"github.com/your-org/ai-tools-atlas-backend/internal/domain"
)

// MockService is a mock implementation of categories.Service
type MockService struct {
	mock.Mock
}

func (m *MockService) ListCategories() ([]domain.Category, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Category), args.Error(1)
}

func (m *MockService) GetCategoryBySlug(slug string) (*domain.Category, error) {
	args := m.Called(slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Category), args.Error(1)
}

func (m *MockService) ListToolsByCategory(slug string, page, pageSize int) ([]domain.Tool, int64, error) {
	args := m.Called(slug, page, pageSize)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]domain.Tool), args.Get(1).(int64), args.Error(2)
}

func (m *MockService) ListCategoriesWithCount() ([]categories.CategoryWithCount, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]categories.CategoryWithCount), args.Error(1)
}

func (m *MockService) GetCategoryByID(id uint) (*domain.Category, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Category), args.Error(1)
}

func (m *MockService) CreateCategory(input categories.CreateCategoryInput) (*domain.Category, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Category), args.Error(1)
}

func (m *MockService) UpdateCategory(id uint, input categories.UpdateCategoryInput) (*domain.Category, error) {
	args := m.Called(id, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Category), args.Error(1)
}

func (m *MockService) DeleteCategory(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func setupTestRouter(service categories.Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	v1 := r.Group("/api/v1")
	handler := categories.NewHandler(service)
	handler.RegisterRoutes(v1)
	return r
}

func TestListCategories(t *testing.T) {
	t.Run("returns categories successfully", func(t *testing.T) {
		mockService := new(MockService)
		expectedCategories := []domain.Category{
			{ID: 1, Slug: "ai-writing", Name: "AI Writing", DisplayOrder: 1},
			{ID: 2, Slug: "ai-image", Name: "AI Image", DisplayOrder: 2},
		}
		mockService.On("ListCategories").Return(expectedCategories, nil)

		router := setupTestRouter(mockService)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/categories", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response["data"])

		data := response["data"].([]interface{})
		assert.Len(t, data, 2)

		mockService.AssertExpectations(t)
	})

	t.Run("returns error on service failure", func(t *testing.T) {
		mockService := new(MockService)
		mockService.On("ListCategories").Return(nil, assert.AnError)

		router := setupTestRouter(mockService)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/categories", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestListToolsByCategory(t *testing.T) {
	t.Run("returns tools for valid category", func(t *testing.T) {
		mockService := new(MockService)
		expectedTools := []domain.Tool{
			{ID: 1, Slug: "chatgpt", Name: "ChatGPT"},
			{ID: 2, Slug: "claude", Name: "Claude"},
		}
		mockService.On("ListToolsByCategory", "ai-writing", 1, 20).Return(expectedTools, int64(2), nil)

		router := setupTestRouter(mockService)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/categories/ai-writing/tools", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response["data"])
		assert.NotNil(t, response["meta"])

		meta := response["meta"].(map[string]interface{})
		assert.Equal(t, float64(1), meta["page"])
		assert.Equal(t, float64(20), meta["page_size"])
		assert.Equal(t, float64(2), meta["total"])

		mockService.AssertExpectations(t)
	})

	t.Run("returns 404 for non-existent category", func(t *testing.T) {
		mockService := new(MockService)
		mockService.On("ListToolsByCategory", "non-existent", 1, 20).Return(nil, int64(0), categories.ErrCategoryNotFound)

		router := setupTestRouter(mockService)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/categories/non-existent/tools", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("uses custom pagination parameters", func(t *testing.T) {
		mockService := new(MockService)
		mockService.On("ListToolsByCategory", "ai-writing", 2, 10).Return([]domain.Tool{}, int64(0), nil)

		router := setupTestRouter(mockService)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/categories/ai-writing/tools?page=2&page_size=10", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})
}
