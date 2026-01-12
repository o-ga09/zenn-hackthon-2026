'use client'

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import apiClient from './client'

// VLog 型定義
export interface Vlog {
  id: string
  video_id: string
  video_url: string
  share_url: string
  duration: number
  thumbnail: string
  created_at: string
  updated_at: string
}

export interface VlogsResponse {
  items: Vlog[]
  total: number
  message?: string
}

// キャッシュキー
export const VLOGS_QUERY_KEY = ['vlogs']
export const VLOG_QUERY_KEY = (id: string) => ['vlogs', id]

/** VLog一覧取得 */
export const useGetVlogs = () => {
  return useQuery({
    queryKey: VLOGS_QUERY_KEY,
    queryFn: async (): Promise<VlogsResponse> => {
      const res = await apiClient.get('/vlogs')
      return res.data
    },
    staleTime: 1000 * 60 * 2,
  })
}

/** VLogをIDで取得 */
export const useGetVlogById = (vlogId: string) => {
  return useQuery({
    queryKey: VLOG_QUERY_KEY(vlogId),
    queryFn: async (): Promise<Vlog> => {
      const res = await apiClient.get(`/vlogs/${vlogId}`)
      return res.data
    },
    enabled: !!vlogId,
    staleTime: 1000 * 60 * 5,
  })
}

/** VLog削除 */
export const useDeleteVlog = (vlogId?: string) => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (): Promise<void> => {
      if (!vlogId) throw new Error('vlogId is required')
      await apiClient.delete(`/vlogs/${vlogId}`)
    },
    onSuccess: () => {
      if (vlogId) {
        queryClient.removeQueries({ queryKey: VLOG_QUERY_KEY(vlogId) })
      }
      queryClient.invalidateQueries({ queryKey: VLOGS_QUERY_KEY })
    },
    onError: error => {
      console.error('VLog削除エラー:', error)
    },
  })
}
