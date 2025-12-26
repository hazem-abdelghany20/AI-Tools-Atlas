// API Response types
export interface ApiResponse<T> {
  data: T
  meta?: {
    page: number
    page_size: number
    total: number
  }
}

export interface ApiError {
  code: string
  message: string
  details?: any
}

// Domain types
export interface User {
  id: number
  email: string
  display_name?: string
  role: string
  created_at: string
  updated_at: string
}

export interface Category {
  id: number
  slug: string
  name: string
  description?: string
  icon_url?: string
  display_order: number
}

export interface ToolFeature {
  name: string
  description?: string
  icon?: string
}

export interface Tool {
  id: number
  slug: string
  name: string
  logo_url?: string
  tagline?: string
  description?: string
  best_for?: string
  primary_use_cases?: string
  pricing_summary?: string
  target_roles?: string
  platforms?: string
  has_free_tier: boolean
  official_url?: string
  primary_category_id: number
  primary_category?: Category
  avg_rating_overall: number
  review_count: number
  bookmark_count: number
  trending_score: number
  tags?: Tag[]
  media?: Media[]
  badges?: Badge[]
  features?: ToolFeature[] | string
  archived_at?: string | null
}

// Admin types
export interface CreateToolInput {
  slug: string
  name: string
  logo_url?: string
  tagline?: string
  description?: string
  best_for?: string
  primary_use_cases?: string
  pricing_summary?: string
  target_roles?: string
  platforms?: string
  has_free_tier: boolean
  official_url?: string
  primary_category_id: number
}

export interface UpdateToolInput {
  name?: string
  logo_url?: string
  tagline?: string
  description?: string
  best_for?: string
  primary_use_cases?: string
  pricing_summary?: string
  target_roles?: string
  platforms?: string
  has_free_tier?: boolean
  official_url?: string
  primary_category_id?: number
}

export interface Tag {
  id: number
  slug: string
  name: string
}

export interface Media {
  id: number
  tool_id: number
  type: 'screenshot' | 'video'
  url: string
  thumbnail_url?: string
  display_order: number
}

export interface Badge {
  id: number
  slug: string
  name: string
  description?: string
  icon_url?: string
}

export interface ReviewUser {
  id: number
  display_name: string
}

export interface Review {
  id: number
  tool_id?: number
  user_id?: number
  rating_overall: number
  rating_ease_of_use?: number
  rating_value?: number
  rating_accuracy?: number
  rating_speed?: number
  rating_support?: number
  pros?: string
  cons?: string
  primary_use_case?: string
  reviewer_role?: string
  company_size?: string
  usage_context?: string
  helpful_count: number
  created_at: string
  user: ReviewUser
}

export interface Bookmark {
  id: number
  user_id: number
  tool_id: number
  created_at: string
}

// Admin category types
export interface CategoryWithCount extends Category {
  tool_count: number
}

export interface CreateCategoryInput {
  slug: string
  name: string
  description?: string
  icon_url?: string
  display_order: number
}

export interface UpdateCategoryInput {
  name?: string
  description?: string
  icon_url?: string
  display_order?: number
}

// Admin tag types
export interface TagWithCount extends Tag {
  tool_count: number
}

export interface CreateTagInput {
  slug: string
  name: string
}

export interface UpdateTagInput {
  name?: string
}
