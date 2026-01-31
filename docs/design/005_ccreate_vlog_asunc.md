# VLog作成機能の非同期化設計

## 背景

現在のVLog作成APIは同期処理で最大120秒待機させる実装です。即座にジョブIDを返し、バックグラウンドで処理を行う非同期アーキテクチャに変更します。

## 結論

**Goroutine単独はCloud Run環境では推奨しない。Cloud Tasksを使用した非同期構成が最適。**

Cloud RunはHTTPレスポンス返却後、リクエストベース課金ではCPU割り当てがゼロになりGoroutineが停止します。インスタンスベース課金でも、スケールイン時に強制終了されるリスクがあります。Cloud Tasksは月100万回まで無料で、自動リトライ機能も組み込まれており、本ユースケースに最適です。

## スケーラビリティ比較

| 構成 | スケーラビリティ | リスク |
|------|-----------------|--------|
| **Goroutine単独** | ❌ 低い | レスポンス後にCPU停止、スケールイン時に強制終了 |
| **Cloud Tasks** | ✅ 高い | 自動リトライ、レート制限可能、キューイング組み込み |
| **Cloud Run Jobs** | ✅ 高い | 実装コストが高い（別コンテナ必要） |

## コスト比較

| 月間VLog数 | Cloud Tasks | Cloud Run（処理時間） | 合計 |
|------------|------------|----------------------|------|
| 100件 | **$0.00**（無料枠内） | $0.33 | **$0.33** |
| 1,000件 | **$0.00**（無料枠内） | $3.33 | **$3.33** |
| 10,000件 | **$0.00**（無料枠内） | $33.33 | **$33.33** |

※ Cloud Tasksは月100万回まで無料。VLog生成コストはCloud Runの実行時間のみ。

## アーキテクチャ

```
┌─────────────────────────────────────────────────────────────────┐
│                        ユーザーリクエスト                          │
└─────────────────────────────────────────────────────────────────┘
                                ↓
┌─────────────────────────────────────────────────────────────────┐
│  Cloud Run (API)                                                 │
│  POST /api/agent/create-vlog                                     │
│   1. メディアアップロード (R2)                                     │
│   2. VLogレコード作成 (status=pending)                            │
│   3. Cloud Tasksにタスク登録                                       │
│   4. 即座にレスポンス返却 (vlogId, status: "processing")          │
└─────────────────────────────────────────────────────────────────┘
                                ↓
┌─────────────────────────────────────────────────────────────────┐
│  Cloud Tasks Queue                                               │
│  - 自動リトライ (最大3回)                                          │
│  - レート制限 (10件/秒)                                            │
└─────────────────────────────────────────────────────────────────┘
                                ↓
┌─────────────────────────────────────────────────────────────────┐
│  Cloud Run (API)                                                 │
│  POST /internal/tasks/create-vlog (内部エンドポイント)             │
│   1. VLogステータス更新 (processing)                               │
│   2. Gemini Vision でメディア分析                                  │
│   3. Veo3 で動画生成                                              │
│   4. GCS → R2 転送                                                │
│   5. VLogステータス更新 (completed/failed)                         │
└─────────────────────────────────────────────────────────────────┘
                                ↓
┌─────────────────────────────────────────────────────────────────┐
│  フロントエンド                                                   │
│  - ポーリングで進捗確認 GET /api/vlogs/:id                        │
│  - または SSE で完了通知                                          │
└─────────────────────────────────────────────────────────────────┘
```

## 実装ステップ

### Step 1: VLogテーブルにステータス管理を追加

**ファイル**: `backend/internal/domain/vlog.go`, DBマイグレーション

```go
type VlogStatus string

const (
    VlogStatusPending    VlogStatus = "pending"
    VlogStatusProcessing VlogStatus = "processing"
    VlogStatusCompleted  VlogStatus = "completed"
    VlogStatusFailed     VlogStatus = "failed"
)

type Vlog struct {
    BaseModel
    VideoID      string     `gorm:"column:video_id"`
    VideoURL     string     `gorm:"column:video_url"`
    ShareURL     string     `gorm:"column:share_url"`
    Duration     float64    `gorm:"column:duration"`
    Thumbnail    string     `gorm:"column:thumbnail"`
    // 追加フィールド
    Status       VlogStatus `gorm:"column:status;default:pending"`
    ErrorMessage string     `gorm:"column:error_message"`
    Progress     float64    `gorm:"column:progress;default:0"`
    StartedAt    *time.Time `gorm:"column:started_at"`
    CompletedAt  *time.Time `gorm:"column:completed_at"`
}
```

### Step 2: Cloud Tasks クライアント実装

**ファイル**: `backend/pkg/cloudtasks/client.go`

- Cloud Tasks クライアントを追加
- タスク登録機能を実装
- OIDC認証設定

### Step 3: CreateVLog APIを非同期化

**ファイル**: `backend/internal/handler/agent.go`

1. メディアアップロード
2. VLogレコード作成（status=pending）
3. Cloud Tasksにタスク登録
4. 即座にVLog IDを返却

### Step 4: 内部タスクエンドポイント実装

**ファイル**: `backend/internal/server/router.go`, `backend/internal/handler/agent.go`

- `POST /internal/tasks/create-vlog` を追加
- Cloud Tasksからの呼び出し用
- 認証はCloud Tasks署名検証

### Step 5: 進捗取得エンドポイント更新

**ファイル**: `backend/internal/handler/vlog.go`

- `GET /api/vlogs/:id` でステータスと進捗を含むレスポンスを返す

### Step 6: フロントエンドのポーリング実装

**ファイル**: `frontend/api/createVlog.ts`, `frontend/app/upload/page.tsx`

- VLog ID取得後、完了までポーリングで監視
- 進捗バーの表示

## 実装優先度

| 優先度 | 項目 | 工数（目安） |
|--------|------|-------------|
| P0 | DBにVLogステータスカラム追加 | 0.5日 |
| P0 | Cloud Tasks タスク登録処理 | 0.5日 |
| P0 | 内部エンドポイント実装 | 1日 |
| P1 | フロントエンドポーリング実装 | 0.5日 |
| P2 | エラーハンドリング強化 | 0.5日 |
| P3 | SSE通知（オプション） | 1日 |

**合計工数**: 約3〜4日

## 検討事項

### Cloud Tasks のキュー設定

- リトライ回数: 3回
- リトライ間隔: 指数バックオフ（最小10秒、最大600秒）
- レート制限: 10件/秒

### 内部エンドポイントの認証

Cloud Tasksからの呼び出しは以下の方法で保護:
- **推奨**: OIDC認証（サービスアカウントトークン検証）
- 代替: X-CloudTasks-* ヘッダー検証

### SSE対応の優先度

- **初期リリース**: ポーリングのみ（3秒間隔）
- **後続フェーズ**: SSEでリアルタイム通知
