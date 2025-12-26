<template>
  <div class="admin-analytics-view min-h-screen bg-dark-bg" dir="rtl">
    <!-- Header -->
    <div class="bg-dark-surface border-b border-gray-800 px-6 py-4">
      <div class="max-w-7xl mx-auto">
        <h1 class="text-xl font-bold text-white">Admin - Analytics</h1>
        <p class="text-gray-400 text-sm">Overview of platform statistics</p>
      </div>
    </div>

    <!-- Content -->
    <div class="max-w-7xl mx-auto px-6 py-6 space-y-6">
      <!-- Loading -->
      <div v-if="loading" class="flex justify-center py-12">
        <Spinner />
      </div>

      <template v-else>
        <!-- Overview Stats -->
        <div class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-5 gap-4">
          <StatCard label="Total Tools" :value="stats?.total_tools || 0" icon="tools" />
          <StatCard label="Total Categories" :value="stats?.total_categories || 0" icon="folder" />
          <StatCard label="Total Reviews" :value="stats?.total_reviews || 0" icon="star" />
          <StatCard label="Total Bookmarks" :value="stats?.total_bookmarks || 0" icon="bookmark" />
          <StatCard label="Total Users" :value="stats?.total_users || 0" icon="users" />
        </div>

        <!-- Activity Stats -->
        <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
          <StatCard label="New Tools (Week)" :value="stats?.new_tools_week || 0" color="green" />
          <StatCard label="New Tools (Month)" :value="stats?.new_tools_month || 0" color="blue" />
          <StatCard label="New Reviews (Week)" :value="stats?.new_reviews_week || 0" color="yellow" />
          <StatCard label="New Users (Week)" :value="stats?.new_users_week || 0" color="purple" />
        </div>

        <!-- Top Tools -->
        <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <!-- By Bookmarks -->
          <div class="bg-dark-surface border border-gray-800 rounded-xl overflow-hidden">
            <div class="px-4 py-3 border-b border-gray-800 flex items-center gap-2">
              <svg class="w-5 h-5 text-red-400" fill="currentColor" viewBox="0 0 20 20">
                <path d="M3.172 5.172a4 4 0 015.656 0L10 6.343l1.172-1.171a4 4 0 115.656 5.656L10 17.657l-6.828-6.829a4 4 0 010-5.656z" />
              </svg>
              <h3 class="font-medium text-white">Top by Bookmarks</h3>
            </div>
            <div class="divide-y divide-gray-800">
              <div
                v-for="(tool, index) in topTools?.by_bookmarks || []"
                :key="tool.id"
                class="px-4 py-3 flex items-center gap-3"
              >
                <span class="text-gray-500 text-sm w-6">{{ index + 1 }}.</span>
                <img
                  v-if="tool.logo_url"
                  :src="tool.logo_url"
                  :alt="tool.name"
                  class="w-8 h-8 rounded object-cover"
                />
                <div v-else class="w-8 h-8 rounded bg-neon-blue/20 flex items-center justify-center">
                  <span class="text-neon-blue text-xs font-bold">{{ tool.name.charAt(0) }}</span>
                </div>
                <div class="flex-1 min-w-0">
                  <div class="text-white text-sm truncate">{{ tool.name }}</div>
                </div>
                <span class="text-gray-400 text-sm">{{ tool.bookmark_count }}</span>
              </div>
              <div v-if="!topTools?.by_bookmarks?.length" class="px-4 py-6 text-center text-gray-500">
                No data
              </div>
            </div>
          </div>

          <!-- By Rating -->
          <div class="bg-dark-surface border border-gray-800 rounded-xl overflow-hidden">
            <div class="px-4 py-3 border-b border-gray-800 flex items-center gap-2">
              <svg class="w-5 h-5 text-yellow-400" fill="currentColor" viewBox="0 0 20 20">
                <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"/>
              </svg>
              <h3 class="font-medium text-white">Top by Rating</h3>
            </div>
            <div class="divide-y divide-gray-800">
              <div
                v-for="(tool, index) in topTools?.by_rating || []"
                :key="tool.id"
                class="px-4 py-3 flex items-center gap-3"
              >
                <span class="text-gray-500 text-sm w-6">{{ index + 1 }}.</span>
                <img
                  v-if="tool.logo_url"
                  :src="tool.logo_url"
                  :alt="tool.name"
                  class="w-8 h-8 rounded object-cover"
                />
                <div v-else class="w-8 h-8 rounded bg-neon-blue/20 flex items-center justify-center">
                  <span class="text-neon-blue text-xs font-bold">{{ tool.name.charAt(0) }}</span>
                </div>
                <div class="flex-1 min-w-0">
                  <div class="text-white text-sm truncate">{{ tool.name }}</div>
                </div>
                <span class="text-yellow-400 text-sm">{{ tool.avg_rating?.toFixed(1) }}</span>
              </div>
              <div v-if="!topTools?.by_rating?.length" class="px-4 py-6 text-center text-gray-500">
                No data
              </div>
            </div>
          </div>

          <!-- By Reviews -->
          <div class="bg-dark-surface border border-gray-800 rounded-xl overflow-hidden">
            <div class="px-4 py-3 border-b border-gray-800 flex items-center gap-2">
              <svg class="w-5 h-5 text-blue-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
              </svg>
              <h3 class="font-medium text-white">Top by Reviews</h3>
            </div>
            <div class="divide-y divide-gray-800">
              <div
                v-for="(tool, index) in topTools?.by_reviews || []"
                :key="tool.id"
                class="px-4 py-3 flex items-center gap-3"
              >
                <span class="text-gray-500 text-sm w-6">{{ index + 1 }}.</span>
                <img
                  v-if="tool.logo_url"
                  :src="tool.logo_url"
                  :alt="tool.name"
                  class="w-8 h-8 rounded object-cover"
                />
                <div v-else class="w-8 h-8 rounded bg-neon-blue/20 flex items-center justify-center">
                  <span class="text-neon-blue text-xs font-bold">{{ tool.name.charAt(0) }}</span>
                </div>
                <div class="flex-1 min-w-0">
                  <div class="text-white text-sm truncate">{{ tool.name }}</div>
                </div>
                <span class="text-gray-400 text-sm">{{ tool.review_count }}</span>
              </div>
              <div v-if="!topTools?.by_reviews?.length" class="px-4 py-6 text-center text-gray-500">
                No data
              </div>
            </div>
          </div>
        </div>

        <!-- Top Categories -->
        <div class="bg-dark-surface border border-gray-800 rounded-xl overflow-hidden">
          <div class="px-4 py-3 border-b border-gray-800">
            <h3 class="font-medium text-white">Top Categories by Tool Count</h3>
          </div>
          <div class="p-4">
            <div class="space-y-3">
              <div
                v-for="category in topCategories"
                :key="category.id"
                class="flex items-center gap-4"
              >
                <div class="w-32 text-gray-300 truncate">{{ category.name }}</div>
                <div class="flex-1 bg-dark-bg rounded-full h-4 overflow-hidden">
                  <div
                    class="h-full bg-neon-blue rounded-full transition-all duration-500"
                    :style="{ width: `${getCategoryWidth(category.tool_count)}%` }"
                  ></div>
                </div>
                <div class="w-12 text-right text-gray-400 text-sm">{{ category.tool_count }}</div>
              </div>
              <div v-if="!topCategories?.length" class="py-6 text-center text-gray-500">
                No data
              </div>
            </div>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useSessionStore } from '@/stores/session'
import { storeToRefs } from 'pinia'
import { apiClient } from '@/lib/apiClient'
import Spinner from '@/components/common/Spinner.vue'
import StatCard from '@/components/admin/StatCard.vue'

const router = useRouter()
const sessionStore = useSessionStore()
const { user, isAuthenticated } = storeToRefs(sessionStore)

// Check admin access
watch([isAuthenticated, user], ([authenticated, currentUser]) => {
  if (!authenticated || currentUser?.role !== 'admin') {
    router.push('/')
  }
}, { immediate: true })

interface OverviewStats {
  total_tools: number
  total_categories: number
  total_reviews: number
  total_bookmarks: number
  total_users: number
  new_tools_week: number
  new_tools_month: number
  new_reviews_week: number
  new_users_week: number
}

interface TopTool {
  id: number
  slug: string
  name: string
  logo_url?: string
  bookmark_count: number
  review_count: number
  avg_rating: number
}

interface TopToolsResponse {
  by_bookmarks: TopTool[]
  by_rating: TopTool[]
  by_reviews: TopTool[]
}

interface TopCategory {
  id: number
  slug: string
  name: string
  tool_count: number
}

const loading = ref(false)
const stats = ref<OverviewStats | null>(null)
const topTools = ref<TopToolsResponse | null>(null)
const topCategories = ref<TopCategory[]>([])

const maxToolCount = computed(() => {
  if (!topCategories.value.length) return 1
  return Math.max(...topCategories.value.map(c => c.tool_count), 1)
})

function getCategoryWidth(count: number): number {
  return (count / maxToolCount.value) * 100
}

async function fetchData() {
  loading.value = true
  try {
    const [overviewRes, topToolsRes, topCatsRes] = await Promise.all([
      apiClient.get<{ data: OverviewStats }>('/admin/analytics/overview'),
      apiClient.get<{ data: TopToolsResponse }>('/admin/analytics/top-tools?limit=5'),
      apiClient.get<{ data: TopCategory[] }>('/admin/analytics/top-categories?limit=10'),
    ])
    stats.value = overviewRes.data
    topTools.value = topToolsRes.data
    topCategories.value = topCatsRes.data
  } catch (error) {
    console.error('Failed to fetch analytics:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchData()
})
</script>
