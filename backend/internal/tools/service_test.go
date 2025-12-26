package tools_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/your-org/ai-tools-atlas-backend/internal/domain"
	"github.com/your-org/ai-tools-atlas-backend/internal/tools"
	"gorm.io/gorm"
)

// MockRepository is a mock implementation of tools.Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) ListTools(filters tools.ToolFilters, page, pageSize int) ([]domain.Tool, int64, error) {
	args := m.Called(filters, page, pageSize)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]domain.Tool), args.Get(1).(int64), args.Error(2)
}

func (m *MockRepository) SearchTools(query string, filters tools.ToolFilters, page, pageSize int) ([]domain.Tool, int64, error) {
	args := m.Called(query, filters, page, pageSize)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]domain.Tool), args.Get(1).(int64), args.Error(2)
}

func (m *MockRepository) GetToolBySlug(slug string) (*domain.Tool, error) {
	args := m.Called(slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Tool), args.Error(1)
}

func (m *MockRepository) GetToolByID(id uint) (*domain.Tool, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Tool), args.Error(1)
}

func (m *MockRepository) GetToolAlternatives(toolID uint, limit int) (*tools.AlternativesResult, error) {
	args := m.Called(toolID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tools.AlternativesResult), args.Error(1)
}

// Admin methods
func (m *MockRepository) ListToolsAdmin(search string, includeArchived bool, page, pageSize int) ([]domain.Tool, int64, error) {
	args := m.Called(search, includeArchived, page, pageSize)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]domain.Tool), args.Get(1).(int64), args.Error(2)
}

func (m *MockRepository) GetToolByIDAdmin(id uint) (*domain.Tool, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Tool), args.Error(1)
}

func (m *MockRepository) Create(tool *domain.Tool) error {
	args := m.Called(tool)
	return args.Error(0)
}

func (m *MockRepository) Update(tool *domain.Tool) error {
	args := m.Called(tool)
	return args.Error(0)
}

func (m *MockRepository) Archive(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepository) SlugExists(slug string, excludeID uint) (bool, error) {
	args := m.Called(slug, excludeID)
	return args.Bool(0), args.Error(1)
}

func TestServiceListTools(t *testing.T) {
	t.Run("returns tools from repository", func(t *testing.T) {
		mockRepo := new(MockRepository)
		expectedTools := []domain.Tool{{ID: 1, Slug: "chatgpt", Name: "ChatGPT"}}

		filters := tools.ToolFilters{Sort: tools.SortTopRated}
		mockRepo.On("ListTools", filters, 1, 20).Return(expectedTools, int64(1), nil)

		service := tools.NewService(mockRepo)
		result, total, err := service.ListTools(filters, 1, 20)

		assert.NoError(t, err)
		assert.Equal(t, expectedTools, result)
		assert.Equal(t, int64(1), total)
		mockRepo.AssertExpectations(t)
	})

	t.Run("validates pagination defaults", func(t *testing.T) {
		mockRepo := new(MockRepository)

		filters := tools.ToolFilters{Sort: tools.SortTopRated}
		mockRepo.On("ListTools", filters, 1, 20).Return([]domain.Tool{}, int64(0), nil)

		service := tools.NewService(mockRepo)
		_, _, err := service.ListTools(tools.ToolFilters{}, 0, 0)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("limits page size to 100", func(t *testing.T) {
		mockRepo := new(MockRepository)

		filters := tools.ToolFilters{Sort: tools.SortTopRated}
		mockRepo.On("ListTools", filters, 1, 100).Return([]domain.Tool{}, int64(0), nil)

		service := tools.NewService(mockRepo)
		_, _, err := service.ListTools(tools.ToolFilters{}, 1, 200)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("validates sort option to default", func(t *testing.T) {
		mockRepo := new(MockRepository)

		filters := tools.ToolFilters{Sort: tools.SortTopRated}
		mockRepo.On("ListTools", filters, 1, 20).Return([]domain.Tool{}, int64(0), nil)

		service := tools.NewService(mockRepo)
		_, _, err := service.ListTools(tools.ToolFilters{Sort: "invalid"}, 1, 20)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestServiceSearchTools(t *testing.T) {
	t.Run("searches tools via repository", func(t *testing.T) {
		mockRepo := new(MockRepository)
		expectedTools := []domain.Tool{{ID: 1, Slug: "chatgpt", Name: "ChatGPT"}}

		filters := tools.ToolFilters{Sort: tools.SortTopRated}
		mockRepo.On("SearchTools", "chat", filters, 1, 20).Return(expectedTools, int64(1), nil)

		service := tools.NewService(mockRepo)
		result, total, err := service.SearchTools("chat", filters, 1, 20)

		assert.NoError(t, err)
		assert.Equal(t, expectedTools, result)
		assert.Equal(t, int64(1), total)
		mockRepo.AssertExpectations(t)
	})
}

func TestServiceGetToolBySlug(t *testing.T) {
	t.Run("returns tool when found", func(t *testing.T) {
		mockRepo := new(MockRepository)
		expectedTool := &domain.Tool{ID: 1, Slug: "chatgpt", Name: "ChatGPT"}
		mockRepo.On("GetToolBySlug", "chatgpt").Return(expectedTool, nil)

		service := tools.NewService(mockRepo)
		result, err := service.GetToolBySlug("chatgpt")

		assert.NoError(t, err)
		assert.Equal(t, expectedTool, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns ErrToolNotFound when not found", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockRepo.On("GetToolBySlug", "non-existent").Return(nil, gorm.ErrRecordNotFound)

		service := tools.NewService(mockRepo)
		result, err := service.GetToolBySlug("non-existent")

		assert.Nil(t, result)
		assert.ErrorIs(t, err, tools.ErrToolNotFound)
		mockRepo.AssertExpectations(t)
	})
}

func TestServiceGetToolByID(t *testing.T) {
	t.Run("returns tool when found", func(t *testing.T) {
		mockRepo := new(MockRepository)
		expectedTool := &domain.Tool{ID: 1, Slug: "chatgpt", Name: "ChatGPT"}
		mockRepo.On("GetToolByID", uint(1)).Return(expectedTool, nil)

		service := tools.NewService(mockRepo)
		result, err := service.GetToolByID(1)

		assert.NoError(t, err)
		assert.Equal(t, expectedTool, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns ErrToolNotFound when not found", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockRepo.On("GetToolByID", uint(999)).Return(nil, gorm.ErrRecordNotFound)

		service := tools.NewService(mockRepo)
		result, err := service.GetToolByID(999)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, tools.ErrToolNotFound)
		mockRepo.AssertExpectations(t)
	})
}
