import {apiClient} from './client'
import {
  MediaUploadRequest,
  MediaUploadResponse,
  MediaImageUploadRequest,
  MediaImageUploadResponse,
  MediaVideoUploadRequest,
  MediaVideoUploadResponse,
  MediaGetResponse,
  MediaListResponse,
  MediaAnalysisBatchResponse,
} from './types'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'

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
} as const

/**
 * メディア（画像・動画）アップロードのフック（統合版）
 */
export const useUploadMedia = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (request: MediaUploadRequest) => uploadMedia(request),
    onSuccess: data => {
      // アップロード成功時にキャッシュを無効化して最新データを取得
      queryClient.invalidateQueries({ queryKey: MEDIA_QUERY_KEYS.images() })

      // アップロードしたメディアを個別にキャッシュ
      queryClient.setQueryData(MEDIA_QUERY_KEYS.image(data.file_id), {
        file_id: data.file_id,
        url: data.url,
      })
    },
    onError: error => {
      console.error('メディアアップロードエラー:', error)
    },
  })
}

/**
 * 画像アップロードのフック（後方互換性のため保持）
 */
export const useUploadMediaImage = () => {
  return useUploadMedia()
}

/**
 * 動画アップロードのフック（後方互換性のため保持）
 */
export const useUploadMediaVideo = () => {
  return useUploadMedia()
}

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
 * メディアをアップロードし、同時にAI分析を実行するフック
 */
export const useAnalyzeMedia = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (files: File[]): Promise<MediaAnalysisBatchResponse> => {
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
