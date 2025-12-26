package categories_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/your-org/ai-tools-atlas-backend/internal/categories"
	"github.com/your-org/ai-tools-atlas-backend/internal/domain"
	"gorm.io/gorm"
)

// MockRepository is a mock implementation of categories.Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) ListCategories() ([]domain.Category, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Category), args.Error(1)
}

func (m *MockRepository) GetCategoryBySlug(slug string) (*domain.Category, error) {
	args := m.Called(slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Category), args.Error(1)
}

func (m *MockRepository) ListToolsByCategory(categoryID uint, page, pageSize int) ([]domain.Tool, int64, error) {
	args := m.Called(categoryID, page, pageSize)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]domain.Tool), args.Get(1).(int64), args.Error(2)
}

func (m *MockRepository) GetCategoryByID(id uint) (*domain.Category, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Category), args.Error(1)
}

func (m *MockRepository) Create(category *domain.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockRepository) Update(category *domain.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepository) SlugExists(slug string, excludeID uint) (bool, error) {
	args := m.Called(slug, excludeID)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) GetToolCount(categoryID uint) (int64, error) {
	args := m.Called(categoryID)
	return args.Get(0).(int64), args.Error(1)
}

func TestServiceListCategories(t *testing.T) {
	t.Run("returns categories from repository", func(t *testing.T) {
		mockRepo := new(MockRepository)
		expectedCategories := []domain.Category{
			{ID: 1, Slug: "ai-writing", Name: "AI Writing"},
		}
		mockRepo.On("ListCategories").Return(expectedCategories, nil)

		service := categories.NewService(mockRepo)
		result, err := service.ListCategories()

		assert.NoError(t, err)
		assert.Equal(t, expectedCategories, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestServiceGetCategoryBySlug(t *testing.T) {
	t.Run("returns category when found", func(t *testing.T) {
		mockRepo := new(MockRepository)
		expectedCategory := &domain.Category{ID: 1, Slug: "ai-writing", Name: "AI Writing"}
		mockRepo.On("GetCategoryBySlug", "ai-writing").Return(expectedCategory, nil)

		service := categories.NewService(mockRepo)
		result, err := service.GetCategoryBySlug("ai-writing")

		assert.NoError(t, err)
		assert.Equal(t, expectedCategory, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns ErrCategoryNotFound when not found", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockRepo.On("GetCategoryBySlug", "non-existent").Return(nil, gorm.ErrRecordNotFound)

		service := categories.NewService(mockRepo)
		result, err := service.GetCategoryBySlug("non-existent")

		assert.Nil(t, result)
		assert.ErrorIs(t, err, categories.ErrCategoryNotFound)
		mockRepo.AssertExpectations(t)
	})
}

func TestServiceListToolsByCategory(t *testing.T) {
	t.Run("returns tools for valid category", func(t *testing.T) {
		mockRepo := new(MockRepository)
		category := &domain.Category{ID: 1, Slug: "ai-writing", Name: "AI Writing"}
		expectedTools := []domain.Tool{{ID: 1, Slug: "chatgpt", Name: "ChatGPT"}}

		mockRepo.On("GetCategoryBySlug", "ai-writing").Return(category, nil)
		mockRepo.On("ListToolsByCategory", uint(1), 1, 20).Return(expectedTools, int64(1), nil)

		service := categories.NewService(mockRepo)
		result, total, err := service.ListToolsByCategory("ai-writing", 1, 20)

		assert.NoError(t, err)
		assert.Equal(t, expectedTools, result)
		assert.Equal(t, int64(1), total)
		mockRepo.AssertExpectations(t)
	})

	t.Run("applies default pagination for invalid values", func(t *testing.T) {
		mockRepo := new(MockRepository)
		category := &domain.Category{ID: 1, Slug: "ai-writing", Name: "AI Writing"}

		mockRepo.On("GetCategoryBySlug", "ai-writing").Return(category, nil)
		mockRepo.On("ListToolsByCategory", uint(1), 1, 20).Return([]domain.Tool{}, int64(0), nil)

		service := categories.NewService(mockRepo)
		_, _, err := service.ListToolsByCategory("ai-writing", 0, 0)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("limits page size to 100", func(t *testing.T) {
		mockRepo := new(MockRepository)
		category := &domain.Category{ID: 1, Slug: "ai-writing", Name: "AI Writing"}

		mockRepo.On("GetCategoryBySlug", "ai-writing").Return(category, nil)
		mockRepo.On("ListToolsByCategory", uint(1), 1, 20).Return([]domain.Tool{}, int64(0), nil)

		service := categories.NewService(mockRepo)
		_, _, err := service.ListToolsByCategory("ai-writing", 1, 200)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns error when category not found", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockRepo.On("GetCategoryBySlug", "non-existent").Return(nil, gorm.ErrRecordNotFound)

		service := categories.NewService(mockRepo)
		result, total, err := service.ListToolsByCategory("non-existent", 1, 20)

		assert.Nil(t, result)
		assert.Equal(t, int64(0), total)
		assert.ErrorIs(t, err, categories.ErrCategoryNotFound)
		mockRepo.AssertExpectations(t)
	})
}
