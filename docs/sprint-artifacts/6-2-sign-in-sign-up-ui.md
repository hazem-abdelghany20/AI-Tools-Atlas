### Story 6.2: Sign In / Sign Up UI

**As a** user,
**I want** to sign in or create an account,
**So that** I can save my bookmarks and reviews across devices.

**Acceptance Criteria:**

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

**New Files:**
- `frontend/src/components/auth/AuthModal.vue`
- `frontend/src/components/auth/SignInForm.vue`
- `frontend/src/components/auth/SignUpForm.vue`

**Files Modified:**
- `frontend/src/components/layout/HeaderNav.vue` (add auth UI)

---
