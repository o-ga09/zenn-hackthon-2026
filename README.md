# AI旅行Vlog自動生成アプリ

**第4回 Agentic AI Hackathon with Google Cloud 提出作品**

ユーザーがアップロードした旅行の写真・動画から、AIが自動的に魅力的な旅行Vlogを生成するWebアプリケーション。

[公式ページ](https://zenn.dev/hackathons/google-cloud-japan-ai-hackathon-vol4)

## コンセプト

「旅行の振り返りを AI で、簡単に、便利に、そして、もっとエモく」

## 主要機能

- **自動Vlog生成**: AIエージェントがアップロードされたメディアから自律的にショート動画を生成
- **音声合成**: ずんだもんの声でナレーション
- **ユーザー認証**: Firebase Authによる安全なユーザー管理
- **メディア管理**: 旅行写真・動画のアップロード、整理、管理
- **トークンベース課金**: フリーミアム課金モデル

## アーキテクチャ

| 構成要素        | 技術スタック                       |
| --------------- | ---------------------------------- |
| フロントエンド  | Next.js 15 + React 19              |
| バックエンド    | Go + Echo                          |
| AI エージェント | Firebase Genkit for Go + Vertex AI |
| データベース    | TiDB (MySQL互換)                   |
| ユーザー認証    | Firebase Auth                      |
| ストレージ      | Cloudflare R2 (S3互換)             |
| ホスティング    | Google Cloud Run                   |

## プロジェクト構造

```
├── backend/           # Goバックエンドアプリケーション
├── frontend/          # Next.jsフロントエンドアプリケーション
├── docs/             # プロジェクトドキュメント
│   ├── design/       # システム設計ドキュメント
│   └── requirements/ # 要件・仕様書
└── .kiro/            # Kiro AI アシスタント設定
    ├── hooks/        # 自動化フック設定
    └── steering/     # 開発ガイドライン・ルール
```

## ドキュメント

### 要件・仕様書 (`docs/requirements/`)

- [001_project_overview.md](docs/requirements/001_project_overview.md) - プロジェクト概要
- [002_features.md](docs/requirements/002_features.md) - 機能要件
- [003_technical_requirements.md](docs/requirements/003_technical_requirements.md) - 技術要件
- その他の詳細仕様書

### 設計ドキュメント (`docs/design/`)

- [001_system_design.md](docs/design/001_system_design.md) - システム設計
- [002_development_process.md](docs/design/002_development_process.md) - 開発プロセス設計
- [003_pr_template_design.md](docs/design/003_pr_template_design.md) - PRテンプレート設計

## 開発環境セットアップ

### バックエンド

```bash
cd backend
make up              # Dockerサービス起動
make migrate-up      # データベースマイグレーション実行
make seed           # データベースシード
```

### フロントエンド

```bash
cd frontend
pnpm install        # 依存関係インストール
pnpm dev           # 開発サーバー起動
```

## 技術スタック詳細

### フロントエンド

- **フレームワーク**: Next.js 15.5.9 (App Router)
- **ランタイム**: React 19.1.0
- **スタイリング**: TailwindCSS 4 + shadcn/ui
- **状態管理**: TanStack Query
- **フォーム**: React Hook Form + Zod
- **パッケージマネージャー**: pnpm

### バックエンド

- **言語**: Go 1.25.5
- **Webフレームワーク**: Echo
- **ORM**: GORM + MySQLドライバー
- **AIフレームワーク**: Firebase Genkit for Go
- **監視**: Sentry, OpenTelemetry

### 開発ツール

- **AI開発支援**: Kiro AI アシスタント + 自動化フック
- **開発フロー**: 自動PR作成ガイド、ドキュメント同期
- **コンテナ化**: Docker Compose
- **ホットリロード**: Air (Goバックエンド)
- **テスト**: Go標準テスト, Vitest (フロントエンド)

## ライセンス

Copyright © 2026 o-ga09
