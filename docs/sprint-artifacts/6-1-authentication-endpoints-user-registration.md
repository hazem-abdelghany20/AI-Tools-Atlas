### Story 6.1: Authentication Endpoints & User Registration

**As a** developer,
**I want** authentication endpoints for registration, login, and logout,
**So that** users can create accounts and authenticate.

**Acceptance Criteria:**

**Endpoint: `POST /api/v1/auth/register`**
- Request: `{ email, password, display_name }`
- Validation: email format, password min 8 chars, display_name required
- Hashes password with bcrypt
- Creates user with role: "user"
- Generates JWT token, sets HTTP-only cookie
- Migrates anonymous session bookmarks to new user
- Returns: `{ data: { user: { id, email, display_name, role } } }`
- Errors: 422 if validation fails, 409 if email already exists

**Endpoint: `POST /api/v1/auth/login`**
- Request: `{ email, password }`
- Validates credentials, checks password hash
- Generates JWT token, sets HTTP-only cookie
- Migrates anonymous session bookmarks to user
- Returns: `{ data: { user: { id, email, display_name, role } } }`
- Errors: 401 if credentials invalid

**Endpoint: `POST /api/v1/auth/logout`**
- Clears auth_token cookie
- Returns: 204 No Content

**Endpoint: `GET /api/v1/me`** (Authenticated)
- Returns current user profile
- Response: `{ data: { id, email, display_name, role, created_at } }`

**And:**
- Auth service: register, login, logout, getCurrentUser
- JWT claims: user_id, email, role, exp (7 days)
- HTTP-only cookie: auth_token, secure in production, SameSite: Lax

**Prerequisites:** Epic 1 (Auth setup, User model)

**Files Modified:**
- `backend/internal/auth/service.go` (add register, login methods)
- `backend/internal/auth/handler.go` (add registration/login handlers)

---
