### Story 4.2: Review Display on Tool Profile

**As a** user,
**I want** to see structured reviews on the tool profile,
**So that** I can learn from others' experiences.

**Acceptance Criteria:**

**Reviews Section on ToolProfileView:**
- Heading: "المراجعات" (Reviews) with overall rating summary
- Sort controls: Newest, Most helpful, Highest rated, Lowest rated
- Each review card shows:
  - Reviewer display name + role + company size
  - Overall rating (stars) + submission date
  - Pros/Cons in structured format ("الإيجابيات" / "السلبيات")
  - Primary use case badge
  - Usage context tags
  - Helpful count with "مفيد؟" (Helpful?) button (future feature)
- Pagination: Load more button or infinite scroll
- Empty state: "لا توجد مراجعات بعد. كن أول من يراجع!" with "Add Review" button

**Visual Design:**
- Review cards with dark background, subtle borders
- Pros in green-tinted section, cons in red-tinted section
- Reviewer info in lighter text
- Stars displayed clearly for rating

**Technical Implementation:**
- Fetch reviews: `apiClient.get(\`/tools/${slug}/reviews\`, { params: { sort, page } })`
- ReviewList.vue component iterates over reviews
- Sort state managed locally or in URL params

**Prerequisites:** Story 4.1, Story 3.3

**Files Modified:**
- `frontend/src/views/ToolProfileView.vue`

**New Files:**
- `frontend/src/components/reviews/ReviewList.vue`
- `frontend/src/components/reviews/ReviewCard.vue`

---
