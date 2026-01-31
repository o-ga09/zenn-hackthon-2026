package genkit

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/storage"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"google.golang.org/genai"

	"github.com/o-ga09/zenn-hackthon-2026/internal/agent"
	"github.com/o-ga09/zenn-hackthon-2026/internal/domain"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/config"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/logger"
)

// GenkitAgent はGenkit AIエージェントの実装
// agent.IAgentインターフェースを実装する
type GenkitAgent struct {
	// FlowContext - 依存性を保持
	flowContext *FlowContext

	// 登録されたツール
	tools *RegisteredTools

	// VLog生成フロー
	vlogFlow VlogFlow

	// 設定
	baseURL string
}

// GenkitAgentOption はGenkitAgentを設定するためのオプション関数
type GenkitAgentOption func(*GenkitAgent)

// WithBaseURL はベースURLを設定するオプション
func WithBaseURL(baseURL string) GenkitAgentOption {
	return func(ga *GenkitAgent) {
		ga.baseURL = baseURL
	}
}

// WithAgentStorage はStorageを設定するオプション
func WithAgentStorage(storage domain.IImageStorage) GenkitAgentOption {
	return func(ga *GenkitAgent) {
		ga.flowContext.Storage = storage
	}
}

// WithAgentMediaRepository はMediaRepositoryを設定するオプション
func WithAgentMediaRepository(repo domain.IMediaRepository) GenkitAgentOption {
	return func(ga *GenkitAgent) {
		ga.flowContext.MediaRepo = repo
	}
}

// WithAgentMediaAnalyticsRepository はMediaAnalyticsRepositoryを設定するオプション
func WithAgentMediaAnalyticsRepository(repo domain.IMediaAnalyticsRepository) GenkitAgentOption {
	return func(ga *GenkitAgent) {
		ga.flowContext.MediaAnalyticsRepo = repo
	}
}

// WithAgentVlogRepository はVlogRepositoryを設定するオプション
func WithAgentVlogRepository(repo domain.IVLogRepository) GenkitAgentOption {
	return func(ga *GenkitAgent) {
		ga.flowContext.VlogRepo = repo
	}
}

// WithAgentGCSClient はGCSClientを設定するオプション
func WithAgentGCSClient(client *storage.Client) GenkitAgentOption {
	return func(ga *GenkitAgent) {
		ga.flowContext.GCSClient = client
	}
}

// WithAgentGenAIClient はGenAIClientを設定するオプション
func WithAgentGenAIClient(client *genai.Client) GenkitAgentOption {
	return func(ga *GenkitAgent) {
		ga.flowContext.GenAI = client
	}
}

// NewGenkitAgent は新しいGenkitAgentを作成する
func NewGenkitAgent(ctx context.Context, opts ...GenkitAgentOption) *GenkitAgent {
	g := config.GetGenkitCtx(ctx)

	// dotpromptファイルをロード
	genkit.LoadPromptDir(g, "prompts", "tavinikkiy")

	// FlowContextを初期化
	fc := NewFlowContext(
		WithGenkit(g),
		WithFlowConfig(DefaultFlowConfig()),
	)

	ga := &GenkitAgent{
		flowContext: fc,
		baseURL:     "https://tavinikkiy.example.com", // デフォルト
	}

	// オプションを適用
	for _, opt := range opts {
		opt(ga)
	}

	// ツールを登録
	ga.tools = RegisterAllTools(g, ga.baseURL)

	// VLog生成フローを登録
	ga.vlogFlow = RegisterVlogFlow(g, ga.tools)

	return ga
}

// CreateVlog はメディアからVLogを生成する
func (ga *GenkitAgent) CreateVlog(ctx context.Context, input *agent.VlogInput) (*agent.VlogOutput, error) {
	// FlowContextをコンテキストに設定
	ctx = WithFlowContext(ctx, ga.flowContext)

	// フローを実行
	output, err := ga.vlogFlow.Run(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("vlog flow failed: %w", err)
	}

	return output, nil
}

// CreateVlogWithProgress はメディアからVLogを生成し、進捗をコールバックで通知する
func (ga *GenkitAgent) CreateVlogWithProgress(ctx context.Context, input *agent.VlogInput, onProgress func(agent.FlowProgress)) (*agent.VlogOutput, error) {
	// FlowContextをコンテキストに設定
	ctx = WithFlowContext(ctx, ga.flowContext)

	// 初期化中を通知
	if onProgress != nil {
		onProgress(agent.FlowProgress{
			Step:     string(agent.StepInitializing),
			Progress: 0,
			Message:  "VLog生成を開始しています...",
		})
	}

	// 分析中を通知
	if onProgress != nil {
		onProgress(agent.FlowProgress{
			Step:       string(agent.StepAnalyzing),
			Progress:   10,
			Message:    "メディアを分析しています...",
			TotalItems: len(input.MediaItems),
		})
	}

	// フローを実行
	output, err := ga.vlogFlow.Run(ctx, input)
	if err != nil {
		if onProgress != nil {
			onProgress(agent.FlowProgress{
				Step:     string(agent.StepFailed),
				Progress: 0,
				Message:  fmt.Sprintf("エラーが発生しました: %v", err),
			})
		}
		return nil, fmt.Errorf("vlog flow failed: %w", err)
	}

	// 完了を通知
	if onProgress != nil {
		onProgress(agent.FlowProgress{
			Step:     string(agent.StepCompleted),
			Progress: 100,
			Message:  "VLog生成が完了しました！",
		})
	}

	return output, nil
}

// AnalyzeMediaBatch は複数のメディアを分析する
func (ga *GenkitAgent) AnalyzeMediaBatch(ctx context.Context, input *agent.MediaAnalysisBatchInput) (*agent.MediaAnalysisBatchOutput, error) {
	// FlowContextをコンテキストに設定
	ctx = WithFlowContext(ctx, ga.flowContext)

	// 各メディアを分析
	results := make([]agent.MediaAnalysisOutput, 0, len(input.Items))
	successfulItems := 0
	failedItems := 0

	locationMap := make(map[string]bool)
	activityMap := make(map[string]bool)

	for _, item := range input.Items {
		outputRaw, err := ga.tools.AnalyzeMedia.RunRaw(ctx, item)
		if err != nil {
			logger.Warn(ctx, fmt.Sprintf("FileID: %s, URL: %s", item.FileID, item.URL), "error", err.Error())
			failedItems++
			continue
		}
		output, err := convertToStruct[agent.MediaAnalysisOutput](outputRaw)
		if err != nil {
			logger.Warn(ctx, fmt.Sprintf("FileID: %s, URL: %s", item.FileID, item.URL), "error", "conversion failed: "+err.Error())
			failedItems++
			continue
		}
		results = append(results, output)
		successfulItems++

		// DBに保存
		if ga.flowContext.MediaAnalyticsRepo != nil {
			analytics := &domain.MediaAnalytics{
				FileID:      output.FileID,
				Description: output.Description,
				Objects:     output.Objects,
				Landmarks:   output.Landmarks,
				Activities:  output.Activities,
				Mood:        output.Mood,
			}
			if err := ga.flowContext.MediaAnalyticsRepo.Save(ctx, analytics); err != nil {
				logger.Warn(ctx, fmt.Sprintf("failed to save media analytics for file %s: %v", output.FileID, err))
			}
		}

		for _, l := range output.Landmarks {
			locationMap[l] = true
		}
		for _, a := range output.Activities {
			activityMap[a] = true
		}
	}

	uniqueLocations := make([]string, 0, len(locationMap))
	for l := range locationMap {
		uniqueLocations = append(uniqueLocations, l)
	}
	uniqueActivities := make([]string, 0, len(activityMap))
	for a := range activityMap {
		uniqueActivities = append(uniqueActivities, a)
	}

	return &agent.MediaAnalysisBatchOutput{
		Results: results,
		Summary: agent.MediaAnalysisSummary{
			TotalItems:       len(input.Items),
			SuccessfulItems:  successfulItems,
			FailedItems:      failedItems,
			UniqueLocations:  uniqueLocations,
			UniqueActivities: uniqueActivities,
			OverallMood:      "", // 必要に応じて設定
		},
	}, nil
}

// GetFlowContext は内部のFlowContextを取得する（テスト用）
func (ga *GenkitAgent) GetFlowContext() *FlowContext {
	return ga.flowContext
}

// GetTools は登録されたツールを取得する（テスト用）
func (ga *GenkitAgent) GetTools() *RegisteredTools {
	return ga.tools
}

// GenerateWithTools はツールを使用してAI生成を行うヘルパー関数
func GenerateWithTools[T any](ctx context.Context, ga *GenkitAgent, prompt string, toolList []ai.Tool) (*T, error) {
	// []ai.Tool を []ai.ToolRef に変換
	toolRefs := make([]ai.ToolRef, len(toolList))
	for i, t := range toolList {
		toolRefs[i] = t
	}

	result, _, err := genkit.GenerateData[T](ctx, ga.flowContext.Genkit,
		ai.WithPrompt(prompt),
		ai.WithTools(toolRefs...),
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// convertToStruct はinterface{}（通常はmap[string]interface{}）を指定した構造体に変換する
func convertToStruct[T any](raw interface{}) (T, error) {
	var result T

	// すでに目的の型の場合はそのまま返す
	if typed, ok := raw.(T); ok {
		return typed, nil
	}

	// JSONを経由して変換
	jsonBytes, err := json.Marshal(raw)
	if err != nil {
		return result, fmt.Errorf("failed to marshal to JSON: %w", err)
	}

	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		return result, fmt.Errorf("failed to unmarshal from JSON: %w", err)
	}

	return result, nil
}
