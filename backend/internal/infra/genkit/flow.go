package genkit

import (
	"context"
	"fmt"

	"github.com/firebase/genkit/go/core"
	"github.com/firebase/genkit/go/genkit"
	"github.com/o-ga09/zenn-hackthon-2026/internal/agent"
	"github.com/o-ga09/zenn-hackthon-2026/internal/domain"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/errors"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/logger"
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
			thumbnailResult, _ = convertToStruct[GenerateThumbnailOutput](thumbnailRaw)
		}

		// Step 4: 共有URL生成
		shareRaw, err := registeredTools.GenerateShareURL.RunRaw(ctx, GenerateShareURLInput{
			VideoID: videoResult.VideoID,
			UserID:  input.UserID,
		})
		if err != nil {
			return nil, fmt.Errorf("share URL generation failed: %w", err)
		}
		shareResult, err := convertToStruct[GenerateShareURLOutput](shareRaw)
		if err != nil {
			return nil, fmt.Errorf("failed to parse share result: %w", err)
		}

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
	fc := GetFlowContext(ctx)
	results := make([]agent.MediaAnalysisOutput, 0, len(items))
	var allErrors []error

	for _, item := range items {
		resultRaw, analyzeErr := registeredTools.AnalyzeMedia.RunRaw(ctx, agent.MediaAnalysisInput{
			FileID:      item.FileID,
			URL:         item.URL,
			Type:        item.Type,
			ContentType: item.ContentType,
		})
		if analyzeErr != nil {
			allErrors = append(allErrors, analyzeErr)
			continue
		}

		// RunRawはmap[string]interface{}を返すのでJSONを経由して変換
		result, err := convertToStruct[agent.MediaAnalysisOutput](resultRaw)
		if err != nil {
			allErrors = append(allErrors, fmt.Errorf("failed to convert result: %w", err))
			continue
		}
		results = append(results, result)

		// DBに保存
		if fc != nil && fc.MediaAnalyticsRepo != nil {
			analytics := &domain.MediaAnalytics{
				FileID:      result.FileID,
				Description: result.Description,
				Objects:     make([]domain.DetectedObject, len(result.Objects)),
				Landmarks:   make([]domain.Landmark, len(result.Landmarks)),
				Activities:  make([]domain.Activity, len(result.Activities)),
				Mood:        result.Mood,
			}
			for i, o := range result.Objects {
				analytics.Objects[i] = domain.DetectedObject{Name: o}
			}
			for i, l := range result.Landmarks {
				analytics.Landmarks[i] = domain.Landmark{Name: l}
			}
			for i, a := range result.Activities {
				analytics.Activities[i] = domain.Activity{Name: a}
			}

			if err := fc.MediaAnalyticsRepo.Save(ctx, analytics); err != nil {
				// 保存失敗はログ出力のみで続行
				logger.Warn(ctx, fmt.Sprintf("failed to save media analytics for file %s: %v", result.FileID, err))
			}
		}
	}

	if len(results) == 0 {
		if len(allErrors) > 0 {
			return nil, errors.Join(allErrors...)
		}
		return nil, errors.ErrNoMediaItems
	}

	return results, nil
}

// generateVlog はAIモデルを使用してVLogを生成する
func generateVlog(ctx context.Context, g *genkit.Genkit, input *agent.VlogInput, analysisResults []agent.MediaAnalysisOutput, registeredTools *RegisteredTools) (*GenerateVlogVideoOutput, error) {
	// 直接ツールを呼び出してVeo3で動画生成
	resultRaw, err := registeredTools.GenerateVlogVideo.RunRaw(ctx, GenerateVlogVideoInput{
		AnalysisResults: analysisResults,
		Style:           input.Style,
		Title:           input.Title,
		MediaItems:      input.MediaItems,
		UserID:          input.UserID,
	})
	if err != nil {
		return nil, fmt.Errorf("vlog generation failed: %w", err)
	}

	// RunRawはmap[string]interface{}を返すのでJSONを経由して変換
	result, err := convertToStruct[GenerateVlogVideoOutput](resultRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to convert result: %w", err)
	}
	return &result, nil
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

// TODO: refactor は internal/infra/genkit/agent.go の convertToStruct を共通利用
