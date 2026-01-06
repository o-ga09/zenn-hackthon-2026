'use client'

import React from 'react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'

/**
 * APIプロバイダーコンポーネント
 * アプリケーション全体をラップして、TanStack Queryの機能を提供します
 */
export const ApiProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  // クライアントコンポーネント内でQueryClientを初期化
  const [queryClient] = React.useState(
    () =>
      new QueryClient({
        defaultOptions: {
          queries: {
            staleTime: 5 * 60 * 1000, // 5分間キャッシュを保持
            refetchOnWindowFocus: false, // ウィンドウフォーカス時の再取得を無効化
            retry: 1, // エラー時の再試行回数
          },
        },
      })
  )

  return <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
}
