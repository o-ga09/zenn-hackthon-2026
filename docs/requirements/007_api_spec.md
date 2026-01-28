# API仕様

## 6. API仕様

### 6.1 エンドポイント

#### POST /api/agent/create-vlog
**説明:** 旅行Vlogを自動生成（Veo3使用）

**認証:** Firebase ID Token

**リクエスト:**
```
Content-Type: multipart/form-data
Authorization: Bearer <firebase_id_token>

フォームフィールド:
- files: File[]          (必須) メディアファイル（複数可、画像・動画）
- title: string          (任意) VLogのタイトル
- travelDate: string     (任意) 旅行日（YYYY-MM-DD形式）
- destination: string    (任意) 旅行先
- theme: string          (任意) テーマ（adventure/relaxing/romantic/family）
- musicMood: string      (任意) BGMの雰囲気
- duration: int          (任意) 目標再生時間（秒、デフォルト8）
- transition: string     (任意) トランジション効果（fade/slide/zoom）
```

**curlでの使用例:**
```bash
curl -X POST http://localhost:8080/api/agent/create-vlog \
  -H "Authorization: Bearer <firebase_id_token>" \
  -F "files=@photo1.jpg" \
  -F "files=@photo2.jpg" \
  -F "files=@video1.mp4" \
  -F "title=沖縄旅行の思い出" \
  -F "destination=沖縄" \
  -F "travelDate=2026-01-15" \
  -F "theme=adventure" \
  -F "musicMood=upbeat" \
  -F "duration=8"
```

**レスポンス (JSON):**
```json
{
  "videoId": "01HXYZ123ABC",
  "videoUrl": "https://r2.example.com/users/user123/vlogs/01HXYZ123ABC.mp4",
  "shareUrl": "https://tavinikkiy.example.com/share/ABCDEF",
  "thumbnailUrl": "https://r2.example.com/thumbnails/01HXYZ123ABC.jpg",
  "duration": 8.0,
  "title": "沖縄旅行の思い出",
  "description": "青い海と白い砂浜で過ごした最高の休日",
  "subtitles": [
    {
      "startTime": 0.0,
      "endTime": 2.5,
      "text": "美しい沖縄の海で癒しのひととき"
    }
  ],
  "analytics": {
    "locations": ["沖縄", "ビーチ"],
    "activities": ["海水浴", "シュノーケリング"],
    "mood": "relaxing",
    "highlights": ["美しい沖縄の海で癒しのひととき"],
    "mediaCount": 3
  }
}
```

**エラーレスポンス:**
```json
{
  "error": "No files uploaded. Please upload at least one media file."
}
```

#### POST /api/agent/analyze-media
**説明:** 単一メディアを分析

**認証:** Firebase ID Token

**リクエスト:**
```json
{
  "fileId": "media_001",
  "url": "https://storage.example.com/images/photo1.jpg",
  "type": "image",
  "contentType": "image/jpeg"
}
```

**レスポンス:**
```json
{
  "description": "青い海と白い砂浜の風景",
  "landmarks": ["沖縄", "ビーチ"],
  "activities": ["海水浴", "シュノーケリング"],
  "mood": "relaxing",
  "suggestedCaption": "美しい沖縄の海で癒しのひととき",
  "fileId": "media_001"
}
```

#### POST /api/create-vlog (レガシー・SSE対応)
**説明:** 旅行Vlogを自動生成（要認証・トークン消費）- SSE進捗通知版

**認証:** Firebase ID Token

**リクエスト:**
```
Content-Type: multipart/form-data
Authorization: Bearer <firebase_id_token>

files: File[] (画像・動画ファイル)
```

**レスポンス (Server-Sent Events):**
```
Content-Type: text/event-stream

data: {"type": "token_check", "required": 150, "available": 200}
data: {"type": "status", "message": "画像をアップロード中..."}
data: {"type": "tool_call", "tool": "upload_media", "status": "running"}
data: {"type": "tool_call", "tool": "upload_media", "status": "completed", "result": {...}}
data: {"type": "status", "message": "画像を分析中..."}
...
data: {"type": "token_consumed", "amount": 150, "remaining": 50}
data: {"type": "completed", "video_url": "https://...", "share_url": "https://..."}
```

**エラーレスポンス:**
```json
{
  "error": "insufficient_tokens",
  "message": "トークンが不足しています",
  "required": 150,
  "available": 30
}
```

#### GET /api/user/tokens
**説明:** ユーザーのトークン残高を取得

**認証:** Firebase ID Token

**レスポンス:**
```json
{
  "balance": 230,
  "plan": "monthly",
  "next_refill": "2026-02-01T00:00:00Z"
}
```

#### GET /api/user/transactions
**説明:** トークン利用履歴を取得

**認証:** Firebase ID Token

**クエリパラメータ:**
- `limit`: 取得件数（デフォルト: 20）
- `offset`: オフセット

**レスポンス:**
```json
{
  "transactions": [
    {
      "id": "txn_123",
      "type": "consumption",
      "amount": -150,
      "balance": 230,
      "description": "動画生成（3ファイル）",
      "created_at": "2026-01-06T10:30:00Z"
    },
    {
      "id": "txn_122",
      "type": "purchase",
      "amount": 500,
      "balance": 380,
      "description": "500トークンパック購入",
      "created_at": "2026-01-05T15:20:00Z"
    }
  ],
  "total": 45
}
```

#### POST /api/billing/create-checkout
**説明:** Stripe Checkout Sessionを作成

**認証:** Firebase ID Token

**リクエスト:**
```json
{
  "type": "token_purchase",  // or "subscription"
  "plan": "500_tokens",       // or "monthly", "yearly"
  "success_url": "https://app.example.com/success",
  "cancel_url": "https://app.example.com/cancel"
}
```

**レスポンス:**
```json
{
  "checkout_url": "https://checkout.stripe.com/...",
  "session_id": "cs_test_..."
}
```

#### POST /api/billing/create-portal
**説明:** Stripeカスタマーポータルセッションを作成

**認証:** Firebase ID Token

**レスポンス:**
```json
{
  "portal_url": "https://billing.stripe.com/..."
}
```

#### POST /api/webhook/stripe
**説明:** Stripe Webhookエンドポイント

**認証:** Stripe Signature

**処理するイベント:**
- `checkout.session.completed`: 決済完了時
- `customer.subscription.created`: サブスク作成時
- `customer.subscription.updated`: サブスク更新時
- `customer.subscription.deleted`: サブスク削除時
- `invoice.payment_succeeded`: 請求成功時
- `invoice.payment_failed`: 請求失敗時

### 6.2 Genkit Flow

#### 6.2.1 createVlog Flow

**入力スキーマ:**
```json
{
  "uid": "user123",
  "files": ["file1.jpg", "file2.jpg", "video1.mp4"]
}
```

**出力スキーマ:**
```json
{
  "video_id": "abc123",
  "video_url": "https://storage.googleapis.com/.../video.mp4",
  "share_url": "https://your-app.run.app/v/abc123",
  "duration": 45.2,
  "thumbnail": "https://storage.googleapis.com/.../thumb.jpg",
  "tokens_used": 150,
  "log": "エージェントの実行ログ"
}
```
