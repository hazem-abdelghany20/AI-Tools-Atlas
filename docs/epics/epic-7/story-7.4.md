### Story 7.4: Admin Badges & Analytics

**As an** admin,
**I want** to assign badges and view basic analytics,
**So that** I can highlight top tools and track engagement.

**Acceptance Criteria:**

**Badges Management:**
- `POST /api/v1/admin/tools/:id/badges` - assign badge to tool
- `DELETE /api/v1/admin/tools/:id/badges/:badge_id` - remove badge
- UI: badge assignment interface on tool edit form or dedicated page
- Predefined badges: "Top in Category", "Rising Star", "Editor's Pick"

**Analytics Dashboard (`/admin/analytics`):**
- **Overview stats:**
  - Total tools, categories, reviews, bookmarks
  - New tools this week/month
  - Active users (if tracking)
- **Top Tools:**
  - By views (if tracking page views)
  - By bookmarks
  - By rating
- **Top Categories:**
  - By tool count
  - By activity
- **Charts:** Simple bar/line charts for trends over time

**Endpoints:**
- `GET /api/v1/admin/analytics/overview`
- `GET /api/v1/admin/analytics/top-tools`
- `GET /api/v1/admin/analytics/top-categories`

**Technical Implementation:**
- Badge assignment: many-to-many relationship management
- Analytics: aggregate queries on tools, reviews, bookmarks tables
- Simple charts using a Vue chart library (e.g., Chart.js, ApexCharts)

**Prerequisites:** Story 7.1, Story 7.2

**New Files:**
- `frontend/src/views/admin/AdminAnalyticsView.vue`
- `backend/internal/admin/analytics_handler.go` (or in tools package)

---
