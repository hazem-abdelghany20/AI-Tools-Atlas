<template>
  <div class="admin-tags-view min-h-screen bg-dark-bg" dir="rtl">
    <!-- Header -->
    <div class="bg-dark-surface border-b border-gray-800 px-6 py-4">
      <div class="max-w-7xl mx-auto flex justify-between items-center">
        <div>
          <h1 class="text-xl font-bold text-white">Admin - Tags</h1>
          <p class="text-gray-400 text-sm">Manage tool tags</p>
        </div>
        <button
          @click="openCreateModal"
          class="px-4 py-2 bg-neon-blue text-dark-bg font-medium rounded-lg hover:bg-neon-blue-hover transition-colors flex items-center gap-2"
        >
          <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          Add Tag
        </button>
      </div>
    </div>

    <!-- Table -->
    <div class="max-w-7xl mx-auto px-6 py-6">
      <div class="bg-dark-surface border border-gray-800 rounded-xl overflow-hidden">
        <!-- Loading -->
        <div v-if="loading" class="p-8 text-center">
          <Spinner />
        </div>

        <!-- Empty State -->
        <div v-else-if="tags.length === 0" class="p-8 text-center text-gray-400">
          No tags found
        </div>

        <!-- Data Table -->
        <table v-else class="w-full">
          <thead class="bg-dark-bg/50">
            <tr>
              <th class="px-4 py-3 text-right text-xs font-medium text-gray-400 uppercase">ID</th>
              <th class="px-4 py-3 text-right text-xs font-medium text-gray-400 uppercase">Name</th>
              <th class="px-4 py-3 text-right text-xs font-medium text-gray-400 uppercase">Slug</th>
              <th class="px-4 py-3 text-right text-xs font-medium text-gray-400 uppercase">Tools</th>
              <th class="px-4 py-3 text-right text-xs font-medium text-gray-400 uppercase">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-800">
            <tr
              v-for="tag in tags"
              :key="tag.id"
              class="hover:bg-dark-bg/30 transition-colors"
            >
              <!-- ID -->
              <td class="px-4 py-3 text-gray-500 text-sm">
                {{ tag.id }}
              </td>

              <!-- Name -->
              <td class="px-4 py-3">
                <span class="px-3 py-1 text-sm rounded-full bg-neon-blue/10 text-neon-blue border border-neon-blue/30">
                  {{ tag.name }}
                </span>
              </td>

              <!-- Slug -->
              <td class="px-4 py-3 text-gray-500 text-sm font-mono">
                {{ tag.slug }}
              </td>

              <!-- Tool Count -->
              <td class="px-4 py-3 text-gray-300">
                {{ tag.tool_count || 0 }}
              </td>

              <!-- Actions -->
              <td class="px-4 py-3">
                <div class="flex items-center gap-2">
                  <button
                    @click="openEditModal(tag)"
                    class="p-2 text-gray-400 hover:text-neon-blue transition-colors"
                    title="Edit"
                  >
                    <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                    </svg>
                  </button>
                  <button
                    @click="deleteTag(tag)"
                    class="p-2 text-gray-400 hover:text-red-400 transition-colors"
                    title="Delete"
                  >
                    <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Tag Form Modal -->
    <TagFormModal
      v-if="showModal"
      :tag="editingTag"
      @close="closeModal"
      @saved="onTagSaved"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useSessionStore } from '@/stores/session'
import { storeToRefs } from 'pinia'
import { apiClient } from '@/lib/apiClient'
import type { TagWithCount } from '@/lib/types'
import Spinner from '@/components/common/Spinner.vue'
import TagFormModal from '@/components/admin/TagFormModal.vue'

const router = useRouter()
const sessionStore = useSessionStore()
const { user, isAuthenticated } = storeToRefs(sessionStore)

// Check admin access
watch([isAuthenticated, user], ([authenticated, currentUser]) => {
  if (!authenticated || currentUser?.role !== 'admin') {
    router.push('/')
  }
}, { immediate: true })

const tags = ref<TagWithCount[]>([])
const loading = ref(false)
const showModal = ref(false)
const editingTag = ref<TagWithCount | null>(null)

async function fetchTags() {
  loading.value = true
  try {
    const response = await apiClient.get<{ data: TagWithCount[] }>('/admin/tags')
    tags.value = response.data
  } catch (error) {
    console.error('Failed to fetch tags:', error)
  } finally {
    loading.value = false
  }
}

function openCreateModal() {
  editingTag.value = null
  showModal.value = true
}

function openEditModal(tag: TagWithCount) {
  editingTag.value = tag
  showModal.value = true
}

function closeModal() {
  showModal.value = false
  editingTag.value = null
}

function onTagSaved() {
  closeModal()
  fetchTags()
}

async function deleteTag(tag: TagWithCount) {
  if (!confirm(`Are you sure you want to delete "${tag.name}"? This will remove the tag from all tools.`)) return

  try {
    await apiClient.delete(`/admin/tags/${tag.id}`)
    fetchTags()
  } catch (error) {
    console.error('Failed to delete tag:', error)
  }
}

onMounted(() => {
  fetchTags()
})
</script>
