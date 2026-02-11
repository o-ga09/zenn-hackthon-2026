'use client'

import { useQuery } from '@tanstack/react-query'
import { TravelsResponse } from './types'

// キャッシュキー
export const TRAVELS_QUERY_KEY = (userId: string) => ['travels', userId]

/** ユーザーIDで旅行一覧を取得 */
export const useGetTravelsByUserId = (userId: string) => {
  return useQuery({
    queryKey: TRAVELS_QUERY_KEY(userId),
    queryFn: async (): Promise<TravelsResponse> => {
      // TODO: バックエンドに旅行取得エンドポイントが実装されたら、ここを実装する
      // const res = await apiClient.get(`/travels?userId=${userId}`)
      // return res.data

      // ダミーデータを返す（ビルドを通すため）
      return {
        travels: [],
        total: 0,
        message: 'Travel API is not implemented yet',
      }
    },
    enabled: !!userId,
    staleTime: 1000 * 60 * 5,
  })
}
