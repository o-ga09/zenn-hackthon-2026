package handler

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/o-ga09/zenn-hackthon-2026/internal/domain"
	"github.com/o-ga09/zenn-hackthon-2026/internal/handler/request"
	"github.com/o-ga09/zenn-hackthon-2026/internal/handler/response"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/context"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/errors"
)

type IImageServer interface {
	List(c echo.Context) error
	Upload(c echo.Context) error
	GetByKey(c echo.Context) error
	Delete(c echo.Context) error
}

type ImageServer struct {
	imageRepo domain.IMediaRepository
	storage   domain.IImageStorage
}

func NewImageServer(imageRepo domain.IMediaRepository, storage domain.IImageStorage) *ImageServer {
	return &ImageServer{
		imageRepo: imageRepo,
		storage:   storage,
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

	mediaResponses := make([]*response.MediaListItem, 0, len(medias))
	for _, media := range medias {
		mediaResponses = append(mediaResponses, &response.MediaListItem{
			ID:          media.ID,
			Type:        media.Type,
			ContentType: media.ContentType,
			Size:        media.Size,
			URL:         media.URL,
			ImageData:   base64Images[media.URL],
			CreatedAt:   media.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return c.JSON(http.StatusOK, response.MediaListResponse{
		Media: mediaResponses,
		Total: len(mediaResponses),
	})
}

func (s *ImageServer) Upload(c echo.Context) error {
	ctx := c.Request().Context()
	var req request.MediaImageUploadRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	user := context.GetCtxFromUser(ctx)
	key, err := s.storage.Upload(ctx, fmt.Sprintf("media/%s/", user), req.Base64Data)
	if err != nil {
		return err
	}

	// base64DataからContentTypeやSizeを取得する
	// base64データをデコードしてバイナリデータを取得
	decodedData, err := base64.StdEncoding.DecodeString(req.Base64Data)
	if err != nil {
		return errors.Wrap(ctx, err)
	}

	// ファイルサイズを取得
	fileSize := int64(len(decodedData))

	// MIMEタイプを検出
	contentType := http.DetectContentType(decodedData)

	// ファイル形式を判定
	var fileType string
	switch contentType {
	case "image/jpeg":
		fileType = "jpeg"
	case "image/png":
		fileType = "png"
	case "image/gif":
		fileType = "gif"
	case "image/webp":
		fileType = "webp"
	default:
		fileType = "unknown"
	}

	model := &domain.Media{
		Type:        fileType,
		ContentType: contentType,
		Size:        fileSize,
		URL:         key,
	}
	if err := s.imageRepo.Save(ctx, model); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response.MediaImageUploadResponse{
		ID:  model.ID,
		URL: model.URL,
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
