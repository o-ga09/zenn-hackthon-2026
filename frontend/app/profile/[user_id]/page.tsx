'use client'
import React from 'react'
import MainLayout from '@/components/layout/MainLayout'
import UserProfile from './_components/UserProfile'
import TravelMemoryCard from './_components/TravelMemoryCard'
import { useParams } from 'next/navigation'
import { useGetUserByName, useGetUserPhotoCount } from '@/api/user'
import { useGetTravelsByUserId } from '@/api/travelApi'
import { useAuth } from '@/context/authContext'
import { TravelMemory } from '@/api/extendedTypes'

export default function UserProfilePage() {
  const { user_id } = useParams()
  const targetUserId = user_id as string
  const { user: currentUser } = useAuth()

  // 対象ユーザーの情報を取得
  const { data: profileUser, isLoading: isLoadingProfile } = useGetUserByName(targetUserId)
  // アップロード数の取得
  const { data: photoCount } = useGetUserPhotoCount(targetUserId)

  // 旅行情報（メモリ）の取得
  const { data: travelsData, isLoading: isLoadingMemories } = useGetTravelsByUserId(targetUserId)

  // 取得した旅行データを TravelMemory 形式に変換
  const memories: TravelMemory[] =
    travelsData?.travels.map(travel => ({
      id: travel.id,
      title: travel.title,
      location: '旅の思い出', // 既存データに場所情報がないためのプレースホルダー
      date: new Date(travel.startDate).toLocaleDateString('ja-JP'),
      thumbnailUrl: travel.thumbnail || '/placeholder.webp',
      likes: Math.floor(Math.random() * 10), // 本来はAPIから取得するが、現状はダミー
      description: travel.description,
    })) || []

  if (isLoadingProfile) {
    return (
      <MainLayout>
        <div className="max-w-screen-lg mx-auto p-4 flex justify-center py-20">
          <div className="animate-spin rounded-full h-10 w-10 border-b-2 border-primary"></div>
        </div>
      </MainLayout>
    )
  }

  if (!profileUser) {
    return (
      <MainLayout>
        <div className="max-w-screen-lg mx-auto p-4">
          <p>指定されたユーザーは存在しません。</p>
        </div>
      </MainLayout>
    )
  }

  // 公開設定のチェック (自分自身であれば非公開でも表示)
  const isOwner = currentUser?.id === profileUser.id
  if (!profileUser.isPublic && !isOwner) {
    return (
      <MainLayout>
        <div className="max-w-screen-lg mx-auto p-4">
          <p>このプロフィールは非公開です。</p>
        </div>
      </MainLayout>
    )
  }

  return (
    <MainLayout>
      {/* プロフィールセクション */}
      <div className="border-b">
        <UserProfile
          userId={profileUser.id}
          name={profileUser.displayName ?? ''}
          occupation={profileUser.name || 'ユーザー'}
          avatarUrl={profileUser.profileImage || '/placeholder.webp'}
          bio={profileUser.bio || ''}
          totalMemories={photoCount?.videoCount || 0}
          followers={profileUser.followersCount || 0}
          following={profileUser.followingCount || 0}
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
                <div className="w-full aspect-[1.91/1] bg-muted rounded-lg animate-pulse" />
                <div className="space-y-2">
                  <div className="h-4 w-3/4 bg-muted rounded animate-pulse" />
                  <div className="h-4 w-1/2 bg-muted rounded animate-pulse" />
                </div>
              </div>
            ))}
          </div>
        ) : memories.length > 0 ? (
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            {memories.map(memory => (
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
