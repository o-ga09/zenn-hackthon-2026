package response

type MediaListResponse struct {
	Media []*MediaListItem `json:"media"` // メディア一覧
	Total int              `json:"total"` // メディアの総数
}

type MediaImageUploadResponse struct {
	ID  string `json:"id"`  // アップロードされた画像のファイルID
	URL string `json:"url"` // 画像の取得URL
}

type MediaGetResponse struct {
	ID  string `json:"id"`  // 画像のファイルID
	URL string `json:"url"` // 画像の取得URL
}

type MediaListItem struct {
	ID          string `json:"id"`                   // ファイルID
	Type        string `json:"type"`                 // メディアタイプ
	ContentType string `json:"content_type"`         // コンテンツタイプ
	Size        int64  `json:"size"`                 // ファイルサイズ（バイト単位）
	URL         string `json:"url"`                  // 取得URL
	ImageData   string `json:"image_data,omitempty"` // Base64エンコードされた画像データ
	CreatedAt   string `json:"created_at"`           // 作成日時
}
