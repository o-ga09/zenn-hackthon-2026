package domain

import (
	"context"
	"database/sql"
	"reflect"
)

// MediaStatus はメディアの処理状態を表す
type MediaStatus string

func (m MediaStatus) Equals(b MediaStatus) bool {
	return reflect.DeepEqual(m, b)
}

func (m MediaStatus) String() string {
	return string(m)
}

const (
	MediaStatusPending   MediaStatus = "pending"   // アップロード待機中
	MediaStatusUploading MediaStatus = "uploading" // アップロード中
	MediaStatusAnalyzing MediaStatus = "analyzing" // 分析中
	MediaStatusCompleted MediaStatus = "completed" // アップロード完了
	MediaStatusFailed    MediaStatus = "failed"    // アップロード失敗
)

type Media struct {
	BaseModel
	ContentType  string         `gorm:"column:content_type" json:"content_type"`             // MIMEタイプ
	Size         int64          `gorm:"column:size" json:"size"`                             // ファイルサイズ（バイト単位）
	URL          sql.NullString `gorm:"column:url" json:"url"`                               // ファイルのURL
	Status       MediaStatus    `gorm:"column:status;default:completed" json:"status"`       // 処理状態
	Progress     float64        `gorm:"column:progress;default:1.0" json:"progress"`         // 進捗率（0.0〜1.0）
	ErrorMessage string         `gorm:"column:error_message" json:"error_message,omitempty"` // エラーメッセージ
}

type IMediaRepository interface {
	List(ctx context.Context, opts *ListOpts) ([]*Media, error)
	GetByID(ctx context.Context, id string) (*Media, error)
	Save(ctx context.Context, media *Media) error
	FindByFileID(ctx context.Context, media *Media) (*Media, error)
	DeleteByFileID(ctx context.Context, media *Media) error
}

type IImageStorage interface {
	Upload(ctx context.Context, key string, base64Data string) (string, error)
	UploadFile(ctx context.Context, key string, file []byte, contentType string) (string, error) // ファイルアップロード用
	Delete(ctx context.Context, key string) error
	Get(ctx context.Context, key string) (string, error)
	List(ctx context.Context, prefix string) (map[string]string, error)
}

type ListOpts struct {
	Limit  int
	Offset int
}
