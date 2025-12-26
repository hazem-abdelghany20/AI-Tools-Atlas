<template>
  <div class="fixed inset-0 z-50 overflow-y-auto" dir="rtl">
    <!-- Backdrop -->
    <div class="fixed inset-0 bg-black/70" @click="$emit('close')"></div>

    <!-- Modal -->
    <div class="relative min-h-screen flex items-center justify-center p-4">
      <div class="relative bg-dark-surface border border-gray-700 rounded-xl w-full max-w-md max-h-[90vh] overflow-y-auto">
        <!-- Header -->
        <div class="sticky top-0 bg-dark-surface border-b border-gray-800 px-6 py-4 flex justify-between items-center">
          <h2 class="text-lg font-bold text-white">
            {{ isEditing ? 'Edit Tag' : 'Create Tag' }}
          </h2>
          <button @click="$emit('close')" class="text-gray-400 hover:text-white">
            <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <!-- Form -->
        <form @submit.prevent="handleSubmit" class="p-6 space-y-6">
          <!-- Slug (only for create) -->
          <div v-if="!isEditing">
            <label class="block text-sm font-medium text-gray-300 mb-2">
              Slug <span class="text-red-400">*</span>
            </label>
            <input
              v-model="form.slug"
              type="text"
              required
              pattern="[a-z0-9-]+"
              class="w-full px-4 py-2 bg-dark-bg border border-gray-700 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-neon-blue"
              placeholder="e.g., machine-learning"
            />
            <p class="text-xs text-gray-500 mt-1">Lowercase letters, numbers, and hyphens only</p>
          </div>

          <!-- Name -->
          <div>
            <label class="block text-sm font-medium text-gray-300 mb-2">
              Name <span class="text-red-400">*</span>
            </label>
            <input
              v-model="form.name"
              type="text"
              required
              class="w-full px-4 py-2 bg-dark-bg border border-gray-700 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-neon-blue"
              placeholder="Tag name"
            />
          </div>

          <!-- Error Message -->
          <div v-if="error" class="p-3 bg-red-500/20 border border-red-500/50 rounded-lg text-red-400 text-sm">
            {{ error }}
          </div>

          <!-- Actions -->
          <div class="flex justify-end gap-3 pt-4 border-t border-gray-800">
            <button
              type="button"
              @click="$emit('close')"
              class="px-4 py-2 text-gray-400 hover:text-white transition-colors"
            >
              Cancel
            </button>
            <button
              type="submit"
              :disabled="saving"
              class="px-6 py-2 bg-neon-blue text-dark-bg font-medium rounded-lg hover:bg-neon-blue-hover transition-colors disabled:opacity-50 flex items-center gap-2"
            >
              <span v-if="saving">Saving...</span>
              <span v-else>{{ isEditing ? 'Update Tag' : 'Create Tag' }}</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { apiClient } from '@/lib/apiClient'
import type { TagWithCount, CreateTagInput, UpdateTagInput } from '@/lib/types'

const props = defineProps<{
  tag: TagWithCount | null
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'saved'): void
}>()

const isEditing = computed(() => !!props.tag)

const form = ref({
  slug: '',
  name: '',
})

const saving = ref(false)
const error = ref('')

onMounted(() => {
  if (props.tag) {
    form.value = {
      slug: props.tag.slug,
      name: props.tag.name,
    }
  }
})

async function handleSubmit() {
  error.value = ''
  saving.value = true

  try {
    if (isEditing.value && props.tag) {
      // Update
      const input: UpdateTagInput = {
        name: form.value.name,
      }
      await apiClient.patch(`/admin/tags/${props.tag.id}`, input)
    } else {
      // Create
      const input: CreateTagInput = {
        slug: form.value.slug,
        name: form.value.name,
      }
      await apiClient.post('/admin/tags', input)
    }
    emit('saved')
  } catch (err: any) {
    error.value = err.message || 'Failed to save tag'
  } finally {
    saving.value = false
  }
}
</script>
