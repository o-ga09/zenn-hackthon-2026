package domain

import (
	"time"

	"gorm.io/gorm"
)

type MediaAnalytics struct {
	gorm.Model
	FileID      string    `gorm:"column:file_id" json:"file_id"`
	Type        string    `gorm:"column:type" json:"type"`               // "image" or "video"
	Description string    `gorm:"column:description" json:"description"` // 全体的な説明
	Objects     []string  `gorm:"column:objects" json:"objects"`         // 検出されたオブジェクト
	Landmarks   []string  `gorm:"column:landmarks" json:"landmarks"`     // 観光地・ランドマーク
	Activities  []string  `gorm:"column:activities" json:"activities"`   // アクティビティ
	Mood        string    `gorm:"column:mood" json:"mood"`               // 雰囲気（楽しい、穏やか、など）
	Timestamp   time.Time `gorm:"column:timestamp" json:"timestamp"`
}
