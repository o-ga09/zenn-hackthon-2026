package domain

import "gorm.io/gorm"

type SubtitleSegment struct {
	gorm.Model
	Index int    `gorm:"column:index" json:"index"`
	Start string `gorm:"column:start" json:"start"` // "00:00:01,000"
	End   string `gorm:"column:end" json:"end"`     // "00:00:04,000"
	Text  string `gorm:"column:text" json:"text"`   // 表示テキスト
}
