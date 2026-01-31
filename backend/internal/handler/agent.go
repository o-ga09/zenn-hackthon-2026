package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/o-ga09/zenn-hackthon-2026/internal/agent"
	"github.com/o-ga09/zenn-hackthon-2026/internal/domain"
	"github.com/o-ga09/zenn-hackthon-2026/internal/queue"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/config"
	Ctx "github.com/o-ga09/zenn-hackthon-2026/pkg/context"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/errors"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/ulid"
)

type IAgentServer interface {
	CreateVLog(echo.Context) error
	AnalyzeMedia(echo.Context) error
	ProcessVLogTask(echo.Context) error
}

type AgentServer struct {
	storage            domain.IImageStorage
	agent              agent.IAgent
	vlogRepo           domain.IVLogRepository
	mediaRepo          domain.IMediaRepository
	mediaAnalyticsRepo domain.IMediaAnalyticsRepository
	taskClient         queue.IQueue
	txManager          domain.ITransactionManager
}

func NewAgentServer(ctx context.Context, storage domain.IImageStorage, agentInstance agent.IAgent, vlogRepo domain.IVLogRepository, mediaRepo domain.IMediaRepository, mediaAnalyticsRepo domain.IMediaAnalyticsRepository, taskClient queue.IQueue, txManager domain.ITransactionManager) *AgentServer {
	return &AgentServer{
		storage:            storage,
		agent:              agentInstance,
		vlogRepo:           vlogRepo,
		mediaRepo:          mediaRepo,
		mediaAnalyticsRepo: mediaAnalyticsRepo,
		taskClient:         taskClient,
		txManager:          txManager,
	}
}

type CreateVLogRequest struct {
	MediaItems  []agent.MediaItem `json:"mediaItems" validate:"required,min=1"`
	Title       string            `json:"title,omitempty"`
	TravelDate  string            `json:"travelDate,omitempty"`
	Destination string            `json:"destination,omitempty"`
	Style       agent.VlogStyle   `json:"style,omitempty"`
}

// CreateVLogResponse はVLog生成APIのレスポンス
type CreateVLogResponse struct {
	VlogID string `json:"vlogId"`
	Status string `json:"status"`
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
	userIDStr := Ctx.GetCtxFromUser(ctx)
	if userIDStr == "" {
		userIDStr = "anonymous"
	}

	// マルチパートフォームをパース
	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid multipart form",
		})
	}

	// ファイルを取得
	files := form.File["files"]

	// 既存メディアIDを取得
	var selectedMediaIds []string
	if mediaIdsStr := c.FormValue("mediaIds"); mediaIdsStr != "" {
		if err := json.Unmarshal([]byte(mediaIdsStr), &selectedMediaIds); err != nil {
			// fallback to comma separated if not JSON
			selectedMediaIds = strings.Split(mediaIdsStr, ",")
		}
	}

	if len(files) == 0 && len(selectedMediaIds) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "No files uploaded and no media IDs provided. Please provide at least one source.",
		})
	}

	var mediaItems []agent.MediaItem

	// 1. 新規ファイルをアップロード
	if len(files) > 0 {
		items, err := s.uploadMediaFiles(ctx, userIDStr, files)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": fmt.Sprintf("Failed to upload files: %v", err),
			})
		}
		mediaItems = append(mediaItems, items...)
	}

	// 2. 既存メディアを取得
	if len(selectedMediaIds) > 0 {
		for _, id := range selectedMediaIds {
			media, err := s.mediaRepo.GetByID(ctx, id)
			if err != nil {
				continue // またはエラーハンドリング
			}
			mediaItems = append(mediaItems, agent.MediaItem{
				FileID:      media.ID,
				URL:         media.URL,
				Type:        media.Type,
				ContentType: media.ContentType,
			})
		}
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

	// VLogレコードをPENDINGステータスで作成
	vlog := &domain.Vlog{
		Status: domain.VlogStatusPending,
	}
	if err := s.vlogRepo.Create(ctx, vlog); err != nil {
		return errors.Wrap(ctx, err)
	}

	// Cloud Tasksにタスクを登録
	payload := &queue.Task{
		ID:     vlog.ID,
		Type:   "ProcessVLogTask",
		Data:   input,
		Status: "pending",
	}

	cfg := config.GetCtxEnv(ctx)
	if cfg.Env == "local" {
		go func() {
			// ローカル環境ではGoroutineで直接実行
			if err := s.executeVLogGeneration(ctx, payload); err != nil {
				// エラーはログに出力（DBなどのステータスはexecuteVLogGeneration内で更新済み）
				fmt.Printf("Local VLog generation failed: %v\n", err)
			}
		}()
	} else {
		if err := s.taskClient.Enqueue(ctx, payload); err != nil {
			return errors.Wrap(ctx, err)
		}
	}

	return c.JSON(http.StatusAccepted, CreateVLogResponse{
		VlogID: vlog.ID,
		Status: string(domain.VlogStatusProcessing),
	})
}

// ProcessVLogTask はCloud Tasksからのリクエストを受け取り、VLog生成を非同期に実行する
// POST /internal/tasks/create-vlog
func (s *AgentServer) ProcessVLogTask(c echo.Context) error {
	ctx := c.Request().Context()

	var task queue.Task
	if err := c.Bind(&task); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if err := s.executeVLogGeneration(ctx, &task); err != nil {
		return errors.Wrap(ctx, err)
	}

	return c.NoContent(http.StatusOK)
}

// executeVLogGeneration はVLog生成のコアロジックを実行する
func (s *AgentServer) executeVLogGeneration(ctx context.Context, task *queue.Task) error {
	// ステータスをPROCESSINGに更新
	now := time.Now()
	if err := s.vlogRepo.UpdateStatus(ctx, task.ID, domain.VlogStatusProcessing, "", 0.1); err != nil {
		return err
	}

	vlogRef, err := s.vlogRepo.GetByID(ctx, &domain.Vlog{BaseModel: domain.BaseModel{ID: task.ID}})
	if err != nil {
		return err
	}
	vlogRef.StartedAt = &now
	_ = s.vlogRepo.Update(ctx, vlogRef)

	// VLog生成を実行
	var res *agent.VlogOutput
	err = s.txManager.Do(ctx, func(ctx context.Context) error {
		var err error
		res, err = s.agent.CreateVlogWithProgress(ctx, task.Data, func(p agent.FlowProgress) {
			// 進捗をDBに更新
			_ = s.vlogRepo.UpdateStatus(ctx, task.ID, domain.VlogStatusProcessing, "", p.Progress)
		})
		return err
	})

	if err != nil {
		// 失敗ステータスに更新
		_ = s.vlogRepo.UpdateStatus(ctx, task.ID, domain.VlogStatusFailed, err.Error(), 0)
		return err
	}

	// 成功ステータスと生成された情報を更新
	vlogRef.VideoID = res.VideoID
	vlogRef.VideoURL = res.VideoURL
	vlogRef.ShareURL = res.ShareURL
	vlogRef.Duration = res.Duration
	vlogRef.Thumbnail = res.ThumbnailURL
	vlogRef.Status = domain.VlogStatusCompleted
	vlogRef.Progress = 1.0
	completedAt := time.Now()
	vlogRef.CompletedAt = &completedAt

	if err := s.vlogRepo.Update(ctx, vlogRef); err != nil {
		return err
	}

	return nil
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
		key := fmt.Sprintf("users/%s/uploads/", userID)
		fileID, _ := ulid.GenerateULID()

		// ストレージにアップロード
		path := key + fileID + ext
		url, err := s.storage.UploadFile(ctx, path, data, contentType)
		if err != nil {
			return nil, fmt.Errorf("failed to upload file %s: %w", fileHeader.Filename, err)
		}

		// MediaレコードをDBに保存
		media := &domain.Media{
			BaseModel: domain.BaseModel{
				ID:           fileID,
				CreateUserID: &userID,
			},
			ContentType: contentType,
			Type:        mediaType,
			Size:        int64(len(data)),
			URL:         url,
		}
		if err := s.mediaRepo.Save(ctx, media); err != nil {
			return nil, fmt.Errorf("failed to save media record: %w", err)
		}

		// MediaItemを作成
		mediaItems = append(mediaItems, agent.MediaItem{
			FileID:      media.ID,
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

// AnalyzeMedia はメディアを分析する
// POST /api/agent/analyze-media
// Content-Type: multipart/form-data
//
// フォームフィールド:
//   - files: メディアファイル（複数可）
func (s *AgentServer) AnalyzeMedia(c echo.Context) error {
	ctx := c.Request().Context()
	// ユーザーIDをコンテキストから取得
	userIDStr := Ctx.GetCtxFromUser(ctx)
	if userIDStr == "" {
		userIDStr = "anonymous"
	}

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

	// メディア分析の入力を準備
	analysisItems := make([]agent.MediaAnalysisInput, len(mediaItems))
	for i, item := range mediaItems {
		analysisItems[i] = agent.MediaAnalysisInput{
			FileID:      item.FileID,
			URL:         item.URL,
			Type:        item.Type,
			ContentType: item.ContentType,
		}
	}

	input := &agent.MediaAnalysisBatchInput{
		Items: analysisItems,
	}

	// 分析を実行
	var res *agent.MediaAnalysisBatchOutput
	err = s.txManager.Do(ctx, func(ctx context.Context) error {
		var err error
		res, err = s.agent.AnalyzeMediaBatch(ctx, input)
		return err
	})
	if err != nil {
		return errors.Wrap(ctx, err)
	}

	return c.JSON(http.StatusOK, res)
}
