import axios from 'axios'
import { getAuth } from 'firebase/auth'

// ユーザーIDを格納するための変数
let currentUserID: string | null = null

// APIからのユーザーID設定用の関数
export const setCurrentUserID = (userID: string) => {
  currentUserID = userID
}

// Axiosインスタンスを作成
const apiClient = axios.create({
  baseURL: '/api',
  headers: {
    'Content-Type': 'application/json',
  },
})

// リクエストインターセプターでFirebaseの認証トークンをヘッダーに追加
apiClient.interceptors.request.use(async config => {
  try {
    const auth = getAuth()
    const user = auth.currentUser

    if (user) {
      const token = await user.getIdToken()
      config.headers.Authorization = `Bearer ${token}`

      // APIから取得したユーザーIDをヘッダーに追加
      if (currentUserID) {
        config.headers['X-Tavinikkiy-User-Id'] = currentUserID
      } else {
        // ユーザーIDがまだ設定されていない場合は、FirebaseのUIDを使用
        config.headers['X-Tavinikkiy-User-Id'] = user.uid
      }
    }

    return config
  } catch (error) {
    console.error('認証トークン取得エラー:', error)
    return config
  }
})

export default apiClient
