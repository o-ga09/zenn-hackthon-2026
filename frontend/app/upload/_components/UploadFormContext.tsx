'use client'

import React, { createContext, useContext, useState } from 'react'
import { useForm, FormProvider } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { TravelFormValues, UploadStep, travelFormSchema } from './form-schema'
import { useRouter } from 'next/navigation'
import { useAuth } from '@/context/authContext'
import { apiClient } from '@/api/client'

interface UploadFormContextProps {
  step: UploadStep
  setStep: (step: UploadStep) => void
  nextStep: () => void
  prevStep: () => void
  handleGenerateVideo: () => Promise<void>
  addFiles: (files: File[]) => void
  removeFile: (index: number) => void
  uploadedFiles: File[]
  selectedMediaIds: string[]
  toggleMediaId: (id: string) => void
  isGenerating: boolean
  generationError: string | null
  vlogId: string | null
}

const UploadFormContext = createContext<UploadFormContextProps | undefined>(undefined)

export const useUploadForm = () => {
  const context = useContext(UploadFormContext)
  if (!context) {
    throw new Error('useUploadForm must be used within a UploadFormProvider')
  }
  return context
}

interface UploadFormProviderProps {
  children: React.ReactNode
}

export function UploadFormProvider({ children }: UploadFormProviderProps) {
  const router = useRouter()
  const { user } = useAuth()
  const [step, setStep] = useState<UploadStep>('upload')
  const [uploadedFiles, setUploadedFiles] = useState<File[]>([])
  const [selectedMediaIds, setSelectedMediaIds] = useState<string[]>([])

  const [isGenerating, setIsGenerating] = useState(false)
  const [generationError, setGenerationError] = useState<string | null>(null)
  const [vlogId, setVlogId] = useState<string | null>(null)

  const methods = useForm<TravelFormValues>({
    resolver: zodResolver(travelFormSchema),
    defaultValues: {
      travelTitle: '',
      travelDate: '',
      travelLocation: '',
      travelDescription: '',
      uploadedFiles: [],
      mediaIds: [],
    },
    mode: 'onChange',
  })

  const nextStep = () => {
    if (step === 'upload') {
      if (uploadedFiles.length > 0 || selectedMediaIds.length > 0) {
        setStep('info')
      }
    } else if (step === 'info') {
      // 必須入力がなくなったので、強制的に次へ進む
      setStep('confirm')
    }
  }

  const prevStep = () => {
    if (step === 'confirm') setStep('info')
    else if (step === 'info') setStep('upload')
  }

  const handleGenerateVideo = async () => {
    try {
      // 生成中の状態を設定
      setIsGenerating(true)
      setGenerationError(null)

      // フォームデータの取得
      const formDataValues = methods.getValues()

      // ユーザーが認証されているか確認
      if (!user || !user.id) {
        throw new Error('ユーザーが認証されていません。再度ログインしてください。')
      }

      // 1. VLog作成APIを呼び出し
      const formData = new FormData()

      // 新規アップロードファイル
      uploadedFiles.forEach(file => {
        formData.append('files', file)
      })

      // 既存メディアIDの追加 (配列の各要素を個別にappend)
      selectedMediaIds.forEach(mediaId => {
        formData.append('mediaIds', mediaId)
      })

      formData.append('title', formDataValues.travelTitle || '')
      formData.append('travelDate', formDataValues.travelDate || '')
      formData.append('destination', formDataValues.travelLocation || '')
      formData.append('theme', 'adventure') // デフォルト

      const res = await apiClient.post('/agent/create-vlog', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      })

      setVlogId(res.data.vlogId)

      // 生成プロセスが開始されたので、進捗表示フェーズへ（実装が必要）
      // ここではひとまず生成中フラグを維持し、SSEによる進捗監視に移行する準備をする
    } catch (error) {
      // エラーハンドリング
      console.error('動画生成中にエラーが発生しました', error)
      const errorMessage =
        error instanceof Error
          ? error.message
          : '動画生成中にエラーが発生しました。もう一度お試しください。'
      setGenerationError(errorMessage)
      setIsGenerating(false)
    }
  }

  const addFiles = (files: File[]) => {
    setUploadedFiles(prev => {
      const newFiles = [...prev, ...files]
      methods.setValue('uploadedFiles', newFiles)
      return newFiles
    })
  }

  const removeFile = (index: number) => {
    setUploadedFiles(prev => {
      const newFiles = [...prev]
      newFiles.splice(index, 1)
      methods.setValue('uploadedFiles', newFiles)
      return newFiles
    })
  }

  const toggleMediaId = (id: string) => {
    setSelectedMediaIds(prev => {
      const next = prev.includes(id) ? prev.filter(i => i !== id) : [...prev, id]
      methods.setValue('mediaIds', next)
      return next
    })
  }

  const value = {
    step,
    setStep,
    nextStep,
    prevStep,
    handleGenerateVideo,
    addFiles,
    removeFile,
    uploadedFiles,
    selectedMediaIds,
    toggleMediaId,
    isGenerating,
    generationError,
    vlogId,
  }

  return (
    <UploadFormContext.Provider value={value}>
      <FormProvider {...methods}>{children}</FormProvider>
    </UploadFormContext.Provider>
  )
}
