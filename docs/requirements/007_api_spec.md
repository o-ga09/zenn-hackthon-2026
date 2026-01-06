# API仕様

## 6. API仕様

### 6.1 エンドポイント

#### POST /api/create-vlog
**説明:** 旅行Vlogを自動生成（要認証・トークン消費）

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
