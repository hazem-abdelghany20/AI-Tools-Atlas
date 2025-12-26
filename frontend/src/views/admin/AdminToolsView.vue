<template>
  <div class="admin-tools-view min-h-screen bg-dark-bg" dir="rtl">
    <!-- Header -->
    <div class="bg-dark-surface border-b border-gray-800 px-6 py-4">
      <div class="max-w-7xl mx-auto flex justify-between items-center">
        <div>
          <h1 class="text-xl font-bold text-white">Admin - Tools</h1>
          <p class="text-gray-400 text-sm">Manage AI tools catalog</p>
        </div>
        <button
          @click="openCreateModal"
          class="px-4 py-2 bg-neon-blue text-dark-bg font-medium rounded-lg hover:bg-neon-blue-hover transition-colors flex items-center gap-2"
        >
          <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          Add Tool
        </button>
      </div>
    </div>

    <!-- Filters & Search -->
    <div class="max-w-7xl mx-auto px-6 py-4">
      <div class="flex flex-wrap gap-4 items-center">
        <!-- Search -->
        <div class="flex-1 min-w-[200px]">
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search by name or slug..."
            class="w-full px-4 py-2 bg-dark-surface border border-gray-700 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-neon-blue"
            @input="debouncedSearch"
          />
        </div>

        <!-- Status Filter -->
        <div class="flex gap-2">
          <button
            v-for="status in statusFilters"
            :key="status.value"
            @click="currentFilter = status.value"
            class="px-4 py-2 rounded-lg text-sm font-medium transition-colors"
            :class="currentFilter === status.value
              ? 'bg-neon-blue text-dark-bg'
              : 'bg-dark-surface text-gray-400 hover:text-white border border-gray-700'"
          >
            {{ status.label }}
          </button>
        </div>
      </div>
    </div>

    <!-- Table -->
    <div class="max-w-7xl mx-auto px-6 pb-6">
      <div class="bg-dark-surface border border-gray-800 rounded-xl overflow-hidden">
        <!-- Loading -->
        <div v-if="loading" class="p-8 text-center">
          <Spinner />
        </div>

        <!-- Empty State -->
        <div v-else-if="tools.length === 0" class="p-8 text-center text-gray-400">
          No tools found
        </div>

        <!-- Data Table -->
        <table v-else class="w-full">
          <thead class="bg-dark-bg/50">
            <tr>
              <th class="px-4 py-3 text-right text-xs font-medium text-gray-400 uppercase">Tool</th>
              <th class="px-4 py-3 text-right text-xs font-medium text-gray-400 uppercase">Category</th>
              <th class="px-4 py-3 text-right text-xs font-medium text-gray-400 uppercase">Rating</th>
              <th class="px-4 py-3 text-right text-xs font-medium text-gray-400 uppercase">Reviews</th>
              <th class="px-4 py-3 text-right text-xs font-medium text-gray-400 uppercase">Bookmarks</th>
              <th class="px-4 py-3 text-right text-xs font-medium text-gray-400 uppercase">Status</th>
              <th class="px-4 py-3 text-right text-xs font-medium text-gray-400 uppercase">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-800">
            <tr
              v-for="tool in tools"
              :key="tool.id"
              class="hover:bg-dark-bg/30 transition-colors"
              :class="{ 'opacity-60': tool.archived_at }"
            >
              <!-- Tool Info -->
              <td class="px-4 py-3">
                <div class="flex items-center gap-3">
                  <img
                    v-if="tool.logo_url"
                    :src="tool.logo_url"
                    :alt="tool.name"
                    class="w-10 h-10 rounded-lg object-cover"
                  />
                  <div
                    v-else
                    class="w-10 h-10 rounded-lg bg-neon-blue/20 flex items-center justify-center"
                  >
                    <span class="text-neon-blue font-bold">{{ tool.name.charAt(0) }}</span>
                  </div>
                  <div>
                    <div class="font-medium text-white">{{ tool.name }}</div>
                    <div class="text-sm text-gray-500">{{ tool.slug }}</div>
                  </div>
                </div>
              </td>

              <!-- Category -->
              <td class="px-4 py-3 text-gray-300">
                {{ tool.primary_category?.name || '-' }}
              </td>

              <!-- Rating -->
              <td class="px-4 py-3">
                <div class="flex items-center gap-1">
                  <svg class="w-4 h-4 text-yellow-400" fill="currentColor" viewBox="0 0 20 20">
                    <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"/>
                  </svg>
                  <span class="text-white">{{ tool.avg_rating_overall?.toFixed(1) || '0.0' }}</span>
                </div>
              </td>

              <!-- Reviews -->
              <td class="px-4 py-3 text-gray-300">
                {{ tool.review_count || 0 }}
              </td>

              <!-- Bookmarks -->
              <td class="px-4 py-3 text-gray-300">
                {{ tool.bookmark_count || 0 }}
              </td>

              <!-- Status -->
              <td class="px-4 py-3">
                <span
                  class="px-2 py-1 text-xs rounded-full"
                  :class="tool.archived_at
                    ? 'bg-red-500/20 text-red-400'
                    : 'bg-green-500/20 text-green-400'"
                >
                  {{ tool.archived_at ? 'Archived' : 'Active' }}
                </span>
              </td>

              <!-- Actions -->
              <td class="px-4 py-3">
                <div class="flex items-center gap-2">
                  <button
                    @click="openEditModal(tool)"
                    class="p-2 text-gray-400 hover:text-neon-blue transition-colors"
                    title="Edit"
                  >
                    <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                    </svg>
                  </button>
                  <button
                    @click="tool.archived_at ? restoreTool(tool.id) : archiveTool(tool.id)"
                    class="p-2 transition-colors"
                    :class="tool.archived_at ? 'text-green-400 hover:text-green-300' : 'text-gray-400 hover:text-red-400'"
                    :title="tool.archived_at ? 'Restore' : 'Archive'"
                  >
                    <svg v-if="tool.archived_at" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                    </svg>
                    <svg v-else class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4" />
                    </svg>
                  </button>
                  <router-link
                    :to="`/tools/${tool.slug}`"
                    target="_blank"
                    class="p-2 text-gray-400 hover:text-neon-blue transition-colors"
                    title="View on site"
                  >
                    <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
                    </svg>
                  </router-link>
                </div>
              </td>
            </tr>
          </tbody>
        </table>

        <!-- Pagination -->
        <div v-if="totalPages > 1" class="px-4 py-3 border-t border-gray-800 flex items-center justify-between">
          <div class="text-sm text-gray-400">
            Showing {{ (page - 1) * pageSize + 1 }} - {{ Math.min(page * pageSize, total) }} of {{ total }}
          </div>
          <div class="flex gap-2">
            <button
              @click="page--"
              :disabled="page === 1"
              class="px-3 py-1 rounded bg-dark-bg text-gray-400 disabled:opacity-50"
            >
              Previous
            </button>
            <button
              @click="page++"
              :disabled="page >= totalPages"
              class="px-3 py-1 rounded bg-dark-bg text-gray-400 disabled:opacity-50"
            >
              Next
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Tool Form Modal -->
    <ToolFormModal
      v-if="showModal"
      :tool="editingTool"
      :categories="categories"
      @close="closeModal"
      @saved="onToolSaved"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useSessionStore } from '@/stores/session'
import { storeToRefs } from 'pinia'
import { apiClient } from '@/lib/apiClient'
import type { Tool, Category } from '@/lib/types'
import Spinner from '@/components/common/Spinner.vue'
import ToolFormModal from '@/components/admin/ToolFormModal.vue'

const router = useRouter()
const sessionStore = useSessionStore()
const { user, isAuthenticated } = storeToRefs(sessionStore)

// Check admin access
watch([isAuthenticated, user], ([authenticated, currentUser]) => {
  if (!authenticated || currentUser?.role !== 'admin') {
    router.push('/')
  }
}, { immediate: true })

const tools = ref<Tool[]>([])
const categories = ref<Category[]>([])
const loading = ref(false)
const searchQuery = ref('')
const currentFilter = ref<'all' | 'active' | 'archived'>('all')
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)

const showModal = ref(false)
const editingTool = ref<Tool | null>(null)

const statusFilters: { value: 'all' | 'active' | 'archived'; label: string }[] = [
  { value: 'all', label: 'All' },
  { value: 'active', label: 'Active' },
  { value: 'archived', label: 'Archived' },
]

const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

let searchTimeout: ReturnType<typeof setTimeout> | null = null

function debouncedSearch() {
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    page.value = 1
    fetchTools()
  }, 300)
}

async function fetchTools() {
  loading.value = true
  try {
    const params = new URLSearchParams()
    params.set('page', page.value.toString())
    params.set('page_size', pageSize.value.toString())
    if (searchQuery.value) params.set('search', searchQuery.value)
    if (currentFilter.value === 'archived') params.set('archived', 'true')

    const response = await apiClient.get<{ data: Tool[]; meta: { total: number } }>(
      `/admin/tools?${params.toString()}`
    )

    let fetchedTools = response.data
    // Filter based on status
    if (currentFilter.value === 'active') {
      fetchedTools = fetchedTools.filter(t => !t.archived_at)
    } else if (currentFilter.value === 'archived') {
      fetchedTools = fetchedTools.filter(t => t.archived_at)
    }

    tools.value = fetchedTools
    total.value = response.meta.total
  } catch (error) {
    console.error('Failed to fetch tools:', error)
  } finally {
    loading.value = false
  }
}

async function fetchCategories() {
  try {
    const response = await apiClient.get<{ data: Category[] }>('/categories')
    categories.value = response.data
  } catch (error) {
    console.error('Failed to fetch categories:', error)
  }
}

function openCreateModal() {
  editingTool.value = null
  showModal.value = true
}

function openEditModal(tool: Tool) {
  editingTool.value = tool
  showModal.value = true
}

function closeModal() {
  showModal.value = false
  editingTool.value = null
}

function onToolSaved() {
  closeModal()
  fetchTools()
}

async function archiveTool(id: number) {
  if (!confirm('Are you sure you want to archive this tool?')) return
  try {
    await apiClient.delete(`/admin/tools/${id}`)
    fetchTools()
  } catch (error) {
    console.error('Failed to archive tool:', error)
  }
}

async function restoreTool(_id: number) {
  // For restore, we need to update the tool - future enhancement
  // For now just refresh
  alert('Restore functionality will be implemented with PATCH endpoint update')
}

watch([currentFilter, page], () => {
  fetchTools()
})

onMounted(() => {
  fetchTools()
  fetchCategories()
})
</script>
