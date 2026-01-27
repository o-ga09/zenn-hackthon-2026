package agent

import "context"

// IAgent はAIエージェントのインターフェース
// VLog生成のユースケース層として機能する
type IAgent interface {
	// CreateVlog はメディアからVLogを生成する
	CreateVlog(ctx context.Context, input *VlogInput) (*VlogOutput, error)

	// CreateVlogWithProgress はメディアからVLogを生成し、進捗をコールバックで通知する
	CreateVlogWithProgress(ctx context.Context, input *VlogInput, onProgress func(FlowProgress)) (*VlogOutput, error)

	// AnalyzeMedia は単一のメディアを分析する（ツール単体テスト用）
	AnalyzeMedia(ctx context.Context, input *MediaAnalysisInput) (*MediaAnalysisOutput, error)
}

// ProgressCallback は進捗通知用のコールバック型
type ProgressCallback func(progress FlowProgress)

// Agent はデフォルトのエージェント実装（スタブ）
type Agent struct{}

func NewAgent() *Agent {
	return &Agent{}
}
