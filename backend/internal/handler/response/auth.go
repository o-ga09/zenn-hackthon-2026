package response

// SessionResponse セッションレスポンス
type SessionResponse struct {
	Message   string `json:"message"`
	ExpiresIn int64  `json:"expires_in"`
}

// DeleteSessionResponse セッション削除レスポンス
type DeleteSessionResponse struct {
	Message string `json:"message"`
}
