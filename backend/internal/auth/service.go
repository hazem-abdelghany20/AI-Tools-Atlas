package auth

import (
	"errors"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/your-org/ai-tools-atlas-backend/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrPasswordTooShort   = errors.New("password must be at least 8 characters")
	ErrDisplayNameRequired = errors.New("display name is required")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

// JWT Claims structure
type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// RegisterInput contains the fields for user registration
type RegisterInput struct {
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password" binding:"required"`
	DisplayName string `json:"display_name" binding:"required"`
}

// LoginInput contains the fields for user login
type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserResponse is the safe user data returned to clients
type UserResponse struct {
	ID          uint      `json:"id"`
	Email       string    `json:"email"`
	DisplayName string    `json:"display_name"`
	Role        string    `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
}

// Service handles authentication operations
type Service struct {
	jwtSecret     []byte
	tokenDuration time.Duration
	repo          Repository
}

// NewService creates a new auth service
func NewService() *Service {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET environment variable is required")
	}

	return &Service{
		jwtSecret:     []byte(secret),
		tokenDuration: 7 * 24 * time.Hour, // 7 days
	}
}

// NewServiceWithRepo creates a new auth service with a repository
func NewServiceWithRepo(repo Repository) *Service {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET environment variable is required")
	}

	return &Service{
		jwtSecret:     []byte(secret),
		tokenDuration: 7 * 24 * time.Hour, // 7 days
		repo:          repo,
	}
}

// Register creates a new user account
func (s *Service) Register(input RegisterInput) (*domain.User, string, error) {
	// Validate input
	if err := s.validateRegistration(input); err != nil {
		return nil, "", err
	}

	// Check if email already exists
	exists, err := s.repo.EmailExists(input.Email)
	if err != nil {
		return nil, "", err
	}
	if exists {
		return nil, "", ErrEmailAlreadyExists
	}

	// Hash password
	hashedPassword, err := s.HashPassword(input.Password)
	if err != nil {
		return nil, "", err
	}

	// Create user
	user := &domain.User{
		Email:        strings.ToLower(strings.TrimSpace(input.Email)),
		PasswordHash: hashedPassword,
		DisplayName:  strings.TrimSpace(input.DisplayName),
		Role:         "user",
	}

	if err := s.repo.Create(user); err != nil {
		return nil, "", err
	}

	// Generate token
	token, err := s.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// Login authenticates a user and returns a token
func (s *Service) Login(input LoginInput) (*domain.User, string, error) {
	// Get user by email
	user, err := s.repo.GetByEmail(strings.ToLower(strings.TrimSpace(input.Email)))
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, "", ErrInvalidCredentials
		}
		return nil, "", err
	}

	// Check password
	if err := s.CheckPassword(input.Password, user.PasswordHash); err != nil {
		return nil, "", ErrInvalidCredentials
	}

	// Generate token
	token, err := s.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// GetCurrentUser returns the user by ID
func (s *Service) GetCurrentUser(userID uint) (*domain.User, error) {
	return s.repo.GetByID(userID)
}

// ToUserResponse converts a User to a safe response
func ToUserResponse(user *domain.User) UserResponse {
	return UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		Role:        user.Role,
		CreatedAt:   user.CreatedAt,
	}
}

// validateRegistration validates registration input
func (s *Service) validateRegistration(input RegisterInput) error {
	// Validate email format
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(input.Email) {
		return ErrInvalidEmail
	}

	// Validate password length
	if len(input.Password) < 8 {
		return ErrPasswordTooShort
	}

	// Validate display name
	if strings.TrimSpace(input.DisplayName) == "" {
		return ErrDisplayNameRequired
	}

	return nil
}

// GenerateToken creates a new JWT token for a user
func (s *Service) GenerateToken(userID uint, email string, role string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// ValidateToken validates a JWT token and returns the claims
func (s *Service) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// HashPassword creates a bcrypt hash of the password
func (s *Service) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPassword compares a password with a hash
func (s *Service) CheckPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
