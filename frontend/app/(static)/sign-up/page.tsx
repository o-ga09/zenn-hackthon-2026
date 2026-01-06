'use client'

import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import Link from 'next/link'
import React, { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useCreateUser } from '@/api/userApi'
import { useRouter } from 'next/navigation'
import { useAuth } from '@/context/authContext'
import GoogleLoginPromptDialog from '@/components/dialog/GoogleLoginPromptDialog'
import { UserInputFrontend } from '@/api/types'

// フォームのバリデーションスキーマ
const signUpSchema = z.object({
  name: z
    .string()
    .min(3, 'ユーザーIDは3文字以上で入力してください')
    .max(20, 'ユーザーIDは20文字以下で入力してください')
    .regex(/^[a-z0-9_-]+$/, 'ユーザーIDは半角英数字、ハイフン、アンダースコアのみ使用できます'),
  display_name: z
    .string()
    .min(1, '表示名は必須です')
    .max(30, '表示名は30文字以下で入力してください'),
})

type SignUpFormValues = z.infer<typeof signUpSchema>

export default function SignUp() {
  const router = useRouter()
  const { currentUser } = useAuth()
  const [error, setError] = useState<string | null>(null)
  const [showLoginPrompt, setShowLoginPrompt] = useState(false)

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<SignUpFormValues>({
    resolver: zodResolver(signUpSchema),
    defaultValues: {
      name: '',
      display_name: '',
    },
  })

  // Googleログイン済みか確認し、未ログインならダイアログを表示
  React.useEffect(() => {
    if (!currentUser) {
      setShowLoginPrompt(true)
    }
  }, [currentUser])

  const createUser = useCreateUser()

  const onSubmit = async (data: SignUpFormValues) => {
    if (!currentUser) {
      setError('認証情報が見つかりません。再度ログインしてください。')
      setShowLoginPrompt(true)
      return
    }

    try {
      setError(null)

      // フロントエンドの形式でデータを作成
      const userData: UserInputFrontend = {
        firebase_id: currentUser.uid,
        name: data.name,
        display_name: data.display_name,
      }

      // APIが型変換を内部で行うように修正したuseCreateUserを使用
      await createUser.mutateAsync(userData)
      router.push('/dashboard') // 登録成功後ダッシュボードへリダイレクト
    } catch (err: any) {
      setError(err.response?.data?.message || 'ユーザー登録中にエラーが発生しました')
      console.error('ユーザー登録エラー:', err)
    }
  }

  return (
    <div className="flex items-center justify-center h-full py-4 px-4">
      {/* Googleログイン促進ダイアログ */}
      <GoogleLoginPromptDialog isOpen={showLoginPrompt} onClose={() => setShowLoginPrompt(false)} />

      <div className="max-w-md w-full bg-white/60 backdrop-blur-sm rounded-2xl p-6 shadow-lg">
        <div className="text-center mb-4">
          <h1 className="text-2xl font-bold mb-1 text-primary">プロフィール設定</h1>
          <p className="text-sm text-gray-600">
            旅の記録を始める前に、プロフィールを設定しましょう
          </p>
        </div>

        {error && (
          <div className="bg-red-50 border border-red-200 text-red-600 px-4 py-3 rounded-lg text-sm mb-4">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-3">
          <div>
            <label htmlFor="name" className="block text-sm font-medium text-gray-700 mb-1">
              ユーザーID
            </label>
            <Input
              id="name"
              type="text"
              placeholder="ユーザーID（英数字）"
              {...register('name')}
              className="w-full"
            />
            {errors.name && <p className="text-xs text-red-500 mt-0.5">{errors.name.message}</p>}
            <p className="text-xs text-gray-500 mt-0.5">
              ※ユーザーIDは一度設定すると変更できません
            </p>
          </div>

          <div>
            <label htmlFor="display_name" className="block text-sm font-medium text-gray-700 mb-1">
              表示名
            </label>
            <Input
              id="display_name"
              type="text"
              placeholder="表示名"
              {...register('display_name')}
              className="w-full"
            />
            {errors.display_name && (
              <p className="text-xs text-red-500 mt-0.5">{errors.display_name.message}</p>
            )}
          </div>

          <Button
            type="submit"
            className="w-full text-white py-2 px-4 rounded-lg font-semibold hover:opacity-90 transition-opacity disabled:opacity-50"
            disabled={isSubmitting}
          >
            {isSubmitting ? '登録中...' : 'プロフィールを登録する'}
          </Button>

          <div className="text-center text-xs text-gray-600">
            <Link href="/register" className="text-blue-600 hover:text-blue-800 underline">
              戻る
            </Link>
          </div>
        </form>
      </div>
    </div>
  )
}
