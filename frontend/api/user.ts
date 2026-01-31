import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import {apiClient} from './client'
import z from 'zod'

// キャッシュのキー
export const USERS_QUERY_KEY = ['users']
export const USER_QUERY_KEY = (userId: string) => ['users', userId]
export const FIREBASE_USER_QUERY_KEY = (firebaseId: string) => ['firebaseUsers', firebaseId]
export const USER_PHOTO_COUNT_QUERY_KEY = (userId: string) => ['users', userId, 'upload']
/**
 * ユーザー一覧を取得するフック
 */
export const useGetUsers = () => {
  return useQuery({
    queryKey: USERS_QUERY_KEY,
    queryFn: async (): Promise<User[]> => {
      const response = await apiClient.get('/users')
      return response.data
    },
  })
}

/**
 * ユーザーをIDで取得するフック
 * @param userId ユーザーID
 */
export const useGetUserById = (userId: string) => {
  return useQuery({
    queryKey: USER_QUERY_KEY(userId),
    queryFn: async (): Promise<User> => {
      const response = await apiClient.get(`/users/${userId}`)
      return response.data
    },
    // ユーザーIDが空の場合はクエリを無効化
    enabled: !!userId,
  })
}

/**
 * FirebaseのIDでユーザーを取得するフック
 * @param firebaseId FirebaseのユーザーID
 */
export const useGetUserByFirebaseId = (firebaseId: string) => {
  return useQuery({
    queryKey: FIREBASE_USER_QUERY_KEY(firebaseId),
    queryFn: async (): Promise<User> => {
      const response = await apiClient.get(`/firebaseUsers/${firebaseId}`)
      return response.data
    },
    // FirebaseIDが空の場合はクエリを無効化
    enabled: !!firebaseId,
  })
}

/**
 * 新規ユーザーを作成するフック
 */
export const useCreateUser = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (userData: User): Promise<User> => {
      // フロントエンドの形式をAPIの形式に変換
      const apiUserData =
        'firebase_id' in userData
          ? {
              uid: userData.firebase_id,
              id: userData.id, // フォームのname = ID
              name: userData.name, // nameフィールドも送信
              displayName: userData.displayName,
              ...(userData.profileImage && { profileImage: userData.profileImage }),
              ...(userData.birthday && { birthday: userData.birthday }),
              ...(userData.gender && { gender: userData.gender }),
            }
          : userData

      const response = await apiClient.post('/users', apiUserData)
      return response.data
    },
    // ミューテーション成功後にユーザー一覧のキャッシュを無効化
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: USERS_QUERY_KEY })
    },
  })
}

/**
 * ユーザーを更新するフック
 */
export const useUpdateUser = (userId: string) => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (userData: User): Promise<User> => {
      const response = await apiClient.put(`/users/${userId}`, userData)
      return response.data
    },
    // ミューテーション成功後に関連するキャッシュを無効化
    onSuccess: (data: User) => {
      queryClient.invalidateQueries({ queryKey: USER_QUERY_KEY(userId) })
      queryClient.invalidateQueries({ queryKey: FIREBASE_USER_QUERY_KEY(data.uid) })
      queryClient.invalidateQueries({ queryKey: USERS_QUERY_KEY })
    },
  })
}

/**
 * ユーザーを削除するフック
 */
export const useDeleteUser = (userId: string) => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (version: number): Promise<void> => {
      await apiClient.delete(`/users/${userId}`, {
        params: { version },
      })
    },
    // ミューテーション成功後に関連するキャッシュを無効化
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: USER_QUERY_KEY(userId) })
      queryClient.invalidateQueries({ queryKey: USERS_QUERY_KEY })
    },
  })
}

/**
 * ユーザーがアップロードした写真の総数を取得するフック
 * @param userId ユーザーID
 */
export const useGetUserPhotoCount = (userId: string) => {
  return useQuery({
    queryKey: USER_PHOTO_COUNT_QUERY_KEY(userId),
    queryFn: async (): Promise<{ videoCount: number; uploadCount: number }> => {
      const response = await apiClient.get(`/users/${userId}/upload`)
      return response.data
    },
    // ユーザーIDが空の場合はクエリを無効化
    enabled: !!userId,
  })
}

export type User = {
  id: string
  version: number
  uid: string
  type: string
  tokenBalance: number
  name?: string
  displayName?: string
  bio?: string
  plan: string
  profileImage?: string
  birthday?: string
  gender?: string
  followersCount: number
  followingCount: number
  isPublic?: boolean
  created_at: string
  updated_at: string
}
// バリデーションスキーマ
export const profileFormSchema = z.object({
  displayName: z
    .string()
    .min(1, 'ユーザー名は必須です')
    .max(50, 'ユーザー名は50文字以内で入力してください'),
  isPublic: z.boolean().optional(),
  profileImage: z.string().optional(),
})

export type ProfileFormData = z.infer<typeof profileFormSchema>
