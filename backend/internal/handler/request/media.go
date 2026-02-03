package request

import "mime/multipart"

// multipart/form-dataに対応したメディアアップロードリクエスト
type MediaUploadRequest struct {
	File []*multipart.FileHeader `form:"file" validate:"required,dive"` // アップロードされるファイル
}

type MediaGetRequest struct {
	Key string `param:"key" validate:"required"` // 画像キー
}

type MediaDeleteRequest struct {
	Key string `param:"key" validate:"required"` // 画像キー
}
