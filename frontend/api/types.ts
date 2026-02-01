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

// メディアアップロード関連の型定義（画像・動画統合）
export interface MediaUploadRequest {
  file: File // アップロードするファイル（画像または動画）
}

export interface MediaUploadResponse {
  file_id: string // アップロードされたメディアのファイルID
  url: string // メディアの取得URL
}

// 画像管理関連の型定義（後方互換性のため保持）
export interface MediaImageUploadRequest {
  file: File // 画像ファイル
}

export interface MediaImageUploadResponse {
  file_id: string // アップロードされた画像のファイルID
  url: string // 画像の取得URL
}

// 動画管理関連の型定義（後方互換性のため保持）
export interface MediaVideoUploadRequest {
  file: File // 動画ファイル
}

export interface MediaVideoUploadResponse {
  file_id: string // アップロードされた動画のファイルID
  url: string // 動画の取得URL
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
  status: 'pending' | 'uploading' | 'completed' | 'failed' // アップロード状態
  progress: number // 進捗率（0.0〜1.0）
  error_message?: string // エラーメッセージ
  created_at: string
  updated_at: string
}

export interface MediaListResponse {
  media: Media[]
  total: number
  message?: string
}

// メディア分析関連の型定義
export interface MediaAnalysisOutput {
  fileId: string
  type: string
  description: string
  objects: string[]
  landmarks: string[]
  activities: string[]
  mood: string
  suggestedCaption: string
}

export interface MediaAnalysisBatchResponse {
  results: MediaAnalysisOutput[]
  summary: {
    totalItems: number
    successfulItems: number
    failedItems: number
    uniqueLocations: string[]
    uniqueActivities: string[]
    overallMood: string
  }
}

// メディア分析レスポンスの型定義
export interface AnalyzeMediaResponse {
  media_ids: string[]
  status: string
}

// メディア分析SSEレスポンスの型定義
export interface MediaStatusResponse {
  medias: Media[]
  total_items: number
  completed_items: number
  failed_items: number
  all_completed: boolean
}

// メディア分析結果の型定義
export interface MediaAnalyticsResponse {
  file_id: string
  description: string
  mood: string
  objects: string[]
  landmarks: string[]
  activities: string[]
}

// メディア分析結果更新リクエストの型定義
export interface UpdateMediaAnalyticsRequest {
  description?: string
  mood?: string
  objects?: string[]
  landmarks?: string[]
  activities?: string[]
}
