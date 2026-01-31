'use client'

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import {apiClient} from './client'

// VLog 型定義
export interface Vlog {
  id: string
  video_id: string
  video_url: string
  share_url: string
  duration: number
  thumbnail: string
  status: 'pending' | 'processing' | 'completed' | 'failed'
  error_message?: string
  progress: number
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

/** 
 * VLog作成の進捗をSSEで監視するフック
 * 接続が切れた場合は自動的にポーリングにフォールバックする
 */
import { useEffect, useState } from 'react'

export const useVlogSSE = (vlogId: string | null) => {
  const [status, setStatus] = useState<Vlog | null>(null)
  const [error, setError] = useState<Error | null>(null)
  const queryClient = useQueryClient()

  useEffect(() => {
    if (!vlogId) return

    let pollInterval: NodeJS.Timeout | null = null

    const startPolling = () => {
      if (pollInterval) return
      console.log('Falling back to polling for vlog:', vlogId)
      pollInterval = setInterval(async () => {
        try {
          const res = await apiClient.get(`/vlogs/${vlogId}`)
          const data = res.data as Vlog
          setStatus(data)
          if (data.status === 'completed' || data.status === 'failed') {
            if (pollInterval) clearInterval(pollInterval)
            queryClient.invalidateQueries({ queryKey: VLOGS_QUERY_KEY })
          }
        } catch (err) {
          console.error('Polling error:', err)
          if (pollInterval) clearInterval(pollInterval)
        }
      }, 3000)
    }

    const sseUrl = `${process.env.NEXT_PUBLIC_API_URL || ''}/api/vlogs/${vlogId}/stream`
    const eventSource = new EventSource(sseUrl, { withCredentials: true })

    eventSource.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data) as Vlog
        setStatus(data)
        if (data.status === 'completed' || data.status === 'failed') {
          eventSource.close()
          queryClient.invalidateQueries({ queryKey: VLOGS_QUERY_KEY })
        }
      } catch (err) {
        console.error('Failed to parse SSE message:', err)
      }
    }

    eventSource.onerror = (err) => {
      console.error('SSE connection error:', err)
      eventSource.close()
      // SSEが失敗したらポーリングに切り替え
      startPolling()
    }

    return () => {
      eventSource.close()
      if (pollInterval) clearInterval(pollInterval)
    }
  }, [vlogId, queryClient])

  return { status, error }
}
