import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useFiltersStore = defineStore('filters', () => {
  const query = ref('')
  const category = ref<string | null>(null)
  const price = ref<string | null>(null)
  const minRating = ref<number | null>(null)
  const platform = ref<string | null>(null)
  const sort = ref('top_rated')

  function setFilter(key: string, value: any) {
    switch (key) {
      case 'query':
        query.value = value
        break
      case 'category':
        category.value = value
        break
      case 'price':
        price.value = value
        break
      case 'minRating':
        minRating.value = value
        break
      case 'platform':
        platform.value = value
        break
      case 'sort':
        sort.value = value
        break
    }
  }

  function clearFilters() {
    query.value = ''
    category.value = null
    price.value = null
    minRating.value = null
    platform.value = null
    sort.value = 'top_rated'
  }

  return {
    query,
    category,
    price,
    minRating,
    platform,
    sort,
    setFilter,
    clearFilters,
  }
})
