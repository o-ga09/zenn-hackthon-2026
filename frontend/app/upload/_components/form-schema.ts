'use client'

import { z } from 'zod'

export const travelFormSchema = z.object({
  travelTitle: z.string().optional(),
  travelDate: z.string().optional(),
  travelLocation: z.string().optional(),
  travelDescription: z.string().optional(),
  uploadedFiles: z.array(z.instanceof(File)),
  mediaIds: z.array(z.string()),
}).refine(data => data.uploadedFiles.length > 0 || data.mediaIds.length > 0, {
  message: '少なくとも1つの素材を選択またはアップロードしてください',
  path: ['uploadedFiles'],
})

export type TravelFormValues = z.infer<typeof travelFormSchema>

export type UploadStep = 'upload' | 'info' | 'confirm'
