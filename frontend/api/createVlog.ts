import { useMutation, useQueryClient } from '@tanstack/react-query'
import {apiClient} from './client'
import { VLOGS_QUERY_KEY } from './vlogAPi'

// VLog作成レスポンス型定義
export type CreateVlogResponse = {
  vlogId: string
  status: string
}

export const useCreateVlog = () => {
  const queryClient = useQueryClient()

  return useMutation<CreateVlogResponse>({
    mutationFn: async (): Promise<CreateVlogResponse> => {
      // 実際には FormData を送信するはずだが、既存のコードに合わせて post() を使用
      const res = await apiClient.post('/agent/create-vlog')
      return res.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: VLOGS_QUERY_KEY })
    },
    onError: error => {
      console.error('Error creating vlog:', error)
    },
  })
}
