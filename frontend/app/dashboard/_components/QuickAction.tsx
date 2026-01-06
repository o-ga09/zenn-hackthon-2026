import { Button } from '@/components/ui/button'
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card'
import { Plus, Upload } from 'lucide-react'
import Link from 'next/link'
import React from 'react'
import Statistics from './Statistics'

export default function QuickAction() {
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
          <Link href="/upload">
            <Button className="w-full bg-primary hover:bg-primary/90 text-sm md:text-base py-1.5 md:py-2 h-auto">
              <Upload className="w-4 h-4 mr-2" />
              写真をアップロード
            </Button>
          </Link>
        </CardContent>
      </Card>
      <Statistics />
    </div>
  )
}
