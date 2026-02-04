package domain

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID           string         `gorm:"column:id"`
	Version      int            `gorm:"column:version"`
	CreateUserID *string        `gorm:"column:create_user_id"`
	UpdateUserID *string        `gorm:"column:update_user_id"`
	CreatedAt    time.Time      `gorm:"column:created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at"`
}
