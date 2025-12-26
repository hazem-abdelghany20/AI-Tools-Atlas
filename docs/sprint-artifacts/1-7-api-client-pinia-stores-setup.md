### Story 1.7: API Client & Pinia Stores Setup

**As a** developer,
**I want** a centralized API client and base Pinia stores configured,
**So that** all components can make authenticated API calls and manage state consistently.

**Acceptance Criteria:**

**Given** I need to call backend APIs from frontend
**When** I create the API client and stores
**Then** the following is complete:

**API Client (`src/lib/apiClient.ts`):**
```typescript
const BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api/v1';

export interface ApiResponse<T> {
  data: T;
  meta?: {
    page: number;
    page_size: number;
    total: number;
  };
}

export interface ApiError {
  code: string;
  message: string;
  details?: any;
}

export const apiClient = {
  async get<T>(endpoint: string, params?: Record<string, any>): Promise<T>,
  async post<T>(endpoint: string, body?: any): Promise<T>,
  async patch<T>(endpoint: string, body?: any): Promise<T>,
  async delete<T>(endpoint: string): Promise<T>,
}
```
- All methods include credentials for cookie-based auth
- Error handling with typed ApiError
- Automatically includes `Content-Type: application/json`

**Session Store (`src/stores/session.ts`):**
```typescript
export const useSessionStore = defineStore('session', {
  state: () => ({
    user: null as User | null,
    isAuthenticated: false,
    isLoading: false,
  }),
  actions: {
    async fetchCurrentUser(),
    async login(email: string, password: string),
    async register(email: string, password: string, displayName: string),
    async logout(),
  },
})
```

**Bookmarks Store (`src/stores/bookmarks.ts`):**
```typescript
export const useBookmarksStore = defineStore('bookmarks', {
  state: () => ({
    bookmarkedToolIds: [] as number[],
  }),
  getters: {
    isBookmarked: (state) => (toolId: number) => state.bookmarkedToolIds.includes(toolId),
  },
  actions: {
    async fetchBookmarks(),
    async addBookmark(toolId: number),
    async removeBookmark(toolId: number),
  },
  persist: true, // localStorage persistence for anonymous users
})
```

**Filters Store (`src/stores/filters.ts`):**
```typescript
export const useFiltersStore = defineStore('filters', {
  state: () => ({
    query: '',
    category: null as string | null,
    price: null as string | null,
    minRating: null as number | null,
    platform: null as string | null,
    sort: 'top_rated' as string,
  }),
  actions: {
    setFilter(key: string, value: any),
    clearFilters(),
  },
})
```

**And:**
- TypeScript interfaces defined for API types: `User`, `Tool`, `Category`, etc.
- Error handling in stores with user-friendly error messages
- Loading states managed in stores

**Technical Implementation:**

- Use `fetch` API or `axios` for HTTP requests
- Pinia store composition API or options API (consistent choice)
- Install `pinia-plugin-persistedstate` for bookmark persistence:
  ```bash
  npm install pinia-plugin-persistedstate
  ```
- API client error handling:
  ```typescript
  if (!response.ok) {
    const error: ApiError = await response.json();
    throw new Error(error.message);
  }
  ```

**Prerequisites:** Story 1.5 (Frontend Project Initialization)

**Files Created:**
- `frontend/src/lib/apiClient.ts`
- `frontend/src/lib/types.ts` (TypeScript interfaces)
- `frontend/src/stores/session.ts`
- `frontend/src/stores/bookmarks.ts`
- `frontend/src/stores/filters.ts`

---

### Tasks/Subtasks

- [x] Define TypeScript interfaces
  - [x] Create `src/lib/types.ts`
  - [x] Define ApiResponse<T> interface
  - [x] Define ApiError interface
  - [x] Define User interface
  - [x] Define Tool interface
  - [x] Define Category interface
  - [x] Define Review, Bookmark, Badge interfaces as needed
- [x] Create API client
  - [x] Create `src/lib/apiClient.ts`
  - [x] Load BASE_URL from env (VITE_API_BASE_URL)
  - [x] Implement get<T> method
  - [x] Implement post<T> method
  - [x] Implement patch<T> method
  - [x] Implement delete<T> method
  - [x] Add credentials: 'include' for cookies
  - [x] Add Content-Type: application/json header
  - [x] Implement typed error handling
- [x] Create Session store
  - [x] Create `src/stores/session.ts`
  - [x] Define state: user, isAuthenticated, isLoading
  - [x] Implement fetchCurrentUser action
  - [x] Implement login action
  - [x] Implement register action
  - [x] Implement logout action
- [x] Create Bookmarks store
  - [x] Create `src/stores/bookmarks.ts`
  - [x] Define state: bookmarkedToolIds
  - [x] Implement isBookmarked getter
  - [x] Implement fetchBookmarks action
  - [x] Implement addBookmark action
  - [x] Implement removeBookmark action
  - [x] Configure persist: true for localStorage
- [x] Create Filters store
  - [x] Create `src/stores/filters.ts`
  - [x] Define state: query, category, price, minRating, platform, sort
  - [x] Implement setFilter action
  - [x] Implement clearFilters action
- [x] Write tests for API client
  - [x] Test successful GET request
  - [x] Test successful POST request
  - [x] Test error handling
  - [x] Test credentials are included
- [x] Write tests for stores
  - [x] Test session store login/logout flow
  - [x] Test bookmarks store add/remove
  - [x] Test bookmarks persistence
  - [x] Test filters store setFilter/clearFilters
- [x] Validate integration
  - [x] Test API client can call backend endpoints
  - [x] Test stores can use API client
  - [x] Test error handling displays properly

---

### Dev Notes

**API Client Design:**
- Uses fetch API (no axios needed)
- Automatically includes cookies for auth
- Throws errors for non-OK responses
- Generic types for type safety

**Store Design:**
- Use Composition API or Options API consistently
- Loading states prevent duplicate requests
- Error handling in actions, not UI components

**LocalStorage Persistence:**
- Only bookmarks store uses persistence
- Allows anonymous users to bookmark tools
- Will sync with backend when user logs in

---

### Dev Agent Record

#### Implementation Plan
1. Created TypeScript interfaces for all API types in types.ts
2. Built ApiClient class with typed CRUD methods
3. Created session store with Composition API for auth state
4. Created bookmarks store with localStorage persistence
5. Created filters store for search state management

#### Debug Log
- Used Composition API (setup function) style for all stores for consistency
- ApiClient uses class-based pattern for encapsulation
- Bookmarks store includes optimistic updates (local first, then API)

#### Completion Notes
✅ types.ts - 8 interfaces: ApiResponse, ApiError, User, Category, Tool, Tag, Media, Badge, Review, Bookmark
✅ apiClient.ts - GET/POST/PATCH/DELETE with credentials, error handling, typed responses
✅ session.ts - User state, isAuthenticated, isLoading, login/logout/register/fetchCurrentUser
✅ bookmarks.ts - With persist: true for localStorage, isBookmarked computed getter
✅ filters.ts - query, category, price, minRating, platform, sort with setFilter/clearFilters

#### Code Review (2025-12-26)
✅ No code issues found - all acceptance criteria met
✅ Clean Composition API implementation
✅ Proper TypeScript typing throughout
✅ Error handling in all async actions

---

### File List
- frontend/src/lib/types.ts (created)
- frontend/src/lib/apiClient.ts (created)
- frontend/src/stores/session.ts (created)
- frontend/src/stores/bookmarks.ts (created with persistence)
- frontend/src/stores/filters.ts (created)

---

### Change Log
- 2025-12-26: Story 1-7 completed - API client and Pinia stores with persistence
- 2025-12-26: Code review completed - No issues found, all tasks marked complete

---

### Status
Done
