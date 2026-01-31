'use client'

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import {apiClient} from './client'
import { Travel, TravelInput, TravelsResponse } from './types'

// キャッシュのキー
export const TRAVELS_QUERY_KEY = ['travels']
export const TRAVEL_QUERY_KEY = (travelId: string) => ['travels', travelId]
export const USER_TRAVELS_QUERY_KEY = (userId: string) => ['userTravels', userId]

/**
 * 旅行情報一覧を取得するフック
 */
export const useGetTravels = () => {
  return useQuery({
    queryKey: TRAVELS_QUERY_KEY,
    queryFn: async (): Promise<TravelsResponse> => {
      const response = await apiClient.get('/travels')
      return response.data
    },
  })
}

/**
 * 旅行情報をIDで取得するフック
 * @param travelId 旅行ID
 */
export const useGetTravelById = (travelId: string) => {
  return useQuery({
    queryKey: TRAVEL_QUERY_KEY(travelId),
    queryFn: async (): Promise<Travel> => {
      const response = await apiClient.get(`/travels/${travelId}`)
      return response.data
    },
    // 旅行IDが空の場合はクエリを無効化
    enabled: !!travelId,
  })
}

/**
 * ユーザーIDに紐づく旅行情報一覧を取得するフック
 * @param userId ユーザーID
 */
export const useGetTravelsByUserId = (userId: string) => {
  return useQuery({
    queryKey: USER_TRAVELS_QUERY_KEY(userId),
    queryFn: async (): Promise<TravelsResponse> => {
      const response = await apiClient.get(`/travels`, {
        params: { userId }
      })
      return response.data
    },
    // ユーザーIDが空の場合はクエリを無効化
    enabled: !!userId,
  })
}

/**
 * 新規旅行情報を作成するフック
 */
export const useCreateTravel = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (travelData: TravelInput): Promise<Travel> => {
      console.log('旅行データ:', travelData)
      // すべてのフィールドが文字列であることを確認
      const validatedData = {
        ...travelData,
        userId: String(travelData.userId),
        title: String(travelData.title || '無題の旅行'),
        description: String(travelData.description || '旅の思い出'),
        startDate: String(travelData.startDate || new Date().toISOString().split('T')[0]),
        endDate: String(travelData.endDate || new Date().toISOString().split('T')[0]),
        sharedId: String(travelData.sharedId || `share_${Date.now()}`),
        thumbnail: String(travelData.thumbnail || '/placeholder.webp'),
      }

      const response = await apiClient.post('/travels', validatedData)
      return response.data
    },
    // ミューテーション成功後に旅行情報一覧のキャッシュを無効化
    onSuccess: (data: Travel) => {
      queryClient.invalidateQueries({ queryKey: TRAVELS_QUERY_KEY })
      queryClient.invalidateQueries({ queryKey: USER_TRAVELS_QUERY_KEY(data.userId) })
    },
  })
}

/**
 * 旅行情報を更新するフック
 */
export const useUpdateTravel = (travelId: string) => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (travelData: Partial<TravelInput>): Promise<Travel> => {
      const response = await apiClient.put(`/travels/${travelId}`, travelData)
      return response.data
    },
    // ミューテーション成功後に関連するキャッシュを無効化
    onSuccess: (data: Travel) => {
      queryClient.invalidateQueries({ queryKey: TRAVEL_QUERY_KEY(travelId) })
      queryClient.invalidateQueries({ queryKey: USER_TRAVELS_QUERY_KEY(data.userId) })
      queryClient.invalidateQueries({ queryKey: TRAVELS_QUERY_KEY })
    },
  })
}

/**
 * 旅行情報を削除するフック
 */
export const useDeleteTravel = (travelId: string, userId: string) => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (version: number): Promise<void> => {
      await apiClient.delete(`/travels/${travelId}`, {
        params: { version },
      })
    },
    // ミューテーション成功後に関連するキャッシュを無効化
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: TRAVEL_QUERY_KEY(travelId) })
      queryClient.invalidateQueries({ queryKey: USER_TRAVELS_QUERY_KEY(userId) })
      queryClient.invalidateQueries({ queryKey: TRAVELS_QUERY_KEY })
    },
  })
}
