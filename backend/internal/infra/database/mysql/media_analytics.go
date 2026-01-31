package mysql

import (
	"context"

	"github.com/o-ga09/zenn-hackthon-2026/internal/domain"
	Ctx "github.com/o-ga09/zenn-hackthon-2026/pkg/context"
)

type MediaAnalyticsRepository struct{}

func (r *MediaAnalyticsRepository) Save(ctx context.Context, analytics *domain.MediaAnalytics) error {
	if err := Ctx.GetDB(ctx).Save(analytics).Error; err != nil {
		return err
	}
	return nil
}

func (r *MediaAnalyticsRepository) FindByFileID(ctx context.Context, fileID string) (*domain.MediaAnalytics, error) {
	var analytics domain.MediaAnalytics
	if err := Ctx.GetDB(ctx).
		Preload("Objects").
		Preload("Landmarks").
		Preload("Activities").
		Where("file_id = ?", fileID).
		First(&analytics).Error; err != nil {
		return nil, err
	}
	return &analytics, nil
}
