// API統合用の仮実装（プロダクション環境では実際のデータを取得するAPIを使用）
import { Travel } from '@/api/types'

// 旅行データを強化する関数
export function enhanceTravelData(travel: Travel): Travel {
  // サムネイルがない場合にデフォルト画像を設定
  if (!travel.thumbnail || travel.thumbnail.trim() === '') {
    const defaultImages = [
      '/placeholder.webp',
      'https://placehold.co/600x400/e2e8f0/475569?text=Travel+Video',
    ]

    // サムネイルがなければデフォルト画像を設定
    return {
      ...travel,
      thumbnail: defaultImages[0],
    }
  }

  // 相対パスの場合、絶対URLに変換
  if (
    travel.thumbnail &&
    !travel.thumbnail.startsWith('http') &&
    !travel.thumbnail.startsWith('/')
  ) {
    return {
      ...travel,
      thumbnail: `/${travel.thumbnail}`,
    }
  }

  return travel
}
