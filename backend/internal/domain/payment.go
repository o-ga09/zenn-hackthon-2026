package domain

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model
	UID             string     `gorm:"column:uid"`
	Type            string     `gorm:"column:type"`           // "token_purchase", "subscription"
	Amount          int        `gorm:"column:amount"`         // 金額（円）
	TokensGranted   int        `gorm:"column:tokens_granted"` // 付与トークン数
	Status          string     `gorm:"column:status"`         // "pending", "completed", "failed"
	StripePaymentID string     `gorm:"column:stripe_payment_id"`
	CreatedAt       time.Time  `gorm:"column:created_at"`
	CompletedAt     *time.Time `gorm:"column:completed_at"`
}
