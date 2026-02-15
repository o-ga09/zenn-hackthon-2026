package genkit

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/genai"

	pkgStorage "github.com/o-ga09/zenn-hackthon-2026/internal/infra/storage"
	pkgConfig "github.com/o-ga09/zenn-hackthon-2026/pkg/config"
	pkgerrors "github.com/o-ga09/zenn-hackthon-2026/pkg/errors"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/ulid"
)

// VeoGenerateConfig はVeo動画生成の設定
type VeoGenerateConfig struct {
	Prompt          string
	DurationSeconds int32
	AspectRatio     string // "16:9" or "9:16"
	UserID          string
}

// VeoGenerateResult はVeo動画生成の結果
type VeoGenerateResult struct {
	VideoID  string
	VideoURL string // R2のURL
	Duration float64
}

// GenerateVideoWithVeo はVeo3を使用して動画を生成し、R2にアップロードする
func GenerateVideoWithVeo(ctx context.Context, fc *FlowContext, config VeoGenerateConfig) (*VeoGenerateResult, error) {
	if fc.GenAI == nil {
		return nil, fmt.Errorf("%w: GenAI client not initialized", pkgerrors.ErrGenkitNotInitialized)
	}
	if fc.GCSClient == nil {
		return nil, fmt.Errorf("%w: GCS client not initialized", pkgerrors.ErrGenkitNotInitialized)
	}
	if fc.Storage == nil {
		return nil, pkgerrors.ErrStorageNotInitialized
	}

	// 動画IDを生成
	videoID, err := ulid.GenerateULID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate video ID: %w", err)
	}

	// GCS一時出力パスを設定
	gcsOutputPath := fmt.Sprintf("gs://%s/temp/%s/", fc.Config.GCSTempBucket, videoID)

	// デフォルト値設定
	duration := config.DurationSeconds
	if duration == 0 {
		duration = 8
	}
	aspectRatio := config.AspectRatio
	if aspectRatio == "" {
		aspectRatio = "16:9"
	}

	// Veo動画生成オペレーションを開始
	op, err := fc.GenAI.Models.GenerateVideos(ctx,
		fc.Config.VeoModel,
		config.Prompt,
		nil, // 画像入力なし（テキストのみ）
		&genai.GenerateVideosConfig{
			DurationSeconds:  genai.Ptr(duration),
			AspectRatio:      aspectRatio,
			Resolution:       "720p",
			NumberOfVideos:   1,
			OutputGCSURI:     gcsOutputPath,
			GenerateAudio:    genai.Ptr(true),
			PersonGeneration: "allow_adult",
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start video generation: %w", err)
	}

	// オペレーション完了を待機（バックグラウンドコンテキストを使用してHTTPタイムアウトを回避）
	maxWait := time.Duration(fc.Config.VeoMaxWaitTime) * time.Second
	pollInterval := time.Duration(fc.Config.VeoPollingInterval) * time.Second
	startTime := time.Now()

	// Veoポーリング用に独立したコンテキストを使用
	veoCtx := context.Background()

	for !op.Done {
		if time.Since(startTime) > maxWait {
			return nil, fmt.Errorf("video generation timed out after %v", maxWait)
		}
		time.Sleep(pollInterval)
		op, err = fc.GenAI.Operations.GetVideosOperation(veoCtx, op, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get operation status: %w", err)
		}
	}

	// オペレーション全体をデバッグ出力
	fmt.Printf("[Veo] Operation Name: %s\n", op.Name)
	if op.Error != nil {
		return nil, fmt.Errorf("veo operation error: error=%v", op.Error)
	}

	// エラーチェック
	if op.Response == nil {
		return nil, fmt.Errorf("no response from Veo API")
	}

	if len(op.Response.GeneratedVideos) == 0 {
		// エラー詳細を出力
		return nil, fmt.Errorf("no video generated (response empty)")
	}

	generatedVideo := op.Response.GeneratedVideos[0]
	gcsVideoURI := generatedVideo.Video.URI

	// GCSから動画データを取得
	videoData, err := downloadFromGCS(ctx, fc.GCSClient, gcsVideoURI)
	if err != nil {
		return nil, fmt.Errorf("failed to download video from GCS: %w", err)
	}

	// R2にアップロード
	r2Key := fmt.Sprintf("users/%s/vlogs/%s.mp4", config.UserID, videoID)
	objectKey, err := fc.Storage.UploadFile(ctx, r2Key, videoData, "video/mp4")
	if err != nil {
		return nil, fmt.Errorf("failed to upload video to R2: %w", err)
	}

	// GCS一時ファイルを削除
	if err := deleteFromGCS(ctx, fc.GCSClient, gcsVideoURI); err != nil {
		// 削除失敗はログのみ（致命的ではない）
		fmt.Printf("warning: failed to delete temp file from GCS: %v\n", err)
	}

	env := pkgConfig.GetCtxEnv(ctx)
	url := pkgStorage.ObjectURKFromKey(env.CLOUDFLARE_R2_PUBLIC_URL, objectKey)
	if env.Env == "local" {
		url = fmt.Sprintf("%s/%s/%s", env.CLOUDFLARE_R2_PUBLIC_URL, env.CLOUDFLARE_R2_BUCKET_NAME, objectKey)
	}
	return &VeoGenerateResult{
		VideoID:  videoID,
		VideoURL: url,
		Duration: float64(duration),
	}, nil
}

// downloadFromGCS はGCSからファイルをダウンロードする
func downloadFromGCS(ctx context.Context, client *storage.Client, gcsURI string) ([]byte, error) {
	// gs://bucket/path/file.mp4 形式をパース
	bucket, object, err := parseGCSURI(gcsURI)
	if err != nil {
		return nil, err
	}

	reader, err := client.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCS reader: %w", err)
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read from GCS: %w", err)
	}

	return data, nil
}

// deleteFromGCS はGCSからファイルを削除する
func deleteFromGCS(ctx context.Context, client *storage.Client, gcsURI string) error {
	bucket, object, err := parseGCSURI(gcsURI)
	if err != nil {
		return err
	}

	return client.Bucket(bucket).Object(object).Delete(ctx)
}

// parseGCSURI はGCS URIをバケット名とオブジェクトパスに分解する
func parseGCSURI(gcsURI string) (bucket, object string, err error) {
	// gs://bucket/path/to/file.mp4
	if len(gcsURI) < 5 || gcsURI[:5] != "gs://" {
		return "", "", fmt.Errorf("invalid GCS URI: %s", gcsURI)
	}
	path := gcsURI[5:]
	idx := strings.Index(path, "/")
	if idx == -1 {
		return "", "", fmt.Errorf("invalid GCS URI format: %s", gcsURI)
	}
	return path[:idx], path[idx+1:], nil
}

// BuildVlogPrompt はVLog生成用のプロンプトを構築する
func BuildVlogPrompt(analysisResults []MediaAnalysisSummary, style VlogStyleConfig) string {
	// Veo向けに安全で効果的なプロンプトを構築
	prompt := "Create a beautiful travel vlog video. "

	// シーン情報を英語ベースで構築（安全性のため）
	if len(analysisResults) > 0 {
		// ランドマークと活動を収集
		var landmarks []string
		var activities []string
		var moods []string

		for _, result := range analysisResults {
			landmarks = append(landmarks, result.Landmarks...)
			activities = append(activities, result.Activities...)
			if result.Mood != "" {
				moods = append(moods, result.Mood)
			}
		}

		if len(landmarks) > 0 {
			// 最大3つのランドマークを使用
			if len(landmarks) > 3 {
				landmarks = landmarks[:3]
			}
			prompt += fmt.Sprintf("Feature these locations: %s. ", strings.Join(landmarks, ", "))
		}

		if len(activities) > 0 {
			// 最大3つのアクティビティを使用
			if len(activities) > 3 {
				activities = activities[:3]
			}
			prompt += fmt.Sprintf("Show activities like: %s. ", strings.Join(activities, ", "))
		}

		if len(moods) > 0 {
			prompt += fmt.Sprintf("The overall mood should be %s. ", moods[0])
		}
	}

	// スタイル設定（英語のみ）
	themeDescriptions := map[string]string{
		"adventure": "exciting and adventurous",
		"relaxing":  "calm and peaceful",
		"romantic":  "warm and romantic",
		"family":    "joyful and family-friendly",
	}
	if style.Theme != "" {
		if desc, ok := themeDescriptions[style.Theme]; ok {
			prompt += fmt.Sprintf("Style: %s. ", desc)
		}
	}

	prompt += "Use smooth camera movements, vibrant colors, and cinematic transitions. "
	prompt += "Professional travel documentary style with natural lighting."
	prompt += "Make it feel like a professional travel vlog that captures the essence of the journey."

	return prompt
}

// MediaAnalysisSummary はプロンプト生成用の分析サマリー
type MediaAnalysisSummary struct {
	Description string
	Landmarks   []string
	Activities  []string
	Mood        string
}

// VlogStyleConfig はVLogスタイル設定
type VlogStyleConfig struct {
	Theme      string
	MusicMood  string
	Duration   int
	Transition string
}
