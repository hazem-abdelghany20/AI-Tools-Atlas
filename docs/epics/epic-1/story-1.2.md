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
