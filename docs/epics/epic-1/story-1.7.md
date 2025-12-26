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
