package agent

import "context"

// IAgent はAIエージェントのインターフェース
// VLog生成のユースケース層として機能する
type IAgent interface {
	// CreateVlog はメディアからVLogを生成する
	CreateVlog(ctx context.Context, input *VlogInput) (*VlogOutput, error)

	// CreateVlogWithProgress はメディアからVLogを生成し、進捗をコールバックで通知する
	CreateVlogWithProgress(ctx context.Context, input *VlogInput, onProgress func(FlowProgress)) (*VlogOutput, error)

	// AnalyzeMediaBatch は複数のメディアを分析する
	AnalyzeMediaBatch(ctx context.Context, input *MediaAnalysisBatchInput) (*MediaAnalysisBatchOutput, error)
}

// ProgressCallback は進捗通知用のコールバック型
type ProgressCallback func(progress FlowProgress)

// Agent はデフォルトのエージェント実装（スタブ）
type Agent struct{}

func NewAgent() *Agent {
	return &Agent{}
}
