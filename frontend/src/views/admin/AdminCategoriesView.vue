<template>
  <div class="admin-categories-view min-h-screen bg-dark-bg" dir="rtl">
    <!-- Header -->
    <div class="bg-dark-surface border-b border-gray-800 px-6 py-4">
      <div class="max-w-7xl mx-auto flex justify-between items-center">
        <div>
          <h1 class="text-xl font-bold text-white">Admin - Categories</h1>
          <p class="text-gray-400 text-sm">Manage tool categories</p>
        </div>
        <button
          @click="openCreateModal"
          class="px-4 py-2 bg-neon-blue text-dark-bg font-medium rounded-lg hover:bg-neon-blue-hover transition-colors flex items-center gap-2"
        >
          <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          Add Category
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
        <div v-else-if="categories.length === 0" class="p-8 text-center text-gray-400">
          No categories found
        </div>

        <!-- Data Table -->
        <table v-else class="w-full">
          <thead class="bg-dark-bg/50">
            <tr>
              <th class="px-4 py-3 text-right text-xs font-medium text-gray-400 uppercase">Order</th>
              <th class="px-4 py-3 text-right text-xs font-medium text-gray-400 uppercase">Category</th>
              <th class="px-4 py-3 text-right text-xs font-medium text-gray-400 uppercase">Slug</th>
              <th class="px-4 py-3 text-right text-xs font-medium text-gray-400 uppercase">Description</th>
              <th class="px-4 py-3 text-right text-xs font-medium text-gray-400 uppercase">Tools</th>
              <th class="px-4 py-3 text-right text-xs font-medium text-gray-400 uppercase">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-800">
            <tr
              v-for="category in categories"
              :key="category.id"
              class="hover:bg-dark-bg/30 transition-colors"
            >
              <!-- Order -->
              <td class="px-4 py-3 text-gray-300">
                {{ category.display_order }}
              </td>

              <!-- Category Info -->
              <td class="px-4 py-3">
                <div class="flex items-center gap-3">
                  <img
                    v-if="category.icon_url"
                    :src="category.icon_url"
                    :alt="category.name"
                    class="w-8 h-8 rounded object-cover"
                  />
                  <div
                    v-else
                    class="w-8 h-8 rounded bg-neon-blue/20 flex items-center justify-center"
                  >
                    <span class="text-neon-blue font-bold text-sm">{{ category.name.charAt(0) }}</span>
                  </div>
                  <div class="font-medium text-white">{{ category.name }}</div>
                </div>
              </td>

              <!-- Slug -->
              <td class="px-4 py-3 text-gray-500 text-sm font-mono">
                {{ category.slug }}
              </td>

              <!-- Description -->
              <td class="px-4 py-3 text-gray-300 max-w-xs truncate">
                {{ category.description || '-' }}
              </td>

              <!-- Tool Count -->
              <td class="px-4 py-3">
                <span class="px-2 py-1 text-xs rounded-full bg-neon-blue/20 text-neon-blue">
                  {{ category.tool_count || 0 }} tools
                </span>
              </td>

              <!-- Actions -->
              <td class="px-4 py-3">
                <div class="flex items-center gap-2">
                  <button
                    @click="openEditModal(category)"
                    class="p-2 text-gray-400 hover:text-neon-blue transition-colors"
                    title="Edit"
                  >
                    <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                    </svg>
                  </button>
                  <button
                    @click="deleteCategory(category)"
                    :disabled="category.tool_count > 0"
                    class="p-2 transition-colors"
                    :class="category.tool_count > 0 ? 'text-gray-600 cursor-not-allowed' : 'text-gray-400 hover:text-red-400'"
                    :title="category.tool_count > 0 ? 'Cannot delete - has tools' : 'Delete'"
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

    <!-- Category Form Modal -->
    <CategoryFormModal
      v-if="showModal"
      :category="editingCategory"
      @close="closeModal"
      @saved="onCategorySaved"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useSessionStore } from '@/stores/session'
import { storeToRefs } from 'pinia'
import { apiClient } from '@/lib/apiClient'
import type { CategoryWithCount } from '@/lib/types'
import Spinner from '@/components/common/Spinner.vue'
import CategoryFormModal from '@/components/admin/CategoryFormModal.vue'

const router = useRouter()
const sessionStore = useSessionStore()
const { user, isAuthenticated } = storeToRefs(sessionStore)

// Check admin access
watch([isAuthenticated, user], ([authenticated, currentUser]) => {
  if (!authenticated || currentUser?.role !== 'admin') {
    router.push('/')
  }
}, { immediate: true })

const categories = ref<CategoryWithCount[]>([])
const loading = ref(false)
const showModal = ref(false)
const editingCategory = ref<CategoryWithCount | null>(null)

async function fetchCategories() {
  loading.value = true
  try {
    const response = await apiClient.get<{ data: CategoryWithCount[] }>('/admin/categories')
    categories.value = response.data
  } catch (error) {
    console.error('Failed to fetch categories:', error)
  } finally {
    loading.value = false
  }
}

function openCreateModal() {
  editingCategory.value = null
  showModal.value = true
}

function openEditModal(category: CategoryWithCount) {
  editingCategory.value = category
  showModal.value = true
}

function closeModal() {
  showModal.value = false
  editingCategory.value = null
}

function onCategorySaved() {
  closeModal()
  fetchCategories()
}

async function deleteCategory(category: CategoryWithCount) {
  if (category.tool_count > 0) {
    alert('Cannot delete category with tools. Reassign tools first.')
    return
  }
  if (!confirm(`Are you sure you want to delete "${category.name}"?`)) return

  try {
    await apiClient.delete(`/admin/categories/${category.id}`)
    fetchCategories()
  } catch (error) {
    console.error('Failed to delete category:', error)
  }
}

onMounted(() => {
  fetchCategories()
})
</script>
