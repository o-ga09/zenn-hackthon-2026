package domain

import (
	"context"
	"reflect"
	"time"
)

// VlogStatus はVLogの処理状態を表す型
type VlogStatus string

func (s VlogStatus) String() string {
	return string(s)
}

func (s VlogStatus) Equals(other VlogStatus) bool {
	return reflect.DeepEqual(s, other)
}

const (
	VlogStatusPending    VlogStatus = "pending"
	VlogStatusProcessing VlogStatus = "processing"
	VlogStatusCompleted  VlogStatus = "completed"
	VlogStatusFailed     VlogStatus = "failed"
)

type Vlog struct {
	BaseModel
	VideoID      string     `gorm:"column:video_id" json:"video_id"`
	VideoURL     string     `gorm:"column:video_url" json:"video_url"`
	ShareURL     string     `gorm:"column:share_url" json:"share_url"`
	Duration     float64    `gorm:"column:duration" json:"duration"`
	Thumbnail    string     `gorm:"column:thumbnail" json:"thumbnail"`
	Status       VlogStatus `gorm:"column:status;default:pending" json:"status"`
	ErrorMessage string     `gorm:"column:error_message" json:"error_message,omitempty"`
	Progress     float64    `gorm:"column:progress;default:0" json:"progress"`
	StartedAt    *time.Time `gorm:"column:started_at" json:"started_at,omitempty"`
	CompletedAt  *time.Time `gorm:"column:completed_at" json:"completed_at,omitempty"`
}

type IVLogRepository interface {
	List(ctx context.Context, opts *ListOptions) ([]*Vlog, error)
	GetByID(ctx context.Context, model *Vlog) (*Vlog, error)
	Delete(ctx context.Context, model *Vlog) error
	Create(ctx context.Context, vlog *Vlog) error
	Update(ctx context.Context, vlog *Vlog) error
	UpdateStatus(ctx context.Context, id string, status VlogStatus, errorMsg string, progress float64) error
}

type ListOptions struct {
	Offset *int
	Limit  *int
}
