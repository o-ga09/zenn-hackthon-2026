'use client'

import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Sparkles, ArrowRight, Loader2 } from 'lucide-react'
import React, { useEffect, useState } from 'react'
import { useFormContext } from 'react-hook-form'
import { TravelFormValues } from './form-schema'
import { useUploadForm } from './UploadFormContext'
import GenerationProgress from './GenerationProgress'

export default function VideoGenerationConfirm() {
  const { getValues } = useFormContext<TravelFormValues>()
  const { handleGenerateVideo, uploadedFiles, selectedMediaIds, isGenerating, generationError, vlogId } = useUploadForm()
  const [formValues, setFormValues] = useState({
    travelTitle: '',
    travelDate: '',
    travelLocation: '指定なし',
    travelDescription: '指定なし',
  })

  // コンポーネントがマウントされた時と再レンダリング時に最新のフォーム値を取得
  useEffect(() => {
    const values = getValues()
    setFormValues({
      travelTitle: values.travelTitle || '',
      travelDate: values.travelDate || '',
      travelLocation: values.travelLocation || '指定なし',
      travelDescription: values.travelDescription || '指定なし',
    })
  }, [getValues, uploadedFiles, selectedMediaIds]) // uploadedFilesまたはselectedMediaIdsが変更された時も再取得

  if (vlogId) {
    return <GenerationProgress vlogId={vlogId} />
  }

  return (
    <div className="space-y-4 md:space-y-6">
      <div className="bg-muted p-3 md:p-4 rounded-lg">
        <h3 className="font-semibold mb-2 text-sm md:text-base">確認事項</h3>
        <ul className="list-disc list-inside space-y-1 md:space-y-2 text-xs md:text-sm">
          <li>新規アップロード: {uploadedFiles.length}枚</li>
          <li>ライブラリから選択: {selectedMediaIds.length}枚</li>
          <li>旅行タイトル: {formValues.travelTitle || 'AIにお任せ'}</li>
          <li>旅行日: {formValues.travelDate || 'AIにお任せ'}</li>
          <li>場所: {formValues.travelLocation}</li>
          <li>
            説明: {formValues.travelDescription.substring(0, 30)}
            {formValues.travelDescription.length > 30 ? '...' : ''}
          </li>
        </ul>
      </div>

      {generationError && (
        <div className="bg-destructive/10 text-destructive p-3 md:p-4 rounded-lg text-xs md:text-sm">
          <p>{generationError}</p>
        </div>
      )}

      <Card className="border-0 shadow-lg bg-gradient-to-r from-primary to-secondary p-4 md:p-6 text-white">
        <div className="mb-3 md:mb-4">
          <h3 className="text-base md:text-lg lg:text-2xl font-bold mb-1 md:mb-2">
            動画生成の準備完了！
          </h3>
          <p className="opacity-90 text-xs md:text-sm lg:text-base">
            AIが{uploadedFiles.length + selectedMediaIds.length}枚の素材から素敵な動画を作成します
          </p>
        </div>
        <Button
          size="lg"
          variant="secondary"
          className="bg-white text-primary hover:bg-white/90 text-sm md:text-lg px-4 md:px-8 py-3 md:py-6 w-full"
          onClick={handleGenerateVideo}
          disabled={
            (uploadedFiles.length === 0 && selectedMediaIds.length === 0) ||
            isGenerating
          }
        >
          {isGenerating ? (
            <>
              <Loader2 className="mr-2 h-4 w-4 md:h-5 md:w-5 animate-spin" />
              リクエスト送信中...
            </>
          ) : (
            <>
              <Sparkles className="w-4 h-4 md:w-5 md:h-5 mr-2" />
              AI動画生成を開始
              <ArrowRight className="w-4 h-4 md:w-5 md:h-5 ml-2" />
            </>
          )}
        </Button>
      </Card>
    </div>
  )
}
