'use client'

import React from 'react'
import { useFormContext } from 'react-hook-form'
import StepperLayout from '@/components/layout/StepperLayout'
import PhotoUpload from './PhotoUpload'
import TravelInfo from './TravelInfo'
import VideoGenerationConfirm from './VideoGenerationConfirm'
import { UploadFormProvider, useUploadForm } from './UploadFormContext'

function UploadStepper() {
  const { step, nextStep, prevStep, uploadedFiles } = useUploadForm()
  const { formState } = useFormContext()

  const steps = [
    {
      title: '写真のアップロード',
      description: '旅行の思い出の写真を選択',
      content: <PhotoUpload />,
      isNextDisabled: uploadedFiles.length === 0,
    },
    {
      title: '旅行情報の入力',
      description: '動画に含める旅行の詳細',
      content: <TravelInfo />,
      isNextDisabled:
        !formState.isValid ||
        Boolean(formState.errors.travelTitle) ||
        Boolean(formState.errors.travelDate),
    },
    {
      title: '動画生成',
      description: '確認して生成を開始',
      content: <VideoGenerationConfirm />,
      showNextButton: false,
    },
  ]

  const activeStepIndex = step === 'upload' ? 0 : step === 'info' ? 1 : 2

  return (
    <StepperLayout steps={steps} activeStep={activeStepIndex} onNext={nextStep} onBack={prevStep} />
  )
}

export default function StepperWrapper() {
  return (
    <UploadFormProvider>
      <UploadStepper />
    </UploadFormProvider>
  )
}
