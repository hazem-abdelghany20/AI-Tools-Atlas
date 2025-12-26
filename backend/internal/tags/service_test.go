package tags_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/your-org/ai-tools-atlas-backend/internal/domain"
	"github.com/your-org/ai-tools-atlas-backend/internal/tags"
	"gorm.io/gorm"
)

// MockRepository is a mock implementation of tags.Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) ListTags() ([]domain.Tag, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Tag), args.Error(1)
}

func (m *MockRepository) GetTagByID(id uint) (*domain.Tag, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Tag), args.Error(1)
}

func (m *MockRepository) Create(tag *domain.Tag) error {
	args := m.Called(tag)
	return args.Error(0)
}

func (m *MockRepository) Update(tag *domain.Tag) error {
	args := m.Called(tag)
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

func (m *MockRepository) GetToolCount(tagID uint) (int64, error) {
	args := m.Called(tagID)
	return args.Get(0).(int64), args.Error(1)
}

func TestServiceListTagsWithCount(t *testing.T) {
	t.Run("returns tags with counts", func(t *testing.T) {
		mockRepo := new(MockRepository)
		expectedTags := []domain.Tag{
			{ID: 1, Slug: "machine-learning", Name: "Machine Learning"},
			{ID: 2, Slug: "nlp", Name: "NLP"},
		}
		mockRepo.On("ListTags").Return(expectedTags, nil)
		mockRepo.On("GetToolCount", uint(1)).Return(int64(5), nil)
		mockRepo.On("GetToolCount", uint(2)).Return(int64(3), nil)

		service := tags.NewService(mockRepo)
		result, err := service.ListTagsWithCount()

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, int64(5), result[0].ToolCount)
		assert.Equal(t, int64(3), result[1].ToolCount)
		mockRepo.AssertExpectations(t)
	})
}

func TestServiceGetTagByID(t *testing.T) {
	t.Run("returns tag when found", func(t *testing.T) {
		mockRepo := new(MockRepository)
		expectedTag := &domain.Tag{ID: 1, Slug: "machine-learning", Name: "Machine Learning"}
		mockRepo.On("GetTagByID", uint(1)).Return(expectedTag, nil)

		service := tags.NewService(mockRepo)
		result, err := service.GetTagByID(1)

		assert.NoError(t, err)
		assert.Equal(t, expectedTag, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns ErrTagNotFound when not found", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockRepo.On("GetTagByID", uint(999)).Return(nil, gorm.ErrRecordNotFound)

		service := tags.NewService(mockRepo)
		result, err := service.GetTagByID(999)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, tags.ErrTagNotFound)
		mockRepo.AssertExpectations(t)
	})
}

func TestServiceCreateTag(t *testing.T) {
	t.Run("creates tag successfully", func(t *testing.T) {
		mockRepo := new(MockRepository)
		input := tags.CreateTagInput{Slug: "new-tag", Name: "New Tag"}
		mockRepo.On("SlugExists", "new-tag", uint(0)).Return(false, nil)
		mockRepo.On("Create", mock.AnythingOfType("*domain.Tag")).Return(nil)

		service := tags.NewService(mockRepo)
		result, err := service.CreateTag(input)

		assert.NoError(t, err)
		assert.Equal(t, "new-tag", result.Slug)
		assert.Equal(t, "New Tag", result.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns error when slug is empty", func(t *testing.T) {
		mockRepo := new(MockRepository)
		input := tags.CreateTagInput{Slug: "", Name: "New Tag"}

		service := tags.NewService(mockRepo)
		result, err := service.CreateTag(input)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, tags.ErrSlugRequired)
	})

	t.Run("returns error when name is empty", func(t *testing.T) {
		mockRepo := new(MockRepository)
		input := tags.CreateTagInput{Slug: "new-tag", Name: ""}

		service := tags.NewService(mockRepo)
		result, err := service.CreateTag(input)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, tags.ErrNameRequired)
	})

	t.Run("returns error when slug exists", func(t *testing.T) {
		mockRepo := new(MockRepository)
		input := tags.CreateTagInput{Slug: "existing", Name: "Existing Tag"}
		mockRepo.On("SlugExists", "existing", uint(0)).Return(true, nil)

		service := tags.NewService(mockRepo)
		result, err := service.CreateTag(input)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, tags.ErrSlugExists)
		mockRepo.AssertExpectations(t)
	})
}

func TestServiceUpdateTag(t *testing.T) {
	t.Run("updates tag successfully", func(t *testing.T) {
		mockRepo := new(MockRepository)
		existingTag := &domain.Tag{ID: 1, Slug: "ml", Name: "Machine Learning"}
		newName := "ML Updated"
		input := tags.UpdateTagInput{Name: &newName}

		mockRepo.On("GetTagByID", uint(1)).Return(existingTag, nil)
		mockRepo.On("Update", existingTag).Return(nil)

		service := tags.NewService(mockRepo)
		result, err := service.UpdateTag(1, input)

		assert.NoError(t, err)
		assert.Equal(t, "ML Updated", result.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns error when tag not found", func(t *testing.T) {
		mockRepo := new(MockRepository)
		newName := "New Name"
		input := tags.UpdateTagInput{Name: &newName}
		mockRepo.On("GetTagByID", uint(999)).Return(nil, gorm.ErrRecordNotFound)

		service := tags.NewService(mockRepo)
		result, err := service.UpdateTag(999, input)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, tags.ErrTagNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns error when name is empty", func(t *testing.T) {
		mockRepo := new(MockRepository)
		existingTag := &domain.Tag{ID: 1, Slug: "ml", Name: "Machine Learning"}
		emptyName := ""
		input := tags.UpdateTagInput{Name: &emptyName}

		mockRepo.On("GetTagByID", uint(1)).Return(existingTag, nil)

		service := tags.NewService(mockRepo)
		result, err := service.UpdateTag(1, input)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, tags.ErrNameRequired)
		mockRepo.AssertExpectations(t)
	})
}

func TestServiceDeleteTag(t *testing.T) {
	t.Run("deletes tag successfully", func(t *testing.T) {
		mockRepo := new(MockRepository)
		existingTag := &domain.Tag{ID: 1, Slug: "ml", Name: "Machine Learning"}

		mockRepo.On("GetTagByID", uint(1)).Return(existingTag, nil)
		mockRepo.On("Delete", uint(1)).Return(nil)

		service := tags.NewService(mockRepo)
		err := service.DeleteTag(1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns error when tag not found", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockRepo.On("GetTagByID", uint(999)).Return(nil, gorm.ErrRecordNotFound)

		service := tags.NewService(mockRepo)
		err := service.DeleteTag(999)

		assert.ErrorIs(t, err, tags.ErrTagNotFound)
		mockRepo.AssertExpectations(t)
	})
}
