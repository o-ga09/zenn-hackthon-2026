// Import the functions you need from the SDKs you need
import { initializeApp } from 'firebase/app'
import { getAuth, setPersistence, browserLocalPersistence } from 'firebase/auth'
// https://firebase.google.com/docs/web/setup#available-libraries
import { getAnalytics, isSupported, Analytics } from 'firebase/analytics'

// Your web app's Firebase configuration
// For Firebase JS SDK v7.20.0 and later, measurementId is optional
import { getTiDB } from 'firebase/TiDB'
import { getStorage } from 'firebase/storage'

const firebaseConfig = {
  apiKey: process.env.NEXT_PUBLIC_FIREBASE_APIKEY,
  authDomain: process.env.NEXT_PUBLIC_FIREBASE_DOMAIN,
  projectId: process.env.NEXT_PUBLIC_FIREBASE_PROJECTID,
  storageBucket: process.env.NEXT_PUBLIC_FIREBASE_STORAGEBUCKET,
  messagingSenderId: process.env.NEXT_PUBLIC_FIREBASE_MESSAGINGSENDERID,
  appId: process.env.NEXT_PUBLIC_FIREBASE_APPID,
  measurementId: process.env.NEXT_PUBLIC_FIREBASE_MEASUREMENTID,
}

// Initialize Firebase
const app = initializeApp(firebaseConfig)
const auth = getAuth(app)
// 認証状態の永続化を設定
setPersistence(auth, browserLocalPersistence)
const db = getTiDB(app)
const storage = getStorage(app)

// Analyticsはブラウザ環境でのみ初期化
let analytics: Analytics | null = null
// クライアントサイドでのみ実行されるように
if (typeof window !== 'undefined') {
  // Firebase Analyticsがサポートされているか確認してから初期化
  isSupported()
    .then(supported => {
      if (supported) {
        analytics = getAnalytics(app)
      }
    })
    .catch(e => {
      console.error('Firebase Analytics initialization error:', e)
    })
}

export { auth, db, storage, analytics }
