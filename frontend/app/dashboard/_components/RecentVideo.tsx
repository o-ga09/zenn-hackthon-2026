'use client'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Play, Calendar, Download, Share2, Video, Image, Upload, X, Loader2 } from 'lucide-react'
import React, { useState } from 'react'
import { useAnalyzeMedia, useGetMediaList, MEDIA_QUERY_KEYS } from '@/api/mediaApi'
import { useGetVlogs } from '@/api/vlogAPi'
import { toast } from 'sonner'

export default function RecentVideo() {
  const [isUploadDialogOpen, setIsUploadDialogOpen] = useState(false)
  const [selectedFiles, setSelectedFiles] = useState<FileList | null>(null)
  const [uploadProgress, setUploadProgress] = useState<{ [key: string]: boolean }>({})

  const analyzeMutation = useAnalyzeMedia()
  const { data: mediaListData, isLoading: isMediaLoading } = useGetMediaList()
  const { data: vlogsData, isLoading: isVlogsLoading } = useGetVlogs()

  // 実際のメディアデータから素材動画を生成（最新3件のみ）
  const originalMedia = (mediaListData?.media || [])
    .sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
    .slice(0, 3)
    .map(media => ({
      id: media.id,
      title: `${media.type === 'video' ? '動画' : '画像'}_${media.id.slice(-6)}`, // タイプに応じたタイトル
      date: new Date(media.created_at).toLocaleDateString('ja-JP'),
      thumbnail: media.url,
      image_data: media.image_data,
      video_url: media.type === 'video' ? media.url : undefined, // 動画の場合はvideo_urlを設定
      duration: media.type === 'video' ? '不明' : '', // 動画の場合は不明、画像の場合は固定値
      type: 'original',
      media_type: media.type, // メディアタイプを保持
    }))

  const handleFileSelect = (event: React.ChangeEvent<HTMLInputElement>) => {
    setSelectedFiles(event.target.files)
  }

  const handleUpload = async () => {
    if (!selectedFiles || selectedFiles.length === 0) return

    try {
      const files = Array.from(selectedFiles)

      // 各ファイルの分析状態を初期化
      const initialProgress = files.reduce((acc, file) => {
        acc[file.name] = true
        return acc
      }, {} as { [key: string]: boolean })
      setUploadProgress(initialProgress)

      // すべてのファイルを一括でアップロードして分析
      toast.info(`${files.length}個のファイルのアップロードと分析を開始します...`)
      
      await analyzeMutation.mutateAsync(files)

      toast.success(`${files.length}個のメディアファイルのアップロードと分析が完了しました。`)

      // ダイアログを閉じて状態をリセット
      setIsUploadDialogOpen(false)
      setSelectedFiles(null)
      setUploadProgress({})
    } catch (error) {
      console.error('アップロード・分析エラー:', error)
      toast.error('アップロードまたは分析中にエラーが発生しました。')
      setUploadProgress({})
    }
  }

  const isUploading = Object.keys(uploadProgress).length > 0

  // 生成動画（VLog）を VideoGrid 用の形式に変換
  const generatedVideos = (vlogsData?.items || [])
    .sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
    .map(v => ({
      id: (v as any).id || String(Math.random()),
      title: (v as any).share_url || `VLog_${String((v as any).id || '').slice(-6)}`,
      date: new Date((v as any).created_at).toLocaleDateString('ja-JP'),
      thumbnail: (v as any).thumbnail || (v as any).video_url || '/placeholder.webp',
      video_url: (v as any).video_url,
      duration:
        typeof (v as any).duration === 'number'
          ? `${Math.floor((v as any).duration / 60)}:${String(
              Math.floor((v as any).duration % 60)
            ).padStart(2, '0')}`
          : '0:00',
      type: 'generated',
    }))

  return (
    <div className="mb-6">
      <div className="flex items-center justify-between mb-3 md:mb-4">
        <h2 className="text-lg md:text-xl font-semibold text-foreground">最近の動画</h2>
        <div className="flex items-center gap-2">
          <Dialog open={isUploadDialogOpen} onOpenChange={setIsUploadDialogOpen}>
            <DialogTrigger asChild>
              <Button variant="default" size="sm" className="text-xs md:text-sm h-8 px-3 gap-1.5">
                <Upload className="w-3 h-3 md:w-4 md:h-4" />
                素材アップロード
              </Button>
            </DialogTrigger>
            <DialogContent className="sm:max-w-[425px]">
              <DialogHeader>
                <DialogTitle className="flex items-center gap-2">
                  <Upload className="w-4 h-4" />
                  素材をアップロード
                </DialogTitle>
                <DialogDescription>
                  旅行の写真や動画をアップロードすると、AIが内容を分析してタグ付けや要約を自動で行います。
                </DialogDescription>
              </DialogHeader>
              <div className="grid gap-4 py-4">
                <div className="grid gap-2">
                  <Label htmlFor="file-upload">ファイルを選択</Label>
                  <Input
                    id="file-upload"
                    type="file"
                    multiple
                    accept="image/*,video/*"
                    onChange={handleFileSelect}
                    className="cursor-pointer"
                    disabled={isUploading}
                  />
                  {selectedFiles && (
                    <div className="text-sm text-muted-foreground">
                      {selectedFiles.length}個のファイルが選択されています
                    </div>
                  )}
                  {isUploading && (
                    <div className="space-y-2">
                      <div className="text-sm font-medium">アップロードと分析を行っています...</div>
                      <div className="space-y-1">
                        {Object.keys(uploadProgress).map(fileName => (
                          <div
                            key={fileName}
                            className="flex items-center gap-2 text-xs text-muted-foreground"
                          >
                            <Loader2 className="w-3 h-3 animate-spin" />
                            {fileName}
                          </div>
                        ))}
                      </div>
                    </div>
                  )}
                </div>
                <div className="bg-muted/50 rounded-lg p-4 text-sm">
                  <div className="font-medium mb-2">対応ファイル形式:</div>
                  <div className="text-muted-foreground">
                    • 画像: JPG, PNG, WEBP, GIF
                    <br />• 動画: MP4, MOV, AVI, WebM
                  </div>
                </div>
              </div>
              <div className="flex justify-end gap-2">
                <Button
                  variant="outline"
                  onClick={() => setIsUploadDialogOpen(false)}
                  disabled={isUploading}
                >
                  キャンセル
                </Button>
                <Button
                  onClick={handleUpload}
                  disabled={!selectedFiles || selectedFiles.length === 0 || isUploading}
                >
                  {isUploading ? (
                    <>
                      <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                      アップロード中
                    </>
                  ) : (
                    'アップロード開始'
                  )}
                </Button>
              </div>
            </DialogContent>
          </Dialog>
          <Button variant="outline" size="sm" className="text-xs md:text-sm h-8 px-3">
            すべて表示
          </Button>
        </div>
      </div>

      <Tabs defaultValue="generated" className="w-full">
        <TabsList className="grid w-full grid-cols-2 mb-4 bg-muted/50 border border-border">
          <TabsTrigger
            value="generated"
            className="flex items-center gap-2 text-xs md:text-sm data-[state=active]:bg-primary data-[state=active]:text-primary-foreground data-[state=active]:shadow-sm font-medium transition-all"
          >
            <Video className="w-3 h-3 md:w-4 md:h-4" />
            生成した動画
          </TabsTrigger>
          <TabsTrigger
            value="original"
            className="flex items-center gap-2 text-xs md:text-sm data-[state=active]:bg-primary data-[state=active]:text-primary-foreground data-[state=active]:shadow-sm font-medium transition-all"
          >
            <Image className="w-3 h-3 md:w-4 md:h-4" />
            素材メディア
          </TabsTrigger>
        </TabsList>

        <TabsContent value="generated" className="mt-0">
          <VideoGrid
            videos={vlogsData && generatedVideos.length > 0 ? generatedVideos : []}
            isLoading={isVlogsLoading}
          />
        </TabsContent>

        <TabsContent value="original" className="mt-0">
          <VideoGrid videos={originalMedia} isLoading={isMediaLoading} />
        </TabsContent>
      </Tabs>
    </div>
  )
}

// 動画グリッドコンポーネント
interface VideoGridProps {
  videos: Array<{
    id: string
    title: string
    date: string
    thumbnail: string
    image_data?: string
    video_url?: string
    duration: string
    type: string
    media_type?: string // メディアタイプ (image or video)
  }>
  isLoading?: boolean
}

function VideoGrid({ videos, isLoading = false }: VideoGridProps) {
  if (isLoading) {
    return (
      <div className="flex justify-center py-8">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
      </div>
    )
  }

  return (
    <div>
      {videos.length === 0 ? (
        <div className="text-center py-8 text-muted-foreground">
          動画がありません。新しい動画を作成してみましょう！
        </div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-3 md:gap-4">
          {videos.map(video => (
            <div
              key={video.id}
              className="relative aspect-[4/3] bg-gradient-to-br from-primary/10 to-secondary/10 overflow-hidden rounded-lg"
            >
              {video.video_url || video.media_type === 'video' ? (
                <video
                  src={video.video_url || video.image_data}
                  poster={video.image_data}
                  className="w-full h-full object-cover"
                  controls
                  playsInline
                  preload="metadata"
                />
              ) : (
                <img
                  src={video.image_data || video.thumbnail}
                  alt={video.title}
                  className="w-full h-full object-cover"
                  loading="lazy"
                />
              )}

              {/* タイトルオーバーレイ */}
              <div className="absolute top-0 left-0 right-0 bg-gradient-to-b from-black/70 to-transparent p-2 pointer-events-none">
                <h3 className="text-white font-medium text-sm line-clamp-1">{video.title}</h3>
                <div className="flex items-center text-xs text-white/80">
                  <Calendar className="w-3 h-3 mr-1" />
                  {video.date}
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
