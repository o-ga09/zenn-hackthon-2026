# プロジェクト構造

## リポジトリレイアウト

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

## バックエンド構造 (`backend/`)

Goプロジェクトレイアウト規約 + クリーンアーキテクチャ原則に従う:

```
backend/
├── cmd/                    # アプリケーションエントリーポイント
│   ├── api/               # メインAPIサーバー
│   └── migration/         # データベースマイグレーションツール
├── internal/              # プライベートアプリケーションコード
│   ├── domain/           # ドメインモデル・インターフェース
│   ├── handler/          # HTTPハンドラー (コントローラー)
│   │   ├── request/      # リクエストDTO
│   │   └── response/     # レスポンスDTO
│   ├── infra/           # インフラストラクチャ層
│   │   ├── database/    # データベース実装
│   │   └── storage/     # ファイルストレージ実装
│   └── server/          # サーバー設定・ルーティング
├── pkg/                 # パブリックパッケージ (インポート可能)
│   ├── config/         # 設定管理
│   ├── errors/         # エラーハンドリングユーティリティ
│   ├── logger/         # ログユーティリティ
│   └── [その他utils]   # 各種ユーティリティパッケージ
├── db/                 # データベース関連ファイル
│   ├── migrations/     # SQLマイグレーションファイル
│   ├── seed/          # データベースシードファイル
│   └── sql/           # SQL初期化スクリプト
└── [設定ファイル]      # Docker, Makefile等
```

## フロントエンド構造 (`frontend/`)

Next.js 15 App Router規約に従う:

```
frontend/
├── app/                    # App Routerページ
│   ├── dashboard/         # ダッシュボードページ・コンポーネント
│   ├── profile/[user_id]/ # 動的ユーザープロフィールページ
│   ├── upload/           # メディアアップロードフロー
│   ├── videos/           # 動画ライブラリ
│   └── [ルートファイル]   # レイアウト、グローバル等
├── components/            # 再利用可能UIコンポーネント
│   ├── ui/               # shadcn/uiベースコンポーネント
│   ├── header/           # ヘッダーコンポーネント
│   ├── footer/           # フッターコンポーネント
│   ├── dialog/           # モーダルダイアログ
│   └── lp/              # ランディングページコンポーネント
├── api/                  # APIクライアント・型定義
├── context/              # Reactコンテキスト
├── lib/                  # ユーティリティライブラリ
├── types/                # TypeScript型定義
└── [設定ファイル]        # Next.js, Tailwind等
```

## 命名規約

### バックエンド (Go)

- **パッケージ**: 小文字、可能な限り単語1つ (`user`, `config`)
- **ファイル**: snake_case (`user_repository.go`, `media_analytics.go`)
- **型**: PascalCase (`User`, `MediaAnalytics`)
- **関数・メソッド**: エクスポート時PascalCase、プライベート時camelCase
- **定数**: PascalCase または パッケージレベルでSCREAMING_SNAKE_CASE

### フロントエンド (TypeScript/React)

- **ファイル**: コンポーネントはkebab-case (`user-profile.tsx`)
- **コンポーネント**: PascalCase (`UserProfile`, `TravelMemoryCard`)
- **関数**: camelCase (`getUserData`, `handleSubmit`)
- **型・インターフェース**: PascalCase (`User`, `ApiResponse`)
- **定数**: SCREAMING_SNAKE_CASE (`API_BASE_URL`)

## アーキテクチャパターン

### バックエンド

- **クリーンアーキテクチャ**: ドメイン → ユースケース → インフラストラクチャ
- **リポジトリパターン**: データアクセス抽象化 (`IUserRepository`)
- **依存性注入**: コンストラクタパラメータ経由
- **ミドルウェア**: 認証、ログ、エラーハンドリング
- **ドメインモデル**: ビジネスロジックを含むリッチドメインオブジェクト

### フロントエンド

- **コンポーネント合成**: 小さく、焦点を絞ったコンポーネント
- **カスタムフック**: ビジネスロジックの抽出
- **Context API**: グローバル状態管理
- **サーバーコンポーネント**: データフェッチングのデフォルト
- **クライアントコンポーネント**: インタラクティブUI要素

## ファイル整理ルール

1. **機能でグループ化**: 関連ファイルをまとめる
2. **関心の分離**: レイヤー間の明確な分離
3. **デフォルトでプライベート**: Goでは`internal/`使用、不要なエクスポート避ける
4. **一貫した命名**: 言語規約に一貫して従う
5. **ドキュメント**: 複雑なモジュールにはREADMEファイル含める

## ドキュメントファイル命名規約

### プロジェクトドキュメント (`docs/`)

#### 要件・仕様書 (`docs/requirements/`)

- **形式**: `{連番3桁}_{内容}.md`
- **例**: `001_project_overview.md`, `002_features.md`, `012_billing_and_pricing.md`
- **連番ルール**:
  - 001-099: 基本要件・概要
  - 100-199: 機能仕様
  - 200-299: 技術仕様
  - 900-999: 付録・参考資料

#### 設計ドキュメント (`docs/design/`)

- **形式**: `{連番3桁}_{設計内容}.md`
- **例**: `001_system_design.md`, `002_database_design.md`
- **連番ルール**:
  - 001-099: システム設計・アーキテクチャ
  - 100-199: データベース・API設計
  - 200-299: UI/UX設計
  - 300-399: セキュリティ・運用設計

### Steeringファイル (`.kiro/steering/`)

- **形式**: `{機能名}.md` (小文字、ハイフン区切り)
- **例**: `tech.md`, `structure.md`, `product.md`, `workflow.md`
- **命名原則**:
  - 技術関連: `tech.md`, `architecture.md`
  - プロジェクト構造: `structure.md`, `naming.md`
  - プロダクト: `product.md`, `business.md`
  - 開発プロセス: `workflow.md`, `testing.md`
  - コード品質: `quality.md`, `review.md`

### その他ドキュメント

- **README**: 各ディレクトリのルートに `README.md`
- **API仕様**: `api-spec.md` または OpenAPI形式 `openapi.yaml`
- **変更履歴**: `CHANGELOG.md`
- **貢献ガイド**: `CONTRIBUTING.md`

## インポート規約

### Go

```go
// 標準ライブラリ最初
import (
    "context"
    "fmt"

    // サードパーティパッケージ
    "github.com/labstack/echo"

    // ローカルパッケージ
    "github.com/o-ga09/zenn-hackthon-2026/internal/domain"
)
```

### TypeScript

```typescript
// React/Next.js最初
import React from "react";
import { NextPage } from "next";

// サードパーティライブラリ
import { useQuery } from "@tanstack/react-query";

// ローカルインポート (絶対パス推奨)
import { Button } from "@/components/ui/button";
import { getUserData } from "@/api/user";
```
