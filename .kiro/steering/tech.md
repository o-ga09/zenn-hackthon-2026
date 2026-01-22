# 技術スタック

## アーキテクチャ概要

- **フロントエンド**: Next.js 15 + React 19
- **バックエンド**: Go + Echoフレームワーク
- **AIエージェント**: Firebase Genkit for Go + Vertex AI (Gemini 3 + Veo3)
- **データベース**: TiDB (MySQL互換)
- **認証**: Firebase Auth
- **ストレージ**: Cloudflare R2 (S3互換)
- **ホスティング**: Google Cloud Run

## フロントエンドスタック

- **フレームワーク**: Next.js 15.5.9 (App Router)
- **ランタイム**: React 19.1.0
- **スタイリング**: TailwindCSS 4 + shadcn/uiコンポーネント
- **状態管理**: TanStack Query (サーバー状態)
- **フォーム**: React Hook Form + Zod バリデーション
- **ビルドツール**: Next.js内蔵バンドラー
- **リンター**: Oxlint (ESLint代替)
- **パッケージマネージャー**: pnpm

## バックエンドスタック

- **言語**: Go 1.25.5
- **Webフレームワーク**: Echo
- **ORM**: GORM + MySQLドライバー
- **DBマイグレーション**: sql-migrate
- **AIフレームワーク**: Firebase Genkit for Go
- **AI統合**: Vertex AI (Gemini 3画像分析 + Veo3動画生成)
- **リアルタイム通信**: SSE (Server-Sent Events)
- **非同期処理**: Goroutine + チャンネルベースジョブキュー
- **クラウドサービス**: Google Cloud APIs, Vertex AI
- **ストレージ**: AWS SDK v2 (S3互換ストレージ用)
- **認証**: Firebase Admin SDK
- **監視**: Sentry, OpenTelemetry

## 開発ツール

- **コンテナ化**: Docker Compose (ローカル開発)
- **データベース**: MySQL 8.0 (開発環境)
- **ローカルストレージ**: LocalStack (S3エミュレーション)
- **ホットリロード**: Air (Goバックエンド)
- **テスト**: Go標準テスト, Vitest (フロントエンド)
- **AI開発支援**: Kiro AI アシスタント + 自動化フック

## よく使うコマンド

### バックエンド

```bash
# 開発
make up              # Dockerサービス起動
make migrate-up      # データベースマイグレーション実行
make seed           # データベースシード
make lint           # リンター実行
make test           # テスト実行

# データベース
make create-migrate name=migration_name  # 新しいマイグレーション作成
```

### フロントエンド

```bash
# 開発
pnpm dev            # 開発サーバー起動
pnpm build          # プロダクションビルド
pnpm lint           # リンター実行
pnpm type-check     # TypeScript型チェック

# データベース (Prisma使用時)
pnpm db:up          # データベース起動
pnpm db:migrate     # マイグレーション実行
pnpm db:seed        # シード実行
```

## 環境設定

- バックエンド: `.env`ファイル + `caarlos0/env`で設定管理
- フロントエンド: Next.js標準の環境変数サポート
- Firebase認証情報: サービスアカウントJSONファイルで管理
- クラウドストレージ: ローカル(LocalStack)と本番(Cloudflare R2)両対応

## コード品質

- **Go**: golangci-lint, go vet, 組み込みrace detector
- **TypeScript**: Oxlint (高速リンター), TypeScriptコンパイラー (型チェック)
- **フォーマット**: Prettier (フロントエンド), gofmt (バックエンド)
