'use client'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Play, Calendar, Download, Share2 } from 'lucide-react'
import React from 'react'
import { useGetTravelsByUserId, useUpdateTravel } from '@/api/travelApi'
import { useAuth } from '@/context/authContext'
import { Travel, TravelInput } from '@/api/types'
import TravelDetailSheet from './TravelDetailSheet'

export default function RecentVideo() {
  const { user, loading } = useAuth()
  const userId = user?.id || ''
  const [selectedTravel, setSelectedTravel] = React.useState<Travel | null>(null)
  const [isSheetOpen, setIsSheetOpen] = React.useState(false)
  const { mutate: updateTravel } = useUpdateTravel(selectedTravel?.id || '')

  const { data: travelsData, isLoading: isTravelsLoading } = useGetTravelsByUserId(userId)
  const travels = travelsData?.travels || []

  // 最新の3件のみ表示（更新日時でソート）
  const recentTravels = [...travels]
    .sort((a, b) => new Date(b.updated_at).getTime() - new Date(a.updated_at).getTime())
    .slice(0, 3)

  const isLoading = loading || isTravelsLoading

  // 旅行をクリックしたときの処理
  const handleTravelClick = (travel: Travel) => {
    setSelectedTravel(travel)
    setIsSheetOpen(true)
  }

  // 期間から動画の長さを簡易計算（実際は動画の長さはAPIから取得するべき）
  const calculateDuration = (startDate: string, endDate: string): string => {
    const start = new Date(startDate)
    const end = new Date(endDate)
    const days = Math.ceil((end.getTime() - start.getTime()) / (1000 * 60 * 60 * 24))
    // 分と秒の表記として、日数に基づいて簡易的に生成
    return `${Math.max(1, Math.min(5, days))}:${String(Math.floor(Math.random() * 60)).padStart(
      2,
      '0'
    )}`
  }

  const onSave = async (updatedTravel: Partial<TravelInput>) => {
    await updateTravel({
      ...updatedTravel,
    })
  }

  return (
    <div className="mb-6">
      <div className="flex items-center justify-between mb-3 md:mb-4">
        <h2 className="text-lg md:text-xl font-semibold text-foreground">最近の動画</h2>
        <Button variant="outline" size="sm" className="text-xs md:text-sm h-8 px-3">
          すべて表示
        </Button>
      </div>

      {isLoading ? (
        <div className="flex justify-center py-8">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
        </div>
      ) : recentTravels.length === 0 ? (
        <div className="text-center py-8 text-muted-foreground">
          動画がありません。新しい旅行を追加してみましょう！
        </div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-3 md:gap-4">
          {recentTravels.map(travel => (
            <div
              key={travel.id}
              className="relative aspect-[4/3] bg-gradient-to-br from-primary/10 to-secondary/10 overflow-hidden cursor-pointer group"
              onClick={() => handleTravelClick(travel)}
            >
              <img
                src={travel.thumbnail || '/placeholder.webp'}
                alt={travel.title}
                className="w-full h-full object-cover"
                loading="lazy"
              />
              {/* Gradient overlay for better text visibility */}
              <div className="absolute inset-0 bg-gradient-to-t from-black/70 to-transparent opacity-70"></div>
              {/* Content overlay (positioned at the bottom only) */}
              <div className="absolute bottom-0 left-0 right-0 p-2">
                {/* Top-right badge for duration */}
                <Badge className="absolute top-2 right-2 text-2xs md:text-xs bg-black/60 text-white">
                  {calculateDuration(travel.startDate, travel.endDate)}
                </Badge>

                {/* Bottom content */}
                <div className="text-white">
                  <h3 className="font-medium text-sm md:text-base line-clamp-1 mb-0.5">
                    {travel.title}
                  </h3>
                  <div className="flex items-center text-2xs md:text-xs text-white/80">
                    <div className="flex items-center">
                      <Calendar className="w-3 h-3 mr-1" />
                      {new Date(travel.startDate).toLocaleDateString('ja-JP', {
                        year: 'numeric',
                        month: 'long',
                        day: 'numeric',
                      })}
                    </div>
                  </div>
                </div>
              </div>

              {/* Play button overlay */}
              <div className="absolute inset-0 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
                <Button
                  size="sm"
                  variant="secondary"
                  className="bg-white/90 text-black hover:bg-white rounded-full w-9 h-9 md:w-10 md:h-10 p-0"
                >
                  <Play className="w-4 h-4 md:w-5 md:h-5 ml-0.5" fill="currentColor" />
                </Button>
              </div>

              {/* Action buttons (visible on hover) - positioned at the top right */}
              <div className="absolute top-2 right-2 flex flex-col gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                <Button
                  size="sm"
                  variant="outline"
                  className="p-1 h-7 w-7 md:h-8 md:w-8 bg-black/40 hover:bg-black/60 text-white border-white/20 rounded-full"
                >
                  <Download className="w-3.5 h-3.5" />
                </Button>
                <Button
                  size="sm"
                  variant="outline"
                  className="p-1 h-7 w-7 md:h-8 md:w-8 bg-black/40 hover:bg-black/60 text-white border-white/20 rounded-full"
                >
                  <Share2 className="w-3.5 h-3.5" />
                </Button>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* 旅行詳細シート */}
      <TravelDetailSheet
        travel={selectedTravel}
        isOpen={isSheetOpen}
        onOpenChange={setIsSheetOpen}
        onSave={onSave}
      />
    </div>
  )
}
