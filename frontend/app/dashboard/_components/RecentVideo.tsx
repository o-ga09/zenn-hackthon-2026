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
import { useUploadMediaImage, useGetMediaList, MEDIA_QUERY_KEYS } from '@/api/mediaApi'
import { useGetVlogs } from '@/api/vlogAPi'
import { toast } from 'sonner'

export default function RecentVideo() {
  const [isUploadDialogOpen, setIsUploadDialogOpen] = useState(false)
  const [selectedFiles, setSelectedFiles] = useState<FileList | null>(null)
  const [uploadProgress, setUploadProgress] = useState<{ [key: string]: boolean }>({})

  const uploadMutation = useUploadMediaImage()
  const { data: mediaListData, isLoading: isMediaLoading } = useGetMediaList()
  const { data: vlogsData, isLoading: isVlogsLoading } = useGetVlogs()

  // 実際のメディアデータから素材動画を生成（最新3件のみ）
  const originalMedia = (mediaListData?.media || [])
    .sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
    .slice(0, 3)
    .map(media => ({
      id: media.id,
      title: `画像_${media.id.slice(-6)}`, // ファイル名の代わりに簡易タイトル
      date: new Date(media.created_at).toLocaleDateString('ja-JP'),
      thumbnail: media.url,
      image_data: media.image_data,
      duration: '0:30', // 画像なので固定値
      type: 'original',
    }))

  const handleFileSelect = (event: React.ChangeEvent<HTMLInputElement>) => {
    setSelectedFiles(event.target.files)
  }

  // ファイルをBase64に変換する関数
  const convertFileToBase64 = (file: File): Promise<string> => {
    return new Promise((resolve, reject) => {
      const reader = new FileReader()
      reader.onload = () => {
        const result = reader.result as string
        // "data:image/jpeg;base64," のようなプリフィックスを取り除く
        const base64Data = result.split(',')[1]
        resolve(base64Data)
      }
      reader.onerror = reject
      reader.readAsDataURL(file)
    })
  }

  const handleUpload = async () => {
    if (!selectedFiles || selectedFiles.length === 0) return

    try {
      const files = Array.from(selectedFiles)

      // 各ファイルのアップロード状態を初期化
      const initialProgress = files.reduce((acc, file) => {
        acc[file.name] = true
        return acc
      }, {} as { [key: string]: boolean })
      setUploadProgress(initialProgress)

      // すべてのファイルを並列でアップロード
      const uploadPromises = files.map(async file => {
        try {
          // 画像ファイルのみ処理（動画はAPIが実装されたら対応）
          if (!file.type.startsWith('image/')) return null

          const base64Data = await convertFileToBase64(file)
          const result = await uploadMutation.mutateAsync({
            base64_data: base64Data,
          })

          return result
        } catch (error) {
          toast.error(`${file.name} のアップロードに失敗しました`)
          return null
        } finally {
          // アップロード完了後にプログレス状態を更新
          setUploadProgress(prev => {
            const updated = { ...prev }
            delete updated[file.name]
            return updated
          })
        }
      })

      const results = await Promise.all(uploadPromises)
      const successCount = results.filter(result => result !== null).length
      const failedCount = files.length - successCount

      if (successCount > 0) {
        toast.success(`${successCount}個のファイルがアップロードされました。`)
      }
      if (failedCount > 0) {
        toast.error(`${failedCount}個のファイルのアップロードに失敗しました。`)
      }

      // ダイアログを閉じて状態をリセット
      setIsUploadDialogOpen(false)
      setSelectedFiles(null)
      setUploadProgress({})
    } catch (error) {
      console.error('アップロードエラー:', error)
      toast.error('アップロード中にエラーが発生しました。')
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
                  旅行の写真や動画をアップロードして、AIが素敵な動画を作成します。
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
                      <div className="text-sm font-medium">アップロード中...</div>
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
                    • 画像: JPG, PNG, WEBP (現在対応中)
                    <br />• 動画: MP4, MOV, AVI (今後対応予定)
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
            素材動画
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
  }>
  isLoading?: boolean
}

function VideoGrid({ videos, isLoading = false }: VideoGridProps) {
  const handleVideoClick = (video: any) => {
    console.log('Video clicked:', video)
    // 実際のAPIが実装されたら、ここで動画再生処理を行う
  }

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
              className="relative aspect-[4/3] bg-gradient-to-br from-primary/10 to-secondary/10 overflow-hidden cursor-pointer group rounded-lg"
              onClick={() => handleVideoClick(video)}
            >
              {video.video_url ? (
                <video
                  src={video.video_url}
                  poster={video.image_data}
                  className="w-full h-full object-cover"
                  preload="metadata"
                  controls
                  playsInline
                />
              ) : (
                <img
                  src={video.image_data || video.thumbnail}
                  alt={video.title}
                  className="w-full h-full object-cover"
                  loading="lazy"
                />
              )}
              {/* Gradient overlay for better text visibility */}
              <div className="absolute inset-0 bg-gradient-to-t from-black/70 to-transparent opacity-70"></div>

              {/* Content overlay (positioned at the bottom only) */}
              <div className="absolute bottom-0 left-0 right-0 p-2">
                {/* Duration badge */}
                <Badge className="absolute top-2 right-2 text-2xs md:text-xs bg-black/60 text-white">
                  {video.duration}
                </Badge>

                {/* Bottom content */}
                <div className="text-white">
                  <h3 className="font-medium text-sm md:text-base line-clamp-1 mb-0.5">
                    {video.title}
                  </h3>
                  <div className="flex items-center text-2xs md:text-xs text-white/80">
                    <div className="flex items-center">
                      <Calendar className="w-3 h-3 mr-1" />
                      {video.date}
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
    </div>
  )
}
