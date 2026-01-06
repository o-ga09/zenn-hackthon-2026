import { Video } from 'lucide-react'
import React from 'react'

export default function LPFooter() {
  return (
    <footer className="container mx-auto px-4 py-4 border-t border-border">
      <div className="text-center">
        <div className="flex items-center justify-center space-x-2 mb-2">
          <div className="w-6 h-6 bg-primary rounded-lg flex items-center justify-center">
            <Video className="w-4 h-4 text-primary-foreground" />
          </div>
          <span className="text-lg font-bold text-foreground">TravelMoments</span>
        </div>
        <p className="text-xs text-muted-foreground mb-2">
          旅の記録を、AIでもっと簡単に、もっとエモく。
        </p>
        <div className="flex flex-row justify-center space-x-4 text-xs text-muted-foreground">
          <a href="/privacy-policy" className="hover:text-foreground transition-colors">
            プライバシーポリシー
          </a>
          <a href="/terms-of-service" className="hover:text-foreground transition-colors">
            利用規約
          </a>
          <a href="/contact" className="hover:text-foreground transition-colors">
            お問い合わせ
          </a>
        </div>
      </div>
    </footer>
  )
}
