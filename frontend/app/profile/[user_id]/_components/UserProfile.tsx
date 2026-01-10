import React from 'react'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'

type UserProfileProps = {
  userId: string
  name: string
  avatarUrl: string
  bio: string
  totalMemories: number
  followers: number
  following: number
  occupation?: string
}

export default function UserProfile({
  name,
  avatarUrl,
  bio,
  totalMemories,
  followers,
  following,
  occupation,
}: UserProfileProps) {
  return (
    <div className="max-w-screen-lg mx-auto">
      <div className="flex flex-col md:flex-row items-start gap-8 p-4">
        {/* プロフィール画像 */}
        <div className="w-24 h-24 rounded-full overflow-hidden flex-shrink-0">
          <img
            alt={`${name}のプロフィール画像`}
            src={avatarUrl}
            className="w-full h-full object-cover"
          />
        </div>

        {/* プロフィール情報 */}
        <div className="flex-grow">
          <div className="flex items-center gap-4 mb-4">
            <h1 className="text-2xl font-bold">{name}</h1>
            <Button variant="outline" size="sm">
              フォローする
            </Button>
          </div>

          {/* ステータス情報 */}
          <div className="flex items-center gap-6 mb-4 text-sm text-muted-foreground">
            <div className="flex items-center gap-1">
              <span className="font-semibold text-foreground">{totalMemories}</span>
              <span>思い出</span>
            </div>
            <div className="flex items-center gap-1">
              <span className="font-semibold text-foreground">{followers}</span>
              <span>フォロワー</span>
            </div>
            <div className="flex items-center gap-1">
              <span className="font-semibold text-foreground">{following}</span>
              <span>フォロー中</span>
            </div>
          </div>

          {/* プロフィール詳細 */}
          <div className="space-y-2">
            {occupation && <p className="text-sm text-muted-foreground">{occupation}</p>}
            <p className="text-sm whitespace-pre-wrap">{bio}</p>
          </div>
        </div>
      </div>
    </div>
  )
}
