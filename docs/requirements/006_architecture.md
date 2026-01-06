# システムアーキテクチャ

## 5. システムアーキテクチャ

### 5.1 全体構成

ユーザー → HTTPS → Cloud Run (Next.js フロントエンド: UI/アップロード/進捗表示/プレーヤー)
→ HTTP (SSE) → Cloud Run (Genkit Go Backend: VlogCreatorAgent / Tools / FFmpeg)

外部サービス:
- Vertex AI (Gemini 1.5 Pro/Flash)
- VOICEVOX Engine (別コンテナ)
- Cloudflare R2
- TiDB (ユーザー・トークン管理)
- Firebase Authentication
- Stripe API (決済)

### 5.2 ディレクトリ構造（提案）

ai-travel-vlog-generator/
├── frontend/                   # Next.js アプリ
├── backend/                    # Genkit for Go
├── voicevox/                   # VOICEVOX Engine
├── docs/
└── ...
