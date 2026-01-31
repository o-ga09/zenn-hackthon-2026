'use client'
import { Menu, X } from 'lucide-react'
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
    <header className="sticky top-0 z-50 bg-white/80 dark:bg-[#101922]/80 backdrop-blur-md border-b border-border">
      {error && (
        <div className="bg-red-50 border border-red-200 text-red-600 px-4 py-3 rounded-lg text-sm">
          {error}
        </div>
      )}
      <div className="max-w-[1200px] mx-auto px-6 h-16 flex items-center justify-between">
        <div className="flex items-center gap-2">
          <div className="size-8 bg-primary rounded-lg flex items-center justify-center text-white">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="20"
              height="20"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
            >
              <path d="M12 3v3m6.366-.366-2.12 2.12M21 12h-3m.366 6.366-2.12-2.12M12 18v3m-4.246-4.246-2.12 2.12M6 12H3m4.246-4.246-2.12-2.12" />
              <circle cx="12" cy="12" r="4" />
            </svg>
          </div>
          <Link href="/" className="text-lg font-extrabold tracking-tight text-foreground">
            Tavinikkiy <span className="text-primary">Agent</span>
          </Link>
        </div>

        {/* デスクトップナビ */}
        <nav className="hidden md:flex items-center gap-8">
          <a
            href="#features"
            className="text-sm font-semibold hover:text-primary transition-colors"
          >
            Features
          </a>
          <a
            href="#pricing"
            className="text-sm font-semibold hover:text-primary transition-colors"
          >
            Pricing
          </a>
          <a
            href="#contact"
            className="text-sm font-semibold hover:text-primary transition-colors"
          >
            Support
          </a>
        </nav>

        {/* デスクトップボタン */}
        <div className="hidden md:flex items-center gap-3">
          {user ? (
            <>
              <Link href="/dashboard">
                <Button
                  variant="outline"
                  className="text-sm font-bold px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-full transition-all"
                >
                  ダッシュボードへ
                </Button>
              </Link>
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  {user.profileImage ? (
                    <Image
                      src={user.profileImage}
                      alt="User Profile"
                      width={40}
                      height={40}
                      className="rounded-full cursor-pointer"
                    />
                  ) : (
                    <span className="size-10 rounded-full bg-gray-200 flex items-center justify-center text-gray-500 font-bold cursor-pointer">
                      {user.displayName ? user.displayName.charAt(0) : 'U'}
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
                          href={`/settings/${user.name}`}
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
            <>
              <button
                type="button"
                className="text-sm font-bold px-4 py-2 text-foreground hover:bg-gray-100 dark:hover:bg-gray-800 rounded-full transition-all"
                onClick={handleGoogleLogin}
                disabled={loading}
              >
                Login
              </button>
              <Button
                className="bg-primary hover:bg-primary/90 text-white text-sm font-bold px-5 py-2 rounded-full kawaii-shadow transition-all"
                onClick={handleGoogleLogin}
                disabled={loading}
              >
                <FcGoogle className="h-4 w-4 mr-2" />
                Start for Free
              </Button>
            </>
          )}
        </div>

        {/* モバイルメニューボタン */}
        <button
          className="md:hidden p-2 rounded-full hover:bg-gray-100 dark:hover:bg-gray-800"
          onClick={toggleMenu}
          aria-label={isMenuOpen ? 'メニューを閉じる' : 'メニューを開く'}
        >
          {isMenuOpen ? (
            <X className="w-6 h-6 text-foreground" />
          ) : (
            <Menu className="w-6 h-6 text-foreground" />
          )}
        </button>
      </div>

      {/* モバイルメニュー */}
      <AnimatePresence>
        {isMenuOpen && (
          <motion.div
            className="fixed inset-0 bg-background dark:bg-[#101922] z-40 md:hidden"
            initial={{ opacity: 0, x: '100%' }}
            animate={{ opacity: 1, x: 0 }}
            exit={{ opacity: 0, x: '100%' }}
            transition={{ type: 'spring', stiffness: 300, damping: 30 }}
          >
            {/* 閉じるボタン */}
            <button
              className="absolute top-6 right-4 p-2 rounded-full hover:bg-gray-100 dark:hover:bg-gray-800"
              onClick={() => setIsMenuOpen(false)}
              aria-label="メニューを閉じる"
            >
              <X className="w-6 h-6 text-foreground" />
            </button>

            <div className="flex flex-col items-center justify-center h-full space-y-8">
              <a
                href="#features"
                className="text-lg font-semibold hover:text-primary transition-colors"
                onClick={() => setIsMenuOpen(false)}
              >
                Features
              </a>
              <a
                href="#pricing"
                className="text-lg font-semibold hover:text-primary transition-colors"
                onClick={() => setIsMenuOpen(false)}
              >
                Pricing
              </a>
              <a
                href="#contact"
                className="text-lg font-semibold hover:text-primary transition-colors"
                onClick={() => setIsMenuOpen(false)}
              >
                Support
              </a>
              <div className="pt-8 flex flex-col gap-4">
                <button
                  type="button"
                  className="flex items-center justify-center gap-2 px-6 py-3 border border-gray-300 rounded-full bg-white hover:bg-gray-50 disabled:opacity-50 transition-colors"
                  onClick={() => {
                    handleGoogleLogin()
                    setIsMenuOpen(false)
                  }}
                  disabled={loading}
                >
                  <FcGoogle className="h-5 w-5" />
                  <span className="font-medium">Googleでログイン</span>
                </button>
              </div>
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </header>
  )
}
