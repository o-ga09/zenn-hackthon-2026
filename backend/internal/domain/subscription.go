package domain

import (
	"time"

	"gorm.io/gorm"
)

type Subscription struct {
	gorm.Model
	UID              string    `gorm:"column:uid"`
	Plan             string    `gorm:"column:plan"`   // "monthly", "yearly"
	Status           string    `gorm:"column:status"` // "active", "cancelled", "expired"
	StripeCustomerID string    `gorm:"column:stripe_customer_id"`
	StripeSubID      string    `gorm:"column:stripe_subscription_id"`
	CurrentPeriodEnd time.Time `gorm:"column:current_period_end"`
	CreatedAt        time.Time `gorm:"column:created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at"`
}
