### Story 6.3: User Profile & Activity

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
