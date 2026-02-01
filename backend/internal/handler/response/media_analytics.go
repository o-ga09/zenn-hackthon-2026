package response

// 分析結果レスポンス
type MediaAnalyticsResponse struct {
	FileID      string   `json:"file_id"`   // メディアID
	Description string   `json:"description"` // 全体的な説明
	Mood        string   `json:"mood"`      // 雰囲気
	Objects     []string `json:"objects"`   // 検出オブジェクト
	Landmarks   []string `json:"landmarks"` // ランドマーク
	Activities  []string `json:"activities"` // アクティビティ
}
