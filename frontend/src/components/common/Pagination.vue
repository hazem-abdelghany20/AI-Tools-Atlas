<template>
  <nav class="pagination flex items-center justify-center gap-2" dir="ltr">
    <!-- Previous Button -->
    <button
      @click="goToPage(currentPage - 1)"
      :disabled="currentPage <= 1"
      class="pagination-btn"
      :class="{ 'opacity-50 cursor-not-allowed': currentPage <= 1 }"
    >
      <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
      </svg>
    </button>

    <!-- Page Numbers -->
    <template v-for="page in displayedPages" :key="page">
      <span v-if="page === '...'" class="px-2 text-gray-500">...</span>
      <button
        v-else
        @click="goToPage(page as number)"
        class="pagination-btn"
        :class="{ 'bg-neon-blue text-white': currentPage === page }"
      >
        {{ page }}
      </button>
    </template>

    <!-- Next Button -->
    <button
      @click="goToPage(currentPage + 1)"
      :disabled="currentPage >= totalPages"
      class="pagination-btn"
      :class="{ 'opacity-50 cursor-not-allowed': currentPage >= totalPages }"
    >
      <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
      </svg>
    </button>
  </nav>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  currentPage: number
  totalPages: number
  maxVisiblePages?: number
}>()

const emit = defineEmits<{
  (e: 'page-change', page: number): void
}>()

const maxVisible = props.maxVisiblePages || 5

const displayedPages = computed(() => {
  const pages: (number | string)[] = []
  const total = props.totalPages
  const current = props.currentPage

  if (total <= maxVisible) {
    for (let i = 1; i <= total; i++) {
      pages.push(i)
    }
  } else {
    // Always show first page
    pages.push(1)

    const startPage = Math.max(2, current - 1)
    const endPage = Math.min(total - 1, current + 1)

    if (startPage > 2) {
      pages.push('...')
    }

    for (let i = startPage; i <= endPage; i++) {
      pages.push(i)
    }

    if (endPage < total - 1) {
      pages.push('...')
    }

    // Always show last page
    if (total > 1) {
      pages.push(total)
    }
  }

  return pages
})

function goToPage(page: number) {
  if (page >= 1 && page <= props.totalPages && page !== props.currentPage) {
    emit('page-change', page)
  }
}
</script>

<style scoped>
.pagination-btn {
  @apply px-3 py-2 text-sm font-medium text-gray-300 bg-dark-surface
         border border-gray-700 rounded-lg
         hover:border-neon-blue hover:text-neon-blue
         transition-colors disabled:hover:border-gray-700 disabled:hover:text-gray-300;
}
</style>
