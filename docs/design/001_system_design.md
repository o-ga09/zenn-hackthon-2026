# AI Travel Vlog Generator — System Design

この設計書は `docs/requirements` の要件定義を元に、実装のための詳細設計をまとめたものです。

## 目的と範囲

- 目的: 要件に記載されたVlog自動生成ワークフローを確実に実装するためのアーキテクチャ、データ設計、API仕様、決済フロー、運用要件を定義する。
- 範囲: フロントエンド(Next.js)、バックエンド(Genkit Go)、TiDB、Cloudflare R2、Vertex AI、VOICEVOX、FFmpeg、Stripe連携を含む。

## 高レベルアーキテクチャ

コンポーネント:

- フロントエンド (Next.js): UI、ファイルアップロード、SSE受信、Stripe Checkout トリガ
- Firebase Genkit Go バックエンド (Cloud Run): `VlogCreatorAgent` 実行、Gemini 3 + Veo3統合、非同期ジョブ管理、SSEイベント送出
- TiDB: `users`, `subscriptions`, `token_transactions`, `payments`, `vlogs`, `vlog_jobs`, `anonymous_sessions` コレクション
- Cloudflare R2: メディア保存、生成済み動画保存
- Vertex AI: Gemini 3による画像分析・ストーリー生成、Veo3による動画生成
- VOICEVOX: TTS（ずんだもん）
- FFmpeg: 最終動画合成
- Stripe + Cloud Functions: Checkout 作成、Webhook 処理

シーケンス概略 (Create Vlog):

1. フロントエンドが `POST /api/agent/create-vlog` にファイルと認証情報で要求
2. バックエンドは認証後（Firebase Auth または 未認証セッション）、VlogJobを作成し非同期処理開始
3. SSE で `image_analysis` → `story_generation` → `video_generation` → `audio_synthesis` → `final_composition` の進捗を送出
4. Firebase Genkit Agentが各ステップを自律実行：Gemini 3分析 → Veo3動画生成 → 音声合成 → FFmpeg合成
5. 処理完了でVlogレコード更新、SSEで完了通知
6. エラー発生時はジョブステータス更新し、SSEでエラー情報を返す

## コンポーネント設計

フロントエンド (Next.js)

- ルート: `/api/create-vlog` (multipart/form-data) — SSE を返す
- UI: `FileUploader`, `ProgressViewer`, `TokenDisplay`, `PaymentForm`
- Stripe Checkout の呼び出しは `POST /api/billing/create-checkout` 経由

Genkit Go バックエンド

- Flow: `createVlog` (Genkit flow) をエントリとする。
- Tools:
  - `upload_media`: Cloudflare R2へアップロード、返り値は `file_ids` とサイズ情報
  - `analyze_media`: Vertex AI 呼び出し、`MediaAnalysis` を返す
  - `generate_script`: Gemini を使いずんだモン風スクリプトを生成
  - `synthesize_voice`: VOICEVOX へリクエストして音声を取得
  - `create_video`: FFmpeg で合成、SRTを焼き込み、MP4出力
  - `share_video`: 公開URLとサムネ生成
- トークン管理: 処理開始前に `reserveTokens(uid, request_id, amount)`、処理成功で `confirmTokens(...)`、失敗で `rollbackTokens(...)`

Cloud Functions (Stripe Webhook)

- `checkout.session.completed` 受信時に Payment レコードを作成し、`TokenTransaction` を記録して `User.TokenBalance` を増加させる。
- 署名検証とべき等チェックを行う。

## データモデル（TiDB）

主要コレクション（要点）:

- `users/{uid}`
  - `uid`, `email`, `display_name`, `plan`, `token_balance`, `created_at`, `updated_at`
- `subscriptions/{id}`
  - `uid`, `plan`, `status`, `stripe_customer_id`, `stripe_sub_id`, `current_period_end`...
- `token_transactions/{id}`
  - `id`, `uid`, `type` (purchase|consumption|bonus|refund), `amount`, `balance`, `description`, `metadata`, `created_at`
- `payments/{id}`
  - `id`, `uid`, `type` (token_purchase|subscription), `amount`, `tokens_granted`, `status`, `stripe_payment_id`, `created_at`, `completed_at`
- `vlogs/{video_id}`
  - `video_id`, `uid`, `video_url`, `share_url`, `duration`, `thumbnail`, `tokens_used`, `watermarked`, `created_at`

トランザクション要件:

- トークン残高更新は TiDB トランザクションで実行し、一貫性を保つ。
- 仮引き（reserve）を `token_transactions` にタイプ `reserve` で記録し、確定時に `consumption` に変換または `rollback` を発行する。
- べき等性: `request_id` を使い同一リクエストの重複処理を防止。

## API 設計（重要エンドポイント）

- `POST /api/create-vlog`
  - 認証: Firebase ID Token
  - 入力: multipart files + optional `template`, `request_id`
  - 処理: token計算 → reserve → Genkit フロー開始 → SSE で進捗送出
  - SSE イベント例:
    - `token_check` {required, available}
    - `tool_call` {tool, status, result?}
    - `status` {message}
    - `token_consumed` {amount, remaining}
    - `completed` {video_url, share_url}
    - `error` {code, message}

- `GET /api/user/tokens`
  - 認証: Firebase ID Token
  - 出力: `{ balance, plan, next_refill }`

- `GET /api/user/transactions`
  - クエリ: `limit`, `offset`
  - 出力: トランザクションリスト

- `POST /api/billing/create-checkout`
  - 認証: Firebase ID Token
  - 入力: `{ type, plan, success_url, cancel_url }`
  - 出力: `{ checkout_url, session_id }`

- `POST /api/webhook/stripe` (Cloud Functions)
  - 認証: Stripe signature
  - 処理: イベント検証 → Payment/Subscription処理 → TokenTransaction更新

## トークン計算ロジック（概略）

- ファイルリストを受け取り、各ファイルで消費トークンを合算: 画像は `10`、動画は秒数に応じて計算
- スクリプト/音声/動画生成コストを加算して `required` を算出
- `required` が `user.TokenBalance` を超える場合は SSE で `insufficient_tokens` を返す

## 決済フロー設計（Stripe）

- フロントから `create-checkout` を呼び Checkout URL を返す
- 決済成功後、Stripe が `checkout.session.completed` を Webhook で送信
- Cloud Functions が signature 検証 → TiDB 更新（Payments, TokenTransactions, Users）
- 失敗や遅延がある場合は Cloud Functions 側でリトライ＆通知を行う

べき等性・重複対策:

- Webhook は Stripe Event ID を記録して一度だけ処理
- Checkout 作成リクエストは `client_request_id` を受けて二重作成を防止

## エラーハンドリング & リトライ

- 各ツール呼び出しはタイムアウトとリトライポリシーを持つ（指数バックオフ）
- クリティカル失敗時は Sentry / Cloud Logging にエラーログを投げる
- トークン操作は必ずトランザクションで処理し、処理失敗時は仮引きをロールバック

## セキュリティ

- 認証: Firebase Authentication を必須化、ID Token の検証をサーバーで実施
- ストレージ: Cloudflare R2 は署名付きURLでアクセス制御
- Webhook: Stripe Signature 検証
- シークレット管理: GCP Secret Manager (Stripe keys, VOICEVOX endpoint secret)

## スケーラビリティ & 非機能要件

- Cloud Run: 自動スケーリング、MVPは `min: 0, max: 10`、FFmpeg処理用にメモリ高めのインスタンスを設定
- バッチ/非同期処理: 大きな動画処理はジョブ化して Cloud Tasks/Cloud Run にオフロード
- キャッシュ: Vertex AI 呼び出し結果は短期キャッシュしてコスト削減
- 目標: 1分動画を30秒以内で生成（目安）、同時10ユーザー処理をサポート

## 監視・ロギング

- Cloud Logging: 各サービスのログを一元収集
- メトリクス: 動画生成時間、トークン消費率、決済成功率、エラー率
- アラート: 生成失敗率が閾値を超えた場合、決済失敗率が高い場合に通知

## テスト計画（設計視点）

- 単体: 各ツールの入出力のモックテスト
- 統合: Genkit Flow のステップをローカルで順次実行して挙動確認
- E2E: Next.js UI から実際のファイルアップロード→動画生成→再生
- 決済テスト: Stripe のテストモードで Checkout / Webhook の動作を検証

## デプロイと運用手順（要点）

- 環境変数を GCP Secret Manager / Cloud Build に設定
- 前提: `STRIPE_*`, `VERTEX_AI_*`, `GCS_BUCKET`, `VOICEVOX_ENDPOINT`, `FIREBASE_*` を設定
- デプロイ: `gcloud builds submit --config cloudbuild.yaml` または個別 `gcloud run deploy` を使用

## Appendix

- 参照: `docs/requirements` の各ファイル（features, api_spec, billing_and_pricing, data_models など）

---

作成日: 2026-01-06
