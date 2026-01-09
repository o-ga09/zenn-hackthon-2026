import React from 'react'
import { Card } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import Link from 'next/link'

type TravelMemoryProps = {
  id: string
  title: string
  location: string
  date: string
  thumbnailUrl: string
  likes: number
  description?: string
}

export default function TravelMemoryCard({
  id,
  title,
  location,
  date,
  thumbnailUrl,
  likes,
  description,
}: TravelMemoryProps) {
  return (
    <Link href={`/memory/${id}`}>
      <Card className="overflow-hidden hover:shadow-lg transition-shadow duration-200">
        {/* サムネイル */}
        <div className="relative aspect-[1.91/1]">
          <img src={thumbnailUrl} alt={title} className="w-full h-full object-cover" />
          <div className="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/60 to-transparent p-4">
            <Badge variant="secondary" className="bg-white/90">
              {location}
            </Badge>
          </div>
        </div>

        {/* コンテンツ */}
        <div className="p-4">
          <h3 className="text-lg font-semibold line-clamp-2 mb-2 hover:text-blue-500">{title}</h3>
          {description && (
            <p className="text-sm text-muted-foreground line-clamp-2 mb-3">{description}</p>
          )}
          <div className="flex items-center justify-between text-sm text-muted-foreground">
            <time dateTime={date}>{date}</time>
            <div className="flex items-center gap-1">
              <span>❤️</span>
              <span>{likes}</span>
            </div>
          </div>
        </div>
      </Card>
    </Link>
  )
}
