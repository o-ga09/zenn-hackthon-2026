# AI開発ガイドライン

## Firebase Genkit for Go開発パターン

### エージェント設計原則

- **単一責任**: 各エージェントは明確に定義された単一の目的を持つ
- **ステートレス**: エージェント間で状態を共有せず、必要な情報はすべて入力として受け取る
- **エラーハンドリング**: 外部API呼び出しは必ずリトライとサーキットブレーカーで保護
- **進行状況通知**: 長時間処理はSSEで進行状況をリアルタイム通知

### Genkit Agent実装パターン

```go
// エージェントインターフェース定義
type VlogGenerationAgent interface {
    GenerateVlog(ctx context.Context, req VlogGenerationRequest) (*VlogGenerationResponse, error)
}

// 実装例
func (a *VlogAgentImpl) GenerateVlog(ctx context.Context, req VlogGenerationRequest) (*VlogGenerationResponse, error) {
    // 1. 入力検証
    if err := a.validateRequest(req); err != nil {
        return nil, fmt.Errorf("invalid request: %w", err)
    }

    // 2. ジョブ作成・進行状況初期化
    job, err := a.jobService.CreateJob(ctx, req.UserID, req.SessionToken, "vlog_generation")
    if err != nil {
        return nil, fmt.Errorf("failed to create job: %w", err)
    }

    // 3. 非同期処理開始
    go a.processVlogGeneration(ctx, job.ID, req)

    return &VlogGenerationResponse{
        JobID:       job.ID,
        VlogID:      req.VlogID,
        Status:      "queued",
        SSEEndpoint: fmt.Sprintf("/api/sse/connect/%s", job.ID),
    }, nil
}
```

### 外部API統合パターン

#### Vertex AI (Gemini 3) 統合

```go
type GeminiClient struct {
    client     *genai.Client
    retryConfig RetryConfig
    circuitBreaker *CircuitBreaker
}

func (g *GeminiClient) AnalyzeImages(ctx context.Context, images []ImageData) (*ImageAnalysisResult, error) {
    return g.circuitBreaker.Execute(func() (*ImageAnalysisResult, error) {
        return g.callGeminiWithRetry(ctx, images)
    })
}
```

#### Veo3 動画生成統合

```go
type Veo3Client struct {
    endpoint string
    apiKey   string
    retryConfig RetryConfig
}

func (v *Veo3Client) GenerateVideo(ctx context.Context, prompt string, images []ImageData) (*VideoResult, error) {
    // 非同期動画生成リクエスト
    jobID, err := v.submitGenerationJob(ctx, prompt, images)
    if err != nil {
        return nil, err
    }

    // ポーリングで完了待ち
    return v.pollForCompletion(ctx, jobID)
}
```

## SSE (Server-Sent Events) 開発パターン

### SSE Manager実装

```go
type SSEManager struct {
    connections map[string]*Connection
    mutex       sync.RWMutex
}

func (m *SSEManager) AddConnection(jobID, userID string, w http.ResponseWriter) *Connection {
    m.mutex.Lock()
    defer m.mutex.Unlock()

    conn := &Connection{
        JobID:   jobID,
        UserID:  userID,
        Channel: make(chan SSEEvent, 100),
        Context: context.Background(),
    }

    m.connections[jobID] = conn
    go m.handleConnection(conn, w)

    return conn
}
```

### イベント送信パターン

```go
// 進行状況更新
func (s *VlogService) notifyProgress(jobID string, step string, progress int) {
    event := SSEEvent{
        Type: "progress",
        Data: map[string]interface{}{
            "current_step": step,
            "progress_percentage": progress,
            "timestamp": time.Now(),
        },
    }
    s.sseManager.SendEvent(jobID, event)
}

// エラー通知
func (s *VlogService) notifyError(jobID string, err error) {
    event := SSEEvent{
        Type: "error",
        Data: map[string]interface{}{
            "error_code": "processing_failed",
            "message": "動画生成中にエラーが発生しました",
            "timestamp": time.Now(),
        },
    }
    s.sseManager.SendEvent(jobID, event)
}
```

## 非同期処理パターン

### ジョブキュー実装

```go
type JobQueue struct {
    jobs    chan Job
    workers int
    wg      sync.WaitGroup
}

func (q *JobQueue) Start() {
    for i := 0; i < q.workers; i++ {
        q.wg.Add(1)
        go q.worker()
    }
}

func (q *JobQueue) worker() {
    defer q.wg.Done()
    for job := range q.jobs {
        q.processJob(job)
    }
}
```

### タイムアウト管理

```go
func (s *VlogService) processWithTimeout(ctx context.Context, job *VlogJob) error {
    // 画像枚数に応じたタイムアウト設定
    timeout := s.calculateTimeout(len(job.InputData.Images))
    ctx, cancel := context.WithTimeout(ctx, timeout)
    defer cancel()

    // 処理実行
    return s.executeVlogGeneration(ctx, job)
}
```

## エラーハンドリングパターン

### リトライ機能

```go
func (c *ExternalAPIClient) CallWithRetry(ctx context.Context, apiCall func() error) error {
    var lastErr error
    delay := c.config.InitialDelay

    for i := 0; i <= c.config.MaxRetries; i++ {
        if err := apiCall(); err == nil {
            return nil
        } else {
            lastErr = err
            if i < c.config.MaxRetries {
                select {
                case <-time.After(delay):
                    delay = time.Duration(float64(delay) * c.config.BackoffFactor)
                case <-ctx.Done():
                    return ctx.Err()
                }
            }
        }
    }
    return lastErr
}
```

### サーキットブレーカー

```go
type CircuitBreaker struct {
    maxFailures int
    resetTimeout time.Duration
    failures int
    lastFailureTime time.Time
    state CircuitState
    mutex sync.RWMutex
}

func (cb *CircuitBreaker) Execute(fn func() error) error {
    cb.mutex.RLock()
    state := cb.state
    cb.mutex.RUnlock()

    if state == CircuitOpen {
        if time.Since(cb.lastFailureTime) > cb.resetTimeout {
            cb.setState(CircuitHalfOpen)
        } else {
            return ErrCircuitOpen
        }
    }

    err := fn()
    cb.recordResult(err)
    return err
}
```

## テストパターン

### プロパティベーステスト

```go
func TestVlogGenerationProperties(t *testing.T) {
    properties := gopter.NewProperties(gopter.DefaultTestParameters())

    // プロパティ1: AIエージェント処理開始
    properties.Property("Feature: ai-vlog-generation, Property 1: AIエージェント処理開始",
        prop.ForAll(
            func(images []ImageData) bool {
                req := VlogGenerationRequest{Images: images}
                resp, err := agent.GenerateVlog(context.Background(), req)
                return err == nil && resp.JobID != ""
            },
            genValidImageSet(),
        ))

    properties.TestingRun(t, gopter.ConsoleReporter(false))
}
```

### モックパターン

```go
type MockGeminiClient struct {
    responses map[string]*ImageAnalysisResult
    errors    map[string]error
}

func (m *MockGeminiClient) AnalyzeImages(ctx context.Context, images []ImageData) (*ImageAnalysisResult, error) {
    key := generateImageKey(images)
    if err, exists := m.errors[key]; exists {
        return nil, err
    }
    return m.responses[key], nil
}
```

## パフォーマンス最適化

### 並行処理

```go
func (s *VlogService) processImagesParallel(ctx context.Context, images []ImageData) ([]*ImageAnalysis, error) {
    results := make([]*ImageAnalysis, len(images))
    errors := make([]error, len(images))

    var wg sync.WaitGroup
    for i, image := range images {
        wg.Add(1)
        go func(idx int, img ImageData) {
            defer wg.Done()
            result, err := s.geminiClient.AnalyzeImage(ctx, img)
            results[idx] = result
            errors[idx] = err
        }(i, image)
    }

    wg.Wait()

    // エラーチェック
    for _, err := range errors {
        if err != nil {
            return nil, err
        }
    }

    return results, nil
}
```

### キャッシュ活用

```go
type CachedGeminiClient struct {
    client GeminiClient
    cache  map[string]*ImageAnalysisResult
    mutex  sync.RWMutex
    ttl    time.Duration
}

func (c *CachedGeminiClient) AnalyzeImages(ctx context.Context, images []ImageData) (*ImageAnalysisResult, error) {
    key := generateCacheKey(images)

    c.mutex.RLock()
    if cached, exists := c.cache[key]; exists {
        c.mutex.RUnlock()
        return cached, nil
    }
    c.mutex.RUnlock()

    result, err := c.client.AnalyzeImages(ctx, images)
    if err == nil {
        c.mutex.Lock()
        c.cache[key] = result
        c.mutex.Unlock()
    }

    return result, err
}
```

## ベストプラクティス

### 構造化ログ

```go
func (s *VlogService) logJobProgress(jobID, step string, metadata map[string]interface{}) {
    logger.Info("job_progress",
        zap.String("job_id", jobID),
        zap.String("step", step),
        zap.Any("metadata", metadata),
        zap.Time("timestamp", time.Now()),
    )
}
```

### メトリクス収集

```go
func (s *VlogService) recordMetrics(jobID string, duration time.Duration, success bool) {
    metrics.RecordJobDuration(duration)
    if success {
        metrics.IncrementJobSuccess()
    } else {
        metrics.IncrementJobFailure()
    }
}
```

### 設定管理

```go
type AIConfig struct {
    GeminiAPIKey     string        `env:"GEMINI_API_KEY,required"`
    Veo3Endpoint     string        `env:"VEO3_ENDPOINT,required"`
    MaxRetries       int           `env:"MAX_RETRIES" envDefault:"3"`
    RequestTimeout   time.Duration `env:"REQUEST_TIMEOUT" envDefault:"30s"`
    CircuitThreshold int           `env:"CIRCUIT_THRESHOLD" envDefault:"5"`
}
```

## 開発フロー

1. **仕様定義**: .kiro/specs/ でプロパティベーステスト仕様を定義
2. **インターフェース設計**: ドメイン層でインターフェースを定義
3. **実装**: インフラ層で具体実装
4. **テスト**: ユニットテスト + プロパティベーステスト
5. **統合**: サービス層で組み合わせ
6. **E2Eテスト**: 完全フローの動作確認

## 注意事項

- 外部API呼び出しは必ずタイムアウトとリトライを設定
- 長時間処理は必ずSSEで進行状況を通知
- エラーメッセージはセキュリティを考慮してユーザーフレンドリーに
- プロパティベーステストは最低100回実行
- ログには個人情報を含めない
