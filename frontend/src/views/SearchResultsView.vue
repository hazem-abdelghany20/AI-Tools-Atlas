<template>
  <div class="search-results-view min-h-screen bg-dark-bg py-8">
    <div class="container mx-auto px-4">
      <div class="flex flex-col lg:flex-row gap-6" dir="rtl">
        <!-- Filters Panel (Desktop) -->
        <div class="hidden lg:block w-72 flex-shrink-0">
          <FiltersPanel
            :filters="filters"
            @update:filters="handleFilterChange"
          />
        </div>

        <!-- Mobile Filter Button -->
        <div class="lg:hidden flex items-center justify-between mb-4">
          <button
            @click="showMobileFilters = true"
            class="flex items-center gap-2 px-4 py-2 bg-dark-surface border border-gray-700 rounded-lg text-gray-300"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z" />
            </svg>
            تصفية
          </button>
        </div>

        <!-- Mobile Filters Drawer -->
        <div
          v-if="showMobileFilters"
          class="fixed inset-0 z-50 lg:hidden"
        >
          <div class="absolute inset-0 bg-black/50" @click="showMobileFilters = false"></div>
          <div class="absolute left-0 top-0 h-full w-80 bg-dark-bg p-4 overflow-y-auto">
            <div class="flex items-center justify-between mb-4">
              <h3 class="text-lg font-semibold text-white">تصفية النتائج</h3>
              <button @click="showMobileFilters = false" class="text-gray-400 hover:text-white">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
            <FiltersPanel
              :filters="filters"
              @update:filters="handleFilterChange"
            />
          </div>
        </div>

        <!-- Results Area -->
        <div class="flex-1">
          <!-- Header -->
          <div class="flex items-center justify-between mb-6">
            <div>
              <h1 class="text-2xl font-bold text-white">
                {{ pageTitle }}
              </h1>
              <p class="text-gray-400 mt-1" v-if="!loading">
                {{ totalResults }} نتيجة
              </p>
            </div>

            <!-- Sort Dropdown -->
            <select
              v-model="filters.sort"
              @change="handleFilterChange(filters)"
              class="px-4 py-2 bg-dark-surface border border-gray-700 rounded-lg
                     text-white text-sm focus:border-neon-blue focus:outline-none"
            >
              <option value="top_rated">الأعلى تقييماً</option>
              <option value="most_bookmarked">الأكثر حفظاً</option>
              <option value="newest">الأحدث</option>
              <option value="trending">الرائج</option>
            </select>
          </div>

          <!-- Loading Skeletons -->
          <div
            v-if="loading"
            class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4"
          >
            <SkeletonCard
              v-for="i in 9"
              :key="i"
              height="200px"
              rounded
            />
          </div>

          <!-- Results Grid -->
          <div
            v-else-if="tools.length > 0"
            class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4"
          >
            <ToolCard
              v-for="tool in tools"
              :key="tool.id"
              :tool="tool"
              :is-in-compare="isInCompare(tool.id)"
              @add-to-compare="addToCompare"
              @remove-from-compare="removeFromCompare"
            />
          </div>

          <!-- Empty State -->
          <div
            v-else
            class="text-center py-16"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="w-16 h-16 mx-auto text-gray-600 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <h3 class="text-xl font-semibold text-white mb-2">لا توجد نتائج</h3>
            <p class="text-gray-400 mb-6">لم نجد أدوات تطابق بحثك</p>
            <RouterLink
              to="/"
              class="inline-block px-6 py-3 bg-neon-blue text-white rounded-lg hover:bg-neon-blue/90 transition-colors"
            >
              استعرض الفئات
            </RouterLink>
          </div>

          <!-- Pagination -->
          <div v-if="totalPages > 1" class="mt-8">
            <Pagination
              :current-page="currentPage"
              :total-pages="totalPages"
              @page-change="handlePageChange"
            />
          </div>
        </div>
      </div>
    </div>

    <!-- Compare Bar -->
    <CompareBar
      v-if="comparisonStore.count > 0"
      :tools="comparisonStore.selectedTools"
      @remove="(tool) => comparisonStore.removeTool(tool.id)"
      @clear="clearCompare"
      @compare="goToCompare"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { apiClient } from '@/lib/apiClient'
import { useComparisonStore } from '@/stores/comparison'
import FiltersPanel, { type FilterValues } from '@/components/search/FiltersPanel.vue'
import ToolCard from '@/components/tools/ToolCard.vue'
import SkeletonCard from '@/components/common/SkeletonCard.vue'
import Pagination from '@/components/common/Pagination.vue'
import CompareBar from '@/components/common/CompareBar.vue'
import type { Tool } from '@/lib/types'

interface ToolsResponse {
  data: Tool[]
  meta: {
    page: number
    page_size: number
    total: number
  }
}

const route = useRoute()
const router = useRouter()
const comparisonStore = useComparisonStore()

const tools = ref<Tool[]>([])
const loading = ref(true)
const totalResults = ref(0)
const currentPage = ref(1)
const pageSize = 20
const showMobileFilters = ref(false)

const filters = reactive<FilterValues>({
  category: '',
  price: '',
  minRating: 0,
  platforms: [],
  sort: 'top_rated'
})

const totalPages = computed(() => Math.ceil(totalResults.value / pageSize))

const pageTitle = computed(() => {
  if (route.params.slug) {
    return `أدوات: ${route.params.slug}`
  }
  if (route.query.q) {
    return `نتائج البحث: ${route.query.q}`
  }
  return 'جميع الأدوات'
})

function isInCompare(toolId: number): boolean {
  return comparisonStore.isInComparison(toolId)
}

function addToCompare(tool: Tool) {
  comparisonStore.addTool(tool)
}

function removeFromCompare(tool: Tool) {
  comparisonStore.removeTool(tool.id)
}

function clearCompare() {
  comparisonStore.clear()
}

function goToCompare() {
  router.push(comparisonStore.getCompareUrl())
}

function syncFiltersFromRoute() {
  filters.category = (route.params.slug as string) || (route.query.category as string) || ''
  filters.price = (route.query.price as string) || ''
  filters.minRating = route.query.min_rating ? parseFloat(route.query.min_rating as string) : 0
  filters.platforms = route.query.platform ? (route.query.platform as string).split(',') : []
  filters.sort = (route.query.sort as string) || 'top_rated'
  currentPage.value = route.query.page ? parseInt(route.query.page as string) : 1
}

function updateRouteQuery() {
  const query: Record<string, string> = {}

  if (route.query.q) {
    query.q = route.query.q as string
  }
  if (filters.category && !route.params.slug) {
    query.category = filters.category
  }
  if (filters.price) {
    query.price = filters.price
  }
  if (filters.minRating > 0) {
    query.min_rating = filters.minRating.toString()
  }
  if (filters.platforms.length > 0) {
    query.platform = filters.platforms.join(',')
  }
  if (filters.sort !== 'top_rated') {
    query.sort = filters.sort
  }
  if (currentPage.value > 1) {
    query.page = currentPage.value.toString()
  }

  router.replace({ query })
}

async function fetchTools() {
  loading.value = true

  try {
    const params: Record<string, any> = {
      page: currentPage.value,
      page_size: pageSize,
      sort: filters.sort
    }

    if (filters.category) {
      params.category = filters.category
    }
    if (filters.price) {
      params.price = filters.price
    }
    if (filters.minRating > 0) {
      params.min_rating = filters.minRating
    }
    if (filters.platforms.length > 0) {
      params.platform = filters.platforms[0] // API supports single platform
    }

    let response: ToolsResponse

    if (route.query.q) {
      // Search endpoint
      params.q = route.query.q
      response = await apiClient.get<ToolsResponse>('/search/tools', params)
    } else if (route.params.slug) {
      // Category tools endpoint
      response = await apiClient.get<ToolsResponse>(`/categories/${route.params.slug}/tools`, params)
    } else {
      // All tools endpoint
      response = await apiClient.get<ToolsResponse>('/tools', params)
    }

    tools.value = response.data
    totalResults.value = response.meta.total
  } catch (e) {
    console.error('Failed to fetch tools:', e)
    tools.value = []
    totalResults.value = 0
  } finally {
    loading.value = false
  }
}

function handleFilterChange(newFilters: FilterValues) {
  Object.assign(filters, newFilters)
  currentPage.value = 1 // Reset to first page
  updateRouteQuery()
  fetchTools()
  showMobileFilters.value = false
}

function handlePageChange(page: number) {
  currentPage.value = page
  updateRouteQuery()
  fetchTools()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

// Watch route changes
watch(
  () => [route.query, route.params],
  () => {
    syncFiltersFromRoute()
    fetchTools()
  },
  { deep: true }
)

onMounted(() => {
  syncFiltersFromRoute()
  fetchTools()
})
</script>

<style scoped>
.search-results-view {
  background-color: #05060A;
}
</style>
