<template>
  <div class="profile-view max-w-5xl mx-auto px-4 py-8" dir="rtl">
    <!-- Header -->
    <div class="mb-8">
      <div class="flex items-center gap-4 mb-4">
        <!-- Avatar -->
        <div class="w-16 h-16 bg-neon-blue/20 rounded-full flex items-center justify-center">
          <span class="text-neon-blue font-bold text-2xl">{{ userInitials }}</span>
        </div>
        <div>
          <h1 class="text-2xl font-bold text-white">{{ user?.display_name }}</h1>
          <p class="text-gray-400">{{ user?.email }}</p>
          <p class="text-gray-500 text-sm mt-1">
            عضو منذ {{ formatDate(user?.created_at) }}
          </p>
        </div>
      </div>
    </div>

    <!-- Tabs -->
    <div class="border-b border-gray-800 mb-6">
      <div class="flex gap-6">
        <button
          @click="activeTab = 'reviews'"
          class="pb-3 text-sm font-medium transition-colors border-b-2"
          :class="activeTab === 'reviews'
            ? 'text-neon-blue border-neon-blue'
            : 'text-gray-400 border-transparent hover:text-white'"
        >
          مراجعاتي ({{ reviewsTotal }})
        </button>
        <button
          @click="activeTab = 'bookmarks'"
          class="pb-3 text-sm font-medium transition-colors border-b-2"
          :class="activeTab === 'bookmarks'
            ? 'text-neon-blue border-neon-blue'
            : 'text-gray-400 border-transparent hover:text-white'"
        >
          المفضلة ({{ bookmarksCount }})
        </button>
      </div>
    </div>

    <!-- Reviews Tab -->
    <div v-if="activeTab === 'reviews'">
      <!-- Loading -->
      <div v-if="reviewsLoading" class="space-y-4">
        <SkeletonCard v-for="n in 3" :key="n" />
      </div>

      <!-- Empty State -->
      <div v-else-if="reviews.length === 0" class="text-center py-12">
        <svg class="w-16 h-16 mx-auto text-gray-600 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
        </svg>
        <h3 class="text-lg font-semibold text-white mb-2">لم تكتب أي مراجعات بعد</h3>
        <p class="text-gray-400 mb-6">شارك رأيك في الأدوات التي استخدمتها</p>
        <router-link
          to="/"
          class="inline-block px-6 py-3 bg-neon-blue text-dark-bg font-medium rounded-lg hover:bg-neon-blue-hover transition-colors"
        >
          استكشف الأدوات
        </router-link>
      </div>

      <!-- Reviews List -->
      <div v-else class="space-y-4">
        <div
          v-for="review in reviews"
          :key="review.id"
          class="bg-dark-surface border border-gray-800 rounded-xl p-5"
        >
          <!-- Tool Info -->
          <router-link
            :to="`/tools/${review.tool.slug}`"
            class="flex items-center gap-3 mb-4 hover:opacity-80 transition-opacity"
          >
            <img
              v-if="review.tool.logo_url"
              :src="review.tool.logo_url"
              :alt="review.tool.name"
              class="w-10 h-10 rounded-lg object-cover"
            />
            <div
              v-else
              class="w-10 h-10 rounded-lg bg-neon-blue/20 flex items-center justify-center"
            >
              <span class="text-neon-blue font-bold">{{ review.tool.name?.charAt(0) }}</span>
            </div>
            <div>
              <h4 class="font-medium text-white">{{ review.tool.name }}</h4>
              <p class="text-sm text-gray-500">{{ formatDate(review.created_at) }}</p>
            </div>
          </router-link>

          <!-- Rating -->
          <div class="flex items-center gap-2 mb-3">
            <div class="flex">
              <svg
                v-for="i in 5"
                :key="i"
                class="w-4 h-4"
                :class="i <= review.rating_overall ? 'text-yellow-400' : 'text-gray-600'"
                fill="currentColor"
                viewBox="0 0 20 20"
              >
                <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"/>
              </svg>
            </div>
            <span
              class="text-xs px-2 py-0.5 rounded-full"
              :class="getStatusClass(review.moderation_status)"
            >
              {{ getStatusLabel(review.moderation_status) }}
            </span>
          </div>

          <!-- Pros/Cons -->
          <div class="space-y-2 text-sm">
            <div v-if="review.pros" class="flex gap-2">
              <span class="text-green-400">+</span>
              <span class="text-gray-300">{{ review.pros }}</span>
            </div>
            <div v-if="review.cons" class="flex gap-2">
              <span class="text-red-400">-</span>
              <span class="text-gray-300">{{ review.cons }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Bookmarks Tab -->
    <div v-if="activeTab === 'bookmarks'">
      <!-- Loading -->
      <div v-if="bookmarksLoading" class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <SkeletonCard v-for="n in 4" :key="n" />
      </div>

      <!-- Empty State -->
      <div v-else-if="bookmarkedTools.length === 0" class="text-center py-12">
        <svg class="w-16 h-16 mx-auto text-gray-600 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
        </svg>
        <h3 class="text-lg font-semibold text-white mb-2">لا توجد أدوات محفوظة</h3>
        <p class="text-gray-400 mb-6">ابدأ بحفظ الأدوات التي تهمك للرجوع إليها لاحقاً</p>
        <router-link
          to="/"
          class="inline-block px-6 py-3 bg-neon-blue text-dark-bg font-medium rounded-lg hover:bg-neon-blue-hover transition-colors"
        >
          استكشف الأدوات
        </router-link>
      </div>

      <!-- Bookmarks Grid -->
      <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <ToolCard
          v-for="tool in bookmarkedTools"
          :key="tool.id"
          :tool="tool"
          :show-compare="false"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useSessionStore } from '@/stores/session'
import { useBookmarksStore } from '@/stores/bookmarks'
import { storeToRefs } from 'pinia'
import { apiClient } from '@/lib/apiClient'
import ToolCard from '@/components/tools/ToolCard.vue'
import SkeletonCard from '@/components/common/SkeletonCard.vue'

interface UserReview {
  id: number
  rating_overall: number
  pros: string
  cons: string
  helpful_count: number
  moderation_status: string
  created_at: string
  tool: {
    slug: string
    name: string
    logo_url?: string
  }
}

const router = useRouter()
const sessionStore = useSessionStore()
const bookmarksStore = useBookmarksStore()

const { user, isAuthenticated } = storeToRefs(sessionStore)
const { bookmarkedTools, loading: bookmarksLoading } = storeToRefs(bookmarksStore)

const activeTab = ref<'reviews' | 'bookmarks'>('reviews')
const reviews = ref<UserReview[]>([])
const reviewsTotal = ref(0)
const reviewsLoading = ref(false)

const userInitials = computed(() => {
  if (!user.value?.display_name) return '?'
  return user.value.display_name.charAt(0).toUpperCase()
})

const bookmarksCount = computed(() => bookmarkedTools.value.length)

function formatDate(dateStr?: string): string {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString('ar-EG', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}

function getStatusClass(status: string): string {
  switch (status) {
    case 'approved':
      return 'bg-green-500/20 text-green-400'
    case 'pending':
      return 'bg-yellow-500/20 text-yellow-400'
    case 'rejected':
      return 'bg-red-500/20 text-red-400'
    default:
      return 'bg-gray-500/20 text-gray-400'
  }
}

function getStatusLabel(status: string): string {
  switch (status) {
    case 'approved':
      return 'موافق عليها'
    case 'pending':
      return 'قيد المراجعة'
    case 'rejected':
      return 'مرفوضة'
    default:
      return status
  }
}

async function fetchReviews() {
  reviewsLoading.value = true
  try {
    const response = await apiClient.get<{
      data: UserReview[]
      meta: { total: number }
    }>('/me/reviews')
    reviews.value = response.data
    reviewsTotal.value = response.meta.total
  } catch (error) {
    console.error('Failed to fetch user reviews:', error)
  } finally {
    reviewsLoading.value = false
  }
}

// Redirect if not authenticated
watch(isAuthenticated, (authenticated) => {
  if (!authenticated) {
    router.push('/')
  }
}, { immediate: true })

onMounted(() => {
  if (isAuthenticated.value) {
    fetchReviews()
    bookmarksStore.fetchBookmarks()
  }
})
</script>
