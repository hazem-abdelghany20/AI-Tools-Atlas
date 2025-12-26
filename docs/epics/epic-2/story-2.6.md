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
