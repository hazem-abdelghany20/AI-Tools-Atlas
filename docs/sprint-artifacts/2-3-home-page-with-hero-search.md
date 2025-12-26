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
