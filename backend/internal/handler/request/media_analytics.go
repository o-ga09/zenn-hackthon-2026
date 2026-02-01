package request

// 分析結果取得用パスパラメータ
type GetMediaAnalyticsParam struct {
	ID string `param:"id" validate:"required"` // メディアID
}

// 分析結果更新リクエスト
type UpdateMediaAnalyticsRequest struct {
	ID          string   `param:"id" validate:"required"`                                  // メディアID
	Description *string  `json:"description,omitempty" validate:"omitempty,max=500"`       // 説明（最大500文字）
	Mood        *string  `json:"mood,omitempty" validate:"omitempty"`                      // 雰囲気
	Objects     []string `json:"objects,omitempty" validate:"omitempty,dive,min=1,max=50"` // 検出オブジェクト
	Landmarks   []string `json:"landmarks,omitempty" validate:"omitempty,dive,min=1,max=50"` // ランドマーク
	Activities  []string `json:"activities,omitempty" validate:"omitempty,dive,min=1,max=50"` // アクティビティ
}
