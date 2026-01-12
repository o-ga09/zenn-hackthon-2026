package handler

import (
	"net/http"

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
