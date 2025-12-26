<template>
  <div class="home-view min-h-screen bg-dark-bg">
    <!-- Hero Section with Search -->
    <HeroSearch />

    <!-- Category Grid -->
    <CategoryGrid />

    <!-- Popular Tools Section -->
    <section v-if="popularTools.length > 0" class="popular-section py-12">
      <div class="container mx-auto px-4">
        <h2 class="text-2xl font-bold text-white mb-8 text-right">الأدوات الأكثر شعبية</h2>
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4" dir="rtl">
          <ToolCard
            v-for="tool in popularTools"
            :key="tool.id"
            :tool="tool"
            :is-in-compare="comparisonStore.isInComparison(tool.id)"
            @add-to-compare="(t) => comparisonStore.addTool(t)"
            @remove-from-compare="(t) => comparisonStore.removeTool(t.id)"
          />
        </div>
      </div>
    </section>

    <!-- Compare Bar -->
    <CompareBar
      v-if="comparisonStore.count > 0"
      :tools="comparisonStore.selectedTools"
      @remove="(tool) => comparisonStore.removeTool(tool.id)"
      @clear="() => comparisonStore.clear()"
      @compare="goToCompare"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import HeroSearch from '@/components/home/HeroSearch.vue'
import CategoryGrid from '@/components/home/CategoryGrid.vue'
import ToolCard from '@/components/tools/ToolCard.vue'
import CompareBar from '@/components/common/CompareBar.vue'
import { apiClient } from '@/lib/apiClient'
import { useComparisonStore } from '@/stores/comparison'
import type { Tool } from '@/lib/types'

interface ToolsResponse {
  data: Tool[]
  meta: {
    page: number
    page_size: number
    total: number
  }
}

const router = useRouter()
const comparisonStore = useComparisonStore()
const popularTools = ref<Tool[]>([])

function goToCompare() {
  router.push(comparisonStore.getCompareUrl())
}

async function loadPopularTools() {
  try {
    const response = await apiClient.get<ToolsResponse>('/tools', {
      sort: 'top_rated',
      page_size: '4'
    })
    popularTools.value = response.data
  } catch (e) {
    console.error('Failed to load popular tools:', e)
  }
}

onMounted(() => {
  loadPopularTools()
})
</script>

<style scoped>
.home-view {
  background-color: #05060A;
}

.popular-section {
  background: linear-gradient(180deg, #0A0B10 0%, #05060A 100%);
}
</style>
