import { useQuery } from '@tanstack/react-query'
import apiClient from '@/api/client'
import { TravelMemoriesResponse } from '@/api/extendedTypes'

// キャッシュのキー
export const USER_MEMORIES_QUERY_KEY = (userId: string) => ['users', userId, 'memories']

/**
 * ユーザーの旅行メモリーを取得するフック
 */
export const useGetUserMemories = (userId: string) => {
  return useQuery({
    queryKey: USER_MEMORIES_QUERY_KEY(userId),
    queryFn: async (): Promise<TravelMemoriesResponse> => {
      const response = await apiClient.get(`/users/${userId}/memories`)
      return response.data
    },
    enabled: !!userId,
  })
}
