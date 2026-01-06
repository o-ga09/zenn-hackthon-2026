'use client'
import React, { useRef, useEffect, useState } from 'react'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Heart, MessageCircle, Share2, User, Volume2, VolumeX } from 'lucide-react'

interface Video {
  id: number
  title: string
  description: string
  date: string
  duration: string
  videoUrl: string
  thumbnail: string
  userName: string
  userAvatar?: string
  likes: number
  comments: number
}

interface ShortsVideoCardProps {
  video: Video
  isActive: boolean
}

export default function ShortsVideoCard({ video, isActive }: ShortsVideoCardProps) {
  const videoRef = useRef<HTMLVideoElement>(null)
  const [isMuted, setIsMuted] = useState(true)

  // 表示されたらビデオを再生、非表示になったら一時停止
  useEffect(() => {
    if (!videoRef.current) return

    if (isActive) {
      videoRef.current.currentTime = 0
      videoRef.current.play().catch(error => {
        console.error('自動再生エラー:', error)
      })
    } else {
      videoRef.current.pause()
    }
  }, [isActive])

  const toggleMute = () => {
    if (videoRef.current) {
      videoRef.current.muted = !videoRef.current.muted
      setIsMuted(!isMuted)
    }
  }

  return (
    <div className="relative h-full w-full bg-black flex items-center justify-center">
      {/* ビデオ */}
      <div className="relative w-full h-full">
        {isActive ? (
          <video
            ref={videoRef}
            src={video.videoUrl}
            className="w-full h-full object-cover"
            loop
            playsInline
            muted={isMuted}
            poster={video.thumbnail}
          />
        ) : (
          <img src={video.thumbnail} alt={video.title} className="w-full h-full object-cover" />
        )}

        {/* タップエリア - ビデオを一時停止/再生 */}
        <div
          className="absolute inset-0 z-10"
          onClick={() => {
            if (videoRef.current) {
              if (videoRef.current.paused) {
                videoRef.current.play()
              } else {
                videoRef.current.pause()
              }
            }
          }}
        />

        {/* オーバーレイグラデーション */}
        <div className="absolute inset-0 bg-gradient-to-b from-black/30 via-transparent to-black/50 pointer-events-none"></div>
      </div>

      {/* 動画情報 - 下部 */}
      <div className="absolute bottom-4 left-4 right-16 text-white z-20 pointer-events-none">
        <h3 className="text-lg font-bold mb-1">{video.title}</h3>
        <p className="text-sm mb-2 line-clamp-2">{video.description}</p>

        {/* ユーザー情報 */}
        <div className="flex items-center gap-2 pointer-events-auto">
          <div className="h-9 w-9 rounded-full bg-gray-500 flex items-center justify-center overflow-hidden">
            {video.userAvatar ? (
              <img
                src={video.userAvatar}
                alt={video.userName}
                className="w-full h-full object-cover"
              />
            ) : (
              <User className="h-5 w-5 text-white" />
            )}
          </div>
          <span className="font-medium text-sm">{video.userName}</span>
          <Button
            size="sm"
            variant="secondary"
            className="ml-2 h-8 text-xs px-3 py-1 rounded-full pointer-events-auto"
          >
            フォロー
          </Button>
        </div>
      </div>

      {/* アクションボタン - 右側 */}
      <div className="absolute right-4 bottom-24 flex flex-col items-center gap-6 z-20">
        {/* いいねボタン */}
        <div className="flex flex-col items-center">
          <Button
            size="sm"
            variant="ghost"
            className="h-12 w-12 rounded-full bg-black/30 text-white p-0 pointer-events-auto"
          >
            <Heart className="h-7 w-7" />
          </Button>
          <span className="text-white text-xs mt-1">{video.likes}</span>
        </div>

        {/* コメントボタン */}
        <div className="flex flex-col items-center">
          <Button
            size="sm"
            variant="ghost"
            className="h-12 w-12 rounded-full bg-black/30 text-white p-0 pointer-events-auto"
          >
            <MessageCircle className="h-7 w-7" />
          </Button>
          <span className="text-white text-xs mt-1">{video.comments}</span>
        </div>

        {/* 共有ボタン */}
        <div className="flex flex-col items-center">
          <Button
            size="sm"
            variant="ghost"
            className="h-12 w-12 rounded-full bg-black/30 text-white p-0 pointer-events-auto"
          >
            <Share2 className="h-7 w-7" />
          </Button>
          <span className="text-white text-xs mt-1">共有</span>
        </div>
      </div>

      {/* 動画コントロール - 上部 */}
      <div className="absolute top-4 right-4 flex items-center gap-2 z-20">
        {/* 音声コントロール */}
        <Button
          size="sm"
          variant="ghost"
          className="h-10 w-10 rounded-full bg-black/50 text-white p-0 pointer-events-auto"
          onClick={toggleMute}
        >
          {isMuted ? <VolumeX className="h-5 w-5" /> : <Volume2 className="h-5 w-5" />}
        </Button>

        {/* 動画時間 */}
        <Badge className="bg-black/60 text-white">{video.duration}</Badge>
      </div>

      {/* 関連動画表示部分（YouTube Shortsスタイルの下部ナビゲーション） */}
      <div className="absolute bottom-0 left-0 right-0 h-1 bg-gray-700 z-20">
        <div
          className="absolute bottom-0 left-0 h-1 bg-white w-1/3"
          style={{
            width: `${Math.random() * 100}%`, // 実際にはビデオの進行状態に応じて調整
          }}
        ></div>
      </div>
    </div>
  )
}
