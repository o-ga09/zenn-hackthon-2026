package agent

// ============================================================
// VLog生成用スキーマ
// ============================================================

// VlogInput はVLog生成フローの入力スキーマ
type VlogInput struct {
	UserID      string      `json:"userId" jsonschema:"description=ユーザーID,required"`
	MediaItems  []MediaItem `json:"mediaItems" jsonschema:"description=分析対象のメディアアイテム,required"`
	Title       string      `json:"title,omitempty" jsonschema:"description=VLogのタイトル（省略時は自動生成）"`
	TravelDate  string      `json:"travelDate,omitempty" jsonschema:"description=旅行日（YYYY-MM-DD形式）"`
	Destination string      `json:"destination,omitempty" jsonschema:"description=旅行先"`
	Style       VlogStyle   `json:"style,omitempty" jsonschema:"description=VLogのスタイル設定"`
}

// MediaItem は個別のメディアファイル情報
type MediaItem struct {
	FileID      string `json:"fileId" jsonschema:"description=ファイルID,required"`
	URL         string `json:"url" jsonschema:"description=メディアのURL,required"`
	Type        string `json:"type" jsonschema:"description=メディアタイプ（image または video）,required"`
	ContentType string `json:"contentType" jsonschema:"description=MIMEタイプ"`
	Timestamp   string `json:"timestamp,omitempty" jsonschema:"description=撮影日時（ISO 8601形式）"`
	Order       int    `json:"order,omitempty" jsonschema:"description=表示順序"`
	IsAnalyzed  bool   `json:"isAnalyzed,omitempty" jsonschema:"description=分析済みかどうか"`
}

// VlogStyle はVLog生成スタイルの設定
type VlogStyle struct {
	Theme      string `json:"theme,omitempty" jsonschema:"description=テーマ（adventure/relaxing/romantic/family）"`
	MusicMood  string `json:"musicMood,omitempty" jsonschema:"description=BGMの雰囲気"`
	Duration   int    `json:"duration,omitempty" jsonschema:"description=目標再生時間（秒）,default=60"`
	Transition string `json:"transition,omitempty" jsonschema:"description=トランジション効果（fade/slide/zoom）"`
}

// VlogOutput はVLog生成フローの出力スキーマ
type VlogOutput struct {
	VideoID      string          `json:"videoId" jsonschema:"description=生成されたVLogのID"`
	VideoURL     string          `json:"videoUrl" jsonschema:"description=VLog動画のURL"`
	ShareURL     string          `json:"shareUrl" jsonschema:"description=共有用URL"`
	ThumbnailURL string          `json:"thumbnailUrl" jsonschema:"description=サムネイル画像のURL"`
	Duration     float64         `json:"duration" jsonschema:"description=動画の長さ（秒）"`
	Title        string          `json:"title" jsonschema:"description=VLogのタイトル"`
	Description  string          `json:"description" jsonschema:"description=VLogの説明文"`
	Subtitles    []SubtitleEntry `json:"subtitles,omitempty" jsonschema:"description=字幕データ"`
	Analytics    VlogAnalytics   `json:"analytics" jsonschema:"description=分析結果サマリー"`
}

// SubtitleEntry は字幕の1エントリ
type SubtitleEntry struct {
	StartTime float64 `json:"startTime" jsonschema:"description=開始時間（秒）"`
	EndTime   float64 `json:"endTime" jsonschema:"description=終了時間（秒）"`
	Text      string  `json:"text" jsonschema:"description=字幕テキスト"`
}

// VlogAnalytics はVLog全体の分析結果サマリー
type VlogAnalytics struct {
	Locations  []string `json:"locations" jsonschema:"description=検出されたロケーション"`
	Activities []string `json:"activities" jsonschema:"description=検出されたアクティビティ"`
	Mood       string   `json:"mood" jsonschema:"description=全体の雰囲気"`
	Highlights []string `json:"highlights" jsonschema:"description=ハイライトシーンの説明"`
	MediaCount int      `json:"mediaCount" jsonschema:"description=使用されたメディア数"`
}

// ============================================================
// メディア分析用スキーマ（Tool用）
// ============================================================

// MediaAnalysisInput はメディア分析ツールの入力
type MediaAnalysisInput struct {
	FileID      string `json:"fileId" jsonschema:"description=ファイルID,required"`
	URL         string `json:"url" jsonschema:"description=分析対象のURL,required"`
	Type        string `json:"type" jsonschema:"description=メディアタイプ（image/video）,required"`
	ContentType string `json:"contentType,omitempty" jsonschema:"description=MIMEタイプ"`
}

// MediaAnalysisOutput はメディア分析ツールの出力
type MediaAnalysisOutput struct {
	FileID           string   `json:"fileId" jsonschema:"description=ファイルID"`
	Type             string   `json:"type" jsonschema:"description=メディアタイプ"`
	Description      string   `json:"description" jsonschema:"description=シーンの説明"`
	Objects          []string `json:"objects" jsonschema:"description=検出されたオブジェクト"`
	Landmarks        []string `json:"landmarks" jsonschema:"description=検出されたランドマーク・観光地"`
	Activities       []string `json:"activities" jsonschema:"description=検出されたアクティビティ"`
	Mood             string   `json:"mood" jsonschema:"description=シーンの雰囲気"`
	SuggestedCaption string   `json:"suggestedCaption" jsonschema:"description=提案されるキャプション"`
}

// MediaAnalysisBatchInput は複数メディア分析の入力
type MediaAnalysisBatchInput struct {
	Items []MediaAnalysisInput `json:"items" jsonschema:"description=分析対象のメディアリスト"`
}

// MediaAnalysisBatchOutput は複数メディア分析の出力
type MediaAnalysisBatchOutput struct {
	Results []MediaAnalysisOutput `json:"results" jsonschema:"description=分析結果のリスト"`
	Summary MediaAnalysisSummary  `json:"summary" jsonschema:"description=分析結果の全体サマリー"`
}

// MediaAnalysisSummary は分析結果の全体サマリー
type MediaAnalysisSummary struct {
	TotalItems       int      `json:"totalItems"`
	SuccessfulItems  int      `json:"successfulItems"`
	FailedItems      int      `json:"failedItems"`
	UniqueLocations  []string `json:"uniqueLocations"`
	UniqueActivities []string `json:"uniqueActivities"`
	OverallMood      string   `json:"overallMood"`
}

// ============================================================
// 進捗通知用スキーマ（SSE用）
// ============================================================

// FlowProgress はフロー実行中の進捗状態
type FlowProgress struct {
	Step        string  `json:"step" jsonschema:"description=現在のステップ名"`
	Progress    float64 `json:"progress" jsonschema:"description=進捗率（0-100）"`
	Message     string  `json:"message" jsonschema:"description=ユーザー向けメッセージ"`
	CurrentItem int     `json:"currentItem,omitempty" jsonschema:"description=現在処理中のアイテム番号"`
	TotalItems  int     `json:"totalItems,omitempty" jsonschema:"description=総アイテム数"`
}

// FlowStep は進捗のステップ定数
type FlowStep string

const (
	StepInitializing    FlowStep = "initializing"
	StepAnalyzing       FlowStep = "analyzing"
	StepGeneratingVideo FlowStep = "generating_video"
	StepUploadingVideo  FlowStep = "uploading_video"
	StepFinalizing      FlowStep = "finalizing"
	StepCompleted       FlowStep = "completed"
	StepFailed          FlowStep = "failed"
)

// ============================================================
// レガシースキーマ（後方互換性のため保持）
// ============================================================

// RecipeInput はサンプル用のレシピ入力スキーマ（非推奨）
// Deprecated: VlogInputを使用してください
type RecipeInput struct {
	Ingredient          string `json:"ingredient" jsonschema:"description=Main ingredient or cuisine type"`
	DietaryRestrictions string `json:"dietaryRestrictions,omitempty" jsonschema:"description=Any dietary restrictions"`
}

// Recipe はサンプル用のレシピ出力スキーマ（非推奨）
// Deprecated: VlogOutputを使用してください
type Recipe struct {
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	PrepTime     string   `json:"prepTime"`
	CookTime     string   `json:"cookTime"`
	Servings     int      `json:"servings"`
	Ingredients  []string `json:"ingredients"`
	Instructions []string `json:"instructions"`
	Tips         []string `json:"tips,omitempty"`
}
