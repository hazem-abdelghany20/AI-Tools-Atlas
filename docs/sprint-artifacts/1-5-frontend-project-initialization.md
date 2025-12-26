### Story 1.5: Frontend Project Initialization

**As a** developer,
**I want** the Vue 3 frontend project initialized with Vite, TypeScript, Pinia, and Tailwind,
**So that** I can start building UI components and views.

**Acceptance Criteria:**

**Given** I need the frontend project structure
**When** I initialize the Vue 3 project
**Then** the following is complete:

**Project Created:**
- `npm create vite@latest frontend -- --template vue-ts` executed
- Dependencies installed: `vue`, `vue-router@4`, `pinia`, `tailwindcss`, `postcss`, `autoprefixer`
- Project structure:
  ```
  frontend/
  ├── src/
  │   ├── main.ts
  │   ├── App.vue
  │   ├── router/index.ts
  │   ├── stores/
  │   ├── views/
  │   ├── components/
  │   ├── lib/
  │   └── assets/
  ├── index.html
  ├── vite.config.ts
  ├── tailwind.config.cjs
  ├── postcss.config.cjs
  ├── tsconfig.json
  ├── .env.example
  └── package.json
  ```

**Tailwind Configuration (`tailwind.config.cjs`):**
- Dark theme colors configured:
  ```js
  theme: {
    extend: {
      colors: {
        'dark-bg': '#05060A',
        'dark-surface': '#0A0B10',
        'dark-border': '#1F2933',
        'primary': {
          500: '#3B82F6',
          600: '#2563EB',
        },
      },
    },
  },
  ```
- RTL plugin configured: `require('tailwindcss-rtl')`
- Content paths: `./index.html`, `./src/**/*.{vue,js,ts,jsx,tsx}`

**Router Setup (`src/router/index.ts`):**
- Vue Router 4 configured with routes:
  - `/` - HomeView
  - `/tools/:slug` - ToolProfileView
  - `/categories/:slug` - SearchResultsView (category filter)
  - `/search` - SearchResultsView
  - `/compare` - CompareView
  - `/bookmarks` - BookmarksView
- History mode enabled

**Pinia Setup (`src/main.ts`):**
- Pinia initialized and registered
- App mounted to `#app`

**And:**
- TypeScript configured with Vue 3 support
- Vite configured to proxy API requests to backend (default: `http://localhost:8080`)
- Dev server starts on port 5173
- `.env.example` with `VITE_API_BASE_URL=http://localhost:8080/api/v1`

**Technical Implementation:**

- Install commands:
  ```bash
  cd frontend
  npm install pinia vue-router@4
  npm install -D tailwindcss postcss autoprefixer tailwindcss-rtl
  npx tailwindcss init -p
  ```
- Vite config proxy:
  ```ts
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
  ```

**Prerequisites:** None (parallel to backend setup)

**Files Created:**
- `frontend/src/main.ts`
- `frontend/src/router/index.ts`
- `frontend/tailwind.config.cjs`
- `frontend/vite.config.ts`
- `frontend/.env.example`

---
