'use client'

import { useEffect, useState } from 'react'
import { useGetMediaAnalytics, useUpdateMediaAnalytics } from '@/api/mediaApi'
import { Media } from '@/api/types'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { TagInput } from '@/components/ui/tag-input'
import { Loader2, AlertCircle } from 'lucide-react'

interface MediaAnalyticsDialogProps {
  media: Media | null
  open: boolean
  onOpenChange: (open: boolean) => void
}

const MOOD_OPTIONS = ['楽しい', '穏やか', 'エキサイティング', 'ロマンチック'] as const

export function MediaAnalyticsDialog({ media, open, onOpenChange }: MediaAnalyticsDialogProps) {
  const [description, setDescription] = useState('')
  const [mood, setMood] = useState('')
  const [objects, setObjects] = useState<string[]>([])
  const [landmarks, setLandmarks] = useState<string[]>([])
  const [activities, setActivities] = useState<string[]>([])

  const {
    data: analytics,
    isLoading,
    error,
  } = useGetMediaAnalytics(media?.id || '', open && media?.status === 'completed' && !!media?.id)

  const updateMutation = useUpdateMediaAnalytics(media?.id || '')

  // 分析結果を取得したらフォームにセット
  useEffect(() => {
    if (analytics) {
      setDescription(analytics.description || '')
      setMood(analytics.mood || '')
      setObjects(analytics.objects || [])
      setLandmarks(analytics.landmarks || [])
      setActivities(analytics.activities || [])
    }
  }, [analytics])

  const handleSave = () => {
    updateMutation.mutate(
      {
        description,
        mood,
        objects,
        landmarks,
        activities,
      },
      {
        onSuccess: () => {
          onOpenChange(false)
        },
      }
    )
  }

  const handleCancel = () => {
    // 元の値に戻す
    if (analytics) {
      setDescription(analytics.description || '')
      setMood(analytics.mood || '')
      setObjects(analytics.objects || [])
      setLandmarks(analytics.landmarks || [])
      setActivities(analytics.activities || [])
    }
    onOpenChange(false)
  }

  // 分析中または失敗の場合
  if (media?.status === 'pending' || media?.status === 'uploading') {
    return (
      <Dialog open={open} onOpenChange={onOpenChange}>
        <DialogContent className="max-w-lg sm:max-w-full">
          <DialogHeader>
            <DialogTitle>分析中です</DialogTitle>
            <DialogDescription>
              メディアを分析しています。しばらくお待ちください。
            </DialogDescription>
          </DialogHeader>
          <div className="flex items-center justify-center py-8">
            <Loader2 className="w-8 h-8 animate-spin text-primary" />
          </div>
        </DialogContent>
      </Dialog>
    )
  }

  if (media?.status === 'failed') {
    return (
      <Dialog open={open} onOpenChange={onOpenChange}>
        <DialogContent className="max-w-lg sm:max-w-full">
          <DialogHeader>
            <DialogTitle>分析失敗</DialogTitle>
            <DialogDescription>メディアの分析に失敗しました。</DialogDescription>
          </DialogHeader>
          <div className="flex items-center justify-center py-8 text-destructive">
            <AlertCircle className="w-8 h-8 mr-2" />
            <span>{media.error_message || '分析中にエラーが発生しました'}</span>
          </div>
          <div className="flex justify-end mt-4">
            <Button onClick={() => onOpenChange(false)}>閉じる</Button>
          </div>
        </DialogContent>
      </Dialog>
    )
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-lg max-h-[90vh] overflow-y-auto sm:max-w-full">
        <DialogHeader>
          <DialogTitle>メディアのタグを編集</DialogTitle>
          <DialogDescription>分析結果を確認・編集して保存してください。</DialogDescription>
        </DialogHeader>

        {isLoading ? (
          <div className="flex items-center justify-center py-8">
            <Loader2 className="w-8 h-8 animate-spin text-primary" />
          </div>
        ) : error ? (
          <div className="flex items-center justify-center py-8 text-destructive">
            <AlertCircle className="w-8 h-8 mr-2" />
            <span>分析結果の取得に失敗しました</span>
          </div>
        ) : (
          <div className="space-y-4">
            {/* Description */}
            <div>
              <Label htmlFor="description">説明</Label>
              <Textarea
                id="description"
                value={description}
                onChange={e => setDescription(e.target.value)}
                placeholder="メディアの説明を入力してください（最大500文字）"
                maxLength={500}
                className="min-h-[100px] resize-none"
              />
              <p className="text-xs text-muted-foreground mt-1">{description.length}/500文字</p>
            </div>

            {/* Mood */}
            <div>
              <Label htmlFor="mood">雰囲気</Label>
              <Select value={mood} onValueChange={setMood}>
                <SelectTrigger id="mood" className="w-full">
                  <SelectValue placeholder="雰囲気を選択してください" />
                </SelectTrigger>
                <SelectContent>
                  {MOOD_OPTIONS.map(option => (
                    <SelectItem key={option} value={option}>
                      {option}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            {/* Objects */}
            <div>
              <Label htmlFor="objects">Objects（物体）</Label>
              <TagInput tags={objects} onChange={setObjects} placeholder="例: 海, 砂浜, 人" />
            </div>

            {/* Landmarks */}
            <div>
              <Label htmlFor="landmarks">Landmarks（場所）</Label>
              <TagInput
                tags={landmarks}
                onChange={setLandmarks}
                placeholder="例: 沖縄, 美ら海水族館"
              />
            </div>

            {/* Activities */}
            <div>
              <Label htmlFor="activities">Activities（活動）</Label>
              <TagInput
                tags={activities}
                onChange={setActivities}
                placeholder="例: 海水浴, シュノーケリング"
              />
            </div>
          </div>
        )}

        <div className="flex justify-end gap-2 mt-4">
          <Button variant="outline" onClick={handleCancel}>
            キャンセル
          </Button>
          <Button onClick={handleSave} disabled={updateMutation.isPending || isLoading || !!error}>
            {updateMutation.isPending ? (
              <>
                <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                保存中...
              </>
            ) : (
              '保存'
            )}
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  )
}
