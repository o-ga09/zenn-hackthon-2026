'use client'
import React, { ReactNode } from 'react'
import CommonHeader from '@/components/header/CommonHeader'
import CommonFooter from '@/components/footer/CommonFooter'
import { cn } from '@/lib/utils'

interface MainLayoutProps {
  children: ReactNode
  description?: string
  className?: string
  showHeader?: boolean
  showFooter?: boolean
  showTitle?: boolean
  bgClassName?: string
  disableBodyScroll?: boolean
}

export default function MainLayout({
  children,
  description,
  className,
  showHeader = true,
  showFooter = true,
  bgClassName = 'bg-gradient-to-br from-pink-100 via-purple-50 to-blue-100',
  disableBodyScroll = false,
}: MainLayoutProps) {
  // ボディのスクロールを無効化するエフェクト
  React.useEffect(() => {
    if (disableBodyScroll) {
      document.body.style.overflow = 'hidden'
    } else {
      document.body.style.overflow = ''
    }

    return () => {
      document.body.style.overflow = ''
    }
  }, [disableBodyScroll])

  return (
    <div
      className={`flex flex-col min-h-screen ${bgClassName} ${
        disableBodyScroll ? 'overflow-hidden' : ''
      }`}
    >
      {showHeader && <CommonHeader />}
      {/* Main Content */}
      <main className={cn('container mx-auto px-4 py-6 md:py-8 flex-grow', className)}>
        {children}
      </main>
      {showFooter && <CommonFooter />}
    </div>
  )
}
