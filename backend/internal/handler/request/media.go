package request

type MediaImageUploadRequest struct {
	Base64Data string `json:"base64_data" validate:"required"` // Base64エンコードされた画像データ
}

type MediaGetRequest struct {
	Key string `param:"key" validate:"required"` // 画像キー
}

type MediaDeleteRequest struct {
	Key string `param:"key" validate:"required"` // 画像キー
}
