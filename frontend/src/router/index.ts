import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/tools/:slug',
      name: 'tool-profile',
      component: () => import('../views/ToolProfileView.vue'),
    },
    {
      path: '/categories/:slug',
      name: 'category',
      component: () => import('../views/SearchResultsView.vue'),
      props: true,
    },
    {
      path: '/search',
      name: 'search',
      component: () => import('../views/SearchResultsView.vue'),
    },
    {
      path: '/compare',
      name: 'compare',
      component: () => import('../views/CompareView.vue'),
    },
    {
      path: '/bookmarks',
      name: 'bookmarks',
      component: () => import('../views/BookmarksView.vue'),
    },
    {
      path: '/profile',
      name: 'profile',
      component: () => import('../views/ProfileView.vue'),
      meta: { requiresAuth: true },
    },
    // Admin routes
    {
      path: '/admin/tools',
      name: 'admin-tools',
      component: () => import('../views/admin/AdminToolsView.vue'),
      meta: { requiresAuth: true, requiresAdmin: true },
    },
    {
      path: '/admin/categories',
      name: 'admin-categories',
      component: () => import('../views/admin/AdminCategoriesView.vue'),
      meta: { requiresAuth: true, requiresAdmin: true },
    },
    {
      path: '/admin/tags',
      name: 'admin-tags',
      component: () => import('../views/admin/AdminTagsView.vue'),
      meta: { requiresAuth: true, requiresAdmin: true },
    },
    {
      path: '/admin/analytics',
      name: 'admin-analytics',
      component: () => import('../views/admin/AdminAnalyticsView.vue'),
      meta: { requiresAuth: true, requiresAdmin: true },
    },
    // Moderation routes
    {
      path: '/moderation/queue',
      name: 'moderation-queue',
      component: () => import('../views/moderation/ModerationQueueView.vue'),
      meta: { requiresAuth: true, requiresAdmin: true },
    },
  ],
})

export default router
