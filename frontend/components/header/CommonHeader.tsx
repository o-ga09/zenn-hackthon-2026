import { Users } from 'lucide-react'
import React from 'react'
import { Button } from '../ui/button'
import Link from 'next/link'
import { Bell } from 'lucide-react'
import { useAuth } from '@/context/authContext'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '../ui/dropdown-menu'
import Image from 'next/image'
import { useRouter } from 'next/navigation'

export default function CommonHeader() {
  const { user, logout } = useAuth()
  const router = useRouter()

  console.log('Current User in Header:', user)
  const handleLogout = async () => {
    try {
      await logout()
      router.push('/')
    } catch (error) {
      console.error('Logout error:', error)
    }
  }

  return (
    <header className="bg-white/80 backdrop-blur-sm shadow-sm sticky top-0 z-10">
      <div className="container mx-auto px-4 py-3 flex items-center justify-between">
        <div className="flex items-center space-x-4">
          <Link href="/dashboard" className="text-xl font-bold text-primary">
            Tavinikkiy
          </Link>

          {/* ナビゲーションリンク */}
          <nav className="hidden md:flex items-center space-x-4 ml-6">
            <Link
              href="/dashboard"
              className="text-sm font-medium text-gray-700 hover:text-primary"
            >
              ダッシュボード
            </Link>
            <Link href="/videos" className="text-sm font-medium text-gray-700 hover:text-primary">
              マイ動画
            </Link>
            <Link
              href="/any-one"
              className="text-sm font-medium text-gray-700 hover:text-primary flex items-center"
            >
              <Users className="h-4 w-4 mr-1" />
              みんなの動画
            </Link>
          </nav>
        </div>
        <div className="flex items-center space-x-2">
          <Button variant="ghost" size="icon">
            <Bell className="h-5 w-5" />
          </Button>
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
                      src={user?.photoURL || '/no-avatar.webp'}
                      alt="User Profile"
                      width={40}
                      height={40}
                      className="rounded-full"
                    />
                  ) : (
                    <span className="w-8 h-8 rounded-full bg-gray-200 flex items-center justify-center text-gray-500 font-bold">
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
                          href={`/profile/${user.name}/setting`}
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
            <></>
          )}
        </div>
      </div>
    </header>
  )
}
