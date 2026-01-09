import * as admin from 'firebase-admin'

if (!admin.apps.length) {
  // 本番環境（Cloud Run）
  if (process.env.NODE_ENV === 'production') {
    admin.initializeApp({
      credential: admin.credential.applicationDefault(),
    })
  }
  // ローカル開発環境
  else {
    try {
      admin.initializeApp({
        credential: admin.credential.cert(process.env.FIREBASE_SERVICE_ACCOUNT_KEY ?? ''),
        projectId: process.env.FIREBASE_PROJECT_ID ?? '',
      })
    } catch (error) {
      throw new Error(
        'Failed to parse FIREBASE_SERVICE_ACCOUNT_KEY. ' + 'Make sure it is valid JSON.'
      )
    }
  }
}

export const adminAuth = admin.auth()
