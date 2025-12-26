<template>
  <section class="category-grid py-12">
    <div class="container mx-auto px-4">
      <h2 class="text-2xl font-bold text-white mb-8 text-right">تصفح حسب الفئة</h2>

      <!-- Loading Skeleton -->
      <div v-if="loading" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4" dir="rtl">
        <div v-for="i in 8" :key="i" class="category-skeleton">
          <SkeletonCard height="120px" rounded />
        </div>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="text-center py-8">
        <p class="text-red-400">{{ error }}</p>
        <button @click="loadCategories" class="mt-4 text-neon-blue hover:underline">
          إعادة المحاولة
        </button>
      </div>

      <!-- Categories Grid -->
      <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4" dir="rtl">
        <RouterLink
          v-for="category in categories"
          :key="category.id"
          :to="`/categories/${category.slug}`"
          class="category-card group"
        >
          <div class="flex items-center gap-4">
            <div class="category-icon">
              <img
                v-if="category.icon_url"
                :src="category.icon_url"
                :alt="category.name"
                class="w-8 h-8 object-contain"
              />
              <div v-else class="w-8 h-8 bg-neon-blue/20 rounded-lg flex items-center justify-center">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-neon-blue" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
                </svg>
              </div>
            </div>
            <div class="flex-1 text-right">
              <h3 class="font-semibold text-white group-hover:text-neon-blue transition-colors">
                {{ category.name }}
              </h3>
              <p v-if="category.description" class="text-sm text-gray-400 line-clamp-1">
                {{ category.description }}
              </p>
            </div>
          </div>
          <div class="mt-3 flex justify-end">
            <span class="text-xs text-gray-500 group-hover:text-gray-400">
              عرض الأدوات
              <svg xmlns="http://www.w3.org/2000/svg" class="inline w-3 h-3 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
              </svg>
            </span>
          </div>
        </RouterLink>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { apiClient } from '@/lib/apiClient'
import type { Category, ApiResponse } from '@/lib/types'
import SkeletonCard from '@/components/common/SkeletonCard.vue'

const categories = ref<Category[]>([])
const loading = ref(true)
const error = ref<string | null>(null)

async function loadCategories() {
  loading.value = true
  error.value = null

  try {
    const response = await apiClient.get<ApiResponse<Category[]>>('/categories')
    categories.value = response.data
  } catch (e) {
    error.value = 'فشل في تحميل الفئات'
    console.error('Failed to load categories:', e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadCategories()
})
</script>

<style scoped>
.category-card {
  @apply block p-4 bg-dark-surface border border-gray-800 rounded-xl
         hover:border-neon-blue/50 hover:bg-dark-surface/80
         transition-all duration-200;
}

.category-skeleton {
  @apply bg-dark-surface rounded-xl;
}
</style>
