package mysql

import (
	"context"

	"github.com/o-ga09/zenn-hackthon-2026/internal/domain"
	Ctx "github.com/o-ga09/zenn-hackthon-2026/pkg/context"
)

type VLogRepository struct{}

func (r *VLogRepository) List(ctx context.Context, opts *domain.ListOptions) ([]*domain.Vlog, error) {
	var vlogs []*domain.Vlog
	userID := Ctx.GetCtxFromUser(ctx)
	if err := Ctx.GetDB(ctx).Where("create_user_id = ?", userID).Find(&vlogs).Error; err != nil {
		return nil, err
	}
	return vlogs, nil
}

func (r *VLogRepository) GetByID(ctx context.Context, model *domain.Vlog) (*domain.Vlog, error) {
	var vlog *domain.Vlog
	if err := Ctx.GetDB(ctx).Where("id = ?", model.ID).First(&vlog).Error; err != nil {
		return nil, err
	}
	return vlog, nil
}

func (r *VLogRepository) Delete(ctx context.Context, model *domain.Vlog) error {
	if err := Ctx.GetDB(ctx).Where("id = ?", model.ID).Delete(model).Error; err != nil {
		return err
	}
	return nil
}

func (r *VLogRepository) Create(ctx context.Context, vlog *domain.Vlog) error {
	if err := Ctx.GetDB(ctx).Create(vlog).Error; err != nil {
		return err
	}
	return nil
}

func (r *VLogRepository) Update(ctx context.Context, vlog *domain.Vlog) error {
	if err := Ctx.GetDB(ctx).Save(vlog).Error; err != nil {
		return err
	}
	return nil
}

func (r *VLogRepository) UpdateStatus(ctx context.Context, vlog *domain.Vlog) error {
	if err := Ctx.GetDB(ctx).Updates(vlog).Error; err != nil {
		return err
	}
	return nil
}
