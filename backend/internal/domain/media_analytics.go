package domain

import (
	"context"
)

type MediaAnalytics struct {
	BaseModel
	FileID      string           `gorm:"column:file_id" json:"file_id"`
	Description string           `gorm:"column:description" json:"description"` // 全体的な説明
	Mood        string           `gorm:"column:mood" json:"mood"`               // 雰囲気（楽しい、穏やか、など）
	Objects     []DetectedObject `gorm:"foreignKey:MediaAnalyticsID" json:"objects"`
	Landmarks   []Landmark       `gorm:"foreignKey:MediaAnalyticsID" json:"landmarks"`
	Activities  []Activity       `gorm:"foreignKey:MediaAnalyticsID" json:"activities"`
}

type DetectedObject struct {
	BaseModel
	MediaAnalyticsID string `gorm:"column:media_analytics_id" json:"media_analytics_id"`
	Name             string `gorm:"column:name" json:"name"`
}

func (DetectedObject) TableName() string {
	return "objects"
}

type Landmark struct {
	BaseModel
	MediaAnalyticsID string `gorm:"column:media_analytics_id" json:"media_analytics_id"`
	Name             string `gorm:"column:name" json:"name"`
}

type Activity struct {
	BaseModel
	MediaAnalyticsID string `gorm:"column:media_analytics_id" json:"media_analytics_id"`
	Name             string `gorm:"column:name" json:"name"`
}

type IMediaAnalyticsRepository interface {
	Save(ctx context.Context, analytics *MediaAnalytics) error
	FindByFileID(ctx context.Context, fileID string) (*MediaAnalytics, error)
	Update(ctx context.Context, analytics *MediaAnalytics) error
}
