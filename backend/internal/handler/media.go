package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/o-ga09/zenn-hackthon-2026/internal/domain"
	"github.com/o-ga09/zenn-hackthon-2026/internal/handler/request"
	"github.com/o-ga09/zenn-hackthon-2026/internal/handler/response"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/config"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/context"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/date"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/errors"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/ptr"
	"gorm.io/gorm"
)

type IImageServer interface {
	List(c echo.Context) error
	GetByKey(c echo.Context) error
	Delete(c echo.Context) error
	GetAnalytics(c echo.Context) error
	UpdateAnalytics(c echo.Context) error
}

type ImageServer struct {
	imageRepo     domain.IMediaRepository
	storage       domain.IImageStorage
	analyticsRepo domain.IMediaAnalyticsRepository
}

func NewImageServer(imageRepo domain.IMediaRepository, storage domain.IImageStorage, analyticsRepo domain.IMediaAnalyticsRepository) *ImageServer {
	return &ImageServer{
		imageRepo:     imageRepo,
		storage:       storage,
		analyticsRepo: analyticsRepo,
	}
}

func (s *ImageServer) List(c echo.Context) error {
	ctx := c.Request().Context()
	medias, err := s.imageRepo.List(ctx, nil)
	if err != nil {
		return err
	}

	userId := context.GetCtxFromUser(ctx)
	prefix := fmt.Sprintf("media/%s/", userId)
	base64Images, err := s.storage.List(ctx, prefix)
	if err != nil {
		return err
	}

	env := config.GetCtxEnv(ctx)

	mediaResponses := make([]*response.MediaListItem, 0, len(medias))
	for _, media := range medias {
		var url string
		// NOTE: ローカル環境でフロントエンドで取得できるようにURLを置換
		if env.Env == "local" {
			url = strings.ReplaceAll(media.URL.String, "localstack", "localhost")
		}
		mediaResponses = append(mediaResponses, &response.MediaListItem{
			ID:          media.ID,
			ContentType: media.ContentType,
			Size:        media.Size,
			URL:         ptr.StringToPtr(url),
			Status:      string(media.Status),
			ImageData:   base64Images[media.URL.String],
			CreatedAt:   date.Format(media.CreatedAt),
		})
	}

	return c.JSON(http.StatusOK, response.MediaListResponse{
		Media: mediaResponses,
		Total: len(mediaResponses),
	})
}

func (s *ImageServer) GetByKey(c echo.Context) error {
	ctx := c.Request().Context()
	var req request.MediaGetRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	key, err := s.storage.Get(ctx, req.Key)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response.MediaGetResponse{
		ID:  req.Key,
		URL: key,
	})
}

func (s *ImageServer) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	var req request.MediaDeleteRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	err := s.storage.Delete(ctx, req.Key)
	if err != nil {
		return err
	}
	err = s.imageRepo.DeleteByFileID(ctx, &domain.Media{BaseModel: domain.BaseModel{ID: req.Key}})
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

// GetAnalytics メディアの分析結果を取得
func (s *ImageServer) GetAnalytics(c echo.Context) error {
	ctx := c.Request().Context()
	var req request.GetMediaAnalyticsParam
	if err := c.Bind(&req); err != nil {
		return errors.Wrap(ctx, err)
	}
	if err := c.Validate(&req); err != nil {
		return errors.Wrap(ctx, err)
	}

	analytics, err := s.analyticsRepo.FindByFileID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.MakeNotFoundError(ctx, "分析結果を取得できませんでした")
		}
		return errors.Wrap(ctx, err)
	}

	// ドメインモデルからレスポンスに変換
	objects := make([]string, len(analytics.Objects))
	for i, obj := range analytics.Objects {
		objects[i] = obj.Name
	}
	landmarks := make([]string, len(analytics.Landmarks))
	for i, landmark := range analytics.Landmarks {
		landmarks[i] = landmark.Name
	}
	activities := make([]string, len(analytics.Activities))
	for i, activity := range analytics.Activities {
		activities[i] = activity.Name
	}

	return c.JSON(http.StatusOK, response.MediaAnalyticsResponse{
		FileID:      analytics.FileID,
		Description: analytics.Description,
		Mood:        analytics.Mood,
		Objects:     objects,
		Landmarks:   landmarks,
		Activities:  activities,
	})
}

// UpdateAnalytics メディアの分析結果を更新
func (s *ImageServer) UpdateAnalytics(c echo.Context) error {
	ctx := c.Request().Context()
	var req request.UpdateMediaAnalyticsRequest
	if err := c.Bind(&req); err != nil {
		return errors.Wrap(ctx, err)
	}
	if err := c.Validate(&req); err != nil {
		return errors.Wrap(ctx, err)
	}

	// 既存の分析結果を取得
	analytics, err := s.analyticsRepo.FindByFileID(ctx, req.ID)
	if err != nil {
		return errors.Wrap(ctx, err)
	}

	// リクエストから更新
	if req.Description != nil {
		analytics.Description = *req.Description
	}
	if req.Mood != nil {
		analytics.Mood = *req.Mood
	}
	if req.Objects != nil {
		analytics.Objects = make([]domain.DetectedObject, len(req.Objects))
		for i, name := range req.Objects {
			analytics.Objects[i] = domain.DetectedObject{Name: name}
		}
	}
	if req.Landmarks != nil {
		analytics.Landmarks = make([]domain.Landmark, len(req.Landmarks))
		for i, name := range req.Landmarks {
			analytics.Landmarks[i] = domain.Landmark{Name: name}
		}
	}
	if req.Activities != nil {
		analytics.Activities = make([]domain.Activity, len(req.Activities))
		for i, name := range req.Activities {
			analytics.Activities[i] = domain.Activity{Name: name}
		}
	}

	// 更新を実行
	if err := s.analyticsRepo.Update(ctx, analytics); err != nil {
		return errors.Wrap(ctx, err)
	}

	// 更新後のレスポンスを返す
	objects := make([]string, len(analytics.Objects))
	for i, obj := range analytics.Objects {
		objects[i] = obj.Name
	}
	landmarks := make([]string, len(analytics.Landmarks))
	for i, landmark := range analytics.Landmarks {
		landmarks[i] = landmark.Name
	}
	activities := make([]string, len(analytics.Activities))
	for i, activity := range analytics.Activities {
		activities[i] = activity.Name
	}

	return c.JSON(http.StatusOK, response.MediaAnalyticsResponse{
		FileID:      analytics.FileID,
		Description: analytics.Description,
		Mood:        analytics.Mood,
		Objects:     objects,
		Landmarks:   landmarks,
		Activities:  activities,
	})
}
