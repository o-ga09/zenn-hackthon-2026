package mysql

import (
	"context"

	"github.com/o-ga09/zenn-hackthon-2026/internal/domain"
	Ctx "github.com/o-ga09/zenn-hackthon-2026/pkg/context"
)

type MediaRepository struct{}

func (r *MediaRepository) List(ctx context.Context, opts *domain.ListOpts) ([]*domain.Media, error) {
	var medias []*domain.Media
	if err := Ctx.GetDB(ctx).Find(&medias).Error; err != nil {
		return nil, err
	}
	return medias, nil
}

func (r *MediaRepository) GetByID(ctx context.Context, id string) (*domain.Media, error) {
	var media *domain.Media
	if err := Ctx.GetDB(ctx).Where("id = ?", id).First(&media).Error; err != nil {
		return nil, err
	}
	return media, nil
}

func (r *MediaRepository) Save(ctx context.Context, media *domain.Media) error {
	if err := Ctx.GetDB(ctx).Save(media).Error; err != nil {
		return err
	}
	return nil
}

func (r *MediaRepository) FindByFileID(ctx context.Context, model *domain.Media) (*domain.Media, error) {
	var media *domain.Media
	if err := Ctx.GetDB(ctx).First(&media, model).Error; err != nil {
		return nil, err
	}
	return media, nil
}

func (r *MediaRepository) DeleteByFileID(ctx context.Context, model *domain.Media) error {
	if err := Ctx.GetDB(ctx).Delete(model).Error; err != nil {
		return err
	}
	return nil
}
