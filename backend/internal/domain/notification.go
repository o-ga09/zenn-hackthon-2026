package domain

import (
	"context"
	"database/sql"
)

// 通知タイプ定数
const (
	NotificationTypeMediaCompleted = "media_completed"
	NotificationTypeMediaFailed    = "media_failed"
	NotificationTypeVlogCompleted  = "vlog_completed"
	NotificationTypeVlogFailed     = "vlog_failed"
)

// Notification - 通知ドメインモデル
type Notification struct {
	BaseModel
	UserID  string         `json:"user_id" gorm:"type:varchar(255);not null;index"`
	Type    string         `json:"type" gorm:"type:varchar(50);not null"`
	Title   string         `json:"title" gorm:"type:varchar(255);not null"`
	Message string         `json:"message" gorm:"type:text;not null"`
	MediaID sql.NullString `json:"media_id,omitempty" gorm:"type:varchar(26);index"`
	VlogID  sql.NullString `json:"vlog_id,omitempty" gorm:"type:varchar(26);index"`
	Read    bool           `json:"read" gorm:"not null;default:false"`
}

// TableName - テーブル名を指定
func (Notification) TableName() string {
	return "notifications"
}

// INotificationRepository - 通知リポジトリインターフェース
type INotificationRepository interface {
	Create(ctx context.Context, notification *Notification) error
	FindByID(ctx context.Context, id string) (*Notification, error)
	FindByUserID(ctx context.Context, userID string) ([]*Notification, error)
	MarkAsRead(ctx context.Context, notification *Notification) error
	MarkAllAsRead(ctx context.Context, notification *Notification) (int64, error)
	Delete(ctx context.Context, notification *Notification) error
	DeleteAllByUserID(ctx context.Context, userID string) error
	CountUnread(ctx context.Context, userID string) (int64, error)
}
