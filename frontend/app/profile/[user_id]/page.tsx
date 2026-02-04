'use client'
import React from 'react'
import MainLayout from '@/components/layout/MainLayout'
import UserProfile from './_components/UserProfile'
import TravelMemoryCard from './_components/TravelMemoryCard'
import { useParams } from 'next/navigation'
import { useGetUserByName } from '@/api/user'
import { useAuth } from '@/context/authContext'
import { TravelMemory } from '@/api/extendedTypes'

export default function UserProfilePage() {
  const { user_id } = useParams()
  const targetUserId = user_id as string
  const { user: currentUser } = useAuth()

  // 対象ユーザーの情報を取得
  const { data: profileUser, isLoading: isLoadingProfile } = useGetUserByName(targetUserId)

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
          totalMemories={profileUser?.videoCount || 0}
          followers={profileUser.followersCount || 0}
          following={profileUser.followingCount || 0}
        />
      </div>

      {/* タブナビゲーション */}
      <div className="border-b">
        <div className="max-w-screen-lg mx-auto">
          <div className="flex gap-8 px-4">
            <button className="py-4 text-sm text-primary border-b-2 border-primary -mb-px">
              思い出 <span>{profileUser?.videoCount || 0}</span>
            </button>
          </div>
        </div>
      </div>

      {/* Vlog 一覧 */}
      <div className="text-center py-12">
        <p className="text-muted-foreground">まだ思い出が投稿されていません。</p>
      </div>
    </MainLayout>
  )
}
