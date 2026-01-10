// フォロー関連の型定義
export interface Follow {
  user_id: string
  follow_count: number
  follower_count: number
}

// ソーシャルアカウント関連の型定義
export interface SocialAccount {
  x_url: string
  instagram_url: string
  facebook_url: string
  tiktok_url: string
  youtube_url: string
}

// 旅行情報関連の型定義
export interface Travel {
  id: string
  userId: string
  title: string
  description: string
  startDate: string
  endDate: string
  sharedId: string
  thumbnail: string
  created_at: string
  updated_at: string
  version: number
}

export interface TravelInput {
  id?: string // サーバー側で生成される場合はオプショナル
  userId: string
  title: string
  description: string
  startDate: string
  endDate: string
  sharedId: string
  thumbnail: string
}

export interface TravelsResponse {
  travels: Travel[]
  total: number
  message?: string
  next_page_token?: string
}
