package mysql

import (
	"context"

	"github.com/o-ga09/zenn-hackthon-2026/internal/domain"
	Ctx "github.com/o-ga09/zenn-hackthon-2026/pkg/context"
	"gorm.io/gorm"
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

func (r *MediaAnalyticsRepository) Update(ctx context.Context, analytics *domain.MediaAnalytics) error {
	db := Ctx.GetDB(ctx)

	// トランザクション内で既存の関連データを削除し、新しいデータを保存
	return db.Transaction(func(tx *gorm.DB) error {
		// 既存のObjects, Landmarks, Activitiesを削除
		if err := tx.Where("media_analytics_id = ?", analytics.ID).Delete(&domain.DetectedObject{}).Error; err != nil {
			return err
		}
		if err := tx.Where("media_analytics_id = ?", analytics.ID).Delete(&domain.Landmark{}).Error; err != nil {
			return err
		}
		if err := tx.Where("media_analytics_id = ?", analytics.ID).Delete(&domain.Activity{}).Error; err != nil {
			return err
		}

		// MediaAnalytics本体を更新
		if err := tx.Model(analytics).Updates(map[string]interface{}{
			"description": analytics.Description,
			"mood":        analytics.Mood,
		}).Error; err != nil {
			return err
		}

		// 新しい関連データを保存
		for _, obj := range analytics.Objects {
			obj.MediaAnalyticsID = analytics.ID
			if err := tx.Create(&obj).Error; err != nil {
				return err
			}
		}
		for _, landmark := range analytics.Landmarks {
			landmark.MediaAnalyticsID = analytics.ID
			if err := tx.Create(&landmark).Error; err != nil {
				return err
			}
		}
		for _, activity := range analytics.Activities {
			activity.MediaAnalyticsID = analytics.ID
			if err := tx.Create(&activity).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
