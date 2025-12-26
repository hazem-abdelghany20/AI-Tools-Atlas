### Story 6.1: Authentication Endpoints & User Registration

**As a** developer,
**I want** authentication endpoints for registration, login, and logout,
**So that** users can create accounts and authenticate.

**Status:** done

---

## Acceptance Criteria:

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

---

## Tasks/Subtasks:

- [x] **Task 1: Create User Repository**
  - [x] 1.1 Define Repository interface
  - [x] 1.2 Implement Create method
  - [x] 1.3 Implement GetByID method
  - [x] 1.4 Implement GetByEmail method
  - [x] 1.5 Implement EmailExists method
  - [x] 1.6 Implement Update method

- [x] **Task 2: Extend Auth Service**
  - [x] 2.1 Add RegisterInput and LoginInput types
  - [x] 2.2 Add UserResponse type
  - [x] 2.3 Implement Register with validation
  - [x] 2.4 Implement Login with credential check
  - [x] 2.5 Implement GetCurrentUser
  - [x] 2.6 Add NewServiceWithRepo constructor

- [x] **Task 3: Create Auth Handler**
  - [x] 3.1 Implement Register endpoint
  - [x] 3.2 Implement Login endpoint
  - [x] 3.3 Implement Logout endpoint
  - [x] 3.4 Implement GetCurrentUser endpoint
  - [x] 3.5 Handle HTTP-only cookie setting
  - [x] 3.6 Integrate session bookmark migration

- [x] **Task 4: Update Router**
  - [x] 4.1 Add auth repository initialization
  - [x] 4.2 Register auth routes
  - [x] 4.3 Verify middleware integration

- [x] **Task 5: Write Tests**
  - [x] 5.1 Test Register success
  - [x] 5.2 Test Register validation errors
  - [x] 5.3 Test Login success
  - [x] 5.4 Test Login invalid credentials
  - [x] 5.5 Test Logout
  - [x] 5.6 Test GetCurrentUser

---

## File List:
- `backend/internal/auth/repository.go` (created)
- `backend/internal/auth/service.go` (modified)
- `backend/internal/auth/handler.go` (created)
- `backend/internal/auth/handler_test.go` (created)
- `backend/internal/platform/http/router.go` (modified)

---

## Change Log:
| Date | Change |
|------|--------|
| 2025-12-26 | Story implementation completed |
| 2025-12-26 | Code Review: CRITICAL - Fixed context key "userID" to "user_id", added type assertion safety |
