import { clsx, type ClassValue } from 'clsx'
import { twMerge } from 'tailwind-merge'
import type { NextTopLoaderProps } from 'nextjs-toploader'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export const topLoaderConfig: NextTopLoaderProps = {
  color: '#FCC31C',
  initialPosition: 0.3,
  crawlSpeed: 300,
  height: 4,
  crawl: true,
  showSpinner: false,
  easing: 'cubic-bezier(0.4, 0, 0.2, 1)',
  speed: 250,
  shadow: '0 2px 15px rgba(233, 100, 43, 0.5)',
  zIndex: 2000,
}

/**
 * UTF-8文字列を安全にBase64エンコードする
 */
export function encodeBase64UTF8(str: string): string {
  try {
    // UTF-8バイト列に変換してからBase64エンコード
    return btoa(
      encodeURIComponent(str).replace(/%([0-9A-F]{2})/g, (match, p1) => {
        return String.fromCharCode(parseInt(p1, 16))
      })
    )
  } catch (error) {
    console.error('Failed to encode Base64 UTF-8:', error)
    throw error
  }
}

/**
 * Base64エンコードされたUTF-8文字列をデコードする
 */
export function decodeBase64UTF8(str: string): string {
  try {
    // Base64デコードしてからUTF-8文字列に変換
    return decodeURIComponent(
      Array.prototype.map
        .call(atob(str), (c: string) => {
          return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2)
        })
        .join('')
    )
  } catch (error) {
    console.error('Failed to decode Base64 UTF-8:', error)
    throw error
  }
}

/**
 * ユーザーデータをクッキーに安全に保存するためのエンコード
 */
export function encodeUserDataForCookie(userData: Record<string, unknown>): string {
  const jsonString = JSON.stringify(userData)
  return encodeBase64UTF8(jsonString)
}

/**
 * クッキーからユーザーデータを安全にデコード
 */
export function decodeUserDataFromCookie(cookieValue: string): Record<string, unknown> {
  const jsonString = decodeBase64UTF8(cookieValue)
  return JSON.parse(jsonString)
}

// CSRFトークン管理クラス
class CSRFManager {
  private static instance: CSRFManager | null = null
  private csrfToken: string | null = null
  private tokenExpiry: number = 0

  private constructor() {
    // プライベートコンストラクタでシングルトンパターンを強制
  }

  static getInstance(): CSRFManager {
    if (!CSRFManager.instance) {
      CSRFManager.instance = new CSRFManager()
    }
    return CSRFManager.instance
  }

  async getCSRFToken(): Promise<string> {
    // キャッシュされたトークンが有効な場合はそれを返す
    if (this.csrfToken && Date.now() < this.tokenExpiry) {
      return this.csrfToken
    }

    try {
      // サーバーからCSRFトークンを取得
      const response = await fetch('/api/auth/csrf', {
        method: 'GET',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
      })

      if (!response.ok) {
        throw new Error(`Failed to fetch CSRF token: ${response.status} ${response.statusText}`)
      }

      const data = (await response.json()) as { token?: string }

      if (!data.token) {
        throw new Error('CSRF token not found in response')
      }

      // トークンをキャッシュ（30分有効）
      this.csrfToken = data.token
      this.tokenExpiry = Date.now() + 30 * 60 * 1000

      return data.token
    } catch (error) {
      console.error('Error fetching CSRF token:', error)
      // トークンの取得に失敗した場合はキャッシュをクリア
      this.clearToken()
      throw error
    }
  }

  // トークンをクリア（ログアウト時など）
  clearToken(): void {
    this.csrfToken = null
    this.tokenExpiry = 0
  }

  // トークンが有効かどうかをチェック
  isTokenValid(): boolean {
    return this.csrfToken !== null && Date.now() < this.tokenExpiry
  }

  // 現在のトークンを取得（キャッシュのみ、APIコールなし）
  getCurrentToken(): string | null {
    return this.isTokenValid() ? this.csrfToken : null
  }
}

export default CSRFManager
