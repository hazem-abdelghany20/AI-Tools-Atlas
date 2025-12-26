### Story 6.3: User Profile & Activity

**Status:** Done ✅

**As a** logged-in user,
**I want** to view my submitted reviews and bookmarks in one place,
**So that** I can manage my activity.

**Acceptance Criteria:**

**Endpoint: `GET /api/v1/me/reviews`** (Authenticated)
- Returns user's submitted reviews
- Response: `{ data: [...reviews] }` with tool info populated
- Each review includes: tool (slug, name, logo_url), rating, created_at

**User Profile View (`/profile` or `/me`):**
- Page title: "حسابي" (My Account)
- Tabs or sections:
  - **My Reviews:** List of submitted reviews with links to tool profiles
    - Edit/delete options (future feature)
  - **My Bookmarks:** Same as BookmarksView
- User info displayed: display name, email, member since

**Technical Implementation:**
- Fetch user reviews: `apiClient.get('/me/reviews')`
- Display using ReviewCard component with tool context
- Bookmarks: reuse BookmarksView component or fetch

**Prerequisites:** Story 6.1, Story 4.1 (Reviews), Story 5.1 (Bookmarks)

**New Files:**
- `frontend/src/views/ProfileView.vue`

---

## Tasks

### Task 1: Backend - User Reviews Endpoint ✅
- [x] Add `ListReviewsByUser` to reviews repository with Tool preload
- [x] Add `UserReviewResponse` type to service with tool info
- [x] Add `ListUserReviews` method to service interface
- [x] Add `GetUserReviews` handler at GET /me/reviews
- [x] Register route in handler's RegisterRoutes method

### Task 2: Frontend - Profile View ✅
- [x] Create `ProfileView.vue` with user header section
- [x] Implement tabs component (Reviews/Bookmarks)
- [x] Fetch and display user reviews with tool context
- [x] Display review moderation status badges
- [x] Integrate bookmarks store for bookmarks tab
- [x] Add empty states for both tabs
- [x] Add loading skeletons

### Task 3: Router Integration ✅
- [x] Add /profile route with requiresAuth meta
- [x] Redirect to home if not authenticated

---

## Change Log:
| Date | Change |
|------|--------|
| 2025-12-26 | Story implementation completed |
| 2025-12-26 | Code Review: Passed - only LOW issues |
