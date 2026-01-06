# 旅行Vlog自動生成エージェント - 技術要件

## 3. 技術要件

### 3.1 技術スタック

#### 3.1.1 フロントエンド
- フレームワーク: Next.js 14+ (App Router)
- 言語: TypeScript
- UI: Tailwind CSS
- 状態管理: React Hooks (useState, useEffect)
- ファイルアップロード: FormData API
- ストリーミング: Server-Sent Events (SSE)

#### 3.1.2 バックエンド
- フレームワーク: Firebase Genkit for Go
- 言語: Go 1.25
- 動画処理: FFmpeg
- 音声合成: VOICEVOX Engine

#### 3.1.3 AIエージェント
- フレームワーク: Firebase Genkit for Go
- アーキテクチャ: シングルエージェント + ツール
- エージェント名: VlogCreatorAgent
- LLMモデル: Vertex AI Gemini 1.5 Pro

#### 3.1.4 Google Cloud サービス
- Vertex AI: Gemini 1.5 Pro（エージェント推論）、Gemini 1.5 Flash（画像・動画分析）
- Cloudflare R2: メディアファイル保存
- Cloud Run: Next.jsアプリとGenkit Goサーバーのホスティング
- TiDB: ユーザー情報管理、トークン残高管理、利用履歴保存、プラン情報管理
- Firebase Authentication: ユーザー認証
- Cloud Run: 決済Webhook処理、トークン付与処理
- Stripe API: 決済処理
- リージョン: asia-northeast1 (東京)

### 3.2 AIエージェント設計

#### 3.2.1 エージェントタイプ
シングルエージェント + ツール

理由:
- ワークフローが線形（逐次処理が多い）
- 文脈保持が重要
- 実装・デバッグが容易
- ハッカソン期間に適している

#### 3.2.2 エージェント構成

VlogCreatorAgent (Gemini 1.5 Pro)
├─ Tool: upload_media
├─ Tool: analyze_media
├─ Tool: generate_script
├─ Tool: synthesize_voice
├─ Tool: create_video
└─ Tool: share_video
