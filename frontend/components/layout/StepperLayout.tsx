'use client'

import React, { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Stepper } from '@/components/ui/stepper'
import { ArrowLeft, ArrowRight } from 'lucide-react'

interface StepperLayoutProps {
  steps: {
    title: string
    description?: string
    content: React.ReactNode
    showNextButton?: boolean
    showBackButton?: boolean
    nextButtonText?: string
    backButtonText?: string
    isNextDisabled?: boolean
  }[]
  onComplete?: () => void
  className?: string
  activeStep?: number
  onNext?: () => void
  onBack?: () => void
}

export default function StepperLayout({
  steps,
  onComplete,
  className,
  activeStep: externalActiveStep,
  onNext,
  onBack,
}: StepperLayoutProps) {
  const [internalActiveStep, setInternalActiveStep] = useState(0)

  // 外部から制御されているかどうかを判断
  const isControlled = externalActiveStep !== undefined
  const activeStep = isControlled ? externalActiveStep : internalActiveStep

  const handleNext = () => {
    if (onNext) {
      onNext()
    } else if (activeStep < steps.length - 1) {
      setInternalActiveStep(prev => prev + 1)
    } else if (onComplete) {
      onComplete()
    }
  }

  const handleBack = () => {
    if (onBack) {
      onBack()
    } else if (activeStep > 0) {
      setInternalActiveStep(prev => prev - 1)
    }
  }

  const currentStep = steps[activeStep]

  return (
    <div className="flex flex-col md:flex-row gap-4 md:gap-8 pb-16 md:pb-0">
      <div className="w-full md:w-1/4">
        <Card className="p-4 md:p-6 md:sticky md:top-24 overflow-hidden">
          <div className="relative">
            {/* ステッパーコンポーネント */}
            <Stepper
              steps={steps.map(step => ({
                title: step.title,
                description: step.description,
              }))}
              activeStep={activeStep}
            />
          </div>
        </Card>
      </div>
      <div className="w-full md:w-3/4">
        <Card className="p-4 md:p-6">
          <div className="mb-6 md:mb-8">{currentStep.content}</div>
          {/* デスクトップ用のナビゲーションボタン */}
          <div className="hidden md:flex justify-between mt-6 md:mt-8">
            {(currentStep.showBackButton ?? true) && activeStep > 0 ? (
              <Button
                variant="outline"
                onClick={handleBack}
                className="flex items-center gap-1 md:gap-2 px-2 md:px-4 py-1 md:py-2 text-sm md:text-base"
              >
                <ArrowLeft className="w-3 h-3 md:w-4 md:h-4" />
                {currentStep.backButtonText || '前へ戻る'}
              </Button>
            ) : (
              <div></div>
            )}
            {(currentStep.showNextButton ?? true) && (
              <Button
                onClick={handleNext}
                className="flex items-center gap-1 md:gap-2 px-2 md:px-4 py-1 md:py-2 text-sm md:text-base"
                disabled={currentStep.isNextDisabled}
              >
                {currentStep.nextButtonText || '次へ進む'}
                <ArrowRight className="w-3 h-3 md:w-4 md:h-4" />
              </Button>
            )}
          </div>
        </Card>
      </div>

      {/* モバイル用の固定ナビゲーションバー */}
      <div className="fixed bottom-0 left-0 right-0 bg-background border-t border-border p-3 md:hidden z-50">
        <div className="flex justify-between items-center container mx-auto">
          {(currentStep.showBackButton ?? true) && activeStep > 0 ? (
            <Button
              variant="outline"
              onClick={handleBack}
              className="flex items-center gap-1 px-3 py-2 text-sm"
            >
              <ArrowLeft className="w-3 h-3" />
              {currentStep.backButtonText || '前へ'}
            </Button>
          ) : (
            <div></div>
          )}
          <div className="text-xs text-muted-foreground">
            {activeStep + 1} / {steps.length}
          </div>
          {(currentStep.showNextButton ?? true) && (
            <Button
              onClick={handleNext}
              className="flex items-center gap-1 px-3 py-2 text-sm"
              disabled={currentStep.isNextDisabled}
            >
              {currentStep.nextButtonText || '次へ'}
              <ArrowRight className="w-3 h-3" />
            </Button>
          )}
        </div>
      </div>
    </div>
  )
}
