# 料金・課金要件（Billing & Pricing）

## 課金概要

トークンベースの課金システムを採用し、生成処理ごとにトークンを消費する。プランと従量課金を併用し、Stripeで決済を行う。

## トークン消費（主要項目）

| 処理 | 消費トークン | 備考 |
|------|------------:|------|
| 画像アップロード（1枚） | 5 | Cloudflare R2費用を反映 |
| 画像分析（1枚） | 10 | Vertex AI Vision API |
| 動画アップロード（10秒） | 15 | サイズに応じて変動 |
| 動画分析（10秒） | 20 | Vertex AI Video API |
| スクリプト生成 | 20 | Gemini 1.5 Pro 使用を想定 |
| 音声合成（10秒） | 15 | VOICEVOX処理 |
| 動画生成（1分） | 50 | FFmpeg処理 + Storage |

> 注: 上記は初期設計の想定値。運用に応じてトークンコストは調整する。

## プラン比較（要約）

| 項目 | フリー | 月額 | 年額 | 従量課金 |
|------|-------:|-----:|-----:|--------:|
| 月間トークン | 50 | 300 | 300 | 都度購入 |
| 動画最大長 | 30秒 | 3分 | 5分 | 1分 |
| 透かし | あり | なし | なし | あり |
| 優先処理 | - | ✓ | ✓✓ | - |
| 字幕編集 | - | ✓ | ✓ | - |
| 共有機能 | 制限 | ✓ | ✓ | 制限 |
| 月額料金 | 無料 | ¥1,980 | ¥1,650/月相当 (年額¥19,800) | - |
| トークン単価 | - | ¥6.6/token | ¥5.5/token | ¥10/token |

## トークン購入例
- 100トークン: ¥980
- 500トークン: ¥4,500（10%オフ）
- 1,000トークン: ¥8,000（20%オフ）

## 決済・購買フロー（概要）
1. フロントエンドでプラン／トークンを選択し `POST /api/billing/create-checkout` を呼び出す
2. バックエンドはStripe Checkout Sessionを作成して `checkout_url` を返す
3. ユーザーが決済を完了すると、Stripeが `checkout.session.completed` イベントをWebhookで送信
4. Cloud FunctionsがWebhookを受信し、決済を検証してTiDBに `Payment` レコードを作成、`TokenTransaction` を記録して `User.TokenBalance` を更新する

## トークン消費と返却ポリシー
- トークンは、処理開始前に仮引き（reserve）し、処理完了で確定する運用を推奨
- 処理失敗時は仮引きを解除してトークンを返却する
- 二重消費防止のため、べき等性キー（request_id）を利用する

## 監査・ログ
- すべてのトークン変動は `TokenTransaction` に記録
- 重要な決済イベントは `Payment` に保存し、Stripe Payment ID 等を保持

## 環境変数（課金関連）
- `STRIPE_SECRET_KEY`, `STRIPE_WEBHOOK_SECRET`, `STRIPE_PUBLISHABLE_KEY`
- StripeのPrice IDは環境変数で管理（`STRIPE_PRICE_100_TOKENS` 等）
