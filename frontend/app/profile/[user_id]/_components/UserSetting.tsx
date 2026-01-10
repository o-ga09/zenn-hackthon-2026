'use client'

import React, { useState, useEffect } from 'react'
import { useAuth } from '@/context/authContext'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { toast } from 'sonner'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Label } from '@/components/ui/label'
import { Switch } from '@/components/ui/switch'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { ProfileFormData, profileFormSchema, useUpdateUser } from '@/api/user'
import { User, Mail, Coins, Globe, Lock, Save, X, Edit, Loader2 } from 'lucide-react'

export default function UserSetting() {
  const { user, refetchUser } = useAuth()
  const [isEditingProfile, setIsEditingProfile] = useState(false)
  const [isEditingPrivacy, setIsEditingPrivacy] = useState(false)

  // プロフィール更新用のフォーム（プライバシー設定も含む）
  const form = useForm<ProfileFormData>({
    resolver: zodResolver(profileFormSchema),
    defaultValues: {
      displayName: user?.displayName || '',
      isPublic: user?.isPublic || false,
    },
  })

  // ユーザーデータが更新された時にフォームのデフォルト値を更新
  useEffect(() => {
    if (user) {
      form.reset({
        displayName: user.displayName || '',
        isPublic: user.isPublic || false,
      })
    }
  }, [user, form])

  // ユーザー更新のミューテーション
  const updateUser = useUpdateUser(user?.id || '')

  // プロフィール情報の更新処理
  const onSubmitProfile = async (data: ProfileFormData) => {
    if (!user) return
    try {
      await updateUser.mutateAsync({
        ...user,
        displayName: data.displayName,
        isPublic: data.isPublic,
      })

      toast.success('プロフィールを更新しました')
      refetchUser()
      setIsEditingProfile(false)
    } catch {
      toast.error('プロフィールの更新に失敗しました')
    }
  }

  // プライバシー設定の更新処理
  const onSubmitPrivacy = async (data: ProfileFormData) => {
    try {
      if (!user) return
      await updateUser.mutateAsync({
        ...user,
        isPublic: data.isPublic,
      })

      setIsEditingPrivacy(false)
      refetchUser()
      toast.success('プライバシー設定を更新しました')
    } catch {
      toast.error('プライバシー設定の更新に失敗しました')
    }
  }
  if (!user) {
    return (
      <div className="flex items-center justify-center h-[60vh]">
        <p className="text-gray-500">ユーザー情報を読み込み中...</p>
      </div>
    )
  }

  // プランの表示を決定
  const getPlanBadgeVariant = (plan: string) => {
    switch (plan?.toLowerCase()) {
      case 'premium':
        return 'default' as const
      case 'pro':
        return 'secondary' as const
      default:
        return 'outline' as const
    }
  }

  // ログイン方法を判定（簡易的な実装）
  const getLoginMethod = () => {
    return 'Google'
  }

  return (
    <div className="container mx-auto p-3 max-w-4xl">
      <div className="space-y-4 max-h-[calc(100vh-8rem)] overflow-y-auto">
        {/* ユーザープロフィール画像セクション */}
        <Card className="h-fit">
          <CardHeader className="pb-3">
            <div className="flex items-center justify-between">
              <CardTitle className="flex items-center gap-2 text-lg">
                <User className="h-4 w-4" />
                プロフィール
              </CardTitle>
              {!isEditingProfile && (
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => setIsEditingProfile(true)}
                  className="h-8 w-8 p-0"
                >
                  <Edit className="h-4 w-4" />
                </Button>
              )}
            </div>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center gap-4">
              <Avatar className="h-16 w-16">
                <AvatarImage src={user.photoURL || ''} alt={user.name || ''} />
                <AvatarFallback className="text-xl">{user.name?.charAt(0) || 'U'}</AvatarFallback>
              </Avatar>
              <div className="flex-1 min-w-0">
                <h3 className="text-xl font-semibold truncate">{user.displayName || 'ユーザー'}</h3>
                <p className="text-gray-600 text-sm truncate">@{user.name || user.id}</p>
              </div>
            </div>

            {isEditingProfile && (
              <Form {...form}>
                <form onSubmit={form.handleSubmit(onSubmitProfile)} className="space-y-4">
                  <FormField
                    control={form.control}
                    name="displayName"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>ユーザー名</FormLabel>
                        <FormControl>
                          <Input {...field} placeholder="ユーザー名を入力" />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <div className="flex gap-2">
                    <Button
                      type="submit"
                      size="sm"
                      disabled={updateUser.isPending}
                      className="flex items-center gap-2"
                    >
                      {updateUser.isPending ? (
                        <Loader2 className="h-4 w-4 animate-spin" />
                      ) : (
                        <Save className="h-4 w-4" />
                      )}
                      保存
                    </Button>
                    <Button
                      type="button"
                      variant="outline"
                      size="sm"
                      onClick={() => {
                        setIsEditingProfile(false)
                        form.reset()
                      }}
                    >
                      <X className="h-4 w-4" />
                      キャンセル
                    </Button>
                  </div>
                </form>
              </Form>
            )}
          </CardContent>
        </Card>

        {/* プライバシー設定 */}
        <Card className="h-fit">
          <CardHeader className="pb-3">
            <div className="flex items-center justify-between">
              <CardTitle className="flex items-center gap-2 text-lg">
                {form.watch('isPublic') ? (
                  <Globe className="h-4 w-4" />
                ) : (
                  <Lock className="h-4 w-4" />
                )}
                プライバシー設定
              </CardTitle>
              {!isEditingPrivacy && (
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => setIsEditingPrivacy(true)}
                  className="h-8 w-8 p-0"
                >
                  <Edit className="h-4 w-4" />
                </Button>
              )}
            </div>
          </CardHeader>
          <CardContent className="space-y-3">
            {isEditingPrivacy ? (
              <Form {...form}>
                <form onSubmit={form.handleSubmit(onSubmitPrivacy)} className="space-y-4">
                  <FormField
                    control={form.control}
                    name="isPublic"
                    render={({ field }) => (
                      <FormItem>
                        <div className="flex items-center justify-between">
                          <div className="space-y-1">
                            <FormLabel className="text-sm font-medium">プロフィール公開</FormLabel>
                            <p className="text-xs text-gray-500">
                              {field.value
                                ? '他のユーザーがあなたのプロフィールを閲覧できます'
                                : 'プロフィールは非公開に設定されています'}
                            </p>
                          </div>
                          <FormControl>
                            <Switch checked={field.value} onCheckedChange={field.onChange} />
                          </FormControl>
                        </div>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <div className="flex gap-2 pt-2">
                    <Button
                      type="submit"
                      size="sm"
                      disabled={updateUser.isPending}
                      className="flex items-center gap-2"
                    >
                      {updateUser.isPending ? (
                        <Loader2 className="h-4 w-4 animate-spin" />
                      ) : (
                        <Save className="h-4 w-4" />
                      )}
                      保存
                    </Button>
                    <Button
                      type="button"
                      variant="outline"
                      size="sm"
                      onClick={() => {
                        setIsEditingPrivacy(false)
                        form.reset({
                          displayName: user?.displayName || '',
                          isPublic: user?.isPublic || false,
                        })
                      }}
                    >
                      <X className="h-4 w-4" />
                      キャンセル
                    </Button>
                  </div>
                </form>
              </Form>
            ) : (
              <div className="flex items-center justify-between">
                <div className="space-y-1">
                  <Label className="text-sm font-medium">プロフィール公開</Label>
                  <p className="text-xs text-gray-500">
                    {form.watch('isPublic')
                      ? '他のユーザーがあなたのプロフィールを閲覧できます'
                      : 'プロフィールは非公開に設定されています'}
                  </p>
                </div>
                <Switch checked={form.watch('isPublic')} disabled />
              </div>
            )}
          </CardContent>
        </Card>

        {/* 一般セクション */}
        <Card className="h-fit">
          <CardHeader className="pb-3">
            <CardTitle className="text-lg">一般設定</CardTitle>
          </CardHeader>
          <CardContent className="space-y-1">
            {/* 現在のプラン */}
            <div className="flex items-center justify-between">
              <Label className="text-sm font-medium">プラン</Label>
              <Badge variant={getPlanBadgeVariant(user.type)} className="text-xs">
                {user.type === 'general' ? 'フリー' : user.type}
              </Badge>
            </div>

            {/* 使用量（残トークン） */}
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <Coins className="h-4 w-4 text-yellow-500" />
                <Label className="text-sm font-medium">残りトークン</Label>
              </div>
              <span className="text-sm font-mono">
                {user.tokenBalance?.toLocaleString() || '0'}
              </span>
            </div>
          </CardContent>
        </Card>

        {/* アカウントセクション */}
        <Card className="h-fit">
          <CardHeader className="pb-3">
            <CardTitle className="flex items-center gap-2 text-lg">
              <Mail className="h-4 w-4" />
              アカウント設定
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            {/* ログイン元 */}
            <div className="flex items-center justify-between">
              <Label className="text-sm font-medium">ログイン方法</Label>
              <Badge variant="secondary" className="text-xs">
                {getLoginMethod()}
              </Badge>
            </div>

            {/* ユーザーID */}
            <div className="space-y-1">
              <Label className="text-sm font-medium">ユーザーID</Label>
              <p className="text-xs text-gray-600 font-mono px-2 py-1 bg-gray-50 rounded border truncate">
                {user.id}
              </p>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
