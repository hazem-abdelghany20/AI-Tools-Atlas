package tags_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/your-org/ai-tools-atlas-backend/internal/domain"
	"github.com/your-org/ai-tools-atlas-backend/internal/tags"
)

// MockService is a mock implementation of tags.Service
type MockService struct {
	mock.Mock
}

func (m *MockService) ListTagsWithCount() ([]tags.TagWithCount, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]tags.TagWithCount), args.Error(1)
}

func (m *MockService) GetTagByID(id uint) (*domain.Tag, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Tag), args.Error(1)
}

func (m *MockService) CreateTag(input tags.CreateTagInput) (*domain.Tag, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Tag), args.Error(1)
}

func (m *MockService) UpdateTag(id uint, input tags.UpdateTagInput) (*domain.Tag, error) {
	args := m.Called(id, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Tag), args.Error(1)
}

func (m *MockService) DeleteTag(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func setupTestRouter(service tags.Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	admin := r.Group("/api/v1/admin")
	handler := tags.NewHandler(service)
	handler.RegisterAdminRoutes(admin)
	return r
}

func TestAdminListTags(t *testing.T) {
	t.Run("returns tags successfully", func(t *testing.T) {
		mockService := new(MockService)
		expectedTags := []tags.TagWithCount{
			{Tag: domain.Tag{ID: 1, Slug: "ml", Name: "Machine Learning"}, ToolCount: 5},
			{Tag: domain.Tag{ID: 2, Slug: "nlp", Name: "NLP"}, ToolCount: 3},
		}
		mockService.On("ListTagsWithCount").Return(expectedTags, nil)

		router := setupTestRouter(mockService)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/admin/tags", nil)
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
}

func TestAdminGetTag(t *testing.T) {
	t.Run("returns tag when found", func(t *testing.T) {
		mockService := new(MockService)
		expectedTag := &domain.Tag{ID: 1, Slug: "ml", Name: "Machine Learning"}
		mockService.On("GetTagByID", uint(1)).Return(expectedTag, nil)

		router := setupTestRouter(mockService)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/admin/tags/1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("returns 404 when not found", func(t *testing.T) {
		mockService := new(MockService)
		mockService.On("GetTagByID", uint(999)).Return(nil, tags.ErrTagNotFound)

		router := setupTestRouter(mockService)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/admin/tags/999", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("returns 400 for invalid ID", func(t *testing.T) {
		mockService := new(MockService)

		router := setupTestRouter(mockService)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/admin/tags/invalid", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestAdminCreateTag(t *testing.T) {
	t.Run("creates tag successfully", func(t *testing.T) {
		mockService := new(MockService)
		input := tags.CreateTagInput{Slug: "new-tag", Name: "New Tag"}
		createdTag := &domain.Tag{ID: 1, Slug: "new-tag", Name: "New Tag"}
		mockService.On("CreateTag", input).Return(createdTag, nil)

		router := setupTestRouter(mockService)
		body, _ := json.Marshal(input)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/admin/tags", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("returns 422 when slug is required", func(t *testing.T) {
		mockService := new(MockService)
		input := tags.CreateTagInput{Slug: "", Name: "New Tag"}
		mockService.On("CreateTag", input).Return(nil, tags.ErrSlugRequired)

		router := setupTestRouter(mockService)
		body, _ := json.Marshal(input)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/admin/tags", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("returns 409 when slug exists", func(t *testing.T) {
		mockService := new(MockService)
		input := tags.CreateTagInput{Slug: "existing", Name: "Existing"}
		mockService.On("CreateTag", input).Return(nil, tags.ErrSlugExists)

		router := setupTestRouter(mockService)
		body, _ := json.Marshal(input)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/admin/tags", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestAdminUpdateTag(t *testing.T) {
	t.Run("updates tag successfully", func(t *testing.T) {
		mockService := new(MockService)
		newName := "Updated Name"
		input := tags.UpdateTagInput{Name: &newName}
		updatedTag := &domain.Tag{ID: 1, Slug: "ml", Name: "Updated Name"}
		mockService.On("UpdateTag", uint(1), input).Return(updatedTag, nil)

		router := setupTestRouter(mockService)
		body, _ := json.Marshal(input)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PATCH", "/api/v1/admin/tags/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("returns 404 when not found", func(t *testing.T) {
		mockService := new(MockService)
		newName := "Updated Name"
		input := tags.UpdateTagInput{Name: &newName}
		mockService.On("UpdateTag", uint(999), input).Return(nil, tags.ErrTagNotFound)

		router := setupTestRouter(mockService)
		body, _ := json.Marshal(input)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PATCH", "/api/v1/admin/tags/999", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestAdminDeleteTag(t *testing.T) {
	t.Run("deletes tag successfully", func(t *testing.T) {
		mockService := new(MockService)
		mockService.On("DeleteTag", uint(1)).Return(nil)

		router := setupTestRouter(mockService)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/admin/tags/1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("returns 404 when not found", func(t *testing.T) {
		mockService := new(MockService)
		mockService.On("DeleteTag", uint(999)).Return(tags.ErrTagNotFound)

		router := setupTestRouter(mockService)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/admin/tags/999", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}
