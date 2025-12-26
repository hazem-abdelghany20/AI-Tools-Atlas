import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Tool } from '../lib/types'

const MAX_COMPARE_TOOLS = 4

export const useComparisonStore = defineStore('comparison', () => {
  const selectedTools = ref<Tool[]>([])

  const count = computed(() => selectedTools.value.length)
  const isEmpty = computed(() => selectedTools.value.length === 0)
  const isFull = computed(() => selectedTools.value.length >= MAX_COMPARE_TOOLS)

  function isInComparison(toolId: number): boolean {
    return selectedTools.value.some(t => t.id === toolId)
  }

  function addTool(tool: Tool): boolean {
    if (selectedTools.value.length >= MAX_COMPARE_TOOLS) {
      return false
    }
    if (!isInComparison(tool.id)) {
      selectedTools.value.push(tool)
      return true
    }
    return false
  }

  function removeTool(toolId: number) {
    selectedTools.value = selectedTools.value.filter(t => t.id !== toolId)
  }

  function clear() {
    selectedTools.value = []
  }

  function getCompareUrl(): string {
    if (selectedTools.value.length === 0) return '/compare'
    const slugs = selectedTools.value.map(t => t.slug).join(',')
    return `/compare?tools=${slugs}`
  }

  return {
    selectedTools,
    count,
    isEmpty,
    isFull,
    isInComparison,
    addTool,
    removeTool,
    clear,
    getCompareUrl,
    MAX_COMPARE_TOOLS,
  }
}, {
  persist: true, // Persist to localStorage
})
