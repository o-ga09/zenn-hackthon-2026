# Stripe決済導入ガイド
## Go (Echo) + Next.js Webアプリケーション

---

## 目次

1. [概要](#概要)
2. [前提条件](#前提条件)
3. [Stripeアカウントのセットアップ](#stripeアカウントのセットアップ)
4. [バックエンド実装（Go + Echo）](#バックエンド実装go--echo)
5. [フロントエンド実装（Next.js）](#フロントエンド実装nextjs)
6. [環境変数の設定](#環境変数の設定)
7. [Webhookの設定とテスト](#webhookの設定とテスト)
8. [デプロイ前チェックリスト](#デプロイ前チェックリスト)
9. [本番環境へのデプロイ](#本番環境へのデプロイ)
10. [トラブルシューティング](#トラブルシューティング)

---

## 概要

このドキュメントでは、Go（Echoフレームワーク）をバックエンド、Next.jsをフロントエンドとするWebアプリケーションにStripe決済を導入する手順を説明します。

### システム構成

- **バックエンド**: Go言語 + Echoフレームワーク（RESTful API）
- **フロントエンド**: Next.js（React）
- **決済処理**: Stripe Checkout + Webhook

### 処理フロー

1. ユーザーがフロントエンド（Next.js）で購入ボタンをクリック
2. バックエンド（Go）にリクエストを送信
3. バックエンドがStripe Checkout Sessionを作成
4. ユーザーをStripe Checkoutページにリダイレクト
5. ユーザーが決済情報を入力
6. 決済完了後、Stripeがバックエンドに結果を通知（Webhook）
7. 成功/失敗ページにリダイレクト

---

## 前提条件

### 開発環境

- Go 1.20以上
- Node.js 18以上
- npm または yarn
- Git

### 必要なアカウント

- Stripeアカウント（無料で作成可能）

---

## Stripeアカウントのセットアップ

### 3.1 アカウント作成

1. [Stripe公式サイト](https://dashboard.stripe.com/register) にアクセス
2. メールアドレスとパスワードを入力してアカウント作成
3. メール認証を完了

### 3.2 APIキーの取得

1. Stripeダッシュボードにログイン
2. 左メニューから「開発者」→「APIキー」を選択
3. 以下の2つのキーを確認・コピー
   - **公開可能キー（Publishable key）**: `pk_test_...` で始まる
   - **シークレットキー（Secret key）**: `sk_test_...` で始まる

> ⚠️ **注意**: シークレットキーは絶対に公開しないでください。Gitリポジトリにコミットすることも避けてください。

### 3.3 商品と価格の作成

1. 左メニューから「商品カタログ」→「商品」を選択
2. 「商品を追加」ボタンをクリック
3. 商品名、説明、価格を入力
4. 作成した価格のIDをコピー（`price_...` で始まる）

---

## バックエンド実装（Go + Echo）

### 4.1 必要なパッケージのインストール

```bash
go get github.com/labstack/echo/v4
go get github.com/stripe/stripe-go/v76
go get github.com/joho/godotenv
```

### 4.2 プロジェクト構成

```
backend/
├── main.go
├── handlers/
│   ├── stripe.go
│   └── webhook.go
├── .env
└── go.mod
```

### 4.3 main.go の実装

```go
package main

import (
    "log"
    "net/http"
    "os"

    "github.com/joho/godotenv"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "github.com/stripe/stripe-go/v76"
)

func main() {
    // 環境変数の読み込み
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }

    // Stripe APIキーの設定
    stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

    e := echo.New()

    // ミドルウェア
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())
    e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
        AllowOrigins: []string{"http://localhost:3000"},
        AllowMethods: []string{http.MethodGet, http.MethodPost},
    }))

    // ルーティング
    e.POST("/create-checkout-session", createCheckoutSession)
    e.POST("/webhook", handleWebhook)

    e.Logger.Fatal(e.Start(":8080"))
}
```

### 4.4 Checkout Session作成ハンドラー

```go
package main

import (
    "net/http"

    "github.com/labstack/echo/v4"
    "github.com/stripe/stripe-go/v76"
    "github.com/stripe/stripe-go/v76/checkout/session"
)

type CreateCheckoutSessionRequest struct {
    PriceID string `json:"priceId"`
}

func createCheckoutSession(c echo.Context) error {
    var req CreateCheckoutSessionRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": err.Error(),
        })
    }

    params := &stripe.CheckoutSessionParams{
        Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
        LineItems: []*stripe.CheckoutSessionLineItemParams{
            {
                Price:    stripe.String(req.PriceID),
                Quantity: stripe.Int64(1),
            },
        },
        SuccessURL: stripe.String("http://localhost:3000/success?session_id={CHECKOUT_SESSION_ID}"),
        CancelURL:  stripe.String("http://localhost:3000/cancel"),
    }

    s, err := session.New(params)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": err.Error(),
        })
    }

    return c.JSON(http.StatusOK, map[string]string{
        "sessionId": s.ID,
        "url":       s.URL,
    })
}
```

### 4.5 Webhookハンドラー

```go
package main

import (
    "io"
    "log"
    "net/http"
    "os"

    "github.com/labstack/echo/v4"
    "github.com/stripe/stripe-go/v76/webhook"
)

func handleWebhook(c echo.Context) error {
    payload, err := io.ReadAll(c.Request().Body)
    if err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": "Invalid payload",
        })
    }

    signature := c.Request().Header.Get("Stripe-Signature")
    webhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")

    event, err := webhook.ConstructEvent(payload, signature, webhookSecret)
    if err != nil {
        log.Printf("Webhook signature verification failed: %v", err)
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": err.Error(),
        })
    }

    switch event.Type {
    case "checkout.session.completed":
        log.Println("Payment successful!")
        // データベースに保存、メール送信などの処理
    case "payment_intent.succeeded":
        log.Println("Payment intent succeeded")
    default:
        log.Printf("Unhandled event type: %s", event.Type)
    }

    return c.JSON(http.StatusOK, map[string]string{
        "status": "success",
    })
}
```

---

## フロントエンド実装（Next.js）

### 5.1 必要なパッケージのインストール

```bash
npm install @stripe/stripe-js
```

### 5.2 プロジェクト構成

```
frontend/
├── app/
│   ├── page.tsx
│   ├── success/
│   │   └── page.tsx
│   └── cancel/
│       └── page.tsx
├── .env.local
└── package.json
```

### 5.3 app/page.tsx の実装

```typescript
'use client';

import { useState } from 'react';
import { loadStripe } from '@stripe/stripe-js';

const stripePromise = loadStripe(
  process.env.NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY!
);

export default function Home() {
  const [loading, setLoading] = useState(false);

  const handleCheckout = async () => {
    setLoading(true);

    try {
      const response = await fetch(
        'http://localhost:8080/create-checkout-session',
        {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            priceId: 'price_xxxxxxxxxxxxx', // StripeのPrice ID
          }),
        }
      );

      const { url } = await response.json();
      
      // Stripe Checkoutへリダイレクト
      window.location.href = url;
    } catch (error) {
      console.error('Error:', error);
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center">
      <button
        onClick={handleCheckout}
        disabled={loading}
        className="bg-blue-600 text-white px-6 py-3 rounded-lg disabled:opacity-50"
      >
        {loading ? '処理中...' : '購入する'}
      </button>
    </div>
  );
}
```

### 5.4 app/success/page.tsx

```typescript
'use client';

import { useSearchParams } from 'next/navigation';

export default function Success() {
  const searchParams = useSearchParams();
  const sessionId = searchParams.get('session_id');

  return (
    <div className="min-h-screen flex items-center justify-center">
      <div className="text-center">
        <h1 className="text-3xl font-bold mb-4">支払いが完了しました</h1>
        <p className="text-gray-600">ありがとうございます！</p>
        {sessionId && (
          <p className="text-sm text-gray-400 mt-2">
            セッションID: {sessionId}
          </p>
        )}
      </div>
    </div>
  );
}
```

### 5.5 app/cancel/page.tsx

```typescript
export default function Cancel() {
  return (
    <div className="min-h-screen flex items-center justify-center">
      <div className="text-center">
        <h1 className="text-3xl font-bold mb-4">支払いがキャンセルされました</h1>
        <p className="text-gray-600">再度お試しください。</p>
        <a
          href="/"
          className="inline-block mt-4 bg-blue-600 text-white px-6 py-3 rounded-lg"
        >
          ホームに戻る
        </a>
      </div>
    </div>
  );
}
```

---

## 環境変数の設定

### 6.1 バックエンド（.env）

`backend/.env` ファイルを作成し、以下を記述します。

```env
STRIPE_SECRET_KEY=sk_test_xxxxxxxxxxxxx
STRIPE_WEBHOOK_SECRET=whsec_xxxxxxxxxxxxx
```

### 6.2 フロントエンド（.env.local）

`frontend/.env.local` ファイルを作成し、以下を記述します。

```env
NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY=pk_test_xxxxxxxxxxxxx
```

> 💡 **重要**: `.env` および `.env.local` ファイルは `.gitignore` に追加し、Gitリポジトリにコミットしないようにしてください。

---

## Webhookの設定とテスト

### 7.1 Stripe CLIのインストール

開発環境でWebhookをテストするため、Stripe CLIをインストールします。

#### macOS
```bash
brew install stripe/stripe-cli/stripe
```

#### Windows
```bash
scoop bucket add stripe https://github.com/stripe/scoop-stripe-cli
scoop install stripe
```

#### Linux
公式サイトからバイナリをダウンロードしてインストールします。

### 7.2 Stripe CLIでログイン

```bash
stripe login
```

ブラウザが開くので、Stripeアカウントでログインします。

### 7.3 Webhookのフォワーディング

以下のコマンドでWebhookをローカル環境に転送します。

```bash
stripe listen --forward-to localhost:8080/webhook
```

実行すると、Webhook署名シークレット（`whsec_...`）が表示されます。これを `.env` ファイルの `STRIPE_WEBHOOK_SECRET` に設定します。

### 7.4 Webhookのテスト

別のターミナルで以下のコマンドを実行し、テストイベントを送信します。

```bash
stripe trigger checkout.session.completed
```

バックエンドのログに「Payment successful!」と表示されれば成功です。

---

## デプロイ前チェックリスト

- [ ] **環境変数の確認**: 本番用のStripe APIキーに変更したか
- [ ] **HTTPS化**: 本番環境ではHTTPSを使用しているか
- [ ] **CORS設定**: 本番環境のドメインをAllowOriginsに追加したか
- [ ] **エラーハンドリング**: 適切なエラーメッセージを表示しているか
- [ ] **ログ記録**: 決済ログを適切に記録しているか
- [ ] **Webhook署名検証**: 必ず有効化しているか
- [ ] **タイムアウト設定**: API呼び出しのタイムアウトを設定したか
- [ ] **セキュリティ**: 本番環境でデバッグモードが無効化されているか

---

## 本番環境へのデプロイ

### 9.1 本番用APIキーの取得

1. Stripeダッシュボードで、テストモードを本番モードに切り替え
2. 「開発者」→「APIキー」から本番用キーを取得
3. 環境変数に本番用キーを設定

### 9.2 Webhookエンドポイントの登録

1. Stripeダッシュボードで「開発者」→「Webhook」を選択
2. 「エンドポイントを追加」をクリック
3. 本番環境のWebhook URL（例: `https://api.example.com/webhook`）を入力
4. リッスンするイベントを選択
   - `checkout.session.completed`
   - `payment_intent.succeeded`
   - `payment_intent.payment_failed`
5. 作成後、署名シークレットをコピーして環境変数に設定

### 9.3 URLの更新

コード内のURLを本番環境用に更新します。

**バックエンド（Go）**:
```go
SuccessURL: stripe.String("https://yoursite.com/success?session_id={CHECKOUT_SESSION_ID}"),
CancelURL:  stripe.String("https://yoursite.com/cancel"),
```

**フロントエンド（Next.js）**:
```typescript
const response = await fetch(
  'https://api.yoursite.com/create-checkout-session',
  // ...
);
```

**CORS設定**:
```go
AllowOrigins: []string{"https://yoursite.com"},
```

---

## トラブルシューティング

### よくある問題と解決方法

| 問題 | 解決方法 |
|------|----------|
| CORS エラー | バックエンドのCORS設定でフロントエンドのドメインを許可 |
| Webhook署名検証失敗 | `STRIPE_WEBHOOK_SECRET` が正しいか確認。Stripe CLIで取得した値を使用 |
| Checkout Session作成失敗 | Price IDが正しいか確認。Stripeダッシュボードで商品と価格が作成されているか確認 |
| リダイレクトが動作しない | Success URL と Cancel URL が正しいか確認 |
| 環境変数が読み込まれない | `.env` ファイルがプロジェクトルートに存在するか確認。`godotenv.Load()` が正しく呼ばれているか確認 |

### デバッグ方法

- ブラウザの開発者ツールでネットワークタブを確認
- バックエンドのログを確認（`e.Use(middleware.Logger())`）
- Stripeダッシュボードの「ログ」セクションでイベントを確認
- `stripe listen` コマンドでリアルタイムイベントを監視

### ログの例

```bash
# 成功時
POST /create-checkout-session 200
Payment successful!

# エラー時
POST /create-checkout-session 400
Error: Price not found
```

---

## まとめ

このドキュメントでは、Go（Echo）とNext.jsのWebアプリケーションにStripe決済を導入する手順を説明しました。開発環境での実装からテスト、本番環境へのデプロイまで、段階的に進めることで安全に導入できます。

### 次のステップ

- サブスクリプション機能の実装
- 複数の支払い方法のサポート
- 返金処理の実装
- 請求書の自動生成

### 参考リンク

- [Stripe公式ドキュメント](https://stripe.com/docs)
- [Stripe Go SDK](https://github.com/stripe/stripe-go)
- [Echo フレームワーク](https://echo.labstack.com/)
- [Next.js ドキュメント](https://nextjs.org/docs)
