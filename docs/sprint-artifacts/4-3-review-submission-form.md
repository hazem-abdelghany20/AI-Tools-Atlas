### Story 4.3: Review Submission Form

**As a** logged-in user,
**I want** to submit a structured review for a tool,
**So that** I can share my experience with others.

**Status:** done

---

## Acceptance Criteria:

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
- Success: close modal, show success toast, refresh reviews list
- Error: show error message, keep form open

---

## Tasks/Subtasks:

- [x] **Task 1: Create DimensionRating Component**
  - [x] 1.1 Create compact star rating for dimension ratings
  - [x] 1.2 Support v-model binding

- [x] **Task 2: Create ReviewForm Component**
  - [x] 2.1 Create modal overlay structure
  - [x] 2.2 Implement overall rating with interactive stars
  - [x] 2.3 Add dimension ratings section (optional)
  - [x] 2.4 Add pros/cons textareas with character counters
  - [x] 2.5 Add dropdown selects for use case, role, company size
  - [x] 2.6 Add usage context textarea
  - [x] 2.7 Implement form validation with Arabic messages
  - [x] 2.8 Handle form submission with loading state
  - [x] 2.9 Display error messages from server

- [x] **Task 3: Integrate into ToolProfileView**
  - [x] 3.1 Add showReviewForm state
  - [x] 3.2 Update handleAddReview to show form (auth check)
  - [x] 3.3 Handle form close and success events
  - [x] 3.4 Refresh reviews list on success

---

## Dev Notes:

**Technical Implementation:**
- ReviewForm.vue component with reactive form state
- useSessionStore to check authentication
- Form validation with manual validation
- API call with error handling

**Prerequisites:** Story 4.1 (Reviews API), Epic 6 (Auth required)

---

## Dev Agent Record:

### Implementation Plan:
1. Create DimensionRating component for optional dimension ratings
2. Create ReviewForm modal component with all form fields
3. Implement validation and submission logic
4. Integrate into ToolProfileView

### Completion Notes:
- Created DimensionRating.vue for compact star selection
- Created ReviewForm.vue with full modal form
- Implemented client-side validation with Arabic messages
- Added form submission with loading state and error handling
- Integrated into ToolProfileView with authentication check
- Reviews list refreshes automatically on successful submission

---

## File List:
- `frontend/src/components/reviews/DimensionRating.vue` (created)
- `frontend/src/components/reviews/ReviewForm.vue` (created)
- `frontend/src/views/ToolProfileView.vue` (modified)

---

## Change Log:
| Date | Change |
|------|--------|
| 2025-12-26 | Story implementation completed |
| 2025-12-26 | Created review form components |
| 2025-12-26 | Integrated form into tool profile |
| 2025-12-26 | Code Review: Fixed Unicode char count, added Escape key handler, visual char limit indicator |
