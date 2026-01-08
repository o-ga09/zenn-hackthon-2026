package response

import "github.com/o-ga09/zenn-hackthon-2026/internal/domain"

// UserResponse ユーザーレスポンス
type UserResponse struct {
	ID           string `json:"id"`
	UID          string `json:"uid"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Plan         string `json:"plan"`
	TokenBalance *int64 `json:"token_balance,omitempty"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

func ToResponse(user *domain.User) *UserResponse {
	resp := &UserResponse{
		ID:        user.ID,
		UID:       user.UID,
		Name:      user.Name,
		Type:      user.Type,
		Plan:      user.Plan,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if user.TokenBalance.Valid {
		resp.TokenBalance = &user.TokenBalance.Int64
	}

	return resp
}
