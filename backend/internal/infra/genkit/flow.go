package genkit

import (
	"context"
	"fmt"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/core"
	"github.com/firebase/genkit/go/genkit"
	"github.com/o-ga09/zenn-hackthon-2026/internal/agent"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/errors"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/ulid"
)

// ============================================================
// VLog生成Flow
// ============================================================

// VlogFlow はVLog生成フローの型エイリアス
type VlogFlow = *core.Flow[*agent.VlogInput, *agent.VlogOutput, struct{}]

// RegisterVlogFlow はVLog生成フローを登録する
func RegisterVlogFlow(g *genkit.Genkit, registeredTools *RegisteredTools) VlogFlow {
	return genkit.DefineFlow(g, "createVlogFlow", func(ctx context.Context, input *agent.VlogInput) (*agent.VlogOutput, error) {
		// FlowContextを取得
		fc := GetFlowContext(ctx)
		if fc == nil {
			return nil, errors.ErrFlowContextNotFound
		}

		// 入力バリデーション
		if err := validateVlogInput(input, fc.Config); err != nil {
			return nil, err
		}

		// Step 1: メディア分析
		analysisResults, err := analyzeAllMedia(ctx, input.MediaItems, registeredTools)
		if err != nil {
			return nil, fmt.Errorf("media analysis failed: %w", err)
		}

		// Step 2: VLog動画生成
		videoResult, err := generateVlog(ctx, fc.Genkit, input, analysisResults, registeredTools)
		if err != nil {
			return nil, fmt.Errorf("video generation failed: %w", err)
		}

		// Step 3: サムネイル生成
		thumbnailRaw, err := registeredTools.GenerateThumbnail.RunRaw(ctx, GenerateThumbnailInput{
			VideoURL: videoResult.VideoURL,
			VideoID:  videoResult.VideoID,
		})
		var thumbnailResult GenerateThumbnailOutput
		if err != nil {
			// サムネイル生成失敗は致命的ではない
			thumbnailResult = GenerateThumbnailOutput{
				ThumbnailURL: "",
			}
		} else {
			thumbnailResult = thumbnailRaw.(GenerateThumbnailOutput)
		}

		// Step 4: 共有URL生成
		shareRaw, err := registeredTools.GenerateShareURL.RunRaw(ctx, GenerateShareURLInput{
			VideoID: videoResult.VideoID,
			UserID:  input.UserID,
		})
		if err != nil {
			return nil, fmt.Errorf("share URL generation failed: %w", err)
		}
		shareResult := shareRaw.(GenerateShareURLOutput)

		// 分析サマリーを構築
		analytics := buildAnalyticsSummary(analysisResults, len(input.MediaItems))

		return &agent.VlogOutput{
			VideoID:      videoResult.VideoID,
			VideoURL:     videoResult.VideoURL,
			ShareURL:     shareResult.ShareURL,
			ThumbnailURL: thumbnailResult.ThumbnailURL,
			Duration:     videoResult.Duration,
			Title:        videoResult.Title,
			Description:  videoResult.Description,
			Subtitles:    videoResult.Subtitles,
			Analytics:    analytics,
		}, nil
	})
}

// validateVlogInput は入力を検証する
func validateVlogInput(input *agent.VlogInput, config *FlowConfig) error {
	if input == nil {
		return errors.ErrInvalidInput
	}
	if len(input.MediaItems) == 0 {
		return errors.ErrNoMediaItems
	}
	if len(input.MediaItems) > config.MaxMediaItems {
		return errors.ErrMaxMediaItemsExceeded
	}
	if input.UserID == "" {
		return fmt.Errorf("%w: userID is required", errors.ErrInvalidInput)
	}
	return nil
}

// analyzeAllMedia は全メディアを分析する
func analyzeAllMedia(ctx context.Context, items []agent.MediaItem, registeredTools *RegisteredTools) ([]agent.MediaAnalysisOutput, error) {
	results := make([]agent.MediaAnalysisOutput, 0, len(items))

	for _, item := range items {
		resultRaw, err := registeredTools.AnalyzeMedia.RunRaw(ctx, agent.MediaAnalysisInput{
			FileID:      item.FileID,
			URL:         item.URL,
			Type:        item.Type,
			ContentType: item.ContentType,
		})
		if err != nil {
			// 分析失敗は警告として続行
			continue
		}
		result := resultRaw.(agent.MediaAnalysisOutput)
		results = append(results, result)
	}

	if len(results) == 0 {
		return nil, errors.ErrMediaAnalysisFailed
	}

	return results, nil
}

// generateVlog はAIモデルを使用してVLogを生成する
func generateVlog(ctx context.Context, g *genkit.Genkit, input *agent.VlogInput, analysisResults []agent.MediaAnalysisOutput, registeredTools *RegisteredTools) (*GenerateVlogVideoOutput, error) {
	// プロンプトを構築
	prompt := buildVlogGenerationPrompt(input, analysisResults)

	// AIモデルにVLog生成を依頼（ツールを使用）
	_, err := genkit.Generate(ctx, g,
		ai.WithPrompt(prompt),
		ai.WithTools(registeredTools.GenerateVlogVideo),
	)
	if err != nil {
		// 直接ツールを呼び出す（フォールバック）
		resultRaw, err := registeredTools.GenerateVlogVideo.RunRaw(ctx, GenerateVlogVideoInput{
			AnalysisResults: analysisResults,
			Style:           input.Style,
			Title:           input.Title,
			MediaItems:      input.MediaItems,
		})
		if err != nil {
			return nil, err
		}
		result := resultRaw.(GenerateVlogVideoOutput)
		return &result, nil
	}

	// 生成されたビデオIDを生成
	videoID, _ := ulid.GenerateULID()

	return &GenerateVlogVideoOutput{
		VideoID:     videoID,
		VideoURL:    fmt.Sprintf("https://storage.example.com/videos/%s.mp4", videoID),
		Duration:    float64(input.Style.Duration),
		Title:       input.Title,
		Description: "",
	}, nil
}

// buildVlogGenerationPrompt はVLog生成用のプロンプトを構築する
func buildVlogGenerationPrompt(input *agent.VlogInput, analysisResults []agent.MediaAnalysisOutput) string {
	// 分析結果をサマリー化
	var locations, activities []string
	for _, r := range analysisResults {
		locations = append(locations, r.Landmarks...)
		activities = append(activities, r.Activities...)
	}

	prompt := fmt.Sprintf(`あなたは旅行VLog制作のエキスパートです。

以下の情報をもとに、感動的なVLogを生成してください。

## 旅行情報
- 旅行先: %s
- 旅行日: %s
- 訪問した場所: %v
- アクティビティ: %v
- メディア数: %d枚

## スタイル設定
- テーマ: %s
- BGMの雰囲気: %s
- 目標再生時間: %d秒
- トランジション: %s

## タスク
1. generateVlogVideoツールを使用してVLogを生成してください
2. 思い出に残る、エモいVLogを作成してください
3. 各シーンに適切な字幕を付けてください

ユーザーIDは「%s」です。`,
		input.Destination,
		input.TravelDate,
		locations,
		activities,
		len(input.MediaItems),
		input.Style.Theme,
		input.Style.MusicMood,
		input.Style.Duration,
		input.Style.Transition,
		input.UserID,
	)

	return prompt
}

// buildAnalyticsSummary は分析結果からサマリーを構築する
func buildAnalyticsSummary(results []agent.MediaAnalysisOutput, mediaCount int) agent.VlogAnalytics {
	locationsMap := make(map[string]struct{})
	activitiesMap := make(map[string]struct{})
	highlights := make([]string, 0)
	moodCounts := make(map[string]int)

	for _, r := range results {
		for _, loc := range r.Landmarks {
			locationsMap[loc] = struct{}{}
		}
		for _, act := range r.Activities {
			activitiesMap[act] = struct{}{}
		}
		if r.SuggestedCaption != "" {
			highlights = append(highlights, r.SuggestedCaption)
		}
		moodCounts[r.Mood]++
	}

	// マップをスライスに変換
	locations := make([]string, 0, len(locationsMap))
	for loc := range locationsMap {
		locations = append(locations, loc)
	}
	activities := make([]string, 0, len(activitiesMap))
	for act := range activitiesMap {
		activities = append(activities, act)
	}

	// 最も多いムードを選択
	overallMood := ""
	maxCount := 0
	for mood, count := range moodCounts {
		if count > maxCount {
			maxCount = count
			overallMood = mood
		}
	}

	return agent.VlogAnalytics{
		Locations:  locations,
		Activities: activities,
		Mood:       overallMood,
		Highlights: highlights,
		MediaCount: mediaCount,
	}
}
