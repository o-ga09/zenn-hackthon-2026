import { apiClient } from './client'
import {
  MediaUploadRequest,
  MediaUploadResponse,
  MediaImageUploadRequest,
  MediaImageUploadResponse,
  MediaVideoUploadRequest,
  MediaVideoUploadResponse,
  MediaGetResponse,
  MediaListResponse,
  Media,
  AnalyzeMediaResponse,
  MediaStatusResponse,
  MediaAnalyticsResponse,
  UpdateMediaAnalyticsRequest,
} from './types'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useEffect, useState } from 'react'
import { toast } from 'sonner'

/**
 * メディア（画像・動画）をアップロードするAPI関数（統合版）
 */
export const uploadMedia = async (request: MediaUploadRequest): Promise<MediaUploadResponse> => {
  const formData = new FormData()
  formData.append('file', request.file)

  const response = await apiClient.post('/media', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  })
  return response.data
}

/**
 * 画像をアップロードするAPI関数（後方互換性のため保持）
 */
export const uploadMediaImage = async (
  request: MediaImageUploadRequest
): Promise<MediaImageUploadResponse> => {
  return uploadMedia(request)
}

/**
 * 動画をアップロードするAPI関数（後方互換性のため保持）
 */
export const uploadMediaVideo = async (
  request: MediaVideoUploadRequest
): Promise<MediaVideoUploadResponse> => {
  return uploadMedia(request)
}

/**
 * メディア一覧を取得するAPI関数
 */
export const getMediaList = async (): Promise<MediaListResponse> => {
  const response = await apiClient.get('/media')
  return response.data
}

/**
 * 画像を取得するAPI関数
 */
export const getMediaImage = async (key: string): Promise<MediaGetResponse> => {
  const response = await apiClient.get(`/media/${key}`)
  return response.data
}

/**
 * 画像を削除するAPI関数
 */
export const deleteMediaImage = async (key: string): Promise<void> => {
  await apiClient.delete(`/media/${key}`)
}

// キャッシュのキー定数
export const MEDIA_QUERY_KEYS = {
  image: (key: string) => ['media', 'image', key],
  images: () => ['media', 'images'],
  analytics: (fileId: string) => ['media', 'analytics', fileId],
} as const

/**
 * 画像取得のフック
 */
export const useGetMediaImage = (key: string, enabled: boolean = true) => {
  return useQuery({
    queryKey: MEDIA_QUERY_KEYS.image(key),
    queryFn: () => getMediaImage(key),
    enabled: !!key && enabled,
    staleTime: 1000 * 60 * 5, // 5分間はキャッシュを使用
    gcTime: 1000 * 60 * 10, // 10分後にガベージコレクション
  })
}

/**
 * 画像削除のフック
 */
export const useDeleteMediaImage = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (key: string) => deleteMediaImage(key),
    onSuccess: (_, key) => {
      // 削除成功時に該当画像のキャッシュを削除
      queryClient.removeQueries({ queryKey: MEDIA_QUERY_KEYS.image(key) })

      // 画像一覧のキャッシュも無効化
      queryClient.invalidateQueries({ queryKey: MEDIA_QUERY_KEYS.images() })
    },
    onError: error => {
      console.error('画像削除エラー:', error)
    },
  })
}

/**
 * メディア一覧取得のフック
 */
export const useGetMediaList = () => {
  return useQuery({
    queryKey: MEDIA_QUERY_KEYS.images(),
    queryFn: getMediaList,
    staleTime: 1000 * 60 * 2, // 2分間はキャッシュを使用
    gcTime: 1000 * 60 * 5, // 5分後にガベージコレクション
  })
}

/**
 * メディアをアップロードし、同時にAI分析を実行するフック（非同期SSE対応）
 */
export const useAnalyzeMedia = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (files: File[]): Promise<AnalyzeMediaResponse> => {
      const formData = new FormData()
      files.forEach(file => {
        formData.append('files', file)
      })

      const response = await apiClient.post('/agent/analyze-media', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      })
      return response.data
    },
    onSuccess: () => {
      // 分析完了後にメディア一覧キャッシュを無効化
      queryClient.invalidateQueries({ queryKey: MEDIA_QUERY_KEYS.images() })
    },
    onError: error => {
      console.error('メディア分析エラー:', error)
    },
  })
}

/**
 * メディア分析の進捗をSSEで監視するフック
 */
export const useMediaAnalysisSSE = (mediaIds: string[] | null) => {
  const [status, setStatus] = useState<MediaStatusResponse | null>(null)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!mediaIds || mediaIds.length === 0) return

    const baseUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api'
    const sseUrl = `${baseUrl}/agent/analyze-media/stream?ids=${mediaIds.join(',')}`

    const eventSource = new EventSource(sseUrl, {
      withCredentials: true,
    })

    eventSource.onmessage = event => {
      try {
        const data = JSON.parse(event.data) as MediaStatusResponse
        setStatus(data)

        if (data.all_completed) {
          if (data.failed_items > 0) {
            setError(`${data.failed_items}件のメディアの分析に失敗しました`)
          }
          eventSource.close()
        }
      } catch (err) {
        console.error('SSEデータのパースエラー:', err)
      }
    }

    eventSource.onerror = err => {
      console.error('SSE接続エラー:', err)
      eventSource.close()

      // フォールバック: ポーリング
      startAnalysisPolling(mediaIds, setStatus, setError)
    }

    return () => {
      eventSource.close()
    }
  }, [mediaIds?.join(',')])

  return { status, error }
}

/**
 * SSE接続失敗時のポーリングフォールバック（分析用）
 */
const startAnalysisPolling = (
  mediaIds: string[],
  setStatus: (status: MediaStatusResponse) => void,
  setError: (error: string | null) => void
) => {
  const pollInterval = setInterval(async () => {
    try {
      // 各メディアの状態を取得
      const medias: Media[] = []
      let completedCount = 0
      let failedCount = 0

      for (const id of mediaIds) {
        const response = await apiClient.get(`/media/${id}`)
        const media = response.data as Media
        medias.push(media)

        if (media.status === 'completed') completedCount++
        if (media.status === 'failed') failedCount++
      }

      const allCompleted = completedCount + failedCount === mediaIds.length

      const data: MediaStatusResponse = {
        medias,
        total_items: mediaIds.length,
        completed_items: completedCount,
        failed_items: failedCount,
        all_completed: allCompleted,
      }
      setStatus(data)

      if (allCompleted) {
        if (failedCount > 0) {
          setError(`${failedCount}件のメディアの分析に失敗しました`)
        }
        clearInterval(pollInterval)
      }
    } catch (err) {
      console.error('ポーリングエラー:', err)
      setError('進捗の取得に失敗しました')
      clearInterval(pollInterval)
    }
  }, 2000)

  // 5分後にタイムアウト
  setTimeout(
    () => {
      clearInterval(pollInterval)
    },
    5 * 60 * 1000
  )
}

/**
 * メディアの分析結果を取得するAPI関数
 */
export const getMediaAnalytics = async (fileId: string): Promise<MediaAnalyticsResponse> => {
  const response = await apiClient.get(`/media/${fileId}/analytics`)
  return response.data
}

/**
 * メディアの分析結果を更新するAPI関数
 */
export const updateMediaAnalytics = async (
  fileId: string,
  request: UpdateMediaAnalyticsRequest
): Promise<MediaAnalyticsResponse> => {
  const response = await apiClient.put(`/media/${fileId}/analytics`, request)
  return response.data
}

/**
 * メディア分析結果取得のフック
 */
export const useGetMediaAnalytics = (fileId: string, enabled: boolean = true) => {
  return useQuery({
    queryKey: MEDIA_QUERY_KEYS.analytics(fileId),
    queryFn: () => getMediaAnalytics(fileId),
    enabled: !!fileId && enabled,
    staleTime: 1000 * 60 * 5, // 5分間はキャッシュを使用
    gcTime: 1000 * 60 * 10, // 10分後にガベージコレクション
  })
}

/**
 * メディア分析結果更新のフック（楽観的UI更新パターン）
 */
export const useUpdateMediaAnalytics = (fileId: string) => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (request: UpdateMediaAnalyticsRequest) => updateMediaAnalytics(fileId, request),

    // 楽観的更新: API呼び出し前にローカル状態を更新
    onMutate: async (newData: UpdateMediaAnalyticsRequest) => {
      // 進行中のクエリをキャンセル
      await queryClient.cancelQueries({ queryKey: MEDIA_QUERY_KEYS.analytics(fileId) })

      // 現在のキャッシュデータを取得（ロールバック用）
      const previousAnalytics = queryClient.getQueryData<MediaAnalyticsResponse>(
        MEDIA_QUERY_KEYS.analytics(fileId)
      )

      // 楽観的にキャッシュを更新
      if (previousAnalytics) {
        queryClient.setQueryData<MediaAnalyticsResponse>(MEDIA_QUERY_KEYS.analytics(fileId), {
          ...previousAnalytics,
          ...(newData.description !== undefined && { description: newData.description }),
          ...(newData.mood !== undefined && { mood: newData.mood }),
          ...(newData.objects !== undefined && { objects: newData.objects }),
          ...(newData.landmarks !== undefined && { landmarks: newData.landmarks }),
          ...(newData.activities !== undefined && { activities: newData.activities }),
        })
      }

      // ロールバック用のデータを返す
      return { previousAnalytics }
    },

    // 成功時: トースト表示とキャッシュ無効化
    onSuccess: () => {
      toast.success('タグを更新しました')
      queryClient.invalidateQueries({ queryKey: MEDIA_QUERY_KEYS.analytics(fileId) })
    },

    // エラー時: ロールバック
    onError: (_err, _newData, context) => {
      if (context?.previousAnalytics) {
        queryClient.setQueryData(MEDIA_QUERY_KEYS.analytics(fileId), context.previousAnalytics)
      }
      toast.error('タグの更新に失敗しました')
    },

    // 完了時: 念のため再取得
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: MEDIA_QUERY_KEYS.analytics(fileId) })
    },
  })
}
