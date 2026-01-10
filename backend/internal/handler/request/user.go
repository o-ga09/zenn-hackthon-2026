package request

import (
	"github.com/o-ga09/zenn-hackthon-2026/internal/domain"
	nullvalue "github.com/o-ga09/zenn-hackthon-2026/pkg/null_value"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/ptr"
)

const (
	DefaultType = "tavinikkiy-agent"
	DefaultPlan = "free"
)

// ListQuery ユーザー一覧取得のクエリパラメータ
type ListQuery struct {
	Limit  *int `query:"limit" validate:"omitempty,gte=0,lte=100"`
	Offset *int `query:"offset" validate:"omitempty,gte=0"`
}

// GetByIDParam IDでユーザー取得のパスパラメータ
type GetByIDParam struct {
	ID string `param:"id" validate:"required" ja:"ユーザーID"`
}

// GetByUIDQuery UIDでユーザー取得のクエリパラメータ
type GetByUIDQuery struct {
	UID string `query:"uid" validate:"required,min=1,max=255"`
}

// CreateUserRequest ユーザー作成リクエスト
type CreateUserRequest struct {
	Plan           string  `json:"plan,omitempty" validate:"omitempty"`
	Type           string  `json:"type,omitempty" validate:"omitempty"`
	UID            string  `json:"uid,omitempty" validate:"required,min=1,max=255"`
	Name           *string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	TokenBalance   *int64  `json:"token_balance,omitempty" validate:"omitempty,gte=0"`
	IsPublic       *bool   `json:"is_public,omitempty" validate:"omitempty"`
	DisplayName    *string `json:"display_name,omitempty" validate:"omitempty,min=1,max=100"`
	Bio            *string `json:"bio,omitempty" validate:"omitempty,max=500"`
	ProfileImage   *string `json:"profile_image,omitempty" validate:"omitempty,url"`
	BirthDay       *string `json:"birth_day,omitempty" validate:"omitempty,datetime=2006-01-02"`
	Gender         *string `json:"gender,omitempty" validate:"omitempty"`
	FollowersCount *int    `json:"followers_count,omitempty" validate:"omitempty"`
	FollowingCount *int    `json:"following_count,omitempty" validate:"omitempty"`
}

// UpdateUserRequest ユーザー更新リクエスト
type UpdateUserRequest struct {
	ID             string  `param:"id" validate:"required,gte=1"`
	Version        int     `json:"version" validate:"required"`
	Plan           string  `json:"plan,omitempty" validate:"required"`
	Type           string  `json:"type,omitempty" validate:"required"`
	UID            string  `json:"uid,omitempty" validate:"required,min=1,max=255"`
	Name           *string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	TokenBalance   *int64  `json:"token_balance,omitempty" validate:"omitempty,gte=0"`
	IsPublic       *bool   `json:"isPublic,omitempty" validate:"omitempty"`
	DisplayName    *string `json:"displayName,omitempty" validate:"omitempty,min=1,max=100"`
	Bio            *string `json:"bio,omitempty" validate:"omitempty,max=500"`
	ProfileImage   *string `json:"profileImage,omitempty" validate:"omitempty,url"`
	BirthDay       *string `json:"birthday,omitempty" validate:"omitempty,datetime=2006-01-02"`
	Gender         *string `json:"gender,omitempty" validate:"omitempty"`
	FollowersCount *int    `json:"followersCount,omitempty" validate:"omitempty"`
	FollowingCount *int    `json:"followingCount,omitempty" validate:"omitempty"`
}

// DeleteUserParam ユーザー削除のパスパラメータ
type DeleteUserParam struct {
	ID string `param:"id" validate:"required,gte=1"`
}

func (req *CreateUserRequest) ToUser() *domain.User {

	// TODO: 共通処理化
	// UIDの先頭10文字をNameに設定
	name := req.UID
	if len(req.UID) > 10 {
		name = req.UID[:10]
	}

	return &domain.User{
		UID:            req.UID,
		Name:           name,
		Type:           DefaultType,
		Plan:           DefaultPlan,
		IsPublic:       nullvalue.ToNullBool(false),
		DisplayName:    ptr.PtrStringToNullString(req.DisplayName),
		Bio:            ptr.PtrStringToNullString(req.Bio),
		ProfileImage:   ptr.PtrStringToNullString(req.ProfileImage),
		BirthDay:       ptr.PtrStringToNullString(req.BirthDay),
		Gender:         ptr.PtrStringToNullString(req.Gender),
		FollowersCount: ptr.PtrIntToNullInt64(req.FollowersCount),
		FollowingCount: ptr.PtrIntToNullInt64(req.FollowingCount),
	}
}

func (req *UpdateUserRequest) ToUser() *domain.User {
	return &domain.User{
		BaseModel: domain.BaseModel{
			ID:      req.ID,
			Version: req.Version,
		},
		UID:            req.UID,
		Name:           ptr.PtrToString(req.Name),
		Type:           req.Type,
		Plan:           req.Plan,
		IsPublic:       ptr.PtrBoolToNullBool(req.IsPublic),
		DisplayName:    ptr.PtrStringToNullString(req.DisplayName),
		Bio:            ptr.PtrStringToNullString(req.Bio),
		ProfileImage:   ptr.PtrStringToNullString(req.ProfileImage),
		BirthDay:       ptr.PtrStringToNullString(req.BirthDay),
		Gender:         ptr.PtrStringToNullString(req.Gender),
		FollowersCount: ptr.PtrIntToNullInt64(req.FollowersCount),
		FollowingCount: ptr.PtrIntToNullInt64(req.FollowingCount),
	}
}
