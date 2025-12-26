package auth

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
)

// MockRepository is a mock implementation of Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockRepository) GetByID(id uint) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockRepository) GetByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockRepository) EmailExists(email string) (bool, error) {
	args := m.Called(email)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) Update(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

// MockBookmarkService is a mock implementation of bookmarks.Service
type MockBookmarkService struct {
	mock.Mock
}

func (m *MockBookmarkService) GetBookmarks(userID uint, sessionID string) ([]interface{}, error) {
	args := m.Called(userID, sessionID)
	return nil, args.Error(1)
}

func (m *MockBookmarkService) AddBookmark(userID uint, sessionID string, toolID uint) (interface{}, error) {
	args := m.Called(userID, sessionID, toolID)
	return nil, args.Error(1)
}

func (m *MockBookmarkService) RemoveBookmark(userID uint, sessionID string, toolID uint) error {
	args := m.Called(userID, sessionID, toolID)
	return args.Error(0)
}

func (m *MockBookmarkService) IsBookmarked(userID uint, sessionID string, toolID uint) (bool, error) {
	args := m.Called(userID, sessionID, toolID)
	return args.Bool(0), args.Error(1)
}

func (m *MockBookmarkService) MigrateSessionBookmarks(userID uint, sessionID string) error {
	args := m.Called(userID, sessionID)
	return args.Error(0)
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestHandler_Register_Success(t *testing.T) {
	mockRepo := new(MockRepository)

	// Setup expectations
	mockRepo.On("EmailExists", "test@example.com").Return(false, nil)
	mockRepo.On("Create", mock.AnythingOfType("*domain.User")).Run(func(args mock.Arguments) {
		user := args.Get(0).(*domain.User)
		user.ID = 1 // Simulate database setting ID
	}).Return(nil)

	service := NewServiceWithRepo(mockRepo)
	handler := NewHandler(service, nil)

	router := setupTestRouter()
	v1 := router.Group("/api/v1")
	handler.RegisterRoutes(v1, nil)

	body := RegisterInput{
		Email:       "test@example.com",
		Password:    "password123",
		DisplayName: "Test User",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response, "data")

	data := response["data"].(map[string]interface{})
	user := data["user"].(map[string]interface{})
	assert.Equal(t, "test@example.com", user["email"])
	assert.Equal(t, "Test User", user["display_name"])
	assert.Equal(t, "user", user["role"])

	mockRepo.AssertExpectations(t)
}

func TestHandler_Register_InvalidEmail(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewServiceWithRepo(mockRepo)
	handler := NewHandler(service, nil)

	router := setupTestRouter()
	v1 := router.Group("/api/v1")
	handler.RegisterRoutes(v1, nil)

	body := RegisterInput{
		Email:       "invalid-email",
		Password:    "password123",
		DisplayName: "Test User",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	err := response["error"].(map[string]interface{})
	assert.Equal(t, "INVALID_EMAIL", err["code"])
}

func TestHandler_Register_PasswordTooShort(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewServiceWithRepo(mockRepo)
	handler := NewHandler(service, nil)

	router := setupTestRouter()
	v1 := router.Group("/api/v1")
	handler.RegisterRoutes(v1, nil)

	body := RegisterInput{
		Email:       "test@example.com",
		Password:    "short",
		DisplayName: "Test User",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	err := response["error"].(map[string]interface{})
	assert.Equal(t, "PASSWORD_TOO_SHORT", err["code"])
}

func TestHandler_Register_EmailAlreadyExists(t *testing.T) {
	mockRepo := new(MockRepository)
	mockRepo.On("EmailExists", "existing@example.com").Return(true, nil)

	service := NewServiceWithRepo(mockRepo)
	handler := NewHandler(service, nil)

	router := setupTestRouter()
	v1 := router.Group("/api/v1")
	handler.RegisterRoutes(v1, nil)

	body := RegisterInput{
		Email:       "existing@example.com",
		Password:    "password123",
		DisplayName: "Test User",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	err := response["error"].(map[string]interface{})
	assert.Equal(t, "EMAIL_EXISTS", err["code"])

	mockRepo.AssertExpectations(t)
}

func TestHandler_Login_Success(t *testing.T) {
	mockRepo := new(MockRepository)

	// Create a user with a hashed password
	service := NewServiceWithRepo(mockRepo)
	hashedPassword, _ := service.HashPassword("password123")
	user := &domain.User{
		ID:           1,
		Email:        "test@example.com",
		PasswordHash: hashedPassword,
		DisplayName:  "Test User",
		Role:         "user",
	}

	mockRepo.On("GetByEmail", "test@example.com").Return(user, nil)

	handler := NewHandler(service, nil)

	router := setupTestRouter()
	v1 := router.Group("/api/v1")
	handler.RegisterRoutes(v1, nil)

	body := LoginInput{
		Email:    "test@example.com",
		Password: "password123",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response, "data")

	data := response["data"].(map[string]interface{})
	respUser := data["user"].(map[string]interface{})
	assert.Equal(t, "test@example.com", respUser["email"])

	// Check that auth cookie was set
	cookies := w.Result().Cookies()
	var authCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "auth_token" {
			authCookie = c
			break
		}
	}
	assert.NotNil(t, authCookie)
	assert.True(t, authCookie.HttpOnly)

	mockRepo.AssertExpectations(t)
}

func TestHandler_Login_InvalidCredentials(t *testing.T) {
	mockRepo := new(MockRepository)
	mockRepo.On("GetByEmail", "test@example.com").Return(nil, ErrUserNotFound)

	service := NewServiceWithRepo(mockRepo)
	handler := NewHandler(service, nil)

	router := setupTestRouter()
	v1 := router.Group("/api/v1")
	handler.RegisterRoutes(v1, nil)

	body := LoginInput{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	err := response["error"].(map[string]interface{})
	assert.Equal(t, "INVALID_CREDENTIALS", err["code"])

	mockRepo.AssertExpectations(t)
}

func TestHandler_Logout(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewServiceWithRepo(mockRepo)
	handler := NewHandler(service, nil)

	router := setupTestRouter()
	v1 := router.Group("/api/v1")
	handler.RegisterRoutes(v1, nil)

	req, _ := http.NewRequest("POST", "/api/v1/auth/logout", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	// Check that auth cookie was cleared
	cookies := w.Result().Cookies()
	var authCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "auth_token" {
			authCookie = c
			break
		}
	}
	assert.NotNil(t, authCookie)
	assert.Equal(t, "", authCookie.Value)
	assert.True(t, authCookie.MaxAge < 0)
}

func TestHandler_GetCurrentUser_Unauthorized(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewServiceWithRepo(mockRepo)
	handler := NewHandler(service, nil)

	router := setupTestRouter()
	v1 := router.Group("/api/v1")

	// Auth middleware that rejects
	authMiddleware := func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{"code": "UNAUTHORIZED"},
		})
	}

	handler.RegisterRoutes(v1, authMiddleware)

	req, _ := http.NewRequest("GET", "/api/v1/me", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestHandler_GetCurrentUser_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	user := &domain.User{
		ID:          1,
		Email:       "test@example.com",
		DisplayName: "Test User",
		Role:        "user",
	}
	mockRepo.On("GetByID", uint(1)).Return(user, nil)

	service := NewServiceWithRepo(mockRepo)
	handler := NewHandler(service, nil)

	router := setupTestRouter()
	v1 := router.Group("/api/v1")

	// Auth middleware that sets user ID
	authMiddleware := func(c *gin.Context) {
		c.Set("userID", uint(1))
		c.Next()
	}

	handler.RegisterRoutes(v1, authMiddleware)

	req, _ := http.NewRequest("GET", "/api/v1/me", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response, "data")

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "test@example.com", data["email"])
	assert.Equal(t, "Test User", data["display_name"])

	mockRepo.AssertExpectations(t)
}
