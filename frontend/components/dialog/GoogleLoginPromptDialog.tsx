'use client'

import React, { useState } from 'react'
import { useRouter } from 'next/navigation'
import { FcGoogle } from 'react-icons/fc'
import { useAuth } from '@/context/authContext'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'

/**
 * GoogleLoginPromptDialog - Googleログインを促すダイアログコンポーネント
 *
 * @param isOpen - ダイアログの表示状態
 * @param onClose - ダイアログを閉じる関数
 */
export default function GoogleLoginPromptDialog({
  isOpen,
  onClose,
}: {
  isOpen: boolean
  onClose: () => void
}) {
  const { googleLogin } = useAuth()
  const router = useRouter()
  const [isLoading, setIsLoading] = useState(false)

  const handleGoogleLogin = async () => {
    try {
      setIsLoading(true)
      await googleLogin()
      // ログイン処理は authContext 内で自動的にハンドリングされるため、
      // ここでのリダイレクトは不要
    } catch (error) {
      console.error('Google login error:', error)
    } finally {
      setIsLoading(false)
      onClose()
    }
  }

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle className="text-xl font-bold text-center text-primary">
            Googleログインが必要です
          </DialogTitle>
          <DialogDescription className="text-center">
            プロフィール設定を行う前に、Googleアカウントでのログインが必要です。
          </DialogDescription>
        </DialogHeader>
        <div className="flex flex-col items-center justify-center p-4">
          <div className="w-16 h-16 rounded-full bg-red-50 flex items-center justify-center mb-4">
            <FcGoogle className="w-10 h-10" />
          </div>
          <p className="text-sm text-gray-600 mb-4 text-center">
            Googleアカウントでログインすると、プロフィール設定が可能になります。
            また、作成した旅行動画を安全に保存することができます。
          </p>
        </div>
        <DialogFooter className="flex flex-col sm:flex-row gap-2">
          <Button
            variant="outline"
            onClick={onClose}
            className="w-full sm:w-auto"
            disabled={isLoading}
          >
            キャンセル
          </Button>
          <Button
            type="button"
            onClick={handleGoogleLogin}
            className="w-full sm:w-auto"
            disabled={isLoading}
          >
            {isLoading ? (
              <span className="flex items-center justify-center">
                <svg
                  className="animate-spin -ml-1 mr-2 h-4 w-4 text-white"
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                >
                  <circle
                    className="opacity-25"
                    cx="12"
                    cy="12"
                    r="10"
                    stroke="currentColor"
                    strokeWidth="4"
                  ></circle>
                  <path
                    className="opacity-75"
                    fill="currentColor"
                    d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                  ></path>
                </svg>
                ログイン中...
              </span>
            ) : (
              <span className="flex items-center justify-center">
                <FcGoogle className="mr-2 h-5 w-5" />
                Googleでログイン
              </span>
            )}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
