# Epic 9: Bug Fixes & Polish

## Epic Overview

**Goal:** Address minor bugs and UI issues discovered during implementation to ensure a polished, production-ready user experience.

**Business Value:** Improves user experience, reduces friction in key user flows (reviews, ratings), and ensures data consistency across the application.

**Dependencies:** Epics 1-8 (all core features completed)

---

## Stories

### Story 9.1: Review UI & Rating Statistics Fixes

**Priority:** High
**Complexity:** Low-Medium

**As a** user viewing a tool profile,
**I want** to be able to add a review even when reviews already exist and see rating statistics update immediately after submission,
**So that** I can contribute my feedback seamlessly and see the impact of my review right away.

#### Acceptance Criteria

1. **AC1:** "Add Review" button is visible and accessible when reviews already exist on the tool profile page
2. **AC2:** After submitting a review, the rating statistics (overall average, distribution bars, dimension ratings) update immediately without requiring a page refresh
3. **AC3:** The total review count updates in real-time after submission
4. **AC4:** The newly submitted review appears in the reviews list immediately

#### Technical Notes

- **Bug 1 Location:** `frontend/src/components/reviews/ReviewList.vue` - "Add Review" button only shown in empty state (lines 68-75), missing from reviews list view (lines 79-106)
- **Bug 2 Location:** `frontend/src/views/ToolProfileView.vue` - `handleReviewSuccess()` only refreshes reviews array, not the tool's `avg_rating_overall` prop
- **Solution Approach:**
  - Add persistent "Add Review" button in ReviewList header when reviews exist
  - Re-fetch tool data OR calculate ratings locally from reviews array after submission
