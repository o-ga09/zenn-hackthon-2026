package request

// CreateSessionRequest セッション作成リクエスト
type CreateSessionRequest struct {
	IDToken   string `json:"id_token" validate:"required,min=1" ja:"IDトークン"`
	ExpiresIn *int64 `json:"expires_in,omitempty" validate:"omitempty,gte=300,lte=1209600" ja:"有効期限"` // 5分〜14日（デフォルト14日）
}
