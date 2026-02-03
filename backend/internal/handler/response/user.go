package response

import (
	"github.com/o-ga09/zenn-hackthon-2026/internal/domain"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/date"
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
	TokenBalance   *int    `json:"tokenBalance,omitempty"`
	DisplayName    *string `json:"displayName,omitempty"`
	Bio            *string `json:"bio,omitempty"`
	ProfileImage   *string `json:"profileImage,omitempty"`
	BirthDay       *string `json:"birthDay,omitempty"`
	Gender         *string `json:"gender,omitempty"`
	FollowersCount *int    `json:"followersCount,omitempty"`
	FollowingCount *int    `json:"followingCount,omitempty"`
	IsPublic       *bool   `json:"isPublic,omitempty"`
	CreatedAt      string  `json:"createdAt"`
	UpdatedAt      string  `json:"updatedAt"`
}

func ToResponse(user *domain.User) *UserResponse {
	resp := &UserResponse{
		ID:             user.ID,
		Version:        user.Version,
		UID:            user.UID,
		Name:           user.Name,
		Type:           user.Type,
		Plan:           user.Plan,
		CreatedAt:      date.Format(user.CreatedAt),
		UpdatedAt:      date.Format(user.UpdatedAt),
		TokenBalance:   ptr.Int64ToPtr(user.TokenBalance.Int64),
		DisplayName:    ptr.StringToPtr(user.DisplayName.String),
		Bio:            ptr.StringToPtr(user.Bio.String),
		ProfileImage:   ptr.StringToPtr(user.ProfileImage.String),
		BirthDay:       ptr.StringToPtr(user.BirthDay.String),
		Gender:         ptr.StringToPtr(user.Gender.String),
		FollowersCount: ptr.Int64ToPtr(user.FollowersCount.Int64),
		FollowingCount: ptr.Int64ToPtr(user.FollowingCount.Int64),
	}

	if user.IsPublic.Valid {
		resp.IsPublic = ptr.BoolToPtr(user.IsPublic.Bool)
	}

	return resp
}
