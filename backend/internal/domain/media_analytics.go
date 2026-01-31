package domain

import (
	"context"
)

type MediaAnalytics struct {
	BaseModel
	Version      int      `gorm:"column:version" json:"version"`
	CreateUserID *string  `gorm:"column:create_user_id" json:"create_user_id"`
	UpdateUserID *string  `gorm:"column:update_user_id" json:"update_user_id"`
	FileID       string   `gorm:"column:file_id" json:"file_id"`
	Description  string   `gorm:"column:description" json:"description"`               // 全体的な説明
	Objects      []string `gorm:"column:objects;serializer:json" json:"objects"`       // 検出されたオブジェクト
	Landmarks    []string `gorm:"column:landmarks;serializer:json" json:"landmarks"`   // 観光地・ランドマーク
	Activities   []string `gorm:"column:activities;serializer:json" json:"activities"` // アクティビティ
	Mood         string   `gorm:"column:mood" json:"mood"`                             // 雰囲気（楽しい、穏やか、など）
}

type IMediaAnalyticsRepository interface {
	Save(ctx context.Context, analytics *MediaAnalytics) error
	FindByFileID(ctx context.Context, fileID string) (*MediaAnalytics, error)
}
