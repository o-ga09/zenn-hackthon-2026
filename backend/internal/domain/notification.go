package domain

import (
	"context"
	"database/sql"
	"fmt"
	"time"
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
	ID        string         `json:"id" gorm:"primaryKey;type:varchar(26)"`
	UserID    string         `json:"user_id" gorm:"type:varchar(255);not null;index"`
	Type      string         `json:"type" gorm:"type:varchar(50);not null"`
	Title     string         `json:"title" gorm:"type:varchar(255);not null"`
	Message   string         `json:"message" gorm:"type:text;not null"`
	MediaID   sql.NullString `json:"media_id,omitempty" gorm:"type:varchar(26);index"`
	VlogID    sql.NullString `json:"vlog_id,omitempty" gorm:"type:varchar(26);index"`
	Read      bool           `json:"read" gorm:"not null;default:false"`
	CreatedAt time.Time      `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"not null"`
}

// TableName - テーブル名を指定
func (Notification) TableName() string {
	return "notifications"
}

// Validate - 通知のバリデーション
func (n *Notification) Validate() error {
	if n.UserID == "" {
		return fmt.Errorf("user_id is required")
	}
	if n.Type == "" {
		return fmt.Errorf("type is required")
	}
	if n.Title == "" {
		return fmt.Errorf("title is required")
	}
	if n.Message == "" {
		return fmt.Errorf("message is required")
	}

	// タイプの検証
	validTypes := []string{
		NotificationTypeMediaCompleted,
		NotificationTypeMediaFailed,
		NotificationTypeVlogCompleted,
		NotificationTypeVlogFailed,
	}
	isValidType := false
	for _, validType := range validTypes {
		if n.Type == validType {
			isValidType = true
			break
		}
	}
	if !isValidType {
		return fmt.Errorf("invalid notification type: %s", n.Type)
	}

	return nil
}

// INotificationRepository - 通知リポジトリインターフェース
type INotificationRepository interface {
	Create(ctx context.Context, notification *Notification) error
	FindByID(ctx context.Context, id string) (*Notification, error)
	FindByUserID(ctx context.Context, userID string) ([]*Notification, error)
	MarkAsRead(ctx context.Context, id string) error
	MarkAllAsRead(ctx context.Context, userID string) (int64, error)
	Delete(ctx context.Context, id string) error
	CountUnread(ctx context.Context, userID string) (int64, error)
}
