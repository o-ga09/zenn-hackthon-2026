'use client'
import React, { use } from 'react'
import MainLayout from '@/components/layout/MainLayout'
import UserProfile from './_components/UserProfile'
import TravelMemoryCard from './_components/TravelMemoryCard'
import { useGetUserById, useGetUserPhotoCount } from '@/api/user'
import { useGetUserMemories } from '@/hooks/useUserData'
import { Skeleton } from '@/components/ui/skeleton'
import { useAuth } from '@/context/authContext'

export default function UserProfilePage() {
  const { user } = useAuth()
  // アップロード数の取得
  const { data: photoCount } = useGetUserPhotoCount(user?.id || '')
  // メモリーの取得
  const { data: memories, isLoading: isLoadingMemories } = useGetUserMemories(user?.id || '')

  if (!user) {
    return (
      <MainLayout>
        <div className="max-w-screen-lg mx-auto p-4">
          <p>指定されたユーザーは存在しません。</p>
        </div>
      </MainLayout>
    )
  }

  return (
    <MainLayout>
      {/* プロフィールセクション */}
      <div className="border-b">
        <UserProfile
          userId={user?.id || ''}
          name={user?.displayName ?? ''}
          occupation={user?.name || 'ユーザー'}
          avatarUrl={user?.profileImage || '/placeholder.webp'}
          bio={user?.bio || ''}
          totalMemories={photoCount?.videoCount || 0}
          followers={user?.followersCount || 0}
          following={user?.followingCount || 0}
        />
      </div>

      {/* タブナビゲーション */}
      <div className="border-b">
        <div className="max-w-screen-lg mx-auto">
          <div className="flex gap-8 px-4">
            <button className="py-4 text-sm text-primary border-b-2 border-primary -mb-px">
              思い出 <span>{photoCount?.videoCount || 0}</span>
            </button>
          </div>
        </div>
      </div>

      {/* 旅行メモリーグリッド */}
      <div className="max-w-screen-lg mx-auto px-4 py-8">
        {isLoadingMemories ? (
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            {[1, 2, 3].map(i => (
              <div key={i} className="space-y-4">
                <div className="w-full aspect-[1.91/1] bg-muted rounded-lg" />
                <div className="space-y-2">
                  <div className="h-4 w-3/4 bg-muted rounded" />
                  <div className="h-4 w-1/2 bg-muted rounded" />
                </div>
              </div>
            ))}
          </div>
        ) : memories?.memories && memories.memories.length > 0 ? (
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            {memories.memories.map(memory => (
              <TravelMemoryCard
                key={memory.id}
                id={memory.id}
                title={memory.title}
                location={memory.location}
                date={memory.date}
                thumbnailUrl={memory.thumbnailUrl}
                likes={memory.likes}
                description={memory.description}
              />
            ))}
          </div>
        ) : (
          <div className="text-center py-12">
            <p className="text-muted-foreground">まだ思い出が投稿されていません。</p>
          </div>
        )}
      </div>
    </MainLayout>
  )
}
