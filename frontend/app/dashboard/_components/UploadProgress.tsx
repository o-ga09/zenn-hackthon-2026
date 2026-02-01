'use client'

import { useMediaAnalysisSSE } from '@/api/mediaApi'
import { Progress } from '@/components/ui/progress'
import { Skeleton } from '@/components/ui/skeleton'
import { CheckCircle2, Loader2, AlertCircle, Upload } from 'lucide-react'

interface UploadProgressProps {
  mediaIds: string[] | null
}

export const UploadProgress = ({ mediaIds }: UploadProgressProps) => {
  const { status, error } = useMediaAnalysisSSE(mediaIds)

  if (!mediaIds || mediaIds.length === 0 || !status) {
    return <Skeleton className="h-32 w-full rounded-lg" />
  }

  // 全体の進捗を計算（各メディアの進捗の平均）
  const totalProgress = status.medias.reduce((sum, m) => sum + m.progress, 0)
  const progressPercent = Math.round((totalProgress / status.total_items) * 100)

  // 全体のステータスを判定
  const getOverallStatus = () => {
    if (status.all_completed) {
      return status.failed_items > 0 ? 'partial' : 'completed'
    }
    const hasUploading = status.medias.some(m => m.status === 'uploading')
    if (hasUploading) return 'uploading'
    return 'analyzing'
  }
  const overallStatus = getOverallStatus()

  return (
    <div className="border rounded-lg p-4 bg-card">
      <div className="flex items-center justify-between mb-3">
        <div className="flex items-center gap-2">
          {overallStatus === 'completed' && <CheckCircle2 className="h-5 w-5 text-green-600" />}
          {overallStatus === 'partial' && <AlertCircle className="h-5 w-5 text-yellow-600" />}
          {overallStatus === 'analyzing' && (
            <Loader2 className="h-5 w-5 text-blue-600 animate-spin" />
          )}
          {overallStatus === 'uploading' && (
            <Upload className="h-5 w-5 text-blue-600 animate-pulse" />
          )}

          <span className="text-sm font-medium">
            {overallStatus === 'completed' && '✓ 分析完了'}
            {overallStatus === 'partial' && '一部失敗'}
            {overallStatus === 'analyzing' && '分析中...'}
            {overallStatus === 'uploading' && 'アップロード中...'}
          </span>
        </div>

        <div className="text-right">
          <span className="text-sm text-muted-foreground">{progressPercent}%</span>
          {status.total_items > 0 && (
            <div className="text-xs text-muted-foreground">
              {status.completed_items}/{status.total_items} 件
            </div>
          )}
        </div>
      </div>

      <Progress value={progressPercent} className="h-2" />

      {error && (
        <div className="mt-3 p-3 bg-destructive/10 border border-destructive/20 rounded-md flex items-start gap-2">
          <AlertCircle className="h-4 w-4 text-destructive mt-0.5" />
          <p className="text-sm text-destructive">{error}</p>
        </div>
      )}

      {status.all_completed && (
        <div
          className={`mt-3 flex items-center gap-2 text-sm ${status.failed_items > 0 ? 'text-yellow-600' : 'text-green-600'}`}
        >
          {status.failed_items > 0 ? (
            <AlertCircle className="h-4 w-4" />
          ) : (
            <CheckCircle2 className="h-4 w-4" />
          )}
          <span>
            {status.completed_items}件の画像を分析しました
            {status.failed_items > 0 && ` (${status.failed_items}件失敗)`}
          </span>
        </div>
      )}
    </div>
  )
}
