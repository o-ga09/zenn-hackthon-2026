// ユーザー関連の拡張型定義
import { User as BaseUser } from './types'

export interface ExtendedUser extends BaseUser {
  bio?: string
  occupation?: string
  avatar_url?: string
  followers_count?: number
  following_count?: number
}

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
