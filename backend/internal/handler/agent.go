package handler

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/labstack/echo"
	"github.com/o-ga09/zenn-hackthon-2026/internal/agent"
	"github.com/o-ga09/zenn-hackthon-2026/internal/domain"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/errors"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/ulid"
)

type IAgentServer interface {
	CreateVLog(echo.Context) error
	AnalyzeMedia(echo.Context) error
}

type AgentServer struct {
	storage domain.IImageStorage
	agent   agent.IAgent
}

func NewAgentServer(ctx context.Context, storage domain.IImageStorage, agentInstance agent.IAgent) *AgentServer {
	return &AgentServer{
		storage: storage,
		agent:   agentInstance,
	}
}

// CreateVLogRequest はVLog生成APIのリクエスト（JSON形式の場合）
type CreateVLogRequest struct {
	MediaItems  []agent.MediaItem `json:"mediaItems" validate:"required,min=1"`
	Title       string            `json:"title,omitempty"`
	TravelDate  string            `json:"travelDate,omitempty"`
	Destination string            `json:"destination,omitempty"`
	Style       agent.VlogStyle   `json:"style,omitempty"`
}

// CreateVLog はメディアからVLogを生成する
// POST /api/agent/create-vlog
// Content-Type: multipart/form-data
//
// フォームフィールド:
//   - files: メディアファイル（複数可）
//   - title: VLogのタイトル（任意）
//   - travelDate: 旅行日（YYYY-MM-DD形式、任意）
//   - destination: 旅行先（任意）
//   - theme: テーマ（adventure/relaxing/romantic/family、任意）
//   - musicMood: BGMの雰囲気（任意）
//   - duration: 目標再生時間（秒、任意、デフォルト60）
//   - transition: トランジション効果（fade/slide/zoom、任意）
func (s *AgentServer) CreateVLog(c echo.Context) error {
	ctx := c.Request().Context()

	// ユーザーIDをコンテキストから取得
	userID := c.Get("userID")
	if userID == nil {
		userID = "anonymous"
	}
	userIDStr := userID.(string)

	// マルチパートフォームをパース
	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid multipart form",
		})
	}

	// ファイルを取得
	files := form.File["files"]
	if len(files) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "No files uploaded. Please upload at least one media file.",
		})
	}

	// ファイルをストレージにアップロードしてMediaItemsを構築
	mediaItems, err := s.uploadMediaFiles(ctx, userIDStr, files)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to upload files: %v", err),
		})
	}

	// フォームフィールドを取得
	title := c.FormValue("title")
	travelDate := c.FormValue("travelDate")
	destination := c.FormValue("destination")

	// スタイル設定を取得
	style := agent.VlogStyle{
		Theme:      c.FormValue("theme"),
		MusicMood:  c.FormValue("musicMood"),
		Duration:   60, // デフォルト
		Transition: c.FormValue("transition"),
	}
	if durationStr := c.FormValue("duration"); durationStr != "" {
		var duration int
		if _, err := fmt.Sscanf(durationStr, "%d", &duration); err == nil && duration > 0 {
			style.Duration = duration
		}
	}

	// 入力を構築
	input := &agent.VlogInput{
		UserID:      userIDStr,
		MediaItems:  mediaItems,
		Title:       title,
		TravelDate:  travelDate,
		Destination: destination,
		Style:       style,
	}

	// VLog生成を実行
	res, err := s.agent.CreateVlog(ctx, input)
	if err != nil {
		return errors.Wrap(ctx, err)
	}

	return c.JSON(http.StatusOK, res)
}

// uploadMediaFiles はマルチパートファイルをストレージにアップロードしてMediaItemsを返す
func (s *AgentServer) uploadMediaFiles(ctx context.Context, userID string, files []*multipart.FileHeader) ([]agent.MediaItem, error) {
	mediaItems := make([]agent.MediaItem, 0, len(files))

	for i, fileHeader := range files {
		// ファイルを開く
		file, err := fileHeader.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open file %s: %w", fileHeader.Filename, err)
		}
		defer file.Close()

		// ファイルデータを読み込む
		data, err := io.ReadAll(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %w", fileHeader.Filename, err)
		}

		// ファイルIDを生成
		fileID, err := ulid.GenerateULID()
		if err != nil {
			return nil, fmt.Errorf("failed to generate file ID: %w", err)
		}

		// コンテンツタイプを取得
		contentType := fileHeader.Header.Get("Content-Type")
		if contentType == "" {
			contentType = detectContentType(fileHeader.Filename, data)
		}

		// メディアタイプを判定
		mediaType := detectMediaType(contentType)

		// ストレージキーを生成
		ext := filepath.Ext(fileHeader.Filename)
		if ext == "" {
			ext = getExtensionFromContentType(contentType)
		}
		key := fmt.Sprintf("users/%s/uploads/%s%s", userID, fileID, ext)

		// ストレージにアップロード
		url, err := s.storage.UploadFile(ctx, key, data, contentType)
		if err != nil {
			return nil, fmt.Errorf("failed to upload file %s: %w", fileHeader.Filename, err)
		}

		// MediaItemを作成
		mediaItems = append(mediaItems, agent.MediaItem{
			FileID:      fileID,
			URL:         url,
			Type:        mediaType,
			ContentType: contentType,
			Order:       i + 1,
		})
	}

	return mediaItems, nil
}

// detectContentType はファイル名とデータからコンテンツタイプを検出する
func detectContentType(filename string, data []byte) string {
	// まずファイル拡張子から判定
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".heic", ".heif":
		return "image/heic"
	case ".mp4":
		return "video/mp4"
	case ".mov":
		return "video/quicktime"
	case ".avi":
		return "video/x-msvideo"
	case ".webm":
		return "video/webm"
	}

	// データからマジックナンバーで判定
	if len(data) > 0 {
		return http.DetectContentType(data)
	}

	return "application/octet-stream"
}

// detectMediaType はコンテンツタイプからメディアタイプ（image/video）を判定する
func detectMediaType(contentType string) string {
	if strings.HasPrefix(contentType, "video/") {
		return "video"
	}
	return "image"
}

// getExtensionFromContentType はコンテンツタイプから拡張子を取得する
func getExtensionFromContentType(contentType string) string {
	switch contentType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	case "image/heic":
		return ".heic"
	case "video/mp4":
		return ".mp4"
	case "video/quicktime":
		return ".mov"
	case "video/x-msvideo":
		return ".avi"
	case "video/webm":
		return ".webm"
	default:
		return ""
	}
}

// AnalyzeMediaRequest はメディア分析APIのリクエストボディ
type AnalyzeMediaRequest struct {
	FileID      string `json:"fileId" validate:"required"`
	URL         string `json:"url" validate:"required"`
	Type        string `json:"type" validate:"required,oneof=image video"`
	ContentType string `json:"contentType,omitempty"`
}

// AnalyzeMedia は単一メディアを分析する
// POST /api/agent/analyze-media
func (s *AgentServer) AnalyzeMedia(c echo.Context) error {
	ctx := c.Request().Context()

	var req AnalyzeMediaRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	input := &agent.MediaAnalysisInput{
		FileID:      req.FileID,
		URL:         req.URL,
		Type:        req.Type,
		ContentType: req.ContentType,
	}

	res, err := s.agent.AnalyzeMedia(ctx, input)
	if err != nil {
		return errors.Wrap(ctx, err)
	}

	return c.JSON(http.StatusOK, res)
}
