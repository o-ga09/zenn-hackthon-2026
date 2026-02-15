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
	"github.com/o-ga09/zenn-hackthon-2026/internal/handler/request"
	"github.com/o-ga09/zenn-hackthon-2026/internal/handler/response"
	"github.com/o-ga09/zenn-hackthon-2026/internal/infra/storage"
	"github.com/o-ga09/zenn-hackthon-2026/internal/queue"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/config"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/constant"
	Ctx "github.com/o-ga09/zenn-hackthon-2026/pkg/context"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/errors"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/image"
	nullvalue "github.com/o-ga09/zenn-hackthon-2026/pkg/null_value"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/ptr"
)

type IAgentServer interface {
	CreateVLog(echo.Context) error
	AnalyzeMedia(echo.Context) error
	StreamAnalysisStatus(echo.Context) error
	ProcessVLogTask(echo.Context) error
	ProcessMediaAnalysisTask(echo.Context) error
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

// CreateVLog はメディアからVLogを生成する
func (s *AgentServer) CreateVLog(c echo.Context) error {
	ctx := c.Request().Context()
	// ユーザーIDをコンテキストから取得
	userIDStr := Ctx.GetCtxFromUser(ctx)
	if userIDStr == "" {
		userIDStr = "anonymous"
	}

	var req request.CreateVLogRequest
	if err := c.Bind(&req); err != nil {
		return errors.Wrap(ctx, err)
	}

	if err := c.Validate(&req); err != nil {
		return errors.Wrap(ctx, err)
	}

	// ファイルを取得
	files := req.Files

	if len(files) == 0 && len(req.MediaIDs) == 0 {
		return errors.MakeBusinessError(ctx, "新規メディアも既存メディも指定されていないため、新しいVlogを生成できません")
	}

	var mediaItems []agent.MediaItem

	// 1. 新規ファイルをアップロード
	if len(files) > 0 {
		items, err := s.uploadMediaFiles(ctx, userIDStr, files)
		if err != nil {
			return errors.Wrap(ctx, err)
		}
		mediaItems = append(mediaItems, items...)
	}

	// 2. 既存メディアを取得
	if len(req.MediaIDs) > 0 {
		for _, id := range req.MediaIDs {
			media, err := s.mediaRepo.GetByID(ctx, id)
			if err != nil {
				continue
			}
			mediaItems = append(mediaItems, agent.MediaItem{
				FileID:      media.ID,
				URL:         media.URL.String,
				ContentType: media.ContentType,
				Type:        detectMediaType(media.ContentType),
			})
		}
	}

	// スタイル設定を取得
	d := constant.DefaultVLogDurationSeconds
	if req.Duration != nil && ptr.PtrToInt(req.Duration) < d {
		d = ptr.PtrToInt(req.Duration)
	}

	style := agent.VlogStyle{
		Theme:      ptr.PtrToString(req.Theme),
		MusicMood:  ptr.PtrToString(req.MusicMood),
		Duration:   d,
		Transition: ptr.PtrToString(req.Transition),
	}

	// 入力を構築
	input := &agent.VlogInput{
		UserID:      userIDStr,
		MediaItems:  mediaItems,
		Title:       ptr.PtrToString(req.Title),
		TravelDate:  ptr.PtrToString(req.TravelDate),
		Destination: ptr.PtrToString(req.Destination),
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
		ID:      vlog.ID,
		Version: vlog.Version,
		Type:    "ProcessVLogTask",
		Data:    input,
		Status:  domain.MediaStatusPending.String(),
	}

	cfg := config.GetCtxEnv(ctx)
	if cfg.Env == "local" {
		// 新しいバックグラウンドコンテキストを作成
		bgCtx := context.Background()
		bgCtx = Ctx.SetConfig(bgCtx, cfg)
		bgCtx = Ctx.SetCtxFromUser(bgCtx, userIDStr)
		bgCtx = Ctx.SetRequestTime(bgCtx, time.Now())
		bgCtx = Ctx.SetDB(bgCtx, Ctx.GetDB(ctx))

		go func() {
			// ローカル環境ではGoroutineで直接実行
			if err := s.executeVLogGeneration(bgCtx, payload); err != nil {
				// エラーはログに出力（DBなどのステータスはexecuteVLogGeneration内で更新済み）
				fmt.Printf("Local VLog generation failed: %v\n", err)
			}
		}()
	} else {
		if err := s.taskClient.Enqueue(ctx, payload); err != nil {
			return errors.Wrap(ctx, err)
		}
	}

	return c.JSON(http.StatusAccepted, response.CreateVLogResponse{
		VlogID: vlog.ID,
		Status: string(domain.VlogStatusProcessing),
	})
}

// ProcessVLogTask はCloud Tasksからのリクエストを受け取り、VLog生成を非同期に実行する
func (s *AgentServer) ProcessVLogTask(c echo.Context) error {
	ctx := c.Request().Context()

	var task queue.Task
	if err := c.Bind(&task); err != nil {
		return errors.Wrap(ctx, err)
	}

	if err := s.executeVLogGeneration(ctx, &task); err != nil {
		return errors.Wrap(ctx, err)
	}

	return c.NoContent(http.StatusNoContent)
}

// executeVLogGeneration はVLog生成のコアロジックを実行する
func (s *AgentServer) executeVLogGeneration(ctx context.Context, task *queue.Task) error {
	// タスクデータから*agent.VlogInputを取得
	vlogInput, ok := task.Data.(*agent.VlogInput)
	if !ok {
		return errors.MakeBusinessError(ctx, "invalid task data type for VLog generation: expected *agent.VlogInput")
	}

	// 最新のVlogレコードを取得
	vlogRef, err := s.vlogRepo.GetByID(ctx, &domain.Vlog{BaseModel: domain.BaseModel{ID: task.ID}})
	if err != nil {
		return errors.Wrap(ctx, err)
	}

	// ステータスをPROCESSINGに更新
	now := time.Now()
	vlogRef.Status = domain.VlogStatusProcessing
	vlogRef.ErrorMessage = ""
	vlogRef.Progress = 0.1
	vlogRef.StartedAt = &now
	if err := s.vlogRepo.Update(ctx, vlogRef); err != nil {
		return errors.Wrap(ctx, err)
	}

	// VLog生成を実行
	var res *agent.VlogOutput
	err = s.txManager.Do(ctx, func(ctx context.Context) error {
		var err error
		res, err = s.agent.CreateVlogWithProgress(ctx, vlogInput, func(p agent.FlowProgress) {
			// 進捗をDBに更新（最新のレコードを取得してから更新）
			latestVlog, getErr := s.vlogRepo.GetByID(ctx, &domain.Vlog{BaseModel: domain.BaseModel{ID: vlogRef.ID}})
			if getErr == nil {
				latestVlog.Status = domain.VlogStatusProcessing
				latestVlog.ErrorMessage = ""
				latestVlog.Progress = p.Progress
				_ = s.vlogRepo.Update(ctx, latestVlog)
			}
		})
		return errors.Wrap(ctx, err)
	})

	if err != nil {
		// 最新のレコードを取得してから失敗ステータスに更新
		latestVlog, getErr := s.vlogRepo.GetByID(ctx, &domain.Vlog{BaseModel: domain.BaseModel{ID: vlogRef.ID}})
		if getErr == nil {
			latestVlog.Status = domain.VlogStatusFailed
			latestVlog.ErrorMessage = err.Error()
			latestVlog.Progress = 0
			_ = s.vlogRepo.Update(ctx, latestVlog)
		}

		// VLog生成失敗の通知を作成
		if vlogRef.CreateUserID != nil {
			notification := &domain.Notification{
				UserID:  *vlogRef.CreateUserID,
				Type:    domain.NotificationTypeVlogFailed,
				Title:   "VLog生成失敗",
				Message: fmt.Sprintf("VLogの生成に失敗しました: %s", err.Error()),
				VlogID:  nullvalue.ToNullString(vlogRef.ID),
				Read:    false,
			}
			if notifErr := s.notificationRepo.Create(ctx, notification); notifErr != nil {
				fmt.Printf("[executeVLogGeneration] Failed to create notification: %v\n", notifErr)
			}
		}

		return errors.Wrap(ctx, err)
	}

	// 成功ステータスと生成された情報を更新
	// 最新のレコードを取得
	latestVlog, err := s.vlogRepo.GetByID(ctx, &domain.Vlog{BaseModel: domain.BaseModel{ID: vlogRef.ID}})
	if err != nil {
		return errors.Wrap(ctx, err)
	}

	latestVlog.VideoID = res.VideoID
	latestVlog.VideoURL = res.VideoURL
	latestVlog.ShareURL = res.ShareURL
	latestVlog.Duration = res.Duration
	latestVlog.Thumbnail = res.ThumbnailURL
	latestVlog.Status = domain.VlogStatusCompleted
	latestVlog.Progress = 1.0
	completedAt := time.Now()
	latestVlog.CompletedAt = &completedAt

	if err := s.vlogRepo.Update(ctx, latestVlog); err != nil {
		return errors.Wrap(ctx, err)
	}

	// VLog生成完了の通知を作成
	if latestVlog.CreateUserID != nil {
		notification := &domain.Notification{
			UserID:  *latestVlog.CreateUserID,
			Type:    domain.NotificationTypeVlogCompleted,
			Title:   "VLog生成完了",
			Message: "VLogの生成が完了しました",
			VlogID:  nullvalue.ToNullString(latestVlog.ID),
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
			contentType = image.DetectContentType(data)
		}

		// メディアタイプを判定
		mediaType := detectMediaType(contentType)

		// ストレージキーを生成
		ext := filepath.Ext(fileHeader.Filename)
		if ext == "" {
			ext = image.GetExtensionFromContentType(contentType)
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

		env := config.GetCtxEnv(ctx)
		// ストレージにアップロード
		path := key + media.ID + ext
		objectKey, err := s.storage.UploadFile(ctx, path, data, contentType)
		if err != nil {
			return nil, fmt.Errorf("failed to upload file %s: %w", fileHeader.Filename, err)
		}

		url := storage.ObjectURKFromKey(env.CLOUDFLARE_R2_PUBLIC_URL, objectKey)
		if env.Env == "local" {
			url = fmt.Sprintf("%s/%s/%s", "http://localstack:4566", env.CLOUDFLARE_R2_BUCKET_NAME, objectKey)
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

// detectMediaType はコンテンツタイプからメディアタイプ（image/video）を判定する
func detectMediaType(contentType string) string {
	if strings.HasPrefix(contentType, "video/") {
		return "video"
	}
	return "image"
}

// AnalyzeMedia はメディアを分析する（非同期処理版）
func (s *AgentServer) AnalyzeMedia(c echo.Context) error {
	ctx := c.Request().Context()
	userIDStr := Ctx.GetCtxFromUser(ctx)
	if userIDStr == "" {
		userIDStr = "anonymous"
	}

	var req request.AnalyzeMediaRequest
	if err := c.Bind(&req); err != nil {
		return errors.Wrap(ctx, err)
	}

	if err := c.Validate(&req); err != nil {
		return errors.Wrap(ctx, err)
	}

	// ファイルを取得
	// 各ファイルのMediaレコードをPENDING状態で作成し、先にアップロードを実行
	files := req.Files
	mediaIDs := make([]string, 0, len(files))

	for i, fileHeader := range files {
		contentType := fileHeader.Header.Get("Content-Type")
		if contentType == "" {
			file, err := fileHeader.Open()
			if err != nil {
				continue
			}
			data, err := io.ReadAll(file)
			file.Close()
			if err != nil {
				continue
			}
			contentType = image.DetectContentType(data)
		}

		media := &domain.Media{
			BaseModel: domain.BaseModel{
				CreateUserID: &userIDStr,
			},
			ContentType: contentType,
			Size:        fileHeader.Size,
			Status:      domain.MediaStatusUploading,
			Progress:    0.0,
		}
		if err := s.mediaRepo.Save(ctx, media); err != nil {
			continue
		}

		// ファイルをアップロード
		file, err := fileHeader.Open()
		if err != nil {
			media.Status = domain.MediaStatusFailed
			media.ErrorMessage = fmt.Sprintf("Failed to open file: %v", err)
			s.mediaRepo.Save(ctx, media)
			continue
		}

		data, err := io.ReadAll(file)
		file.Close()
		if err != nil {
			media.Status = domain.MediaStatusFailed
			media.ErrorMessage = fmt.Sprintf("Failed to read file: %v", err)
			s.mediaRepo.Save(ctx, media)
			continue
		}

		ext := filepath.Ext(fileHeader.Filename)
		if ext == "" {
			ext = image.GetExtensionFromContentType(contentType)
		}
		key := fmt.Sprintf("users/%s/uploads/%s%s", userIDStr, media.ID, ext)

		objectKey, err := s.storage.UploadFile(ctx, key, data, contentType)
		if err != nil {
			media.Status = domain.MediaStatusFailed
			media.ErrorMessage = fmt.Sprintf("Failed to upload: %v", err)
			s.mediaRepo.Save(ctx, media)
			continue
		}

		// アップロード成功 - URLを更新
		env := config.GetCtxEnv(ctx)
		url := storage.ObjectURKFromKey(env.CLOUDFLARE_R2_PUBLIC_URL, objectKey)
		if env.Env == "local" {
			url = fmt.Sprintf("%s/%s/%s", "http://localstack:4566", env.CLOUDFLARE_R2_BUCKET_NAME, key)
		}

		media.URL = nullvalue.ToNullString(url)
		media.Progress = 0.5                     // アップロード完了で50%
		media.Status = domain.MediaStatusPending // 分析待ちに戻す
		if err := s.mediaRepo.Save(ctx, media); err != nil {
			continue
		}

		mediaIDs = append(mediaIDs, media.ID)
		_ = i // 未使用変数の警告を回避
	}

	if len(mediaIDs) == 0 {
		return errors.MakeBusinessError(ctx, "No media files were successfully uploaded")
	}

	// Cloud Tasksにタスクを登録
	payload := &queue.Task{
		Type: "ProcessMediaAnalysisTask",
		Data: map[string]interface{}{
			"user_id":   userIDStr,
			"media_ids": mediaIDs,
		},
		Status: domain.MediaStatusPending.String(),
	}

	env := config.GetCtxEnv(ctx)
	if env.Env == "local" {
		// ローカル環境ではGoroutineで直接実行
		bgCtx := context.Background()
		bgCtx = Ctx.SetConfig(bgCtx, env)
		bgCtx = Ctx.SetCtxFromUser(bgCtx, userIDStr)
		bgCtx = Ctx.SetRequestTime(bgCtx, time.Now())
		bgCtx = Ctx.SetDB(bgCtx, Ctx.GetDB(ctx))
		go func() {
			if err := s.executeMediaAnalysis(bgCtx, payload); err != nil {
				fmt.Printf("Media analysis failed: %v\n", err)
			}
		}()
	} else {
		if err := s.taskClient.Enqueue(ctx, payload); err != nil {
			return errors.Wrap(ctx, err)
		}
	}

	// 即座にmediaIDリストを返却
	return c.JSON(http.StatusOK, response.AnalyzeMediaResponse{
		MediaIDs: mediaIDs,
		Status:   string(domain.MediaStatusPending),
	})
}

// ProcessMediaAnalysisTask はCloud Tasksからのリクエストを受け取り、メディア分析を非同期に実行する
func (s *AgentServer) ProcessMediaAnalysisTask(c echo.Context) error {
	ctx := c.Request().Context()

	var task queue.Task
	if err := c.Bind(&task); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if err := s.executeMediaAnalysis(ctx, &task); err != nil {
		return errors.Wrap(ctx, err)
	}

	return c.NoContent(http.StatusOK)
}

// executeMediaAnalysis はメディア分析のコアロジックを実行する
func (s *AgentServer) executeMediaAnalysis(ctx context.Context, task *queue.Task) error {
	// タスクデータから必要な情報を取得
	data, ok := task.Data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid task data format")
	}

	userID, _ := data["user_id"].(string)
	mediaIDsInterface, _ := data["media_ids"].([]interface{})
	mediaIDs := make([]string, len(mediaIDsInterface))
	for i, v := range mediaIDsInterface {
		mediaIDs[i], _ = v.(string)
	}
	return s.processMediaAnalysisFromIDs(ctx, userID, mediaIDs)
}

// processMediaAnalysisFromIDs はメディアIDから分析を実行する
func (s *AgentServer) processMediaAnalysisFromIDs(ctx context.Context, userID string, mediaIDs []string) error {
	if len(mediaIDs) == 0 {
		return errors.MakeBusinessError(ctx, "No media IDs provided for analysis")
	}

	// メディアレコードを取得してMediaItemsを構築
	mediaItems := make([]agent.MediaAnalysisInput, 0, len(mediaIDs))
	for _, mediaID := range mediaIDs {
		media, err := s.mediaRepo.GetByID(ctx, mediaID)
		if err != nil {
			// メディアが見つからない場合はスキップ
			continue
		}

		if !media.URL.Valid {
			// URLが無効な場合は失敗ステータスに更新
			media.Status = domain.MediaStatusFailed
			media.ErrorMessage = "Media URL is invalid"
			s.mediaRepo.Save(ctx, media)
			continue
		}

		// Status: ANALYZING
		media.Status = domain.MediaStatusAnalyzing
		media.Progress = 0.6
		if err := s.mediaRepo.Save(ctx, media); err != nil {
			fmt.Printf("Failed to update media status: %v\n", err)
		}

		mediaType := detectMediaType(media.ContentType)
		mediaItems = append(mediaItems, agent.MediaAnalysisInput{
			FileID:      media.ID,
			URL:         media.URL.String,
			Type:        mediaType,
			ContentType: media.ContentType,
		})
	}

	if len(mediaItems) == 0 {
		return errors.MakeBusinessError(ctx, "No valid media items found for analysis")
	}

	// 分析を実行
	input := &agent.MediaAnalysisBatchInput{
		Items: mediaItems,
	}

	err := s.txManager.Do(ctx, func(ctx context.Context) error {
		_, err := s.agent.AnalyzeMediaBatch(ctx, input)
		if err != nil {
			return errors.Wrap(ctx, err)
		}
		return nil
	})

	// 結果を各メディアに反映
	for _, item := range mediaItems {
		media, getErr := s.mediaRepo.GetByID(ctx, item.FileID)
		if getErr != nil {
			continue
		}

		if err != nil {
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
				fmt.Printf("Failed to create notification: %v\n", notifErr)
			}
		} else {
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
				fmt.Printf("Failed to create notification: %v\n", notifErr)
			}
		}
		if err := s.mediaRepo.Save(ctx, media); err != nil {
			fmt.Printf("Failed to save media: %v\n", err)
		}
	}
	return nil
}

// StreamAnalysisStatus はメディア分析の進捗をSSEでストリーミングする
func (s *AgentServer) StreamAnalysisStatus(c echo.Context) error {
	ctx := c.Request().Context()
	var req request.AnalyzeMediaStreamRequest
	if err := c.Bind(&req); err != nil {
		return errors.Wrap(ctx, err)
	}

	if err := c.Validate(&req); err != nil {
		return errors.Wrap(ctx, err)
	}

	mediaIDs := make([]string, 0, len(req.IDs))
	for _, id := range req.IDs {
		mediaIDs = append(mediaIDs, id)
	}

	// SSEヘッダー設定
	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().Header().Set("X-Accel-Buffering", "no")

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	iterationCount := 0
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			iterationCount++

			medias := make([]*domain.Media, 0, len(mediaIDs))
			completedCount := 0
			failedCount := 0

			for _, id := range mediaIDs {
				media, err := s.mediaRepo.GetByID(ctx, strings.TrimSpace(id))
				if err != nil {
					continue
				}
				medias = append(medias, media)

				if media.Status == domain.MediaStatusCompleted {
					completedCount++
				} else if media.Status == domain.MediaStatusFailed {
					failedCount++
				}
			}

			allCompleted := (completedCount + failedCount) == len(mediaIDs)

			response := response.MediaStatusResponse{
				Medias:         medias,
				TotalItems:     len(mediaIDs),
				CompletedItems: completedCount,
				FailedItems:    failedCount,
				AllCompleted:   allCompleted,
			}

			// JSON送信
			data, err := json.Marshal(response)
			if err != nil {
				return errors.Wrap(ctx, err)
			}

			_, err = fmt.Fprintf(c.Response(), "data: %s\n\n", data)
			if err != nil {
				return errors.Wrap(ctx, err)
			}
			c.Response().Flush()

			// 全て完了で終了
			if allCompleted {
				// 明示的な終了イベントを送信
				_, err = fmt.Fprintf(c.Response(), "event: complete\ndata: {\"status\":\"done\"}\n\n")
				if err != nil {
					c.Response().Flush()
					return errors.Wrap(ctx, err)
				}

				// クライアントがメッセージを処理する時間を確保
				time.Sleep(100 * time.Millisecond)

				// 強制クローズはせず、return nil で自然に接続を閉じる
				return nil
			}
		}
	}
}
