import { useMutation, useQueryClient } from '@tanstack/react-query'
import {apiClient} from './client'
import { VLOGS_QUERY_KEY } from './vlogAPi'

// サンプルレスポンス型定義
type Recipe = {
  title: string
  description: string
  prepTime: string
  cookTime: string
  servings: number
  ingredients: string[]
  instructions: string[]
  tips?: string[]
}

export const useCreateVlog = () => {
  const queryClient = useQueryClient()

  return useMutation<Recipe>({
    mutationFn: async (): Promise<Recipe> => {
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
