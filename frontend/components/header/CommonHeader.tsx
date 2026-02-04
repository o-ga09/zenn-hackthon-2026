'use client'

import { Users, Check, Trash2 } from 'lucide-react'
import React from 'react'
import { Button } from '../ui/button'
import Link from 'next/link'
import { Bell } from 'lucide-react'
import { useAuth } from '@/context/authContext'
import { useNotifications } from '@/context/notificationContext'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
  DropdownMenuSeparator,
} from '../ui/dropdown-menu'
import { Badge } from '../ui/badge'
import Image from 'next/image'
import { useRouter } from 'next/navigation'

export default function CommonHeader() {
  const { user, logout } = useAuth()
  const { notifications, unreadCount, markAsRead, markAllAsRead, removeNotification, clearAll } =
    useNotifications()
  const router = useRouter()

  console.log('通知', notifications)
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
          {/* 通知ベル */}
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon" className="relative">
                <Bell className="h-5 w-5" />
                {unreadCount > 0 && (
                  <Badge
                    variant="destructive"
                    className="absolute -top-1 -right-1 h-5 w-5 flex items-center justify-center p-0 text-xs"
                  >
                    {unreadCount > 9 ? '9+' : unreadCount}
                  </Badge>
                )}
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent
              align="end"
              className="w-[320px] max-h-[400px] overflow-y-auto bg-white shadow-lg z-50"
            >
              <div className="flex items-center justify-between p-3 border-b">
                <h3 className="font-semibold text-sm">通知</h3>
                <div className="flex gap-1">
                  {unreadCount > 0 && (
                    <Button
                      variant="ghost"
                      size="sm"
                      className="h-7 px-2 text-xs"
                      onClick={markAllAsRead}
                    >
                      <Check className="h-3 w-3 mr-1" />
                      すべて既読
                    </Button>
                  )}
                  {notifications.length > 0 && (
                    <Button
                      variant="ghost"
                      size="sm"
                      className="h-7 px-2 text-xs text-red-500 hover:text-red-600"
                      onClick={clearAll}
                    >
                      <Trash2 className="h-3 w-3 mr-1" />
                      すべて削除
                    </Button>
                  )}
                </div>
              </div>

              {notifications.length === 0 ? (
                <div className="p-4 text-center text-sm text-muted-foreground">
                  通知はありません
                </div>
              ) : (
                notifications.map(notification => (
                  <div key={notification.id}>
                    <DropdownMenuItem
                      className={`flex flex-col items-start p-3 cursor-pointer focus:bg-gray-50 ${
                        !notification.read ? 'bg-blue-50/50' : ''
                      }`}
                      onClick={() => markAsRead(notification.id, notification.version)}
                    >
                      <div className="flex items-start justify-between w-full">
                        <div className="flex-1 min-w-0">
                          <div className="flex items-center gap-2 mb-1">
                            <span
                              className={`text-xs font-semibold ${
                                notification.type === 'success'
                                  ? 'text-green-600'
                                  : notification.type === 'error'
                                    ? 'text-red-600'
                                    : 'text-blue-600'
                              }`}
                            >
                              {notification.title}
                            </span>
                            {!notification.read && (
                              <div className="w-2 h-2 bg-blue-500 rounded-full"></div>
                            )}
                          </div>
                          <p className="text-xs text-gray-600 break-words">
                            {notification.message}
                          </p>
                          <p className="text-xs text-gray-400 mt-1">
                            {new Date(notification.timestamp).toLocaleTimeString('ja-JP', {
                              hour: '2-digit',
                              minute: '2-digit',
                            })}
                          </p>
                        </div>
                        <Button
                          variant="ghost"
                          size="sm"
                          className="h-6 w-6 p-0 ml-2 text-gray-400 hover:text-red-500"
                          onClick={e => {
                            e.stopPropagation()
                            removeNotification(notification.id)
                          }}
                        >
                          <Trash2 className="h-3 w-3" />
                        </Button>
                      </div>
                    </DropdownMenuItem>
                    <DropdownMenuSeparator />
                  </div>
                ))
              )}
            </DropdownMenuContent>
          </DropdownMenu>

          {user ? (
            <>
              <Link href="/dashboard">
                <Button variant="outline" className="border rounded px-4 py-2 hover:bg-gray-50">
                  ダッシュボードへ
                </Button>
              </Link>
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  {user.profileImage ? (
                    <Image
                      src={user?.profileImage || '/no-avatar.webp'}
                      alt="User Profile"
                      width={50}
                      height={50}
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
                        <Link href={`/profile/${user.name}`} className="w-full px-3 py-2 text-sm">
                          プロフィール設定
                        </Link>
                      </DropdownMenuItem>
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
