### Story 5.4: Comparison View

**As a** user,
**I want** to see tools side-by-side in a comparison table,
**So that** I can evaluate options and make a decision.

**Status:** done

---

## Acceptance Criteria:

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
- Parse tool IDs from URL: `route.query.ids.split(',')`
- Fetch all tools: `Promise.all(ids.map(id => apiClient.get(\`/tools/${id}\`)))`
- CompareTable.vue component builds table from tool data
- Remove tool: update URL, refetch remaining tools

**Prerequisites:** Story 5.3 (Comparison selection), Epic 3 (Tool data)

---

## Tasks/Subtasks:

- [x] **Task 1: Create CompareView Component**
  - [x] 1.1 Create view with RTL layout
  - [x] 1.2 Parse IDs from URL query params
  - [x] 1.3 Fetch tools by ID in parallel
  - [x] 1.4 Add loading state with Spinner
  - [x] 1.5 Add error state with retry
  - [x] 1.6 Add empty state with CTAs

- [x] **Task 2: Create CompareTable Component**
  - [x] 2.1 Create table with sticky first column
  - [x] 2.2 Add header row with tool logos and names
  - [x] 2.3 Add Best For row
  - [x] 2.4 Add Category row
  - [x] 2.5 Add Rating row with stars
  - [x] 2.6 Add Pricing row
  - [x] 2.7 Add Platforms row
  - [x] 2.8 Add Target Roles row
  - [x] 2.9 Add Bookmark Count row
  - [x] 2.10 Add Badges row
  - [x] 2.11 Add Actions row with buttons

- [x] **Task 3: Compare Actions**
  - [x] 3.1 Implement remove tool (update URL)
  - [x] 3.2 Implement clear all
  - [x] 3.3 Implement share comparison (copy URL)
  - [x] 3.4 Add success toast for share

- [x] **Task 4: Route Integration**
  - [x] 4.1 Verify /compare route exists
  - [x] 4.2 Watch for URL changes and refetch

---

## File List:
- `frontend/src/views/CompareView.vue` (modified)
- `frontend/src/components/compare/CompareTable.vue` (created)

---

## Change Log:
| Date | Change |
|------|--------|
| 2025-12-26 | Story implementation completed |
| 2025-12-26 | Code Review: Fixed to use ?tools=slugs param, fetch by slug |
