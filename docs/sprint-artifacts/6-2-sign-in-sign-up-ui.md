### Story 6.2: Sign In / Sign Up UI

**As a** user,
**I want** to sign in or create an account,
**So that** I can save my bookmarks and reviews across devices.

**Status:** done

---

## Acceptance Criteria:

**Sign In Modal:**
- Triggered from "تسجيل الدخول" button in header
- Form fields: Email, Password
- Submit button: "تسجيل الدخول"
- Link to sign up: "ليس لديك حساب؟ أنشئ حساباً"
- Validation: email format, password required
- Error handling: "البريد الإلكتروني أو كلمة المرور غير صحيحة"
- Success: close modal, update header UI, show success toast

**Sign Up Modal:**
- Triggered from "إنشاء حساب" button or link in sign in modal
- Form fields: Display Name, Email, Password
- Submit button: "إنشاء حساب"
- Link to sign in: "لديك حساب؟ سجّل الدخول"
- Validation: all fields required, email format, password min 8 chars
- Error handling: "البريد الإلكتروني مستخدم بالفعل"
- Success: close modal, update header UI, show success toast
- Anonymous bookmark migration: show toast "تم حفظ X أدوات إلى حسابك"

**Header UI Updates:**
- When authenticated: show avatar/initials + display name
- Dropdown menu: "مراجعاتي" (My Reviews), "المفضلة" (My Bookmarks), "تسجيل الخروج" (Sign Out)
- When not authenticated: "تسجيل الدخول", "إنشاء حساب" buttons

**Technical Implementation:**
- AuthModal.vue component with tabs for sign in/sign up
- useSessionStore actions: login, register, logout, fetchCurrentUser
- On mount: call fetchCurrentUser to restore session
- Form validation with reactive error state

**Prerequisites:** Story 6.1 (Auth endpoints)

---

## Tasks/Subtasks:

- [x] **Task 1: Fix Session Store**
  - [x] 1.1 Fix fetchCurrentUser endpoint to /me
  - [x] 1.2 Fix login response structure
  - [x] 1.3 Fix register response structure
  - [x] 1.4 Add return types to methods

- [x] **Task 2: Create SignInForm Component**
  - [x] 2.1 Create form with email/password fields
  - [x] 2.2 Add validation and error display
  - [x] 2.3 Handle submit with session store
  - [x] 2.4 Add loading state with spinner
  - [x] 2.5 Add switch to sign up link

- [x] **Task 3: Create SignUpForm Component**
  - [x] 3.1 Create form with name/email/password fields
  - [x] 3.2 Add client-side validation
  - [x] 3.3 Handle submit with session store
  - [x] 3.4 Add loading state with spinner
  - [x] 3.5 Add switch to sign in link

- [x] **Task 4: Create AuthModal Component**
  - [x] 4.1 Create modal with backdrop
  - [x] 4.2 Add mode switching (signin/signup)
  - [x] 4.3 Integrate SignInForm and SignUpForm
  - [x] 4.4 Add success toast notification
  - [x] 4.5 Handle close on success

- [x] **Task 5: Update HeaderNav**
  - [x] 5.1 Add authenticated state with avatar
  - [x] 5.2 Add dropdown menu
  - [x] 5.3 Add guest state with auth buttons
  - [x] 5.4 Integrate AuthModal
  - [x] 5.5 Fetch current user on mount
  - [x] 5.6 Add mobile menu auth options

---

## File List:
- `frontend/src/stores/session.ts` (modified)
- `frontend/src/components/auth/SignInForm.vue` (created)
- `frontend/src/components/auth/SignUpForm.vue` (created)
- `frontend/src/components/auth/AuthModal.vue` (created)
- `frontend/src/components/layout/HeaderNav.vue` (modified)

---

## Change Log:
| Date | Change |
|------|--------|
| 2025-12-26 | Story implementation completed |
| 2025-12-26 | Code Review: Added Escape key handler to AuthModal |
