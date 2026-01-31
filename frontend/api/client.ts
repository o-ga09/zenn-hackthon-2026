import axios from 'axios'
import { getAuth } from 'firebase/auth'

const baseURL = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080'

// ユーザーIDを格納するための変数
let currentUserID: string | null = null

// APIからのユーザーID設定用の関数
export const setCurrentUserID = (userID: string) => {
  currentUserID = userID
}

// Axiosインスタンスを作成
// 認証が必要なAPI
const apiClient = axios.create({
  baseURL: `${baseURL}/api`,
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: true, // Cookieを送受信するために必要
})

// 認証不要なAPI
const noCredentialApiClient = axios.create({
  baseURL: `${baseURL}/api`,
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

export {
  apiClient,
  noCredentialApiClient,
}
