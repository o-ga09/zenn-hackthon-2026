package genkit

import (
	"fmt"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/o-ga09/zenn-hackthon-2026/internal/agent"
	pkgerrors "github.com/o-ga09/zenn-hackthon-2026/pkg/errors"
)

// AnalyzeMediaPromptInput はanalyze_media.prompt用の入力
type AnalyzeMediaPromptInput struct {
	MediaType string `json:"mediaType"`
	FileID    string `json:"fileId"`
}

// DefineAnalyzeMediaTool はメディア分析ツールを定義する
// Gemini Vision APIを使用して画像/動画を分析し、情報を抽出する
func DefineAnalyzeMediaTool(g *genkit.Genkit) ai.Tool {
	return genkit.DefineTool(g, "analyzeMedia",
		"メディア（画像/動画）を分析し、シーンの説明、オブジェクト、ランドマーク、アクティビティ、雰囲気を抽出する",
		func(ctx *ai.ToolContext, input agent.MediaAnalysisInput) (agent.MediaAnalysisOutput, error) {
			fc := GetFlowContext(ctx)
			if fc == nil {
				return agent.MediaAnalysisOutput{}, pkgerrors.ErrFlowContextNotFound
			}

			// dotpromptを使用してメディアを分析
			prompt := genkit.LookupPrompt(fc.Genkit, "tavinikkiy/analyze_media")
			if prompt == nil {
				return agent.MediaAnalysisOutput{}, fmt.Errorf("prompt 'tavinikkiy/analyze_media' not found")
			}

			// プロンプト入力を準備
			mediaType := "画像"
			if input.Type == "video" {
				mediaType = "動画"
			}

			promptInput := AnalyzeMediaPromptInput{
				MediaType: mediaType,
				FileID:    input.FileID,
			}

			// メディアパーツを追加
			var mediaParts []*ai.Part
			if input.Type == "image" || input.Type == "video" {
				mediaParts = append(mediaParts, ai.NewResourcePart(input.URL))
			}

			// プロンプトを実行
			resp, err := prompt.Execute(ctx, ai.WithInput(promptInput), ai.WithMessages(ai.NewUserMessage(mediaParts...)))
			if err != nil {
				return agent.MediaAnalysisOutput{}, fmt.Errorf("%w: %v", pkgerrors.ErrMediaAnalysisFailed, err)
			}

			// レスポンスをパース
			var result agent.MediaAnalysisOutput
			if err := resp.Output(&result); err != nil {
				return agent.MediaAnalysisOutput{}, fmt.Errorf("%w: failed to parse output: %v", pkgerrors.ErrMediaAnalysisFailed, err)
			}

			result.FileID = input.FileID
			result.Type = input.Type

			return result, nil
		},
	)
}

// BatchAnalysisInput は複数メディアの一括分析入力
type BatchAnalysisInput struct {
	Items []agent.MediaAnalysisInput `json:"items" jsonschema:"description=分析対象のメディアリスト"`
}

// BatchAnalysisOutput は複数メディアの一括分析出力
type BatchAnalysisOutput struct {
	Results []agent.MediaAnalysisOutput `json:"results" jsonschema:"description=分析結果のリスト"`
	Summary AnalysisSummary             `json:"summary" jsonschema:"description=全体のサマリー"`
}

// AnalysisSummary は分析結果の全体サマリー
type AnalysisSummary struct {
	TotalItems       int      `json:"totalItems"`
	SuccessfulItems  int      `json:"successfulItems"`
	FailedItems      int      `json:"failedItems"`
	UniqueLocations  []string `json:"uniqueLocations"`
	UniqueActivities []string `json:"uniqueActivities"`
	OverallMood      string   `json:"overallMood"`
}

// DefineAnalyzeMediaBatchTool は複数メディアの一括分析ツールを定義する
// NOTE: バッチ分析は現在Flow層で実装されているため、このツールは将来の拡張用
func DefineAnalyzeMediaBatchTool(g *genkit.Genkit, _ ai.Tool) ai.Tool {
	return genkit.DefineTool(g, "analyzeMediaBatch",
		"複数のメディアを一括で分析し、全体のサマリーを生成する",
		func(ctx *ai.ToolContext, input BatchAnalysisInput) (BatchAnalysisOutput, error) {
			// NOTE: 現在はFlow層でループ処理しているため、このツールは空実装
			// 将来的にはここで並列分析を実装可能
			return BatchAnalysisOutput{
				Results: []agent.MediaAnalysisOutput{},
				Summary: AnalysisSummary{
					TotalItems:       len(input.Items),
					SuccessfulItems:  0,
					FailedItems:      len(input.Items),
					UniqueLocations:  []string{},
					UniqueActivities: []string{},
					OverallMood:      "",
				},
			}, nil
		},
	)
}
