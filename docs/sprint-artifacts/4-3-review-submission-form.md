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
