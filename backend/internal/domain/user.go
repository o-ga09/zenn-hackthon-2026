package domain

import (
	"context"
	"database/sql"
)

type User struct {
	BaseModel
	UID            string         `gorm:"column:uid"`             // Firebase UID
	Name           string         `gorm:"column:name"`            // Display name
	Type           string         `gorm:"column:type"`            // User type: admin, tavinikkiy, tavinikkiy-agent
	Plan           string         `gorm:"column:plan"`            // Subscription plan: free, premium
	TokenBalance   sql.NullInt64  `gorm:"column:token_balance"`   // Token balance
	IsPublic       sql.NullBool   `gorm:"column:is_public"`       // Is profile public
	DisplayName    sql.NullString `gorm:"column:display_name"`    // Display name
	Bio            sql.NullString `gorm:"column:bio"`             // Biography
	ProfileImage   sql.NullString `gorm:"column:profile_image"`   // Profile image URL
	BirthDay       sql.NullString `gorm:"column:birth_day"`       // Birth
	Gender         sql.NullString `gorm:"column:gender"`          // Gender
	FollowersCount sql.NullInt64  `gorm:"column:followers_count"` // Number of followers
	FollowingCount sql.NullInt64  `gorm:"column:following_count"` // Number of following
}

// IUserRepository ユーザーリポジトリのインターフェース
type IUserRepository interface {
	// Create 新規ユーザーを作成
	Create(ctx context.Context, user *User) error
	// FindByID IDでユーザーを検索
	FindByID(ctx context.Context, cond *User) (*User, error)
	// FindByUID Firebase UIDでユーザーを検索
	FindByUID(ctx context.Context, cond *User) (*User, error)
	// FindAll 全ユーザーを取得
	FindAll(ctx context.Context, opts *FindOptions) ([]*User, error)
	// Update ユーザー情報を更新
	Update(ctx context.Context, user *User) error
	// Delete ユーザーを削除（論理削除）
	Delete(ctx context.Context, cond *User) error
}

type FindOptions struct {
	Limit  int
	Offset int
}
