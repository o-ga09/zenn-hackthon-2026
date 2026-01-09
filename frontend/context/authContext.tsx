'use client'

import { createContext, useContext, useEffect, useState, ReactNode } from 'react'
import { signInWithPopup, GoogleAuthProvider, signOut } from 'firebase/auth'
import { auth } from '@/lib/firebase'
import { useRouter } from 'next/navigation'

const baseURL = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080'

type User = {
  id: string
  uid: string
  type: string
  name: string
  token_balance: number
  username?: string
  photoURL?: string
  displayname?: string
  created_at: string
  updated_at: string
}

type AuthContextType = {
  user: User | null
  loading: boolean
  login: () => Promise<void>
  logout: () => Promise<void>
  refetchUser: () => Promise<void>
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)
  const router = useRouter()

  // 初期ロード時にユーザー情報を取得
  const fetchUser = async () => {
    try {
      const response = await fetch(`${baseURL}/api/auth/user`, { credentials: 'include' })

      if (response.ok) {
        const data = await response.json()
        setUser(data.user)
      } else {
        setUser(null)
      }
    } catch (error) {
      console.error('Failed to fetch user:', error)
      setUser(null)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchUser()
  }, [])

  const login = async () => {
    try {
      const provider = new GoogleAuthProvider()
      const result = await signInWithPopup(auth, provider)
      const idToken = await result.user.getIdToken()

      // セッションクッキー作成
      const sessionRes = await fetch(`${baseURL}/api/auth/session`, {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ id_token: idToken }),
      })

      if (!sessionRes.ok) throw new Error('Session creation failed')

      // ユーザーAPI確認
      let userRes: Response = await fetch(`${baseURL}/api/auth/user`, { credentials: 'include' })

      if (!userRes.ok) {
        // ユーザーを新規作成
        userRes = await fetch(`${baseURL}/api/users`, {
          method: 'POST',
          credentials: 'include',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            uid: result.user.uid,
            name: result.user.displayName,
            type: 'general',
            plan: 'free',
          }),
        })
      }

      const userData: User = await userRes.json()
      console.log('Logged in user data:', userData)
      setUser(userData)
      router.push(`/profile/${userData.id}`)

      // クライアント側のFirebase Authセッションをクリア
      await signOut(auth)
    } catch (error) {
      setUser(null)
      throw error
    }
  }

  const logout = async () => {
    try {
      await fetch(`${baseURL}/api/auth/logout`, { method: 'POST' })
      setUser(null)
    } catch (error) {
      console.error('Logout failed:', error)
    }
  }

  const refetchUser = async () => {
    setLoading(true)
    await fetchUser()
  }

  return (
    <AuthContext.Provider value={{ user, loading, login, logout, refetchUser }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
