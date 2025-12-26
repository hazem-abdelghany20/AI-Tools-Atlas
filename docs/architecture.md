---
inputDocuments:
  - docs/prd.md
  - docs/ux-design-specification.md
  - docs/analysis/brainstorming-session-2025-12-01T13-32-48.md
workflowType: 'architecture'
stepsCompleted: [1, 2, 3, 4, 5, 6, 7]
lastStep: 7
project_name: 'AI Tools Atlas'
user_name: 'Hazzouma'
date: '2025-12-02T11:56:45+02:00'
---

# Architecture Decision Document

_This document builds collaboratively through step-by-step discovery. Sections are appended as we work through each architectural decision together._

## Project Context Analysis

### Requirements Overview

**Functional Requirements:**
- Provide a public, web-based AI tools directory where users can discover tools by free-text search, categories, and pre-defined use cases.
- Offer structured search and filtering (category, sub-use case, price, minimum rating, platform, experience level) with sort options (top rated, most bookmarked, trending, newest).
- Expose rich, SaaS-style tool profiles with:
  - “Best for…” positioning, long descriptions, primary use cases, target roles/industries.
  - Pricing summaries (free tier, starting price, billing model).
  - Media (screenshots, YouTube videos) and key features.
  - Structured alternatives (“Similar tools” and “Alternatives to X” driven by category/tags/behavior).
- Support structured reviews and ratings with:
  - Overall rating and dimensional ratings (ease of use, value, accuracy/quality, speed/performance, support/community).
  - Pros/cons, primary use case, reviewer role, company size, and usage context.
- Allow users to bookmark tools as a shortlist and later evolve this into named “stacks” of tools per workflow.
- Enable side-by-side comparison of 2–4 tools with a comparison table across:
  - Overview (“best for…”, primary use cases).
  - Features.
  - Pricing (free tier, starting price, billing model).
  - Ratings and social proof (review counts, bookmark counts, badges).
- Provide basic account/light identity model to support persistent bookmarks (with a simpler local-storage fallback for v1 if needed).
- Prepare for future phases:
  - Private and shareable stacks built on top of bookmarks.
  - Updates/feed surfaces based on followed tools/categories/stacks.
  - Admin and moderation tooling for catalog changes and reviews.

**Non-Functional Requirements:**
- **Performance & Perceived Performance:**
  - Search, filters, and navigation between results and profiles should feel snappy; heavy media (YouTube, screenshots) must not block core interactions.
  - Use techniques like skeleton states, lazy loading, and incremental data fetching to keep the UI responsive.
- **Trust & Transparency:**
  - Rankings, badges, and “best for…” messaging should be clearly grounded in visible data (ratings, reviews, bookmarks) to avoid a “black box” feel.
- **Responsiveness & Accessibility:**
  - Full-feature parity on desktop and mobile with a responsive layout.
  - Arabic-first, RTL-aware design with readable typography and contrast; aim for solid accessibility baselines.
- **Scalability:**
  - Data model and APIs should handle growth in tools, reviews, and bookmarks without major refactors, especially with stacks and feeds in mind.
- **Reliability & Data Integrity:**
  - Safe handling of review and bookmark submissions, with guardrails against obvious abuse in later phases.
- **SEO & Shareability:**
  - Tool profiles and category pages should be linkable and indexable to support organic discovery.

**Scale & Complexity:**
- Primary domain: public web directory + JSON API.
- Complexity level: Low–Medium for the MVP; tends toward Medium as stacks, feeds, and social features mature.
- Estimated architectural components (initially within a single deployable backend app):
  - Catalog and metadata (tools, categories, tags, media, alternatives, badges).
  - Search and filtering (queries, ranking, facets).
  - Reviews and ratings.
  - Bookmarks (and future stacks).
  - User/account and session handling (for light auth).
  - Analytics and basic signals for popularity/trending.
  - Admin/moderation capabilities (at least minimal at first).

### Technical Constraints & Dependencies

- Backend:
  - Language: Go.
  - Framework: Gin for JSON APIs and HTTP routing.
  - Database: Postgres implementing the relational model (Tool, Category, Tag, Review, Bookmark, Media, Badge, ToolAlternative, etc.).
  - Data access: GORM as the primary ORM for queries and persistence (using raw SQL only where needed for performance).
  - Auth: JWT-based auth with cookie support for the web client (light accounts; anonymous session IDs for early bookmark behavior if needed).
- Frontend:
  - Framework: Vue 3 + Vite + TypeScript.
  - State management: Pinia for bookmarks, filters, and session state.
  - Styling: Tailwind CSS or a Vue UI library for fast, card-heavy UI implementation.
  - Requirements: responsive, RTL-aware, dark-theme design with strong perceived performance.
- Platform:
  - Primary surface is a responsive web app; future mobile app or browser extension should reuse the same backend and data model.
- Roadmap-aware constraints:
  - The data model must anticipate stacks and feeds (collections referencing tools and users, events for updates) even if these features are not exposed in MVP.
  - Search and ranking logic should be designed so it can evolve (e.g., from simple sort + filters to more advanced relevance or recommendation signals).

### Cross-Cutting Concerns Identified

- Search and ranking logic powering:
  - “Top rated”, “Most bookmarked”, “Trending”, and “Newest” sorts.
  - Badges like “Top in Category”, “Rising Star”, “Most Bookmarked”.
- Localization and RTL:
  - Arabic-first interface, with possible English/other language content mixed in; affects routing, layout, typography, and validation/error messaging.
- Analytics and instrumentation:
  - Tracking search usage, filters, bookmarks, compares, profile visits, and review submissions to validate the decision-first experience.
- Content quality and moderation:
  - Managing tool catalog correctness and review quality; becomes more important as volume grows.
- SEO and linkability:
  - Clean slugs and structured metadata for tool and category pages.
- Performance and caching:
  - Potential caching of popular searches, categories, and tool profiles as catalog and traffic grow.

## Starter Template Evaluation

### Primary Technology Domain

Web application with a Go + Gin + GORM JSON API backend, Postgres as the primary data store, and a Vue 3 + Vite + TypeScript SPA frontend styled with Tailwind CSS.

### Starter Options Considered

- **External starters (Gin + Vue templates):**
  - Considered using community starter kits that wire up Gin APIs with a Vue/Vite frontend.
  - Not selected because we want explicit, minimal architecture decisions documented in this repo and we are not relying on web search to validate individual starter health or versions.
- **Custom baseline (selected):**
  - A lightweight, hand-rolled Gin + GORM backend and Vite/Vue 3 frontend that align directly with the PRD data model and UX flows.
  - Keeps dependencies boring and transparent, and makes it easier for future AI agents and contributors to understand and extend the system.

### Selected Starter: Custom Gin + GORM + Vue 3 Baseline

**Rationale for Selection:**
- Matches the hard technical preferences (Go, Gin, GORM, Postgres, Vue 3, Pinia, Tailwind).
- Avoids hidden architectural decisions from third-party starters; all structure and conventions are explicitly captured in this architecture.
- Simple enough to support a Low–Medium complexity MVP, but capable of scaling into stacks, feeds, and additional surfaces later.
- Works well with an “API-first” mindset: backend provides a clean JSON API that can serve web, future mobile, and browser-extension clients.

**Initialization Command:**

```bash
# Backend (from backend folder)
mkdir ai-tools-atlas-backend && cd ai-tools-atlas-backend
go mod init github.com/your-org/ai-tools-atlas-backend
go get github.com/gin-gonic/gin
go get gorm.io/gorm
go get gorm.io/driver/postgres

# Frontend (from project root or separate folder)
npm create vite@latest ai-tools-atlas-frontend -- --template vue-ts
cd ai-tools-atlas-frontend
npm install
npm install pinia vue-router@4 tailwindcss postcss autoprefixer
npx tailwindcss init -p
```

**Architectural Decisions Provided by Starter:**

**Language & Runtime:**
- Backend in Go (Gin + GORM, Postgres driver), exposing JSON APIs.
- Frontend in TypeScript with Vue 3 and Vite, consuming the API.

**Styling Solution:**
- Tailwind CSS as the primary styling system, configured for a dark, Arabic-first UI and responsive breakpoints.

**Build Tooling:**
- Go modules for backend dependency management.
- Vite for fast dev server and optimized production builds of the Vue SPA.

**Testing Framework:**
- Go `testing` (and `net/http/httptest`) for API and handler tests.
- Frontend tests using the tooling provided by the Vite + Vue + TS template, extendable with Vitest and Vue Testing Library as needed.

**Code Organization:**
- Backend organized with a clean structure (e.g., `cmd/api`, `internal/app`, `internal/tools`, `internal/reviews`, `internal/platform/http`, `internal/platform/db`).
- Frontend organized into `src/views` (Home, Results, ToolProfile, Compare, Bookmarks), `src/components`, `src/stores` (Pinia), and `src/router`.

**Development Experience:**
- Local development using `go run ./cmd/api` (or a file-watcher tool) for the backend and `npm run dev` for the frontend.
- Clear separation of concerns between API and client, with a simple path to add additional clients (mobile, extension) against the same backend.

## Core Architectural Decisions

### Decision Priority Analysis

**Critical Decisions (Block Implementation):**
- Single Postgres database as the system of record, modeled via GORM entities that mirror the PRD data model (Tool, Category, Tag, Media, Review, Bookmark, Badge, ToolAlternative).
- Gin-based JSON REST API with a versioned prefix (`/api/v1`) exposing read-heavy endpoints for tools, categories, search, reviews, bookmarks, and compare.
- JWT-based authentication with HTTP-only cookies for logged-in users, plus anonymous session identifiers for pre-account bookmark flows.
- Vue 3 SPA (Vite + TypeScript) as the primary client, using Pinia for global state (session, filters, bookmarks) and Vue Router for navigation.

**Important Decisions (Shape Architecture):**
- Use GORM for data access with explicit model definitions and indexes tuned for search/filter use cases.
- Apply database migrations via an external migration tool (e.g., golang-migrate) managed alongside the Go codebase.
- Standardized API error envelope (e.g., `{ error: { code, message, details } }`) and pagination/filter conventions across all list endpoints.
- Tailwind CSS as the primary styling system, configured for dark, Arabic-first layouts and responsive breakpoints.
- Centralized logging and request instrumentation via Gin middleware, with structured logs suitable for cloud hosting.

**Deferred Decisions (Post-MVP):**
- Dedicated search engine or external search service (e.g., OpenSearch/Algolia) for more advanced relevance; MVP relies on Postgres + indices.
- Distributed caching layer (e.g., Redis) for hot lists, trending queries, or heavy profile pages; MVP relies on DB indexes and basic in-process caching if needed.
- Full-featured RBAC beyond simple roles (`user`, `admin`) and more granular permissions for moderation workflows.

### Data Architecture

- **Database & Modeling:**
  - Primary store: Postgres, single logical database for MVP.
  - Entities: `Tool`, `Category`, `Tag`, `ToolTag`, `Media`, `Review`, `Bookmark`, `Badge`, `ToolBadge`, `ToolAlternative`, plus basic `User` and `Session`/`Auth` tables.
  - IDs: use numeric or UUID primary keys, with external-facing slugs for SEO on tools and categories.
  - GORM models in an internal `data` or `models` package, with explicit column definitions and indexes for common filters (category, pricing_model, avg_rating, bookmark_count, trending_score).
- **Validation & Integrity:**
  - Request-level validation using Gin’s binding + a validation library (e.g., go-playground/validator) for incoming JSON.
  - Database-level constraints (NOT NULL, foreign keys, uniqueness on slugs) to maintain catalog integrity.
- **Migrations:**
  - Schema changes managed via a migration tool (e.g., golang-migrate) with versioned migration files committed to the repo.
  - Migrations applied as part of deployment, separate from runtime.
- **Caching Strategy:**
  - MVP: rely on tuned SQL queries and indexes; optionally use simple in-memory caching inside the process for “popular this week/top rated” strips.
  - Future: introduce Redis or similar if query load or latency becomes an issue, especially for home page and category aggregates.

### Authentication & Security

- **Authentication:**
  - Email/password–based accounts for users who want persistent bookmarks and reviews.
  - JWT tokens issued on login and stored in HTTP-only, SameSite cookies for the SPA, with short-lived access tokens and optional refresh tokens.
  - Support for anonymous session IDs (stored in cookie/local storage) for early bookmark usage before account creation.
- **Authorization:**
  - Simple role model: `user` (default) and `admin` (for catalog management and moderation).
  - Gin middleware for enforcing auth on protected routes (review submission, bookmarks, admin operations).
- **Security Middleware & Practices:**
  - Gin middleware for logging, recovery, CORS, and auth.
  - HTTPS enforced at the edge (load balancer or hosting platform).
  - Basic rate limiting (via middleware or ingress-level config) on write-heavy endpoints (reviews, bookmarks) to reduce abuse.
- **API Security:**
  - All APIs under `/api/v1` expecting and returning JSON.
  - Input sanitization and consistent error responses; no sensitive internals leaked through error messages.

### API & Communication Patterns

- **API Style:**
  - REST-like JSON API with clear resource naming: `/tools`, `/categories`, `/reviews`, `/bookmarks`, `/compare`, etc.
  - Query-based filtering and sorting on list endpoints (e.g., `?category=…&price=…&min_rating=…&sort=top_rated`).
- **Key Endpoints (MVP):**
  - Catalog: `GET /tools`, `GET /tools/:slug`, `GET /categories`, `GET /categories/:slug/tools`.
  - Search: `GET /search/tools?q=…` with parameters mapping free-text queries into categories/use cases.
  - Reviews: `GET /tools/:slug/reviews`, `POST /tools/:slug/reviews` (authenticated).
  - Bookmarks: `GET/POST/DELETE /me/bookmarks` and a compare endpoint taking multiple tool IDs/slugs.
  - Metadata: endpoints for tags, badges, and alternatives to support profile and comparison views.
- **Error Handling & Documentation:**
  - Centralized error-handling helper in Gin, mapping domain errors to HTTP status codes.
  - API documentation generated via Swagger/OpenAPI tooling compatible with Gin, shared with frontend and future clients.

### Frontend Architecture

- **Routing & Views:**
  - Vue Router with routes for Home/Discovery, Search/Results, Tool Profile, Compare, and Bookmarks/Shortlist.
  - SEO-friendly routes (e.g., `/tools/:slug`, `/categories/:slug`) with client-side fetching.
- **State Management:**
  - Pinia stores for user session, bookmarks, filters/sorts, and comparison selections.
  - API client module encapsulating calls to the Gin backend with consistent error handling.
- **Components & Performance:**
  - Reusable components for tool cards, filter panels, comparison tables, rating displays, and review blocks.
  - Skeleton loaders and lazy loading of heavy components (e.g., embedded YouTube) to maintain perceived performance.
  - RTL-aware layout and typography with Tailwind, respecting the dark visual design defined in the UX spec.

### Infrastructure & Deployment

- **Hosting Strategy (MVP):**
  - Backend: containerized Go service (Gin + GORM + Postgres driver) deployed as a single app behind HTTPS on a managed platform or VM.
  - Database: managed Postgres instance for durability and backups.
  - Frontend: static build of the Vue app served via a CDN/static hosting service, configured to talk to the `/api/v1` backend origin.
- **Configuration & Environments:**
  - Environment-specific configuration via environment variables (DB URL, JWT secrets, allowed origins, feature flags).
  - Separate `dev`, `staging`, and `prod` environments planned, even if only `dev` and `prod` exist initially.
- **Monitoring & Logging:**
  - Structured logs from Gin and application code, shipped to a log aggregation service in production.
  - Basic health check endpoint (e.g., `GET /health`) for uptime monitoring.
- **Scaling Strategy:**
  - MVP: vertical scaling on a single backend instance with Postgres.
  - Future: horizontal scaling via multiple backend replicas behind a load balancer and database tuning as catalog and traffic grow.

## Implementation Patterns & Consistency Rules

### Pattern Categories Defined

**Critical Conflict Points Identified:**
- Naming conventions for DB tables/columns, API routes and payloads, and Vue components/files.
- API response and error envelope shapes.
- Placement of domain logic vs. handlers vs. helpers in the Go codebase.
- State management patterns in Pinia stores (for bookmarks, filters, session).
- Loading/error handling patterns in the frontend, especially around search, results, and profiles.

### Naming Patterns

**Database Naming Conventions:**
- Tables: all lowercase, plural snake_case (e.g., `tools`, `categories`, `tool_tags`, `reviews`, `bookmarks`).
- Columns: snake_case (e.g., `tool_id`, `user_id`, `primary_category_id`, `avg_rating_overall`).
- Foreign keys: `<referenced>_id` (e.g., `tool_id`, `category_id`, `user_id`).
- Indexes: `idx_<table>_<column1>[_column2]` (e.g., `idx_tools_primary_category_id`, `idx_reviews_tool_id`).

**API Naming Conventions:**
- Base path: `/api/v1`.
- Resources: plural kebab-case paths (e.g., `/tools`, `/categories`, `/bookmarks`, `/reviews`).
- Resource identifiers: prefer slugs in URLs where appropriate (e.g., `/tools/:slug`, `/categories/:slug`) and IDs for internal references.
- Query params: snake_case (e.g., `min_rating`, `price`, `platform`, `sort`, `page`, `page_size`).
- HTTP methods follow REST semantics: `GET` for retrieval, `POST` for creation, `PATCH` for partial update (if needed later), `DELETE` for deletion.

**Code Naming Conventions:**
- Go packages: simple lowercase (e.g., `tools`, `reviews`, `bookmarks`, `platform`, `db`).
- Go types: PascalCase (e.g., `Tool`, `Category`, `Review`, `BookmarkService`).
- Go functions: camelCase (e.g., `getTool`, `listTools`, `createReview`), with exported functions for handlers/services.
- Vue components: PascalCase (e.g., `ToolCard.vue`, `ToolProfileView.vue`, `SearchResultsView.vue`).
- Vue/TS files: match component name (e.g., `ToolCard.vue`), utilities in camelCase filenames (e.g., `formatRatings.ts`).

### Structure Patterns

**Project Organization (Backend):**
- Single Go module with a clear internal layout, for example:
  - `cmd/api` – main entrypoint wiring Gin, config, DB, and routes.
  - `internal/platform/db` – DB connection, migrations bootstrap, GORM setup.
  - `internal/platform/http` – Gin router setup, middleware, common response helpers.
  - `internal/tools`, `internal/categories`, `internal/reviews`, `internal/bookmarks` – each containing models, repository functions, and service logic.
- Tests co-located with code using `_test.go` files in the same packages.

**Project Organization (Frontend):**
- `src/views` – route-level pages (Home, Results, ToolProfile, Compare, Bookmarks).
- `src/components` – reusable UI components (cards, filters, badges, tables).
- `src/stores` – Pinia stores for `useSessionStore`, `useBookmarksStore`, `useFiltersStore`, etc.
- `src/router` – Vue Router configuration with named routes.
- `src/lib` or `src/utils` – shared formatting/helpers (e.g., rating formatting, API client).

### Format Patterns

**API Response Formats:**
- Successful list responses:
  - `{ "data": [...], "meta": { "page": number, "page_size": number, "total": number } }`
- Successful single resource responses:
  - `{ "data": { ...resourceFields } }`
- Error responses:
  - `{ "error": { "code": string, "message": string, "details": object|null } }`
- Dates and timestamps:
  - ISO 8601 strings in UTC in APIs; formatting for display handled in the frontend.

**Data Exchange Formats:**
- JSON field names in snake_case to match query param style (e.g., `avg_rating_overall`, `has_free_tier`).
- Booleans for flags (e.g., `has_free_tier`, `is_featured`), numerics for counts and ratings.

### Communication Patterns

**State Management Patterns:**
- Pinia stores own:
  - Session/auth state (current user, tokens where applicable).
  - Bookmark state (list of tool IDs/objects and labels).
  - Search/filter state (current query, filters, sort).
- All API calls go through a centralized client module that:
  - Handles base URL configuration.
  - Applies auth headers/cookies where needed.
  - Normalizes errors into a consistent shape for the UI.

**Event/Interaction Patterns:**
- Within the app, prefer explicit function calls and store actions rather than ad-hoc event buses.
- Component outputs use clearly named events/props (e.g., `@add-to-compare`, `@bookmark`, `@filter-change`).

### Process Patterns

**Error Handling Patterns:**
- Backend:
  - Handlers use a shared responder for success and error, always following the agreed envelope.
  - Domain errors mapped to appropriate HTTP status codes (400, 401, 403, 404, 422, 500).
- Frontend:
  - Errors surfaced via a small set of UI patterns: inline messages on forms and a global toast/notification system.
  - User-facing messages concise, avoiding leaking technical details.

**Loading State Patterns:**
- For key flows (search, results, profile, compare), components:
  - Maintain explicit `isLoading` and `hasError` flags.
  - Use skeleton loaders for cards and profiles instead of spinners where possible.
  - Avoid overlapping parallel loaders that confuse the user; one primary loader per view.

### Enforcement Guidelines

**All AI Agents MUST:**
- Follow the database, API, and code naming conventions exactly as defined here.
- Use the shared API response/error envelope for all new endpoints.
- Place new domain logic, models, and handlers within the established backend and frontend folders, extending existing patterns instead of inventing new ones.

**Pattern Enforcement:**
- New PRs and AI-generated changes should be checked against this section of `docs/architecture.md`.
- Any deviation from patterns must either:
  - Be corrected to match the existing patterns, or
  - Be explicitly discussed and result in an update to this document.

## Project Structure & Boundaries

### Complete Project Directory Structure

```text
ai-tools-atlas/
├── README.md
├── docs/
│   ├── prd.md
│   ├── ux-design-specification.md
│   ├── architecture.md
│   └── analysis/
│       └── brainstorming-session-2025-12-01T13-32-48.md
├── backend/
│   ├── go.mod
│   ├── go.sum
│   ├── Makefile
│   ├── .env.example
│   ├── cmd/
│   │   └── api/
│   │       └── main.go
│   ├── internal/
│   │   ├── platform/
│   │   │   ├── config/
│   │   │   │   └── config.go
│   │   │   ├── db/
│   │   │   │   └── db.go
│   │   │   └── http/
│   │   │       ├── router.go
│   │   │       ├── middleware.go
│   │   │       └── responses.go
│   │   ├── tools/
│   │   │   ├── model.go
│   │   │   ├── repository.go
│   │   │   ├── service.go
│   │   │   └── handler.go
│   │   ├── categories/
│   │   │   ├── model.go
│   │   │   ├── repository.go
│   │   │   ├── service.go
│   │   │   └── handler.go
│   │   ├── reviews/
│   │   │   ├── model.go
│   │   │   ├── repository.go
│   │   │   ├── service.go
│   │   │   └── handler.go
│   │   ├── bookmarks/
│   │   │   ├── model.go
│   │   │   ├── repository.go
│   │   │   ├── service.go
│   │   │   └── handler.go
│   │   ├── auth/
│   │   │   ├── model.go
│   │   │   ├── service.go
│   │   │   └── handler.go
│   │   └── badges/
│   │       ├── model.go
│   │       ├── repository.go
│   │       ├── service.go
│   │       └── handler.go
│   └── migrations/
│       └── (timestamped migration files)
└── frontend/
    ├── package.json
    ├── tsconfig.json
    ├── vite.config.ts
    ├── tailwind.config.cjs
    ├── postcss.config.cjs
    ├── index.html
    ├── .env.example
    └── src/
        ├── main.ts
        ├── router/
        │   └── index.ts
        ├── stores/
        │   ├── session.ts
        │   ├── bookmarks.ts
        │   └── filters.ts
        ├── views/
        │   ├── HomeView.vue
        │   ├── SearchResultsView.vue
        │   ├── ToolProfileView.vue
        │   ├── CompareView.vue
        │   └── BookmarksView.vue
        ├── components/
        │   ├── layout/
        │   │   ├── AppShell.vue
        │   │   └── HeaderNav.vue
        │   ├── tools/
        │   │   ├── ToolCard.vue
        │   │   ├── ToolSummary.vue
        │   │   └── ToolMediaGallery.vue
        │   ├── reviews/
        │   │   ├── ReviewList.vue
        │   │   └── ReviewForm.vue
        │   ├── compare/
        │   │   └── CompareTable.vue
        │   └── common/
        │       ├── RatingStars.vue
        │       ├── BadgePill.vue
        │       ├── SkeletonCard.vue
        │       └── Spinner.vue
        ├── lib/
        │   ├── apiClient.ts
        │   ├── toolMappers.ts
        │   └── formatting.ts
        └── assets/
            └── (images, logos, icons)
```

### Architectural Boundaries

**API Boundaries:**
- All external API endpoints exposed under `/api/v1` from the backend `cmd/api` server.
- Feature-specific handlers live with their domain packages (`tools`, `reviews`, `bookmarks`, etc.), all using the shared response helpers in `internal/platform/http`.

**Component Boundaries:**
- Backend domain packages (`tools`, `categories`, `reviews`, `bookmarks`, `auth`, `badges`) encapsulate models, repositories, services, and HTTP handlers for their slice of the domain.
- Frontend views correspond to high-level UX flows; components in `components/tools`, `components/reviews`, and `components/compare` are shared across views.

**Service & Data Boundaries:**
- Data access is centralized via GORM in each domain package and wired through `internal/platform/db`.
- Cross-cutting concerns (config, logging, HTTP middleware) live in `internal/platform`.
- The database schema is owned by the backend `migrations` directory; no direct DB access from the frontend.

### Requirements to Structure Mapping

- Discovery (home, search, categories) → `backend/internal/tools`, `backend/internal/categories`, `frontend/src/views/HomeView.vue`, `SearchResultsView.vue`, and `components/tools/*`.
- Tool profile (rich details, media, alternatives) → `backend/internal/tools`, `backend/internal/badges`, `frontend/src/views/ToolProfileView.vue`, `components/tools/*`, `components/reviews/ReviewList.vue`.
- Reviews & ratings → `backend/internal/reviews`, `frontend/components/reviews/*`, and `frontend/src/views/ToolProfileView.vue`.
- Bookmarks/shortlist and compare → `backend/internal/bookmarks`, `frontend/src/stores/bookmarks.ts`, `CompareView.vue`, and `components/compare/CompareTable.vue`.
- Auth & accounts (for persistent bookmarks/reviews) → `backend/internal/auth`, `frontend/src/stores/session.ts`, and any auth-related components.

### Integration Points

- Internal communication:
  - Backend: domain services call each other via Go interfaces where needed (e.g., bookmarks querying tools), keeping package dependencies explicit.
  - Frontend: views talk to Pinia stores and the shared `apiClient`, which in turn calls the Gin API.
- External integrations:
  - YouTube and external tool websites are referenced via URLs from the `Media` and `Tool` entities; embedding handled in frontend components.
- Data flow:
  - Requests flow `frontend (views/components) → apiClient → Gin handlers → services → GORM/DB` and back as JSON responses shaped by the response helpers.

### File Organization Patterns

- Configuration files:
  - Backend: `.env`, `.env.example`, and config loading in `internal/platform/config`.
  - Frontend: `.env`, `.env.example`, Vite and Tailwind configs at the project root of `frontend/`.
- Tests:
  - Backend: `_test.go` files next to the code under `internal/*`.
  - Frontend: tests co-located next to components/views or in a parallel `tests/` folder, following the same directory layout.
- Assets:
  - Frontend static assets under `frontend/src/assets`; any shared branding or documentation assets can live under `docs/` as needed.

## Architecture Validation & Readiness

### Coherence Validation ✅

- Technology choices (Go + Gin + GORM + Postgres + Vue 3 + Pinia + Tailwind) are compatible and aligned with a web directory/search product.
- Implementation patterns (naming, structure, API envelopes) match the chosen stack and are consistent across backend and frontend.
- Project structure (backend/internal domain packages and frontend views/components/stores) supports the flows defined in the PRD and UX spec.

### Requirements Coverage ✅

- All MVP functional requirements are mapped:
  - Discovery/search/results → tools, categories, tags, filters, and ranking endpoints plus corresponding views/stores.
  - Tool profiles → tools, media, alternatives, badges, reviews, and related UI components.
  - Reviews and ratings → reviews domain, auth, and profile UI.
  - Bookmarks/shortlist and compare → bookmarks domain, compare logic, and related stores/views.
- Non-functional requirements are addressed architecturally:
  - Performance via indexed queries, simple caching options, and frontend perceived-performance patterns.
  - Security via JWT + cookies, roles, middleware, and consistent API error handling.
  - Scalability via a clean data model, API-first design, and a deployment model that can scale out later.

### Implementation Readiness ✅

- Decisions, patterns, and structure are concrete enough for AI agents and humans to implement consistently without inventing new conventions.
- Key conflict points (naming, responses, layout, state management) are explicitly covered by consistency rules.
- Directory structure specifies where new features, endpoints, and components should live.

### Gap Analysis

- No critical gaps blocking MVP implementation identified.
- Known deferred/Phase 2+ areas:
  - Dedicated search engine or recommendation system beyond Postgres queries.
  - Stacks and feed/update surfaces (data model is prepared; features will be layered later).
  - More advanced moderation tooling and fine-grained RBAC.

### Readiness Assessment

- **Overall Status:** READY FOR IMPLEMENTATION  
- **Confidence Level:** High for MVP as defined in the PRD.

AI agents implementing this project should:

- Follow the documented architecture, patterns, and structure sections in this file as the single source of truth.
- Treat any ambiguity as a prompt to update this document rather than introducing ad-hoc patterns.
