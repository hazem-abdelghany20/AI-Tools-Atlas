<template>
  <div class="tool-alternatives space-y-8">
    <!-- Similar Tools -->
    <section v-if="similar.length > 0" class="bg-dark-surface rounded-lg p-6 md:p-8">
      <h2 class="text-xl font-bold mb-6">أدوات مشابهة</h2>
      <div class="horizontal-scroll-container relative">
        <div
          class="flex gap-4 overflow-x-auto pb-4 scroll-smooth snap-x snap-mandatory"
          :class="{ 'opacity-50': loading }"
        >
          <div
            v-for="tool in similar"
            :key="tool.id"
            class="flex-shrink-0 w-72 snap-start"
          >
            <ToolCard :tool="tool" compact />
          </div>
        </div>
        <!-- Gradient fade indicators -->
        <div class="absolute top-0 start-0 bottom-4 w-8 bg-gradient-to-l from-transparent to-dark-surface pointer-events-none"></div>
        <div class="absolute top-0 end-0 bottom-4 w-8 bg-gradient-to-r from-transparent to-dark-surface pointer-events-none"></div>
      </div>
    </section>

    <!-- Alternative Tools -->
    <section v-if="alternatives.length > 0" class="bg-dark-surface rounded-lg p-6 md:p-8">
      <h2 class="text-xl font-bold mb-6">بدائل لهذه الأداة</h2>
      <div class="horizontal-scroll-container relative">
        <div
          class="flex gap-4 overflow-x-auto pb-4 scroll-smooth snap-x snap-mandatory"
          :class="{ 'opacity-50': loading }"
        >
          <div
            v-for="tool in alternatives"
            :key="tool.id"
            class="flex-shrink-0 w-72 snap-start"
          >
            <ToolCard :tool="tool" compact />
          </div>
        </div>
        <!-- Gradient fade indicators -->
        <div class="absolute top-0 start-0 bottom-4 w-8 bg-gradient-to-l from-transparent to-dark-surface pointer-events-none"></div>
        <div class="absolute top-0 end-0 bottom-4 w-8 bg-gradient-to-r from-transparent to-dark-surface pointer-events-none"></div>
      </div>
    </section>

    <!-- Loading State -->
    <div v-if="loading && !similar.length && !alternatives.length" class="space-y-8">
      <section class="bg-dark-surface rounded-lg p-6 md:p-8">
        <div class="h-6 w-32 bg-dark-border rounded skeleton mb-6"></div>
        <div class="flex gap-4 overflow-hidden">
          <div v-for="i in 4" :key="i" class="flex-shrink-0 w-72 h-48 bg-dark-border rounded-lg skeleton"></div>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { apiClient } from '../../lib/apiClient'
import type { Tool } from '../../lib/types'
import ToolCard from './ToolCard.vue'

interface AlternativesResponse {
  similar: Tool[]
  alternatives: Tool[]
}

const props = defineProps<{
  toolSlug: string
}>()

const similar = ref<Tool[]>([])
const alternatives = ref<Tool[]>([])
const loading = ref(true)
const error = ref<string | null>(null)

const fetchAlternatives = async () => {
  if (!props.toolSlug) return

  loading.value = true
  error.value = null

  try {
    const response = await apiClient.get<{ data: AlternativesResponse }>(`/tools/${props.toolSlug}/alternatives`)
    similar.value = response.data.similar || []
    alternatives.value = response.data.alternatives || []
  } catch (err: any) {
    error.value = err.message
    // Don't show error to user, just hide the section
  } finally {
    loading.value = false
  }
}

watch(() => props.toolSlug, () => {
  fetchAlternatives()
})

onMounted(() => {
  fetchAlternatives()
})
</script>

<style scoped>
.skeleton {
  background: linear-gradient(90deg, #0A0B10 25%, #1F2933 50%, #0A0B10 75%);
  background-size: 200% 100%;
  animation: pulse 1.5s ease-in-out infinite;
}

@keyframes pulse {
  0% {
    background-position: 200% 0;
  }
  100% {
    background-position: -200% 0;
  }
}

/* Custom scrollbar styling */
.horizontal-scroll-container div::-webkit-scrollbar {
  height: 6px;
}

.horizontal-scroll-container div::-webkit-scrollbar-track {
  background: #1F2933;
  border-radius: 3px;
}

.horizontal-scroll-container div::-webkit-scrollbar-thumb {
  background: #3B82F6;
  border-radius: 3px;
}

.horizontal-scroll-container div::-webkit-scrollbar-thumb:hover {
  background: #2563EB;
}
</style>
