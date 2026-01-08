package request

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
	UID  string `json:"uid" validate:"required,min=1,max=255"`
	Name string `json:"name" validate:"required,min=1,max=100"`
	Type string `json:"type" validate:"required"`
	Plan string `json:"plan" validate:"required"`
}

// UpdateUserRequest ユーザー更新リクエスト
type UpdateUserRequest struct {
	ID           string  `param:"id" validate:"required,gte=1"`
	Name         *string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Type         *string `json:"type,omitempty" validate:"omitempty"`
	Plan         *string `json:"plan,omitempty" validate:"omitempty"`
	TokenBalance *int64  `json:"token_balance,omitempty" validate:"omitempty,gte=0"`
}

// DeleteUserParam ユーザー削除のパスパラメータ
type DeleteUserParam struct {
	ID string `param:"id" validate:"required,gte=1"`
}
