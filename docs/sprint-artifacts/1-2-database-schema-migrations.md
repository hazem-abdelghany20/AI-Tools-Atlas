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

### Tasks/Subtasks

- [x] Setup golang-migrate dependency
  - [x] Install `golang-migrate/migrate/v4` and PostgreSQL driver
  - [x] Verify installation
- [x] Create migration directory structure
  - [x] Create `backend/migrations/` directory
- [x] Create initial schema UP migration
  - [x] Create `000001_initial_schema.up.sql`
  - [x] Add core tables: tools, categories, tags, tool_tags, media
  - [x] Add badge tables: badges, tool_badges
  - [x] Add relationship tables: tool_alternatives
  - [x] Add user/auth tables: users, reviews, bookmarks
  - [x] Add all required indexes
  - [x] Add foreign key constraints with ON DELETE behavior
- [x] Create initial schema DOWN migration
  - [x] Create `000001_initial_schema.down.sql`
  - [x] Add DROP TABLE statements in reverse dependency order
- [x] Create/update Makefile with migration targets
  - [x] Add `migrate-up` target
  - [x] Add `migrate-down` target
  - [x] Add `migrate-create` target with name parameter
- [x] Write tests for migration
  - [x] Test migration up successfully creates all tables
  - [x] Test migration down successfully drops all tables
  - [x] Test all indexes exist after migration
  - [x] Test foreign key constraints are enforced
- [x] Run and validate migration
  - [x] Run `make migrate-up` on test database
  - [x] Verify all tables exist with correct schema
  - [x] Verify all indexes created
  - [x] Test rollback with `make migrate-down`

---

### Dev Notes

**Database Connection:**
- Connection string from env: `DATABASE_URL=postgres://user:pass@localhost:5432/ai_tools_atlas?sslmode=disable`
- Use `golang-migrate` CLI for migrations

**Schema Requirements from Architecture.md:**
- All table names use snake_case
- All column names use snake_case
- Foreign keys enforced with appropriate cascading behavior
- Indexes on commonly queried fields (slug, ratings, counts)

**Testing Strategy:**
- Create test database for migration testing
- Verify schema matches exactly what GORM models will expect
- Test both up and down migrations

---

### Dev Agent Record

#### Implementation Plan
1. Verified golang-migrate already installed in go.mod:10
2. Created migrations directory structure
3. Created comprehensive SQL migration files with all tables, indexes, and constraints
4. Enhanced Makefile with migrate-create target
5. Created comprehensive test suite for migrations

#### Debug Log
- golang-migrate dependency was already present from story 1-1
- Makefile already had migrate-up and migrate-down targets from story 1-1
- Added migrate-create target to complete the implementation

#### Completion Notes
✅ All 11 tables created with proper schema matching Architecture.md
✅ All 9 required indexes created
✅ Foreign key constraints with appropriate ON DELETE behavior
✅ Makefile targets: migrate-up, migrate-down, migrate-create
✅ Comprehensive test suite in migrations_test.go
- Tests migration up creates all tables
- Tests migration down drops all tables
- Tests indexes exist
- Tests foreign key enforcement
- Tests CHECK constraints enforcement (media.type, tool_alternatives.relationship_type, users.role, reviews.rating_*, reviews.moderation_status)
- Tests schema validation for all 11 tables
- Tests multiple migration cycles

#### Code Review Fixes Applied (2025-12-26)
🔧 **CRITICAL:** Fixed missing github.com/lib/pq dependency - tests now compile
🔧 **MEDIUM:** Added comprehensive CHECK constraint tests covering all 5 constraint types
🔧 **MEDIUM:** Expanded schema validation to cover all 11 tables (previously only tools)
🔧 **MEDIUM:** Created backend/.env.example with DATABASE_URL and TEST_DATABASE_URL configuration

---

### File List
- backend/migrations/000001_initial_schema.up.sql (created)
- backend/migrations/000001_initial_schema.down.sql (created)
- backend/Makefile (modified)
- backend/internal/platform/db/migrations_test.go (created, enhanced during code review)
- backend/.env.example (created during code review)
- backend/go.mod (modified - added lib/pq dependency)
- backend/go.sum (modified - dependency updates)

---

### Change Log
- 2025-12-26: Story 1-2 completed - Database migrations created with full schema, indexes, and tests
- 2025-12-26: Code review completed - Fixed CRITICAL compilation issue, added comprehensive CHECK constraint tests, expanded schema validation to all tables, created .env.example

---

### Status
Done
