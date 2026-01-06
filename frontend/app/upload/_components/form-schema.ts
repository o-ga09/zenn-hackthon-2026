'use client'

import { z } from 'zod'

export const travelFormSchema = z.object({
  travelTitle: z.string().min(1, { message: '旅行タイトルを入力してください' }),
  travelDate: z.string().min(1, { message: '旅行日を選択してください' }),
  travelLocation: z.string().optional(),
  travelDescription: z.string().optional(),
  uploadedFiles: z
    .array(z.instanceof(File))
    .min(1, { message: '少なくとも1枚の写真をアップロードしてください' }),
})

export type TravelFormValues = z.infer<typeof travelFormSchema>

export type UploadStep = 'upload' | 'info' | 'confirm'
