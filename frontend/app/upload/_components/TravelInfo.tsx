'use client'

import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Label } from '@radix-ui/react-label'
import { MapPin } from 'lucide-react'
import React from 'react'
import { useFormContext } from 'react-hook-form'
import { TravelFormValues } from './form-schema'

export default function TravelInfo() {
  const {
    register,
    formState: { errors },
  } = useFormContext<TravelFormValues>()

  return (
    <Card className="border-0 shadow-lg bg-card/50 backdrop-blur-sm">
      <CardHeader className="px-4 md:px-6 py-4 md:py-6">
        <CardTitle className="flex items-center space-x-2 text-base md:text-lg lg:text-xl">
          <MapPin className="w-4 h-4 md:w-5 md:h-5" />
          <span>旅行情報</span>
        </CardTitle>
        <CardDescription className="text-xs md:text-sm">
          動画に含める旅行の詳細情報を入力してください
        </CardDescription>
      </CardHeader>
      <CardContent className="px-4 md:px-6 py-3 md:py-4">
        <div className="space-y-3 md:space-y-4">
          <div>
            <Label htmlFor="travel-title" className="text-sm md:text-base mb-1 block">
              旅行のタイトル
            </Label>
            <Input
              id="travel-title"
              placeholder="例: 京都の桜旅行"
              className="text-sm md:text-base"
              {...register('travelTitle')}
            />
            {errors.travelTitle && (
              <p className="mt-1 text-xs md:text-sm text-red-500">{errors.travelTitle.message}</p>
            )}
          </div>

          <div>
            <Label htmlFor="travel-date" className="text-sm md:text-base mb-1 block">
              旅行日
            </Label>
            <Input
              id="travel-date"
              type="date"
              className="text-sm md:text-base"
              {...register('travelDate')}
            />
            {errors.travelDate && (
              <p className="mt-1 text-xs md:text-sm text-red-500">{errors.travelDate.message}</p>
            )}
          </div>

          <div>
            <Label htmlFor="travel-location" className="text-sm md:text-base mb-1 block">
              場所
            </Label>
            <Input
              id="travel-location"
              placeholder="例: 京都府京都市"
              className="text-sm md:text-base"
              {...register('travelLocation')}
            />
          </div>

          <div>
            <Label htmlFor="travel-description" className="text-sm md:text-base mb-1 block">
              旅行の説明
            </Label>
            <Textarea
              id="travel-description"
              placeholder="この旅行の思い出や特別な瞬間について教えてください..."
              rows={3}
              className="text-sm md:text-base"
              {...register('travelDescription')}
            />
          </div>
        </div>
      </CardContent>
    </Card>
  )
}
