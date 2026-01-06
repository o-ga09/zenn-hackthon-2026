'use client'
import React, { useState, useRef, useEffect } from 'react'
import ShortsVideoCard from './ShortsVideoCard'
import { anyOneVideos } from '@/mock/mock'
import './shortsviewer.css' // スクロールバー非表示用のCSSをインポート

export default function ShortsViewer() {
  const [currentVideoIndex, setCurrentVideoIndex] = useState(0)
  const containerRef = useRef<HTMLDivElement>(null)

  // ビデオのインデックスを変更する関数
  const handleScroll = () => {
    if (!containerRef.current) return

    const container = containerRef.current
    const containerHeight = container.clientHeight
    const scrollTop = container.scrollTop

    // スクロール位置に基づいて現在のビデオインデックスを計算
    const newIndex = Math.round(scrollTop / containerHeight)

    if (newIndex !== currentVideoIndex && newIndex >= 0 && newIndex < anyOneVideos.length) {
      setCurrentVideoIndex(newIndex)
    }
  }

  // スクロールイベントリスナーを設定
  useEffect(() => {
    const container = containerRef.current
    if (!container) return

    container.addEventListener('scroll', handleScroll)
    return () => {
      container.removeEventListener('scroll', handleScroll)
    }
  }, [currentVideoIndex])

  return (
    <div className="w-full flex justify-center bg-white md:py-6 h-[calc(100vh-64px)] overflow-hidden">
      {/* PC版の場合は白背景でコンテナを表示 */}
      <div className="md:bg-black md:rounded-xl md:overflow-hidden md:shadow-xl w-full max-w-[420px] h-full">
        <div
          ref={containerRef}
          className="h-full w-full overflow-y-auto snap-y snap-mandatory scrollbar-hide bg-black shorts-container"
          style={{ scrollSnapType: 'y mandatory', msOverflowStyle: 'none', scrollbarWidth: 'none' }}
        >
          {anyOneVideos.map((video, index) => (
            <div key={video.id} className="h-full snap-start">
              <ShortsVideoCard video={video} isActive={index === currentVideoIndex} />
            </div>
          ))}

          {/* ナビゲーションヒント - 初回のみ表示するようにしたい場合は状態管理が必要 */}
          <div className="absolute bottom-20 left-1/2 transform -translate-x-1/2 bg-black/70 text-white text-xs px-3 py-1 rounded-full pointer-events-none z-50 opacity-80">
            上下にスワイプして動画を切り替え
          </div>
        </div>
      </div>
    </div>
  )
}
