'use client'
import { Menu, Video, X } from 'lucide-react'
import React, { useState } from 'react'
import { Button } from '../ui/button'
import Link from 'next/link'
import { motion, AnimatePresence } from 'framer-motion'
import { FcGoogle } from 'react-icons/fc'
import { useAuth } from '@/context/authContext'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '../ui/dropdown-menu'
import Image from 'next/image'
import { useRouter } from 'next/navigation'

export default function LPHeader() {
  const [isMenuOpen, setIsMenuOpen] = useState(false)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const { user, login, logout } = useAuth()
  const router = useRouter()

  const toggleMenu = () => {
    setIsMenuOpen(!isMenuOpen)
  }

  const handleGoogleLogin = async () => {
    setLoading(true)
    try {
      // googleLoginにリダイレクト処理が含まれているので、ここではルーティングしない
      await login()
    } catch (err) {
      setError('Googleログインに失敗しました')
      console.error(err)
    } finally {
      setLoading(false)
    }
  }

  const handleLogout = async () => {
    try {
      await logout()
      router.push('/')
    } catch (error) {
      console.error('Logout error:', error)
    }
  }

  return (
    <header className="container mx-auto px-4 py-6 relative z-50">
      {error && (
        <div className="bg-red-50 border border-red-200 text-red-600 px-4 py-3 rounded-lg text-sm">
          {error}
        </div>
      )}
      <nav className="flex items-center justify-between">
        <div className="flex items-center space-x-2">
          <div className="w-8 h-8 bg-primary rounded-lg flex items-center justify-center">
            <Video className="w-5 h-5 text-primary-foreground" />
          </div>
          <Link href="/" className="text-lg font-bold text-foreground">
            <span className="text-xl font-bold text-foreground">TravelMoments</span>
          </Link>
        </div>

        {/* デスクトップメニュー */}
        <div className="hidden md:flex items-center space-x-4">
          {user ? (
            <>
              <Link href="/dashboard">
                <Button variant="outline" className="border rounded px-4 py-2 hover:bg-gray-50">
                  ダッシュボードへ
                </Button>
              </Link>
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  {user.photoURL ? (
                    <Image
                      src={user.photoURL}
                      alt="User Profile"
                      width={40}
                      height={40}
                      className="rounded-full"
                    />
                  ) : (
                    <span className="w-8 h-8 rounded-full bg-gray-200 flex items-center justify-center text-gray-500 font-bold">
                      {user.displayname ? user.displayname.charAt(0) : 'U'}
                    </span>
                  )}
                </DropdownMenuTrigger>
                <DropdownMenuContent
                  align="end"
                  className="w-[200px] bg-white rounded-md shadow-lg p-2 flex flex-col gap-1 z-50"
                >
                  {user ? (
                    <>
                      <DropdownMenuItem asChild className="focus:bg-gray-100 focus:outline-none">
                        <Link
                          href={`/settings/${user.username}`}
                          className="w-full px-3 py-2 text-sm"
                        >
                          ユーザー設定
                        </Link>
                      </DropdownMenuItem>
                    </>
                  ) : (
                    <DropdownMenuItem className="focus:bg-gray-100 focus:outline-none text-gray-400 cursor-not-allowed">
                      <span className="w-full px-3 py-2 text-sm">設定を読み込み中...</span>
                    </DropdownMenuItem>
                  )}
                  <DropdownMenuItem
                    className="focus:bg-gray-100 focus:outline-none text-red-500"
                    onSelect={e => {
                      e.preventDefault()
                      handleLogout()
                    }}
                  >
                    <button onClick={handleLogout} className="w-full px-3 py-2 text-sm text-left">
                      ログアウト
                    </button>
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </>
          ) : (
            <button
              type="button"
              className="w-full flex items-center justify-center gap-2 px-4 py-3 border border-gray-300 rounded-lg bg-white hover:bg-gray-50 disabled:opacity-50 transition-colors"
              onClick={handleGoogleLogin}
              disabled={loading}
            >
              <FcGoogle className="h-5 w-5" />
              Googleアカウントでログイン
            </button>
          )}
        </div>

        {/* モバイルメニューボタン */}
        <button
          className="md:hidden p-2 rounded-full hover:bg-gray-100"
          onClick={toggleMenu}
          aria-label={isMenuOpen ? 'メニューを閉じる' : 'メニューを開く'}
        >
          {isMenuOpen ? (
            <X className="w-6 h-6 text-foreground" />
          ) : (
            <Menu className="w-6 h-6 text-foreground" />
          )}
        </button>
      </nav>

      {/* モバイルメニュー */}
      <AnimatePresence>
        {isMenuOpen && (
          <motion.div
            className="fixed inset-0 bg-background z-40 md:hidden"
            initial={{ opacity: 0, x: '100%' }}
            animate={{ opacity: 1, x: 0 }}
            exit={{ opacity: 0, x: '100%' }}
            transition={{ type: 'spring', stiffness: 300, damping: 30 }}
          >
            {/* 閉じるボタン - 画面右上に配置 */}
            <button
              className="absolute top-6 right-4 p-2 rounded-full hover:bg-gray-100"
              onClick={() => setIsMenuOpen(false)}
              aria-label="メニューを閉じる"
            >
              <X className="w-6 h-6 text-foreground" />
            </button>

            <div className="flex flex-col items-center justify-center h-full space-y-8">
              <Link href="/login" onClick={() => setIsMenuOpen(false)}>
                <Button
                  variant="outline"
                  className="flex items-center gap-2 border rounded px-4 py-2 hover:bg-gray-50 transition-colors"
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    viewBox="0 0 48 48"
                    width="20px"
                    height="20px"
                    className="relative top-[1px]"
                  >
                    <path
                      fill="#FFC107"
                      d="M43.611,20.083H42V20H24v8h11.303c-1.649,4.657-6.08,8-11.303,8c-6.627,0-12-5.373-12-12c0-6.627,5.373-12,12-12c3.059,0,5.842,1.154,7.961,3.039l5.657-5.657C34.046,6.053,29.268,4,24,4C12.955,4,4,12.955,4,24c0,11.045,8.955,20,20,20c11.045,0,20-8.955,20-20C44,22.659,43.862,21.35,43.611,20.083z"
                    />
                    <path
                      fill="#FF3D00"
                      d="M6.306,14.691l6.571,4.819C14.655,15.108,18.961,12,24,12c3.059,0,5.842,1.154,7.961,3.039l5.657-5.657C34.046,6.053,29.268,4,24,4C16.318,4,9.656,8.337,6.306,14.691z"
                    />
                    <path
                      fill="#4CAF50"
                      d="M24,44c5.166,0,9.86-1.977,13.409-5.192l-6.19-5.238C29.211,35.091,26.715,36,24,36c-5.202,0-9.619-3.317-11.283-7.946l-6.522,5.025C9.505,39.556,16.227,44,24,44z"
                    />
                    <path
                      fill="#1976D2"
                      d="M43.611,20.083H42V20H24v8h11.303c-0.792,2.237-2.231,4.166-4.087,5.571c0.001-0.001,0.002-0.001,0.003-0.002l6.19,5.238C36.971,39.205,44,34,44,24C44,22.659,43.862,21.35,43.611,20.083z"
                    />
                  </svg>
                  <span className="text-xl font-medium ml-1">ログイン</span>
                </Button>
              </Link>
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </header>
  )
}
