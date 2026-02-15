package genkit

import (
	"encoding/base64"
	"fmt"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/o-ga09/zenn-hackthon-2026/internal/agent"
	pkgerrors "github.com/o-ga09/zenn-hackthon-2026/pkg/errors"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/http"
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

			// メディアパーツを追加（URLからデータを取得してBase64で渡す）
			var mediaParts []*ai.Part
			fmt.Println("❓", input.Type)
			if input.Type == "image" || input.Type == "video" {
				mediaData, contentType, err := http.FetchMediaData(input.URL, input.ContentType)
				if err != nil {
					return agent.MediaAnalysisOutput{}, fmt.Errorf("%w: failed to fetch media: %v", pkgerrors.ErrMediaAnalysisFailed, err)
				}
				fmt.Println("❓", mediaData[:100], "...", contentType) // デバッグ用ログ
				// Base64エンコードしてdata URIとして渡す
				dataURI := fmt.Sprintf("data:%s;base64,%s", contentType, base64.StdEncoding.EncodeToString(mediaData))
				mediaParts = append(mediaParts, ai.NewMediaPart(contentType, dataURI))
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

// DefineAnalyzeMediaBatchTool は複数メディアの一括分析ツールを定義する
// NOTE: 現在はGenkitAgent.AnalyzeMediaBatchで実装されているため、このツールは将来のリファクタリング用
func DefineAnalyzeMediaBatchTool(g *genkit.Genkit) ai.Tool {
	return genkit.DefineTool(g, "analyzeMediaBatch",
		"複数のメディアを一括で分析し、全体のサマリーを生成する",
		func(ctx *ai.ToolContext, input agent.MediaAnalysisBatchInput) (agent.MediaAnalysisBatchOutput, error) {
			return agent.MediaAnalysisBatchOutput{
				Results: []agent.MediaAnalysisOutput{},
				Summary: agent.MediaAnalysisSummary{
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
