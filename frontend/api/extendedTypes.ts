// 旅行メモリー関連の型定義
export interface TravelMemory {
  id: string
  title: string
  location: string
  date: string
  thumbnailUrl: string
  likes: number
  description?: string
}

export interface TravelMemoriesResponse {
  memories: TravelMemory[]
}
