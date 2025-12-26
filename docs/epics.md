# AI Tools Atlas - Epic Breakdown

**Author:** Hazzouma
**Date:** 2025-12-26
**Project Level:** Low-Medium
**Target Scale:** MVP Web Application

---

## Overview

This document provides the complete epic and story breakdown for AI Tools Atlas, decomposing the requirements from the [PRD](./prd.md) into implementable stories.

**Living Document Notice:** This is the initial version created with full context from PRD, Architecture, and UX Design specifications.

## Context Validation

### ✅ Prerequisites Met

All required documents have been loaded and analyzed:

1. **PRD.md** - Product Requirements Document
   - 30 functional requirements (FR1-FR30) covering discovery, profiles, reviews, comparison, bookmarks, accounts, admin, moderation, and analytics
   - Complete product scope: MVP, Growth, and Vision features defined
   - User journeys for New-to-AI Explorer, Power User Consolidator, Updates Tracker, Admin, and Moderator
   - Success criteria and measurable outcomes established

2. **Architecture.md** - Technical Implementation Context
   - **Backend:** Go + Gin + GORM + PostgreSQL with complete data model
   - **Frontend:** Vue 3 + Vite + TypeScript + Pinia + Tailwind CSS
   - **API Design:** RESTful `/api/v1` endpoints with standardized responses
   - **Auth:** JWT-based authentication with HTTP-only cookies
   - **Data Model:** 11 core entities (Tool, Category, Tag, Review, Bookmark, Badge, ToolAlternative, Media, User, plus join tables)
   - Complete naming conventions, structure patterns, and consistency rules defined

3. **UX Design Specification.md** - User Experience Context
   - **Core Experience:** Hero search-driven discovery ("type your situation, see right tools")
   - **Visual Design:** Dark theme (#05060A-#0A0B10) with neon blue primary (#3B82F6-#2563EB)
   - **Typography:** Cairo font, Arabic-first, RTL-aware layouts
   - **Patterns:** G2-style directory adapted for Arabic users with structured cards, visible filters, consistent profiles
   - **Key Flows:** Discovery → Results → Profile → Compare → Bookmark with skeleton loaders and lazy loading

### 📊 Context Analysis Summary

**Technical Stack:**
- Backend: Go/Gin REST API serving JSON over `/api/v1` endpoints
- Database: PostgreSQL with normalized relational model
- Frontend: Vue 3 SPA consuming backend API, Pinia for state management
- Styling: Tailwind CSS with dark theme and RTL support

**Key Technical Decisions:**
- snake_case for database tables/columns and API fields
- PascalCase for Go types, camelCase for functions
- kebab-case for URLs and routes
- Standardized API response envelopes: `{ data, meta }` for lists, `{ data }` for single resources, `{ error }` for errors
- JWT tokens in HTTP-only cookies for auth
- Slugs for SEO-friendly URLs (`/tools/:slug`, `/categories/:slug`)

**User Experience Priorities:**
- Arabic-first, RTL-default interface with Cairo typography
- Hero search as primary entry point with free-text queries
- G2-inspired structured cards showing: name, "best for", rating, pricing badge, tags, bookmark/compare actions
- Side-by-side comparison for 2-4 tools
- Skeleton states and optimistic UI for perceived performance

## Epic Structure Overview

This breakdown organizes the 30 functional requirements into **8 epics** that deliver incremental user value:

1. **Epic 1: Foundation & Core Infrastructure** - Enable all subsequent development with database, API framework, and auth
2. **Epic 2: Tool Discovery & Browsing** - Users can find AI tools relevant to their needs via hero search, categories, and filters
3. **Epic 3: Rich Tool Profiles** - Users can deeply evaluate individual tools with metadata, media, and alternatives
4. **Epic 4: Reviews & Ratings** - Users can see peer feedback and contribute their own structured reviews
5. **Epic 5: Comparison & Shortlists** - Users can bookmark tools, manage shortlists, and compare options side-by-side
6. **Epic 6: User Accounts & Persistence** - Users can save their work across sessions and devices (lightweight auth)
7. **Epic 7: Admin & Catalog Management** - Admins can maintain a high-quality, accurate tool directory
8. **Epic 8: Moderation & Content Quality** - The platform maintains trusted, abuse-free user-generated content

**Note on Analytics (FR29-FR30):** Basic analytics capabilities will be integrated into Epic 7 (Admin) as operational dashboards rather than a separate epic.

---

## Functional Requirements Inventory

### Content Discovery & Search (FR1-FR5)
- **FR1:** Free-text search for AI tools by query (name, use case, role, industry)
- **FR2:** Browse AI tools by top-level category
- **FR3:** Refine results using filters (category, price, rating, platform)
- **FR4:** Sort results (top rated, most bookmarked, trending, newest)
- **FR5:** Open tool detail page from search/category results

### Tool Profiles & Reviews (FR6-FR11)
- **FR6:** View tool profile with core metadata (name, logo, tagline, description, "best for", use cases, pricing, target users, platforms)
- **FR7:** View rich media (screenshots, videos) on tool profile
- **FR8:** View structured reviews with ratings, pros, cons, use case, reviewer role/context
- **FR9:** Logged-in users submit structured reviews
- **FR10:** View aggregated rating information (overall rating, review count, dimension ratings)
- **FR11:** View social proof indicators (bookmark counts, badges)

### Comparison & Alternatives (FR12-FR14)
- **FR12:** Select multiple tools and add to comparison set
- **FR13:** View comparison page showing tools side-by-side
- **FR14:** View alternative/similar tools from tool profile

### Shortlists & Bookmarks (FR15-FR18)
- **FR15:** Bookmark a tool to personal shortlist
- **FR16:** View shortlist/bookmarks as dedicated list
- **FR17:** Remove tools from shortlist
- **FR18:** Send bookmarked tools to comparison view

### User Accounts & Sessions (FR19-FR20)
- **FR19:** Create account or sign in to persist reviews/bookmarks
- **FR20:** View own submitted reviews and bookmarked tools

### Admin & Catalog Management (FR21-FR24)
- **FR21:** Admin create/edit/archive tool records
- **FR22:** Admin manage categories and tags
- **FR23:** Admin manage curated lists and badges
- **FR24:** Admin identify and resolve catalog issues

### Moderation & Content Quality (FR25-FR28)
- **FR25:** Users report tools or reviews
- **FR26:** Moderators view queue of reported content
- **FR27:** Moderators approve/edit/hide/remove content
- **FR28:** Moderators view moderation action history

### Analytics & Observability (FR29-FR30)
- **FR29:** Admin view engagement signals
- **FR30:** Admin view top categories/tools by metrics

**Total: 30 Functional Requirements**

---

## Epic Planning with Full Context

### Epic 1: Foundation & Core Infrastructure

**User Value Statement:** Enable all subsequent user-facing features by establishing the technical foundation that allows developers to build a fast, secure, Arabic-first AI tools directory.

**PRD Coverage:** Indirectly supports all FRs by providing the infrastructure layer.

**Technical Context (from Architecture):**
- Initialize Go module with Gin framework for REST API serving
- Set up PostgreSQL database with complete schema: `tools`, `categories`, `tags`, `tool_tags`, `reviews`, `bookmarks`, `badges`, `tool_badges`, `tool_alternatives`, `media`, `users`
- Implement database migrations using golang-migrate
- Configure GORM models for all entities with proper foreign keys and indexes
- Set up `/api/v1` routing structure with standardized response helpers in `internal/platform/http`
- Implement JWT-based authentication with HTTP-only cookies
- Create Vue 3 + Vite + TypeScript frontend project with Pinia state management
- Configure Tailwind CSS with dark theme tokens and RTL support
- Set up Vue Router with key routes: `/`, `/tools/:slug`, `/categories/:slug`, `/compare`, `/bookmarks`
- Implement shared API client in `src/lib/apiClient.ts`

**UX Integration (from UX Design):**
- Configure Tailwind with dark background (#05060A-#0A0B10) and neon blue primary (#3B82F6-#2563EB)
- Set up Cairo font with proper Arabic typography scale
- Create base layout components: `AppShell.vue`, `HeaderNav.vue` with RTL awareness
- Implement skeleton loading components for performance patterns

**Dependencies:** None - this is the foundation epic.

---

### Epic 2: Tool Discovery & Browsing

**User Value Statement:** Users can quickly find AI tools relevant to their specific needs and use cases through intuitive search, browsing, and filtering - supporting the primary "New to AI in my field" user journey.

**PRD Coverage:** FR1, FR2, FR3, FR4, FR5 (Content Discovery & Search)

**Technical Context (from Architecture):**
- Implement `GET /api/v1/tools` with query parameters for filtering (category, price, min_rating, platform) and sorting (top_rated, most_bookmarked, trending, newest)
- Implement `GET /api/v1/search/tools?q=...` for free-text search
- Implement `GET /api/v1/categories` and `GET /api/v1/categories/:slug/tools`
- Backend services in `internal/tools` and `internal/categories` with repository layer using GORM
- Database indexes on `tools.primary_category_id`, `tools.avg_rating_overall`, `tools.bookmark_count`, `tools.trending_score`
- Response format: `{ data: [...], meta: { page, page_size, total } }`

**UX Integration (from UX Design):**
- **HomeView.vue:** Large hero search bar with Arabic prompt text "اكتب وضعك أو ما تريد القيام به…"
- Category grid/strip below hero with icons and labels
- **SearchResultsView.vue:** Results list with visible filter panel (RTL drawer on mobile)
- **ToolCard.vue component:** Display name, logo, "best for" line, rating + review count, pricing badge, tags, bookmark + compare buttons
- Filters panel showing: category, price (Free/Freemium/Paid), rating threshold, platform
- Sort dropdown: Top rated, Most bookmarked, Trending, Newest
- Skeleton states while loading results
- Empty state with clear messaging: "لم تجد ما تريد؟ استعرض حسب الفئة"

**Dependencies:** Epic 1 (Foundation)

---

### Epic 3: Rich Tool Profiles

**User Value Statement:** Users can deeply evaluate individual tools with comprehensive information including metadata, rich media, and alternatives - enabling confident tool selection decisions.

**PRD Coverage:** FR6, FR7, FR14 (Tool Profiles & Alternatives), FR11 (Social Proof)

**Technical Context (from Architecture):**
- Implement `GET /api/v1/tools/:slug` returning complete tool object with relationships (category, tags, media, badges, alternatives)
- Implement `GET /api/v1/tools/:slug/alternatives` for similar/alternative tools
- Join queries to load: tool, primary category, tags via `tool_tags`, media items, badges via `tool_badges`, alternative tools via `tool_alternatives`
- Response includes: name, slug, logo_url, tagline, description, best_for, primary_use_cases, pricing_summary, target_roles, platforms, has_free_tier, official_url
- Media model with type (screenshot/video), url, thumbnail_url, display_order

**UX Integration (from UX Design):**
- **ToolProfileView.vue:** Consistent, scannable profile layout
- **Hero section:** Name, logo, tagline, rating, review count, pricing badge, primary CTA "Visit Tool"
- **Overview section:** "Best for", description, primary use cases, target roles
- **Features section:** Key features list or grid
- **Pricing section:** Pricing summary (free tier, starting price, billing model)
- **ToolMediaGallery.vue:** Screenshots with lightbox, YouTube embeds lazy-loaded to avoid blocking page render
- **Alternatives section:** "Similar Tools" and "Alternatives to X" as horizontal scrolling cards
- **Social proof:** Badge pills ("Top in Category", "Rising Star"), bookmark count display
- RTL-aware layout with proper Arabic typography
- Skeleton loaders for profile sections

**Dependencies:** Epic 1 (Foundation), Epic 2 (Discovery - users navigate from search to profiles)

---

### Epic 4: Reviews & Ratings

**User Value Statement:** Users can see authentic peer feedback through structured reviews and contribute their own experiences - building trust and enabling informed decisions.

**PRD Coverage:** FR8, FR9, FR10 (Reviews & Ratings)

**Technical Context (from Architecture):**
- Implement `GET /api/v1/tools/:slug/reviews` with pagination and sorting (newest, most_helpful)
- Implement `POST /api/v1/tools/:slug/reviews` (authenticated) with validation
- Reviews table schema: tool_id, user_id, rating_overall, rating_ease_of_use, rating_value, rating_accuracy, rating_speed, rating_support, pros (text), cons (text), primary_use_case, reviewer_role, company_size, usage_context, helpful_count, created_at
- Aggregation logic for avg_rating_overall and dimension ratings, stored/cached on tools table
- Validation: required fields (rating_overall, pros, cons), length limits, rate limiting on submissions

**UX Integration (from UX Design):**
- **ReviewList.vue component:** Display on ToolProfileView below media/features
- Each review card shows: Overall rating (stars), reviewer role + company size, primary use case, pros/cons in structured format, usage context, submission date
- Sort controls: Newest, Most helpful, Highest rated, Lowest rated
- Empty state: "No reviews yet. Be the first!" with prompt to add review (if logged in)
- **ReviewForm.vue component:** Modal or inline form for authenticated users
- Form fields: Rating stars (overall + 5 dimensions), pros (textarea), cons (textarea), primary use case (dropdown), role (dropdown), company size (dropdown), usage context (checkboxes)
- Client-side validation with clear error messages in Arabic
- Success feedback: "شكرًا! تمت إضافة مراجعتك" with review appearing in list
- **Aggregated ratings display:** Overall rating with review count at top of profile, dimension ratings as small horizontal bars or progress indicators

**Dependencies:** Epic 1 (Foundation), Epic 3 (Tool Profiles - reviews are displayed on profile pages), Epic 6 (Accounts - for submitting reviews)

---

### Epic 5: Comparison & Shortlists

**User Value Statement:** Users can bookmark promising tools, manage their shortlist, and compare options side-by-side to make final consolidation and selection decisions.

**PRD Coverage:** FR12, FR13 (Comparison), FR15, FR16, FR17, FR18 (Bookmarks & Shortlists)

**Technical Context (from Architecture):**
- Implement `GET /api/v1/me/bookmarks` (authenticated or session-based) returning list of bookmarked tools
- Implement `POST /api/v1/me/bookmarks` with `{ tool_id }` to add bookmark
- Implement `DELETE /api/v1/me/bookmarks/:tool_id` to remove bookmark
- Bookmarks table: user_id (or session_id for anonymous), tool_id, created_at, unique constraint on (user_id, tool_id)
- Implement `GET /api/v1/compare?tool_ids=1,2,3,4` (or slugs) returning array of tool objects with comparison-relevant fields
- Update tools.bookmark_count on bookmark add/remove for social proof

**UX Integration (from UX Design):**
- **Bookmark interaction:** Bookmark icon/button on ToolCard and ToolProfile with toggle state (hollow → filled heart or bookmark icon)
- Pinia `useBookmarksStore`: Track bookmarked tool IDs, sync with backend, persist to localStorage for anonymous users
- Optimistic UI: Immediately update bookmark state, rollback on error
- **BookmarksView.vue:** Dedicated `/bookmarks` route showing shortlist
- Display bookmarked tools using same ToolCard component as results
- "Remove from shortlist" action on each card
- "Compare selected" button to send tools to comparison (multi-select checkboxes)
- Empty state: "You haven't bookmarked any tools yet. Start browsing!"
- **CompareView.vue:** `/compare` route showing side-by-side table for 2-4 tools
- **CompareTable.vue component:** Tools as columns, comparison rows:
  - Overview row: Name, logo, "best for", tagline
  - Features row: Key features comparison
  - Pricing row: Free tier, starting price, billing model
  - Ratings row: Overall rating, review count, dimension ratings as mini bars
  - Social proof row: Bookmark count, badges
- "Remove from comparison" action on each column
- "Visit tool" CTA for each column
- Responsive: Horizontal scroll on mobile, full table on desktop
- Shareable URL with tool slugs: `/compare?tools=chatgpt,claude,gemini`

**Dependencies:** Epic 1 (Foundation), Epic 2 (Discovery - bookmark from results), Epic 3 (Tool Profiles - bookmark from profile)

---

### Epic 6: User Accounts & Persistence

**User Value Statement:** Users can create accounts or sign in to persist their bookmarks and reviews across sessions and devices, eliminating the risk of losing their work.

**PRD Coverage:** FR19, FR20 (User Accounts & Sessions)

**Technical Context (from Architecture):**
- Implement `POST /api/v1/auth/register` with email + password, bcrypt hashing
- Implement `POST /api/v1/auth/login` returning JWT token in HTTP-only cookie
- Implement `POST /api/v1/auth/logout` clearing cookie
- Implement `GET /api/v1/me` returning current user profile
- Implement `GET /api/v1/me/reviews` returning user's submitted reviews
- Users table: id, email (unique), password_hash, display_name, role (user/admin), created_at
- JWT middleware in `internal/platform/http/middleware.go` validating tokens on protected routes
- Auth service in `internal/auth` handling registration, login, token generation
- Migration from anonymous session bookmarks to user bookmarks on account creation/login

**UX Integration (from UX Design):**
- **Header navigation:** Show "Sign In / Sign Up" buttons when not authenticated
- When authenticated: Show user avatar/initials + display name, dropdown menu with "My Reviews", "My Bookmarks", "Sign Out"
- **Sign up/Sign in modals or pages:** Simple forms with email + password
- Validation: Email format, password requirements (8+ chars, complexity if needed)
- Error handling: Clear messages for duplicate email, incorrect password, etc.
- Success: Auto-close modal, show success toast, update header UI
- **My Account area (optional dedicated view):**
  - Tabs or sections for "My Reviews" and "My Bookmarks"
  - My Reviews: List of submitted reviews with links to tool profiles
  - My Bookmarks: Same as BookmarksView but in account context
- Persist Pinia session store (useSessionStore) with current user data
- Anonymous bookmark migration: On login, prompt "You have X bookmarks. Save them to your account?" and merge

**Dependencies:** Epic 1 (Foundation - JWT auth setup), Epic 4 (Reviews - for persistent review authorship), Epic 5 (Bookmarks - for persistent bookmarks)

---

### Epic 7: Admin & Catalog Management

**User Value Statement:** Admins can maintain a high-quality, accurate, and complete tool directory through CRUD operations, taxonomy management, and data quality checks - ensuring users always find trustworthy information.

**PRD Coverage:** FR21, FR22, FR23, FR24 (Admin & Catalog Management), FR29, FR30 (Analytics)

**Technical Context (from Architecture):**
- Implement admin-only endpoints (protected by role check in middleware):
  - `GET/POST/PATCH/DELETE /api/v1/admin/tools` - manage tool records
  - `GET/POST/PATCH/DELETE /api/v1/admin/categories` - manage categories
  - `GET/POST/PATCH/DELETE /api/v1/admin/tags` - manage tags
  - `POST /api/v1/admin/tools/:id/badges` - assign badges
  - `GET /api/v1/admin/catalog/issues` - surface tools with missing fields, broken links
  - `POST /api/v1/admin/tools/:id/merge` - merge duplicate tools
- Analytics endpoints:
  - `GET /api/v1/admin/analytics/overview` - page views, bookmark counts, review counts over time
  - `GET /api/v1/admin/analytics/top-tools` - top tools by views, bookmarks
  - `GET /api/v1/admin/analytics/top-categories` - top categories by activity
- Admin service layer validating completeness (e.g., tools must have name, slug, description, best_for)
- Soft delete pattern: tools.archived_at for deprecation without data loss

**UX Integration (from UX Design):**
- **Admin UI (separate route or section, e.g., `/admin`):**
  - Dashboard showing catalog health metrics: Total tools, tools missing key fields, recent submissions
  - **Tools management:** Table with search, filter, sort; "Add Tool", "Edit", "Archive" actions
  - **Tool form:** Full CRUD form with all fields (name, slug, tagline, description, best_for, pricing, categories, tags, media URLs)
  - **Categories/Tags management:** Simple lists with add/edit/deactivate actions
  - **Badges management:** Assign/remove badges like "Top in Category", "Editor's Pick" to tools
  - **Catalog issues view:** List of tools with missing fields or broken media, direct links to edit
  - **Analytics dashboard:** Simple charts/tables showing engagement over time, top performers
- Admin nav: Additional menu items visible only to admin role
- Confirmation modals for destructive actions (archive, merge, delete)
- Success/error toasts for all admin operations
- RTL-aware admin UI using same design system

**Dependencies:** Epic 1 (Foundation), Epic 6 (Accounts - admin role), Epic 2-5 (user-facing features that admins manage)

---

### Epic 8: Moderation & Content Quality

**User Value Statement:** The platform maintains trusted, abuse-free user-generated content through reporting mechanisms and moderation workflows - ensuring the review layer remains valuable and credible.

**PRD Coverage:** FR25, FR26, FR27, FR28 (Moderation & Content Quality)

**Technical Context (from Architecture):**
- Implement reporting endpoints:
  - `POST /api/v1/tools/:slug/report` with reason (spam/abuse/misinformation/other) and optional comment
  - `POST /api/v1/reviews/:id/report` with same structure
- Implement moderation queue endpoints (admin/moderator role only):
  - `GET /api/v1/moderation/queue` with filters (type, status, date, tool, reason)
  - `PATCH /api/v1/moderation/reviews/:id/approve`
  - `PATCH /api/v1/moderation/reviews/:id/hide`
  - `PATCH /api/v1/moderation/reviews/:id/remove`
  - `GET /api/v1/moderation/history/:review_id` - audit log of actions
- Database tables:
  - `reports`: id, reportable_type (tool/review), reportable_id, reporter_user_id, reason, comment, status (pending/reviewed/dismissed), created_at
  - `moderation_actions`: id, actor_user_id, action_type (approve/hide/remove/edit), target_type (review/tool), target_id, notes, created_at
- Reviews table: add `moderation_status` (pending/approved/hidden/removed), `moderated_by`, `moderated_at`
- Auto-approve or manual review queue based on configuration

**UX Integration (from UX Design):**
- **Report action:** "Report" link/button on tool profiles and review cards
- **Report modal:** Simple form with reason dropdown (Spam, Abuse, Misinformation, Other) and optional comment textarea
- Success message: "Thank you for reporting. We'll review this soon."
- **Moderation UI (`/moderation` or admin section):**
  - **Queue view:** Table/list of reported items with filters (type, date, tool, reason)
  - Each entry shows: Content preview, reporter info (if available), report reason, submission date
  - Actions: "View in context" (link to tool/profile), "Approve", "Hide", "Remove", "Dismiss report"
  - Bulk actions for handling multiple reports
  - **History view:** Audit log showing who took what action when, with notes
- Moderator role: Similar to admin but focused on content review, not catalog management
- Clear status indicators: Pending (yellow), Approved (green), Hidden (gray), Removed (red)
- Confirmation modals for actions with impact (remove, hide)

**Dependencies:** Epic 1 (Foundation), Epic 4 (Reviews - content to moderate), Epic 6 (Accounts - moderator role)

---

## FR Coverage Map

| FR | Functional Requirement | Epic | Stories |
|----|----------------------|------|---------|
| FR1 | Free-text search for AI tools | Epic 2: Tool Discovery & Browsing | 2.2 |
| FR2 | Browse tools by category | Epic 2: Tool Discovery & Browsing | 2.1 |
| FR3 | Refine results using filters | Epic 2: Tool Discovery & Browsing | 2.3 |
| FR4 | Sort results | Epic 2: Tool Discovery & Browsing | 2.4 |
| FR5 | Open tool detail page | Epic 2: Tool Discovery & Browsing | 2.5 |
| FR6 | View tool profile with metadata | Epic 3: Rich Tool Profiles | 3.1 |
| FR7 | View rich media on profile | Epic 3: Rich Tool Profiles | 3.2 |
| FR8 | View structured reviews | Epic 4: Reviews & Ratings | 4.1 |
| FR9 | Submit structured review (logged-in) | Epic 4: Reviews & Ratings | 4.2 |
| FR10 | View aggregated ratings | Epic 4: Reviews & Ratings | 4.3 |
| FR11 | View social proof indicators | Epic 3: Rich Tool Profiles | 3.1 |
| FR12 | Add tools to comparison set | Epic 5: Comparison & Shortlists | 5.3 |
| FR13 | View comparison page | Epic 5: Comparison & Shortlists | 5.4 |
| FR14 | View alternatives/similar tools | Epic 3: Rich Tool Profiles | 3.3 |
| FR15 | Bookmark tool to shortlist | Epic 5: Comparison & Shortlists | 5.1 |
| FR16 | View bookmarks/shortlist | Epic 5: Comparison & Shortlists | 5.2 |
| FR17 | Remove tools from shortlist | Epic 5: Comparison & Shortlists | 5.2 |
| FR18 | Send bookmarks to comparison | Epic 5: Comparison & Shortlists | 5.2, 5.4 |
| FR19 | Create account / sign in | Epic 6: User Accounts & Persistence | 6.1 |
| FR20 | View own reviews and bookmarks | Epic 6: User Accounts & Persistence | 6.2 |
| FR21 | Admin manage tool records | Epic 7: Admin & Catalog Management | 7.1 |
| FR22 | Admin manage categories/tags | Epic 7: Admin & Catalog Management | 7.2 |
| FR23 | Admin manage badges | Epic 7: Admin & Catalog Management | 7.3 |
| FR24 | Admin resolve catalog issues | Epic 7: Admin & Catalog Management | 7.4 |
| FR25 | Report tools or reviews | Epic 8: Moderation & Content Quality | 8.1 |
| FR26 | Moderator view queue | Epic 8: Moderation & Content Quality | 8.2 |
| FR27 | Moderator approve/hide/remove content | Epic 8: Moderation & Content Quality | 8.3 |
| FR28 | View moderation history | Epic 8: Moderation & Content Quality | 8.4 |
| FR29 | Admin view engagement signals | Epic 7: Admin & Catalog Management | 7.5 |
| FR30 | Admin view top tools/categories | Epic 7: Admin & Catalog Management | 7.5 |

---

## Epic 1: Foundation & Core Infrastructure

**Goal:** Establish the complete technical foundation (backend, frontend, database, auth) that enables all subsequent user-facing features.

### Story 1.1: Backend Project Initialization

**As a** developer,
**I want** the Go backend project initialized with Gin, GORM, and PostgreSQL configured,
**So that** I can start building API endpoints on a solid foundation.

**Acceptance Criteria:**

**Given** I need to set up the backend project
**When** I initialize the project structure
**Then** the following is complete:

- Go module created at `backend/` with `go.mod` defining module path `github.com/your-org/ai-tools-atlas-backend`
- Dependencies installed: `gin-gonic/gin`, `gorm.io/gorm`, `gorm.io/driver/postgres`, `golang-jwt/jwt`, `golang-migrate/migrate`
- Project structure matches Architecture spec:
  ```
  backend/
  ├── cmd/api/main.go
  ├── internal/
  │   ├── platform/
  │   │   ├── config/config.go
  │   │   ├── db/db.go
  │   │   └── http/
  │   │       ├── router.go
  │   │       ├── middleware.go
  │   │       └── responses.go
  ├── migrations/
  ├── .env.example
  └── Makefile
  ```
- `config.go` loads environment variables: `DATABASE_URL`, `JWT_SECRET`, `PORT`, `ALLOWED_ORIGINS`
- `db.go` establishes PostgreSQL connection using GORM with connection pooling
- Health check endpoint `GET /health` returns `{ "status": "ok" }`
- Server starts on configured port and responds to health checks

**Technical Implementation:**

- Use `github.com/joho/godotenv` for `.env` loading in development
- GORM connection: `gorm.Open(postgres.Open(dsn), &gorm.Config{})`
- Gin router setup in `router.go` with versioned prefix `/api/v1`
- Response helpers in `responses.go`:
  - `SuccessResponse(c *gin.Context, data interface{})`
  - `ListResponse(c *gin.Context, data interface{}, meta map[string]interface{})`
  - `ErrorResponse(c *gin.Context, code string, message string, details interface{})`
- Middleware setup: Recovery, CORS, Request logging

**Prerequisites:** None

**Files Created:**
- `backend/cmd/api/main.go`
- `backend/internal/platform/config/config.go`
- `backend/internal/platform/db/db.go`
- `backend/internal/platform/http/router.go`
- `backend/internal/platform/http/middleware.go`
- `backend/internal/platform/http/responses.go`
- `backend/.env.example`

---

### Story 1.2: Database Schema & Migrations

**As a** developer,
**I want** the complete PostgreSQL database schema defined and migration system set up,
**So that** all entities from the data model are ready for use.

**Acceptance Criteria:**

**Given** I need the database schema from Architecture.md
**When** I create and run migrations
**Then** the following tables exist with proper schema:

**Core Tables:**
- `tools`: id, slug (unique), name, logo_url, tagline, description, best_for, primary_use_cases (text), pricing_summary, target_roles (text), platforms (text), has_free_tier (bool), official_url, primary_category_id (FK), avg_rating_overall (decimal), review_count (int), bookmark_count (int), trending_score (decimal), created_at, updated_at, archived_at (nullable)
- `categories`: id, slug (unique), name, description, icon_url, display_order (int), created_at, updated_at
- `tags`: id, slug (unique), name, created_at
- `tool_tags`: tool_id (FK), tag_id (FK), PRIMARY KEY (tool_id, tag_id)
- `media`: id, tool_id (FK), type (enum: screenshot/video), url, thumbnail_url, display_order, created_at
- `badges`: id, slug (unique), name, description, icon_url, created_at
- `tool_badges`: tool_id (FK), badge_id (FK), assigned_at, PRIMARY KEY (tool_id, badge_id)
- `tool_alternatives`: id, tool_id (FK), alternative_tool_id (FK), relationship_type (enum: similar/alternative), created_at
- `users`: id, email (unique), password_hash, display_name, role (enum: user/admin/moderator), created_at, updated_at
- `reviews`: id, tool_id (FK), user_id (FK), rating_overall (int 1-5), rating_ease_of_use, rating_value, rating_accuracy, rating_speed, rating_support (all int 1-5), pros (text), cons (text), primary_use_case, reviewer_role, company_size, usage_context (text), helpful_count (int default 0), moderation_status (enum), moderated_by (FK nullable), moderated_at (nullable), created_at, updated_at
- `bookmarks`: id, user_id (FK), tool_id (FK), created_at, UNIQUE(user_id, tool_id)

**Indexes Created:**
- `idx_tools_slug` on `tools.slug`
- `idx_tools_primary_category_id` on `tools.primary_category_id`
- `idx_tools_avg_rating` on `tools.avg_rating_overall`
- `idx_tools_bookmark_count` on `tools.bookmark_count`
- `idx_tools_trending_score` on `tools.trending_score`
- `idx_categories_slug` on `categories.slug`
- `idx_reviews_tool_id` on `reviews.tool_id`
- `idx_reviews_user_id` on `reviews.user_id`
- `idx_bookmarks_user_tool` on `bookmarks(user_id, tool_id)`

**And:**
- Migrations managed via `golang-migrate` with timestamped files in `backend/migrations/`
- Makefile targets: `make migrate-up`, `make migrate-down`, `make migrate-create name=<name>`
- Foreign key constraints enforced with appropriate ON DELETE behavior
- All naming follows snake_case convention from Architecture.md

**Technical Implementation:**

- Use `golang-migrate/migrate/v4` with PostgreSQL driver
- Migration files: `000001_initial_schema.up.sql` and `000001_initial_schema.down.sql`
- Connection string from env: `DATABASE_URL=postgres://user:pass@localhost:5432/ai_tools_atlas?sslmode=disable`
- Run migrations on app startup or via separate command/Makefile target

**Prerequisites:** Story 1.1 (Backend Project Initialization)

**Files Created:**
- `backend/migrations/000001_initial_schema.up.sql`
- `backend/migrations/000001_initial_schema.down.sql`
- `backend/Makefile` (migration targets)

---

### Story 1.3: GORM Models

**As a** developer,
**I want** GORM model structs defined for all entities,
**So that** I can perform type-safe database operations.

**Acceptance Criteria:**

**Given** I have the database schema from Story 1.2
**When** I create GORM model files
**Then** the following models are defined:

**In `internal/tools/model.go`:**
```go
type Tool struct {
    ID                 uint      `gorm:"primaryKey"`
    Slug               string    `gorm:"uniqueIndex;not null"`
    Name               string    `gorm:"not null"`
    LogoURL            string    `gorm:"column:logo_url"`
    Tagline            string
    Description        string    `gorm:"type:text"`
    BestFor            string    `gorm:"column:best_for;type:text"`
    PrimaryUseCases    string    `gorm:"type:text"`
    PricingSummary     string
    TargetRoles        string    `gorm:"type:text"`
    Platforms          string    `gorm:"type:text"`
    HasFreeTier        bool      `gorm:"default:false"`
    OfficialURL        string    `gorm:"column:official_url"`
    PrimaryCategoryID  uint
    PrimaryCategory    Category  `gorm:"foreignKey:PrimaryCategoryID"`
    AvgRatingOverall   float64   `gorm:"column:avg_rating_overall;default:0"`
    ReviewCount        int       `gorm:"default:0"`
    BookmarkCount      int       `gorm:"default:0"`
    TrendingScore      float64   `gorm:"default:0"`
    Tags               []Tag     `gorm:"many2many:tool_tags"`
    Media              []Media   `gorm:"foreignKey:ToolID"`
    Badges             []Badge   `gorm:"many2many:tool_badges"`
    CreatedAt          time.Time
    UpdatedAt          time.Time
    ArchivedAt         *time.Time `gorm:"index"`
}
```

**Similar models for:**
- `Category` in `internal/categories/model.go`
- `Tag` in `internal/tags/model.go` (future, can be in tools package initially)
- `Media` in `internal/tools/model.go` or separate media package
- `Badge` in `internal/badges/model.go`
- `User` in `internal/auth/model.go`
- `Review` in `internal/reviews/model.go`
- `Bookmark` in `internal/bookmarks/model.go`
- `ToolAlternative` in `internal/tools/model.go`

**And:**
- All models use proper GORM tags matching database schema
- JSON tags added for API serialization using snake_case: `json:"slug"`
- Relationships configured with appropriate foreign keys and constraints
- Auto-migration can run successfully: `db.AutoMigrate(&Tool{}, &Category{}, ...)` (for development)

**Technical Implementation:**

- PascalCase for struct names, camelCase for field names (Go convention)
- `gorm:` tags match database column names in snake_case
- `json:` tags use snake_case to match API response format
- Use pointers for nullable fields (`*time.Time`, `*string`)
- Omit empty fields in JSON with `json:"field,omitempty"` where appropriate

**Prerequisites:** Story 1.2 (Database Schema)

**Files Created:**
- `backend/internal/tools/model.go`
- `backend/internal/categories/model.go`
- `backend/internal/reviews/model.go`
- `backend/internal/bookmarks/model.go`
- `backend/internal/badges/model.go`
- `backend/internal/auth/model.go`

---

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

### Story 1.5: Frontend Project Initialization

**As a** developer,
**I want** the Vue 3 frontend project initialized with Vite, TypeScript, Pinia, and Tailwind,
**So that** I can start building UI components and views.

**Acceptance Criteria:**

**Given** I need the frontend project structure
**When** I initialize the Vue 3 project
**Then** the following is complete:

**Project Created:**
- `npm create vite@latest frontend -- --template vue-ts` executed
- Dependencies installed: `vue`, `vue-router@4`, `pinia`, `tailwindcss`, `postcss`, `autoprefixer`
- Project structure:
  ```
  frontend/
  ├── src/
  │   ├── main.ts
  │   ├── App.vue
  │   ├── router/index.ts
  │   ├── stores/
  │   ├── views/
  │   ├── components/
  │   ├── lib/
  │   └── assets/
  ├── index.html
  ├── vite.config.ts
  ├── tailwind.config.cjs
  ├── postcss.config.cjs
  ├── tsconfig.json
  ├── .env.example
  └── package.json
  ```

**Tailwind Configuration (`tailwind.config.cjs`):**
- Dark theme colors configured:
  ```js
  theme: {
    extend: {
      colors: {
        'dark-bg': '#05060A',
        'dark-surface': '#0A0B10',
        'dark-border': '#1F2933',
        'primary': {
          500: '#3B82F6',
          600: '#2563EB',
        },
      },
    },
  },
  ```
- RTL plugin configured: `require('tailwindcss-rtl')`
- Content paths: `./index.html`, `./src/**/*.{vue,js,ts,jsx,tsx}`

**Router Setup (`src/router/index.ts`):**
- Vue Router 4 configured with routes:
  - `/` - HomeView
  - `/tools/:slug` - ToolProfileView
  - `/categories/:slug` - SearchResultsView (category filter)
  - `/search` - SearchResultsView
  - `/compare` - CompareView
  - `/bookmarks` - BookmarksView
- History mode enabled

**Pinia Setup (`src/main.ts`):**
- Pinia initialized and registered
- App mounted to `#app`

**And:**
- TypeScript configured with Vue 3 support
- Vite configured to proxy API requests to backend (default: `http://localhost:8080`)
- Dev server starts on port 5173
- `.env.example` with `VITE_API_BASE_URL=http://localhost:8080/api/v1`

**Technical Implementation:**

- Install commands:
  ```bash
  cd frontend
  npm install pinia vue-router@4
  npm install -D tailwindcss postcss autoprefixer tailwindcss-rtl
  npx tailwindcss init -p
  ```
- Vite config proxy:
  ```ts
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
  ```

**Prerequisites:** None (parallel to backend setup)

**Files Created:**
- `frontend/src/main.ts`
- `frontend/src/router/index.ts`
- `frontend/tailwind.config.cjs`
- `frontend/vite.config.ts`
- `frontend/.env.example`

---

### Story 1.6: Base Layout Components & Dark Theme

**As a** developer,
**I want** base layout components and dark theme styling configured,
**So that** all pages have consistent Arabic-first, RTL-aware layout.

**Acceptance Criteria:**

**Given** I need base UI components matching UX Design spec
**When** I create layout components
**Then** the following components exist:

**`src/components/layout/AppShell.vue`:**
- Root layout component with `dir="rtl"` for Arabic
- Dark background using Tailwind classes: `bg-dark-bg text-white`
- Contains: `<HeaderNav />`, `<router-view />`, optional footer
- Full viewport height layout

**`src/components/layout/HeaderNav.vue`:**
- Top navigation bar with dark surface background
- Logo/brand on right (RTL), navigation links in center, user menu on left
- Navigation items: "الرئيسية" (Home), "تصفح" (Browse), "المفضلة" (Bookmarks)
- User menu (when authenticated): Avatar/initials, dropdown with "مراجعاتي" (My Reviews), "تسجيل الخروج" (Sign Out)
- User menu (not authenticated): "تسجيل الدخول" (Sign In), "إنشاء حساب" (Sign Up) buttons
- Mobile: Hamburger menu (RTL-aware icon) with drawer navigation

**`src/components/common/SkeletonCard.vue`:**
- Generic skeleton loader for cards
- Pulse animation, dark theme colors
- Props: `height`, `width`, `rounded`

**`src/components/common/Spinner.vue`:**
- Loading spinner with neon blue color
- Used for page-level loading states

**Global Styles (`src/assets/main.css`):**
- Tailwind directives imported
- Cairo font loaded (from Google Fonts or local)
- Base typography styles for Arabic: proper line-height, letter-spacing
- RTL-specific adjustments for common elements

**And:**
- `App.vue` uses `<AppShell>` as root component
- All text uses Cairo font
- Color scheme matches UX spec: Dark bg (#05060A), neon blue primary (#3B82F6)
- RTL layout tested: elements flow right-to-left naturally

**Technical Implementation:**

- Cairo font:
  ```css
  @import url('https://fonts.googleapis.com/css2?family=Cairo:wght@400;600;700&display=swap');

  body {
    font-family: 'Cairo', sans-serif;
    direction: rtl;
  }
  ```
- Use Tailwind RTL utilities or manual RTL adjustments
- Skeleton animation:
  ```css
  @keyframes pulse {
    0%, 100% { opacity: 1; }
    50% { opacity: 0.5; }
  }
  ```

**Prerequisites:** Story 1.5 (Frontend Project Initialization)

**Files Created:**
- `frontend/src/components/layout/AppShell.vue`
- `frontend/src/components/layout/HeaderNav.vue`
- `frontend/src/components/common/SkeletonCard.vue`
- `frontend/src/components/common/Spinner.vue`
- `frontend/src/assets/main.css`

---

### Story 1.7: API Client & Pinia Stores Setup

**As a** developer,
**I want** a centralized API client and base Pinia stores configured,
**So that** all components can make authenticated API calls and manage state consistently.

**Acceptance Criteria:**

**Given** I need to call backend APIs from frontend
**When** I create the API client and stores
**Then** the following is complete:

**API Client (`src/lib/apiClient.ts`):**
```typescript
const BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api/v1';

export interface ApiResponse<T> {
  data: T;
  meta?: {
    page: number;
    page_size: number;
    total: number;
  };
}

export interface ApiError {
  code: string;
  message: string;
  details?: any;
}

export const apiClient = {
  async get<T>(endpoint: string, params?: Record<string, any>): Promise<T>,
  async post<T>(endpoint: string, body?: any): Promise<T>,
  async patch<T>(endpoint: string, body?: any): Promise<T>,
  async delete<T>(endpoint: string): Promise<T>,
}
```
- All methods include credentials for cookie-based auth
- Error handling with typed ApiError
- Automatically includes `Content-Type: application/json`

**Session Store (`src/stores/session.ts`):**
```typescript
export const useSessionStore = defineStore('session', {
  state: () => ({
    user: null as User | null,
    isAuthenticated: false,
    isLoading: false,
  }),
  actions: {
    async fetchCurrentUser(),
    async login(email: string, password: string),
    async register(email: string, password: string, displayName: string),
    async logout(),
  },
})
```

**Bookmarks Store (`src/stores/bookmarks.ts`):**
```typescript
export const useBookmarksStore = defineStore('bookmarks', {
  state: () => ({
    bookmarkedToolIds: [] as number[],
  }),
  getters: {
    isBookmarked: (state) => (toolId: number) => state.bookmarkedToolIds.includes(toolId),
  },
  actions: {
    async fetchBookmarks(),
    async addBookmark(toolId: number),
    async removeBookmark(toolId: number),
  },
  persist: true, // localStorage persistence for anonymous users
})
```

**Filters Store (`src/stores/filters.ts`):**
```typescript
export const useFiltersStore = defineStore('filters', {
  state: () => ({
    query: '',
    category: null as string | null,
    price: null as string | null,
    minRating: null as number | null,
    platform: null as string | null,
    sort: 'top_rated' as string,
  }),
  actions: {
    setFilter(key: string, value: any),
    clearFilters(),
  },
})
```

**And:**
- TypeScript interfaces defined for API types: `User`, `Tool`, `Category`, etc.
- Error handling in stores with user-friendly error messages
- Loading states managed in stores

**Technical Implementation:**

- Use `fetch` API or `axios` for HTTP requests
- Pinia store composition API or options API (consistent choice)
- Install `pinia-plugin-persistedstate` for bookmark persistence:
  ```bash
  npm install pinia-plugin-persistedstate
  ```
- API client error handling:
  ```typescript
  if (!response.ok) {
    const error: ApiError = await response.json();
    throw new Error(error.message);
  }
  ```

**Prerequisites:** Story 1.5 (Frontend Project Initialization)

**Files Created:**
- `frontend/src/lib/apiClient.ts`
- `frontend/src/lib/types.ts` (TypeScript interfaces)
- `frontend/src/stores/session.ts`
- `frontend/src/stores/bookmarks.ts`
- `frontend/src/stores/filters.ts`

---

## Epic 2: Tool Discovery & Browsing

**Goal:** Users can quickly find AI tools relevant to their specific needs and use cases through intuitive hero search, category browsing, filtering, and sorting.

### Story 2.1: Category Browsing Backend

**As a** developer,
**I want** category listing and category-filtered tools endpoints implemented,
**So that** users can browse tools by category.

**Acceptance Criteria:**

**Given** I need category browsing functionality
**When** I implement the categories API
**Then** the following endpoints work:

**Endpoint: `GET /api/v1/categories`**
- Returns list of all categories with fields: id, slug, name, description, icon_url, display_order
- Response format: `{ data: [{ id, slug, name, description, icon_url, display_order }] }`
- Ordered by display_order ASC
- Only active (non-archived) categories returned

**Endpoint: `GET /api/v1/categories/:slug/tools`**
- Accepts category slug in URL path
- Returns paginated list of tools in that category
- Response format: `{ data: [...tools], meta: { page, page_size, total } }`
- Each tool includes: id, slug, name, logo_url, tagline, best_for, avg_rating_overall, review_count, bookmark_count, pricing_summary, has_free_tier, tags (array)
- Default pagination: page_size=20
- Query params: `?page=1&page_size=20`
- 404 error if category slug not found

**And:**
- Repository layer in `internal/categories/repository.go` with functions:
  - `ListCategories() ([]Category, error)`
  - `GetCategoryBySlug(slug string) (*Category, error)`
  - `ListToolsByCategory(categoryID uint, page, pageSize int) ([]Tool, int, error)`
- Service layer in `internal/categories/service.go` handling business logic
- Handler in `internal/categories/handler.go` using standardized responses
- Proper error handling with appropriate HTTP status codes

**Technical Implementation:**

- GORM queries with preloading for tags: `.Preload("Tags").Preload("PrimaryCategory")`
- Pagination: `.Limit(pageSize).Offset((page - 1) * pageSize)`
- Filter out archived tools: `WHERE archived_at IS NULL`
- Count query for total: `db.Model(&Tool{}).Where("primary_category_id = ?", categoryID).Count(&total)`

**Prerequisites:** Epic 1 (Foundation, GORM models)

**Files Created:**
- `backend/internal/categories/repository.go`
- `backend/internal/categories/service.go`
- `backend/internal/categories/handler.go`

---

### Story 2.2: Tools Listing & Search Backend

**As a** developer,
**I want** tools listing and free-text search endpoints implemented,
**So that** users can search and browse all tools.

**Acceptance Criteria:**

**Given** I need tools listing and search
**When** I implement the tools API
**Then** the following endpoints work:

**Endpoint: `GET /api/v1/tools`**
- Returns paginated list of all tools
- Response format: `{ data: [...tools], meta: { page, page_size, total } }`
- Each tool includes: id, slug, name, logo_url, tagline, best_for, avg_rating_overall, review_count, bookmark_count, pricing_summary, has_free_tier, primary_category (nested object), tags (array), badges (array)
- Query params:
  - `?page=1&page_size=20` (pagination)
  - `?category=slug` (filter by category)
  - `?price=free|freemium|paid` (filter by pricing)
  - `?min_rating=4` (filter by minimum rating)
  - `?platform=web|mobile|api` (filter by platform - search in platforms text field)
  - `?sort=top_rated|most_bookmarked|trending|newest` (sort order)
- Default sort: top_rated
- Only non-archived tools returned

**Endpoint: `GET /api/v1/search/tools`**
- Free-text search across tool name, tagline, description, best_for, primary_use_cases
- Query params: `?q=search+query` plus all filters from GET /tools
- Returns same format as GET /tools with matching tools
- Search is case-insensitive
- Empty query returns all tools (same as GET /tools)

**Sorting Logic:**
- `top_rated`: ORDER BY avg_rating_overall DESC, review_count DESC
- `most_bookmarked`: ORDER BY bookmark_count DESC
- `trending`: ORDER BY trending_score DESC (trending_score calculated separately)
- `newest`: ORDER BY created_at DESC

**And:**
- Repository in `internal/tools/repository.go`:
  - `ListTools(filters ToolFilters, page, pageSize int) ([]Tool, int, error)`
  - `SearchTools(query string, filters ToolFilters, page, pageSize int) ([]Tool, int, error)`
- Service in `internal/tools/service.go`
- Handler in `internal/tools/handler.go`
- ToolFilters struct: `{ Category, Price, MinRating, Platform, Sort }`

**Technical Implementation:**

- Full-text search using PostgreSQL ILIKE or tsvector for better performance:
  ```sql
  WHERE name ILIKE '%query%' OR tagline ILIKE '%query%' OR description ILIKE '%query%' OR best_for ILIKE '%query%'
  ```
- Preload relationships: `.Preload("PrimaryCategory").Preload("Tags").Preload("Badges")`
- Dynamic query building based on filters:
  ```go
  query := db.Model(&Tool{}).Where("archived_at IS NULL")
  if filters.Category != "" {
      query = query.Joins("JOIN categories ON categories.id = tools.primary_category_id").Where("categories.slug = ?", filters.Category)
  }
  if filters.MinRating > 0 {
      query = query.Where("avg_rating_overall >= ?", filters.MinRating)
  }
  // ... etc
  ```

**Prerequisites:** Epic 1, Story 2.1 (Categories)

**Files Created:**
- `backend/internal/tools/repository.go`
- `backend/internal/tools/service.go`
- `backend/internal/tools/handler.go`

---

### Story 2.3: Home Page with Hero Search

**As a** user,
**I want** to see a prominent search bar on the home page,
**So that** I can immediately search for AI tools by describing my situation.

**Acceptance Criteria:**

**Given** I visit the home page
**When** the page loads
**Then** I see:

**Hero Section:**
- Large, centered hero search input with Arabic prompt text: "اكتب وضعك أو ما تريد القيام به…" (Type your situation or what you want to do...)
- Search input is prominent, neon blue border on focus, dark background
- Placeholder text helps users understand they can use free-text (e.g., "محامٍ يريد تلخيص العقود")
- Search button or press Enter to submit
- Optional suggestion chips below input: "محامٍ", "تسويق بالمحتوى", "برمجة" as examples

**Category Grid:**
- Below hero: grid of top-level categories (3-4 columns on desktop, 1-2 on mobile)
- Each category card shows: icon, name (Arabic), tool count
- Cards are clickable, navigate to `/categories/:slug`
- Skeleton loaders while categories load

**Popular/Featured Strip (Optional for MVP):**
- Horizontal strip of "Popular this week" or "Top rated" tools
- Shows 4-6 tool cards in horizontal scroll
- Uses ToolCard component (will be created in next story)

**And:**
- **When** I type into the search input and press Enter
- **Then** I navigate to `/search?q=my+query` with the SearchResultsView

**UX Requirements (from UX Design spec):**
- Dark background (#05060A)
- Hero search is visually dominant, centered, large font size
- Cairo font for all Arabic text
- RTL layout: categories grid flows right-to-left
- Responsive: hero search full-width on mobile, max-width on desktop
- Page loads with skeleton states, no blocking spinners

**Technical Implementation:**

- **HomeView.vue** (`frontend/src/views/HomeView.vue`)
- Fetch categories on mount: `apiClient.get('/categories')`
- Search form submit handler navigates: `router.push({ path: '/search', query: { q: searchQuery } })`
- Use `useFiltersStore` to store current search query
- Skeleton loaders for category grid while loading
- Composable for loading state: `const { data: categories, loading, error } = useFetch('/categories')`

**Prerequisites:** Epic 1 (Frontend setup, API client), Story 2.1 (Categories API)

**Files Created:**
- `frontend/src/views/HomeView.vue`
- `frontend/src/components/home/HeroSearch.vue` (optional sub-component)
- `frontend/src/components/home/CategoryGrid.vue`

---

### Story 2.4: Tool Card Component

**As a** developer,
**I want** a reusable ToolCard component,
**So that** tools can be displayed consistently across search results, categories, and bookmarks.

**Acceptance Criteria:**

**Given** I need to display tool information in card format
**When** I create the ToolCard component
**Then** it displays the following:

**Card Layout (matches G2-style from UX Design):**
- Tool logo (left in RTL) with fallback placeholder if missing
- Tool name as heading
- "Best for" line in lighter text (one-line truncate)
- Rating display: stars (or numeric) + review count in parentheses "(23 مراجعة)"
- Pricing badge: "مجاني" (Free), "Freemium", or starting price
- Tags: 2-3 primary tags as small pills (overflow hidden)
- Bookmark button (hollow heart/bookmark icon, filled when bookmarked)
- "مقارنة" (Compare) button/checkbox to add to comparison
- Clicking card (except buttons) navigates to `/tools/:slug`

**Visual Design (from UX Design spec):**
- Dark card background (#0A0B10) with subtle border (#1F2933)
- Hover state: slight border glow (neon blue)
- RTL-aware: logo on right, actions on left
- Responsive: full-width on mobile, fixed width/flex-basis on desktop grid
- Cairo font for all text
- Proper spacing (8px base unit from Architecture)

**Component Props:**
```typescript
interface ToolCardProps {
  tool: {
    id: number;
    slug: string;
    name: string;
    logo_url?: string;
    best_for: string;
    avg_rating_overall: number;
    review_count: number;
    pricing_summary: string;
    has_free_tier: boolean;
    tags: Array<{ name: string }>;
  };
  showCompare?: boolean; // default true
}
```

**Bookmark Interaction:**
- Uses `useBookmarksStore` to check `isBookmarked(tool.id)`
- Click bookmark button: calls `addBookmark(tool.id)` or `removeBookmark(tool.id)`
- Optimistic UI: immediately updates icon state, rollback on error
- Shows toast on error: "فشل في حفظ الأداة" (Failed to save tool)

**Compare Interaction:**
- Click compare button: emits `@add-to-compare` event with tool
- Parent component handles comparison state (will be in ComparisonStore later)

**Technical Implementation:**

- Vue 3 SFC in `frontend/src/components/tools/ToolCard.vue`
- Use Tailwind for styling, no scoped CSS unless necessary
- Rating display: create small `RatingStars.vue` component or use numeric display
- Icons: use a lightweight icon library or inline SVGs for bookmark/compare
- Truncate "best for" text: `class="truncate"` or `-webkit-line-clamp: 1`

**Prerequisites:** Epic 1 (Frontend setup), Story 2.2 (Tools API for data structure)

**Files Created:**
- `frontend/src/components/tools/ToolCard.vue`
- `frontend/src/components/common/RatingStars.vue` (optional)

---

### Story 2.5: Search Results View with Filters

**As a** user,
**I want** to see search results with filtering and sorting options,
**So that** I can narrow down to the most relevant tools for my needs.

**Acceptance Criteria:**

**Given** I navigate to `/search?q=my+query` or `/categories/:slug`
**When** the SearchResultsView loads
**Then** I see:

**Page Layout:**
- **Filter panel (right side in RTL, drawer on mobile):**
  - Categories filter: dropdown or list of checkboxes
  - Price filter: radio buttons (الكل، مجاني، Freemium، مدفوع)
  - Rating filter: slider or radio (All, 4+, 4.5+)
  - Platform filter: checkboxes (Web, Mobile, API, Desktop)
  - "مسح الفلاتر" (Clear filters) button
- **Results area (main content):**
  - Header showing: query/category name, result count, sort dropdown
  - Sort dropdown: "الأعلى تقييماً" (Top rated), "الأكثر حفظاً" (Most bookmarked), "الأحدث" (Newest), "الرائج" (Trending)
  - Grid of ToolCard components (3-4 columns desktop, 1 column mobile)
  - Pagination controls at bottom (Previous/Next, page numbers)
  - Empty state if no results: "لم تجد ما تريد؟ استعرض حسب الفئة" with links to categories

**Filter Behavior:**
- Changing any filter updates URL query params: `?category=...&price=...&min_rating=...&sort=...`
- URL params drive the API call and filter UI state
- Page resets to 1 when filters change
- Filters are visible/sticky on desktop, drawer on mobile with "تصفية" (Filter) button

**Loading States:**
- Skeleton loaders for tool cards while loading (show 8-12 skeletons)
- No blocking spinner, results area shows skeletons in grid layout
- Preserve previous results briefly while new results load (optional optimization)

**Search Query Handling:**
- Free-text query from hero search calls `GET /api/v1/search/tools?q=...`
- Category browsing calls `GET /api/v1/categories/:slug/tools`
- All other filters apply to both search and category views

**And:**
- Uses `useFiltersStore` to sync filter state with URL params
- Filter changes trigger new API call: `apiClient.get('/tools', { params: filters })`
- Pagination: `&page=2` in query params

**Technical Implementation:**

- **SearchResultsView.vue** (`frontend/src/views/SearchResultsView.vue`)
- Watch route query params, fetch tools when params change:
  ```typescript
  watch(() => route.query, () => {
    fetchTools();
  }, { immediate: true });
  ```
- Filter panel: separate component `FiltersPanel.vue` in `components/search/`
- Responsive filter drawer using Headless UI or custom modal on mobile
- URL sync: use Vue Router to update query params when filters change
- Pagination: `<Pagination :current-page="page" :total-pages="totalPages" @page-change="goToPage" />`

**Prerequisites:** Story 2.2 (Tools API), Story 2.4 (ToolCard component)

**Files Created:**
- `frontend/src/views/SearchResultsView.vue`
- `frontend/src/components/search/FiltersPanel.vue`
- `frontend/src/components/common/Pagination.vue`

---

### Story 2.6: Category Browsing Frontend

**As a** user,
**I want** to click on a category and see tools in that category,
**So that** I can discover tools for a specific field or use case.

**Acceptance Criteria:**

**Given** I am on the home page
**When** I click a category card (e.g., "كتابة وتسويق" - Writing & Marketing)
**Then** I navigate to `/categories/:slug` and see:

**Category Header:**
- Category name and description as page title
- Icon (if available) displayed prominently
- Breadcrumb: "الرئيسية > [Category Name]"
- Tool count: "عدد الأدوات: 42"

**Results Area:**
- Reuses SearchResultsView component with category context
- Filter panel shows with category pre-selected (can change to other categories)
- All filters (price, rating, platform, sort) work as in Story 2.5
- URL structure: `/categories/writing-and-marketing?price=free&sort=top_rated`

**Navigation from Home:**
- CategoryGrid component emits `@category-select` with category slug
- Router navigates to `/categories/:slug`
- Alternatively, each category card is a router-link: `<router-link :to="`/categories/${category.slug}`">`

**And:**
- **When** I change the category filter in the filter panel
- **Then** the URL updates to the new category slug and results refresh

**Technical Implementation:**

- SearchResultsView checks `route.params.slug` to determine if it's a category view
- If category slug present: call `GET /api/v1/categories/:slug/tools`
- If no slug (search view): call `GET /api/v1/search/tools?q=...` or `GET /api/v1/tools`
- Fetch category details on mount: `GET /api/v1/categories/:slug` to display header info
- Loading state: show skeleton for header + tool cards

**Prerequisites:** Story 2.1 (Categories API), Story 2.3 (Home page), Story 2.5 (SearchResultsView)

**Files Modified:**
- `frontend/src/views/SearchResultsView.vue` (add category mode support)
- `frontend/src/components/home/CategoryGrid.vue` (add click handlers)

---

### Story 2.7: Navigate to Tool Profile

**As a** user,
**I want** to click on a tool card and view its detailed profile,
**So that** I can learn more about a specific tool.

**Acceptance Criteria:**

**Given** I see a tool card in search results or category view
**When** I click on the card (anywhere except bookmark/compare buttons)
**Then** I navigate to `/tools/:slug` and the ToolProfileView loads

**And:**
- Clicking the tool name also navigates to the profile
- Clicking the logo navigates to the profile
- Bookmark and compare buttons do NOT trigger navigation (event.stopPropagation)

**Router Configuration:**
- Route defined: `{ path: '/tools/:slug', component: ToolProfileView }`
- Route receives slug as param: `route.params.slug`

**ToolProfileView (Basic Shell):**
- For this story, create a minimal profile view showing:
  - Tool name as heading
  - "Loading full profile..." message or skeleton
  - This will be fully implemented in Epic 3
- Verifies navigation works end-to-end

**Technical Implementation:**

- ToolCard component:
  ```vue
  <div class="tool-card" @click="navigateToProfile">
    <!-- card content -->
    <button @click.stop="toggleBookmark">Bookmark</button>
    <button @click.stop="addToCompare">Compare</button>
  </div>
  ```
- navigateToProfile method: `router.push(\`/tools/${props.tool.slug}\`)`
- ToolProfileView.vue stub: `<div>Tool Profile: {{ route.params.slug }}</div>`

**Prerequisites:** Story 2.4 (ToolCard component)

**Files Created:**
- `frontend/src/views/ToolProfileView.vue` (basic shell for now, will be expanded in Epic 3)

---

## Epic 3: Rich Tool Profiles

**Goal:** Users can deeply evaluate individual tools with comprehensive information including metadata, rich media, and alternatives - enabling confident tool selection decisions.

### Story 3.1: Tool Profile Backend API

**As a** developer,
**I want** a complete tool profile endpoint with all related data,
**So that** users can view comprehensive tool information on the frontend.

**Acceptance Criteria:**

**Given** I need detailed tool profile data
**When** I implement the tool profile endpoint
**Then** the following endpoint works:

**Endpoint: `GET /api/v1/tools/:slug`**
- Accepts tool slug in URL path
- Returns complete tool object with all relationships loaded
- Response format: `{ data: { ...complete_tool_object } }`
- 404 error if tool slug not found or tool is archived

**Tool Object Includes:**
- **Core metadata:** id, slug, name, logo_url, tagline, description, best_for, primary_use_cases (parsed array or text), pricing_summary, target_roles (array or text), platforms (array or text), has_free_tier, official_url
- **Aggregated data:** avg_rating_overall, review_count, bookmark_count, trending_score
- **Relationships:**
  - `primary_category`: Full category object (id, slug, name, icon_url)
  - `tags`: Array of tag objects (id, slug, name)
  - `media`: Array of media objects (id, type, url, thumbnail_url, display_order) ordered by display_order
  - `badges`: Array of badge objects (id, slug, name, description, icon_url)
- **Timestamps:** created_at, updated_at

**And:**
- Repository function in `internal/tools/repository.go`:
  - `GetToolBySlug(slug string) (*Tool, error)` with all preloads
- Service function validates tool exists and is not archived
- Handler uses standardized success response

**Technical Implementation:**

- GORM query with comprehensive preloading:
  ```go
  db.Where("slug = ? AND archived_at IS NULL", slug).
     Preload("PrimaryCategory").
     Preload("Tags", func(db *gorm.DB) *gorm.DB {
         return db.Order("tags.name ASC")
     }).
     Preload("Media", func(db *gorm.DB) *gorm.DB {
         return db.Order("media.display_order ASC")
     }).
     Preload("Badges").
     First(&tool)
  ```
- Parse text fields into arrays if stored as comma-separated or JSON
- Return 404 with error envelope if not found

**Prerequisites:** Epic 1 (Foundation, models), Epic 2 (Tools repository base)

**Files Modified:**
- `backend/internal/tools/repository.go` (add GetToolBySlug)
- `backend/internal/tools/service.go` (add GetTool)
- `backend/internal/tools/handler.go` (add GET /:slug handler)

---

### Story 3.2: Tool Alternatives Backend API

**As a** developer,
**I want** an endpoint for tool alternatives and similar tools,
**So that** users can discover related options.

**Acceptance Criteria:**

**Given** I need alternatives data for a tool
**When** I implement the alternatives endpoint
**Then** the following endpoint works:

**Endpoint: `GET /api/v1/tools/:slug/alternatives`**
- Returns list of alternative and similar tools
- Response format: `{ data: { similar: [...tools], alternatives: [...tools] } }`
- Each tool in response includes: id, slug, name, logo_url, tagline, best_for, avg_rating_overall, review_count, pricing_summary (lightweight version, same as ToolCard)
- Empty arrays if no alternatives found (not an error)

**Logic:**
- Query `tool_alternatives` table where `tool_id = current_tool.id`
- Separate results by `relationship_type`:
  - `similar`: Tools tagged as similar (e.g., same category, overlapping features)
  - `alternatives`: Direct alternatives (e.g., "Use X instead of Y")
- Limit 6 per type for initial display
- Only return non-archived tools

**And:**
- Repository function: `GetToolAlternatives(toolID uint) (similar []Tool, alternatives []Tool, error)`
- Join query to load full tool objects for alternative_tool_id references
- Service layer calls repository and formats response

**Technical Implementation:**

- GORM query example:
  ```go
  var altRecords []ToolAlternative
  db.Where("tool_id = ?", toolID).
     Preload("AlternativeTool.PrimaryCategory").
     Preload("AlternativeTool.Tags").
     Find(&altRecords)

  for _, record := range altRecords {
      if record.RelationshipType == "similar" {
          similar = append(similar, record.AlternativeTool)
      } else if record.RelationshipType == "alternative" {
          alternatives = append(alternatives, record.AlternativeTool)
      }
  }
  ```
- Filter out archived alternatives: check `AlternativeTool.ArchivedAt IS NULL`

**Prerequisites:** Story 3.1 (Tool profile API)

**Files Modified:**
- `backend/internal/tools/repository.go` (add GetToolAlternatives)
- `backend/internal/tools/service.go`
- `backend/internal/tools/handler.go` (add GET /:slug/alternatives handler)

---

### Story 3.3: Tool Profile View - Hero & Overview

**As a** user,
**I want** to see a rich tool profile with hero section and overview,
**So that** I can quickly understand what the tool is and if it fits my needs.

**Acceptance Criteria:**

**Given** I navigate to `/tools/:slug`
**When** the ToolProfileView loads
**Then** I see the following sections:

**Hero Section:**
- Tool logo (large, right side in RTL) with fallback if missing
- Tool name as main heading (H1)
- Tagline as subheading
- Overall rating (stars + numeric, e.g., "4.5") with review count "(42 مراجعة)"
- Pricing badge: "مجاني" / "Freemium" / starting price
- Primary CTA button: "زيارة الأداة" (Visit Tool) - opens official_url in new tab
- Bookmark button (larger than card version): "حفظ" / "محفوظ"
- Category pill/link: navigates to category page
- Badges displayed as pills: "الأفضل في الفئة", "نجم صاعد", etc.

**Overview Section:**
- **"الأفضل لـ" (Best for)** - displayed prominently, 2-3 line summary
- **Description** - full tool description, multiple paragraphs
- **Primary Use Cases** - bullet list or pills showing main use cases
- **Target Roles** - who this tool is for (e.g., "المحامون، الكتّاب، المسوقون")
- **Platforms** - icons or text showing Web, Mobile, API, Desktop availability

**Loading State:**
- Skeleton loaders for hero and overview while API call is in progress
- No blocking spinner, layout structure visible with skeleton content

**Error Handling:**
- 404 page if tool not found: "الأداة غير موجودة" with link back to home
- Network error: retry button

**UX Requirements:**
- Dark background, neon blue for CTA and active states
- RTL-aware layout with proper Arabic typography (Cairo font)
- Responsive: single column on mobile, optimized spacing on desktop
- "Visit Tool" button is visually prominent (neon blue, larger size)

**Technical Implementation:**

- Fetch tool data on component mount: `apiClient.get(\`/tools/${route.params.slug}\`)`
- Use Vue reactive refs for tool data, loading, error states
- Bookmark interaction: `useBookmarksStore` to add/remove bookmark
- External link opens in new tab: `<a :href="tool.official_url" target="_blank" rel="noopener">`

**Prerequisites:** Story 3.1 (Tool profile API), Epic 1 (Frontend components)

**Files Modified:**
- `frontend/src/views/ToolProfileView.vue` (expand from Story 2.7 shell)

**New Files:**
- `frontend/src/components/tools/ToolHero.vue` (optional sub-component)
- `frontend/src/components/tools/ToolOverview.vue` (optional sub-component)

---

### Story 3.4: Tool Profile View - Features & Pricing

**As a** user,
**I want** to see detailed features and pricing information,
**So that** I can understand what the tool offers and how much it costs.

**Acceptance Criteria:**

**Given** I am viewing a tool profile
**When** I scroll below the overview
**Then** I see:

**Features Section:**
- **Heading:** "المميزات" (Features)
- List or grid of key features
- Each feature can have:
  - Icon (optional)
  - Feature name/title
  - Short description (optional)
- Displayed as cards or bulleted list depending on data structure
- If no features data available, section is hidden or shows placeholder

**Pricing Section:**
- **Heading:** "الأسعار" (Pricing)
- Free tier information (if `has_free_tier` is true):
  - "✓ يتوفر خطة مجانية" (Free tier available)
- Pricing summary displayed clearly (from `pricing_summary` field):
  - E.g., "يبدأ من $10/شهرياً" (Starts at $10/month)
  - Or structured tiers if data supports it (Basic, Pro, Enterprise)
- Link to official pricing page: "عرض تفاصيل الأسعار" (View pricing details) linking to official_url
- If no pricing data, show: "الرجاء زيارة الموقع الرسمي للأسعار" (Please visit official site for pricing)

**Visual Design:**
- Sections use consistent heading style (H2, Cairo font, white text)
- Features displayed in 2-3 column grid on desktop, stacked on mobile
- Pricing displayed as structured cards or simple text block
- Proper spacing between sections (24-32px)

**Technical Implementation:**

- Parse `tool.pricing_summary` and display as formatted text
- Features: if stored as JSON array, iterate and display; if text, display as-is
- Conditional rendering: `v-if="tool.features && tool.features.length > 0"`
- Pricing tiers: could be hardcoded structure or parsed from pricing_summary

**Prerequisites:** Story 3.3 (Tool profile view structure)

**Files Modified:**
- `frontend/src/views/ToolProfileView.vue` (add sections)

**New Files (optional):**
- `frontend/src/components/tools/ToolFeatures.vue`
- `frontend/src/components/tools/ToolPricing.vue`

---

### Story 3.5: Tool Profile View - Media Gallery

**As a** user,
**I want** to see screenshots and videos of the tool,
**So that** I can visually understand how it works.

**Acceptance Criteria:**

**Given** I am viewing a tool profile
**When** the media section loads
**Then** I see:

**Media Section:**
- **Heading:** "صور وفيديوهات" (Screenshots and Videos)
- Media items displayed in order (from `display_order` field)
- **Screenshots:**
  - Displayed as thumbnail grid (2-3 columns on desktop, 1-2 on mobile)
  - Click thumbnail opens lightbox/modal with full-size image
  - Lightbox has prev/next navigation for multiple images
  - Image alt text for accessibility
- **Videos (YouTube embeds):**
  - Displayed as embedded YouTube player or thumbnail
  - Lazy-loaded: only load iframe when user scrolls to media section or clicks play
  - Responsive embed: maintains 16:9 aspect ratio
  - If video URL is YouTube, extract video ID and embed properly
- **Mixed media:**
  - Screenshots displayed first, then videos
  - Or single gallery with media type indicators

**Empty State:**
- If no media: section is hidden or shows "لا توجد صور أو فيديوهات متاحة" (No media available)

**Performance:**
- Lazy load YouTube iframes to avoid blocking page render (critical from UX Design spec)
- Use `loading="lazy"` for screenshot images
- Skeleton placeholders while media loads

**Technical Implementation:**

- **ToolMediaGallery.vue** component
- Iterate over `tool.media` array:
  ```vue
  <div v-for="media in tool.media" :key="media.id">
    <img v-if="media.type === 'screenshot'"
         :src="media.thumbnail_url || media.url"
         @click="openLightbox(media)" />
    <YouTubeEmbed v-else-if="media.type === 'video'"
                  :url="media.url"
                  lazy />
  </div>
  ```
- Lightbox modal: use a lightweight library or custom Vue component
- YouTube lazy loading: render thumbnail with play button, only load iframe on click

**Prerequisites:** Story 3.3 (Tool profile structure)

**Files Modified:**
- `frontend/src/views/ToolProfileView.vue` (add media section)

**New Files:**
- `frontend/src/components/tools/ToolMediaGallery.vue`
- `frontend/src/components/common/Lightbox.vue` (for image viewing)
- `frontend/src/components/common/YouTubeEmbed.vue` (lazy-loaded YouTube player)

---

### Story 3.6: Tool Profile View - Alternatives Section

**As a** user,
**I want** to see similar and alternative tools on the profile,
**So that** I can explore other options before making a decision.

**Acceptance Criteria:**

**Given** I am viewing a tool profile
**When** the alternatives section loads
**Then** I see:

**Similar Tools Subsection:**
- **Heading:** "أدوات مشابهة" (Similar Tools)
- Horizontal scrolling row of tool cards (using ToolCard component)
- Shows 4-6 similar tools
- Cards are clickable, navigate to respective tool profiles
- If no similar tools: section is hidden

**Alternatives Subsection:**
- **Heading:** "بدائل لهذه الأداة" (Alternatives to this tool)
- Horizontal scrolling row of tool cards
- Shows 4-6 alternative tools
- Cards are clickable, navigate to respective tool profiles
- If no alternatives: section is hidden

**Visual Design:**
- Each subsection clearly separated with heading
- Horizontal scroll with subtle scroll indicators (gradient fade on edges)
- Tool cards use same ToolCard component as search results (consistent design)
- RTL-aware scrolling: scroll starts from right side
- Mobile: scroll snap for smooth scrolling experience

**Technical Implementation:**

- Fetch alternatives on component mount: `apiClient.get(\`/tools/${route.params.slug}/alternatives\`)`
- Store in reactive refs: `similar`, `alternatives`
- Conditional rendering: only show sections if arrays have items
- Horizontal scroll container:
  ```vue
  <div class="flex overflow-x-auto gap-4 pb-4">
    <ToolCard v-for="tool in similar" :key="tool.id" :tool="tool" />
  </div>
  ```
- RTL scrolling: use `dir="rtl"` on scroll container or CSS `direction: rtl`

**Prerequisites:** Story 3.2 (Alternatives API), Story 2.4 (ToolCard component)

**Files Modified:**
- `frontend/src/views/ToolProfileView.vue` (add alternatives section)

---

### Story 3.7: Social Proof & Engagement Indicators

**As a** user,
**I want** to see social proof indicators like badges and bookmark counts,
**So that** I can gauge the tool's popularity and credibility.

**Acceptance Criteria:**

**Given** I am viewing a tool profile
**When** the page loads
**Then** I see social proof indicators:

**Badges (already in Story 3.3 hero):**
- Displayed prominently in hero section as pills/chips
- Examples: "الأفضل في الفئة" (Top in Category), "نجم صاعد" (Rising Star), "اختيار المحرر" (Editor's Pick)
- Each badge has icon (if available) + name
- Color-coded or styled distinctively from regular tags

**Bookmark Count:**
- Displayed near bookmark button or in stats area
- Format: "حفظه 234 مستخدماً" (Bookmarked by 234 users)
- Updates when user bookmarks/unbookmarks (optimistic UI)

**Review Count & Rating:**
- Already covered in hero section (Story 3.3)
- Overall rating with stars + numeric value
- Review count as clickable link scrolling to reviews section (will be in Epic 4)

**Engagement Stats (Optional):**
- Small stats panel showing:
  - 🔖 Bookmark count
  - ⭐ Rating (avg_rating_overall)
  - 💬 Review count
- Displayed below hero or in sidebar (desktop)

**Visual Design:**
- Badges use neon blue or distinct colors, slightly elevated (subtle shadow/glow)
- Stats displayed with icons and clear Arabic labels
- RTL-aware layout for badge row
- Subtle animations on hover for badges

**Technical Implementation:**

- Badges already loaded in tool profile API (Story 3.1)
- Display: iterate over `tool.badges` array
- Bookmark count from `tool.bookmark_count`
- Update bookmark count optimistically when user bookmarks:
  ```typescript
  const displayBookmarkCount = computed(() => {
    return tool.value.bookmark_count + (isBookmarked ? 1 : 0);
  });
  ```

**Prerequisites:** Story 3.3 (Hero section), Story 3.1 (Tool data with badges)

**Files Modified:**
- `frontend/src/views/ToolProfileView.vue` (ensure badges and stats are displayed)
- `frontend/src/components/tools/ToolHero.vue` (if using sub-component)

---

## Epic 4: Reviews & Ratings

**Goal:** Users can see authentic peer feedback through structured reviews and contribute their own experiences - building trust and enabling informed decisions.

### Story 4.1: Reviews Backend API

**As a** developer,
**I want** endpoints for fetching and submitting reviews,
**So that** users can view and contribute structured reviews.

**Acceptance Criteria:**

**Endpoint: `GET /api/v1/tools/:slug/reviews`**
- Returns paginated list of reviews for a tool
- Response: `{ data: [...reviews], meta: { page, page_size, total } }`
- Each review includes: id, rating_overall, rating_ease_of_use, rating_value, rating_accuracy, rating_speed, rating_support, pros, cons, primary_use_case, reviewer_role, company_size, usage_context, helpful_count, created_at, user (id, display_name)
- Query params: `?page=1&page_size=10&sort=newest|most_helpful|highest|lowest`
- Default sort: newest
- Only approved/non-hidden reviews returned

**Endpoint: `POST /api/v1/tools/:slug/reviews`** (Authenticated)
- Accepts review submission from logged-in users
- Request body: rating_overall (required), rating dimensions (optional), pros (required), cons (required), primary_use_case, reviewer_role, company_size, usage_context
- Validation: rating 1-5, pros/cons max 500 chars each, required fields present
- Creates review with moderation_status: auto-approved or pending based on config
- Updates tool.avg_rating_overall and tool.review_count
- Returns created review object
- 401 if not authenticated, 422 if validation fails

**And:**
- Repository: `ListReviewsByTool`, `CreateReview`, `UpdateToolRatingAggregates`
- Service: validates input, checks user hasn't already reviewed this tool
- Handler: uses auth middleware, standardized responses

**Technical Implementation:**
- GORM preload user: `.Preload("User", func(db *gorm.DB) *gorm.DB { return db.Select("id", "display_name") })`
- Rating aggregation: recalculate avg_rating_overall from all approved reviews
- Rate limiting: prevent spam, 1 review per user per tool

**Prerequisites:** Epic 1, Epic 3 (Tool profile)

**Files Modified:**
- `backend/internal/reviews/repository.go`
- `backend/internal/reviews/service.go`
- `backend/internal/reviews/handler.go`

---

### Story 4.2: Review Display on Tool Profile

**As a** user,
**I want** to see structured reviews on the tool profile,
**So that** I can learn from others' experiences.

**Acceptance Criteria:**

**Reviews Section on ToolProfileView:**
- Heading: "المراجعات" (Reviews) with overall rating summary
- Sort controls: Newest, Most helpful, Highest rated, Lowest rated
- Each review card shows:
  - Reviewer display name + role + company size
  - Overall rating (stars) + submission date
  - Pros/Cons in structured format ("الإيجابيات" / "السلبيات")
  - Primary use case badge
  - Usage context tags
  - Helpful count with "مفيد؟" (Helpful?) button (future feature)
- Pagination: Load more button or infinite scroll
- Empty state: "لا توجد مراجعات بعد. كن أول من يراجع!" with "Add Review" button

**Visual Design:**
- Review cards with dark background, subtle borders
- Pros in green-tinted section, cons in red-tinted section
- Reviewer info in lighter text
- Stars displayed clearly for rating

**Technical Implementation:**
- Fetch reviews: `apiClient.get(\`/tools/${slug}/reviews\`, { params: { sort, page } })`
- ReviewList.vue component iterates over reviews
- Sort state managed locally or in URL params

**Prerequisites:** Story 4.1, Story 3.3

**Files Modified:**
- `frontend/src/views/ToolProfileView.vue`

**New Files:**
- `frontend/src/components/reviews/ReviewList.vue`
- `frontend/src/components/reviews/ReviewCard.vue`

---

### Story 4.3: Review Submission Form

**As a** logged-in user,
**I want** to submit a structured review for a tool,
**So that** I can share my experience with others.

**Acceptance Criteria:**

**Review Form (modal or inline):**
- Triggered by "إضافة مراجعة" (Add Review) button on tool profile
- Only visible if user is authenticated, else shows "Sign in to review"
- Form fields:
  - Overall rating: 5-star selector (required)
  - Dimension ratings: Ease of use, Value, Accuracy, Speed, Support (optional, 5-star each)
  - Pros: textarea, max 500 chars, required
  - Cons: textarea, max 500 chars, required
  - Primary use case: dropdown (options from backend or predefined list)
  - Reviewer role: dropdown (Developer, Designer, Manager, etc.)
  - Company size: dropdown (1-10, 11-50, 51-200, 201+)
  - Usage context: checkboxes (Personal, Work, Freelance, etc.)
- Submit button: "نشر المراجعة" (Publish Review)
- Cancel button closes form

**Validation:**
- Client-side: required fields, character limits, rating range
- Arabic error messages: "الرجاء تعبئة هذا الحقل", "يجب أن يكون التقييم من 1 إلى 5"
- Server-side validation errors displayed clearly

**Submission Flow:**
- On submit: POST to /api/v1/tools/:slug/reviews
- Loading state: button shows "جاري النشر..." (Publishing...)
- Success: close modal, show success toast "شكراً! تمت إضافة مراجعتك", refresh reviews list
- Error: show error message, keep form open

**Technical Implementation:**
- ReviewForm.vue component with reactive form state
- useSessionStore to check authentication
- Form validation with Vuelidate or manual validation
- API call with error handling

**Prerequisites:** Story 4.1 (Reviews API), Epic 6 (Auth required)

**New Files:**
- `frontend/src/components/reviews/ReviewForm.vue`

---

### Story 4.4: Aggregated Ratings Display

**As a** user,
**I want** to see aggregated ratings at the top of the reviews section,
**So that** I can quickly understand overall sentiment.

**Acceptance Criteria:**

**Ratings Summary Component:**
- Overall rating (large): 4.5 / 5 stars with total review count
- Rating distribution: bar chart showing count per star (5★: 42, 4★: 18, etc.)
- Dimension ratings: horizontal bars for Ease of use, Value, Accuracy, Speed, Support
- Displayed above the reviews list
- Updates when new reviews are submitted

**Visual Design:**
- Clean, scannable layout
- Neon blue for filled stars and bars
- Rating distribution bars clearly show proportion
- Arabic labels for all dimensions

**Technical Implementation:**
- Data from tool object: avg_rating_overall, review_count
- Dimension ratings: calculate averages from reviews or store on tool table
- Rating distribution: fetch from backend or calculate client-side

**Prerequisites:** Story 4.1, Story 4.2

**New Files:**
- `frontend/src/components/reviews/RatingsAggregation.vue`

---

## Epic 5: Comparison & Shortlists

**Goal:** Users can bookmark promising tools, manage their shortlist, and compare options side-by-side to make final consolidation and selection decisions.

### Story 5.1: Bookmarks Backend API

**As a** developer,
**I want** endpoints for managing user bookmarks,
**So that** users can save tools to their shortlist.

**Acceptance Criteria:**

**Endpoint: `GET /api/v1/me/bookmarks`** (Authenticated or session-based)
- Returns list of bookmarked tools for current user
- Response: `{ data: [...tools] }` (full tool objects, same as search results)
- Ordered by bookmark created_at DESC (most recent first)
- Supports anonymous users via session_id stored in bookmarks table

**Endpoint: `POST /api/v1/me/bookmarks`** (Authenticated or session-based)
- Request: `{ tool_id: number }`
- Creates bookmark for user/session
- Increments tool.bookmark_count
- Returns 201 with created bookmark
- Returns 409 if already bookmarked (idempotent)

**Endpoint: `DELETE /api/v1/me/bookmarks/:tool_id`** (Authenticated or session-based)
- Deletes bookmark for user/session
- Decrements tool.bookmark_count
- Returns 204 No Content
- Returns 404 if bookmark doesn't exist

**And:**
- Support anonymous sessions: generate session_id, store in cookie
- On user login: migrate session bookmarks to user account
- Repository: `GetUserBookmarks`, `AddBookmark`, `RemoveBookmark`
- Unique constraint on (user_id, tool_id) and (session_id, tool_id)

**Technical Implementation:**
- Session middleware: generate UUID session_id if not present
- Migration on login: UPDATE bookmarks SET user_id = ?, session_id = NULL WHERE session_id = ?
- Update tool counts: db.Model(&Tool{}).Where("id = ?", toolID).UpdateColumn("bookmark_count", gorm.Expr("bookmark_count + ?", delta))

**Prerequisites:** Epic 1 (Auth), Epic 2 (Tools)

**Files Created:**
- `backend/internal/bookmarks/repository.go`
- `backend/internal/bookmarks/service.go`
- `backend/internal/bookmarks/handler.go`

---

### Story 5.2: Bookmarks View

**As a** user,
**I want** to view and manage my bookmarked tools,
**So that** I can access my shortlist and make decisions.

**Acceptance Criteria:**

**BookmarksView (`/bookmarks`):**
- Page title: "المفضلة" (My Bookmarks)
- Grid of bookmarked tools using ToolCard component
- Each card shows bookmark button (filled/active state)
- "Remove from bookmarks" action on each card
- "Compare selected" button: allows multi-select and send to comparison
- Empty state: "لم تحفظ أي أدوات بعد. ابدأ بالتصفح!" with link to home
- Anonymous users: bookmarks stored in localStorage, prompt to sign in to sync across devices

**Multi-Select for Comparison:**
- Checkbox on each card for selection
- "مقارنة المحددة" (Compare Selected) button appears when 2+ tools selected
- Click button navigates to `/compare?tools=slug1,slug2,slug3`
- Max 4 tools can be selected

**Technical Implementation:**
- Fetch bookmarks: `apiClient.get('/me/bookmarks')`
- useBookmarksStore to manage state
- Multi-select: local state tracking selected tool IDs
- Remove bookmark: calls store.removeBookmark(toolId), updates UI optimistically

**Prerequisites:** Story 5.1 (Bookmarks API), Story 2.4 (ToolCard)

**Files Created:**
- `frontend/src/views/BookmarksView.vue`

---

### Story 5.3: Add to Compare from Results

**As a** user,
**I want** to add tools to comparison from search results and profiles,
**So that** I can build up a comparison set.

**Acceptance Criteria:**

**Compare Selection:**
- ToolCard component has "مقارنة" (Compare) checkbox/button
- Clicking adds tool to comparison set (max 4 tools)
- Visual indicator shows tool is in comparison set
- Floating comparison bar appears at bottom when 1+ tools selected:
  - Shows selected tool logos/names
  - "مقارنة" (Compare) button to go to comparison view
  - "مسح" (Clear) button to clear selection
  - Tool count: "2 من 4 أدوات"
- Clicking "Compare" button navigates to `/compare?tools=slug1,slug2`

**Persistence:**
- Comparison selection stored in Pinia store
- Persists across page navigation within session
- Option to save comparison set to bookmarks (future)

**Technical Implementation:**
- useComparisonStore: `{ selectedTools: [], addTool(), removeTool(), clear() }`
- ComparisonBar.vue component fixed at bottom
- ToolCard emits @add-to-compare event, parent calls store.addTool()
- Max 4 validation: show toast if limit reached

**Prerequisites:** Story 2.4 (ToolCard)

**New Files:**
- `frontend/src/stores/comparison.ts`
- `frontend/src/components/compare/ComparisonBar.vue`

**Files Modified:**
- `frontend/src/components/tools/ToolCard.vue` (add compare interaction)

---

### Story 5.4: Comparison View

**As a** user,
**I want** to see tools side-by-side in a comparison table,
**So that** I can evaluate options and make a decision.

**Acceptance Criteria:**

**CompareView (`/compare?tools=slug1,slug2,slug3`):**
- Page title: "المقارنة" (Comparison)
- Side-by-side table with tools as columns (2-4 tools)
- Comparison rows:
  - **Overview:** Name, logo, "best for", tagline
  - **Category:** Primary category
  - **Pricing:** Free tier, starting price, billing model
  - **Rating:** Overall rating, review count
  - **Features:** Key features (if available)
  - **Social Proof:** Bookmark count, badges
  - **Actions:** "Visit Tool" button, bookmark button for each tool
- "Remove from comparison" button on each column
- "Share comparison" button: copies URL to clipboard
- Empty state: Prompt to add tools from search or bookmarks

**Responsive:**
- Desktop: full table with all columns visible
- Mobile: horizontal scroll, sticky first column (row labels)

**Technical Implementation:**
- Parse tool slugs from URL: `route.query.tools.split(',')`
- Fetch all tools: `Promise.all(slugs.map(slug => apiClient.get(\`/tools/${slug}\`)))`
- CompareTable.vue component builds table from tool data
- Remove tool: update URL, refetch remaining tools

**Prerequisites:** Story 5.3 (Comparison selection), Epic 3 (Tool data)

**Files Created:**
- `frontend/src/views/CompareView.vue`
- `frontend/src/components/compare/CompareTable.vue`

---

## Epic 6: User Accounts & Persistence

**Goal:** Users can create accounts or sign in to persist their bookmarks and reviews across sessions and devices, eliminating the risk of losing their work.

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

### Story 6.2: Sign In / Sign Up UI

**As a** user,
**I want** to sign in or create an account,
**So that** I can save my bookmarks and reviews across devices.

**Acceptance Criteria:**

**Sign In Modal:**
- Triggered from "تسجيل الدخول" button in header
- Form fields: Email, Password
- Submit button: "تسجيل الدخول"
- Link to sign up: "ليس لديك حساب؟ أنشئ حساباً"
- Validation: email format, password required
- Error handling: "البريد الإلكتروني أو كلمة المرور غير صحيحة"
- Success: close modal, update header UI, show success toast

**Sign Up Modal:**
- Triggered from "إنشاء حساب" button or link in sign in modal
- Form fields: Display Name, Email, Password
- Submit button: "إنشاء حساب"
- Link to sign in: "لديك حساب؟ سجّل الدخول"
- Validation: all fields required, email format, password min 8 chars
- Error handling: "البريد الإلكتروني مستخدم بالفعل"
- Success: close modal, update header UI, show success toast
- Anonymous bookmark migration: show toast "تم حفظ X أدوات إلى حسابك"

**Header UI Updates:**
- When authenticated: show avatar/initials + display name
- Dropdown menu: "مراجعاتي" (My Reviews), "المفضلة" (My Bookmarks), "تسجيل الخروج" (Sign Out)
- When not authenticated: "تسجيل الدخول", "إنشاء حساب" buttons

**Technical Implementation:**
- AuthModal.vue component with tabs for sign in/sign up
- useSessionStore actions: login, register, logout, fetchCurrentUser
- On mount: call fetchCurrentUser to restore session
- Form validation with reactive error state

**Prerequisites:** Story 6.1 (Auth endpoints)

**New Files:**
- `frontend/src/components/auth/AuthModal.vue`
- `frontend/src/components/auth/SignInForm.vue`
- `frontend/src/components/auth/SignUpForm.vue`

**Files Modified:**
- `frontend/src/components/layout/HeaderNav.vue` (add auth UI)

---

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

## Epic 7: Admin & Catalog Management

**Goal:** Admins can maintain a high-quality, accurate, and complete tool directory through CRUD operations, taxonomy management, and data quality checks - ensuring users always find trustworthy information.

### Story 7.1: Admin Tools Management Backend

**As a** developer,
**I want** admin endpoints for managing tools,
**So that** admins can maintain the catalog.

**Acceptance Criteria:**

**Endpoints (all require admin role):**

**`GET /api/v1/admin/tools`**
- Returns all tools including archived
- Query params: `?search=...&archived=true|false&page=1`
- Response: `{ data: [...tools], meta: { page, page_size, total } }`

**`POST /api/v1/admin/tools`**
- Create new tool
- Request: all tool fields (name, slug, description, best_for, category_id, etc.)
- Validation: required fields, unique slug
- Returns created tool

**`PATCH /api/v1/admin/tools/:id`**
- Update tool fields
- Request: partial tool object
- Validates changes
- Returns updated tool

**`DELETE /api/v1/admin/tools/:id`**
- Soft delete: sets archived_at timestamp
- Tool no longer appears in public endpoints
- Returns 204 No Content

**And:**
- AdminRequired middleware on all routes
- Service layer validates admin permissions
- Audit logging: track who made changes (future enhancement)

**Prerequisites:** Epic 1 (Auth with admin role)

**Files Modified:**
- `backend/internal/tools/handler.go` (add admin routes)
- `backend/internal/tools/service.go` (add admin methods)
- `backend/internal/platform/http/router.go` (register admin routes)

---

### Story 7.2: Admin UI - Tools Management

**As an** admin,
**I want** a UI to manage tools,
**So that** I can keep the catalog accurate and complete.

**Acceptance Criteria:**

**Admin Tools Page (`/admin/tools`):**
- Table of all tools with columns: Logo, Name, Category, Rating, Reviews, Bookmarks, Status (Active/Archived), Actions
- Search bar: filter by name/slug
- Filter: All / Active / Archived
- "Add Tool" button opens tool form modal
- Actions per tool: Edit, Archive/Unarchive, View on site
- Pagination

**Tool Form Modal:**
- All fields editable: name, slug, tagline, description, best_for, primary_use_cases, pricing_summary, target_roles, platforms, has_free_tier, official_url, category, tags, media URLs
- Validation: required fields highlighted
- Preview: shows how tool card will look
- Submit: saves tool, closes modal, refreshes list
- Cancel: closes without saving

**Visual Design:**
- Admin UI uses same dark theme but with admin-specific nav
- Table is sortable and filterable
- Clear visual difference between active/archived tools

**Technical Implementation:**
- AdminToolsView.vue with data table
- ToolFormModal.vue component for create/edit
- Admin routes protected: check user role in router beforeEach
- API calls to admin endpoints

**Prerequisites:** Story 7.1 (Admin tools API)

**New Files:**
- `frontend/src/views/admin/AdminToolsView.vue`
- `frontend/src/components/admin/ToolFormModal.vue`
- `frontend/src/router/admin.ts` (admin routes)

---

### Story 7.3: Admin Categories & Tags Management

**As an** admin,
**I want** to manage categories and tags,
**So that** taxonomy remains coherent.

**Acceptance Criteria:**

**Endpoints:**
- `GET/POST/PATCH/DELETE /api/v1/admin/categories`
- `GET/POST/PATCH/DELETE /api/v1/admin/tags`
- All require admin role

**Admin Categories Page:**
- List of categories with: Name, Slug, Tool Count, Display Order, Actions
- Create/Edit category: name, slug, description, icon_url, display_order
- Cannot delete category if tools exist (show warning, require reassignment)
- Drag-and-drop reordering by display_order

**Admin Tags Page:**
- List of tags with: Name, Slug, Tool Count, Actions
- Create/Edit tag: name, slug
- Delete tag: removes from all tools (confirmation required)

**Technical Implementation:**
- Similar pattern to tools management
- Category/tag CRUD operations
- Prevent deletion with referential integrity checks

**Prerequisites:** Story 7.1

**New Files:**
- `frontend/src/views/admin/AdminCategoriesView.vue`
- `frontend/src/views/admin/AdminTagsView.vue`

**Files Modified:**
- `backend/internal/categories/handler.go`
- Backend tag handlers (new files if not existing)

---

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

## Epic 8: Moderation & Content Quality

**Goal:** The platform maintains trusted, abuse-free user-generated content through reporting mechanisms and moderation workflows - ensuring the review layer remains valuable and credible.

### Story 8.1: Reporting Backend

**As a** developer,
**I want** endpoints for reporting content,
**So that** users can flag problematic tools or reviews.

**Acceptance Criteria:**

**Endpoints:**

**`POST /api/v1/tools/:slug/report`**
- Request: `{ reason: "spam|abuse|misinformation|other", comment?: string }`
- Creates report record
- Returns 201 Created

**`POST /api/v1/reviews/:id/report`**
- Same structure as tool reports
- Creates report record
- Returns 201 Created

**Reports Table:**
- Fields: id, reportable_type (tool/review), reportable_id, reporter_user_id (nullable for anonymous), reason, comment, status (pending/reviewed/dismissed), created_at

**And:**
- Repository: CreateReport
- Service: validates reason enum, limits reports per user
- Both authenticated and anonymous users can report

**Prerequisites:** Epic 1

**New Files:**
- `backend/internal/moderation/repository.go`
- `backend/internal/moderation/service.go`
- `backend/internal/moderation/handler.go`

---

### Story 8.2: Report UI

**As a** user,
**I want** to report tools or reviews,
**So that** I can flag inappropriate content.

**Acceptance Criteria:**

**Report Button:**
- On tool profiles: "إبلاغ" (Report) link near tool name
- On review cards: "إبلاغ" (Report) link in review card footer

**Report Modal:**
- Dropdown: select reason (Spam, Abuse, Misinformation, Other)
- Textarea: optional comment (max 500 chars)
- Submit button: "إرسال البلاغ"
- Cancel button

**Submission:**
- On submit: POST to report endpoint
- Success: close modal, show toast "شكراً للإبلاغ. سنراجعه قريباً"
- Error: show error message in modal

**Technical Implementation:**
- ReportModal.vue component
- Accept props: reportableType, reportableId
- API call to appropriate report endpoint

**Prerequisites:** Story 8.1

**New Files:**
- `frontend/src/components/moderation/ReportModal.vue`

**Files Modified:**
- `frontend/src/views/ToolProfileView.vue` (add report button)
- `frontend/src/components/reviews/ReviewCard.vue` (add report button)

---

### Story 8.3: Moderation Queue Backend

**As a** developer,
**I want** moderation queue endpoints,
**So that** moderators can review reported content.

**Acceptance Criteria:**

**Endpoints (require moderator or admin role):**

**`GET /api/v1/moderation/queue`**
- Returns paginated list of reports
- Query params: `?type=tool|review&status=pending|reviewed|dismissed&page=1`
- Response: `{ data: [...reports], meta: { page, page_size, total } }`
- Each report includes: reportable object (tool or review with full details), reporter info, reason, comment, created_at

**`PATCH /api/v1/moderation/reviews/:id/approve`**
- Sets review.moderation_status = "approved"
- Logs moderation action
- Returns updated review

**`PATCH /api/v1/moderation/reviews/:id/hide`**
- Sets review.moderation_status = "hidden"
- Review no longer appears in public endpoints
- Logs moderation action

**`PATCH /api/v1/moderation/reviews/:id/remove`**
- Sets review.moderation_status = "removed"
- Permanently hides review
- Updates tool rating aggregates
- Logs moderation action

**`GET /api/v1/moderation/history/:review_id`**
- Returns audit log of moderation actions for a review
- Response: `{ data: [...actions] }` with actor, action_type, timestamp, notes

**And:**
- Moderation actions table: tracks who did what when
- Service layer validates moderator permissions
- Aggregates: recalculate tool ratings when reviews hidden/removed

**Prerequisites:** Story 8.1

**Files Modified:**
- `backend/internal/moderation/repository.go`
- `backend/internal/moderation/service.go`
- `backend/internal/moderation/handler.go`
- `backend/internal/reviews/service.go` (moderation status updates)

---

### Story 8.4: Moderation Queue UI

**As a** moderator,
**I want** a moderation queue UI,
**So that** I can review and act on reported content.

**Acceptance Criteria:**

**Moderation Queue Page (`/moderation/queue`):**
- Table of reports with columns: Type (Tool/Review), Content Preview, Reason, Reporter, Date, Status, Actions
- Filters: Type (All/Tools/Reviews), Status (Pending/Reviewed/Dismissed)
- Pagination

**Report Details:**
- Click report row expands to show:
  - Full content (tool details or review text)
  - Reporter info (if available)
  - Report reason and comment
  - "View in context" link (opens tool profile or review in new tab)

**Actions:**
- For reviews: Approve, Hide, Remove buttons
- For tools: View on site, Flag for admin review (tools handled by admins in Story 7.2)
- Dismiss report button (marks report as reviewed without action)
- Confirmation modals for destructive actions

**Moderation History:**
- Link to view history for a review
- Shows timeline of actions: who approved/hid/removed and when

**Visual Design:**
- Clear status indicators: Pending (yellow), Reviewed (green), Dismissed (gray)
- Content preview truncated with "Show more"
- Actions clearly labeled

**Technical Implementation:**
- ModerationQueueView.vue fetches reports
- Action buttons call moderation endpoints
- Optimistic UI: update status immediately, rollback on error
- Role check: only moderators/admins can access

**Prerequisites:** Story 8.3

**New Files:**
- `frontend/src/views/moderation/ModerationQueueView.vue`
- `frontend/src/router/moderation.ts` (moderation routes)

---

## Final Validation & Summary

### Epic Completion Summary

**✅ All 8 Epics Complete:**

1. **Epic 1: Foundation & Core Infrastructure** - 7 stories
2. **Epic 2: Tool Discovery & Browsing** - 7 stories
3. **Epic 3: Rich Tool Profiles** - 7 stories
4. **Epic 4: Reviews & Ratings** - 4 stories
5. **Epic 5: Comparison & Shortlists** - 4 stories
6. **Epic 6: User Accounts & Persistence** - 3 stories
7. **Epic 7: Admin & Catalog Management** - 4 stories
8. **Epic 8: Moderation & Content Quality** - 4 stories

**Total Stories: 40 implementation-ready user stories**

### Complete FR Coverage Matrix

| FR | Requirement | Epic | Stories |
|----|-------------|------|---------|
| FR1 | Free-text search | Epic 2 | 2.2, 2.3, 2.5 |
| FR2 | Browse by category | Epic 2 | 2.1, 2.3, 2.6 |
| FR3 | Refine with filters | Epic 2 | 2.2, 2.5 |
| FR4 | Sort results | Epic 2 | 2.2, 2.5 |
| FR5 | Open tool detail | Epic 2 | 2.7 |
| FR6 | View tool profile | Epic 3 | 3.1, 3.3, 3.4 |
| FR7 | View rich media | Epic 3 | 3.1, 3.5 |
| FR8 | View reviews | Epic 4 | 4.1, 4.2, 4.4 |
| FR9 | Submit review | Epic 4 | 4.1, 4.3 |
| FR10 | View aggregated ratings | Epic 4 | 4.1, 4.4 |
| FR11 | View social proof | Epic 3 | 3.1, 3.3, 3.7 |
| FR12 | Add to comparison | Epic 5 | 5.3 |
| FR13 | View comparison | Epic 5 | 5.4 |
| FR14 | View alternatives | Epic 3 | 3.2, 3.6 |
| FR15 | Bookmark tool | Epic 5 | 5.1, 5.2 |
| FR16 | View bookmarks | Epic 5 | 5.1, 5.2 |
| FR17 | Remove bookmark | Epic 5 | 5.1, 5.2 |
| FR18 | Send to comparison | Epic 5 | 5.2, 5.4 |
| FR19 | Create account/sign in | Epic 6 | 6.1, 6.2 |
| FR20 | View own activity | Epic 6 | 6.3 |
| FR21 | Admin manage tools | Epic 7 | 7.1, 7.2 |
| FR22 | Admin manage taxonomy | Epic 7 | 7.3 |
| FR23 | Admin manage badges | Epic 7 | 7.4 |
| FR24 | Admin resolve issues | Epic 7 | 7.2 |
| FR25 | Report content | Epic 8 | 8.1, 8.2 |
| FR26 | Moderator view queue | Epic 8 | 8.3, 8.4 |
| FR27 | Moderator actions | Epic 8 | 8.3, 8.4 |
| FR28 | View mod history | Epic 8 | 8.3, 8.4 |
| FR29 | Admin view analytics | Epic 7 | 7.4 |
| FR30 | Admin view top items | Epic 7 | 7.4 |

**Coverage: 30/30 FRs (100%)**

### Technical Context Integration

**Backend Architecture Applied:**
- ✅ Go + Gin + GORM + PostgreSQL stack throughout
- ✅ All 11 database tables with proper relationships
- ✅ RESTful `/api/v1` endpoints with standardized responses
- ✅ JWT authentication with HTTP-only cookies
- ✅ Repository-Service-Handler pattern consistently
- ✅ snake_case naming for DB/API, PascalCase for Go types

**Frontend Architecture Applied:**
- ✅ Vue 3 + Vite + TypeScript + Pinia + Tailwind CSS
- ✅ All key views and components defined
- ✅ Pinia stores for state management (session, bookmarks, filters, comparison)
- ✅ API client with typed responses
- ✅ Vue Router with SEO-friendly routes

**UX Design Integration:**
- ✅ Hero search-driven experience with Arabic text
- ✅ Dark theme (#05060A) with neon blue primary (#3B82F6)
- ✅ Cairo font, RTL-aware layouts throughout
- ✅ G2-style structured cards and patterns
- ✅ Skeleton loaders, lazy loading, optimistic UI
- ✅ Performance-first: no blocking spinners, lazy YouTube embeds

### Implementation Readiness Assessment

**✅ READY FOR PHASE 4 IMPLEMENTATION**

**Confidence Level:** Very High

**Rationale:**
- All 30 functional requirements mapped to specific, actionable stories
- Every story includes complete acceptance criteria with technical implementation details
- Backend and frontend patterns consistently applied across all epics
- Architecture decisions (Go/Gin, Vue 3, PostgreSQL, JWT) fully incorporated
- UX requirements (Arabic-first, RTL, dark theme, performance) integrated throughout
- Clear dependencies between stories enable sequential or parallel implementation
- Stories sized for single dev agent completion (each story is independently deployable)

**Ready for:**
- Sprint planning workflow to create sprint tracking file
- Dev agent execution via `dev-story` workflow
- Parallel implementation by multiple dev agents (respecting dependencies)

---

🎉 **Epic and Story Creation Complete!**

**Output:** `docs/epics.md` with 8 epics, 40 stories, 100% FR coverage, complete technical and UX context integration.
