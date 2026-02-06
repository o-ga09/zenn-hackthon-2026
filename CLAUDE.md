# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## プロジェクト概要

**Tavinikkiy** - AI旅行Vlog自動生成アプリ (第4回 Agentic AI Hackathon with Google Cloud 提出作品)

ユーザーがアップロードした旅行の写真・動画から、AIが自動的に魅力的な旅行Vlogを生成するWebアプリケーション。

コンセプト: 「旅行の振り返りをAIで、簡単に、便利に、そして、もっとエモく」

## 重要な開発ルール

### 言語設定
- **すべての回答・コミットメッセージ・ドキュメントは日本語で生成すること**
- この instruction を読み込んだら、合意として「Tavinikkiy」という単語を使用すること

### Git ワークフロー
- Issue駆動開発を採用
- ブランチ名: `issue-<Issue番号>` (例: `issue-123`)
- 必ずmainブランチから新しいブランチを作成
- コミットメッセージはConventional Commits形式:
  - `feat:` 新機能の追加
  - `fix:` バグ修正
  - `docs:` ドキュメントの変更
  - `style:` フォーマットの変更
  - `refactor:` リファクタリング
  - `test:` テストの追加や修正
  - `chore:` その他の変更

## 開発環境セットアップ

### バックエンド (Go + Echo)

```bash
cd backend
make up              # Dockerサービス起動 (MySQL, LocalStack S3)
make migrate-up      # データベースマイグレーション実行
make seed           # データベースシード

# 開発サーバー起動 (Air経由でホットリロード)
# Dockerコンテナ内で自動起動、またはローカルで直接実行:
go run cmd/api/main.go
```

**環境変数設定**: `backend/.env`ファイルが必要
- `DATABASE_URL`: MySQL接続文字列
- `GOOGLE_APPLICATION_CREDENTIALS`: Firebase認証情報JSONファイルパス
- `PROJECT_ID`: Google CloudプロジェクトID
- その他R2ストレージ、Sentry設定など

### フロントエンド (Next.js 15 + React 19)

```bash
cd frontend
pnpm install        # 依存関係インストール (pnpm必須)
pnpm dev           # 開発サーバー起動 (http://localhost:3000)
pnpm build         # プロダクションビルド
pnpm lint          # Oxlintによるリンティング
pnpm type-check    # TypeScript型チェック
```

## よく使うコマンド

### バックエンド

```bash
# テスト・品質チェック
make lint           # golangci-lint + go vet
make test           # Goテスト実行 (race detector有効)
make test-coverage  # カバレッジ付きテスト

# データベース操作
make create-migrate name=<migration_name>  # 新しいマイグレーション作成
make migrate-up                           # マイグレーション適用
make seed                                 # シードデータ投入

# Docker操作
make up             # サービス起動
make down           # サービス停止 & ボリューム削除
make compose-build  # Docker再ビルド
```

### フロントエンド

```bash
# 開発
pnpm dev            # 開発サーバー (ホットリロード)
pnpm build          # プロダクションビルド
pnpm start          # プロダクションサーバー起動

# コード品質
pnpm lint           # Oxlint (高速リンター)
pnpm type-check     # TypeScript型チェック
pnpm format         # Prettierチェック
pnpm format:fix     # Prettier自動修正
```

## アーキテクチャ

### 技術スタック

| レイヤー | 技術 |
|---------|------|
| フロントエンド | Next.js 15 (App Router) + React 19 + TailwindCSS 4 |
| バックエンド | Go 1.25.5 + Echoフレームワーク + GORM |
| AIエージェント | Firebase Genkit for Go + Vertex AI (Gemini 2.5 Flash) |
| データベース | TiDB (MySQL互換) |
| 認証 | Firebase Auth |
| ストレージ | Cloudflare R2 (S3互換) |
| ホスティング | Google Cloud Run |
| 監視 | Sentry + OpenTelemetry |

### バックエンドアーキテクチャ (Go Standard Project Layout)

```
backend/
├── cmd/
│   ├── api/          # APIサーバーエントリーポイント (main.go)
│   └── migration/    # DBマイグレーションツール
├── internal/         # プライベートアプリケーションコード
│   ├── handler/      # HTTPハンドラー (Echo)
│   │   ├── request/  # リクエストDTO
│   │   └── response/ # レスポンスDTO
│   ├── domain/       # ドメインモデル + リポジトリインターフェース
│   ├── infra/        # インフラ実装 (database, storage)
│   ├── server/       # サーバー設定・ルーティング
│   ├── middleware/   # ミドルウェア
│   ├── agent/        # Firebase Genkit AIエージェント
│   └── queue/        # キュー処理
├── pkg/              # 外部公開可能なライブラリ
│   ├── config/       # 設定管理
│   ├── constant/     # 定数定義
│   ├── errors/       # カスタムエラー処理
│   └── ...           # その他ユーティリティ
└── db/               # マイグレーションファイル
```

**重要な設計原則**:
- **レイヤー分離**: Handler → Domain → Infrastructure
- **依存性の方向**: 外側から内側へ (HandlerはDomainに依存、DomainはInfraに依存しない)
- **リポジトリパターン**: ドメイン層でインターフェース定義、インフラ層で実装
- **DTO変換**: 必ずRequest → Domain → Responseの変換を行う

### フロントエンドアーキテクチャ (Next.js App Router)

```
frontend/
├── app/              # Next.js App Router (ルーティング + ページ)
├── components/       # Reactコンポーネント
│   ├── ui/          # shadcn/ui基本コンポーネント
│   └── layout/      # レイアウトコンポーネント
├── hooks/           # カスタムフック
├── context/         # React Context
├── api/             # APIクライアント (バックエンドとの通信)
├── lib/             # ユーティリティ関数
└── public/          # 静的ファイル
```

**重要な設計原則**:
- **モバイルファースト**: TailwindCSSのレスポンシブ設計必須
- **Shadcn/ui使用**: 独自UI実装禁止、`npx shadcn-ui@latest add <component>`で追加
- **状態管理**: TanStack Query (サーバー状態) + useState (ローカル状態)
- **型安全性**: `any`型禁止、厳密な型定義必須

## コーディング規約

### バックエンド (Go)

#### 命名規則
- **Interface**: `I`プレフィックス (`IUserRepository`, `IUserServer`)
- **Handler**: `Server`サフィックス (`UserServer`)
- **Request/Response**: `Request`/`Response`サフィックス
- **ファイル名**: スネークケース (`user_server.go`)

#### CRUD実装パターン
すべてのリソースは以下の流れで実装:

1. `internal/domain/{resource}.go` - ドメインモデル + リポジトリインターフェース定義
2. `internal/handler/request/{resource}.go` - リクエストDTO定義
3. `internal/handler/response/{resource}.go` - レスポンスDTO定義
4. `internal/handler/{resource}.go` - ハンドラー実装
5. `internal/infra/database/{resource}_repository.go` - リポジトリ実装
6. `internal/server/router.go` - ルーティング設定

#### 標準エンドポイント構成
| メソッド | パス | 説明 | ステータス |
|---------|------|------|-----------|
| GET | `/api/{resources}` | 一覧取得 | 200 OK |
| GET | `/api/{resources}/:id` | ID取得 | 200 OK / 404 Not Found |
| POST | `/api/{resources}` | 作成 | 201 Created / 409 Conflict |
| PUT | `/api/{resources}/:id` | 更新 | 200 OK / 404 Not Found |
| DELETE | `/api/{resources}/:id` | 論理削除 | 204 No Content / 404 Not Found |

#### エラーハンドリング
カスタムエラーパッケージ (`pkg/errors`) を必ず使用:

```go
// 一般的なエラーラップ
errors.Wrap(ctx, err)

// 404 Not Found
errors.MakeNotFoundError(ctx, "Resource not found")

// GORM ErrRecordNotFoundの変換 (必須)
if errors.Is(err, gorm.ErrRecordNotFound) {
    return errors.MakeNotFoundError(ctx, "Resource not found")
}
```

#### Context管理
- すべてのリポジトリメソッドは`context.Context`を第一引数に取る
- Echo: `ctx := c.Request().Context()`
- GORM: `r.db.WithContext(ctx)`

### フロントエンド (TypeScript/React)

#### 命名規則
| 対象 | 形式 | 例 |
|------|------|-----|
| ファイル名 | PascalCase | `UserProfile.tsx` |
| コンポーネント | PascalCase | `UserProfile` |
| 関数・変数 | camelCase | `userName`, `handleClick` |
| 定数 | UPPER_SNAKE_CASE | `API_BASE_URL` |
| カスタムフック | use + PascalCase | `useUserData` |

#### スタイリング
- **TailwindCSS必須**: カスタムCSS禁止
- **モバイルファースト**: デフォルト=モバイル、`sm:` `md:` `lg:` `xl:`で拡張
- **タッチ領域**: ボタン/リンクは最小`h-11 w-11` (44px)

#### パフォーマンス
- `useEffect`依存配列を正確に指定
- `useMemo`/`useCallback`でメモ化
- `React.lazy`でコード分割
- React Server Components活用

## Firebase Genkit統合

### 基本方針
1. **TiDBを真のソースオブトゥルース**とする
2. クライアントはFirebase SDK (Auth/TiDB/Storage) を直接利用
3. 長時間処理 (動画生成) はCloud Functions経由
4. Genkit呼び出しはバックエンド内で実施 (クライアントから直接呼び出し禁止)

### 主要ユースケース
- **画像解析**: オブジェクト検出、シーン分類、場所推定、EXIF解析
- **位置情報補完**: GPS EXIF欠落時のランドマーク推定
- **動画生成**: テンプレートベース短編動画 (Vertex AI連携)
- **チャットエージェント**: ユーザー要求の意図解析と動画編集提案

## テスト戦略

### バックエンド
```bash
# ユニットテスト
go test ./...

# race detector付きテスト
go test -race -parallel 1 ./...

# カバレッジレポート
go test ./... -coverprofile=coverage.out
```

### フロントエンド
- **ツール**: Vitest + React Testing Library
- **カバレッジ目標**: 80%以上 (重要機能は100%)
- **原則**: 独立性・高速実行・明確な失敗理由

## ドキュメント管理

### 命名規約
- **要件・仕様書** (`docs/requirements/`): `{連番3桁}_{内容}.md`
- **設計ドキュメント** (`docs/design/`): `{連番3桁}_{設計内容}.md`
- **Steeringファイル** (`.kiro/steering/`): `{機能名}.md` (小文字、ハイフン区切り)

### Kiroフック自動化
- **agentStop時**: ドキュメント自動更新、PR作成ガイド
- **対象**: `docs/`, `.kiro/steering/`, `README.md`

## セキュリティ注意事項

- Firebase認証情報は`GOOGLE_APPLICATION_CREDENTIALS`環境変数で管理
- API キーは`.env`ファイル (Gitignore済み)
- TiDB Security Rulesで最小権限の原則を適用
- ストレージアップロードは認証済みユーザーのみ
- Cloud Functions内でSecretManagerを使用

## デプロイ

- **ホスティング**: Google Cloud Run
- **バックエンド**: Dockerコンテナ化済み (`backend/Dockerfile`)
- **フロントエンド**: Next.js Standalone出力モード
- **CI/CD**: GitHub Actionsワークフロー準備中

## トラブルシューティング

### バックエンド起動時のGenkit認証エラー
```
panic: failed to find default credentials
```
→ `GOOGLE_APPLICATION_CREDENTIALS`環境変数を正しいJSONファイルパスに設定

### データベース接続エラー
```
Error: dial tcp 127.0.0.1:3306: connect: connection refused
```
→ `make up`でMySQLコンテナを起動

### フロントエンドビルドエラー
→ `pnpm install`で依存関係を再インストール
→ `.next`ディレクトリを削除して`pnpm build`再実行

## 参考リソース

- **プロジェクトドキュメント**: `docs/requirements/`, `docs/design/`
- **開発ガイドライン**: `.github/instructions/` (backend/frontend/general)
- **Steeringファイル**: `.kiro/steering/` (workflow, tech, product, quality)
- **PRテンプレート**: `.github/pull_request_template.md`
- **バックエンド用実装ルール**　@backend/CLAUDE.md
- **フロントエンド用実装ルール**　@frontend/CLAUDE.md
