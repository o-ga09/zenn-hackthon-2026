'use client'
import React, { createContext, useContext, useState, useEffect } from 'react'
import { signOut, onAuthStateChanged, GoogleAuthProvider, signInWithPopup } from 'firebase/auth'
import { User } from '@/types/types'
import axios from 'axios'
import { useRouter, usePathname } from 'next/navigation'
import { auth } from '@/lib/firebase'
import { encodeUserDataForCookie } from '@/lib/utils'
import CSRFManager from '@/lib/utils'
import apiClient, { setCurrentUserID } from '@/api/client'

interface AuthContextType {
  user: User | null
  loading: boolean
  setUser: React.Dispatch<React.SetStateAction<User | null>>
  logout: () => Promise<void>
  setAuthInProgress: React.Dispatch<React.SetStateAction<boolean>>
  getToken: (forceRefresh?: boolean) => Promise<string | null>
  currentUser: any | null
  googleLogin: () => Promise<void>
}

const AuthContext = createContext<AuthContextType>({
  user: null,
  loading: true,
  setUser: () => {},
  logout: async () => {},
  setAuthInProgress: () => {},
  getToken: async () => null,
  currentUser: null,
  googleLogin: async () => {},
})

export const useAuth = () => useContext(AuthContext)

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)
  const [authInProgress, setAuthInProgress] = useState(false)

  const [firebaseId, setFirebaseId] = useState<string>('')
  const [cachedToken, setCachedToken] = useState<string | null>(null)
  const [tokenExpiryTime, setTokenExpiryTime] = useState<number>(0)
  const [tokenCreatedTime, setTokenCreatedTime] = useState<number>(0)
  const router = useRouter()
  const pathname = usePathname()

  // APIユーザー状態
  const [apiUser, setApiUser] = useState<User | null>(null)
  const [isUserLoading, setIsUserLoading] = useState(false)
  const [apiError, setApiError] = useState<any>(null)

  // FirebaseIDからユーザー情報を取得する関数
  const fetchUserByFirebaseId = async (fbId: string) => {
    if (!fbId) return

    setIsUserLoading(true)
    setApiError(null)

    try {
      const response = await apiClient.get(`/firebaseUsers/${fbId}`)
      console.log('Fetched API user data:', response.data)
      const userFromApi = {
        id: response.data.user.id,
        userID: response.data.user.userId,
        version: response.data.user.version,
        uid: fbId,
        username: response.data.user.username,
        displayname: response.data.user.displayname,
        photoURL: response.data.user.photoURL,
        email: response.data.user.email,
        birthDay: response.data.user.birthDay,
      }
      setApiUser(userFromApi)

      // APIから取得したユーザーIDを設定
      if (response.data.user.userId) {
        setCurrentUserID(response.data.user.userId)
      }

      return response.data
    } catch (error) {
      setApiError(error)
      return null
    } finally {
      setIsUserLoading(false)
    }
  }

  // firebaseIdが変更されたらユーザー情報を取得
  useEffect(() => {
    if (firebaseId) {
      fetchUserByFirebaseId(firebaseId)
    }
  }, [firebaseId])

  // 公開ルートかどうかを判定するヘルパー関数
  const isPublicRoute = (path: string): boolean => {
    const publicPaths = [
      '/',
      '/privacy',
      '/terms',
      '/help',
      '/about',
      '/features',
      '/pricing',
      '/feedback',
      '/legal',
      '/register',
      '/login',
      '/sign-up',
      '/reset-password',
      '/unsubscribe',
    ]
    return publicPaths.includes(path)
  }

  // Firebase認証状態の監視
  useEffect(() => {
    const unsubscribe = onAuthStateChanged(auth, async firebaseUser => {
      if (!authInProgress) {
        setLoading(true)
      }

      try {
        if (firebaseUser) {
          setFirebaseId(firebaseUser.uid)

          // Firebase認証完了時点で基本的なユーザー情報をクッキーに保存
          // これによりmiddlewareで認証状態を判定できる
          const basicUserData = {
            id: firebaseUser.uid, // 一時的にFirebaseUIDを使用
            uid: firebaseUser.uid,
            username: firebaseUser.displayName || '',
            displayname: firebaseUser.displayName || '',
            photoURL: firebaseUser.photoURL || '',
            email: firebaseUser.email || '',
          }

          if (typeof document !== 'undefined') {
            const expires = new Date(Date.now() + 24 * 60 * 60 * 1000).toUTCString() // 24時間
            const encodedUserData = encodeUserDataForCookie(basicUserData)
            document.cookie = `tavinikkiy-user=${encodedUserData}; expires=${expires}; path=/; secure; samesite=strict`
          }

          // ログインページ、登録ページ、サインアップページでは特別な処理をしない
          if (!['/login', '/register', '/sign-up'].includes(pathname || '')) {
            // トークンを取得（APIリクエスト用）
            await getToken(false)
          }
        } else {
          // ユーザーがログアウトした場合
          setFirebaseId('')
          setUser(null)
          localStorage.removeItem('tavinikkiy-user')

          // クッキーもクリア
          if (typeof document !== 'undefined') {
            document.cookie =
              'tavinikkiy-user=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/; secure; samesite=strict'
          }

          // 公開ルートでない場合のみホームにリダイレクト
          if (pathname && !isPublicRoute(pathname)) {
            router.push('/')
          }
        }
      } catch (error) {
        console.error('Error updating auth state:', error)
      } finally {
        if (!authInProgress) {
          setLoading(false)
        }
      }
    })

    return () => unsubscribe()
  }, [authInProgress, pathname, router])

  // APIユーザーデータの変更を監視
  useEffect(() => {
    // 公開ルートでは実行しない
    if (!pathname || isPublicRoute(pathname)) {
      return
    }

    if (!isUserLoading && firebaseId) {
      const firebaseUser = auth.currentUser
      if (firebaseUser) {
        // APIエラーをチェック
        if (apiError) {
          // @ts-ignore
          const statusCode = apiError.response?.status

          // 404エラーの場合（ユーザーが存在しない）、サインアップページにリダイレクト
          if (statusCode === 404) {
            router.push('/sign-up')
          }
        } else if (apiUser) {
          console.log('API User data fetched:', apiUser)
          // プロフィール画像の取得ロジックを修正
          // APIから画像データが取得できる場合はそれを優先し、取得できない場合はFirebaseのプロフィール画像を使用
          let profileImage = ''

          // 1. APIからのアップロードされた画像データを最優先（Base64形式）
          if (apiUser.photoURL && apiUser.photoURL.length > 0) {
            profileImage = apiUser.photoURL
          }
          // 2. Firebase Authからのプロフィール画像を使用（ない場合はデフォルト画像が使用される）
          else if (firebaseUser.photoURL) {
            profileImage = firebaseUser.photoURL
          }

          // 更新されたユーザー情報を作成（内部型に変換）
          const updatedUser: User = {
            id: apiUser.id,
            userID: apiUser.userID,
            version: apiUser.version || 1,
            uid: firebaseId,
            username: apiUser.username || firebaseUser.displayName || '',
            displayname: apiUser.displayname || firebaseUser.displayName || '',
            photoURL: profileImage,
            email: firebaseUser.email,
            birthDay: apiUser.birthDay,
            gender: apiUser.gender,
          }

          // ローカルストレージと状態を更新
          console.log('Setting user data in auth context:', updatedUser)
          setUser(updatedUser)

          // middlewareとの連携のためCookieにもユーザー情報を保存（詳細版で更新）
          if (typeof document !== 'undefined') {
            const expires = new Date(Date.now() + 24 * 60 * 60 * 1000).toUTCString() // 24時間
            // UTF-8対応のBase64エンコードでクッキーに安全に保存
            const encodedUserData = encodeUserDataForCookie(updatedUser)
            document.cookie = `tavinikkiy-user=${encodedUserData}; expires=${expires}; path=/; secure; samesite=strict`
          }
        }
      }
    }
  }, [apiUser, apiError, isUserLoading, firebaseId, pathname, router])

  // ログアウト処理
  const logout = async () => {
    try {
      // 先にクリーンアップを行う
      delete axios.defaults.headers.common['Authorization']
      delete axios.defaults.headers.common['x-tavinikkiy-user']
      setUser(null)
      setFirebaseId('')
      setCachedToken(null)
      setTokenExpiryTime(0)
      setTokenCreatedTime(0)

      // currentUserIDをリセット
      setCurrentUserID('')

      // ローカルストレージをクリア
      localStorage.removeItem('tavinikkiy-user')

      // クッキーをクリア
      if (typeof document !== 'undefined') {
        document.cookie =
          'tavinikkiy-user=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/; secure; samesite=strict'
      }

      // CSRFトークンもクリア
      CSRFManager.getInstance().clearToken()

      // セッションクリアはFirebase SDKが自動処理

      // 最後にFirebaseからのログアウトを実行
      await signOut(auth)

      // ログアウト後にホームページにリダイレクト
      router.push('/')
    } catch (error) {
      console.error('Logout error:', error)
    }
  }

  // バックグラウンドでのトークン更新
  const refreshTokenInBackground = async () => {
    try {
      const currentUser = auth.currentUser
      if (!currentUser) return

      const newToken = await currentUser.getIdToken(true)
      if (newToken) {
        const currentTime = Date.now()
        setCachedToken(newToken)
        setTokenExpiryTime(currentTime + 30 * 60 * 1000) // 30分
        setTokenCreatedTime(currentTime)
      }
    } catch (error) {
      // バックグラウンド更新の失敗は静かに処理
      console.warn('Background token refresh failed:', error)
    }
  }

  // トークンを取得するメソッド - キャッシュを活用
  const getToken = async (forceRefresh: boolean = false): Promise<string | null> => {
    try {
      const currentUser = auth.currentUser
      if (!currentUser) {
        console.warn('Auth context - No current user found when getting token')
        return null
      }

      const currentTime = Date.now()

      // 1. キャッシュされたトークンがある AND 強制更新でない AND トークンが有効期限内
      if (cachedToken && !forceRefresh && currentTime < tokenExpiryTime) {
        // バックグラウンドで更新が必要かチェック（25分後）
        if (currentTime > tokenCreatedTime + 25 * 60 * 1000) {
          // 非同期でバックグラウンド更新（UXに影響しない）
          refreshTokenInBackground()
        }
        return cachedToken
      }

      // 2. トークンの更新が必要な場合
      try {
        const token = await currentUser.getIdToken(forceRefresh)
        if (token) {
          // トークンをキャッシュし、有効期限を設定（30分）
          setCachedToken(token)
          setTokenExpiryTime(currentTime + 30 * 60 * 1000)
          setTokenCreatedTime(currentTime)

          // セッション状態の管理はFirebase SDKに委任
          // HttpOnlyクッキーは使用しない

          return token
        }
        return null
      } catch (error: any) {
        console.error('Error getting fresh token:', error)

        // クォータ超過エラーの場合の処理
        if (error?.code === 'auth/quota-exceeded') {
          console.warn('Firebase quota exceeded, checking cached token')
          // キャッシュされたトークンがまだ有効な場合はそれを使用
          if (cachedToken && currentTime < tokenExpiryTime) {
            console.info('Using valid cached token')
            return cachedToken
          }
          // キャッシュが無効な場合は、短い待機時間後に再試行
          await new Promise(resolve => setTimeout(resolve, 1000))
          return getToken(false) // 再帰的に呼び出し（forceRefresh=false）
        }

        throw error // その他のエラーは上位に投げる
      }
    } catch (error) {
      console.error('Auth context - Error getting token:', error)
      // キャッシュされたトークンがある場合はそれを返す（エラー回復メカニズム）
      if (cachedToken && Date.now() < tokenExpiryTime) {
        return cachedToken
      }
      return null
    }
  }

  // Googleアカウントでのログイン
  const googleLogin = async (): Promise<void> => {
    try {
      setAuthInProgress(true)
      const provider = new GoogleAuthProvider()

      // Googleリダイレクト認証を開始
      const result = await signInWithPopup(auth, provider)
      const user = result.user

      setFirebaseId(user.uid)
      // 少し待機してAPIリクエストが完了するのを待つ
      await new Promise(resolve => setTimeout(resolve, 1500))

      // apiUserとapiErrorの最新の状態を取得
      const error = apiError

      if (apiUser) {
        // ユーザーが既に存在する場合は/dashboardにリダイレクト
        router.push('/dashboard')
      } else if (error) {
        // @ts-ignore
        const statusCode = error.response?.status

        // 404エラーの場合は新規ユーザーと判断
        if (statusCode === 404) {
          router.push('/sign-up')
        } else {
          // その他のエラーは上位に投げる
          throw error
        }
      } else {
        // データもエラーもない場合は、APIリクエストがまだ完了していない可能性がある
        // sign-upページに進む（APIリクエストが完了し、ユーザーが存在しない場合はそのままになる）
        router.push('/sign-up')
      }

      // 注意: ここで処理は中断され、リダイレクト後に続行されます
      // リダイレクト結果は、onAuthStateChangedで自動的に処理されます
    } catch (error) {
      console.error('Google login error:', error)
      setAuthInProgress(false)
      throw error
    }
  }

  return (
    <AuthContext.Provider
      value={{
        user,
        loading: loading || isUserLoading,
        setUser,
        logout,
        setAuthInProgress,
        getToken,
        currentUser: auth.currentUser,
        googleLogin,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}
