<template>
  <aside class="filters-panel bg-dark-surface border border-gray-800 rounded-xl p-4" dir="rtl">
    <div class="flex items-center justify-between mb-6">
      <h3 class="text-lg font-semibold text-white">تصفية النتائج</h3>
      <button
        v-if="hasActiveFilters"
        @click="clearFilters"
        class="text-sm text-neon-blue hover:underline"
      >
        مسح الفلاتر
      </button>
    </div>

    <!-- Categories Filter -->
    <div class="filter-group mb-6" v-if="categories.length > 0">
      <h4 class="filter-label">الفئة</h4>
      <select
        v-model="localFilters.category"
        @change="emitChange"
        class="filter-select"
      >
        <option value="">جميع الفئات</option>
        <option v-for="cat in categories" :key="cat.slug" :value="cat.slug">
          {{ cat.name }}
        </option>
      </select>
    </div>

    <!-- Price Filter -->
    <div class="filter-group mb-6">
      <h4 class="filter-label">السعر</h4>
      <div class="space-y-2">
        <label class="filter-radio">
          <input
            type="radio"
            v-model="localFilters.price"
            value=""
            @change="emitChange"
          />
          <span>الكل</span>
        </label>
        <label class="filter-radio">
          <input
            type="radio"
            v-model="localFilters.price"
            value="free"
            @change="emitChange"
          />
          <span>مجاني</span>
        </label>
        <label class="filter-radio">
          <input
            type="radio"
            v-model="localFilters.price"
            value="freemium"
            @change="emitChange"
          />
          <span>Freemium</span>
        </label>
        <label class="filter-radio">
          <input
            type="radio"
            v-model="localFilters.price"
            value="paid"
            @change="emitChange"
          />
          <span>مدفوع</span>
        </label>
      </div>
    </div>

    <!-- Rating Filter -->
    <div class="filter-group mb-6">
      <h4 class="filter-label">التقييم</h4>
      <div class="space-y-2">
        <label class="filter-radio">
          <input
            type="radio"
            v-model="localFilters.minRating"
            :value="0"
            @change="emitChange"
          />
          <span>الكل</span>
        </label>
        <label class="filter-radio">
          <input
            type="radio"
            v-model="localFilters.minRating"
            :value="4"
            @change="emitChange"
          />
          <span>4+ نجوم</span>
        </label>
        <label class="filter-radio">
          <input
            type="radio"
            v-model="localFilters.minRating"
            :value="4.5"
            @change="emitChange"
          />
          <span>4.5+ نجوم</span>
        </label>
      </div>
    </div>

    <!-- Platform Filter -->
    <div class="filter-group">
      <h4 class="filter-label">المنصة</h4>
      <div class="space-y-2">
        <label class="filter-checkbox">
          <input
            type="checkbox"
            v-model="localFilters.platforms"
            value="web"
            @change="emitChange"
          />
          <span>ويب</span>
        </label>
        <label class="filter-checkbox">
          <input
            type="checkbox"
            v-model="localFilters.platforms"
            value="mobile"
            @change="emitChange"
          />
          <span>موبايل</span>
        </label>
        <label class="filter-checkbox">
          <input
            type="checkbox"
            v-model="localFilters.platforms"
            value="api"
            @change="emitChange"
          />
          <span>API</span>
        </label>
        <label class="filter-checkbox">
          <input
            type="checkbox"
            v-model="localFilters.platforms"
            value="desktop"
            @change="emitChange"
          />
          <span>سطح المكتب</span>
        </label>
      </div>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { reactive, computed, watch, onMounted, ref } from 'vue'
import { apiClient } from '@/lib/apiClient'
import type { Category, ApiResponse } from '@/lib/types'

export interface FilterValues {
  category: string
  price: string
  minRating: number
  platforms: string[]
  sort: string
}

const props = defineProps<{
  filters: FilterValues
}>()

const emit = defineEmits<{
  (e: 'update:filters', filters: FilterValues): void
}>()

const categories = ref<Category[]>([])

const localFilters = reactive<FilterValues>({
  category: props.filters.category || '',
  price: props.filters.price || '',
  minRating: props.filters.minRating || 0,
  platforms: props.filters.platforms || [],
  sort: props.filters.sort || 'top_rated'
})

// Sync with props
watch(() => props.filters, (newFilters) => {
  localFilters.category = newFilters.category || ''
  localFilters.price = newFilters.price || ''
  localFilters.minRating = newFilters.minRating || 0
  localFilters.platforms = newFilters.platforms || []
  localFilters.sort = newFilters.sort || 'top_rated'
}, { deep: true })

const hasActiveFilters = computed(() => {
  return localFilters.category !== '' ||
         localFilters.price !== '' ||
         localFilters.minRating > 0 ||
         localFilters.platforms.length > 0
})

function emitChange() {
  emit('update:filters', { ...localFilters })
}

function clearFilters() {
  localFilters.category = ''
  localFilters.price = ''
  localFilters.minRating = 0
  localFilters.platforms = []
  emitChange()
}

async function loadCategories() {
  try {
    const response = await apiClient.get<ApiResponse<Category[]>>('/categories')
    categories.value = response.data
  } catch (e) {
    console.error('Failed to load categories:', e)
  }
}

onMounted(() => {
  loadCategories()
})
</script>

<style scoped>
.filter-label {
  @apply text-sm font-medium text-gray-300 mb-3 block;
}

.filter-select {
  @apply w-full px-3 py-2 bg-dark-bg border border-gray-700 rounded-lg
         text-white text-sm focus:border-neon-blue focus:outline-none;
}

.filter-radio {
  @apply flex items-center gap-2 text-gray-400 cursor-pointer hover:text-white;
}

.filter-radio input {
  @apply w-4 h-4 accent-neon-blue;
}

.filter-checkbox {
  @apply flex items-center gap-2 text-gray-400 cursor-pointer hover:text-white;
}

.filter-checkbox input {
  @apply w-4 h-4 accent-neon-blue rounded;
}
</style>
