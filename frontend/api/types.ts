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

// 画像管理関連の型定義
export interface MediaImageUploadRequest {
  base64_data: string // Base64エンコードされた画像データ
}

export interface MediaImageUploadResponse {
  file_id: string // アップロードされた画像のファイルID
  url: string // 画像の取得URL
}

export interface MediaGetRequest {
  key: string // 画像キー
}

export interface MediaGetResponse {
  file_id: string // 画像のファイルID
  url: string // 画像の取得URL
}

export interface MediaDeleteRequest {
  key: string // 画像キー
}

// メディア一覧関連の型定義
export interface Media {
  id: string
  type: string // "image" or "video"
  content_type: string // MIMEタイプ
  size: number // ファイルサイズ（バイト単位）
  url: string // ファイルのURL
  image_data?: string // 画像データ（typeがimageの場合に存在）
  created_at: string
  updated_at: string
}

export interface MediaListResponse {
  media: Media[]
  total: number
  message?: string
}
