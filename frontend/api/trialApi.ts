'use client'

import axios from 'axios'
import { noCredentialApiClient } from './client'

const baseURL = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080'

// VLog生成レスポンスの型定義
export interface TrialVlogResponse {
  jobId?: string
  status?: string
  videoUrl?: string
  thumbnailUrl?: string
  duration?: number
  message?: string
  error?: string
}

export interface TrialVlogRequest {
  files: File[]
  title?: string
  travelDate?: string
  destination?: string
  theme?: string
  musicMood?: string
  duration?: number
  transition?: string
}

/**
 * トライアルVlog生成API（認証不要）
 */
export const createTrialVlog = async (
  request: TrialVlogRequest
): Promise<TrialVlogResponse> => {
  const formData = new FormData()

  // ファイルを追加
  request.files.forEach(file => {
    formData.append('files', file)
  })

  // オプションフィールドを追加
  if (request.title) formData.append('title', request.title)
  if (request.travelDate) formData.append('travelDate', request.travelDate)
  if (request.destination) formData.append('destination', request.destination)
  if (request.theme) formData.append('theme', request.theme)
  if (request.musicMood) formData.append('musicMood', request.musicMood)
  if (request.duration) formData.append('duration', String(request.duration))
  if (request.transition) formData.append('transition', request.transition)

  const response = await noCredentialApiClient.post('/agent/create-vlog', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
    timeout: 120000, // 2分のタイムアウト
  })

  return response.data
}
