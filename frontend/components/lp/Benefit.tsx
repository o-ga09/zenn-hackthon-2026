'use client'
import { Heart, Share2, Star, ChevronLeft, ChevronRight } from 'lucide-react'
import React, { useState, useEffect } from 'react'
import { anyOneVideos } from '@/mock/mock'
import Image from 'next/image'
import { motion, useAnimationControls } from 'framer-motion'

// カードアニメーションバリアント
const cardVariants = {
  hidden: { opacity: 0, y: 20 },
  visible: { opacity: 1, y: 0 },
}

// スクロール制御
const scrollbarHideClass = 'scrollbar-hide'
const scrollbarStyles = {
  msOverflowStyle: 'none',
  scrollbarWidth: 'none',
  WebkitOverflowScrolling: 'touch',
} as const

export default function Benefit() {
  // 表示するカードのインデックスを管理
  const [currentIndex, setCurrentIndex] = useState(0)
  // 自動回転の一時停止状態を管理
  const [isPaused, setIsPaused] = useState(false)
  // アニメーション制御
  const controls = useAnimationControls()

  // カードを次に進める関数
  const nextCard = () => {
    controls.start({ opacity: 0, scale: 0.8 }).then(() => {
      setCurrentIndex(prev => (prev + 1) % anyOneVideos.length)
      controls.start({ opacity: 1, scale: 1 })
    })
  }

  // カードを前に戻す関数
  const prevCard = () => {
    controls.start({ opacity: 0, scale: 0.8 }).then(() => {
      setCurrentIndex(prev => (prev - 1 + anyOneVideos.length) % anyOneVideos.length)
      controls.start({ opacity: 1, scale: 1 })
    })
  }

  // 自動切り替えの制御
  useEffect(() => {
    // 一時停止中は自動切り替えしない
    if (isPaused) return

    // 表示するカードを自動的に変更するタイマー
    const timer = setInterval(() => {
      nextCard()
    }, 3000)

    // コンポーネントのクリーンアップ時やポーズ状態変更時にタイマーを解除
    return () => clearInterval(timer)
  }, [isPaused])

  return (
    <section className="container mx-auto px-4 min-h-screen flex items-center py-10 sm:py-16">
      <div className="grid lg:grid-cols-2 gap-8 sm:gap-12 items-center max-w-5xl mx-auto">
        <div className="order-2 lg:order-1">
          <h3 className="text-2xl sm:text-3xl font-bold text-foreground mb-4 sm:mb-6 text-center lg:text-left">
            なぜTravelMomentsを選ぶのか？
          </h3>
          <div className="grid sm:grid-cols-1 gap-4 sm:gap-6">
            <motion.div
              className="flex items-start space-x-3 sm:space-x-4 p-3 sm:p-4 rounded-lg hover:bg-background/50 transition-colors"
              initial={{ opacity: 0, y: 10 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ delay: 0.1 }}
            >
              <div className="w-8 h-8 sm:w-10 sm:h-10 bg-red-100 rounded-full flex items-center justify-center flex-shrink-0 mt-1 shadow-sm">
                <Heart className="w-4 h-4 sm:w-5 sm:h-5 text-red-500" strokeWidth={2} />
              </div>
              <div>
                <h4 className="font-semibold text-foreground mb-1 sm:mb-2 text-lg">
                  エモーショナルな体験
                </h4>
                <p className="text-sm sm:text-base text-muted-foreground">
                  AIが写真から感情を読み取り、その瞬間にぴったりの音楽とエフェクトを選択します。
                </p>
              </div>
            </motion.div>

            <motion.div
              className="flex items-start space-x-3 sm:space-x-4 p-3 sm:p-4 rounded-lg hover:bg-background/50 transition-colors"
              initial={{ opacity: 0, y: 10 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ delay: 0.2 }}
            >
              <div className="w-8 h-8 sm:w-10 sm:h-10 bg-blue-100 rounded-full flex items-center justify-center flex-shrink-0 mt-1 shadow-sm">
                <Share2 className="w-4 h-4 sm:w-5 sm:h-5 text-blue-500" strokeWidth={2} />
              </div>
              <div>
                <h4 className="font-semibold text-foreground mb-1 sm:mb-2 text-lg">簡単シェア</h4>
                <p className="text-sm sm:text-base text-muted-foreground">
                  YouTube ShortsやInstagramストーリーズに最適化された縦画面動画を生成します。
                </p>
              </div>
            </motion.div>

            <motion.div
              className="flex items-start space-x-3 sm:space-x-4 p-3 sm:p-4 rounded-lg hover:bg-background/50 transition-colors"
              initial={{ opacity: 0, y: 10 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ delay: 0.3 }}
            >
              <div className="w-8 h-8 sm:w-10 sm:h-10 bg-amber-100 rounded-full flex items-center justify-center flex-shrink-0 mt-1 shadow-sm">
                <Star className="w-4 h-4 sm:w-5 sm:h-5 text-amber-500" strokeWidth={2} />
              </div>
              <div>
                <h4 className="font-semibold text-foreground mb-1 sm:mb-2 text-lg">プロ品質</h4>
                <p className="text-sm sm:text-base text-muted-foreground">
                  高度なAI技術により、プロが作ったような高品質な動画を自動生成します。
                </p>
              </div>
            </motion.div>
          </div>
        </div>

        <div className="relative order-1 lg:order-2">
          <div className="relative overflow-hidden py-6 sm:py-8">
            {/* メインのカードショーケース */}
            <div
              className="flex justify-center"
              onMouseEnter={() => setIsPaused(true)}
              onMouseLeave={() => setIsPaused(false)}
              onTouchStart={() => setIsPaused(true)}
              onTouchEnd={() => setIsPaused(false)}
            >
              <motion.div
                className="w-[240px] sm:w-[280px] mx-auto relative"
                animate={controls}
                initial={{ opacity: 0, scale: 0.8 }}
                transition={{
                  type: 'spring',
                  stiffness: 300,
                  damping: 20,
                  duration: 0.4,
                }}
              >
                <div className="aspect-[9/16] bg-gradient-to-br from-primary/10 to-secondary/10 rounded-xl sm:rounded-2xl overflow-hidden relative shadow-lg">
                  <Image
                    src={anyOneVideos[currentIndex].thumbnail || '/placeholder.webp'}
                    alt={anyOneVideos[currentIndex].title}
                    fill
                    className="object-cover"
                    sizes="(max-width: 768px) 100vw, 280px"
                    priority
                  />

                  <motion.div
                    className="absolute inset-0 bg-gradient-to-b from-black/10 to-black/80"
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                  />

                  <motion.div
                    className="absolute bottom-0 left-0 right-0 p-4"
                    initial={{ y: 20, opacity: 0 }}
                    animate={{ y: 0, opacity: 1 }}
                    transition={{ delay: 0.2 }}
                  >
                    <h4 className="text-white text-lg font-medium mb-1">
                      {anyOneVideos[currentIndex].title}
                    </h4>
                    <p className="text-white/90 text-sm">{anyOneVideos[currentIndex].userName}</p>
                  </motion.div>
                </div>
              </motion.div>
            </div>

            {/* ページネーションインジケーターと矢印 */}
            <div className="flex items-center justify-center mt-6 space-x-4">
              {/* 左矢印 */}
              <motion.button
                className="bg-white rounded-full p-2 shadow-md hover:bg-gray-50 transition-colors"
                onClick={prevCard}
                whileHover={{ scale: 1.1 }}
                whileTap={{ scale: 0.95 }}
                aria-label="前へ"
              >
                <ChevronLeft className="w-5 h-5 text-gray-600" />
              </motion.button>

              {/* ページネーションドット */}
              <div className="flex space-x-2">
                {anyOneVideos.map((_, index) => (
                  <motion.button
                    key={index}
                    className={`w-2 h-2 rounded-full ${
                      index === currentIndex ? 'bg-primary' : 'bg-gray-300'
                    }`}
                    onClick={() => {
                      controls.start({ opacity: 0, scale: 0.8 }).then(() => {
                        setCurrentIndex(index)
                        controls.start({ opacity: 1, scale: 1 })
                      })
                    }}
                    whileHover={{ scale: 1.5 }}
                    initial={{ opacity: 0, scale: 0 }}
                    animate={{ opacity: 1, scale: 1 }}
                    transition={{ delay: index * 0.05 }}
                  />
                ))}
              </div>

              {/* 右矢印 */}
              <motion.button
                className="bg-white rounded-full p-2 shadow-md hover:bg-gray-50 transition-colors"
                onClick={nextCard}
                whileHover={{ scale: 1.1 }}
                whileTap={{ scale: 0.95 }}
                aria-label="次へ"
              >
                <ChevronRight className="w-5 h-5 text-gray-600" />
              </motion.button>
            </div>
          </div>
        </div>
      </div>
    </section>
  )
}
