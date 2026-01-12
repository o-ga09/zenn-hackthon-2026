import apiClient from './client'
import {
  MediaImageUploadRequest,
  MediaImageUploadResponse,
  MediaGetResponse,
  MediaListResponse,
} from './types'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'

/**
 * 画像をアップロードするAPI関数
 */
export const uploadMediaImage = async (
  request: MediaImageUploadRequest
): Promise<MediaImageUploadResponse> => {
  const response = await apiClient.post('/media', request)
  return response.data
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
 * 画像アップロードのフック
 */
export const useUploadMediaImage = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (request: MediaImageUploadRequest) => uploadMediaImage(request),
    onSuccess: data => {
      // アップロード成功時にキャッシュを無効化して最新データを取得
      queryClient.invalidateQueries({ queryKey: MEDIA_QUERY_KEYS.images() })

      // アップロードした画像を個別にキャッシュ
      queryClient.setQueryData(MEDIA_QUERY_KEYS.image(data.file_id), {
        file_id: data.file_id,
        url: data.url,
      })
    },
    onError: error => {
      console.error('画像アップロードエラー:', error)
    },
  })
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
