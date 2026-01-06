import { NextRequest, NextResponse } from 'next/server'
import { decodeUserDataFromCookie } from '@/lib/utils'

// 認証が不要な公開ルートのパスパターンを定義
const PUBLIC_PATHS = [
  '/',
  '/privacy',
  '/terms',
  '/help',
  '/about',
  '/features',
  '/pricing',
  '/feedback',
  '/legal',
  '/register',
  '/login',
  '/unsubscribe',
  '/favicon.ico',
  '/android-chrome-192x192.webp',
  '/android-chrome-512x512.webp',
  '/apple-touch-icon.webp',
  '/placeholder.webp',
  '/usage-image.webp',
  '/blogger-image.webp',
  '/information',
  '/information/:path',
]

/**
 * パスが公開ルートかどうかを判定
 */
function isPublicPath(pathname: string): boolean {
  return (
    PUBLIC_PATHS.includes(pathname) ||
    pathname.startsWith('/_next/') ||
    pathname.startsWith('/api/')
  )
}

/**
 * Firebase ID Tokenの形式チェック（簡易版）
 */
function isValidTokenFormat(token: string): boolean {
  if (!token) return false

  // JWT形式の基本チェック（3つのパートがピリオドで区切られている）
  const parts = token.split('.')
  return parts.length === 3
}

/**
 * Firebase ID Tokenの有効期限チェック（簡易版）
 */
function isTokenExpired(token: string): boolean {
  try {
    const payload = JSON.parse(Buffer.from(token.split('.')[1], 'base64').toString('binary'))
    const currentTime = Math.floor(Date.now() / 1000)
    return payload.exp && payload.exp < currentTime
  } catch {
    // デコードに失敗した場合は期限切れとみなす
    return true
  }
}

export async function middleware(request: NextRequest): Promise<NextResponse> {
  const pathname = new URL(request.url).pathname

  // 公開ルートの場合はそのまま通す
  if (isPublicPath(pathname)) {
    // トップページ（ルートパス）の場合、ログイン済みかどうかを確認
    if (pathname === '/') {
      // Cookieから認証情報を取得
      const sessionToken = request.cookies.get('__session')?.value
      const userCookie = request.cookies.get('tavinikkiy-user')?.value

      // セッションが存在しており、有効であれば /travels にリダイレクト
      if (sessionToken && userCookie) {
        try {
          // セッションの正当性確認
          if (isValidTokenFormat(sessionToken) && !isTokenExpired(sessionToken)) {
            // ユーザー情報も確認
            const userData = decodeUserDataFromCookie(userCookie)
            if (userData && userData.id) {
              return NextResponse.redirect(new URL('/travels', request.url))
            }
          }
        } catch (error) {
          console.warn('Session validation error in middleware:', error)
          // エラーが発生した場合はそのままトップページを表示
        }
      }
    }
    return NextResponse.next()
  }

  // 認証状態をチェック
  const userCookie = request.cookies.get('tavinikkiy-user')?.value
  let isAuthenticated = false
  let userId = ''
  let shouldClearUserCookie = false

  // ユーザー情報の取得と認証状態の判定
  if (userCookie) {
    try {
      // UTF-8対応のBase64デコードでユーザーデータを取得
      const userData = decodeUserDataFromCookie(userCookie)
      if (userData.id) {
        userId = userData.id
        isAuthenticated = true
      }
    } catch (error) {
      console.warn('Failed to parse user cookie:', error)
      // 破損したクッキーを削除するフラグを設定
      shouldClearUserCookie = true
    }
  }

  // 認証されていない場合の処理
  if (!isAuthenticated) {
    // ただし、Firebase認証処理中の可能性があるページは例外的に通す
    // これらのページではクライアントサイドで最終的な認証チェックが行われる
    const authProcessingPaths = ['/sign-up', '/register', '/login']
    const isAuthProcessingPath = authProcessingPaths.some(path => pathname.startsWith(path))

    if (!isAuthProcessingPath) {
      const homeUrl = new URL('/', request.url)

      // ログイン後のリダイレクト先として元のURLを保存（オープンリダイレクト対策）
      // 同一オリジンのパスのみ許可
      if (pathname.startsWith('/') && !pathname.startsWith('//')) {
        homeUrl.searchParams.set('redirect', pathname)
      }

      const response = NextResponse.redirect(homeUrl)
      response.headers.set('x-auth-status', 'unauthenticated')

      // 破損したクッキーをクリア
      if (shouldClearUserCookie) {
        response.cookies.set('tavinikkiy-user', '', {
          expires: new Date(0),
          path: '/',
          secure: true,
          sameSite: 'strict',
        })
      }

      return response
    }
  }

  // 認証済みまたは認証処理中の場合、リクエストを続行
  const response = NextResponse.next()

  // レスポンスヘッダーに認証情報を設定
  response.headers.set('x-auth-status', isAuthenticated ? 'authenticated' : 'processing')

  if (userId) {
    response.headers.set('x-tavinikkiy-user-id', userId)
  }

  // 破損したクッキーをクリア
  if (shouldClearUserCookie) {
    response.cookies.set('tavinikkiy-user', '', {
      expires: new Date(0),
      path: '/',
      secure: true,
      sameSite: 'strict',
    })
  }

  return response
}

// ミドルウェアを実行するパスを設定
export const config = {
  matcher: [
    // 静的ファイル・API以外のすべてのパスで認証チェック
    '/((?!_next/static|_next/image|favicon.ico|sitemap.xml|robots.txt).*)',
    // 公開ルート・静的ファイル・API以外はすべて認証必須
    '/((?!_next/static|_next/image|favicon.ico|sitemap.xml|robots.txt).*)',
    // 認証が必要なルート
    '/upload',
    '/dashboard',
    '/videos',
    // トップページにもミドルウェアを適用
    '/',
  ],
}
