package domain

import (
	"gorm.io/gorm"
)

type TokenTransaction struct {
	gorm.Model
	UID         string                 `gorm:"column:uid"`
	Type        string                 `gorm:"column:type"`        // "purchase", "consumption", "bonus", "refund"
	Amount      int                    `gorm:"column:amount"`      // トークン数（消費時はマイナス）
	Balance     int                    `gorm:"column:balance"`     // 取引後の残高
	Description string                 `gorm:"column:description"` // "動画生成", "月額プラン付与"など
	Metadata    map[string]interface{} `gorm:"column:metadata"`
}
