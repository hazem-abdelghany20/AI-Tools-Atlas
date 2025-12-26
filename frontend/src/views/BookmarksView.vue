<template>
  <div class="bookmarks-view max-w-5xl mx-auto px-4 py-8" dir="rtl">
    <!-- Header -->
    <div class="mb-8">
      <h1 class="text-3xl font-bold text-white mb-2">المفضلة</h1>
      <p class="text-gray-400">الأدوات التي قمت بحفظها</p>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <SkeletonCard v-for="n in 4" :key="n" />
    </div>

    <!-- Empty State -->
    <div v-else-if="bookmarkedTools.length === 0" class="text-center py-16">
      <svg class="w-16 h-16 mx-auto text-gray-600 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
      </svg>
      <h2 class="text-xl font-semibold text-white mb-2">لا توجد أدوات محفوظة</h2>
      <p class="text-gray-400 mb-6">ابدأ بحفظ الأدوات التي تهمك للرجوع إليها لاحقاً</p>
      <router-link
        to="/"
        class="px-6 py-3 bg-neon-blue text-dark-bg font-medium rounded-lg hover:bg-neon-blue-hover transition-colors inline-block"
      >
        استكشف الأدوات
      </router-link>
    </div>

    <!-- Bookmarked Tools Grid -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <ToolCard
        v-for="tool in bookmarkedTools"
        :key="tool.id"
        :tool="tool"
        :show-compare="true"
        :is-in-compare="isInCompare(tool.id)"
        @add-to-compare="handleAddToCompare"
        @remove-from-compare="handleRemoveFromCompare"
      />
    </div>

    <!-- Compare Bar (if tools selected) -->
    <CompareBar
      v-if="compareTools.length > 0"
      :tools="compareTools"
      @remove="handleRemoveFromCompare"
      @clear="clearCompare"
      @compare="goToCompare"
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useBookmarksStore } from '../stores/bookmarks'
import { storeToRefs } from 'pinia'
import ToolCard from '../components/tools/ToolCard.vue'
import SkeletonCard from '../components/common/SkeletonCard.vue'
import CompareBar from '../components/common/CompareBar.vue'
import type { Tool } from '../lib/types'

const router = useRouter()
const bookmarksStore = useBookmarksStore()
const { bookmarkedTools, loading } = storeToRefs(bookmarksStore)

const compareTools = ref<Tool[]>([])
const maxCompare = 4

function isInCompare(toolId: number): boolean {
  return compareTools.value.some(t => t.id === toolId)
}

function handleAddToCompare(tool: Tool) {
  if (compareTools.value.length < maxCompare && !isInCompare(tool.id)) {
    compareTools.value.push(tool)
  }
}

function handleRemoveFromCompare(tool: Tool) {
  compareTools.value = compareTools.value.filter(t => t.id !== tool.id)
}

function clearCompare() {
  compareTools.value = []
}

function goToCompare() {
  const slugs = compareTools.value.map(t => t.slug).join(',')
  router.push(`/compare?tools=${slugs}`)
}

onMounted(() => {
  bookmarksStore.fetchBookmarks()
})
</script>
