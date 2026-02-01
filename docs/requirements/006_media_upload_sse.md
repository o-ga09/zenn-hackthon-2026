# メディア分析のSSE対応実装設計

**作成日**: 2026年2月1日  
**更新日**: 2026年2月1日  
**ステータス**: 実装完了  
**優先度**: P1  
**実装工数**: 5日

## 概要

メディア分析処理（/api/agent/analyze-media）を非同期化し、Server-Sent Events (SSE) でリアルタイムに進捗を通知する機能を実装する。各メディアを個別に管理し、並列処理とリトライロジックで信頼性を向上させる。

### 目的

- アップロードと分析の進捗をリアルタイムで表示
- 複数ファイルの並列処理で効率化（最大5並列）
- 分析失敗時の自動リトライで信頼性向上（最大3回、指数バックオフ）
- 処理状態を各Mediaレコードで管理（バッチテーブル不要）

## 実装内容

### アーキテクチャ

```
┌─────────────────────────────────────────────────────────┐
│  フロントエンド                                           │
│  ┌─────────────────────────────────────────────────┐   │
│  │ PhotoUpload.tsx                                  │   │
│  │ - ファイル選択UI                                 │   │
│  │ - 分析開始ボタン                                 │   │
│  │ - useAnalyzeMedia() 呼び出し                     │   │
│  └─────────────────────────────────────────────────┘   │
│                      │                                   │
│                      ▼                                   │
│  ┌─────────────────────────────────────────────────┐   │
│  │ UploadProgress.tsx                               │   │
│  │ - useMediaAnalysisSSE(mediaIds)                  │   │
│  │ - プログレスバー表示                              │   │
│  │ - 状態別アイコン表示                              │   │
│  │ - エラー通知UI                                    │   │
│  └─────────────────────────────────────────────────┘   │
│                      │                                   │
│                      │ SSE Connection                    │
│                      ▼                                   │
└─────────────────────────────────────────────────────────┘
                       │
                       │ GET /api/agent/analyze-media/stream?ids=id1,id2,id3
                       ▼
┌─────────────────────────────────────────────────────────┐
│  バックエンド                                            │
│  ┌─────────────────────────────────────────────────┐   │
│  │ POST /api/agent/analyze-media                    │   │
│  │ 1. 各ファイルのMediaレコード作成 (PENDING)      │   │
│  │ 2. Goroutineで処理開始                          │   │
│  │ 3. media_idsリスト即座に返却                    │   │
│  └─────────────────────────────────────────────────┘   │
│                      │                                   │
│                      ▼                                   │
│  ┌─────────────────────────────────────────────────┐   │
│  │ processMediaAnalysis (Goroutine)                 │   │
│  │ ┌───────────────────────────────────────────┐   │   │
│  │ │ Phase 1: アップロード（各メディアごと）    │   │   │
│  │ │ - Status: UPLOADING                        │   │   │
│  │ │ - R2にファイルアップロード                 │   │   │
│  │ │ - Mediaレコード更新                        │   │   │
│  │ │ - Progress: 0.0 → 0.5                      │   │   │
│  │ └───────────────────────────────────────────┘   │   │
│  │ ┌───────────────────────────────────────────┐   │   │
│  │ │ Phase 2: 分析                              │   │   │
│  │ │ - Genkit AIエージェント呼び出し           │   │   │
│  │ │ - 並列処理 (5並列)                        │   │   │
│  │ │ - Exponential Backoff (最大3回)           │   │   │
│  │ │ - Progress: 0.5 → 1.0                      │   │   │
│  │ └───────────────────────────────────────────┘   │   │
│  │ ┌───────────────────────────────────────────┐   │   │
│  │ │ Phase 3: 完了                              │   │   │
│  │ │ - Status: COMPLETED / FAILED               │   │   │
│  │ │ - 各Mediaレコードを更新                    │   │   │
│  │ └───────────────────────────────────────────┘   │   │
│  └─────────────────────────────────────────────────┘   │
│                      │                                   │
│  ┌─────────────────────────────────────────────────┐   │
│  │ GET /api/agent/analyze-media/stream?ids=...      │   │
│  │ - SSE接続確立                                    │   │
│  │ - 2秒ごとに各メディアの状態をブロードキャスト    │   │
│  │ - 全てCOMPLETED/FAILED時に接続終了              │   │
│  └─────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
```
    MaxBackoff     time.Duration // 最大バックオフ時間
    Multiplier     float64       // バックオフ乗数
}

var DefaultConfig = Config{
    MaxRetries:     3,
    InitialBackoff: 1 * time.Second,
    MaxBackoff:     10 * time.Second,
    Multiplier:     2.0,
}

func Do(ctx context.Context, cfg Config, fn func() error) error {
    var lastErr error
    backoff := cfg.InitialBackoff
    
    for i := 0; i <= cfg.MaxRetries; i++ {
        if i > 0 {
            select {
            case <-ctx.Done():
                return ctx.Err()
            case <-time.After(backoff):
            }
            
            backoff = time.Duration(float64(backoff) * cfg.Multiplier)
            if backoff > cfg.MaxBackoff {
                backoff = cfg.MaxBackoff
            }
        }
        
        if err := fn(); err != nil {
            lastErr = err
            continue
        }
        
        return nil
    }
    
    return lastErr
}
```

### 2. Mediaドメイン拡張

**変更ファイル**: `backend/internal/domain/media.go`

```go
type MediaStatus string

const (
    MediaStatusPending   MediaStatus = "pending"
    MediaStatusUploading MediaStatus = "uploading"
    MediaStatusCompleted MediaStatus = "completed"
    MediaStatusFailed    MediaStatus = "failed"
)

type Media struct {
    BaseModel
    ContentType  string      `json:"content_type"`
    Size         int64       `json:"size"`
    URL          string      `json:"url"`
    Status       MediaStatus `json:"status"`        // 追加
    Progress     float64     `json:"progress"`      // 追加 (0.0~1.0)
    ErrorMessage string      `json:"error_message"` // 追加
}
```

**マイグレーション**: `backend/db/migrations/YYYYMMDDHHMMSS_add_media_status.sql`

```sql
ALTER TABLE media
ADD COLUMN status VARCHAR(20) NOT NULL DEFAULT 'completed',
ADD COLUMN progress DECIMAL(5,4) NOT NULL DEFAULT 1.0,
ADD COLUMN error_message TEXT;

CREATE INDEX idx_media_status ON media(status);
```

### 3. アップロード非同期化

**変更ファイル**: `backend/internal/handler/media_server.go`

#### UploadMedia() の改修

```go
func (s *MediaServer) UploadMedia(c echo.Context) error {
    ctx := c.Request().Context()
    
    // ファイル取得
    file, err := c.FormFile("file")
    if err != nil {
        return err
    }
    
    // Mediaレコードを PENDING で作成
    media := &domain.Media{
        ContentType:  file.Header.Get("Content-Type"),
        Size:         file.Size,
        Status:       domain.MediaStatusPending,
        Progress:     0.0,
    }
    
    if err := s.mediaRepo.Save(ctx, media); err != nil {
        return err
    }
    
    // Cloud Tasksにアップロードタスクを登録
    payload := &UploadTaskPayload{
        MediaID:  media.ID,
        FileName: file.Filename,
    }
    
    cfg := config.GetCtxEnv(ctx)
    if cfg.Env == "local" {
        // ローカル環境: Goroutineで実行
        go s.processMediaUpload(context.Background(), payload)
    } else {
        // 本番環境: Cloud Tasksに登録
        if err := s.taskClient.Enqueue(ctx, payload); err != nil {
            return err
        }
    }
    
    // 即座にmediaIdを返却
    return c.JSON(http.StatusOK, map[string]string{
        "media_id": media.ID,
    })
}
```

#### ProcessMediaUploadTask() の実装

```go
func (s *MediaServer) ProcessMediaUploadTask(c echo.Context) error {
    ctx := c.Request().Context()
    
    var payload UploadTaskPayload
    if err := c.Bind(&payload); err != nil {
        return err
    }
    
    return s.processMediaUpload(ctx, &payload)
}

func (s *MediaServer) processMediaUpload(ctx context.Context, payload *UploadTaskPayload) error {
    // Mediaレコード取得
    media, err := s.mediaRepo.FindByID(ctx, payload.MediaID)
    if err != nil {
        return err
    }
    
    // Status: UPLOADING
    media.Status = domain.MediaStatusUploading
    media.Progress = 0.1
    s.mediaRepo.Save(ctx, media)
    
    // R2アップロード実行
    url, err := s.storageClient.UploadFile(ctx, payload.FileName, payload.FileData)
    if err != nil {
        // 失敗時
        media.Status = domain.MediaStatusFailed
        media.ErrorMessage = err.Error()
        media.Progress = 0.0
        s.mediaRepo.Save(ctx, media)
        return err
    }
    
    // 成功時
    media.URL = url
    media.Status = domain.MediaStatusCompleted
    media.Progress = 1.0
    s.mediaRepo.Save(ctx, media)
    
    return nil
}
```

### 4. SSEエンドポイント実装

**変更ファイル**: `backend/internal/handler/media_server.go`

```go
func (s *MediaServer) StreamMediaStatus(c echo.Context) error {
    ctx := c.Request().Context()
    mediaID := c.Param("id")
    
    // SSEヘッダー設定
    c.Response().Header().Set("Content-Type", "text/event-stream")
    c.Response().Header().Set("Cache-Control", "no-cache")
    c.Response().Header().Set("Connection", "keep-alive")
    
    ticker := time.NewTicker(2 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            return nil
        case <-ticker.C:
            media, err := s.mediaRepo.FindByID(ctx, mediaID)
            if err != nil {
                return err
            }
            
            // JSON送信
            data, _ := json.Marshal(media)
            fmt.Fprintf(c.Response(), "data: %s\n\n", data)
            c.Response().Flush()
            
            // 完了または失敗で終了
            if media.Status == domain.MediaStatusCompleted || 
               media.Status == domain.MediaStatusFailed {
                return nil
            }
        }
    }
}
```

### 5. 画像分析の並列処理化

**変更ファイル**: `backend/internal/agent/agent.go`

```go
func analyzeAllMedia(ctx context.Context, items []*MediaItem, registeredTools *RegisteredTools) error {
    var (
        wg         sync.WaitGroup
        mu         sync.Mutex
        allResults []*domain.MediaAnalytics
        allErrors  []error
    )
    
    // 並列実行数制限（5並列）
    semaphore := make(chan struct{}, 5)
    
    for _, item := range items {
        wg.Add(1)
        
        go func(item *MediaItem) {
            defer wg.Done()
            
            // セマフォ取得
            semaphore <- struct{}{}
            defer func() { <-semaphore }()
            
            // リトライ付き分析実行
            var analytics *domain.MediaAnalytics
            err := retry.Do(ctx, retry.DefaultConfig, func() error {
                result, analyzeErr := registeredTools.AnalyzeMedia.RunRaw(ctx, input)
                if analyzeErr != nil {
                    logger.Warn(ctx, "分析リトライ中", "error", analyzeErr)
                    return analyzeErr
                }
                
                analytics = parseAnalyticsResult(result)
                return nil
            })
            
            mu.Lock()
            defer mu.Unlock()
            
            if err != nil {
                allErrors = append(allErrors, err)
            } else {
                allResults = append(allResults, analytics)
                
                // DB保存
                fc := getFlowContext(ctx)
                fc.MediaAnalyticsRepo.Save(ctx, analytics)
            }
        }(item)
    }
    
    wg.Wait()
    
    // 全て失敗した場合のみエラー
    if len(allErrors) > 0 && len(allResults) == 0 {
        return fmt.Errorf("全ての分析が失敗しました: %v", allErrors)
    }
    
    return nil
}
```

### 6. フロントエンドSSE実装

**新規ファイル**: `frontend/api/mediaApi.ts`

```typescript
export const useMediaUploadSSE = (mediaId: string | null) => {
  const [status, setStatus] = useState<Media | null>(null)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!mediaId) return

    const sseUrl = `${process.env.NEXT_PUBLIC_API_URL}/api/media/${mediaId}/stream`
    const eventSource = new EventSource(sseUrl, { withCredentials: true })

    eventSource.onmessage = (event) => {
      const data = JSON.parse(event.data) as Media
      setStatus(data)

      if (data.status === 'failed') {
        setError(data.error_message)
        eventSource.close()
      } else if (data.status === 'completed') {
        eventSource.close()
      }
    }

    eventSource.onerror = () => {
      eventSource.close()
      // フォールバック: ポーリング
      startPolling(mediaId, setStatus, setError)
    }

    return () => eventSource.close()
  }, [mediaId])

  return { status, error }
}
```

**新規ファイル**: `frontend/app/upload/_components/UploadProgress.tsx`

```typescript
export const UploadProgress = ({ mediaIds }: { mediaIds: string[] }) => {
  return (
    <div className="space-y-4">
      {mediaIds.map(mediaId => (
        <MediaUploadItem key={mediaId} mediaId={mediaId} />
      ))}
    </div>
  )
}

const MediaUploadItem = ({ mediaId }: { mediaId: string }) => {
  const { status, error } = useMediaUploadSSE(mediaId)

  if (!status) return <Skeleton className="h-12" />

  return (
    <div className="border rounded-lg p-4">
      <div className="flex items-center justify-between mb-2">
        <span className="text-sm font-medium">
          {status.status === 'completed' && '✓ 完了'}
          {status.status === 'failed' && '✗ 失敗'}
          {status.status === 'uploading' && '⏳ アップロード中...'}
        </span>
        <span className="text-sm text-gray-500">
          {Math.round(status.progress * 100)}%
        </span>
      </div>
      
      <Progress value={status.progress * 100} />
      
      {error && (
        <Alert variant="destructive" className="mt-2">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}
    </div>
  )
}
```

## テスト計画

### 単体テスト

| ファイル | テストケース |
|---------|-------------|
| `backend/pkg/retry/retry_test.go` | ・成功時のリトライなし<br>・1回失敗後に成功<br>・最大リトライ超過<br>・コンテキストキャンセル |
| `backend/internal/agent/agent_test.go` | ・並列処理の正常動作<br>・部分失敗時の挙動<br>・全失敗時のエラー |

### 統合テスト

| シナリオ | 期待結果 |
|---------|---------|
| 複数ファイルアップロード | 全てのファイルが並列アップロードされる |
| アップロード失敗 | ユーザーにエラー通知 |
| 分析失敗（一時的） | 3回リトライ後に成功 |
| SSE接続失敗 | ポーリングにフォールバック |

## パフォーマンス改善見込み

| 処理 | 改善前 | 改善後 | 改善率 |
|-----|--------|--------|--------|
| **10枚の画像アップロード** | 50秒（順次） | 10秒（5並列） | **80%短縮** |
| **分析成功率** | 85%（リトライなし） | 98%（3回リトライ） | **15%向上** |

## 実装スケジュール

| Day | タスク | 担当ファイル |
|-----|--------|-------------|
| **1日目** | リトライロジック実装 | `pkg/retry/retry.go` |
| **2日目** | Mediaドメイン拡張 + マイグレーション | `domain/media.go`, migrations |
| **3日目** | アップロード非同期化 + SSE | `handler/media_server.go` |
| **4日目** | 画像分析並列処理化 | `agent/agent.go` |
| **5日目** | フロントエンドSSE + UI | `api/mediaApi.ts`, `UploadProgress.tsx` |

## リスクと対策

| リスク | 対策 |
|--------|------|
| Cloud Tasksの遅延 | ローカル環境ではGoroutineで即実行 |
| SSE接続の不安定性 | ポーリングフォールバック実装済み |
| 並列処理のAPIレート制限 | セマフォで5並列に制限 |

## 今後の拡張

- チャンクアップロード対応（100MB超のファイル）
- アップロード進捗の詳細化（バイト単位）
- 分析結果のキャッシュ機能

## 参考実装

- VLog生成SSE: `backend/internal/handler/vlog_server.go#L200-L240`
- フロントエンドSSE: `frontend/api/vlogAPi.ts`
- 非同期タスク: `backend/internal/handler/vlog_server.go#ProcessVLogTask`
