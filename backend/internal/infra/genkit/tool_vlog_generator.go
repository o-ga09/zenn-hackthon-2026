package genkit

import (
	"fmt"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/o-ga09/zenn-hackthon-2026/internal/agent"
	pkgerrors "github.com/o-ga09/zenn-hackthon-2026/pkg/errors"
)

const defaultShortDuration = 15.0 // デフォルトの短い動画の長さ（秒）

// GenerateVlogVideoInput はVLog動画生成ツールの入力
type GenerateVlogVideoInput struct {
	AnalysisResults []agent.MediaAnalysisOutput `json:"analysisResults,omitempty" jsonschema:"description=メディア分析結果のリスト"`
	Style           agent.VlogStyle             `json:"style,omitempty" jsonschema:"description=VLogのスタイル設定"`
	Title           string                      `json:"title,omitempty" jsonschema:"description=VLogのタイトル"`
	MediaItems      []agent.MediaItem           `json:"mediaItems,omitempty" jsonschema:"description=元のメディアアイテム"`
	UserID          string                      `json:"userId,omitempty" jsonschema:"description=ユーザーID"`
}

// GenerateVlogVideoOutput はVLog動画生成ツールの出力
type GenerateVlogVideoOutput struct {
	VideoURL    string                `json:"videoUrl" jsonschema:"description=生成された動画のURL"`
	VideoID     string                `json:"videoId" jsonschema:"description=動画ID"`
	Duration    float64               `json:"duration" jsonschema:"description=動画の長さ（秒）"`
	Title       string                `json:"title" jsonschema:"description=生成されたタイトル"`
	Description string                `json:"description" jsonschema:"description=生成された説明文"`
	Subtitles   []agent.SubtitleEntry `json:"subtitles" jsonschema:"description=字幕データ"`
}

// DefineGenerateVlogVideoTool はVLog動画生成ツールを定義する
func DefineGenerateVlogVideoTool(g *genkit.Genkit) ai.Tool {
	return genkit.DefineTool(g, "generateVlogVideo",
		"分析結果とメディアからVLog動画を生成する（Veo3を使用）",
		func(ctx *ai.ToolContext, input GenerateVlogVideoInput) (GenerateVlogVideoOutput, error) {
			fc := GetFlowContext(ctx)
			if fc == nil {
				return GenerateVlogVideoOutput{}, pkgerrors.ErrFlowContextNotFound
			}

			// タイトルと説明文を生成
			title := input.Title
			description := ""
			if title == "" {
				generated, err := generateTitleAndDescription(ctx, fc.Genkit, input.AnalysisResults)
				if err != nil {
					// タイトル生成失敗時はデフォルトを使用
					title = "Travel Vlog"
					description = ""
				} else {
					title = generated.Title
					description = generated.Description
				}
			}

			// 字幕を生成
			subtitles := generateSubtitles(input.AnalysisResults, input.Style)

			// 分析結果からプロンプト用サマリーを構築
			summaries := make([]MediaAnalysisSummary, 0, len(input.AnalysisResults))
			for _, r := range input.AnalysisResults {
				summaries = append(summaries, MediaAnalysisSummary{
					Description: r.Description,
					Landmarks:   r.Landmarks,
					Activities:  r.Activities,
					Mood:        r.Mood,
				})
			}

			// Veo用プロンプトを構築
			veoPrompt := BuildVlogPrompt(summaries, VlogStyleConfig{
				Theme:      input.Style.Theme,
				MusicMood:  input.Style.MusicMood,
				Duration:   input.Style.Duration,
				Transition: input.Style.Transition,
			})

			// UserIDを取得
			userID := input.UserID
			if userID == "" {
				userID = "anonymous"
			}

			// Veo3で動画生成（サポートされる長さ: 4, 6, 8秒のみ）
			duration := int32(input.Style.Duration)
			// Veo 3.1がサポートするのは 4, 6, 8 秒のみ
			if duration != 4 && duration != 6 && duration != 8 {
				duration = 8 // デフォルトは8秒
			}

			veoResult, err := GenerateVideoWithVeo(ctx, fc, VeoGenerateConfig{
				Prompt:          veoPrompt,
				DurationSeconds: duration,
				AspectRatio:     "16:9",
				UserID:          userID,
			})
			if err != nil {
				return GenerateVlogVideoOutput{}, fmt.Errorf("veo generation failed: %w", err)
			}

			return GenerateVlogVideoOutput{
				VideoURL:    veoResult.VideoURL,
				VideoID:     veoResult.VideoID,
				Duration:    veoResult.Duration,
				Title:       title,
				Description: description,
				Subtitles:   subtitles,
			}, nil
		},
	)
}

// TitleDescription はタイトルと説明文のペア
type TitleDescription struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// TitlePromptInput はgenerate_title.prompt用の入力
type TitlePromptInput struct {
	Locations  []string `json:"locations"`
	Activities []string `json:"activities"`
	Moods      []string `json:"moods"`
}

func generateTitleAndDescription(ctx *ai.ToolContext, g *genkit.Genkit, results []agent.MediaAnalysisOutput) (*TitleDescription, error) {
	var locations, activities, moods []string
	for _, r := range results {
		locations = append(locations, r.Landmarks...)
		activities = append(activities, r.Activities...)
		moods = append(moods, r.Mood)
	}

	// dotpromptを使用
	prompt := genkit.LookupPrompt(g, "tavinikkiy/generate_title")
	if prompt == nil {
		return nil, fmt.Errorf("prompt 'tavinikkiy/generate_title' not found")
	}

	promptInput := TitlePromptInput{
		Locations:  locations,
		Activities: activities,
		Moods:      moods,
	}

	resp, err := prompt.Execute(ctx, ai.WithInput(promptInput))
	if err != nil {
		return nil, fmt.Errorf("failed to execute prompt: %w", err)
	}

	var result TitleDescription
	if err := resp.Output(&result); err != nil {
		return nil, fmt.Errorf("failed to parse output: %w", err)
	}

	return &result, nil
}

func generateSubtitles(results []agent.MediaAnalysisOutput, style agent.VlogStyle) []agent.SubtitleEntry {
	subtitles := make([]agent.SubtitleEntry, 0, len(results))
	duration := float64(style.Duration)
	if duration == 0 {
		duration = defaultShortDuration
	}

	if len(results) == 0 {
		return subtitles
	}

	timePerMedia := duration / float64(len(results))

	for i, result := range results {
		startTime := float64(i) * timePerMedia
		endTime := startTime + timePerMedia - 0.5

		caption := result.SuggestedCaption
		if caption == "" {
			caption = result.Description
			if len(caption) > 50 {
				caption = caption[:47] + "..."
			}
		}

		subtitles = append(subtitles, agent.SubtitleEntry{
			StartTime: startTime,
			EndTime:   endTime,
			Text:      caption,
		})
	}

	return subtitles
}
