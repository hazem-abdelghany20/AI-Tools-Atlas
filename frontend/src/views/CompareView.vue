<template>
  <div class="compare-view max-w-7xl mx-auto px-4 py-8" dir="rtl">
    <!-- Header -->
    <div class="mb-8 flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-white mb-2">المقارنة</h1>
        <p class="text-gray-400">قارن بين الأدوات جنباً إلى جنب</p>
      </div>

      <!-- Share Button -->
      <button
        v-if="tools.length > 0"
        @click="shareComparison"
        class="flex items-center gap-2 px-4 py-2 border border-gray-700 text-gray-300 rounded-lg
               hover:border-neon-blue hover:text-neon-blue transition-colors"
      >
        <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.368 2.684 3 3 0 00-5.368-2.684z" />
        </svg>
        <span>مشاركة المقارنة</span>
      </button>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex justify-center py-16">
      <Spinner />
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="text-center py-16">
      <svg class="w-16 h-16 mx-auto text-red-500 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
              d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
      </svg>
      <h2 class="text-xl font-semibold text-white mb-2">حدث خطأ</h2>
      <p class="text-gray-400 mb-6">{{ error }}</p>
      <button
        @click="fetchTools"
        class="px-6 py-3 bg-neon-blue text-dark-bg font-medium rounded-lg hover:bg-neon-blue-hover transition-colors"
      >
        إعادة المحاولة
      </button>
    </div>

    <!-- Empty State -->
    <div v-else-if="tools.length === 0" class="text-center py-16">
      <svg class="w-16 h-16 mx-auto text-gray-600 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
              d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
      </svg>
      <h2 class="text-xl font-semibold text-white mb-2">لا توجد أدوات للمقارنة</h2>
      <p class="text-gray-400 mb-6">أضف أدوات من نتائج البحث أو المفضلة للمقارنة بينها</p>
      <div class="flex justify-center gap-4">
        <router-link
          to="/"
          class="px-6 py-3 bg-neon-blue text-dark-bg font-medium rounded-lg hover:bg-neon-blue-hover transition-colors"
        >
          استكشف الأدوات
        </router-link>
        <router-link
          to="/bookmarks"
          class="px-6 py-3 border border-gray-700 text-gray-300 rounded-lg hover:border-neon-blue hover:text-neon-blue transition-colors"
        >
          المفضلة
        </router-link>
      </div>
    </div>

    <!-- Comparison Table -->
    <div v-else>
      <!-- Tool Count -->
      <div class="mb-6 flex items-center justify-between">
        <span class="text-gray-400">{{ tools.length }} أدوات في المقارنة</span>
        <button
          @click="clearAll"
          class="text-red-400 hover:text-red-300 text-sm transition-colors"
        >
          مسح الكل
        </button>
      </div>

      <!-- Comparison Table -->
      <div class="bg-dark-surface border border-gray-800 rounded-xl overflow-hidden">
        <CompareTable
          :tools="tools"
          @remove="handleRemoveTool"
        />
      </div>

      <!-- Add More Tools CTA -->
      <div v-if="tools.length < 4" class="mt-6 text-center">
        <p class="text-gray-400 mb-3">يمكنك إضافة حتى {{ 4 - tools.length }} أدوات أخرى</p>
        <router-link
          to="/"
          class="inline-flex items-center gap-2 px-4 py-2 border border-gray-700 text-gray-300 rounded-lg
                 hover:border-neon-blue hover:text-neon-blue transition-colors"
        >
          <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          <span>إضافة أدوات</span>
        </router-link>
      </div>
    </div>

    <!-- Toast for share success -->
    <Transition name="fade">
      <div
        v-if="showShareToast"
        class="fixed bottom-6 left-1/2 -translate-x-1/2 px-6 py-3 bg-green-500 text-white rounded-lg shadow-lg z-50"
      >
        تم نسخ رابط المقارنة
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { apiClient } from '../lib/apiClient'
import CompareTable from '../components/compare/CompareTable.vue'
import Spinner from '../components/common/Spinner.vue'
import type { Tool, ApiResponse } from '../lib/types'

const route = useRoute()
const router = useRouter()

const tools = ref<Tool[]>([])
const loading = ref(false)
const error = ref<string | null>(null)
const showShareToast = ref(false)

async function fetchTools() {
  const toolsParam = route.query.tools as string
  if (!toolsParam) {
    tools.value = []
    return
  }

  const slugs = toolsParam.split(',').map(s => s.trim()).filter(s => s.length > 0)
  if (slugs.length === 0) {
    tools.value = []
    return
  }

  loading.value = true
  error.value = null

  try {
    // Fetch all tools by slug in parallel
    const responses = await Promise.all(
      slugs.map(slug => apiClient.get<ApiResponse<Tool>>(`/tools/${slug}`).catch(() => null))
    )

    tools.value = responses
      .filter((res): res is ApiResponse<Tool> => res !== null)
      .map(res => res.data)
  } catch (err) {
    error.value = 'فشل في تحميل بيانات الأدوات'
    console.error('Failed to fetch tools for comparison:', err)
  } finally {
    loading.value = false
  }
}

function handleRemoveTool(tool: Tool) {
  tools.value = tools.value.filter(t => t.id !== tool.id)
  updateUrl()
}

function clearAll() {
  tools.value = []
  router.push('/compare')
}

function updateUrl() {
  if (tools.value.length === 0) {
    router.push('/compare')
  } else {
    const slugs = tools.value.map(t => t.slug).join(',')
    router.push(`/compare?tools=${slugs}`)
  }
}

async function shareComparison() {
  const url = window.location.href
  try {
    await navigator.clipboard.writeText(url)
    showShareToast.value = true
    setTimeout(() => {
      showShareToast.value = false
    }, 3000)
  } catch (err) {
    console.error('Failed to copy URL:', err)
  }
}

// Watch for route changes
watch(() => route.query.tools, () => {
  fetchTools()
})

onMounted(() => {
  fetchTools()
})
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
