package response

import "github.com/o-ga09/zenn-hackthon-2026/internal/domain"

type VLogListResponse struct {
	Total int        `json:"total"`
	Items []VLogItem `json:"items"`
}

type VLogItem struct {
	ID        string  `json:"id"`
	VideoID   string  `json:"video_id"`
	VideoURL  string  `json:"video_url"`
	ShareURL  string  `json:"share_url"`
	Duration  float64 `json:"duration"`
	Thumbnail string  `json:"thumbnail"`
	CreatedAt string  `json:"created_at"`
}

type VLogGetByIDResponse struct {
	ID        string  `json:"id"`
	VideoID   string  `json:"video_id"`
	VideoURL  string  `json:"video_url"`
	ShareURL  string  `json:"share_url"`
	Duration  float64 `json:"duration"`
	Thumbnail string  `json:"thumbnail"`
	CreatedAt string  `json:"created_at"`
}

func ToVLogItem(vlog *domain.Vlog) VLogItem {
	return VLogItem{
		ID:        vlog.ID,
		VideoID:   vlog.VideoID,
		ShareURL:  vlog.ShareURL,
		Duration:  vlog.Duration,
		VideoURL:  vlog.VideoURL,
		Thumbnail: vlog.Thumbnail,
		CreatedAt: vlog.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func ToVLogGetByIDResponse(vlog *domain.Vlog) VLogGetByIDResponse {
	return VLogGetByIDResponse{
		ID:        vlog.ID,
		VideoID:   vlog.VideoID,
		ShareURL:  vlog.ShareURL,
		Duration:  vlog.Duration,
		VideoURL:  vlog.VideoURL,
		Thumbnail: vlog.Thumbnail,
		CreatedAt: vlog.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
