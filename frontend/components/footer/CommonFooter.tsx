import { Video } from 'lucide-react'
import React from 'react'
import Link from 'next/link'

export default function CommonFooter() {
  return (
    <footer className="bg-white/80 backdrop-blur-sm py-6 border-t">
      <div className="container mx-auto px-4">
        <div className="flex flex-col md:flex-row justify-between items-center">
          <div className="mb-4 md:mb-0">
            <div className="flex items-center justify-center space-x-2 mb-4">
              <div className="w-6 h-6 bg-primary rounded-lg flex items-center justify-center">
                <Video className="w-4 h-4 text-primary-foreground" />
              </div>
              <span className="text-lg font-bold text-foreground">TravelMoments</span>
            </div>
            <p className="text-muted-foreground mb-4">
              旅の記録を、AIでもっと簡単に、もっとエモく。
            </p>
            <p className="text-sm text-muted-foreground">© 2025 Tavinikkiy. All rights reserved.</p>
          </div>
          <div className="flex space-x-6">
            <Link href="/terms" className="text-sm text-muted-foreground hover:text-primary">
              利用規約
            </Link>
            <Link href="/privacy" className="text-sm text-muted-foreground hover:text-primary">
              プライバシーポリシー
            </Link>
            <Link href="/contact" className="text-sm text-muted-foreground hover:text-primary">
              お問い合わせ
            </Link>
          </div>
        </div>
      </div>
    </footer>
  )
}
