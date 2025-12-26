package tools_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/your-org/ai-tools-atlas-backend/internal/domain"
	"github.com/your-org/ai-tools-atlas-backend/internal/tools"
)

// MockService is a mock implementation of tools.Service
type MockService struct {
	mock.Mock
}

func (m *MockService) ListTools(filters tools.ToolFilters, page, pageSize int) ([]domain.Tool, int64, error) {
	args := m.Called(filters, page, pageSize)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]domain.Tool), args.Get(1).(int64), args.Error(2)
}

func (m *MockService) SearchTools(query string, filters tools.ToolFilters, page, pageSize int) ([]domain.Tool, int64, error) {
	args := m.Called(query, filters, page, pageSize)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]domain.Tool), args.Get(1).(int64), args.Error(2)
}

func (m *MockService) GetToolBySlug(slug string) (*domain.Tool, error) {
	args := m.Called(slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Tool), args.Error(1)
}

func (m *MockService) GetToolByID(id uint) (*domain.Tool, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Tool), args.Error(1)
}

func (m *MockService) GetToolAlternatives(slug string) (*tools.AlternativesResult, error) {
	args := m.Called(slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tools.AlternativesResult), args.Error(1)
}

// Admin methods
func (m *MockService) ListToolsAdmin(search string, includeArchived bool, page, pageSize int) ([]domain.Tool, int64, error) {
	args := m.Called(search, includeArchived, page, pageSize)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]domain.Tool), args.Get(1).(int64), args.Error(2)
}

func (m *MockService) GetToolByIDAdmin(id uint) (*domain.Tool, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Tool), args.Error(1)
}

func (m *MockService) CreateTool(input tools.CreateToolInput) (*domain.Tool, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Tool), args.Error(1)
}

func (m *MockService) UpdateTool(id uint, input tools.UpdateToolInput) (*domain.Tool, error) {
	args := m.Called(id, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Tool), args.Error(1)
}

func (m *MockService) ArchiveTool(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func setupTestRouter(service tools.Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	v1 := r.Group("/api/v1")
	handler := tools.NewHandler(service)
	handler.RegisterRoutes(v1)
	return r
}

func TestListTools(t *testing.T) {
	t.Run("returns tools successfully with default filters", func(t *testing.T) {
		mockService := new(MockService)
		expectedTools := []domain.Tool{
			{ID: 1, Slug: "chatgpt", Name: "ChatGPT"},
			{ID: 2, Slug: "claude", Name: "Claude"},
		}

		filters := tools.ToolFilters{Sort: tools.SortTopRated}
		mockService.On("ListTools", filters, 1, 20).Return(expectedTools, int64(2), nil)

		router := setupTestRouter(mockService)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/tools", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response["data"])
		assert.NotNil(t, response["meta"])

		mockService.AssertExpectations(t)
	})

	t.Run("applies filters from query params", func(t *testing.T) {
		mockService := new(MockService)

		filters := tools.ToolFilters{
			Category:  "ai-writing",
			Price:     "free",
			MinRating: 4.0,
			Platform:  "web",
			Sort:      tools.SortNewest,
		}
		mockService.On("ListTools", filters, 2, 10).Return([]domain.Tool{}, int64(0), nil)

		router := setupTestRouter(mockService)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/tools?category=ai-writing&price=free&min_rating=4&platform=web&sort=newest&page=2&page_size=10", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestSearchTools(t *testing.T) {
	t.Run("searches tools with query", func(t *testing.T) {
		mockService := new(MockService)
		expectedTools := []domain.Tool{
			{ID: 1, Slug: "chatgpt", Name: "ChatGPT"},
		}

		filters := tools.ToolFilters{Sort: tools.SortTopRated}
		mockService.On("SearchTools", "chat", filters, 1, 20).Return(expectedTools, int64(1), nil)

		router := setupTestRouter(mockService)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/search/tools?q=chat", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response["data"])

		meta := response["meta"].(map[string]interface{})
		assert.Equal(t, "chat", meta["query"])

		mockService.AssertExpectations(t)
	})

	t.Run("returns all tools with empty query", func(t *testing.T) {
		mockService := new(MockService)

		filters := tools.ToolFilters{Sort: tools.SortTopRated}
		mockService.On("SearchTools", "", filters, 1, 20).Return([]domain.Tool{}, int64(0), nil)

		router := setupTestRouter(mockService)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/search/tools", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetTool(t *testing.T) {
	t.Run("returns tool by slug", func(t *testing.T) {
		mockService := new(MockService)
		expectedTool := &domain.Tool{ID: 1, Slug: "chatgpt", Name: "ChatGPT"}
		mockService.On("GetToolBySlug", "chatgpt").Return(expectedTool, nil)

		router := setupTestRouter(mockService)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/tools/chatgpt", nil)
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
		mockService.On("GetToolBySlug", "non-existent").Return(nil, tools.ErrToolNotFound)

		router := setupTestRouter(mockService)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/tools/non-existent", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}
