# システムアーキテクチャ

## 5. システムアーキテクチャ

### 5.1 全体構成

ユーザー → HTTPS → Cloud Run (Next.js フロントエンド: UI/アップロード/進捗表示/プレーヤー)
→ HTTP → Cloud Run (Genkit Go Backend: VlogCreatorAgent / Tools / Veo3)

外部サービス:
- Vertex AI (Gemini 2.5 Flash/Pro) - メディア分析
- Vertex AI (Veo 3.1) - 動画生成
- Google Cloud Storage (GCS) - Veo一時出力
- Cloudflare R2 - 最終メディア保存
- TiDB (ユーザー・トークン管理)
- Firebase Authentication
- Stripe API (決済)

### 5.2 動画生成フロー

```
1. ユーザーがメディアをアップロード
2. Genkit Agent がメディアを分析（Gemini Vision API）
3. Veo3 API で動画生成 → GCS一時保存
4. GCSから動画取得 → R2にアップロード
5. GCS一時ファイル削除
6. R2のURLをユーザーに返却
```

### 5.3 ディレクトリ構造

```
ai-travel-vlog-generator/
├── frontend/                   # Next.js アプリ
├── backend/                    # Genkit for Go
│   ├── cmd/
│   │   ├── api/               # APIサーバー
│   │   ├── agent/             # エージェントCLI
│   │   └── migration/         # DBマイグレーション
│   ├── internal/
│   │   ├── agent/             # エージェントインターフェース・スキーマ
│   │   ├── domain/            # ドメインモデル
│   │   ├── handler/           # HTTPハンドラー
│   │   ├── infra/
│   │   │   ├── genkit/        # Genkit実装
│   │   │   │   ├── agent.go
│   │   │   │   ├── context.go
│   │   │   │   ├── flow.go
│   │   │   │   ├── tool_*.go  # 各種ツール
│   │   │   │   └── tool_veo_client.go  # Veo3クライアント
│   │   │   ├── database/
│   │   │   └── storage/       # R2ストレージ
│   │   └── server/
│   ├── pkg/
│   │   ├── config/            # 設定・クライアント初期化
│   │   ├── errors/            # エラー定義
│   │   └── ...
│   └── prompts/               # Dotpromptテンプレート
├── docs/
└── ...
```
