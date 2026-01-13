package handler

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/labstack/echo"
	"github.com/o-ga09/zenn-hackthon-2026/internal/domain"
	"github.com/o-ga09/zenn-hackthon-2026/internal/handler/request"
	"github.com/o-ga09/zenn-hackthon-2026/internal/handler/response"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/config"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/context"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/errors"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/ulid"
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

	env := config.GetCtxEnv(ctx)

	mediaResponses := make([]*response.MediaListItem, 0, len(medias))
	for _, media := range medias {
		mediaResponses = append(mediaResponses, &response.MediaListItem{
			ID:          media.ID,
			Type:        media.Type,
			ContentType: media.ContentType,
			Size:        media.Size,
			URL:         fmt.Sprintf("%s/%s/%s", env.CLOUDFLARE_R2_PUBLIC_URL, env.CLOUDFLARE_R2_BUCKET_NAME, media.URL),
			ImageData:   base64Images[media.URL],
			CreatedAt:   media.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return c.JSON(http.StatusOK, response.MediaListResponse{
		Media: mediaResponses,
		Total: len(mediaResponses),
	})
}

// multipart/form-dataでメディアアップロード(画像・動画対応)
func (s *ImageServer) Upload(c echo.Context) error {
	ctx := c.Request().Context()
	// TODO: リクエスト構造体にバインドできるようにする
	// multipart/form-dataからファイルを取得
	file, err := c.FormFile("file")
	if err != nil {
		return errors.Wrap(ctx, err)
	}

	// ファイルオープン
	src, err := file.Open()
	if err != nil {
		return errors.Wrap(ctx, err)
	}
	defer src.Close()

	// ファイルデータを読み取り
	fileData, err := io.ReadAll(src)
	if err != nil {
		return errors.Wrap(ctx, err)
	}

	// MIMEタイプを検出
	contentType := http.DetectContentType(fileData)

	// ファイル拡張子からもMIMEタイプを推定(動画の場合は必要)
	ext := strings.ToLower(filepath.Ext(file.Filename))
	switch ext {
	case ".mp4":
		contentType = "video/mp4"
	case ".mov":
		contentType = "video/quicktime"
	case ".avi":
		contentType = "video/x-msvideo"
	case ".webm":
		contentType = "video/webm"
	case ".mkv":
		contentType = "video/x-matroska"
	}

	// ファイル形式を判定
	var fileType string
	switch {
	case strings.HasPrefix(contentType, "image/"):
		fileType = "image"
	case strings.HasPrefix(contentType, "video/"):
		fileType = "video"
	default:
		return errors.Wrap(ctx, fmt.Errorf("対応していないファイル形式: %s", contentType))
	}

	// ユーザーIDとファイルIDを取得
	userID := context.GetCtxFromUser(ctx)
	fileID := ulid.New()

	// ストレージキーを生成
	storageKey := fmt.Sprintf("media/%s/%s%s", userID, fileID, ext)

	// ストレージにアップロード
	storageURL, err := s.storage.UploadFile(ctx, storageKey, fileData, contentType)
	if err != nil {
		return errors.Wrap(ctx, err)
	}

	// データベースにメタデータを保存
	model := &domain.Media{
		Type:        fileType,
		ContentType: contentType,
		Size:        file.Size,
		URL:         storageURL,
	}
	if err := s.imageRepo.Save(ctx, model); err != nil {
		return errors.Wrap(ctx, err)
	}

	// レスポンスを返す
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
