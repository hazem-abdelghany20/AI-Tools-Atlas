### Story 1.4: JWT Authentication Setup

**As a** developer,
**I want** JWT-based authentication configured with HTTP-only cookies,
**So that** protected endpoints can verify user identity.

**Acceptance Criteria:**

**Given** I need auth for protected endpoints
**When** I implement JWT auth service and middleware
**Then** the following is complete:

**Auth Service (`internal/auth/service.go`):**
- `GenerateToken(userID uint, email string, role string) (string, error)` - creates JWT with claims
- `ValidateToken(tokenString string) (*Claims, error)` - verifies and parses JWT
- `HashPassword(password string) (string, error)` - bcrypt hashing
- `CheckPassword(password, hash string) error` - bcrypt comparison

**JWT Middleware (`internal/platform/http/middleware.go`):**
- `AuthRequired()` middleware reads JWT from cookie named `auth_token`
- Validates token using auth service
- Sets user context: `c.Set("user_id", claims.UserID)`, `c.Set("user_role", claims.Role)`
- Returns 401 with error envelope if token missing/invalid
- `AdminRequired()` middleware checks `user_role == "admin"`

**Cookie Configuration:**
- HTTP-only: true (prevent JavaScript access)
- Secure: true in production (HTTPS only)
- SameSite: Lax or Strict
- Max age: 7 days (configurable)

**And:**
- JWT secret loaded from env var `JWT_SECRET`
- Token expiration set to 7 days (configurable)
- Claims struct includes: `UserID`, `Email`, `Role`, `StandardClaims` (exp, iat)

**Technical Implementation:**

- Use `github.com/golang-jwt/jwt/v5` for JWT handling
- Use `golang.org/x/crypto/bcrypt` for password hashing
- Middleware example:
  ```go
  func AuthRequired() gin.HandlerFunc {
      return func(c *gin.Context) {
          cookie, err := c.Cookie("auth_token")
          if err != nil {
              responses.ErrorResponse(c, "unauthorized", "Authentication required", nil)
              c.Abort()
              return
          }
          claims, err := authService.ValidateToken(cookie)
          if err != nil {
              responses.ErrorResponse(c, "unauthorized", "Invalid token", nil)
              c.Abort()
              return
          }
          c.Set("user_id", claims.UserID)
          c.Set("user_role", claims.Role)
          c.Next()
      }
  }
  ```

**Prerequisites:** Story 1.1 (Backend), Story 1.3 (User model)

**Files Created:**
- `backend/internal/auth/service.go`
- `backend/internal/platform/http/middleware.go` (auth functions added)

---
