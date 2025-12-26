### Story 1.6: Base Layout Components & Dark Theme

**As a** developer,
**I want** base layout components and dark theme styling configured,
**So that** all pages have consistent Arabic-first, RTL-aware layout.

**Acceptance Criteria:**

**Given** I need base UI components matching UX Design spec
**When** I create layout components
**Then** the following components exist:

**`src/components/layout/AppShell.vue`:**
- Root layout component with `dir="rtl"` for Arabic
- Dark background using Tailwind classes: `bg-dark-bg text-white`
- Contains: `<HeaderNav />`, `<router-view />`, optional footer
- Full viewport height layout

**`src/components/layout/HeaderNav.vue`:**
- Top navigation bar with dark surface background
- Logo/brand on right (RTL), navigation links in center, user menu on left
- Navigation items: "الرئيسية" (Home), "تصفح" (Browse), "المفضلة" (Bookmarks)
- User menu (when authenticated): Avatar/initials, dropdown with "مراجعاتي" (My Reviews), "تسجيل الخروج" (Sign Out)
- User menu (not authenticated): "تسجيل الدخول" (Sign In), "إنشاء حساب" (Sign Up) buttons
- Mobile: Hamburger menu (RTL-aware icon) with drawer navigation

**`src/components/common/SkeletonCard.vue`:**
- Generic skeleton loader for cards
- Pulse animation, dark theme colors
- Props: `height`, `width`, `rounded`

**`src/components/common/Spinner.vue`:**
- Loading spinner with neon blue color
- Used for page-level loading states

**Global Styles (`src/assets/main.css`):**
- Tailwind directives imported
- Cairo font loaded (from Google Fonts or local)
- Base typography styles for Arabic: proper line-height, letter-spacing
- RTL-specific adjustments for common elements

**And:**
- `App.vue` uses `<AppShell>` as root component
- All text uses Cairo font
- Color scheme matches UX spec: Dark bg (#05060A), neon blue primary (#3B82F6)
- RTL layout tested: elements flow right-to-left naturally

**Technical Implementation:**

- Cairo font:
  ```css
  @import url('https://fonts.googleapis.com/css2?family=Cairo:wght@400;600;700&display=swap');

  body {
    font-family: 'Cairo', sans-serif;
    direction: rtl;
  }
  ```
- Use Tailwind RTL utilities or manual RTL adjustments
- Skeleton animation:
  ```css
  @keyframes pulse {
    0%, 100% { opacity: 1; }
    50% { opacity: 0.5; }
  }
  ```

**Prerequisites:** Story 1.5 (Frontend Project Initialization)

**Files Created:**
- `frontend/src/components/layout/AppShell.vue`
- `frontend/src/components/layout/HeaderNav.vue`
- `frontend/src/components/common/SkeletonCard.vue`
- `frontend/src/components/common/Spinner.vue`
- `frontend/src/assets/main.css`

---
