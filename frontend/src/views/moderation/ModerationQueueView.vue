<template>
  <div class="moderation-queue max-w-7xl mx-auto px-4 py-8" dir="rtl">
    <div class="flex items-center justify-between mb-8">
      <div>
        <h1 class="text-3xl font-bold text-white">قائمة الإشراف</h1>
        <p class="text-gray-400 mt-1">مراجعة البلاغات المقدمة</p>
      </div>
    </div>

    <!-- Filters -->
    <div class="flex flex-wrap gap-4 mb-6 p-4 bg-dark-surface rounded-lg">
      <div>
        <label class="block text-sm text-gray-400 mb-1">النوع</label>
        <select
          v-model="filters.type"
          @change="fetchReports"
          class="px-4 py-2 bg-dark-bg border border-gray-700 rounded-lg text-white focus:border-neon-blue focus:outline-none"
        >
          <option value="">الكل</option>
          <option value="tool">أدوات</option>
          <option value="review">مراجعات</option>
        </select>
      </div>
      <div>
        <label class="block text-sm text-gray-400 mb-1">الحالة</label>
        <select
          v-model="filters.status"
          @change="fetchReports"
          class="px-4 py-2 bg-dark-bg border border-gray-700 rounded-lg text-white focus:border-neon-blue focus:outline-none"
        >
          <option value="pending">قيد الانتظار</option>
          <option value="reviewed">تمت المراجعة</option>
          <option value="dismissed">مرفوض</option>
        </select>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="space-y-4">
      <div v-for="n in 5" :key="n" class="animate-pulse bg-dark-surface rounded-lg p-4">
        <div class="flex gap-4">
          <div class="h-6 w-20 bg-gray-700 rounded"></div>
          <div class="flex-1 h-6 bg-gray-700 rounded"></div>
          <div class="h-6 w-32 bg-gray-700 rounded"></div>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else-if="reports.length === 0" class="text-center py-16 bg-dark-surface rounded-lg">
      <svg class="w-16 h-16 mx-auto text-gray-600 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      <p class="text-gray-400 text-lg">لا توجد بلاغات</p>
    </div>

    <!-- Reports Table -->
    <div v-else class="bg-dark-surface rounded-lg overflow-hidden">
      <table class="w-full">
        <thead class="border-b border-gray-800">
          <tr class="text-gray-400 text-sm">
            <th class="text-right px-4 py-3">النوع</th>
            <th class="text-right px-4 py-3">المحتوى</th>
            <th class="text-right px-4 py-3">السبب</th>
            <th class="text-right px-4 py-3">المُبلِّغ</th>
            <th class="text-right px-4 py-3">التاريخ</th>
            <th class="text-right px-4 py-3">الحالة</th>
            <th class="text-right px-4 py-3">الإجراءات</th>
          </tr>
        </thead>
        <tbody>
          <template v-for="report in reports" :key="report.id">
            <tr
              class="border-b border-gray-800 hover:bg-dark-bg/50 cursor-pointer"
              @click="toggleExpand(report.id)"
            >
              <td class="px-4 py-4">
                <span
                  class="px-2 py-1 text-xs rounded-full"
                  :class="report.reportable_type === 'tool' ? 'bg-blue-500/20 text-blue-400' : 'bg-purple-500/20 text-purple-400'"
                >
                  {{ report.reportable_type === 'tool' ? 'أداة' : 'مراجعة' }}
                </span>
              </td>
              <td class="px-4 py-4 text-white">
                <div class="max-w-xs truncate">
                  {{ getContentPreview(report) }}
                </div>
              </td>
              <td class="px-4 py-4 text-gray-300">{{ getReasonLabel(report.reason) }}</td>
              <td class="px-4 py-4 text-gray-400">
                {{ report.reporter_user?.display_name || 'مجهول' }}
              </td>
              <td class="px-4 py-4 text-gray-400 text-sm">
                {{ formatDate(report.created_at) }}
              </td>
              <td class="px-4 py-4">
                <span
                  class="px-2 py-1 text-xs rounded-full"
                  :class="getStatusClass(report.status)"
                >
                  {{ getStatusLabel(report.status) }}
                </span>
              </td>
              <td class="px-4 py-4">
                <div class="flex gap-2">
                  <button
                    v-if="report.status === 'pending'"
                    @click.stop="dismissReport(report)"
                    class="px-3 py-1 text-xs bg-gray-700 text-gray-300 rounded hover:bg-gray-600 transition-colors"
                  >
                    رفض
                  </button>
                  <button
                    @click.stop="toggleExpand(report.id)"
                    class="px-3 py-1 text-xs bg-neon-blue/20 text-neon-blue rounded hover:bg-neon-blue/30 transition-colors"
                  >
                    {{ expandedId === report.id ? 'إخفاء' : 'تفاصيل' }}
                  </button>
                </div>
              </td>
            </tr>

            <!-- Expanded Details -->
            <tr v-if="expandedId === report.id" class="bg-dark-bg/30">
              <td colspan="7" class="px-4 py-6">
                <div class="space-y-4">
                  <!-- Content Details -->
                  <div class="grid md:grid-cols-2 gap-6">
                    <div>
                      <h4 class="text-sm font-medium text-gray-400 mb-2">تفاصيل المحتوى</h4>
                      <div class="bg-dark-bg rounded-lg p-4">
                        <template v-if="report.reportable_type === 'tool' && report.tool">
                          <div class="flex items-center gap-3 mb-2">
                            <img
                              v-if="report.tool.logo_url"
                              :src="report.tool.logo_url"
                              :alt="report.tool.name"
                              class="w-10 h-10 rounded-lg"
                            />
                            <div>
                              <div class="text-white font-medium">{{ report.tool.name }}</div>
                              <div class="text-gray-400 text-sm">{{ report.tool.tagline }}</div>
                            </div>
                          </div>
                          <a
                            :href="`/tools/${report.tool.slug}`"
                            target="_blank"
                            class="text-neon-blue text-sm hover:underline"
                          >
                            عرض في الموقع ←
                          </a>
                        </template>
                        <template v-else-if="report.reportable_type === 'review' && report.review">
                          <div class="text-gray-300 mb-2">{{ report.review.pros }}</div>
                          <div class="text-gray-400 text-sm mb-2">
                            بواسطة: {{ report.review.user?.display_name || 'مستخدم' }}
                          </div>
                          <div class="text-sm">
                            التقييم: {{ report.review.rating_overall }}/5
                          </div>
                        </template>
                      </div>
                    </div>

                    <div>
                      <h4 class="text-sm font-medium text-gray-400 mb-2">تفاصيل البلاغ</h4>
                      <div class="bg-dark-bg rounded-lg p-4 space-y-2">
                        <div>
                          <span class="text-gray-400">السبب:</span>
                          <span class="text-white mr-2">{{ getReasonLabel(report.reason) }}</span>
                        </div>
                        <div v-if="report.comment">
                          <span class="text-gray-400">التعليق:</span>
                          <p class="text-gray-300 mt-1">{{ report.comment }}</p>
                        </div>
                        <div>
                          <span class="text-gray-400">المُبلِّغ:</span>
                          <span class="text-white mr-2">
                            {{ report.reporter_user?.display_name || 'مجهول' }}
                          </span>
                        </div>
                      </div>
                    </div>
                  </div>

                  <!-- Actions for Reviews -->
                  <div v-if="report.reportable_type === 'review' && report.review" class="flex gap-3 pt-4 border-t border-gray-800">
                    <button
                      @click="approveReview(report.review.id)"
                      :disabled="actionLoading"
                      class="px-4 py-2 bg-green-500/20 text-green-400 rounded-lg hover:bg-green-500/30 transition-colors disabled:opacity-50"
                    >
                      قبول المراجعة
                    </button>
                    <button
                      @click="hideReview(report.review.id)"
                      :disabled="actionLoading"
                      class="px-4 py-2 bg-yellow-500/20 text-yellow-400 rounded-lg hover:bg-yellow-500/30 transition-colors disabled:opacity-50"
                    >
                      إخفاء المراجعة
                    </button>
                    <button
                      @click="confirmRemoveReview(report.review.id)"
                      :disabled="actionLoading"
                      class="px-4 py-2 bg-red-500/20 text-red-400 rounded-lg hover:bg-red-500/30 transition-colors disabled:opacity-50"
                    >
                      حذف المراجعة
                    </button>
                    <button
                      @click="viewHistory(report.review.id)"
                      class="px-4 py-2 bg-gray-700 text-gray-300 rounded-lg hover:bg-gray-600 transition-colors"
                    >
                      سجل الإجراءات
                    </button>
                  </div>

                  <!-- Mark Report as Reviewed -->
                  <div v-if="report.status === 'pending'" class="flex gap-3 pt-4 border-t border-gray-800">
                    <button
                      @click="markReviewed(report)"
                      :disabled="actionLoading"
                      class="px-4 py-2 bg-neon-blue text-dark-bg font-medium rounded-lg hover:bg-neon-blue-hover transition-colors disabled:opacity-50"
                    >
                      تم المراجعة
                    </button>
                    <button
                      @click="dismissReport(report)"
                      :disabled="actionLoading"
                      class="px-4 py-2 bg-gray-700 text-gray-300 rounded-lg hover:bg-gray-600 transition-colors disabled:opacity-50"
                    >
                      رفض البلاغ
                    </button>
                  </div>
                </div>
              </td>
            </tr>
          </template>
        </tbody>
      </table>

      <!-- Pagination -->
      <div v-if="totalPages > 1" class="flex justify-center gap-2 p-4 border-t border-gray-800">
        <button
          v-for="p in totalPages"
          :key="p"
          @click="goToPage(p)"
          class="px-3 py-1 rounded-lg transition-colors"
          :class="p === page ? 'bg-neon-blue text-dark-bg' : 'bg-gray-700 text-gray-300 hover:bg-gray-600'"
        >
          {{ p }}
        </button>
      </div>
    </div>

    <!-- Remove Confirmation Modal -->
    <div v-if="showRemoveModal" class="fixed inset-0 z-50 flex items-center justify-center">
      <div class="fixed inset-0 bg-black/70" @click="showRemoveModal = false"></div>
      <div class="relative bg-dark-surface border border-gray-700 rounded-xl p-6 max-w-md w-full mx-4">
        <h3 class="text-xl font-bold text-white mb-4">تأكيد الحذف</h3>
        <p class="text-gray-400 mb-6">
          هل أنت متأكد من حذف هذه المراجعة؟ سيتم تحديث تقييمات الأداة.
        </p>
        <div class="mb-4">
          <label class="block text-sm text-gray-400 mb-2">ملاحظات (اختياري)</label>
          <textarea
            v-model="removeNotes"
            class="w-full px-4 py-2 bg-dark-bg border border-gray-700 rounded-lg text-white placeholder-gray-500 focus:border-neon-blue focus:outline-none resize-none"
            rows="2"
            placeholder="سبب الحذف..."
          ></textarea>
        </div>
        <div class="flex gap-3 justify-end">
          <button
            @click="showRemoveModal = false"
            class="px-4 py-2 text-gray-400 hover:text-white transition-colors"
          >
            إلغاء
          </button>
          <button
            @click="removeReview"
            :disabled="actionLoading"
            class="px-4 py-2 bg-red-500 text-white font-medium rounded-lg hover:bg-red-600 transition-colors disabled:opacity-50"
          >
            حذف
          </button>
        </div>
      </div>
    </div>

    <!-- History Modal -->
    <div v-if="showHistoryModal" class="fixed inset-0 z-50 flex items-center justify-center">
      <div class="fixed inset-0 bg-black/70" @click="showHistoryModal = false"></div>
      <div class="relative bg-dark-surface border border-gray-700 rounded-xl p-6 max-w-lg w-full mx-4 max-h-[80vh] overflow-y-auto">
        <h3 class="text-xl font-bold text-white mb-4">سجل الإجراءات</h3>
        <div v-if="historyLoading" class="animate-pulse space-y-3">
          <div class="h-12 bg-gray-700 rounded"></div>
          <div class="h-12 bg-gray-700 rounded"></div>
        </div>
        <div v-else-if="history.length === 0" class="text-gray-400 text-center py-8">
          لا توجد إجراءات سابقة
        </div>
        <div v-else class="space-y-3">
          <div
            v-for="action in history"
            :key="action.id"
            class="bg-dark-bg rounded-lg p-4"
          >
            <div class="flex items-center justify-between mb-2">
              <span
                class="px-2 py-1 text-xs rounded-full"
                :class="getActionClass(action.action_type)"
              >
                {{ getActionLabel(action.action_type) }}
              </span>
              <span class="text-gray-400 text-sm">{{ formatDate(action.created_at) }}</span>
            </div>
            <div class="text-gray-300 text-sm">
              بواسطة: {{ action.moderator?.display_name || 'مشرف' }}
            </div>
            <div v-if="action.notes" class="text-gray-400 text-sm mt-1">
              {{ action.notes }}
            </div>
          </div>
        </div>
        <button
          @click="showHistoryModal = false"
          class="mt-4 w-full px-4 py-2 bg-gray-700 text-gray-300 rounded-lg hover:bg-gray-600 transition-colors"
        >
          إغلاق
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { apiClient } from '@/lib/apiClient'

interface Report {
  id: number
  reportable_type: 'tool' | 'review'
  reportable_id: number
  reason: string
  comment?: string
  status: string
  created_at: string
  reporter_user?: {
    id: number
    display_name: string
  }
  tool?: {
    id: number
    slug: string
    name: string
    tagline?: string
    logo_url?: string
  }
  review?: {
    id: number
    rating_overall: number
    pros: string
    cons: string
    user?: {
      id: number
      display_name: string
    }
  }
}

interface ModerationAction {
  id: number
  action_type: string
  notes?: string
  created_at: string
  moderator?: {
    id: number
    display_name: string
  }
}

const reports = ref<Report[]>([])
const loading = ref(true)
const actionLoading = ref(false)
const page = ref(1)
const pageSize = 20
const total = ref(0)
const expandedId = ref<number | null>(null)
const filters = ref({
  type: '',
  status: 'pending'
})

const showRemoveModal = ref(false)
const removeReviewId = ref<number | null>(null)
const removeNotes = ref('')

const showHistoryModal = ref(false)
const history = ref<ModerationAction[]>([])
const historyLoading = ref(false)

const totalPages = computed(() => Math.ceil(total.value / pageSize))

async function fetchReports() {
  loading.value = true
  try {
    const params: Record<string, string> = {
      page: page.value.toString(),
      page_size: pageSize.toString(),
      status: filters.value.status
    }
    if (filters.value.type) {
      params.type = filters.value.type
    }

    const response = await apiClient.get<{ data: Report[]; meta?: { total: number } }>('/admin/moderation/queue', { params })
    reports.value = response.data || []
    total.value = response.meta?.total || 0
  } catch (err) {
    console.error('Failed to fetch reports:', err)
  } finally {
    loading.value = false
  }
}

function toggleExpand(id: number) {
  expandedId.value = expandedId.value === id ? null : id
}

function getContentPreview(report: Report): string {
  if (report.reportable_type === 'tool' && report.tool) {
    return report.tool.name
  }
  if (report.reportable_type === 'review' && report.review) {
    return report.review.pros.substring(0, 50) + (report.review.pros.length > 50 ? '...' : '')
  }
  return 'محتوى غير متاح'
}

function getReasonLabel(reason: string): string {
  const labels: Record<string, string> = {
    spam: 'محتوى مزعج',
    abuse: 'محتوى مسيء',
    misinformation: 'معلومات مضللة',
    other: 'سبب آخر'
  }
  return labels[reason] || reason
}

function getStatusLabel(status: string): string {
  const labels: Record<string, string> = {
    pending: 'قيد الانتظار',
    reviewed: 'تمت المراجعة',
    dismissed: 'مرفوض'
  }
  return labels[status] || status
}

function getStatusClass(status: string): string {
  const classes: Record<string, string> = {
    pending: 'bg-yellow-500/20 text-yellow-400',
    reviewed: 'bg-green-500/20 text-green-400',
    dismissed: 'bg-gray-500/20 text-gray-400'
  }
  return classes[status] || 'bg-gray-500/20 text-gray-400'
}

function getActionLabel(action: string): string {
  const labels: Record<string, string> = {
    approve: 'قبول',
    hide: 'إخفاء',
    remove: 'حذف',
    restore: 'استعادة'
  }
  return labels[action] || action
}

function getActionClass(action: string): string {
  const classes: Record<string, string> = {
    approve: 'bg-green-500/20 text-green-400',
    hide: 'bg-yellow-500/20 text-yellow-400',
    remove: 'bg-red-500/20 text-red-400',
    restore: 'bg-blue-500/20 text-blue-400'
  }
  return classes[action] || 'bg-gray-500/20 text-gray-400'
}

function formatDate(dateStr: string): string {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  if (isNaN(date.getTime())) return ''
  return date.toLocaleDateString('ar-SA', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

async function approveReview(reviewId: number) {
  actionLoading.value = true
  try {
    await apiClient.patch(`/admin/moderation/reviews/${reviewId}/approve`)
    fetchReports()
  } catch (err) {
    console.error('Failed to approve review:', err)
  } finally {
    actionLoading.value = false
  }
}

async function hideReview(reviewId: number) {
  actionLoading.value = true
  try {
    await apiClient.patch(`/admin/moderation/reviews/${reviewId}/hide`)
    fetchReports()
  } catch (err) {
    console.error('Failed to hide review:', err)
  } finally {
    actionLoading.value = false
  }
}

function confirmRemoveReview(reviewId: number) {
  removeReviewId.value = reviewId
  removeNotes.value = ''
  showRemoveModal.value = true
}

async function removeReview() {
  if (!removeReviewId.value) return
  actionLoading.value = true
  try {
    await apiClient.patch(`/admin/moderation/reviews/${removeReviewId.value}/remove`, {
      notes: removeNotes.value || undefined
    })
    showRemoveModal.value = false
    fetchReports()
  } catch (err) {
    console.error('Failed to remove review:', err)
  } finally {
    actionLoading.value = false
  }
}

async function viewHistory(reviewId: number) {
  showHistoryModal.value = true
  historyLoading.value = true
  try {
    const response = await apiClient.get<{ data: ModerationAction[] }>(`/admin/moderation/history/${reviewId}`)
    history.value = response.data || []
  } catch (err) {
    console.error('Failed to fetch history:', err)
    history.value = []
  } finally {
    historyLoading.value = false
  }
}

async function markReviewed(report: Report) {
  actionLoading.value = true
  try {
    await apiClient.patch(`/admin/moderation/reports/${report.id}`, { status: 'reviewed' })
    fetchReports()
  } catch (err) {
    console.error('Failed to mark reviewed:', err)
  } finally {
    actionLoading.value = false
  }
}

async function dismissReport(report: Report) {
  actionLoading.value = true
  try {
    await apiClient.patch(`/admin/moderation/reports/${report.id}`, { status: 'dismissed' })
    fetchReports()
  } catch (err) {
    console.error('Failed to dismiss report:', err)
  } finally {
    actionLoading.value = false
  }
}

function goToPage(p: number) {
  page.value = p
  fetchReports()
}

// Handle Escape key to close modals
function handleKeyDown(e: KeyboardEvent) {
  if (e.key === 'Escape') {
    if (showRemoveModal.value) {
      showRemoveModal.value = false
    } else if (showHistoryModal.value) {
      showHistoryModal.value = false
    }
  }
}

onMounted(() => {
  fetchReports()
  document.addEventListener('keydown', handleKeyDown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeyDown)
})
</script>
