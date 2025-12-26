<template>
  <section class="hero-section py-16 md:py-24">
    <div class="container mx-auto px-4 text-center">
      <h1 class="text-3xl md:text-5xl font-bold text-white mb-6">
        اكتشف أدوات الذكاء الاصطناعي المناسبة لك
      </h1>
      <p class="text-gray-400 text-lg md:text-xl mb-8 max-w-2xl mx-auto">
        ابحث عن الأدوات التي تناسب احتياجاتك من بين آلاف الأدوات المتاحة
      </p>

      <!-- Search Form -->
      <form @submit.prevent="handleSearch" class="max-w-2xl mx-auto">
        <div class="relative">
          <input
            v-model="searchQuery"
            type="text"
            dir="rtl"
            class="w-full px-6 py-4 pr-14 text-lg bg-dark-surface border-2 border-gray-700
                   rounded-xl text-white placeholder-gray-500
                   focus:border-neon-blue focus:outline-none focus:ring-2 focus:ring-neon-blue/20
                   transition-all duration-200"
            :placeholder="placeholder"
          />
          <button
            type="submit"
            class="absolute left-2 top-1/2 -translate-y-1/2 p-3 bg-neon-blue
                   text-white rounded-lg hover:bg-neon-blue/90 transition-colors"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
          </button>
        </div>
      </form>

      <!-- Suggestion Chips -->
      <div class="flex flex-wrap justify-center gap-3 mt-6">
        <button
          v-for="suggestion in suggestions"
          :key="suggestion"
          @click="selectSuggestion(suggestion)"
          class="px-4 py-2 text-sm bg-dark-surface border border-gray-700
                 rounded-full text-gray-300 hover:border-neon-blue hover:text-neon-blue
                 transition-colors"
        >
          {{ suggestion }}
        </button>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useFiltersStore } from '@/stores/filters'

const router = useRouter()
const filtersStore = useFiltersStore()

const searchQuery = ref('')
const placeholder = 'ابحث باسم الأداة...'

const suggestions: string[] = []

function handleSearch() {
  if (searchQuery.value.trim()) {
    filtersStore.setFilter('query', searchQuery.value.trim())
    router.push({ path: '/search', query: { q: searchQuery.value.trim() } })
  }
}

function selectSuggestion(suggestion: string) {
  searchQuery.value = suggestion
  handleSearch()
}
</script>

<style scoped>
.hero-section {
  background: linear-gradient(180deg, #05060A 0%, #0A0B10 100%);
}
</style>
