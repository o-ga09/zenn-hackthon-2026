package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/o-ga09/zenn-hackthon-2026/internal/domain"
	"github.com/o-ga09/zenn-hackthon-2026/internal/handler/request"
	"github.com/o-ga09/zenn-hackthon-2026/internal/handler/response"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/errors"
)

type IVLogServer interface {
	List(ctx echo.Context) error
	GetByID(ctx echo.Context) error
	Delete(ctx echo.Context) error
	StreamStatus(ctx echo.Context) error
}

type VLogServer struct {
	vlogRepo domain.IVLogRepository
}

func NewVLogServer(vlogRepo domain.IVLogRepository) *VLogServer {
	return &VLogServer{
		vlogRepo: vlogRepo,
	}
}

func (s *VLogServer) List(c echo.Context) error {
	ctx := c.Request().Context()
	var req request.VLogListRequest
	if err := c.Bind(&req); err != nil {
		return errors.Wrap(ctx, err)
	}
	if err := c.Validate(&req); err != nil {
		return errors.Wrap(ctx, err)
	}

	opts := &domain.ListOptions{
		Offset: req.Offset,
		Limit:  req.Limit,
	}
	vlogs, err := s.vlogRepo.List(ctx, opts)
	if err != nil {
		return errors.Wrap(ctx, err)
	}
	items := make([]response.VLogItem, 0, len(vlogs))
	for _, vlog := range vlogs {
		items = append(items, response.ToVLogItem(vlog))
	}
	res := response.VLogListResponse{
		Total: len(items),
		Items: items,
	}
	return c.JSON(http.StatusOK, res)

}

func (s *VLogServer) GetByID(c echo.Context) error {
	ctx := c.Request().Context()
	var req request.VLogGetByIDRequest
	if err := c.Bind(&req); err != nil {
		return errors.Wrap(ctx, err)
	}
	if err := c.Validate(&req); err != nil {
		return errors.Wrap(ctx, err)
	}

	model := &domain.Vlog{
		BaseModel: domain.BaseModel{
			ID: req.ID,
		},
	}
	vlog, err := s.vlogRepo.GetByID(ctx, model)
	if err != nil {
		return errors.Wrap(ctx, err)
	}
	res := response.ToVLogGetByIDResponse(vlog)
	return c.JSON(http.StatusOK, res)
}

func (s *VLogServer) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	var req request.VLogDeleteRequest
	if err := c.Bind(&req); err != nil {
		return errors.Wrap(ctx, err)
	}
	if err := c.Validate(&req); err != nil {
		return errors.Wrap(ctx, err)
	}

	model := &domain.Vlog{
		BaseModel: domain.BaseModel{
			ID: req.ID,
		},
	}
	if err := s.vlogRepo.Delete(ctx, model); err != nil {
		return errors.Wrap(ctx, err)
	}
	return c.NoContent(http.StatusNoContent)
}

func (s *VLogServer) StreamStatus(c echo.Context) error {
	ctx := c.Request().Context()
	vlogID := c.Param("id")

	// SSEヘッダー設定
	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().WriteHeader(http.StatusOK)

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	// 最初に現在の状態を送信
	vlog, err := s.vlogRepo.GetByID(ctx, &domain.Vlog{
		BaseModel: domain.BaseModel{ID: vlogID},
	})
	if err == nil {
		res := response.ToVLogGetByIDResponse(vlog)
		data, _ := json.Marshal(res)
		fmt.Fprintf(c.Response(), "data: %s\n\n", data)
		c.Response().Flush()
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			// VLogステータスを取得
			vlog, err := s.vlogRepo.GetByID(ctx, &domain.Vlog{
				BaseModel: domain.BaseModel{ID: vlogID},
			})
			if err != nil {
				// エラー時は継続を試みるが、あまりにひどい場合は終了
				continue
			}

			// SSEフォーマットでデータ送信
			res := response.ToVLogGetByIDResponse(vlog)
			data, _ := json.Marshal(res)
			fmt.Fprintf(c.Response(), "data: %s\n\n", data)
			c.Response().Flush()

			// 完了または失敗時は最後に1回送信して終了
			if vlog.Status == domain.VlogStatusCompleted || vlog.Status == domain.VlogStatusFailed {
				return nil
			}
		}
	}
}
