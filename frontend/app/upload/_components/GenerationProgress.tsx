'use client'

import React, { useEffect } from 'react'
import { Progress } from '@/components/ui/progress'
import { Card } from '@/components/ui/card'
import { CheckCircle2, Loader2, AlertCircle } from 'lucide-react'
import { useVlogSSE } from '@/api/vlogAPi'
import { useRouter } from 'next/navigation'

interface GenerationProgressProps {
  vlogId: string
}

export default function GenerationProgress({ vlogId }: GenerationProgressProps) {
  const router = useRouter()
  const { status, error } = useVlogSSE(vlogId)

  useEffect(() => {
    if (status?.status === 'completed') {
      const timer = setTimeout(() => {
        router.push('/videos')
      }, 2000)
      return () => clearTimeout(timer)
    }
  }, [status, router])

  if (error) {
    return (
      <div className="bg-destructive/10 text-destructive p-4 rounded-lg flex items-center gap-2">
        <AlertCircle className="w-5 h-5" />
        <p>エラーが発生しました: {error.message}</p>
      </div>
    )
  }

  const progress = status ? Math.round(status.progress * 100) : 0
  const isCompleted = status?.status === 'completed'
  const isFailed = status?.status === 'failed'

  return (
    <Card className="p-6 space-y-6">
      <div className="text-center space-y-2">
        <h3 className="text-xl font-bold">
          {isCompleted ? '動画の生成が完了しました！' : isFailed ? '生成に失敗しました' : 'AI動画を生成中...'}
        </h3>
        <p className="text-muted-foreground text-sm">
          {isCompleted ? 'まもなく一覧ページへ移動します' : isFailed ? status?.error_message : 'これには数分かかる場合があります。このページを閉じてもバックグラウンドで処理は継続されます。'}
        </p>
      </div>

      <div className="space-y-2">
        <div className="flex justify-between text-sm font-medium">
          <span>進捗</span>
          <span>{progress}%</span>
        </div>
        <Progress value={progress} className="h-3" />
      </div>

      <div className="grid grid-cols-1 gap-4">
        <div className={`flex items-center gap-3 p-3 rounded-lg border ${progress > 10 ? 'bg-primary/5 border-primary/20' : 'bg-muted/50 border-transparent'}`}>
          {progress > 10 ? <CheckCircle2 className="w-5 h-5 text-primary" /> : <Loader2 className="w-5 h-5 animate-spin text-muted-foreground" />}
          <span className={progress > 10 ? 'font-medium' : 'text-muted-foreground'}>メディアの分析</span>
        </div>
        <div className={`flex items-center gap-3 p-3 rounded-lg border ${progress > 40 ? 'bg-primary/5 border-primary/20' : 'bg-muted/50 border-transparent'}`}>
          {progress > 40 ? <CheckCircle2 className="w-5 h-5 text-primary" /> : <Loader2 className="w-5 h-5 animate-spin text-muted-foreground" />}
          <span className={progress > 40 ? 'font-medium' : 'text-muted-foreground'}>ストーリーの構成</span>
        </div>
        <div className={`flex items-center gap-3 p-3 rounded-lg border ${progress > 70 ? 'bg-primary/5 border-primary/20' : 'bg-muted/50 border-transparent'}`}>
          {progress > 70 ? <CheckCircle2 className="w-5 h-5 text-primary" /> : <Loader2 className="w-5 h-5 animate-spin text-muted-foreground" />}
          <span className={progress > 70 ? 'font-medium' : 'text-muted-foreground'}>動画の書き出しと保存</span>
        </div>
      </div>
    </Card>
  )
}
