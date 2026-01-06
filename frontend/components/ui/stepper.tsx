'use client'

import React from 'react'
import { cn } from '@/lib/utils'
import { Check } from 'lucide-react'

interface StepProps {
  title: string
  description?: string
  index: number
  active: boolean
  completed: boolean
  last?: boolean
}

const Step = ({ title, description, index, active, completed, last = false }: StepProps) => {
  return (
    <div className={cn('flex items-start', !last && 'w-full')}>
      <div className="flex flex-col items-center">
        <div
          className={cn(
            'flex items-center justify-center w-8 h-8 rounded-full border-2 z-10',
            active
              ? 'border-primary bg-primary text-white'
              : completed
              ? 'border-primary bg-primary text-white'
              : 'border-muted-foreground/40 text-muted-foreground'
          )}
        >
          {completed ? (
            <Check className="w-5 h-5" />
          ) : (
            <span className="text-sm font-medium">{index + 1}</span>
          )}
        </div>
        {!last && (
          <div
            className={cn('w-0.5 h-full mt-2', completed ? 'bg-primary' : 'bg-muted-foreground/40')}
          />
        )}
      </div>
      <div className="ml-4 mt-0.5 pb-8">
        <h3
          className={cn(
            'text-sm font-medium',
            active || completed ? 'text-foreground' : 'text-muted-foreground'
          )}
        >
          {title}
        </h3>
        {description && <p className="text-xs mt-1 text-muted-foreground">{description}</p>}
      </div>
    </div>
  )
}

interface StepperProps {
  steps: {
    title: string
    description?: string
  }[]
  activeStep: number
  className?: string
}

// モバイル用の水平ステッパーのコンポーネント
const MobileStep = ({ title, index, active, completed }: StepProps) => {
  return (
    <div className="flex flex-col items-center z-10">
      <div
        className={cn(
          'flex items-center justify-center rounded-full border-2 mb-1.5 transition-all duration-300 ease-in-out bg-card',
          active
            ? 'border-primary bg-primary text-white w-10 h-10 scale-110 shadow-md'
            : completed
            ? 'border-primary bg-primary text-white w-8 h-8'
            : 'border-muted-foreground/40 text-muted-foreground w-8 h-8',
          active && 'animate-pulse-light'
        )}
      >
        {completed ? (
          <Check className={cn('transition-all duration-300', active ? 'w-5 h-5' : 'w-4 h-4')} />
        ) : (
          <span
            className={cn(
              'font-medium transition-all duration-300',
              active ? 'text-sm' : 'text-xs'
            )}
          >
            {index + 1}
          </span>
        )}
      </div>
      <span
        className={cn(
          'text-center font-medium transition-all duration-300 text-wrap px-1 max-w-16',
          active ? 'text-sm text-foreground' : 'text-xs text-muted-foreground'
        )}
      >
        {title}
      </span>
    </div>
  )
}

export function Stepper({ steps, activeStep, className }: StepperProps) {
  return (
    <>
      {/* モバイル用の水平ステッパー */}
      <div
        className="md:hidden grid w-full pb-2"
        style={{ gridTemplateColumns: `repeat(${steps.length}, 1fr)` }}
      >
        {steps.map((step, index) => (
          <div key={index} className="flex flex-col items-center relative">
            <MobileStep
              title={step.title}
              description={step.description}
              index={index}
              active={index === activeStep}
              completed={index < activeStep}
              last={index === steps.length - 1}
            />

            {/* 接続線 - 各ステップの間に表示 */}
            {index < steps.length - 1 && (
              <div
                className={cn(
                  'absolute top-4 h-0.5 transition-all duration-500 w-full',
                  index < activeStep ? 'bg-primary' : 'bg-muted-foreground/40'
                )}
                style={{
                  left: '50%', // 左端を中央に
                  width: '100%', // 幅を100%に
                  transform: index === activeStep - 1 ? 'scaleX(1.1)' : 'scaleX(1)',
                  transformOrigin: 'left', // 変形の基点を左に
                  zIndex: 0,
                }}
              />
            )}
          </div>
        ))}
      </div>

      {/* デスクトップ用の縦ステッパー */}
      <div className={cn('hidden md:flex md:flex-col', className)}>
        {steps.map((step, index) => (
          <Step
            key={index}
            title={step.title}
            description={step.description}
            index={index}
            active={index === activeStep}
            completed={index < activeStep}
            last={index === steps.length - 1}
          />
        ))}
      </div>
    </>
  )
}
