'use client'

import React, {
  createContext,
  useContext,
  useState,
  useCallback,
  ReactNode,
  useEffect,
} from 'react'
import { Media } from '@/api/types'
import { notificationApi, Notification as ApiNotification } from '@/api/notificationApi'
import { useAuth } from './authContext'

export interface Notification {
  id: string
  type: 'success' | 'error' | 'info'
  title: string
  message: string
  mediaId?: string
  media?: Media
  timestamp: Date
  read: boolean
}

interface NotificationContextType {
  notifications: Notification[]
  unreadCount: number
  addNotification: (notification: Omit<Notification, 'id' | 'timestamp' | 'read'>) => void
  markAsRead: (id: string) => void
  markAllAsRead: () => void
  removeNotification: (id: string) => void
  clearAll: () => void
  fetchNotifications: () => Promise<void>
}

const NotificationContext = createContext<NotificationContextType | undefined>(undefined)

interface NotificationProviderProps {
  children: ReactNode
}

// APIの通知タイプをUIの通知タイプに変換
function mapApiNotificationToUI(apiNotif: ApiNotification): Notification {
  let uiType: 'success' | 'error' | 'info' = 'info'

  if (apiNotif.type === 'media_completed' || apiNotif.type === 'vlog_completed') {
    uiType = 'success'
  } else if (apiNotif.type === 'media_failed' || apiNotif.type === 'vlog_failed') {
    uiType = 'error'
  }

  return {
    id: apiNotif.id,
    type: uiType,
    title: apiNotif.title,
    message: apiNotif.message,
    mediaId: apiNotif.media_id,
    timestamp: new Date(apiNotif.created_at),
    read: apiNotif.read,
  }
}

export function NotificationProvider({ children }: NotificationProviderProps) {
  const [notifications, setNotifications] = useState<Notification[]>([])
  const { user } = useAuth()

  // ログイン時に通知を読み込む
  useEffect(() => {
    if (user) {
      fetchNotifications()
    }
  }, [user])

  const fetchNotifications = useCallback(async () => {
    try {
      const data = await notificationApi.getNotifications()
      const uiNotifications = data.notifications.map(mapApiNotificationToUI)
      setNotifications(uiNotifications)
    } catch (error) {
      console.error('Failed to fetch notifications:', error)
    }
  }, [])

  const addNotification = useCallback(
    (notification: Omit<Notification, 'id' | 'timestamp' | 'read'>) => {
      const newNotification: Notification = {
        ...notification,
        id: `notification-${Date.now()}-${Math.random()}`,
        timestamp: new Date(),
        read: false,
      }
      setNotifications(prev => [newNotification, ...prev])
    },
    []
  )

  const markAsRead = useCallback(async (id: string) => {
    try {
      await notificationApi.markAsRead(id)
      setNotifications(prev =>
        prev.map(notification =>
          notification.id === id ? { ...notification, read: true } : notification
        )
      )
    } catch (error) {
      console.error('Failed to mark notification as read:', error)
    }
  }, [])

  const markAllAsRead = useCallback(async () => {
    try {
      await notificationApi.markAllAsRead()
      setNotifications(prev => prev.map(notification => ({ ...notification, read: true })))
    } catch (error) {
      console.error('Failed to mark all notifications as read:', error)
    }
  }, [])

  const removeNotification = useCallback(async (id: string) => {
    try {
      await notificationApi.deleteNotification(id)
      setNotifications(prev => prev.filter(notification => notification.id !== id))
    } catch (error) {
      console.error('Failed to delete notification:', error)
    }
  }, [])

  const clearAll = useCallback(() => {
    setNotifications([])
  }, [])

  const unreadCount = notifications.filter(n => !n.read).length

  return (
    <NotificationContext.Provider
      value={{
        notifications,
        unreadCount,
        addNotification,
        markAsRead,
        markAllAsRead,
        removeNotification,
        clearAll,
        fetchNotifications,
      }}
    >
      {children}
    </NotificationContext.Provider>
  )
}

export function useNotifications() {
  const context = useContext(NotificationContext)
  if (!context) {
    throw new Error('useNotifications must be used within a NotificationProvider')
  }
  return context
}
