import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'
import { adminAuth } from './lib/firebase-admin'

export const runtime = 'nodejs'

export async function middleware(request: NextRequest) {
  const sessionCookie = request.cookies.get('__session')?.value || ''
  // 保護されたルートの場合
  if (
    request.nextUrl.pathname.startsWith('/dashboard') ||
    request.nextUrl.pathname.startsWith('/profile')
  ) {
    try {
      const decodedClaims = await adminAuth.verifySessionCookie(sessionCookie, true)
      if (!decodedClaims.uid) throw new Error('Unauthorized')
      return NextResponse.next()
    } catch (error) {
      return NextResponse.redirect(new URL('/', request.url))
    }
  }

  return NextResponse.next()
}

export const config = {
  matcher: ['/dashboard/:path*', '/profile/:path*'],
}
