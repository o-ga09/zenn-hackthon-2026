# Genkit AIエージェント実装設計

## 概要

Firebase Genkit を使用したAIエージェントの実装設計書です。VLog自動生成のためのフロー、ツール、依存性注入パターンについて記載します。

## アーキテクチャ

### パッケージ構成

```
backend/internal/infra/genkit/
├── agent.go               # GenkitAgent - メインエントリポイント
├── context.go             # FlowContext - 依存性注入コンテナ
├── flow.go                # VlogFlow - VLog生成フロー
├── tool_media_analyzer.go # メディア分析ツール
├── tool_registry.go       # ツール登録管理
├── tool_storage.go        # ストレージ関連ツール
├── tool_veo_client.go     # Veo3動画生成クライアント
└── tool_vlog_generator.go # VLog動画生成ツール

backend/pkg/errors/
└── genkit_errors.go       # Genkit関連エラー定義

backend/prompts/
├── analyze_media.prompt   # メディア分析プロンプト
├── generate_title.prompt  # タイトル生成プロンプト
└── generate_vlog.prompt   # VLog生成プロンプト
```

### 依存性注入パターン

```go
// FlowContext - Flow内で使用する依存性を保持
type FlowContext struct {
    Genkit    *genkit.Genkit
    Storage   domain.IImageStorage      // R2 Storage
    GCSClient *storage.Client           // GCS Client（Veo一時保存用）
    GenAI     *genai.Client             // Google Gen AIクライアント（Veo用）
    MediaRepo domain.IMediaRepository
    VlogRepo  domain.IVLogRepository
    Config    *FlowConfig
}

// FlowConfig - Flowの設定
type FlowConfig struct {
    DefaultModel         string  // デフォルトAIモデル
    MaxMediaItems        int     // 最大メディア数
    DefaultVideoDuration int     // デフォルト動画長（秒）
    ThumbnailWidth       int
    ThumbnailHeight      int
    // Veo設定
    VeoModel           string  // Veoモデル名
    GCSTempBucket      string  // GCS一時保存バケット
    GCSProjectID       string  // GCPプロジェクトID
    VeoPollingInterval int     // ポーリング間隔（秒）
    VeoMaxWaitTime     int     // 最大待機時間（秒）
}
```

## ツール一覧

### 1. analyzeMedia - メディア分析ツール

Gemini Vision APIを使用して画像・動画を分析します。

**入力:**
```json
{
  "fileId": "media_001",
  "url": "https://storage.example.com/images/photo1.jpg",
  "type": "image",
  "contentType": "image/jpeg"
}
```

**出力:**
```json
{
  "description": "青い海と白い砂浜の風景",
  "landmarks": ["沖縄", "ビーチ"],
  "activities": ["海水浴", "シュノーケリング"],
  "mood": "relaxing",
  "suggestedCaption": "美しい沖縄の海で癒しのひととき"
}
```

### 2. analyzeMediaBatch - バッチメディア分析ツール

複数のメディアを一括分析します。

### 3. uploadMedia - メディアアップロードツール

Cloudflare R2にメディアをアップロードします。

### 4. generateShareURL - 共有URL生成ツール

VLogの共有URLを生成します。

### 5. generateThumbnail - サムネイル生成ツール

動画からサムネイル画像を生成します。

### 6. generateVlogVideo - VLog動画生成ツール

**Veo3を使用して実際の動画を生成します。**

**入力:**
```json
{
  "analysisResults": [...],
  "style": {
    "theme": "adventure",
    "musicMood": "upbeat",
    "duration": 8,
    "transition": "fade"
  },
  "title": "沖縄旅行の思い出",
  "mediaItems": [...],
  "userId": "user123"
}
```

**出力:**
```json
{
  "videoUrl": "https://r2.example.com/users/user123/vlogs/01HXYZ.mp4",
  "videoId": "01HXYZ",
  "duration": 8.0,
  "title": "沖縄旅行の思い出",
  "description": "青い海と白い砂浜で過ごした最高の休日",
  "subtitles": [...]
}
```

## Veo3 動画生成フロー

```
1. VLog生成リクエスト受信
2. メディア分析 → プロンプト生成
3. Veo3 API呼び出し (GenerateVideos)
   ├─ GCS一時バケットに出力
   └─ ポーリングで完了待機（最大120秒）
4. GCSから動画データ取得
5. R2にアップロード
6. GCS一時ファイル削除（即時）
7. R2のURLをレスポンス
```

### Veo3 設定

| パラメータ | 説明 | デフォルト値 |
|-----------|------|-------------|
| VeoModel | Veoモデル名 | `veo-3.1-fast-generate-001` |
| GCSTempBucket | GCS一時保存バケット | `tavinikkiy-temp` |
| VeoPollingInterval | ポーリング間隔 | 5秒 |
| VeoMaxWaitTime | 最大待機時間 | 120秒 |
| DurationSeconds | 動画長 | 8秒 |
| AspectRatio | アスペクト比 | 16:9 |
| Resolution | 解像度 | 720p |

### 利用可能なVeoモデル

| モデルID | 特徴 |
|----------|------|
| `veo-3.1-generate-001` | 最新世代、高品質 |
| `veo-3.1-fast-generate-001` | 高速生成（推奨） |
| `veo-3.0-generate-001` | 安定版 |

## VLog生成Flow

### createVlogFlow

メインのVLog生成フローです。

```go
func RegisterVlogFlow(g *genkit.Genkit, registeredTools *RegisteredTools) VlogFlow {
    return genkit.DefineFlow(g, "createVlogFlow", func(ctx context.Context, input *agent.VlogInput) (*agent.VlogOutput, error) {
        // Step 1: メディア分析
        // Step 2: VLog動画生成（Veo3）
        // Step 3: サムネイル生成
        // Step 4: 共有URL生成
        // Step 5: 分析サマリー構築
        return output, nil
    })
}
```

### 入力スキーマ (VlogInput)

```json
{
  "userId": "user123",
  "mediaItems": [
    {
      "fileId": "media_001",
      "url": "https://storage.example.com/photo1.jpg",
      "type": "image",
      "contentType": "image/jpeg",
      "timestamp": "2026-01-15T10:30:00Z",
      "order": 1
    }
  ],
  "title": "沖縄旅行の思い出",
  "travelDate": "2026-01-15",
  "destination": "沖縄",
  "style": {
    "theme": "adventure",
    "musicMood": "upbeat",
    "duration": 8,
    "transition": "fade"
  }
}
```

### 出力スキーマ (VlogOutput)

```json
{
  "videoId": "01HXYZ",
  "videoUrl": "https://r2.example.com/users/user123/vlogs/01HXYZ.mp4",
  "shareUrl": "https://tavinikkiy.example.com/share/ABCDEF",
  "thumbnailUrl": "https://r2.example.com/thumbnails/01HXYZ.jpg",
  "duration": 8.0,
  "title": "沖縄旅行の思い出",
  "description": "青い海と白い砂浜で過ごした最高の休日",
  "subtitles": [
    {
      "startTime": 0.0,
      "endTime": 2.5,
      "text": "美しい沖縄の海で癒しのひととき"
    }
  ],
  "analytics": {
    "locations": ["沖縄", "ビーチ"],
    "activities": ["海水浴", "シュノーケリング"],
    "mood": "relaxing",
    "highlights": ["美しい沖縄の海で癒しのひととき"],
    "mediaCount": 5
  }
}
```

## Dotprompt テンプレート

### analyze_media.prompt

メディア分析用のプロンプトテンプレート。Gemini Vision APIで画像・動画を分析します。

### generate_title.prompt

分析結果からタイトルと説明文を生成するプロンプトテンプレート。

### generate_vlog.prompt

VLog生成用のプロンプトテンプレート。

## 初期化とDI

### サーバー初期化

```go
// server.go
func New(ctx context.Context) *Server {
    // GCSクライアントの初期化
    gcsClient, err := config.GetGCSClient(ctx)
    if err != nil {
        log.Printf("warning: failed to initialize GCS client: %v", err)
    }

    // GenAIクライアントの初期化
    genaiClient, err := config.GetGenAIClient(ctx)
    if err != nil {
        log.Printf("warning: failed to initialize GenAI client: %v", err)
    }

    // GenkitAgent の初期化（依存性注入）
    genkitAgent := genkit.NewGenkitAgent(ctx,
        genkit.WithAgentStorage(r2Storage),
        genkit.WithAgentGCSClient(gcsClient),
        genkit.WithAgentGenAIClient(genaiClient),
        genkit.WithBaseURL(env.BASE_URL),
    )
}
```

### 環境変数

| 環境変数 | 説明 | デフォルト値 |
|----------|------|-------------|
| `PROJECT_ID` | GCPプロジェクトID | `tavinikkiy` |
| `GCS_TEMP_BUCKET` | GCS一時保存バケット | `tavinikkiy-temp` |
| `GCS_LOCATION` | GCSリージョン | `us-central1` |
| `GOOGLE_APPLICATION_CREDENTIALS` | サービスアカウントJSONパス | - |

## エラーハンドリング

### 定義済みエラー

```go
// pkg/errors/genkit_errors.go

var (
    ErrGenkitNotInitialized   = errors.New("genkit not initialized")
    ErrStorageNotInitialized  = errors.New("storage not initialized")
    ErrFlowContextNotFound    = errors.New("flow context not found in context")
    ErrInvalidInput           = errors.New("invalid input")
    ErrNoMediaItems           = errors.New("no media items provided")
    ErrMaxMediaItemsExceeded  = errors.New("max media items exceeded")
    ErrMediaAnalysisFailed    = errors.New("media analysis failed")
    ErrToolExecutionFailed    = errors.New("tool execution failed")
)
```

## テスト

### ユニットテスト

各ツールとフローのユニットテストを実装します。

### 統合テスト

Veo3 APIの統合テストはモック化して実行します。

## 今後の拡張

1. **非同期処理対応**: WebSocket/SSEによる進捗通知
2. **キャッシュ**: メディア分析結果のキャッシュ
3. **リトライ**: Veo API呼び出しのリトライ機構
4. **コスト管理**: Veo APIの使用量監視とアラート
