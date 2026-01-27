package genkit

import (
	"context"

	"cloud.google.com/go/storage"
	"github.com/firebase/genkit/go/genkit"
	"google.golang.org/genai"

	"github.com/o-ga09/zenn-hackthon-2026/internal/domain"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/errors"
)

// ============================================================
// FlowContext - Flow内で使用する依存性のコンテナ
// ============================================================

// FlowContext はGenkit Flow内で使用する依存性を保持する
type FlowContext struct {
	Genkit    *genkit.Genkit
	Storage   domain.IImageStorage
	GCSClient *storage.Client  // GCSクライアント（Veo一時保存用）
	GenAI     *genai.Client    // Google Gen AIクライアント（Veo用）
	MediaRepo domain.IMediaRepository
	VlogRepo  domain.IVLogRepository
	Config    *FlowConfig
}

// FlowConfig はFlowの設定を保持する
type FlowConfig struct {
	DefaultModel         string
	MaxMediaItems        int
	DefaultVideoDuration int
	ThumbnailWidth       int
	ThumbnailHeight      int
	// Veo設定
	VeoModel           string // Veoモデル名
	GCSTempBucket      string // GCS一時保存バケット
	GCSProjectID       string // GCPプロジェクトID
	VeoPollingInterval int    // ポーリング間隔（秒）
	VeoMaxWaitTime     int    // 最大待機時間（秒）
}

// DefaultFlowConfig はデフォルトのFlowConfigを返す
func DefaultFlowConfig() *FlowConfig {
	return &FlowConfig{
		DefaultModel:         "googleai/gemini-2.5-flash",
		MaxMediaItems:        50,
		DefaultVideoDuration: 8, // Veoは8秒が標準
		ThumbnailWidth:       1280,
		ThumbnailHeight:      720,
		// Veo設定
		VeoModel:           "veo-3.1-fast-generate-001",
		GCSTempBucket:      "tavinikkiy-temp",
		GCSProjectID:       "tavinikkiy",
		VeoPollingInterval: 5,
		VeoMaxWaitTime:     120,
	}
}

// ============================================================
// FlowContextOption - オプションパターンによる依存性注入
// ============================================================

// FlowContextOption はFlowContextを設定するためのオプション関数
type FlowContextOption func(*FlowContext)

// WithStorage はStorageを設定するオプション
func WithStorage(storage domain.IImageStorage) FlowContextOption {
	return func(fc *FlowContext) {
		fc.Storage = storage
	}
}

// WithMediaRepository はMediaRepositoryを設定するオプション
func WithMediaRepository(repo domain.IMediaRepository) FlowContextOption {
	return func(fc *FlowContext) {
		fc.MediaRepo = repo
	}
}

// WithVlogRepository はVlogRepositoryを設定するオプション
func WithVlogRepository(repo domain.IVLogRepository) FlowContextOption {
	return func(fc *FlowContext) {
		fc.VlogRepo = repo
	}
}

// WithGCSClient はGCSClientを設定するオプション
func WithGCSClient(client *storage.Client) FlowContextOption {
	return func(fc *FlowContext) {
		fc.GCSClient = client
	}
}

// WithGenAIClient はGenAIClientを設定するオプション
func WithGenAIClient(client *genai.Client) FlowContextOption {
	return func(fc *FlowContext) {
		fc.GenAI = client
	}
}

// WithFlowConfig はFlowConfigを設定するオプション
func WithFlowConfig(config *FlowConfig) FlowContextOption {
	return func(fc *FlowContext) {
		fc.Config = config
	}
}

// WithGenkitInstance はGenkitインスタンスを設定するオプション
func WithGenkitInstance(g *genkit.Genkit) FlowContextOption {
	return func(fc *FlowContext) {
		fc.Genkit = g
	}
}

// WithGenkit はWithGenkitInstanceのエイリアス
func WithGenkit(g *genkit.Genkit) FlowContextOption {
	return WithGenkitInstance(g)
}

// NewFlowContext は新しいFlowContextを作成する
func NewFlowContext(opts ...FlowContextOption) *FlowContext {
	fc := &FlowContext{
		Config: DefaultFlowConfig(),
	}
	for _, opt := range opts {
		opt(fc)
	}
	return fc
}

// ============================================================
// Context Key
// ============================================================

type flowContextKey struct{}

// WithFlowContext はFlowContextをコンテキストに設定する
func WithFlowContext(ctx context.Context, fc *FlowContext) context.Context {
	return context.WithValue(ctx, flowContextKey{}, fc)
}

// GetFlowContext はコンテキストからFlowContextを取得する
func GetFlowContext(ctx context.Context) *FlowContext {
	if fc, ok := ctx.Value(flowContextKey{}).(*FlowContext); ok {
		return fc
	}
	return nil
}

// MustGetFlowContext はコンテキストからFlowContextを取得し、存在しない場合はpanicする
func MustGetFlowContext(ctx context.Context) *FlowContext {
	fc := GetFlowContext(ctx)
	if fc == nil {
		panic("FlowContext not found in context")
	}
	return fc
}

// ============================================================
// Validation
// ============================================================

// Validate はFlowContextが有効かどうかを検証する
func (fc *FlowContext) Validate() error {
	if fc.Genkit == nil {
		return errors.ErrGenkitNotInitialized
	}
	if fc.Storage == nil {
		return errors.ErrStorageNotInitialized
	}
	return nil
}

// ValidateForVlogGeneration はVLog生成に必要な依存性が設定されているかを検証する
func (fc *FlowContext) ValidateForVlogGeneration() error {
	return fc.Validate()
}
