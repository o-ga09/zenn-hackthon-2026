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
	nullvalue "github.com/o-ga09/zenn-hackthon-2026/pkg/null_value"
)

type IAgentServer interface {
	CreateVLog(echo.Context) error
	AnalyzeMedia(echo.Context) error
	StreamAnalysisStatus(echo.Context) error // SSEエンドポイント
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
	notificationRepo   domain.INotificationRepository
}

func NewAgentServer(ctx context.Context, storage domain.IImageStorage, agentInstance agent.IAgent, vlogRepo domain.IVLogRepository, mediaRepo domain.IMediaRepository, mediaAnalyticsRepo domain.IMediaAnalyticsRepository, taskClient queue.IQueue, txManager domain.ITransactionManager, notificationRepo domain.INotificationRepository) *AgentServer {
	return &AgentServer{
		storage:            storage,
		agent:              agentInstance,
		vlogRepo:           vlogRepo,
		mediaRepo:          mediaRepo,
		mediaAnalyticsRepo: mediaAnalyticsRepo,
		taskClient:         taskClient,
		txManager:          txManager,
		notificationRepo:   notificationRepo,
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

	// TODO: 構造体バインド対応
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
				URL:         media.URL.String,
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

		// VLog生成失敗の通知を作成
		if vlogRef.CreateUserID != nil {
			notification := &domain.Notification{
				UserID:  *vlogRef.CreateUserID,
				Type:    domain.NotificationTypeVlogFailed,
				Title:   "VLog生成失敗",
				Message: fmt.Sprintf("VLogの生成に失敗しました: %s", err.Error()),
				VlogID:  nullvalue.ToNullString(task.ID),
				Read:    false,
			}
			if notifErr := s.notificationRepo.Create(ctx, notification); notifErr != nil {
				fmt.Printf("[executeVLogGeneration] Failed to create notification: %v\n", notifErr)
			}
		}

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

	// VLog生成完了の通知を作成
	if vlogRef.CreateUserID != nil {
		notification := &domain.Notification{
			UserID:  *vlogRef.CreateUserID,
			Type:    domain.NotificationTypeVlogCompleted,
			Title:   "VLog生成完了",
			Message: "VLogの生成が完了しました",
			VlogID:  nullvalue.ToNullString(task.ID),
			Read:    false,
		}
		if notifErr := s.notificationRepo.Create(ctx, notification); notifErr != nil {
			fmt.Printf("[executeVLogGeneration] Failed to create notification: %v\n", notifErr)
		}
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

		// MediaレコードをDBに保存
		media := &domain.Media{
			BaseModel: domain.BaseModel{
				CreateUserID: &userID,
			},
			ContentType: contentType,
			Size:        int64(len(data)),
			URL:         nullvalue.ToNullString(key),
		}
		if err := s.mediaRepo.Save(ctx, media); err != nil {
			return nil, fmt.Errorf("failed to save media record: %w", err)
		}

		// ストレージにアップロード
		path := key + media.ID + ext
		url, err := s.storage.UploadFile(ctx, path, data, contentType)
		if err != nil {
			return nil, fmt.Errorf("failed to upload file %s: %w", fileHeader.Filename, err)
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

// TODO: 以下ユーティリティ関数は共通化検討
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

// TODO: 以下ユーティリティ関数は共通化検討
// detectMediaType はコンテンツタイプからメディアタイプ（image/video）を判定する
func detectMediaType(contentType string) string {
	if strings.HasPrefix(contentType, "video/") {
		return "video"
	}
	return "image"
}

// TODO: 以下ユーティリティ関数は共通化検討
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

// AnalyzeMediaResponse はメディア分析APIのレスポンス
type AnalyzeMediaResponse struct {
	MediaIDs []string `json:"media_ids"`
	Status   string   `json:"status"`
}

// AnalyzeMedia はメディアを分析する（非同期処理版）
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

	// 各ファイルのMediaレコードをPENDING状態で作成
	mediaIDs := make([]string, len(files))
	for i, fileHeader := range files {
		contentType := fileHeader.Header.Get("Content-Type")
		media := &domain.Media{
			BaseModel: domain.BaseModel{
				CreateUserID: &userIDStr,
			},
			ContentType: contentType,
			Size:        fileHeader.Size,
			Status:      domain.MediaStatusPending,
			Progress:    0.0,
		}
		if err := s.mediaRepo.Save(ctx, media); err != nil {
			return errors.Wrap(ctx, err)
		}
		mediaIDs[i] = media.ID
		fmt.Printf("[AnalyzeMedia] Created media record: ID=%s, status=%s\n", media.ID, media.Status)
	}

	fmt.Printf("[AnalyzeMedia] Starting async analysis for %d files\n", len(files))

	// 非同期で分析処理を実行（新しいバックグラウンドコンテキストを使用）
	// Gorutineでの実行は、ローカルのみ
	bgCtx := context.Background()
	bgCtx = Ctx.SetCtxFromUser(bgCtx, userIDStr)
	bgCtx = Ctx.SetRequestTime(bgCtx, time.Now())
	bgCtx = Ctx.SetDB(bgCtx, Ctx.GetDB(ctx))
	go s.processMediaAnalysis(bgCtx, userIDStr, mediaIDs, files)

	// 即座にmediaIDリストを返却
	return c.JSON(http.StatusOK, AnalyzeMediaResponse{
		MediaIDs: mediaIDs,
		Status:   string(domain.MediaStatusPending),
	})
}

// processMediaAnalysis は非同期でメディア分析を処理する
func (s *AgentServer) processMediaAnalysis(ctx context.Context, userID string, mediaIDs []string, files []*multipart.FileHeader) {
	fmt.Printf("[processMediaAnalysis] Started for user=%s, files=%d\n", userID, len(files))

	totalFiles := len(files)
	mediaItems := make([]agent.MediaItem, 0, totalFiles)

	// Phase 1: アップロード（進捗 0.0 → 0.5）
	for i, fileHeader := range files {
		mediaID := mediaIDs[i]
		fmt.Printf("[processMediaAnalysis] Processing file %d/%d, mediaID=%s\n", i+1, totalFiles, mediaID)

		media, err := s.mediaRepo.GetByID(ctx, mediaID)
		if err != nil {
			fmt.Printf("[processMediaAnalysis] Failed to get media %s: %v\n", mediaID, err)
			continue
		}

		// Status: UPLOADING
		media.Status = domain.MediaStatusUploading
		media.Progress = 0.1
		if err := s.mediaRepo.Save(ctx, media); err != nil {
			fmt.Printf("[processMediaAnalysis] Failed to update status to UPLOADING: %v\n", err)
		}
		fmt.Printf("[processMediaAnalysis] Updated status to UPLOADING for %s\n", mediaID)

		// ファイルを開いてアップロード
		file, err := fileHeader.Open()
		if err != nil {
			fmt.Printf("[processMediaAnalysis] Failed to open file: %v\n", err)
			media.Status = domain.MediaStatusFailed
			media.ErrorMessage = fmt.Sprintf("Failed to open file: %v", err)
			s.mediaRepo.Save(ctx, media)
			continue
		}

		data, err := io.ReadAll(file)
		file.Close()
		if err != nil {
			fmt.Printf("[processMediaAnalysis] Failed to read file: %v\n", err)
			media.Status = domain.MediaStatusFailed
			media.ErrorMessage = fmt.Sprintf("Failed to read file: %v", err)
			s.mediaRepo.Save(ctx, media)
			continue
		}

		contentType := fileHeader.Header.Get("Content-Type")
		if contentType == "" {
			contentType = detectContentType(fileHeader.Filename, data)
		}

		ext := filepath.Ext(fileHeader.Filename)
		if ext == "" {
			ext = getExtensionFromContentType(contentType)
		}
		key := fmt.Sprintf("users/%s/uploads/%s%s", userID, mediaID, ext)

		fmt.Printf("[processMediaAnalysis] Uploading to storage: key=%s\n", key)
		url, err := s.storage.UploadFile(ctx, key, data, contentType)
		if err != nil {
			fmt.Printf("[processMediaAnalysis] Failed to upload: %v\n", err)
			media.Status = domain.MediaStatusFailed
			media.ErrorMessage = fmt.Sprintf("Failed to upload: %v", err)
			s.mediaRepo.Save(ctx, media)
			continue
		}

		fmt.Printf("[processMediaAnalysis] Upload successful: url=%s\n", url)

		// アップロード成功 - URLを更新
		media.URL = nullvalue.ToNullString(key)
		media.Progress = 0.5 // アップロード完了で50%
		if err := s.mediaRepo.Save(ctx, media); err != nil {
			fmt.Printf("[processMediaAnalysis] Failed to save media after upload: %v\n", err)
		}
		fmt.Printf("[processMediaAnalysis] Updated media with URL and progress=0.5\n")

		mediaItems = append(mediaItems, agent.MediaItem{
			FileID:      mediaID,
			URL:         url,
			Type:        detectMediaType(contentType),
			ContentType: contentType,
			Order:       i + 1,
		})
	}

	if len(mediaItems) == 0 {
		fmt.Println("[processMediaAnalysis] All uploads failed, aborting analysis")
		return // 全てのアップロードに失敗
	}

	fmt.Printf("[processMediaAnalysis] Starting analysis for %d items\n", len(mediaItems))

	// Phase 2: 分析（進捗 0.5 → 1.0）
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
	err := s.txManager.Do(ctx, func(ctx context.Context) error {
		result, err := s.agent.AnalyzeMediaBatch(ctx, input)
		if err != nil {
			fmt.Printf("[processMediaAnalysis] Analysis failed: %v\n", err)
		} else {
			fmt.Printf("[processMediaAnalysis] Analysis completed successfully: %d results\n", len(result.Results))
		}
		return err
	})

	// 結果を各メディアに反映
	for _, item := range mediaItems {
		media, getErr := s.mediaRepo.GetByID(ctx, item.FileID)
		if getErr != nil {
			fmt.Printf("[processMediaAnalysis] Failed to get media for final update: %v\n", getErr)
			continue
		}

		if err != nil {
			fmt.Printf("[processMediaAnalysis] Setting media %s to FAILED\n", item.FileID)
			media.Status = domain.MediaStatusFailed
			media.ErrorMessage = fmt.Sprintf("Analysis failed: %v", err)
			media.Progress = 0.5

			// メディア分析失敗の通知を作成
			notification := &domain.Notification{
				UserID:  userID,
				Type:    domain.NotificationTypeMediaFailed,
				Title:   "メディア分析失敗",
				Message: fmt.Sprintf("メディアの分析に失敗しました: %s", media.ErrorMessage),
				MediaID: nullvalue.ToNullString(media.ID),
				Read:    false,
			}
			if notifErr := s.notificationRepo.Create(ctx, notification); notifErr != nil {
				fmt.Printf("[processMediaAnalysis] Failed to create notification: %v\n", notifErr)
			}
		} else {
			fmt.Printf("[processMediaAnalysis] Setting media %s to COMPLETED\n", item.FileID)
			media.Status = domain.MediaStatusCompleted
			media.Progress = 1.0

			// メディア分析完了の通知を作成
			notification := &domain.Notification{
				UserID:  userID,
				Type:    domain.NotificationTypeMediaCompleted,
				Title:   "メディア分析完了",
				Message: "メディアの分析が完了しました",
				MediaID: nullvalue.ToNullString(media.ID),
				Read:    false,
			}
			if notifErr := s.notificationRepo.Create(ctx, notification); notifErr != nil {
				fmt.Printf("[processMediaAnalysis] Failed to create notification: %v\n", notifErr)
			}
		}
		if err := s.mediaRepo.Save(ctx, media); err != nil {
			fmt.Printf("[processMediaAnalysis] Failed to save final status: %v\n", err)
		}
	}

	fmt.Println("[processMediaAnalysis] Analysis processing completed")
}

// MediaStatusResponse はメディアステータスSSEのレスポンス
type MediaStatusResponse struct {
	Medias         []*domain.Media `json:"medias"`
	TotalItems     int             `json:"total_items"`
	CompletedItems int             `json:"completed_items"`
	FailedItems    int             `json:"failed_items"`
	AllCompleted   bool            `json:"all_completed"`
}

// StreamAnalysisStatus はメディア分析の進捗をSSEでストリーミングする
// GET /api/agent/analyze-media/stream?ids=id1,id2,id3
func (s *AgentServer) StreamAnalysisStatus(c echo.Context) error {
	ctx := c.Request().Context()
	idsParam := c.QueryParam("ids")

	fmt.Printf("[SSE] リクエスト受信: ids=%s\n", idsParam)

	if idsParam == "" {
		fmt.Println("[SSE] エラー: ids パラメータがありません")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "ids query parameter is required",
		})
	}
	mediaIDs := strings.Split(idsParam, ",")
	fmt.Printf("[SSE] メディアID数: %d, IDs: %v\n", len(mediaIDs), mediaIDs)

	// SSEヘッダー設定
	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().Header().Set("X-Accel-Buffering", "no")

	fmt.Println("[SSE] ヘッダー設定完了、ストリーミング開始")

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	iterationCount := 0
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("[SSE] コンテキストキャンセル (iteration: %d)\n", iterationCount)
			return nil
		case <-ticker.C:
			iterationCount++
			fmt.Printf("[SSE] Tick %d: メディアステータスを取得中...\n", iterationCount)

			medias := make([]*domain.Media, 0, len(mediaIDs))
			completedCount := 0
			failedCount := 0

			for _, id := range mediaIDs {
				media, err := s.mediaRepo.GetByID(ctx, strings.TrimSpace(id))
				if err != nil {
					fmt.Printf("[SSE] メディア取得エラー (ID: %s): %v\n", id, err)
					continue
				}
				medias = append(medias, media)
				fmt.Printf("[SSE] メディア %s: status=%s, progress=%.2f\n", media.ID, media.Status, media.Progress)

				if media.Status == domain.MediaStatusCompleted {
					completedCount++
				} else if media.Status == domain.MediaStatusFailed {
					failedCount++
				}
			}

			allCompleted := (completedCount + failedCount) == len(mediaIDs)
			fmt.Printf("[SSE] 進捗: completed=%d, failed=%d, total=%d, allCompleted=%v\n",
				completedCount, failedCount, len(mediaIDs), allCompleted)

			response := MediaStatusResponse{
				Medias:         medias,
				TotalItems:     len(mediaIDs),
				CompletedItems: completedCount,
				FailedItems:    failedCount,
				AllCompleted:   allCompleted,
			}

			// JSON送信
			data, err := json.Marshal(response)
			if err != nil {
				fmt.Printf("[SSE] JSONマーシャルエラー: %v\n", err)
				return err
			}

			fmt.Printf("[SSE] データ送信: %s\n", string(data))
			_, err = fmt.Fprintf(c.Response(), "data: %s\n\n", data)
			if err != nil {
				fmt.Printf("[SSE] 書き込みエラー: %v\n", err)
				return err
			}
			c.Response().Flush()
			fmt.Println("[SSE] フラッシュ完了")

			// 全て完了で終了
			if allCompleted {
				fmt.Println("[SSE] すべての分析が完了、終了イベントを送信します")

				// 明示的な終了イベントを送信（フロントエンドで確実に完了を検知するため）
				_, err = fmt.Fprintf(c.Response(), "event: complete\ndata: {\"status\":\"done\"}\n\n")
				if err != nil {
					fmt.Printf("[SSE] 終了イベント書き込みエラー: %v\n", err)
				} else {
					c.Response().Flush()
					fmt.Println("[SSE] 終了イベント送信完了")
				}

				// クライアントがメッセージを処理する時間を確保
				time.Sleep(100 * time.Millisecond)

				fmt.Println("[SSE] 接続を正常終了します")
				// 強制クローズはせず、return nil で自然に接続を閉じる
				return nil
			}
		}
	}
}
