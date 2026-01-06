'use client'

import React from 'react'
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetDescription,
  SheetFooter,
} from '@/components/ui/sheet'
import './travel-detail-sheet.css'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Calendar, Download, Share2, Play, Edit, X, Check } from 'lucide-react'
import { Travel } from '@/api/types'
import { enhanceTravelData } from '@/utils/travel-utils'
import Image from 'next/image'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Label } from '@/components/ui/label'
import { useForm, Controller, SubmitHandler } from 'react-hook-form'

// フォームの型定義
interface TravelFormValues {
  title: string
  description: string
  startDate: string
  endDate: string
}

interface TravelDetailSheetProps {
  travel: Travel | null
  isOpen: boolean
  onOpenChange: (open: boolean) => void
  onSave?: (updatedTravel: Partial<Travel>) => Promise<void>
}

export default function TravelDetailSheet({
  travel,
  isOpen,
  onOpenChange,
  onSave,
}: TravelDetailSheetProps) {
  // 画面サイズに応じてシートの向きを調整（モバイルでは下から、デスクトップでは右から）
  const [sheetSide, setSheetSide] = React.useState<'bottom' | 'right'>('right')

  // 編集モード管理
  const [isEditMode, setIsEditMode] = React.useState(false)

  // react-hook-formの設定
  const {
    control,
    handleSubmit,
    formState: { errors },
    reset,
    watch,
  } = useForm<TravelFormValues>({
    defaultValues: {
      title: '',
      description: '',
      startDate: '',
      endDate: '',
    },
  })

  // クライアントサイドでのみ実行される画面サイズ検出
  React.useEffect(() => {
    // 初期ロード時に一度だけ実行
    const updateSheetSide = () => {
      const isMobile = window.innerWidth < 768
      setSheetSide(isMobile ? 'bottom' : 'right')
    }

    // 初期実行
    updateSheetSide()

    // リサイズイベントのリスナー設定
    const handleResize = () => {
      updateSheetSide()
    }

    window.addEventListener('resize', handleResize)
    return () => window.removeEventListener('resize', handleResize)
  }, [])

  // シートが開かれた時にフォーム値を初期化
  React.useEffect(() => {
    if (isOpen && travel) {
      // フォームの値を旅行データで初期化
      reset({
        title: travel.title,
        description: travel.description || '',
        startDate: travel.startDate,
        endDate: travel.endDate,
      })
      // 編集モードはリセット
      setIsEditMode(false)
    }
  }, [isOpen, travel, reset])

  // 編集モードの切り替え
  const toggleEditMode = () => {
    if (isEditMode && travel) {
      // 編集モードを終了する場合は元の値に戻す
      reset({
        title: travel.title,
        description: travel.description || '',
        startDate: travel.startDate,
        endDate: travel.endDate,
      })
    }
    setIsEditMode(!isEditMode)
  }

  // 編集内容の保存処理
  const onSubmit: SubmitHandler<TravelFormValues> = data => {
    if (onSave && travel) {
      onSave({
        ...travel,
        ...data,
      })
    }
    setIsEditMode(false)
  }

  if (!travel) return null

  // 期間の計算
  const calculateDuration = (startDate: string, endDate: string): number => {
    const start = new Date(startDate)
    const end = new Date(endDate)
    return Math.ceil((end.getTime() - start.getTime()) / (1000 * 60 * 60 * 24)) + 1
  }

  const duration = calculateDuration(travel.startDate, travel.endDate)

  // サムネイルなどのデータを強化
  const enhancedTravel = enhanceTravelData(travel)

  return (
    <Sheet open={isOpen} onOpenChange={onOpenChange}>
      <SheetContent
        side={sheetSide}
        className={`w-full max-w-md overflow-y-auto ${
          sheetSide === 'bottom' ? 'p-6 pt-12 pb-8 mobile-sheet-content' : 'p-6'
        }`}
        style={{
          overflowY: 'auto',
          overscrollBehavior: 'contain',
        }}
      >
        <form onSubmit={handleSubmit(onSubmit)}>
          <SheetHeader className="pb-4 space-y-2">
            <div>
              {isEditMode ? (
                <div className="space-y-2">
                  <Label htmlFor="title">タイトル</Label>
                  <Controller
                    name="title"
                    control={control}
                    rules={{ required: 'タイトルは必須です' }}
                    render={({ field }) => (
                      <Input
                        id="title"
                        {...field}
                        placeholder="旅行のタイトルを入力"
                        className={`w-full ${errors.title ? 'border-red-500' : ''}`}
                      />
                    )}
                  />
                  {errors.title && (
                    <p className="text-red-500 text-xs mt-1">{errors.title.message}</p>
                  )}
                </div>
              ) : (
                <SheetTitle className="text-xl sm:text-2xl leading-tight">
                  {enhancedTravel.title}
                </SheetTitle>
              )}
            </div>

            <div>
              {isEditMode ? (
                <div className="grid grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <Label htmlFor="startDate">開始日</Label>
                    <Controller
                      name="startDate"
                      control={control}
                      rules={{ required: '開始日は必須です' }}
                      render={({ field }) => (
                        <Input
                          id="startDate"
                          type="date"
                          {...field}
                          className={`w-full ${errors.startDate ? 'border-red-500' : ''}`}
                        />
                      )}
                    />
                    {errors.startDate && (
                      <p className="text-red-500 text-xs mt-1">{errors.startDate.message}</p>
                    )}
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="endDate">終了日</Label>
                    <Controller
                      name="endDate"
                      control={control}
                      rules={{
                        required: '終了日は必須です',
                        validate: value => {
                          const startDate = watch('startDate')
                          if (startDate && value && new Date(startDate) > new Date(value)) {
                            return '終了日は開始日以降に設定してください'
                          }
                          return true
                        },
                      }}
                      render={({ field }) => (
                        <Input
                          id="endDate"
                          type="date"
                          {...field}
                          className={`w-full ${errors.endDate ? 'border-red-500' : ''}`}
                        />
                      )}
                    />
                    {errors.endDate && (
                      <p className="text-red-500 text-xs mt-1">{errors.endDate.message}</p>
                    )}
                  </div>
                </div>
              ) : (
                <>
                  <SheetDescription className="flex items-center gap-2 text-sm">
                    <Calendar className="w-4 h-4" />
                    <span>
                      {new Date(enhancedTravel.startDate).toLocaleDateString('ja-JP', {
                        year: 'numeric',
                        month: 'long',
                        day: 'numeric',
                      })}
                      {duration > 1 &&
                        ` 〜 ${new Date(enhancedTravel.endDate).toLocaleDateString('ja-JP', {
                          year: 'numeric',
                          month: 'long',
                          day: 'numeric',
                        })}`}
                    </span>
                  </SheetDescription>
                  {duration > 1 && (
                    <Badge variant="outline" className="mt-1 w-fit px-3 py-1">
                      {duration}日間
                    </Badge>
                  )}
                </>
              )}
            </div>
          </SheetHeader>

          {/* サムネイル画像 */}
          <div className="relative">
            <Image
              src={enhancedTravel.thumbnail || '/placeholder.webp'}
              alt={enhancedTravel.title}
              loading="eager"
              width={600}
              height={400}
              className="w-full object-cover rounded-md"
            />
            <div className="absolute inset-0 flex items-center justify-center">
              <Button
                size="sm"
                variant="secondary"
                className="bg-white/90 text-black hover:bg-white rounded-full w-12 h-12 p-0"
                type="button"
              >
                <Play className="w-6 h-6 ml-0.5" fill="currentColor" />
              </Button>
            </div>
          </div>

          {/* 説明文 */}
          <div className="py-4 mb-2">
            <h3 className="text-sm font-medium text-muted-foreground mb-2">旅行の説明</h3>
            <div>
              {isEditMode ? (
                <div className="space-y-2">
                  <Controller
                    name="description"
                    control={control}
                    render={({ field }) => (
                      <Textarea
                        {...field}
                        placeholder="旅行の説明を入力してください"
                        className="min-h-[120px]"
                      />
                    )}
                  />
                </div>
              ) : (
                <p className="text-sm whitespace-pre-line leading-relaxed">
                  {enhancedTravel.description || '説明はありません'}
                </p>
              )}
            </div>
          </div>

          {/* 詳細情報 - 編集モード時は非表示 */}
          {!isEditMode && (
            <>
              <div className="py-4 mb-2">
                <h3 className="text-sm font-medium text-muted-foreground mb-3">旅行情報</h3>
                <div className="grid grid-cols-2 gap-4 text-sm">
                  <div className="bg-muted/40 p-3 rounded-md">
                    <span className="text-xs text-muted-foreground block mb-1">開始日</span>
                    <p className="font-medium">
                      {new Date(enhancedTravel.startDate).toLocaleDateString('ja-JP', {
                        year: 'numeric',
                        month: 'long',
                        day: 'numeric',
                      })}
                    </p>
                  </div>
                  <div className="bg-muted/40 p-3 rounded-md">
                    <span className="text-xs text-muted-foreground block mb-1">終了日</span>
                    <p className="font-medium">
                      {new Date(enhancedTravel.endDate).toLocaleDateString('ja-JP', {
                        year: 'numeric',
                        month: 'long',
                        day: 'numeric',
                      })}
                    </p>
                  </div>
                  <div className="col-span-2 bg-muted/40 p-3 rounded-md">
                    <span className="text-xs text-muted-foreground block mb-1">共有ID</span>
                    <p className="font-mono text-xs overflow-auto">{enhancedTravel.sharedId}</p>
                  </div>
                </div>
              </div>

              {/* 作成・更新情報 */}
              <div className="py-4">
                <h3 className="text-sm font-medium text-muted-foreground mb-3">メタデータ</h3>
                <div className="grid grid-cols-2 gap-4 text-sm">
                  <div className="p-3 rounded-md border border-muted/50">
                    <span className="text-xs text-muted-foreground block mb-1">作成日</span>
                    <p className="font-medium">
                      {new Date(enhancedTravel.created_at).toLocaleDateString('ja-JP', {
                        year: 'numeric',
                        month: 'long',
                        day: 'numeric',
                      })}
                    </p>
                  </div>
                  <div className="p-3 rounded-md border border-muted/50">
                    <span className="text-xs text-muted-foreground block mb-1">更新日</span>
                    <p className="font-medium">
                      {new Date(enhancedTravel.updated_at).toLocaleDateString('ja-JP', {
                        year: 'numeric',
                        month: 'long',
                        day: 'numeric',
                      })}
                    </p>
                  </div>
                </div>
              </div>
            </>
          )}

          <SheetFooter className="flex-row gap-3 pt-6 pb-2 mt-4 border-t">
            {isEditMode ? (
              <>
                <Button
                  type="button"
                  variant="outline"
                  className="flex-1 gap-2 h-11"
                  onClick={toggleEditMode}
                >
                  <X className="w-4 h-4" />
                  <span>キャンセル</span>
                </Button>
                <Button type="submit" className="flex-1 gap-2 h-11">
                  <Check className="w-4 h-4" />
                  <span>保存</span>
                </Button>
              </>
            ) : (
              <>
                <Button type="button" variant="outline" className="flex-1 gap-2 h-11">
                  <Download className="w-4 h-4" />
                  <span>ダウンロード</span>
                </Button>
                <Button type="button" variant="outline" className="flex-1 gap-2 h-11">
                  <Share2 className="w-4 h-4" />
                  <span>共有</span>
                </Button>
                <Button type="button" className="flex-1 gap-2 h-11" onClick={toggleEditMode}>
                  <Edit className="w-4 h-4" />
                  <span>編集</span>
                </Button>
              </>
            )}
          </SheetFooter>
        </form>
      </SheetContent>
    </Sheet>
  )
}
