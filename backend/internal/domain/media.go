package domain

import "context"

type Media struct {
	BaseModel
	Type        string `gorm:"column:type" json:"type"`                 // "image" or "video"
	ContentType string `gorm:"column:content_type" json:"content_type"` // MIMEタイプ
	Size        int64  `gorm:"column:size" json:"size"`                 // ファイルサイズ（バイト単位）
	URL         string `gorm:"column:url" json:"url"`                   // ファイルのURL
}

type IMediaRepository interface {
	List(ctx context.Context, opts *ListOpts) ([]*Media, error)
	Save(ctx context.Context, media *Media) error
	FindByFileID(ctx context.Context, media *Media) (*Media, error)
	DeleteByFileID(ctx context.Context, media *Media) error
}

type IImageStorage interface {
	Upload(ctx context.Context, key string, base64Data string) (string, error)
	Delete(ctx context.Context, key string) error
	Get(ctx context.Context, key string) (string, error)
	List(ctx context.Context, prefix string) (map[string]string, error)
}

type ListOpts struct {
	Limit  int
	Offset int
}
