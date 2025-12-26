<template>
  <div class="tool-profile max-w-5xl mx-auto px-4 py-8">
    <!-- Loading State -->
    <ToolProfileSkeleton v-if="loading" />

    <!-- Error State: 404 -->
    <div v-else-if="notFound" class="text-center py-16">
      <div class="text-6xl mb-4">404</div>
      <h1 class="text-2xl font-bold mb-4">الأداة غير موجودة</h1>
      <p class="text-gray-400 mb-8">عذراً، لم نتمكن من العثور على الأداة التي تبحث عنها.</p>
      <router-link
        to="/"
        class="px-6 py-3 bg-neon-blue text-dark-bg font-bold rounded-lg hover:bg-neon-blue-hover transition-colors inline-flex items-center gap-2"
      >
        العودة للرئيسية
      </router-link>
    </div>

    <!-- Error State: Network Error -->
    <div v-else-if="error" class="text-center py-16">
      <div class="text-4xl mb-4">!</div>
      <h1 class="text-2xl font-bold mb-4">حدث خطأ</h1>
      <p class="text-gray-400 mb-8">{{ error }}</p>
      <button
        @click="fetchTool"
        class="px-6 py-3 bg-neon-blue text-dark-bg font-bold rounded-lg hover:bg-neon-blue-hover transition-colors"
      >
        إعادة المحاولة
      </button>
    </div>

    <!-- Tool Content -->
    <template v-else-if="tool">
      <div class="space-y-6">
        <!-- Hero Section -->
        <ToolHero
          :tool="tool"
          @bookmark-toggle="handleBookmarkToggle"
          @report="handleToolReport"
        />

        <!-- Engagement Stats -->
        <ToolStats :tool="tool" class="bg-dark-surface rounded-lg p-4" />

        <!-- Overview Section -->
        <ToolOverview :tool="tool" />

        <!-- Features Section -->
        <ToolFeatures :tool="tool" />

        <!-- Pricing Section -->
        <ToolPricing :tool="tool" />

        <!-- Media Gallery -->
        <ToolMediaGallery :tool="tool" />

        <!-- Alternatives -->
        <ToolAlternatives :tool-slug="tool.slug" />

        <!-- Reviews Section -->
        <ReviewList
          :reviews="reviews"
          :total-reviews="totalReviews"
          :avg-rating="tool.avg_rating_overall"
          :loading="reviewsLoading"
          :loading-more="reviewsLoadingMore"
          :has-more="hasMoreReviews"
          :initial-sort="reviewSort"
          @sort-change="handleReviewSortChange"
          @load-more="loadMoreReviews"
          @add-review="handleAddReview"
          @report="handleReviewReport"
        />
      </div>

      <!-- Review Form Modal -->
      <ReviewForm
        v-if="showReviewForm && tool"
        :tool-slug="tool.slug"
        @close="handleReviewFormClose"
        @success="handleReviewSuccess"
      />

      <!-- Report Modal -->
      <ReportModal
        v-if="showReportModal"
        :reportable-type="reportType"
        :reportable-id="reportId"
        @close="showReportModal = false"
        @reported="showReportModal = false"
      />
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { apiClient } from '../lib/apiClient'
import type { Tool, Review, ApiResponse } from '../lib/types'
import { useBookmarksStore } from '../stores/bookmarks'
import { useSessionStore } from '../stores/session'
import ToolHero from '../components/tools/ToolHero.vue'
import ToolOverview from '../components/tools/ToolOverview.vue'
import ToolFeatures from '../components/tools/ToolFeatures.vue'
import ToolPricing from '../components/tools/ToolPricing.vue'
import ToolMediaGallery from '../components/tools/ToolMediaGallery.vue'
import ToolAlternatives from '../components/tools/ToolAlternatives.vue'
import ToolStats from '../components/tools/ToolStats.vue'
import ToolProfileSkeleton from '../components/tools/ToolProfileSkeleton.vue'
import ReviewList from '../components/reviews/ReviewList.vue'
import ReviewForm from '../components/reviews/ReviewForm.vue'
import ReportModal from '../components/moderation/ReportModal.vue'

const route = useRoute()
const router = useRouter()
const bookmarksStore = useBookmarksStore()
const sessionStore = useSessionStore()

const tool = ref<Tool | null>(null)
const loading = ref(true)
const error = ref<string | null>(null)
const notFound = ref(false)

// Reviews state
const reviews = ref<Review[]>([])
const totalReviews = ref(0)
const reviewsLoading = ref(false)
const reviewsLoadingMore = ref(false)
const reviewSort = ref('newest')
const reviewPage = ref(1)
const reviewPageSize = 10

const hasMoreReviews = computed(() => reviews.value.length < totalReviews.value)
const showReviewForm = ref(false)
const showReportModal = ref(false)
const reportType = ref<'tool' | 'review'>('tool')
const reportId = ref<string | number>('')

const fetchTool = async () => {
  const slug = route.params.slug as string
  if (!slug) return

  loading.value = true
  error.value = null
  notFound.value = false

  try {
    const response = await apiClient.get<ApiResponse<Tool>>(`/tools/${slug}`)
    tool.value = response.data
  } catch (err: any) {
    if (err.message?.includes('not found') || err.message?.includes('404')) {
      notFound.value = true
    } else {
      error.value = 'حدث خطأ أثناء تحميل بيانات الأداة. يرجى المحاولة مرة أخرى.'
    }
  } finally {
    loading.value = false
  }
}

const handleBookmarkToggle = async (toolId: number) => {
  if (!sessionStore.isAuthenticated) {
    // Could redirect to login or show message
    return
  }

  try {
    if (bookmarksStore.isBookmarked(toolId)) {
      await bookmarksStore.removeBookmark(toolId)
    } else {
      await bookmarksStore.addBookmark(toolId)
    }
  } catch (err) {
    console.error('Failed to toggle bookmark:', err)
  }
}

// Reviews functions
const fetchReviews = async (append = false) => {
  const slug = route.params.slug as string
  if (!slug) return

  if (append) {
    reviewsLoadingMore.value = true
  } else {
    reviewsLoading.value = true
    reviewPage.value = 1
    reviews.value = []
  }

  try {
    const response = await apiClient.get<ApiResponse<Review[]>>(`/tools/${slug}/reviews`, {
      params: {
        sort: reviewSort.value,
        page: reviewPage.value,
        page_size: reviewPageSize
      }
    })
    if (append) {
      reviews.value = [...reviews.value, ...response.data]
    } else {
      reviews.value = response.data
    }
    totalReviews.value = response.meta?.total || 0
  } catch (err) {
    console.error('Failed to fetch reviews:', err)
  } finally {
    reviewsLoading.value = false
    reviewsLoadingMore.value = false
  }
}

const handleReviewSortChange = (sort: string) => {
  reviewSort.value = sort
  fetchReviews()
}

const loadMoreReviews = () => {
  reviewPage.value++
  fetchReviews(true)
}

const handleAddReview = () => {
  if (!sessionStore.isAuthenticated) {
    // Redirect to login
    router.push('/login')
    return
  }
  showReviewForm.value = true
}

const handleReviewFormClose = () => {
  showReviewForm.value = false
}

const handleReviewSuccess = () => {
  showReviewForm.value = false
  // Refresh reviews list
  fetchReviews()
  // Could also show a success toast here
}

const handleToolReport = () => {
  if (!tool.value) return
  reportType.value = 'tool'
  reportId.value = tool.value.slug
  showReportModal.value = true
}

const handleReviewReport = (reviewId: number) => {
  reportType.value = 'review'
  reportId.value = reviewId
  showReportModal.value = true
}

// Watch for route changes to handle navigation between tools
watch(() => route.params.slug, () => {
  fetchTool()
  fetchReviews()
})

onMounted(() => {
  fetchTool()
  fetchReviews()
})
</script>
