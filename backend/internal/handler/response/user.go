package response

import (
	"github.com/o-ga09/zenn-hackthon-2026/internal/domain"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/ptr"
)

// UserResponse ユーザーレスポンス
type UserResponse struct {
	ID             string  `json:"id"`
	Version        int     `json:"version"`
	UID            string  `json:"uid"`
	Name           string  `json:"name"`
	Type           string  `json:"type"`
	Plan           string  `json:"plan"`
	TokenBalance   *int    `json:"token_balance,omitempty"`
	DisplayName    *string `json:"display_name,omitempty"`
	Bio            *string `json:"bio,omitempty"`
	ProfileImage   *string `json:"profile_image,omitempty"`
	BirthDay       *string `json:"birth_day,omitempty"`
	Gender         *string `json:"gender,omitempty"`
	FollowersCount *int    `json:"followers_count,omitempty"`
	FollowingCount *int    `json:"following_count,omitempty"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

func ToResponse(user *domain.User) *UserResponse {
	resp := &UserResponse{
		ID:             user.ID,
		Version:        user.Version,
		UID:            user.UID,
		Name:           user.Name,
		Type:           user.Type,
		Plan:           user.Plan,
		CreatedAt:      user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"), // TODO: 共通化
		UpdatedAt:      user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"), // TODO: 共通化
		TokenBalance:   ptr.Int64ToPtr(user.TokenBalance.Int64),
		DisplayName:    ptr.StringToPtr(user.DisplayName.String),
		Bio:            ptr.StringToPtr(user.Bio.String),
		ProfileImage:   ptr.StringToPtr(user.ProfileImage.String),
		BirthDay:       ptr.StringToPtr(user.BirthDay.String),
		Gender:         ptr.StringToPtr(user.Gender.String),
		FollowersCount: ptr.Int64ToPtr(user.FollowersCount.Int64),
		FollowingCount: ptr.Int64ToPtr(user.FollowingCount.Int64),
	}

	return resp
}
