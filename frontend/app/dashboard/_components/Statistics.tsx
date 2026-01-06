'use client'
import { useGetUserPhotoCount } from '@/api/userApi'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { useAuth } from '@/context/authContext'
import React from 'react'

export default function Statistics() {
  const { user } = useAuth()
  const { data } = useGetUserPhotoCount(user?.userID || '')
  return (
    <>
      <Card className="border-0 shadow-lg bg-card/50 backdrop-blur-sm">
        <CardHeader className="px-4 md:px-6 py-3 md:py-5">
          <CardTitle className="text-base md:text-lg lg:text-xl">統計情報</CardTitle>
        </CardHeader>
        <CardContent className="px-4 md:px-6 py-2 md:py-4">
          <div className="grid grid-cols-2 gap-3 md:gap-6">
            <div className="text-center p-2 bg-primary/5 rounded-lg">
              <div className="text-2xl md:text-3xl font-bold text-primary mb-1">
                {data?.videoCount ?? 0}
              </div>
              <div className="text-2xs md:text-sm text-muted-foreground">作成した動画</div>
            </div>
            <div className="text-center p-2 bg-secondary/5 rounded-lg">
              <div className="text-2xl md:text-3xl font-bold text-primary mb-1">
                {data?.uploadCount ?? 0}
              </div>
              <div className="text-2xs md:text-sm text-muted-foreground">アップロード写真</div>
            </div>
          </div>
        </CardContent>
      </Card>
    </>
  )
}
