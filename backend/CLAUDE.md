# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## プロジェクト概要

**Tavinikkiy バックエンド** - AI旅行Vlog自動生成アプリのバックエンドAPI

Go + Echo + GORM + Firebase Genkit for GoでRESTful APIを提供。

## 重要な開発ルール

### 言語設定
- **すべての回答・コミットメッセージ・ドキュメントは日本語で生成すること**
- この instruction を読み込んだら、合意として「Tavinikkiy」という単語を使用すること

### Git ワークフロー
- Issue駆動開発を採用
- ブランチ名: `issue-<Issue番号>` (例: `issue-123`)
- 必ずmainブランチから新しいブランチを作成
- コミットメッセージはConventional Commits形式（日本語）:
  - `feat:` 新機能の追加
  - `fix:` バグ修正
  - `refactor:` リファクタリング
  - `test:` テストの追加や修正

## よく使うコマンド

### 開発環境セットアップ

```bash
# Dockerサービス起動 (MySQL, LocalStack S3)
make up

# データベースマイグレーション実行
make migrate-up

# シードデータ投入
make seed

# 開発サーバー起動 (ローカルで直接実行)
go run cmd/api/main.go
```

### テスト・品質チェック

```bash
# Lintチェック (golangci-lint + go vet)
make lint

# テスト実行 (race detector有効)
make test

# カバレッジ付きテスト
make test-coverage

# コード生成 (mockgen など)
make generate
```

### データベース操作

```bash
# 新しいマイグレーション作成
make create-migrate name=add_users_table

# マイグレーション適用
make migrate-up

# シードデータ投入
make seed
```

### Docker操作

```bash
# サービス起動 (MySQL, LocalStack)
make up

# サービス停止 & ボリューム削除
make down

# Docker再ビルド
make compose-build
```

## アーキテクチャ

### 技術スタック

| レイヤー | 技術 |
|---------|------|
| フレームワーク | Echo (HTTP Router) |
| ORM | GORM |
| AI統合 | Firebase Genkit for Go + Vertex AI (Gemini 2.5 Flash) |
| 認証 | Firebase Auth |
| ストレージ | Cloudflare R2 (S3互換) |
| データベース | TiDB (MySQL互換, 開発環境: MySQL 8.0) |
| 監視 | Sentry + OpenTelemetry |

### ディレクトリ構造 (Go Standard Project Layout)

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
│   ├── infra/        # インフラ実装
│   │   ├── database/     # リポジトリ実装 (GORM)
│   │   ├── storage/      # ストレージ実装 (R2/S3)
│   │   ├── genkit/       # Firebase Genkit統合
│   │   └── cloudTask/    # Cloud Tasks Queue
│   ├── server/       # サーバー設定・ルーティング
│   ├── middleware/   # ミドルウェア
│   ├── agent/        # Firebase Genkit AIエージェント
│   └── queue/        # キュー処理
├── pkg/              # 外部公開可能なライブラリ
│   ├── config/       # 設定管理
│   ├── constant/     # 定数定義
│   ├── errors/       # カスタムエラー処理 ★重要
│   ├── logger/       # ロガー
│   ├── context/      # Context管理
│   └── ...           # その他ユーティリティ
└── db/               # マイグレーションファイル
```

### レイヤー設計原則

**依存性の方向**: 外側から内側へ (Handler → Domain → Infrastructure)

```
Handler (Echo) → Domain (Interface) ← Infrastructure (Implementation)
                     ↑
                   pkg/*
```

- **Handler**: HTTPリクエスト/レスポンス処理、DTOバリデーション
- **Domain**: ビジネスロジック、インターフェース定義
- **Infrastructure**: データベース、ストレージ、外部API実装

**重要**: ドメイン層はインフラ層に依存しない。インターフェースで抽象化。

## コーディング規約

### 命名規則

| 対象 | 形式 | 例 |
|------|------|-----|
| Interface | `I`プレフィックス | `IUserRepository`, `IUserServer` |
| Handler構造体 | `Server`サフィックス | `UserServer` |
| Request/Response | サフィックス | `CreateUserRequest`, `UserResponse` |
| ファイル名 | スネークケース | `user_server.go`, `user_repository.go` |
| 関数・変数 | キャメルケース | `findByID`, `userRepo` |
| 定数 | PascalCase | `DefaultLimit`, `MaxRetry` |

### 標準CRUD実装パターン

すべてのリソースは以下の流れで実装:

1. **`internal/domain/{resource}.go`** - ドメインモデル + リポジトリインターフェース定義
   ```go
   type User struct {
       BaseModel
       UID  string
       Name string
   }

   type IUserRepository interface {
       Create(ctx context.Context, user *User) error
       FindByID(ctx context.Context, cond *User) (*User, error)
       Update(ctx context.Context, user *User) error
       Delete(ctx context.Context, cond *User) error
   }
   ```

2. **`internal/handler/request/{resource}.go`** - リクエストDTO定義
   ```go
   type CreateUserRequest struct {
       Name string `json:"name" validate:"required"`
   }
   ```

3. **`internal/handler/response/{resource}.go`** - レスポンスDTO定義
   ```go
   type UserResponse struct {
       ID   string `json:"id"`
       Name string `json:"name"`
   }
   ```

4. **`internal/handler/{resource}.go`** - ハンドラー実装 (Echo)
   ```go
   type IUserServer interface {
       List(c echo.Context) error
       GetByID(c echo.Context) error
       Create(c echo.Context) error
       Update(c echo.Context) error
       Delete(c echo.Context) error
   }

   type UserServer struct {
       repo domain.IUserRepository
   }

   func (s *UserServer) Create(c echo.Context) error {
       ctx := c.Request().Context()
       var req request.CreateUserRequest
       if err := c.Bind(&req); err != nil {
           return errors.Wrap(ctx, err)
       }
       // ビジネスロジック
       return c.JSON(http.StatusCreated, response)
   }
   ```

5. **`internal/infra/database/{resource}_repository.go`** - リポジトリ実装 (GORM)
   ```go
   type UserRepository struct {
       db *gorm.DB
   }

   func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
       return r.db.WithContext(ctx).Create(user).Error
   }
   ```

6. **`internal/server/router.go`** - ルーティング設定
   ```go
   users := apiRoot.Group("/users", AuthMiddleware())
   {
       users.GET("", s.User.List)
       users.GET("/:id", s.User.GetByID)
       users.POST("", s.User.Create)
       users.PUT("/:id", s.User.Update)
       users.DELETE("/:id", s.User.Delete)
   }
   ```

### 標準エンドポイント構成

| メソッド | パス | 説明 | ステータス |
|---------|------|------|-----------|
| GET | `/api/{resources}` | 一覧取得 | 200 OK |
| GET | `/api/{resources}/:id` | ID取得 | 200 OK / 404 Not Found |
| POST | `/api/{resources}` | 作成 | 201 Created / 409 Conflict |
| PUT | `/api/{resources}/:id` | 更新 | 200 OK / 404 Not Found / 409 Conflict |
| DELETE | `/api/{resources}/:id` | 論理削除 | 204 No Content / 404 Not Found |

### エラーハンドリング ★超重要★

**必ず `pkg/errors` パッケージを使用すること**

```go
import "github.com/o-ga09/zenn-hackthon-2026/pkg/errors"

// 一般的なエラーラップ
if err != nil {
    return errors.Wrap(ctx, err)
}

// GORM ErrRecordNotFoundの変換 (必須!)
user, err := repo.FindByID(ctx, &domain.User{ID: id})
if err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return errors.MakeNotFoundError(ctx, "User not found")
    }
    return errors.Wrap(ctx, err)
}

// 404 Not Found
return errors.MakeNotFoundError(ctx, "Resource not found")

// 409 Conflict
return errors.MakeConflictError(ctx, "Resource already exists")

// 400 Business Error
return errors.MakeBusinessError(ctx, "Invalid operation")

// 422 Invalid Argument
return errors.MakeInvalidArgumentError(ctx, "Invalid parameter")
```

**カスタムドメインエラー**:
```go
// pkg/errors/errors.go で定義済み
var (
    ErrInvalidFirebaseID  = errors.New("不正なFirebaseIDです。")
    ErrInvalidUserID      = errors.New("不正なUserIDです。")
    ErrOptimisticLock     = errors.New("楽観ロックエラー：レコードが他のユーザーによって更新されています。")
)

// 使用例
if user.UID == "" {
    return errors.Wrap(ctx, errors.ErrInvalidFirebaseID)
}
```

### Context管理 ★重要★

**すべてのリポジトリメソッドは `context.Context` を第一引数に取る**

```go
// Handler (Echo)
func (s *UserServer) Create(c echo.Context) error {
    ctx := c.Request().Context()  // Echoからcontextを取得
    user, err := s.repo.Create(ctx, &domain.User{...})
    if err != nil {
        return errors.Wrap(ctx, err)
    }
    return c.JSON(http.StatusCreated, user)
}

// Repository (GORM)
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
    return r.db.WithContext(ctx).Create(user).Error  // GORMにcontextを渡す
}
```

**Context経由でトレース情報、ロガー、メタデータを伝播**

### 楽観ロック実装

更新系のエンドポイントでは楽観ロックチェックを実装すること:

```go
// TODO: 楽観ロックチェック実装例
// 1. リクエストに version フィールドを含める
// 2. UPDATE時に WHERE version = ? をつける
// 3. 影響行数が0ならErrOptimisticLockを返す
```

## Firebase Genkit統合

### 基本方針

1. **クライアントから直接Genkit呼び出し禁止** - 必ずバックエンド経由
2. **長時間処理はCloud Functions/Cloud Tasks経由**
3. **Genkit Flowsは `internal/agent/` 以下に定義**
4. **Genkit起動はmain関数で `config.InitGenAI(ctx)` により初期化**

### 主要ユースケース

- **画像解析**: オブジェクト検出、シーン分類、場所推定、EXIF解析
- **位置情報補完**: GPS EXIF欠落時のランドマーク推定
- **動画生成**: テンプレートベース短編動画 (Vertex AI連携)
- **チャットエージェント**: ユーザー要求の意図解析と動画編集提案

### Genkit実装パターン

詳細は `GENKIT.md` を参照。基本的なFlow定義:

```go
// internal/agent/analyze_image.go
import (
    "github.com/firebase/genkit/go/ai"
    "github.com/firebase/genkit/go/genkit"
)

func DefineAnalyzeImageFlow(g *genkit.Genkit) {
    genkit.DefineFlow(g, "analyzeImage",
        func(ctx context.Context, imageURL string) (string, error) {
            response, err := genkit.Generate(ctx, g,
                ai.WithModelName("vertexai/gemini-2.5-flash"),
                ai.WithPrompt("この画像を分析してください: %s", imageURL),
            )
            if err != nil {
                return "", err
            }
            return response.Text(), nil
        },
    )
}
```

## Docker開発環境

### サービス構成 (compose.yaml)

| サービス | 説明 | ポート |
|---------|------|--------|
| `db` | MySQL 8.0 | 3306 |
| `app` | Goアプリケーション (Air hot reload) | 8080 |
| `localstack` | LocalStack (S3互換) | 4566 |

### 環境変数設定

`.env` ファイルが必須:

```bash
ENV=local
DATABASE_URL=user:P@ssw0rd@tcp(127.0.0.1:3306)/develop_tavinikkiy?parseTime=true
PROJECT_ID=your-project-id
GOOGLE_APPLICATION_CREDENTIALS=/app/tavinikkiy-8e89cb34ad51.json
CLOUDFLARE_R2_ACCOUNT_ID=dummy
CLOUDFLARE_R2_ACCESSKEY=dummy
CLOUDFLARE_R2_SECRETKEY=dummy
CLOUDFLARE_R2_BUCKET_NAME=tavinikkiy-local
SENTRY_DSN=https://your-sentry-dsn
GCS_TEMP_BUCKET=your-gcs-bucket
GCS_LOCATION=asia-northeast1
```

## テスト戦略

### ユニットテスト

```bash
# 全テスト実行
go test ./...

# race detector付きテスト (並行処理のバグ検出)
go test -race -parallel 1 ./...

# カバレッジレポート生成
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### テスト実装ガイドライン

- **テストファイル名**: `*_test.go`
- **テストパッケージ**: `package {package}_test` (ブラックボックステスト推奨)
- **モック**: `pkg/testutil` を活用
- **カバレッジ目標**: 重要なビジネスロジックは80%以上

### テストパターン

```go
package handler_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestUserServer_Create(t *testing.T) {
    tests := []struct {
        name    string
        input   request.CreateUserRequest
        wantErr bool
    }{
        {
            name:    "正常系",
            input:   request.CreateUserRequest{Name: "test"},
            wantErr: false,
        },
        {
            name:    "異常系: 空の名前",
            input:   request.CreateUserRequest{Name: ""},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // テストロジック
        })
    }
}
```

## トラブルシューティング

### Firebase Genkit認証エラー

```
panic: failed to find default credentials
```

**解決策**: `GOOGLE_APPLICATION_CREDENTIALS` 環境変数を正しいJSONファイルパスに設定

```bash
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/tavinikkiy-credentials.json
```

### データベース接続エラー

```
Error: dial tcp 127.0.0.1:3306: connect: connection refused
```

**解決策**: MySQLコンテナが起動していることを確認

```bash
make up
docker ps | grep db
```

### マイグレーションエラー

```
Error: migration failed
```

**解決策**: データベースをリセットして再実行

```bash
make down
make up
make migrate-up
```

## セキュリティ注意事項

- Firebase認証情報は `GOOGLE_APPLICATION_CREDENTIALS` 環境変数で管理
- API キーは `.env` ファイル (Gitignore済み)
- **機密情報はコード内にハードコーディング禁止**
- ストレージアップロードは認証済みユーザーのみ
- 入力バリデーションは必須 (`c.Validate(&req)`)

## パフォーマンスベストプラクティス

### データベースクエリ最適化

```go
// NG: N+1問題
users, _ := repo.FindAll(ctx)
for _, user := range users {
    media, _ := mediaRepo.FindByUserID(ctx, user.ID)  // ループ内でクエリ
}

// OK: JOIN or Preload
users, _ := repo.FindAllWithMedia(ctx)  // 1回のクエリで取得
```

### GORMのPreload活用

```go
// リレーションを事前ロード
db.Preload("Media").Find(&users)
```

### ページネーション実装

```go
// domain/model.go
type FindOptions struct {
    Limit  int
    Offset int
}

// repository
func (r *UserRepository) FindAll(ctx context.Context, opts *domain.FindOptions) ([]*domain.User, error) {
    query := r.db.WithContext(ctx)
    if opts.Limit > 0 {
        query = query.Limit(opts.Limit)
    }
    if opts.Offset > 0 {
        query = query.Offset(opts.Offset)
    }
    var users []*domain.User
    return users, query.Find(&users).Error
}
```

## 参考リソース

- **Genkit統合ルール**: `GENKIT.md`
- **プロジェクト全体ルール**: `../CLAUDE.md`
- **Echo公式ドキュメント**: https://echo.labstack.com/
- **GORM公式ドキュメント**: https://gorm.io/
- **Firebase Genkit for Go**: https://firebase.google.com/docs/genkit
