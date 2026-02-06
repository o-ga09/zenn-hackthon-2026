import { apiClient } from './client'

export interface Notification {
  id: string
  version: number
  user_id: string
  type: 'media_completed' | 'media_failed' | 'vlog_completed' | 'vlog_failed'
  title: string
  message: string
  media_id?: string
  vlog_id?: string
  read: boolean
  created_at: string
  updated_at: string
}

export interface NotificationListResponse {
  notifications: Notification[]
  unread_count: number
}

export const notificationApi = {
  getNotifications: async (): Promise<NotificationListResponse> => {
    const response = await apiClient.get('/notifications')
    return response.data
  },

  markAsRead: async (id: string, version: number): Promise<Notification> => {
    const response = await apiClient.put(`/notifications/${id}/read`, {
      version,
    })
    return response.data
  },

  markAllAsRead: async (): Promise<{ updated_count: number; message: string }> => {
    const response = await apiClient.put('/notifications/read-all')
    return response.data
  },

  deleteNotification: async (id: string): Promise<void> => {
    await apiClient.delete(`/notifications/${id}`)
  },
  deleteAllNotifications: async (): Promise<void> => {
    await apiClient.delete(`notifications`)
  },
}
