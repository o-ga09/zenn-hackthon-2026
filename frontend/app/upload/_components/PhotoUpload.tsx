'use client'

import { Button } from '@/components/ui/button'
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Camera, Upload } from 'lucide-react'
import React, { useState, useRef, Suspense } from 'react'
import { X, Library as LibraryIcon, CheckCircle2 } from 'lucide-react'
import { useUploadForm } from './UploadFormContext'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { useSearchParams } from 'next/navigation'
import { useGetMediaList } from '@/api/mediaApi'
import { Media } from '@/api/types'
import { MediaAnalyticsDialog } from './MediaAnalyticsDialog'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'

function PhotoUploadInner() {
  const searchParams = useSearchParams()
  const initialSource = searchParams.get('source') || 'upload'

  const { uploadedFiles, addFiles, removeFile, selectedMediaIds, toggleMediaId } = useUploadForm()
  const [isDragging, setIsDragging] = useState(false)
  const [selectedMedia, setSelectedMedia] = useState<Media | null>(null)
  const [isDialogOpen, setIsDialogOpen] = useState(false)
  const fileInputRef = useRef<HTMLInputElement>(null)

  const { data: mediaListData, isLoading: isMediaLoading } = useGetMediaList()

  const handleMediaClick = (media: Media, e: React.MouseEvent) => {
    // 分析完了のメディアのみクリック可能
    if (media.status === 'completed') {
      e.stopPropagation()
      setSelectedMedia(media)
      setIsDialogOpen(true)
    }
  }

  const handleFileUpload = (event: React.ChangeEvent<HTMLInputElement>) => {
    const files = Array.from(event.target.files || [])
    if (files.length > 0) {
      addFiles(files)
    }
  }

  const handleDragEnter = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault()
    e.stopPropagation()
    setIsDragging(true)
  }

  const handleDragLeave = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault()
    e.stopPropagation()
    setIsDragging(false)
  }

  const handleDragOver = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault()
    e.stopPropagation()
    setIsDragging(true)
  }

  const handleDrop = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault()
    e.stopPropagation()
    setIsDragging(false)

    const files = Array.from(e.dataTransfer.files)
    if (files.length > 0) {
      const imageFiles = files.filter(file => file.type.startsWith('image/'))
      if (imageFiles.length > 0) {
        addFiles(imageFiles)
      }
    }
  }

  const handleButtonClick = () => {
    fileInputRef.current?.click()
  }

  return (
    <div className="space-y-6">
      <Tabs defaultValue={initialSource} className="w-full">
        <TabsList className="grid w-full grid-cols-2 mb-6">
          <TabsTrigger value="upload" className="flex items-center gap-2">
            <Upload className="w-4 h-4" />
            新規アップロード
          </TabsTrigger>
          <TabsTrigger value="library" className="flex items-center gap-2">
            <LibraryIcon className="w-4 h-4" />
            ライブラリから選択
          </TabsTrigger>
        </TabsList>

        <TabsContent value="upload">
          <Card className="border-0 shadow-lg bg-card/50 backdrop-blur-sm">
            <CardHeader>
              <CardTitle className="flex items-center space-x-2 text-base md:text-lg lg:text-xl">
                <Camera className="w-4 h-4 md:w-5 md:h-5" />
                <span>写真をアップロード</span>
              </CardTitle>
              <CardDescription className="text-xs md:text-sm">
                旅行の思い出の写真を選択してください（複数選択可能）
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div
                  className={`border-2 border-dashed rounded-lg p-4 md:p-8 text-center transition-colors ${
                    isDragging
                      ? 'border-primary bg-primary/5'
                      : 'border-border hover:border-primary/50'
                  }`}
                  onDragEnter={handleDragEnter}
                  onDragOver={handleDragOver}
                  onDragLeave={handleDragLeave}
                  onDrop={handleDrop}
                  onClick={handleButtonClick}
                >
                  <Upload
                    className={`w-10 h-10 md:w-12 md:h-12 mx-auto mb-3 md:mb-4 ${
                      isDragging ? 'text-primary' : 'text-muted-foreground'
                    }`}
                  />
                  <p className="text-sm md:text-base text-muted-foreground mb-3 md:mb-4">
                    ここに写真をドラッグ&ドロップするか、クリックして選択
                  </p>
                  <Input
                    ref={fileInputRef}
                    type="file"
                    multiple
                    accept="image/*"
                    onChange={handleFileUpload}
                    className="hidden"
                    id="file-upload"
                  />
                  <Button
                    variant="outline"
                    className="cursor-pointer bg-transparent text-sm md:text-base"
                  >
                    写真を選択
                  </Button>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="library">
          <Card className="border-0 shadow-lg bg-card/50 backdrop-blur-sm">
            <CardHeader>
              <CardTitle className="flex items-center space-x-2 text-base md:text-lg lg:text-xl">
                <LibraryIcon className="w-4 h-4 md:w-5 md:h-5" />
                <span>ライブラリから選択</span>
              </CardTitle>
              <CardDescription className="text-xs md:text-sm">
                アップロード済みの写真から動画に使用するものを選んでください
              </CardDescription>
            </CardHeader>
            <CardContent>
              {isMediaLoading ? (
                <div className="flex justify-center py-12">
                  <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
                </div>
              ) : !mediaListData || !mediaListData.media || mediaListData.media.length === 0 ? (
                <div className="text-center py-12 text-muted-foreground">
                  まだメディアがアップロードされていません。
                </div>
              ) : (
                <TooltipProvider>
                  <div className="grid grid-cols-3 sm:grid-cols-4 md:grid-cols-5 gap-2 max-h-[400px] overflow-y-auto p-1">
                    {mediaListData.media.map(media => {
                      const isAnalyzing = media.status === 'pending' || media.status === 'uploading'
                      const isFailed = media.status === 'failed'
                      const isCompleted = media.status === 'completed'

                      return (
                        <Tooltip key={media.id}>
                          <TooltipTrigger asChild>
                            <div
                              className={`relative group transition-all ${
                                selectedMediaIds.includes(media.id)
                                  ? 'ring-2 ring-primary ring-offset-2'
                                  : ''
                              } ${isCompleted ? 'cursor-pointer' : 'cursor-not-allowed opacity-70'}`}
                              onClick={() => {
                                if (isCompleted) {
                                  toggleMediaId(media.id)
                                }
                              }}
                              onContextMenu={e => {
                                if (isCompleted) {
                                  e.preventDefault()
                                  handleMediaClick(media, e)
                                }
                              }}
                            >
                              <img
                                src={media.url}
                                alt="Media"
                                className="w-full aspect-square object-cover rounded-lg"
                              />
                              {selectedMediaIds.includes(media.id) && (
                                <div className="absolute inset-0 bg-primary/20 flex items-center justify-center rounded-lg">
                                  <CheckCircle2 className="w-6 h-6 text-primary drop-shadow-md fill-white" />
                                </div>
                              )}
                              {isAnalyzing && (
                                <div className="absolute inset-0 bg-black/50 flex items-center justify-center rounded-lg">
                                  <div className="text-white text-xs text-center px-2">
                                    分析中...
                                  </div>
                                </div>
                              )}
                              {isFailed && (
                                <div className="absolute top-1 right-1">
                                  <span className="bg-destructive text-destructive-foreground text-xs px-2 py-1 rounded">
                                    分析失敗
                                  </span>
                                </div>
                              )}
                              {isCompleted && (
                                <Button
                                  size="sm"
                                  variant="secondary"
                                  className="absolute bottom-1 right-1 w-6 h-6 p-0 opacity-0 group-hover:opacity-100 transition-opacity text-xs"
                                  onClick={e => handleMediaClick(media, e)}
                                >
                                  詳細
                                </Button>
                              )}
                            </div>
                          </TooltipTrigger>
                          <TooltipContent>
                            <p>
                              {isAnalyzing && '分析中です'}
                              {isFailed && '分析に失敗しました'}
                              {isCompleted && '右クリックで詳細を表示'}
                            </p>
                          </TooltipContent>
                        </Tooltip>
                      )
                    })}
                  </div>
                </TooltipProvider>
              )}
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>

      {/* Selection Summary */}
      {(uploadedFiles.length > 0 || selectedMediaIds.length > 0) && (
        <div className="bg-primary/5 border border-primary/20 rounded-lg p-4 animate-in fade-in slide-in-from-bottom-2">
          <div className="flex items-center justify-between">
            <div className="flex gap-4">
              {uploadedFiles.length > 0 && (
                <div>
                  <div className="text-xs text-muted-foreground mb-1">新規アップロード</div>
                  <div className="text-lg font-bold text-primary">
                    {uploadedFiles.length}{' '}
                    <span className="text-sm font-normal text-muted-foreground">枚</span>
                  </div>
                </div>
              )}
              {selectedMediaIds.length > 0 && (
                <div>
                  <div className="text-xs text-muted-foreground mb-1">ライブラリから選択</div>
                  <div className="text-lg font-bold text-primary">
                    {selectedMediaIds.length}{' '}
                    <span className="text-sm font-normal text-muted-foreground">枚</span>
                  </div>
                </div>
              )}
            </div>
            <div className="text-right">
              <div className="text-xs text-muted-foreground mb-1">合計</div>
              <div className="text-lg font-bold text-primary">
                {uploadedFiles.length + selectedMediaIds.length}{' '}
                <span className="text-sm font-normal text-muted-foreground">枚</span>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Uploaded Files Preview (新規アップロード分のみ表示) */}
      {uploadedFiles.length > 0 && (
        <div className="space-y-4">
          <h4 className="font-semibold text-sm md:text-base">
            新規アップロード写真 ({uploadedFiles.length}枚)
          </h4>

          <div className="grid grid-cols-3 sm:grid-cols-4 md:grid-cols-5 gap-1 md:gap-2 max-h-36 md:max-h-48 overflow-y-auto p-1 md:p-2 border rounded-lg">
            {uploadedFiles.map((file, index) => (
              <div key={index} className="relative group">
                <img
                  src={URL.createObjectURL(file)}
                  alt={`Upload ${index + 1}`}
                  className="w-full h-16 md:h-20 object-cover rounded-lg"
                />
                <Button
                  size="sm"
                  variant="destructive"
                  className="absolute top-1 right-1 w-5 h-5 md:w-6 md:h-6 p-0 opacity-0 group-hover:opacity-100 transition-opacity"
                  onClick={() => removeFile(index)}
                >
                  <X className="w-2 h-2 md:w-3 md:h-3" />
                </Button>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Analytics Dialog */}
      <MediaAnalyticsDialog
        media={selectedMedia}
        open={isDialogOpen}
        onOpenChange={setIsDialogOpen}
      />
    </div>
  )
}

export default function PhotoUpload() {
  return (
    <Suspense fallback={<div>読み込み中...</div>}>
      <PhotoUploadInner />
    </Suspense>
  )
}
