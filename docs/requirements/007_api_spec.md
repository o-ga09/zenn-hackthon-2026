# API仕様

## 6. API仕様

### 6.1 エンドポイント

#### POST /api/agent/create-vlog

**説明:** AI旅行Vlog自動生成（Firebase Genkit for Go + Vertex AI + Veo3）

**認証:** Firebase ID Token または 未認証ユーザー初回利用（セッションベース）

**リクエスト:**

```
Content-Type: multipart/form-data
Authorization: Bearer <firebase_id_token> (認証ユーザーの場合)
X-Session-Token: <anonymous_session_token> (未認証ユーザーの場合)

files: File[] (画像・動画ファイル)
title: string (オプション)
description: string (オプション)
style: string (cinematic|casual|documentary)
duration: number (5-60秒)
```

**レスポンス (非同期 + SSE):**

```json
{
  "job_id": "job_abc123",
  "vlog_id": "vlog_xyz789",
  "status": "queued",
  "sse_endpoint": "/api/sse/connect/job_abc123"
}
```

**SSE進行状況イベント:**

```
Content-Type: text/event-stream

data: {"type": "status", "message": "画像をアップロード中...", "progress": 10}
data: {"type": "step", "current_step": "image_analysis", "estimated_time": "2分"}
data: {"type": "ai_analysis", "tool": "gemini3_analyze", "status": "running"}
data: {"type": "ai_analysis", "tool": "gemini3_analyze", "status": "completed", "result": {...}}
data: {"type": "video_generation", "tool": "veo3_generate", "status": "running"}
data: {"type": "completed", "video_url": "https://...", "thumbnail_url": "https://..."}
```

#### GET /api/agent/vlog/{id}/status

**説明:** Vlog生成ジョブのステータス確認

**レスポンス:**

```json
{
  "job_id": "job_abc123",
  "vlog_id": "vlog_xyz789",
  "status": "processing|completed|failed",
  "progress_percentage": 75,
  "current_step": "video_generation",
  "estimated_completion_time": "2026-01-22T15:30:00Z",
  "video_url": "https://...",
  "error_message": null
}
```

#### GET /api/sse/connect/{jobId}

**説明:** SSE接続エンドポイント（リアルタイム進行状況）

**認証:** Firebase ID Token または セッショントークン

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
  "type": "token_purchase", // or "subscription"
  "plan": "500_tokens", // or "monthly", "yearly"
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
