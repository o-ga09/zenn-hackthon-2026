# 旅行Vlog自動生成エージェント - 技術要件

## 3. 技術要件

### 3.1 技術スタック

#### 3.1.1 フロントエンド
- フレームワーク: Next.js 14+ (App Router)
- 言語: TypeScript
- UI: Tailwind CSS
- 状態管理: React Hooks (useState, useEffect)
- ファイルアップロード: FormData API (multipart/form-data)
- ストリーミング: Server-Sent Events (SSE) ※将来対応予定

#### 3.1.2 バックエンド
- フレームワーク: Echo (HTTP), Gorm (ORM)
- 言語: Go 1.25
- データベース: TiDB

#### 3.1.3 AIエージェント
- フレームワーク: Firebase Genkit for Go v1.3.0
- アーキテクチャ: シングルエージェント + ツール
- エージェント名: VlogCreatorAgent
- LLMモデル: Vertex AI Gemini 2.5 Flash（メディア分析・スクリプト生成）
- 動画生成: Vertex AI Veo 3.1 (veo-3.1-fast-generate-001)
- プロンプト管理: Dotprompt

#### 3.1.4 Google Cloud サービス
- Vertex AI Gemini: メディア分析、スクリプト生成
- Vertex AI Veo 3.1: 動画生成（8秒、720p）
- Google Cloud Storage (GCS): Veo出力の一時保存
- Cloudflare R2: 最終メディアファイル保存
- Cloud Run: Next.jsアプリとGenkit Goサーバーのホスティング
- TiDB: ユーザー情報管理、トークン残高管理、利用履歴保存、プラン情報管理
- Firebase Authentication: ユーザー認証
- Stripe API: 決済処理
- リージョン: asia-northeast1 (東京) / us-central1 (Veo)

### 3.2 AIエージェント設計

#### 3.2.1 エージェントタイプ
シングルエージェント + ツール

理由:
- ワークフローが線形（逐次処理が多い）
- 文脈保持が重要
- 実装・デバッグが容易
- ハッカソン期間に適している

#### 3.2.2 エージェント構成

VlogCreatorAgent (Gemini 2.5 Flash)
├─ Tool: upload_media         # メディアをR2にアップロード
├─ Tool: analyze_media        # Gemini Vision APIで分析
├─ Tool: generate_vlog        # Veo3で動画生成 (GCS→R2転送)
├─ Tool: get_share_url        # 共有URL取得
└─ Tool: generate_thumbnail   # サムネイル生成
