import type { Metadata } from 'next'
import './globals.css'
import NextTopLoader from 'nextjs-toploader'
import { topLoaderConfig } from '@/lib/utils'
import { ApiProvider } from '@/api/apiProvider'
import { AuthProvider } from '@/context/authContext'
import { NotificationProvider } from '@/context/notificationContext'
import { Toaster } from '@/components/ui/sonner'

const APP_URL = process.env.NEXT_PUBLIC_APP_URL || 'http://localhost:3000'

export const metadata: Metadata = {
  title: 'TaviNikkiy - 旅の思い出を記録しよう',
  description:
    'TaviNikkiyは、あなたの旅行の思い出を簡単に記録し、共有できるプラットフォームです。行程、費用、写真を一元管理し、素敵な旅の記録を残しましょう。',
  keywords: '旅行, 記録, 思い出, 写真, 行程管理, 費用管理, 旅程, 旅行計画, 思い出作り',
  authors: [{ name: 'TaviNikkiy Team' }],
  creator: 'TaviNikkiy',
  publisher: 'TaviNikkiy',
  metadataBase: new URL(APP_URL),
  alternates: {
    canonical: '/',
  },
  robots: {
    index: true,
    follow: true,
    googleBot: {
      index: true,
      follow: true,
      'max-image-preview': 'large',
      'max-snippet': -1,
    },
  },
  openGraph: {
    type: 'website',
    locale: 'ja_JP',
    url: APP_URL,
    siteName: 'TaviNikkiy',
    title: 'TaviNikkiy - 旅の思い出を記録しよう',
    description:
      'TaviNikkiyは、あなたの旅行の思い出を簡単に記録し、共有できるプラットフォームです。行程、費用、写真を一元管理し、素敵な旅の記録を残しましょう。',
    images: [
      {
        url: '/og-image.png',
        width: 1200,
        height: 630,
        alt: 'TaviNikkiy - 旅の思い出を記録しよう',
      },
    ],
  },
  twitter: {
    card: 'summary_large_image',
    title: 'TaviNikkiy - 旅の思い出を記録しよう',
    description:
      'TaviNikkiyは、あなたの旅行の思い出を簡単に記録し、共有できるプラットフォームです。行程、費用、写真を一元管理し、素敵な旅の記録を残しましょう。',
    images: ['/og-image.png'],
    creator: '@tavinikkiy',
    site: '@tavinikkiy',
  },
  icons: {
    icon: [
      { url: '/favicon.ico' },
      {
        url: '/android-chrome-192x192.webp',
        sizes: '192x192',
        type: 'image/webp',
      },
      {
        url: '/android-chrome-512x512.webp',
        sizes: '512x512',
        type: 'image/webp',
      },
    ],
    apple: [{ url: '/apple-touch-icon.webp', sizes: '180x180', type: 'image/webp' }],
  },
  manifest: '/manifest.json',
}

export const viewport = 'width=device-width, initial-scale=1'

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode
}>) {
  return (
    <html lang="en">
      <head>
        <link rel="preconnect" href="https://static-agent.tavinikkiy.com" />
        <link rel="dns-prefetch" href="https://static-agent.tavinikkiy.com" />
      </head>
      <body className="overflow-x-hidden" suppressHydrationWarning>
        <ApiProvider>
          <AuthProvider>
            <NotificationProvider>
              <NextTopLoader {...topLoaderConfig} />
              {children}
              <Toaster />
            </NotificationProvider>
          </AuthProvider>
        </ApiProvider>
      </body>
    </html>
  )
}
