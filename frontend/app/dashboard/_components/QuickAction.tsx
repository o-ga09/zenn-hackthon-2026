'use client'
import { Button } from '@/components/ui/button'
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card'
import { Plus, Upload } from 'lucide-react'
import React from 'react'
import Statistics from './Statistics'
import { useRouter } from 'next/navigation'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'
import { Image, Library } from 'lucide-react'

export default function QuickAction() {
  const router = useRouter()

  const handleNavigateToUpload = (source: 'upload' | 'library') => {
    router.push(`/upload?source=${source}`)
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-4 md:gap-6 mb-8 md:mb-12">
      <Card className="border-0 shadow-lg bg-card/50 backdrop-blur-sm hover:shadow-xl transition-shadow">
        <CardHeader className="px-4 md:px-6 py-4 md:py-6">
          <div className="flex items-center space-x-3 md:space-x-4">
            <div className="w-10 h-10 md:w-12 md:h-12 bg-primary/10 rounded-full flex items-center justify-center flex-shrink-0">
              <Plus className="w-5 h-5 md:w-6 md:h-6 text-primary" />
            </div>
            <div>
              <CardTitle className="text-base md:text-lg lg:text-xl">新しい動画を作成</CardTitle>
              <CardDescription className="text-xs md:text-sm">
                写真をアップロードして新しい旅行動画を作成
              </CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent className="px-4 md:px-6 py-3 md:py-4">
          <Dialog>
            <DialogTrigger asChild>
              <Button className="w-full bg-primary hover:bg-primary/90 text-sm md:text-base py-1.5 md:py-2 h-auto">
                <Plus className="w-4 h-4 mr-2" />
                新しい動画を作成
              </Button>
            </DialogTrigger>
            <DialogContent className="sm:max-w-md">
              <DialogHeader>
                <DialogTitle>動画の作成方法を選択</DialogTitle>
                <DialogDescription>
                  どのように動画を作成しますか？
                </DialogDescription>
              </DialogHeader>
              <div className="grid grid-cols-1 gap-4 py-4">
                <Button
                  variant="outline"
                  className="h-24 flex flex-col items-center justify-center gap-2 border-2 hover:border-primary hover:bg-primary/5 transition-all"
                  onClick={() => handleNavigateToUpload('library')}
                >
                  <Library className="w-8 h-8 text-primary" />
                  <div className="text-sm font-semibold">ライブラリの素材から作成</div>
                  <div className="text-xs text-muted-foreground">既にアップロード済みの素材を使用します</div>
                </Button>
                <Button
                  variant="outline"
                  className="h-24 flex flex-col items-center justify-center gap-2 border-2 hover:border-primary hover:bg-primary/5 transition-all"
                  onClick={() => handleNavigateToUpload('upload')}
                >
                  <Upload className="w-8 h-8 text-primary" />
                  <div className="text-sm font-semibold">新しくアップロードして作成</div>
                  <div className="text-xs text-muted-foreground">新しく写真や動画をアップロードして作成します</div>
                </Button>
              </div>
            </DialogContent>
          </Dialog>
        </CardContent>
      </Card>
      <Statistics />
    </div>
  )
}
