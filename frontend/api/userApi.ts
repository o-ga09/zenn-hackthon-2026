import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import apiClient from './client'
import { User, UserInput, UserInputFrontend, UserUpdateInput, UsersResponse } from './types'

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
    queryFn: async (): Promise<UsersResponse> => {
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
    mutationFn: async (userData: UserInput | UserInputFrontend): Promise<User> => {
      // フロントエンドの形式をAPIの形式に変換
      const apiUserData =
        'firebase_id' in userData
          ? {
              uid: userData.firebase_id,
              id: userData.name, // フォームのname = ID
              name: userData.name, // nameフィールドも送信
              displayName: userData.display_name,
              ...(userData.image_data && { image_data: userData.image_data }),
              ...(userData.birth_day && { birth_day: userData.birth_day }),
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
    mutationFn: async (userData: UserUpdateInput): Promise<User> => {
      const response = await apiClient.put(`/users/${userId}`, userData)
      return response.data
    },
    // ミューテーション成功後に関連するキャッシュを無効化
    onSuccess: (data: User) => {
      queryClient.invalidateQueries({ queryKey: USER_QUERY_KEY(userId) })
      queryClient.invalidateQueries({ queryKey: FIREBASE_USER_QUERY_KEY(data.firebase_id) })
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
