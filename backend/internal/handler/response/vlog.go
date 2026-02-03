package response

import "github.com/o-ga09/zenn-hackthon-2026/internal/domain"

type VLogListResponse struct {
	Total int        `json:"total"`
	Items []VLogItem `json:"items"`
}

type VLogItem struct {
	ID           string  `json:"id"`
	VideoID      string  `json:"video_id"`
	VideoURL     string  `json:"video_url"`
	ShareURL     string  `json:"share_url"`
	Duration     float64 `json:"duration"`
	Thumbnail    string  `json:"thumbnail"`
	Status       string  `json:"status"`
	ErrorMessage string  `json:"error_message,omitempty"`
	Progress     float64 `json:"progress"`
	CreatedAt    string  `json:"created_at"`
}

type VLogGetByIDResponse struct {
	ID           string  `json:"id"`
	VideoID      string  `json:"video_id"`
	VideoURL     string  `json:"video_url"`
	ShareURL     string  `json:"share_url"`
	Duration     float64 `json:"duration"`
	Thumbnail    string  `json:"thumbnail"`
	Status       string  `json:"status"`
	ErrorMessage string  `json:"error_message,omitempty"`
	Progress     float64 `json:"progress"`
	CreatedAt    string  `json:"created_at"`
}

// CreateVLogResponse はVLog生成APIのレスポンス
type CreateVLogResponse struct {
	VlogID string `json:"vlogId"`
	Status string `json:"status"`
}

// MediaStatusResponse はメディアステータスSSEのレスポンス
type MediaStatusResponse struct {
	Medias         []*domain.Media `json:"medias"`
	TotalItems     int             `json:"total_items"`
	CompletedItems int             `json:"completed_items"`
	FailedItems    int             `json:"failed_items"`
	AllCompleted   bool            `json:"all_completed"`
}

// AnalyzeMediaResponse はメディア分析APIのレスポンス
type AnalyzeMediaResponse struct {
	MediaIDs []string `json:"media_ids"`
	Status   string   `json:"status"`
}

func ToVLogItem(vlog *domain.Vlog) VLogItem {
	return VLogItem{
		ID:           vlog.ID,
		VideoID:      vlog.VideoID,
		ShareURL:     vlog.ShareURL,
		Duration:     vlog.Duration,
		VideoURL:     vlog.VideoURL,
		Thumbnail:    vlog.Thumbnail,
		Status:       string(vlog.Status),
		ErrorMessage: vlog.ErrorMessage,
		Progress:     vlog.Progress,
		CreatedAt:    vlog.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func ToVLogGetByIDResponse(vlog *domain.Vlog) VLogGetByIDResponse {
	return VLogGetByIDResponse{
		ID:           vlog.ID,
		VideoID:      vlog.VideoID,
		ShareURL:     vlog.ShareURL,
		Duration:     vlog.Duration,
		VideoURL:     vlog.VideoURL,
		Thumbnail:    vlog.Thumbnail,
		Status:       string(vlog.Status),
		ErrorMessage: vlog.ErrorMessage,
		Progress:     vlog.Progress,
		CreatedAt:    vlog.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
