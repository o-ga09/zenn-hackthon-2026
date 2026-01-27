package genkit

import (
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

// RegisteredTools は登録されたすべてのツールを保持する
type RegisteredTools struct {
	// メディア分析
	AnalyzeMedia      ai.Tool
	AnalyzeMediaBatch ai.Tool

	// ストレージ操作
	UploadMedia       ai.Tool
	GenerateShareURL  ai.Tool
	GenerateThumbnail ai.Tool

	// VLog生成
	GenerateVlogVideo ai.Tool
}

// RegisterAllTools はすべてのツールを登録し、RegisteredToolsを返す
func RegisterAllTools(g *genkit.Genkit, baseURL string) *RegisteredTools {
	analyzeMedia := DefineAnalyzeMediaTool(g)
	analyzeMediaBatch := DefineAnalyzeMediaBatchTool(g, analyzeMedia)
	uploadMedia := DefineUploadMediaTool(g)
	generateShareURL := DefineGenerateShareURLTool(g, baseURL)
	generateThumbnail := DefineGenerateThumbnailTool(g)
	generateVlogVideo := DefineGenerateVlogVideoTool(g)

	return &RegisteredTools{
		AnalyzeMedia:      analyzeMedia,
		AnalyzeMediaBatch: analyzeMediaBatch,
		UploadMedia:       uploadMedia,
		GenerateShareURL:  generateShareURL,
		GenerateThumbnail: generateThumbnail,
		GenerateVlogVideo: generateVlogVideo,
	}
}

// AsToolList はai.WithToolsに渡すためのツールリストを返す
func (rt *RegisteredTools) AsToolList() []ai.Tool {
	return []ai.Tool{
		rt.AnalyzeMedia,
		rt.UploadMedia,
		rt.GenerateShareURL,
		rt.GenerateThumbnail,
		rt.GenerateVlogVideo,
	}
}

// AsVlogToolList はVLog生成に必要なツールのみを返す
func (rt *RegisteredTools) AsVlogToolList() []ai.Tool {
	return []ai.Tool{
		rt.AnalyzeMedia,
		rt.GenerateVlogVideo,
		rt.GenerateThumbnail,
		rt.GenerateShareURL,
	}
}
