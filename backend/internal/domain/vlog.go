package domain

import (
	"context"
)

type Vlog struct {
	BaseModel
	VideoID   string  `gorm:"column:video_id" json:"video_id"`
	VideoURL  string  `gorm:"column:video_url" json:"video_url"`
	ShareURL  string  `gorm:"column:share_url" json:"share_url"`
	Duration  float64 `gorm:"column:duration" json:"duration"`
	Thumbnail string  `gorm:"column:thumbnail" json:"thumbnail"`
}

type IVLogRepository interface {
	List(ctx context.Context, opts *ListOptions) ([]*Vlog, error)
	GetByID(ctx context.Context, model *Vlog) (*Vlog, error)
	Delete(ctx context.Context, model *Vlog) error
}

type ListOptions struct {
	Offset *int
	Limit  *int
}
