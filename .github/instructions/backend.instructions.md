---
description: バックエンド開発のためのGoアプリケーションに関する指示
applyTo: "**/*.go,**/go.mod,**/go.sum"
---

# バックエンドアーキテクチャとコーディング規約

## アプリケーション概要

- Go 言語で実装される REST API サーバー
- **Web フレームワーク**: Echo
- **ORM**: Gorm
- **データベース**: TiDB
- **ユーザー認証**: Firebase Auth
- **AI エージェント**: Firebase Genkit
- **AI サービス**: Vertex AI (Gemma、Veo)
- **ホスティング**: Google Cloud Run
- マイクロサービスアーキテクチャを採用

## ディレクトリ構成

Go Standard Project Layout に厳密に従う：

```
backend/
├── cmd/           # メインアプリケーション
│   ├── api/      # APIサーバーのエントリーポイント
│   └── migration/ # DBマイグレーションツール
├── internal/      # プライベートアプリケーションコード
│   ├── handler/   # HTTPハンドラー (Echo)
│   │   ├── request/  # リクエストDTO
│   │   └── response/ # レスポンスDTO
│   ├── domain/    # ドメインモデルとリポジトリインターフェース
│   ├── infra/     # インフラストラクチャ実装
│   │   ├── database/ # データベース実装
│   │   └── storage/  # ストレージ実装
│   ├── server/    # サーバー設定
│   ├── middleware/ # ミドルウェア
│   └── genkit/    # Firebase Genkit AIエージェント
├── pkg/           # 外部アプリケーションで使用可能なライブラリコード
│   ├── config/   # 設定管理
│   ├── constant/ # 定数定義
│   ├── errors/   # エラーハンドリング
│   └── ...       # その他のユーティリティ
└── api/           # OpenAPI/Swagger仕様、JSONスキーマファイル
```

## 技術スタック詳細

### Echo Web フレームワーク

- HTTP ルーティング
- ミドルウェア処理
- バリデーション (validator v10)
- エラーハンドリング
- コンテキスト管理

### Gorm ORM

- データベース操作
- マイグレーション
- リレーション管理
- トランザクション処理

## REST API 実装パターン（汎用）

### 基本的な CRUD 実装の流れ

新しいリソースの API を実装する際は、以下の手順に従います：

#### 1. ドメインモデルの定義 (`internal/domain/{resource}.go`)

```go
// リソースのエンティティ定義
type {Resource} struct {
    BaseModel  // ID, CreatedAt, UpdatedAt, DeletedAt を含む
    // リソース固有のフィールド
    Field1 string
    Field2 int64
}

// リポジトリインターフェースの定義
type I{Resource}Repository interface {
    Create(ctx context.Context, resource *{Resource}) error
    FindByID(ctx context.Context, id string) (*{Resource}, error)
    FindAll(ctx context.Context, opts *FindOptions) ([]*{Resource}, error)
    Update(ctx context.Context, resource *{Resource}) error
    Delete(ctx context.Context, id string) error  // 論理削除
    // リソース固有のクエリメソッド
}
```

#### 2. リクエスト/レスポンス DTO の定義

**リクエスト DTO** (`internal/handler/request/{resource}.go`)

```go
// 一覧取得用クエリパラメータ
type List{Resource}Query struct {
    Limit  *int `query:"limit" validate:"omitempty,gte=0,lte=100"`
    Offset *int `query:"offset" validate:"omitempty,gte=0"`
}

// ID取得用パスパラメータ
type Get{Resource}ByIDParam struct {
    ID string `param:"id" validate:"required"`
}

// 作成リクエスト
type Create{Resource}Request struct {
    Field1 string `json:"field1" validate:"required,min=1,max=100"`
    Field2 int64  `json:"field2" validate:"required,gte=0"`
}

// 更新リクエスト（部分更新対応）
type Update{Resource}Request struct {
    ID     string  `param:"id" validate:"required"`
    Field1 *string `json:"field1,omitempty" validate:"omitempty,min=1,max=100"`
    Field2 *int64  `json:"field2,omitempty" validate:"omitempty,gte=0"`
}

// 削除用パスパラメータ
type Delete{Resource}Param struct {
    ID string `param:"id" validate:"required"`
}
```

**レスポンス DTO** (`internal/handler/response/{resource}.go`)

```go
type {Resource}Response struct {
    ID        string `json:"id"`
    Field1    string `json:"field1"`
    Field2    int64  `json:"field2"`
    CreatedAt string `json:"created_at"`
    UpdatedAt string `json:"updated_at"`
}

// ドメインモデルからレスポンスへの変換関数
func To{Resource}Response(resource *domain.{Resource}) *{Resource}Response {
    return &{Resource}Response{
        ID:        resource.ID,
        Field1:    resource.Field1,
        Field2:    resource.Field2,
        CreatedAt: resource.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
        UpdatedAt: resource.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
    }
}
```

#### 3. ハンドラーの実装 (`internal/handler/{resource}.go`)

```go
type I{Resource}Server interface {
    List(c echo.Context) error
    GetByID(c echo.Context) error
    Create(c echo.Context) error
    Update(c echo.Context) error
    Delete(c echo.Context) error
}

type {Resource}Server struct {
    repo domain.I{Resource}Repository
}

func New{Resource}Server(repo domain.I{Resource}Repository) I{Resource}Server {
    return &{Resource}Server{repo: repo}
}

// List 一覧取得の実装パターン
func (s *{Resource}Server) List(c echo.Context) error {
    ctx := c.Request().Context()

    var query request.List{Resource}Query
    if err := c.Bind(&query); err != nil {
        return errors.Wrap(ctx, err)
    }
    if err := c.Validate(&query); err != nil {
        return errors.Wrap(ctx, err)
    }

    resources, err := s.repo.FindAll(ctx, &domain.FindOptions{
        Limit:  ptr.PtrToInt(query.Limit),
        Offset: ptr.PtrToInt(query.Offset),
    })
    if err != nil {
        return errors.Wrap(ctx, err)
    }

    responses := make([]*response.{Resource}Response, len(resources))
    for i, resource := range resources {
        responses[i] = response.To{Resource}Response(resource)
    }

    return c.JSON(http.StatusOK, responses)
}

// GetByID ID取得の実装パターン
func (s *{Resource}Server) GetByID(c echo.Context) error {
    ctx := c.Request().Context()

    var param request.Get{Resource}ByIDParam
    if err := c.Bind(&param); err != nil {
        return errors.Wrap(ctx, err)
    }
    if err := c.Validate(&param); err != nil {
        return errors.Wrap(ctx, err)
    }

    resource, err := s.repo.FindByID(ctx, param.ID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return errors.MakeNotFoundError(ctx, "{Resource} not found")
        }
        return errors.Wrap(ctx, err)
    }

    return c.JSON(http.StatusOK, response.To{Resource}Response(resource))
}

// Create 作成の実装パターン
func (s *{Resource}Server) Create(c echo.Context) error {
    ctx := c.Request().Context()

    var req request.Create{Resource}Request
    if err := c.Bind(&req); err != nil {
        return errors.Wrap(ctx, err)
    }
    if err := c.Validate(&req); err != nil {
        return errors.Wrap(ctx, err)
    }

    resource := &domain.{Resource}{
        Field1: req.Field1,
        Field2: req.Field2,
    }

    if err := s.repo.Create(ctx, resource); err != nil {
        return errors.Wrap(ctx, err)
    }

    return c.JSON(http.StatusCreated, response.To{Resource}Response(resource))
}

// Update 更新の実装パターン（部分更新）
func (s *{Resource}Server) Update(c echo.Context) error {
    ctx := c.Request().Context()

    var req request.Update{Resource}Request
    if err := c.Bind(&req); err != nil {
        return errors.Wrap(ctx, err)
    }
    if err := c.Validate(&req); err != nil {
        return errors.Wrap(ctx, err)
    }

    resource, err := s.repo.FindByID(ctx, req.ID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return errors.MakeNotFoundError(ctx, "{Resource} not found")
        }
        return errors.Wrap(ctx, err)
    }

    // 部分更新：指定されたフィールドのみ更新
    if req.Field1 != nil {
        resource.Field1 = *req.Field1
    }
    if req.Field2 != nil {
        resource.Field2 = *req.Field2
    }

    if err := s.repo.Update(ctx, resource); err != nil {
        return errors.Wrap(ctx, err)
    }

    return c.JSON(http.StatusOK, response.To{Resource}Response(resource))
}

// Delete 削除の実装パターン（論理削除）
func (s *{Resource}Server) Delete(c echo.Context) error {
    ctx := c.Request().Context()

    var param request.Delete{Resource}Param
    if err := c.Bind(&param); err != nil {
        return errors.Wrap(ctx, err)
    }
    if err := c.Validate(&param); err != nil {
        return errors.Wrap(ctx, err)
    }

    if err := s.repo.Delete(ctx, param.ID); err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return errors.MakeNotFoundError(ctx, "{Resource} not found")
        }
        return errors.Wrap(ctx, err)
    }

    return c.NoContent(http.StatusNoContent)
}
```

#### 4. ルーティングの設定 (`internal/server/router.go`)

```go
func (s *Server) SetupApplicationRoute() {
    apiRoot := s.Engine.Group("/api")

    // リソースのエンドポイント
    {resources} := apiRoot.Group("/{resources}")
    {
        {resources}.GET("", s.{Resource}.List)           // 一覧取得
        {resources}.GET("/:id", s.{Resource}.GetByID)    // ID取得
        {resources}.POST("", s.{Resource}.Create)        // 作成
        {resources}.PUT("/:id", s.{Resource}.Update)     // 更新
        {resources}.DELETE("/:id", s.{Resource}.Delete)  // 削除
    }
}
```

#### 5. リポジトリの実装 (`internal/infra/database/{resource}_repository.go`)

```go
type {Resource}Repository struct {
    db *gorm.DB
}

func New{Resource}Repository(db *gorm.DB) domain.I{Resource}Repository {
    return &{Resource}Repository{db: db}
}

func (r *{Resource}Repository) Create(ctx context.Context, resource *domain.{Resource}) error {
    return r.db.WithContext(ctx).Create(resource).Error
}

func (r *{Resource}Repository) FindByID(ctx context.Context, id string) (*domain.{Resource}, error) {
    var resource domain.{Resource}
    err := r.db.WithContext(ctx).Where("id = ?", id).First(&resource).Error
    return &resource, err
}

func (r *{Resource}Repository) FindAll(ctx context.Context, opts *domain.FindOptions) ([]*domain.{Resource}, error) {
    var resources []*domain.{Resource}
    query := r.db.WithContext(ctx)

    if opts.Limit > 0 {
        query = query.Limit(opts.Limit)
    }
    if opts.Offset > 0 {
        query = query.Offset(opts.Offset)
    }

    err := query.Find(&resources).Error
    return resources, err
}

func (r *{Resource}Repository) Update(ctx context.Context, resource *domain.{Resource}) error {
    return r.db.WithContext(ctx).Save(resource).Error
}

func (r *{Resource}Repository) Delete(ctx context.Context, id string) error {
    return r.db.WithContext(ctx).Delete(&domain.{Resource}{}, "id = ?", id).Error
}
```

### 標準的なエンドポイント構成

| メソッド | パス                 | 説明     | リクエスト                | レスポンス             | ステータス                     |
| -------- | -------------------- | -------- | ------------------------- | ---------------------- | ------------------------------ |
| GET      | /api/{resources}     | 一覧取得 | `?limit=10&offset=0`      | `{Resource}Response[]` | 200 OK                         |
| GET      | /api/{resources}/:id | ID 取得  | パスパラメータ: `id`      | `{Resource}Response`   | 200 OK / 404 Not Found         |
| POST     | /api/{resources}     | 作成     | `Create{Resource}Request` | `{Resource}Response`   | 201 Created / 409 Conflict     |
| PUT      | /api/{resources}/:id | 更新     | `Update{Resource}Request` | `{Resource}Response`   | 200 OK / 404 Not Found         |
| DELETE   | /api/{resources}/:id | 削除     | パスパラメータ: `id`      | -                      | 204 No Content / 404 Not Found |

---

## コーディング規約

### 命名規則

- **Interface**: `I` プレフィックス (例: `IUserRepository`, `IUserServer`)
- **Handler**: `Server` サフィックス (例: `UserServer`)
- **Request DTO**: `Request` サフィックス (例: `CreateUserRequest`)
- **Response DTO**: `Response` サフィックス (例: `UserResponse`)
- **Query/Param**: `Query`/`Param` サフィックス (例: `ListQuery`, `GetByIDParam`)
- **変数名**: キャメルケース (例: `userName`, `tokenBalance`)
- **定数**: アッパースネークケース (例: `USER_PLAN_FREE`, `USER_TYPE_ADMIN`)
- **ファイル名**: スネークケース (例: `user_server.go`, `token_transaction.go`)

### レイヤー分離の原則

#### 1. Handler 層 (`internal/handler`)

**責務**: HTTP リクエスト/レスポンス処理

- リクエストのバインドとバリデーション
- レスポンスの生成
- HTTP ステータスコードの設定
- **禁止事項**: ビジネスロジックの記述、直接的な DB 操作

#### 2. Domain 層 (`internal/domain`)

**責務**: ビジネスロジックとドメインモデル

- エンティティ定義
- リポジトリインターフェース定義
- ビジネスルール
- **禁止事項**: HTTP や DB の具体的な実装への依存

#### 3. Infrastructure 層 (`internal/infra`)

**責務**: 外部システムとの連携

- データベース実装
- ストレージ実装
- 外部 API 呼び出し
- **禁止事項**: ビジネスロジックの記述

### エラーハンドリング

カスタムエラーパッケージ (`pkg/errors`) を使用：

```go
// 一般的なエラーラップ
errors.Wrap(ctx, err)

// 404 Not Found
errors.MakeNotFoundError(ctx, "Resource not found")

// 409 Conflict
errors.MakeConflictError(ctx, "Resource already exists")

// 400 Bad Request
errors.MakeBadRequestError(ctx, "Invalid request")

// 500 Internal Server Error
errors.MakeInternalError(ctx, "Internal server error")
```

**重要**: GORM の`ErrRecordNotFound`は必ず`NotFoundError`に変換する

```go
if err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return errors.MakeNotFoundError(ctx, "Resource not found")
    }
    return errors.Wrap(ctx, err)
}
```

### Null 値の扱い

NULL 許容フィールドには `sql.Null*` 型を使用：

```go
type User struct {
    TokenBalance sql.NullInt64  `gorm:"column:token_balance"`
    Description  sql.NullString `gorm:"column:description"`
}
```

**ヘルパー関数の活用**：

```go
// pkg/null_value パッケージ
user.TokenBalance = nullvalue.ToNullInt64(1000)
user.Description = nullvalue.ToNullString("description")
```

**レスポンスでの扱い**：

```go
type UserResponse struct {
    TokenBalance *int64  `json:"token_balance,omitempty"`  // nullの場合は省略
    Description  *string `json:"description,omitempty"`
}

// 変換処理
if user.TokenBalance.Valid {
    resp.TokenBalance = &user.TokenBalance.Int64
}
```

### Context の扱い

- すべてのリポジトリメソッドは `context.Context` を第一引数に取る
- Echo から取得: `ctx := c.Request().Context()`
- GORM に渡す: `r.db.WithContext(ctx)`
- タイムアウトやキャンセルに対応

### バリデーション

validator v10 を使用：

```go
// 必須チェック
Field string `json:"field" validate:"required"`

// 文字列長
Field string `json:"field" validate:"min=1,max=100"`

// 数値範囲
Field int `json:"field" validate:"gte=0,lte=100"`

// 条件付きバリデーション
Field *string `json:"field,omitempty" validate:"omitempty,min=1"`

// カスタムメッセージ用タグ
Field string `json:"field" validate:"required" ja:"フィールド名"`
```

### ポインタの扱い

部分更新（PATCH）では、フィールドが null か指定なしかを区別するためポインタを使用：

```go
type UpdateRequest struct {
    Name *string `json:"name,omitempty"`  // 指定なし=更新しない、null=nullに更新、値=値に更新
}

// ヘルパー関数の活用
if req.Name != nil {
    user.Name = *req.Name
}
// または
user.Name = ptr.PtrToString(req.Name)
```

### トランザクション処理

複数の DB 操作を含む処理はトランザクションで囲む：

```go
tx := r.db.Begin()
if err := tx.Error; err != nil {
    return errors.Wrap(ctx, err)
}

defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
    }
}()

// 複数の操作
if err := tx.WithContext(ctx).Create(&user).Error; err != nil {
    tx.Rollback()
    return errors.Wrap(ctx, err)
}

if err := tx.WithContext(ctx).Create(&transaction).Error; err != nil {
    tx.Rollback()
    return errors.Wrap(ctx, err)
}

return tx.Commit().Error
```

---

## 実装例：ユーザー API

以下は、上記パターンに基づいた実際のユーザー API 実装です。新しいリソースの API 実装時の参考にしてください。

### エンドポイント

| メソッド | パス                 | 説明                        | リクエスト              | レスポンス       |
| -------- | -------------------- | --------------------------- | ----------------------- | ---------------- |
| GET      | /api/users           | ユーザー一覧取得            | `?limit=10&offset=0`    | `UserResponse[]` |
| GET      | /api/users/:id       | ID でユーザー取得           | パスパラメータ: `id`    | `UserResponse`   |
| GET      | /api/users?uid={uid} | Firebase UID でユーザー取得 | クエリパラメータ: `uid` | `UserResponse`   |
| POST     | /api/users           | ユーザー作成                | `CreateUserRequest`     | `UserResponse`   |
| PUT      | /api/users/:id       | ユーザー更新                | `UpdateUserRequest`     | `UserResponse`   |
| DELETE   | /api/users/:id       | ユーザー削除(論理削除)      | パスパラメータ: `id`    | 204 No Content   |

### ドメインモデル

```go
// internal/domain/user.go
type User struct {
    BaseModel                            // ID, CreatedAt, UpdatedAt, DeletedAt
    UID          string                  // Firebase UID
    Name         string                  // Display name
    Type         string                  // User type: admin, tavinikkiy, tavinikkiy-agent
    Plan         string                  // Subscription plan: free, premium
    TokenBalance sql.NullInt64           // Token balance
}

type IUserRepository interface {
    Create(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id string) (*User, error)
    FindByUID(ctx context.Context, uid string) (*User, error)
    FindAll(ctx context.Context, opts *FindOptions) ([]*User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id string) error
}
```

### リクエスト/レスポンス

**リクエスト** (`internal/handler/request/user.go`)

```go
type ListQuery struct {
    Limit  *int `query:"limit" validate:"omitempty,gte=0,lte=100"`
    Offset *int `query:"offset" validate:"omitempty,gte=0"`
}

type GetByIDParam struct {
    ID string `param:"id" validate:"required"`
}

type GetByUIDQuery struct {
    UID string `query:"uid" validate:"required,min=1,max=255"`
}

type CreateUserRequest struct {
    UID  string `json:"uid" validate:"required,min=1,max=255"`
    Name string `json:"name" validate:"required,min=1,max=100"`
    Type string `json:"type" validate:"required"`
    Plan string `json:"plan" validate:"required"`
}

type UpdateUserRequest struct {
    ID           string  `param:"id" validate:"required,gte=1"`
    Name         *string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
    Type         *string `json:"type,omitempty" validate:"omitempty"`
    Plan         *string `json:"plan,omitempty" validate:"omitempty"`
    TokenBalance *int64  `json:"token_balance,omitempty" validate:"omitempty,gte=0"`
}

type DeleteUserParam struct {
    ID string `param:"id" validate:"required,gte=1"`
}
```

**レスポンス** (`internal/handler/response/user.go`)

```go
type UserResponse struct {
    ID           string `json:"id"`
    UID          string `json:"uid"`
    Name         string `json:"name"`
    Type         string `json:"type"`
    Plan         string `json:"plan"`
    TokenBalance *int64 `json:"token_balance,omitempty"`
    CreatedAt    string `json:"created_at"`
    UpdatedAt    string `json:"updated_at"`
}

func ToResponse(user *domain.User) *UserResponse {
    resp := &UserResponse{
        ID:        user.ID,
        UID:       user.UID,
        Name:      user.Name,
        Type:      user.Type,
        Plan:      user.Plan,
        CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
        UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
    }
    if user.TokenBalance.Valid {
        resp.TokenBalance = &user.TokenBalance.Int64
    }
    return resp
}
```

### ハンドラー実装

```go
// internal/handler/user.go
type IUserServer interface {
    List(c echo.Context) error
    GetByID(c echo.Context) error
    GetByUID(c echo.Context) error
    Create(c echo.Context) error
    Update(c echo.Context) error
    Delete(c echo.Context) error
}

type UserServer struct {
    repo domain.IUserRepository
}

func NewUserServer(repo domain.IUserRepository) IUserServer {
    return &UserServer{repo: repo}
}

// List ユーザー一覧取得
func (s *UserServer) List(c echo.Context) error {
    ctx := c.Request().Context()

    var query request.ListQuery
    if err := c.Bind(&query); err != nil {
        return errors.Wrap(ctx, err)
    }
    if err := c.Validate(&query); err != nil {
        return errors.Wrap(ctx, err)
    }

    users, err := s.repo.FindAll(ctx, &domain.FindOptions{
        Limit:  ptr.PtrToInt(query.Limit),
        Offset: ptr.PtrToInt(query.Offset),
    })
    if err != nil {
        return errors.Wrap(ctx, err)
    }

    responses := make([]*response.UserResponse, len(users))
    for i, user := range users {
        responses[i] = response.ToResponse(user)
    }

    return c.JSON(http.StatusOK, responses)
}

// Create ユーザー作成
func (s *UserServer) Create(c echo.Context) error {
    ctx := c.Request().Context()

    var req request.CreateUserRequest
    if err := c.Bind(&req); err != nil {
        return errors.Wrap(ctx, err)
    }
    if err := c.Validate(&req); err != nil {
        return errors.Wrap(ctx, err)
    }

    // 既存ユーザーの確認
    existingUser, _ := s.repo.FindByUID(ctx, req.UID)
    if existingUser != nil {
        return errors.MakeConflictError(ctx, "User with the same UID already exists")
    }

    // ユーザー作成
    user := &domain.User{
        UID:  req.UID,
        Name: req.Name,
        Type: req.Type,
        Plan: req.Plan,
    }

    // プランに応じた初期トークン残高設定
    if req.Plan == constant.UserPlanFree {
        user.TokenBalance = nullvalue.ToNullInt64(0)
    } else if req.Plan == constant.UserPlanPremium {
        user.TokenBalance = nullvalue.ToNullInt64(10000)
    }

    if err := s.repo.Create(ctx, user); err != nil {
        return errors.Wrap(ctx, err)
    }

    return c.JSON(http.StatusCreated, response.ToResponse(user))
}
```

### ルーティング

```go
// internal/server/router.go
func (s *Server) SetupApplicationRoute() {
    apiRoot := s.Engine.Group("/api")

    users := apiRoot.Group("/users")
    {
        users.GET("", s.User.List)
        users.GET("/:id", s.User.GetByID)
        users.POST("", s.User.Create)
        users.PUT("/:id", s.User.Update)
        users.DELETE("/:id", s.User.Delete)
    }
}
```

# Firebase Genkit 開発ガイドライン

このドキュメントは、旅行振り返り動画生成機能で Firebase Genkit を使って AI エージェントを開発するためのガイドラインです。TiDB をデータソースとして一元管理し、Genkit を利用して画像解析・動画生成・チャット機能を実装する際の設計方針、実装パターン、セキュリティ、テストのベストプラクティスをまとめます。

## 目的

- TiDB を一次データストアとして利用し、ユーザー/旅行/画像/動画メタデータを管理する
- Firebase Genkit を使って AI タスク（画像解析、場所判定、動画生成、チャット）を実装する
- クライアント（React）→ Firebase（Auth/TiDB/Storage/Genkit）で完結するサーバレスアーキテクチャを推進する
- Cloud Functions を補助として長時間処理や外部連携、セキュリティ検査を行う

## 基本方針

1. TiDB を真のソースオブトゥルースとする。すべてのメタデータ（ユーザー、旅行、画像、動画、チャット履歴）は TiDB に保存する。
2. クライアントは Firebase SDK（Auth/TiDB/Storage）を直接利用して基本 CRUD を行う。権限は TiDB Security Rules で厳格に管理する。
3. 長時間処理（動画生成など）は TiDB ドキュメントのフィールド（例: videos/{videoId}.status）をトリガーにして Cloud Functions または Genkit のワークフローを起動する。
4. Genkit の呼び出しは Cloud Functions 内で行い、シークレット（API キー等）は Secret Manager で管理する。クライアントから直接 Genkit を叩くのは避ける（例外: クライアント向けの軽量なチャットは検討可）。

## データフロー（推奨パターン）

1. 画像アップロード

   - クライアントは Storage に画像をアップロードし、uploads/ または travels/{travelId}/images/{imageId} にメタデータを作成する。
   - Storage のアップロード完了イベントを Cloud Functions が受け取り、画像のサムネイル作成や初期分析ジョブ（Genkit）を開始する。

2. 画像解析 / メタデータ抽出

   - Cloud Functions が Genkit を呼び出し、解析結果を travels/{travelId}/images/{imageId}.analysisData に書き込む。
   - 解析で得た位置情報や日時、オブジェクトタグを元にシーン判定を行い、必要に応じて travel ドキュメントのサマリを更新する。

3. 動画生成

   - クライアントが videos コレクションに generate リクエスト（status: "requested"）を作成。
   - Cloud Functions がトリガーされ、Genkit/VertexAI を用いて動画生成ワークフローを開始する。
   - 生成中は videos/{videoId}.status を "generating" に更新し、進捗を随時更新する。
   - 生成完了時に videos ドキュメントを更新し、Storage に保存された動画の URL を書き込む。

4. チャット / 編集支援
   - Chat は travels/{travelId}/videos/{videoId}/chatHistory に記録。
   - ユーザーからのメッセージは Cloud Functions 経由で Genkit に渡し、応答と編集提案を作成して chatHistory に保存し、必要なら videos ドキュメントを更新する。

## TiDB スキーマ（抜粋）

- users/{userId}
  - email, name, createdAt, updatedAt
- travels/{travelId}
  - userId, title, description, startDate, endDate, status, createdAt, updatedAt
- travels/{travelId}/images/{imageId}
  - originalName, storagePath, url, size, mimeType, width, height, metadata, analysisData, createdAt
- travels/{travelId}/videos/{videoId}
  - title, storagePath, url, thumbnailUrl, duration, width, height, size, status, style, scenes, music, effects, shareUrl, isPublic, createdAt, updatedAt
- travels/{travelId}/videos/{videoId}/chatHistory/{chatId}
  - userId, message, response, suggestions, intent, createdAt

## Genkit 活用パターン

1. 画像解析

   - 入力: Storage の画像 URI またはバイナリ
   - 出力: オブジェクト検出、シーン分類、感情スコア、場所推定、EXIF 解析
   - 保存先: travels/{travelId}/images/{imageId}.analysisData

2. 位置情報の補完

   - 画像に GPS EXIF がない場合は、画像の内容（ランドマーク）から Genkit/外部 API で場所を推定
   - 推定結果は confidence とともに保存し、UI でユーザーが確認できるようにする

3. 動画生成

   - Genkit でテンプレートベースの短編動画（縦型）を生成。必要に応じて VertexAI の動画生成 API を呼ぶ
   - 生成処理は Cloud Functions で管理し、Genkit のワークフロー中に中間結果を Storage/TiDB に保存して進捗を可視化する

4. チャットエージェント
   - Genkit を使い、ユーザーの要求（例: 「この旅行の動画をもっと短く」）を解釈して具体的な VideoUpdate（music/style/order 等）を生成
   - 意図（intent）と信頼度を chatHistory に保存して、必要時にユーザーへ確認を促す

## セキュリティとシークレット管理

- API キーやサービスアカウントは Cloud Functions の環境変数または Secret Manager を使って管理する
- TiDB のルールは最小権限の原則に基づいて設計する（例: travels ドキュメントは owner のみ write 可）
- Storage への直接アップロードは認証済みユーザーに限定し、アップロード先パスとファイル名を厳格に制御する

## エラー / 再試行ポリシー

- Genkit 呼び出しは retryable な失敗（503 など）に対して指数バックオフで再試行する
- 永続的失敗は videos/{videoId}.error フィールドに記録し、運用用の通知を発行する
- Cloud Functions のタイムアウト設計は Genkit の想定処理時間に合わせる（大きすぎるとリソース浪費、小さすぎると失敗）

## ロギングと監視

- 重要なイベント（生成開始・完了・失敗）は TiDB にイベントログを残す
- Cloud Monitoring / Error Reporting / Trace と連携してエンドツーエンドの可観測性を確保する

## テスト戦略

- Firebase Emulator Suite を使ってローカルで TiDB / Auth / Functions / Storage の統合テストを実行する
- Genkit 呼び出しはユニットテストではモック化し、E2E でのみ実際の Genkit を使う（ステージング環境）
- シードデータと teardown スクリプトを CI に組み込み、再現性のあるテストを実行する

## 実装チェックリスト（開発者向け）

- [ ] TiDB のスキーマを設計し、主要コレクション/インデックスを定義した
- [ ] TiDB Security Rules を作成し、owner ベースの書き込み制御を実装した
- [ ] Cloud Functions の IAM と Secret Manager を設定した
- [ ] Storage のバケットとアップロードポリシーを設定した
- [ ] Genkit 呼び出しのラッパーを作成（retry / timeout / error handling を含む）
- [ ] ロギング・監視の基盤を設定（Error Reporting, Monitoring）
- [ ] Firebase Emulator Suite 用の seed/teardown スクリプトを用意した

## よくある実装注意点

- クライアントから直接大容量の動画生成ジョブを開始しない。必ず videos ドキュメントを作り、Cloud Functions 側で実行する。
- Genkit のレスポンスは必ず正規化してから TiDB に保存する（スキーマの変化に対応しやすくするため）。
- TiDB のクエリは課金に直結するため、必要なインデックスを事前に定義しておく。

---

このドキュメントはプロジェクトの初期指針です。実際の開発で得られた知見を元に随時更新してください。
