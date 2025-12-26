<template>
  <div class="fixed inset-0 z-50 overflow-y-auto" dir="rtl">
    <!-- Backdrop -->
    <div class="fixed inset-0 bg-black/70" @click="$emit('close')"></div>

    <!-- Modal -->
    <div class="relative min-h-screen flex items-center justify-center p-4">
      <div class="relative bg-dark-surface border border-gray-700 rounded-xl w-full max-w-2xl max-h-[90vh] overflow-y-auto">
        <!-- Header -->
        <div class="sticky top-0 bg-dark-surface border-b border-gray-800 px-6 py-4 flex justify-between items-center">
          <h2 class="text-lg font-bold text-white">
            {{ isEditing ? 'Edit Tool' : 'Create Tool' }}
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
              placeholder="e.g., chatgpt"
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
              placeholder="Tool name"
            />
          </div>

          <!-- Category -->
          <div>
            <label class="block text-sm font-medium text-gray-300 mb-2">
              Category <span class="text-red-400">*</span>
            </label>
            <select
              v-model="form.primary_category_id"
              required
              class="w-full px-4 py-2 bg-dark-bg border border-gray-700 rounded-lg text-white focus:outline-none focus:border-neon-blue"
            >
              <option value="" disabled>Select category</option>
              <option v-for="cat in categories" :key="cat.id" :value="cat.id">
                {{ cat.name }}
              </option>
            </select>
          </div>

          <!-- Logo URL -->
          <div>
            <label class="block text-sm font-medium text-gray-300 mb-2">Logo URL</label>
            <input
              v-model="form.logo_url"
              type="url"
              class="w-full px-4 py-2 bg-dark-bg border border-gray-700 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-neon-blue"
              placeholder="https://..."
            />
          </div>

          <!-- Tagline -->
          <div>
            <label class="block text-sm font-medium text-gray-300 mb-2">Tagline</label>
            <input
              v-model="form.tagline"
              type="text"
              class="w-full px-4 py-2 bg-dark-bg border border-gray-700 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-neon-blue"
              placeholder="Short description"
            />
          </div>

          <!-- Description -->
          <div>
            <label class="block text-sm font-medium text-gray-300 mb-2">Description</label>
            <textarea
              v-model="form.description"
              rows="3"
              class="w-full px-4 py-2 bg-dark-bg border border-gray-700 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-neon-blue resize-none"
              placeholder="Full description"
            ></textarea>
          </div>

          <!-- Best For -->
          <div>
            <label class="block text-sm font-medium text-gray-300 mb-2">Best For</label>
            <input
              v-model="form.best_for"
              type="text"
              class="w-full px-4 py-2 bg-dark-bg border border-gray-700 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-neon-blue"
              placeholder="Who is this tool best for?"
            />
          </div>

          <!-- Primary Use Cases -->
          <div>
            <label class="block text-sm font-medium text-gray-300 mb-2">Primary Use Cases</label>
            <textarea
              v-model="form.primary_use_cases"
              rows="2"
              class="w-full px-4 py-2 bg-dark-bg border border-gray-700 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-neon-blue resize-none"
              placeholder="Main use cases, one per line"
            ></textarea>
          </div>

          <!-- Pricing Summary -->
          <div>
            <label class="block text-sm font-medium text-gray-300 mb-2">Pricing Summary</label>
            <input
              v-model="form.pricing_summary"
              type="text"
              class="w-full px-4 py-2 bg-dark-bg border border-gray-700 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-neon-blue"
              placeholder="e.g., Free / $20/mo / Contact for pricing"
            />
          </div>

          <!-- Has Free Tier -->
          <div class="flex items-center gap-3">
            <input
              v-model="form.has_free_tier"
              type="checkbox"
              id="has_free_tier"
              class="w-4 h-4 rounded bg-dark-bg border-gray-700 text-neon-blue focus:ring-neon-blue"
            />
            <label for="has_free_tier" class="text-sm text-gray-300">Has Free Tier</label>
          </div>

          <!-- Target Roles -->
          <div>
            <label class="block text-sm font-medium text-gray-300 mb-2">Target Roles</label>
            <input
              v-model="form.target_roles"
              type="text"
              class="w-full px-4 py-2 bg-dark-bg border border-gray-700 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-neon-blue"
              placeholder="e.g., Developers, Designers, Marketers"
            />
          </div>

          <!-- Platforms -->
          <div>
            <label class="block text-sm font-medium text-gray-300 mb-2">Platforms</label>
            <input
              v-model="form.platforms"
              type="text"
              class="w-full px-4 py-2 bg-dark-bg border border-gray-700 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-neon-blue"
              placeholder="e.g., Web, iOS, Android, API"
            />
          </div>

          <!-- Official URL -->
          <div>
            <label class="block text-sm font-medium text-gray-300 mb-2">Official URL</label>
            <input
              v-model="form.official_url"
              type="url"
              class="w-full px-4 py-2 bg-dark-bg border border-gray-700 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-neon-blue"
              placeholder="https://..."
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
              <span v-else>{{ isEditing ? 'Update Tool' : 'Create Tool' }}</span>
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
import type { Tool, Category, CreateToolInput, UpdateToolInput } from '@/lib/types'

const props = defineProps<{
  tool: Tool | null
  categories: Category[]
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'saved'): void
}>()

const isEditing = computed(() => !!props.tool)

const form = ref({
  slug: '',
  name: '',
  logo_url: '',
  tagline: '',
  description: '',
  best_for: '',
  primary_use_cases: '',
  pricing_summary: '',
  target_roles: '',
  platforms: '',
  has_free_tier: false,
  official_url: '',
  primary_category_id: '' as number | '',
})

const saving = ref(false)
const error = ref('')

onMounted(() => {
  if (props.tool) {
    form.value = {
      slug: props.tool.slug,
      name: props.tool.name,
      logo_url: props.tool.logo_url || '',
      tagline: props.tool.tagline || '',
      description: props.tool.description || '',
      best_for: props.tool.best_for || '',
      primary_use_cases: props.tool.primary_use_cases || '',
      pricing_summary: props.tool.pricing_summary || '',
      target_roles: props.tool.target_roles || '',
      platforms: props.tool.platforms || '',
      has_free_tier: props.tool.has_free_tier,
      official_url: props.tool.official_url || '',
      primary_category_id: props.tool.primary_category_id,
    }
  }
})

async function handleSubmit() {
  error.value = ''
  saving.value = true

  try {
    if (isEditing.value && props.tool) {
      // Update
      const input: UpdateToolInput = {
        name: form.value.name,
        logo_url: form.value.logo_url || undefined,
        tagline: form.value.tagline || undefined,
        description: form.value.description || undefined,
        best_for: form.value.best_for || undefined,
        primary_use_cases: form.value.primary_use_cases || undefined,
        pricing_summary: form.value.pricing_summary || undefined,
        target_roles: form.value.target_roles || undefined,
        platforms: form.value.platforms || undefined,
        has_free_tier: form.value.has_free_tier,
        official_url: form.value.official_url || undefined,
        primary_category_id: form.value.primary_category_id as number,
      }
      await apiClient.patch(`/admin/tools/${props.tool.id}`, input)
    } else {
      // Create
      const input: CreateToolInput = {
        slug: form.value.slug,
        name: form.value.name,
        logo_url: form.value.logo_url || undefined,
        tagline: form.value.tagline || undefined,
        description: form.value.description || undefined,
        best_for: form.value.best_for || undefined,
        primary_use_cases: form.value.primary_use_cases || undefined,
        pricing_summary: form.value.pricing_summary || undefined,
        target_roles: form.value.target_roles || undefined,
        platforms: form.value.platforms || undefined,
        has_free_tier: form.value.has_free_tier,
        official_url: form.value.official_url || undefined,
        primary_category_id: form.value.primary_category_id as number,
      }
      await apiClient.post('/admin/tools', input)
    }
    emit('saved')
  } catch (err: any) {
    error.value = err.message || 'Failed to save tool'
  } finally {
    saving.value = false
  }
}
</script>
