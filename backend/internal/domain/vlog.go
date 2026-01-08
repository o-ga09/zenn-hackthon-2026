package domain

import (
	"gorm.io/gorm"
)

type Vlog struct {
	gorm.Model
	VideoID   string  `gorm:"column:video_id" json:"video_id"`
	VideoURL  string  `gorm:"column:video_url" json:"video_url"`
	ShareURL  string  `gorm:"column:share_url" json:"share_url"`
	Duration  float64 `gorm:"column:duration" json:"duration"`
	Thumbnail string  `gorm:"column:thumbnail" json:"thumbnail"`
}
