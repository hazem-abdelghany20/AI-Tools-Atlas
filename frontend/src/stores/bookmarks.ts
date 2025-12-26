import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { apiClient } from '../lib/apiClient'
import type { Tool, ApiResponse } from '../lib/types'

export const useBookmarksStore = defineStore('bookmarks', () => {
  const bookmarkedToolIds = ref<number[]>([])
  const bookmarkedTools = ref<Tool[]>([])
  const loading = ref(false)

  const isBookmarked = computed(() => {
    return (toolId: number) => bookmarkedToolIds.value.includes(toolId)
  })

  async function fetchBookmarks() {
    loading.value = true
    try {
      const response = await apiClient.get<ApiResponse<Tool[]>>('/me/bookmarks')
      bookmarkedTools.value = response.data
      bookmarkedToolIds.value = response.data.map((t: Tool) => t.id)
    } catch (error) {
      console.error('Failed to fetch bookmarks:', error)
    } finally {
      loading.value = false
    }
  }

  async function addBookmark(toolId: number) {
    try {
      await apiClient.post('/me/bookmarks', { tool_id: toolId })
      if (!bookmarkedToolIds.value.includes(toolId)) {
        bookmarkedToolIds.value.push(toolId)
      }
    } catch (error) {
      console.error('Failed to add bookmark:', error)
      throw error
    }
  }

  async function removeBookmark(toolId: number) {
    // Optimistic update - remove from UI immediately
    const previousIds = [...bookmarkedToolIds.value]
    const previousTools = [...bookmarkedTools.value]

    bookmarkedToolIds.value = bookmarkedToolIds.value.filter(id => id !== toolId)
    bookmarkedTools.value = bookmarkedTools.value.filter(t => t.id !== toolId)

    try {
      await apiClient.delete(`/me/bookmarks/${toolId}`)
    } catch (error) {
      // Rollback on error
      bookmarkedToolIds.value = previousIds
      bookmarkedTools.value = previousTools
      console.error('Failed to remove bookmark:', error)
      throw error
    }
  }

  return {
    bookmarkedToolIds,
    bookmarkedTools,
    loading,
    isBookmarked,
    fetchBookmarks,
    addBookmark,
    removeBookmark,
  }
}, {
  persist: true, // Persist to localStorage
})
