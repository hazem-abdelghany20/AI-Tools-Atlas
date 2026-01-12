# Story 9.1: Review UI & Rating Statistics Fixes

Status: ready-for-dev

## Story

As a **user viewing a tool profile**,
I want **to be able to add a review even when reviews already exist and see rating statistics update immediately after submission**,
so that **I can contribute my feedback seamlessly and see the impact of my review right away**.

## Acceptance Criteria

1. **AC1:** "Add Review" button is visible and accessible when reviews already exist on the tool profile page
2. **AC2:** After submitting a review, the rating statistics (overall average, distribution bars, dimension ratings) update immediately without requiring a page refresh
3. **AC3:** The total review count updates in real-time after submission
4. **AC4:** The newly submitted review appears in the reviews list immediately

## Tasks / Subtasks

- [ ] Task 1: Add "Add Review" button to ReviewList when reviews exist (AC: #1)
  - [ ] 1.1: Add button in the header section next to sort dropdown in `ReviewList.vue`
  - [ ] 1.2: Ensure button respects `showAddButton` prop
  - [ ] 1.3: Style consistently with existing button design

- [ ] Task 2: Fix real-time rating statistics update (AC: #2, #3)
  - [ ] 2.1: In `ToolProfileView.vue`, re-fetch tool data after review submission OR
  - [ ] 2.2: Calculate `avgRating` locally from reviews array instead of relying on `tool.avg_rating_overall`
  - [ ] 2.3: Ensure `totalReviews` count updates correctly

- [ ] Task 3: Verify review list updates (AC: #4)
  - [ ] 3.1: Confirm `fetchReviews()` is called after submission (already implemented)
  - [ ] 3.2: Test that new review appears at top when sorted by "newest"

## Dev Notes

### Bug Analysis

**Bug 1: Missing "Add Review" Button**
- **File:** `frontend/src/components/reviews/ReviewList.vue`
- **Issue:** The "Add Review" button (lines 68-75) only appears in the empty state (`v-else-if="reviews.length === 0"`)
- **When reviews exist:** The template switches to `v-else` block (line 79) which has NO add review button
- **Fix:** Add a button in the header section (around line 21) that's always visible when `showAddButton` is true

**Bug 2: Rating Statistics Not Updating**
- **File:** `frontend/src/views/ToolProfileView.vue`
- **Issue:** `handleReviewSuccess()` (line 237-242) only calls `fetchReviews()`, but doesn't refresh the tool data
- **The `avgRating` passed to ReviewList (line 64)** comes from `tool.avg_rating_overall` which is stale
- **Options for fix:**
  - Option A: Call `fetchTool()` after review submission to get updated `avg_rating_overall` from server
  - Option B: Calculate average rating locally from the reviews array (more responsive but may differ from server calculation)
- **Recommended:** Option A (call `fetchTool()`) - ensures consistency with server data

### Relevant Code Sections

**ReviewList.vue - Header (add button here):**
```vue
<!-- Line 12-42: Header section -->
<div class="flex items-center justify-between mb-6">
  <div class="flex items-center gap-3">
    <h2 class="text-xl font-bold text-white">المراجعات</h2>
    <!-- ... -->
  </div>
  <!-- Sort dropdown is here - ADD BUTTON NEXT TO THIS -->
</div>
```

**ToolProfileView.vue - Success handler (update this):**
```typescript
// Line 237-242
const handleReviewSuccess = () => {
  showReviewForm.value = false
  fetchReviews()  // Reviews refresh
  // MISSING: fetchTool() to update avg_rating_overall
}
```

### Project Structure Notes

- **Component Pattern:** Vue 3 Composition API with `<script setup>`
- **Styling:** Tailwind CSS with dark theme utilities (`bg-dark-surface`, `text-neon-blue`, etc.)
- **RTL Support:** Components use `dir="rtl"` for Arabic interface
- **State Management:** Pinia stores for bookmarks/session, local refs for component state

### References

- [Source: frontend/src/components/reviews/ReviewList.vue] - Review list component with missing button
- [Source: frontend/src/views/ToolProfileView.vue] - Tool profile view with incomplete success handler
- [Source: frontend/src/components/reviews/RatingsAggregation.vue] - Rating display component (receives avgRating as prop)

## Dev Agent Record

### Context Reference

<!-- Path(s) to story context XML will be added here by context workflow -->

### Agent Model Used

Claude Opus 4.5

### Debug Log References

### Completion Notes List

### File List

- `frontend/src/components/reviews/ReviewList.vue` - Add button, modify template
- `frontend/src/views/ToolProfileView.vue` - Update handleReviewSuccess()
