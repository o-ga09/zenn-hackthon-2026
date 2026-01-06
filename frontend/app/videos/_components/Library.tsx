'use client'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import {
  Search,
  Filter,
  Play,
  Clock,
  Calendar,
  Download,
  Share2,
  Trash2,
  CheckCircle,
  X,
} from 'lucide-react'
import React, { useState, useEffect, useMemo } from 'react'
import { useGetTravelsByUserId } from '@/api/travelApi'
import { useAuth } from '@/context/authContext'

// フィルタータイプの定義
type FilterStatus = '全て' | '完成' | '処理中'
type SortOrder = '新しい順' | '古い順' | '人気順'

export default function Library() {
  // ステート設定
  const [searchTerm, setSearchTerm] = useState('')
  const [displayCount, setDisplayCount] = useState(6)
  const [statusFilter, setStatusFilter] = useState<FilterStatus>('全て')
  const [sortOrder, setSortOrder] = useState<SortOrder>('新しい順')
  const [showFilterMenu, setShowFilterMenu] = useState(false)

  // ユーザー情報取得
  const { user, loading } = useAuth()
  const userId = user?.id || ''

  // 旅行データ取得
  const { data: travelsData, isLoading: isTravelsLoading } = useGetTravelsByUserId(userId)
  const allTravels = travelsData?.travels || []
  console.log(travelsData?.travels.length)

  // APIデータをビデオリスト形式に変換
  const videosFromApi = useMemo(() => {
    return allTravels.map(travel => ({
      id: travel.id,
      title: travel.title,
      date: new Date(travel.startDate).toLocaleDateString('ja-JP', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
      }),
      thumbnail: travel.thumbnail || '/placeholder.webp',
      duration: calculateDuration(travel.startDate, travel.endDate),
      status: '完成', // APIからステータスを取得できるようになったら変更する
      views: Math.floor(Math.random() * 500), // APIからビュー数を取得できるようになったら変更する
      likes: Math.floor(Math.random() * 100), // APIからいいね数を取得できるようになったら変更する
      startDate: travel.startDate,
      endDate: travel.endDate,
      description: travel.description,
    }))
  }, [allTravels])

  // 期間から動画の長さを簡易計算
  function calculateDuration(startDate: string, endDate: string): string {
    const start = new Date(startDate)
    const end = new Date(endDate)
    const days = Math.ceil((end.getTime() - start.getTime()) / (1000 * 60 * 60 * 24))
    return `${Math.max(1, Math.min(5, days))}:${String(Math.floor(Math.random() * 60)).padStart(
      2,
      '0'
    )}`
  }

  // フィルター適用したビデオリスト
  const filteredVideos = useMemo(() => {
    // 1. 検索条件を適用
    let filtered = videosFromApi.filter(
      video =>
        video.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
        video.date.includes(searchTerm)
    )

    // 2. ステータスフィルターを適用
    if (statusFilter !== '全て') {
      filtered = filtered.filter(video => video.status === statusFilter)
    }

    // 3. ソート順を適用
    return filtered.sort((a, b) => {
      if (sortOrder === '新しい順') {
        // 日付を比較（新しい順）- 簡易実装のため日付文字列をそのまま比較
        return b.date.localeCompare(a.date)
      } else if (sortOrder === '古い順') {
        // 日付を比較（古い順）
        return a.date.localeCompare(b.date)
      } else {
        // 人気順（いいね数）
        return (b.likes || 0) - (a.likes || 0)
      }
    })
  }, [searchTerm, statusFilter, sortOrder])

  // 表示する動画
  const displayedVideos = useMemo(() => {
    return filteredVideos.slice(0, displayCount)
  }, [filteredVideos, displayCount])

  // ローディング状態
  const isLoading = loading || isTravelsLoading

  // さらに読み込む機能
  const handleLoadMore = () => {
    setDisplayCount(prevCount => prevCount + 6)
  }

  // 検索入力の処理
  const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSearchTerm(e.target.value)
    // 検索時は表示件数をリセット
    setDisplayCount(6)
  }

  // フィルター適用
  const applyFilter = (status: FilterStatus) => {
    setStatusFilter(status)
    setShowFilterMenu(false)
    setDisplayCount(6) // フィルター変更時にリセット
  }

  // ソート順変更
  const applySort = (order: SortOrder) => {
    setSortOrder(order)
    setShowFilterMenu(false)
    setDisplayCount(6) // ソート変更時にリセット
  }

  // フィルターメニュー外のクリックを検知して閉じる
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      const target = event.target as HTMLElement
      if (showFilterMenu && !target.closest('.filter-container')) {
        setShowFilterMenu(false)
      }
    }

    document.addEventListener('mousedown', handleClickOutside)
    return () => {
      document.removeEventListener('mousedown', handleClickOutside)
    }
  }, [showFilterMenu])

  return (
    <>
      <Card className="border-0 shadow-md bg-card/50 backdrop-blur-sm mb-6 relative z-[996]">
        <CardContent className="p-4">
          <div className="flex flex-col md:flex-row gap-3">
            <div className="flex-1 relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground w-4 h-4" />
              <Input
                placeholder="動画を検索..."
                className="pl-10 h-10"
                value={searchTerm}
                onChange={handleSearchChange}
              />
            </div>
            <div className="relative filter-container z-[997]">
              <Button
                variant="outline"
                className="flex items-center space-x-2 bg-transparent h-10 w-full md:w-auto justify-center md:justify-start relative z-[997]"
                onClick={() => setShowFilterMenu(!showFilterMenu)}
              >
                <Filter className="w-4 h-4" />
                <span>フィルター</span>
                {statusFilter !== '全て' && (
                  <Badge className="ml-2 bg-primary text-white">{statusFilter}</Badge>
                )}
              </Button>

              {/* フィルターメニュー */}
              {showFilterMenu && (
                <>
                  {/* オーバーレイ背景 - モバイル表示のみ */}
                  <div
                    className="fixed inset-0 bg-black/50 z-[998] md:hidden"
                    onClick={() => setShowFilterMenu(false)}
                  ></div>

                  {/* フィルターメニュー本体 */}
                  <div className="fixed md:absolute right-4 md:right-0 left-4 md:left-auto top-[4rem] md:top-full mt-1 w-auto md:w-56 bg-white shadow-lg rounded-md p-2 z-[999] border max-h-[calc(100vh-8rem)] overflow-auto">
                    <div className="sticky top-0 bg-white pt-1 text-sm font-medium text-gray-600 pb-1 mb-1 border-b flex justify-between items-center">
                      <span>フィルター設定</span>
                      <button
                        onClick={() => setShowFilterMenu(false)}
                        className="md:hidden text-gray-500 hover:text-gray-700"
                      >
                        <X className="w-4 h-4" />
                      </button>
                    </div>

                    <div className="text-sm font-medium text-gray-600 pb-1 mb-1">ステータス</div>
                    <div className="flex flex-col space-y-1 mb-3">
                      <button
                        onClick={() => applyFilter('全て')}
                        className={`flex items-center px-2 py-1 text-sm rounded-md transition-colors ${
                          statusFilter === '全て'
                            ? 'bg-primary/10 text-primary'
                            : 'hover:bg-gray-100'
                        }`}
                      >
                        全て
                        {statusFilter === '全て' && <CheckCircle className="ml-auto w-4 h-4" />}
                      </button>
                      <button
                        onClick={() => applyFilter('完成')}
                        className={`flex items-center px-2 py-1 text-sm rounded-md transition-colors ${
                          statusFilter === '完成'
                            ? 'bg-primary/10 text-primary'
                            : 'hover:bg-gray-100'
                        }`}
                      >
                        完成
                        {statusFilter === '完成' && <CheckCircle className="ml-auto w-4 h-4" />}
                      </button>
                      <button
                        onClick={() => applyFilter('処理中')}
                        className={`flex items-center px-2 py-1 text-sm rounded-md transition-colors ${
                          statusFilter === '処理中'
                            ? 'bg-primary/10 text-primary'
                            : 'hover:bg-gray-100'
                        }`}
                      >
                        処理中
                        {statusFilter === '処理中' && <CheckCircle className="ml-auto w-4 h-4" />}
                      </button>
                    </div>

                    <div className="text-sm font-medium text-gray-600 pb-1 mb-1 border-b">
                      並び順
                    </div>
                    <div className="flex flex-col space-y-1">
                      <button
                        onClick={() => applySort('新しい順')}
                        className={`flex items-center px-2 py-1 text-sm rounded-md transition-colors ${
                          sortOrder === '新しい順'
                            ? 'bg-primary/10 text-primary'
                            : 'hover:bg-gray-100'
                        }`}
                      >
                        新しい順
                        {sortOrder === '新しい順' && <CheckCircle className="ml-auto w-4 h-4" />}
                      </button>
                      <button
                        onClick={() => applySort('古い順')}
                        className={`flex items-center px-2 py-1 text-sm rounded-md transition-colors ${
                          sortOrder === '古い順'
                            ? 'bg-primary/10 text-primary'
                            : 'hover:bg-gray-100'
                        }`}
                      >
                        古い順
                        {sortOrder === '古い順' && <CheckCircle className="ml-auto w-4 h-4" />}
                      </button>
                      <button
                        onClick={() => applySort('人気順')}
                        className={`flex items-center px-2 py-1 text-sm rounded-md transition-colors ${
                          sortOrder === '人気順'
                            ? 'bg-primary/10 text-primary'
                            : 'hover:bg-gray-100'
                        }`}
                      >
                        人気順
                        {sortOrder === '人気順' && <CheckCircle className="ml-auto w-4 h-4" />}
                      </button>
                    </div>
                  </div>
                </>
              )}
            </div>
          </div>
        </CardContent>
      </Card>
      {/* フィルター条件の表示 */}
      {(statusFilter !== '全て' || sortOrder !== '新しい順') && (
        <div className="mb-4 flex flex-wrap gap-2 items-center">
          <span className="text-sm text-muted-foreground">フィルター:</span>
          {statusFilter !== '全て' && (
            <Badge variant="outline" className="flex items-center gap-1 py-1">
              ステータス: {statusFilter}
            </Badge>
          )}
          {sortOrder !== '新しい順' && (
            <Badge variant="outline" className="flex items-center gap-1 py-1">
              並び順: {sortOrder}
            </Badge>
          )}
          <Button
            variant="ghost"
            size="sm"
            className="h-7 px-2 text-xs"
            onClick={() => {
              setStatusFilter('全て')
              setSortOrder('新しい順')
              setDisplayCount(6)
            }}
          >
            リセット
          </Button>
        </div>
      )}

      {/* ローディング表示 */}
      {isLoading ? (
        <div className="flex justify-center py-20">
          <div className="animate-spin rounded-full h-10 w-10 border-b-2 border-primary"></div>
        </div>
      ) : (
        <>
          {/* 検索結果がない場合のメッセージ */}
          {filteredVideos.length === 0 && (
            <div className="text-center py-10">
              <p className="text-muted-foreground">検索結果が見つかりませんでした。</p>
            </div>
          )}

          {/* Videos Grid */}
          <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-2 relative z-[1]">
            {displayedVideos.map(video => (
              <div
                key={video.id}
                className="relative overflow-hidden group rounded-lg shadow-sm hover:shadow-md transition-shadow"
              >
                <div className="relative aspect-[4/3] bg-gradient-to-br from-primary/10 to-secondary/10 overflow-hidden">
                  <img
                    src={video.thumbnail || '/placeholder.webp'}
                    alt={video.title}
                    className="w-full h-full object-cover"
                  />

                  {/* Gradient overlay for better text visibility */}
                  <div className="absolute inset-0 bg-gradient-to-t from-black/70 to-transparent opacity-70"></div>

                  {/* Content overlay (positioned at the bottom only) */}
                  <div className="absolute bottom-0 left-0 right-0 p-2">
                    {/* Top-right badge for duration */}
                    <Badge className="absolute top-2 right-2 text-xs bg-black/60 text-white">
                      <Clock className="w-2.5 h-2.5 mr-0.5" />
                      {video.duration}
                    </Badge>

                    {/* Status badge (only shown for non-completed) */}
                    {video.status !== '完成' && (
                      <Badge className="absolute top-2 left-2 text-xs bg-yellow-500/80 text-black">
                        {video.status}
                      </Badge>
                    )}

                    {/* Bottom content */}
                    <div className="text-white">
                      <h3 className="font-medium text-base line-clamp-1">{video.title}</h3>
                      <div className="flex items-center justify-between text-xs text-white/80">
                        <div className="flex items-center">
                          <Calendar className="w-3 h-3 mr-1" />
                          {video.date}
                        </div>
                        {video.status === '完成' && (
                          <div className="text-xs text-white/80">♥ {video.likes}</div>
                        )}
                      </div>
                    </div>
                  </div>

                  {/* Play button overlay */}
                  <div className="absolute inset-0 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
                    <Button
                      size="sm"
                      variant="secondary"
                      className="bg-white/90 text-black hover:bg-white rounded-full w-10 h-10 p-0"
                    >
                      <Play className="w-5 h-5" fill="currentColor" />
                    </Button>
                  </div>

                  {/* Action buttons (visible on hover) - positioned at the top right */}
                  <div className="absolute top-2 right-2 flex flex-col gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                    <Button
                      size="sm"
                      variant="outline"
                      className="p-1 h-7 w-7 bg-black/40 hover:bg-black/60 text-white border-white/20 rounded-full"
                    >
                      <Download className="w-3 h-3" />
                    </Button>
                    <Button
                      size="sm"
                      variant="outline"
                      className="p-1 h-7 w-7 bg-black/40 hover:bg-black/60 text-white border-white/20 rounded-full"
                    >
                      <Share2 className="w-3 h-3" />
                    </Button>
                    <Button
                      size="sm"
                      variant="outline"
                      className="p-1 h-7 w-7 bg-black/40 hover:bg-black/60 text-white border-white/20 rounded-full hover:text-red-400"
                    >
                      <Trash2 className="w-3 h-3" />
                    </Button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </>
      )}

      {/* Load More - 表示できる動画がまだある場合のみ表示 */}
      {displayedVideos.length < filteredVideos.length && (
        <div className="text-center mt-8">
          <Button variant="outline" size="lg" onClick={handleLoadMore}>
            さらに読み込む ({filteredVideos.length - displayedVideos.length}件)
          </Button>
        </div>
      )}
    </>
  )
}
