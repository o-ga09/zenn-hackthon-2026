package genkit

import (
	"fmt"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/o-ga09/zenn-hackthon-2026/internal/infra/storage"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/config"
	pkgerrors "github.com/o-ga09/zenn-hackthon-2026/pkg/errors"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/ulid"
)

// ============================================================
// メディアアップロードツール
// ============================================================

// UploadMediaInput はメディアアップロードツールの入力
type UploadMediaInput struct {
	Key         string `json:"key" jsonschema:"description=保存先のキー（パス）"`
	Data        []byte `json:"data" jsonschema:"description=ファイルデータ"`
	ContentType string `json:"contentType" jsonschema:"description=MIMEタイプ"`
}

// UploadMediaOutput はメディアアップロードツールの出力
type UploadMediaOutput struct {
	URL     string `json:"url" jsonschema:"description=アップロードされたファイルのURL"`
	Key     string `json:"key" jsonschema:"description=保存されたキー"`
	Success bool   `json:"success" jsonschema:"description=成功したかどうか"`
}

// DefineUploadMediaTool はメディアアップロードツールを定義する
func DefineUploadMediaTool(g *genkit.Genkit) ai.Tool {
	return genkit.DefineTool(g, "uploadMedia",
		"メディアファイルをCloudflare R2ストレージにアップロードする",
		func(ctx *ai.ToolContext, input UploadMediaInput) (UploadMediaOutput, error) {
			fc := GetFlowContext(ctx)
			if fc == nil {
				return UploadMediaOutput{}, pkgerrors.ErrFlowContextNotFound
			}
			if fc.Storage == nil {
				return UploadMediaOutput{}, pkgerrors.ErrStorageNotInitialized
			}

			objectKey, err := fc.Storage.UploadFile(ctx, input.Key, input.Data, input.ContentType)
			if err != nil {
				return UploadMediaOutput{
					Success: false,
				}, fmt.Errorf("%w: %v", pkgerrors.ErrToolExecutionFailed, err)
			}

			env := config.GetCtxEnv(ctx)
			return UploadMediaOutput{
				URL:     storage.ObjectURKFromKey(env.CLOUDFLARE_R2_PUBLIC_URL, objectKey),
				Key:     input.Key,
				Success: true,
			}, nil
		},
	)
}

// ============================================================
// 共有URL生成ツール
// ============================================================

// GenerateShareURLInput は共有URL生成ツールの入力
type GenerateShareURLInput struct {
	VideoID string `json:"videoId" jsonschema:"description=動画ID"`
	UserID  string `json:"userId" jsonschema:"description=ユーザーID"`
}

// GenerateShareURLOutput は共有URL生成ツールの出力
type GenerateShareURLOutput struct {
	ShareURL  string `json:"shareUrl" jsonschema:"description=生成された共有URL"`
	ShareCode string `json:"shareCode" jsonschema:"description=共有コード"`
}

// DefineGenerateShareURLTool は共有URL生成ツールを定義する
func DefineGenerateShareURLTool(g *genkit.Genkit, baseURL string) ai.Tool {
	return genkit.DefineTool(g, "generateShareURL",
		"VLogの共有URLを生成する",
		func(ctx *ai.ToolContext, input GenerateShareURLInput) (GenerateShareURLOutput, error) {
			shareCode, err := ulid.GenerateULID()
			if err != nil {
				return GenerateShareURLOutput{}, fmt.Errorf("%w: failed to generate share code", pkgerrors.ErrToolExecutionFailed)
			}

			shareURL := fmt.Sprintf("%s/share/%s", baseURL, shareCode)

			return GenerateShareURLOutput{
				ShareURL:  shareURL,
				ShareCode: shareCode,
			}, nil
		},
	)
}

// ============================================================
// サムネイル生成ツール
// ============================================================

// GenerateThumbnailInput はサムネイル生成ツールの入力
type GenerateThumbnailInput struct {
	VideoURL string `json:"videoUrl" jsonschema:"description=動画のURL"`
	VideoID  string `json:"videoId" jsonschema:"description=動画ID"`
	Width    int    `json:"width,omitempty" jsonschema:"description=サムネイルの幅"`
	Height   int    `json:"height,omitempty" jsonschema:"description=サムネイルの高さ"`
}

// GenerateThumbnailOutput はサムネイル生成ツールの出力
type GenerateThumbnailOutput struct {
	ThumbnailURL string `json:"thumbnailUrl" jsonschema:"description=生成されたサムネイルのURL"`
	Width        int    `json:"width" jsonschema:"description=サムネイルの幅"`
	Height       int    `json:"height" jsonschema:"description=サムネイルの高さ"`
}

// DefineGenerateThumbnailTool はサムネイル生成ツールを定義する
func DefineGenerateThumbnailTool(g *genkit.Genkit) ai.Tool {
	return genkit.DefineTool(g, "generateThumbnail",
		"動画からサムネイル画像を生成する",
		func(ctx *ai.ToolContext, input GenerateThumbnailInput) (GenerateThumbnailOutput, error) {
			fc := GetFlowContext(ctx)
			if fc == nil {
				return GenerateThumbnailOutput{}, pkgerrors.ErrFlowContextNotFound
			}

			width := input.Width
			height := input.Height
			if width == 0 {
				width = fc.Config.ThumbnailWidth
			}
			if height == 0 {
				height = fc.Config.ThumbnailHeight
			}

			// TODO: 実際のサムネイル生成処理を実装
			thumbnailKey := fmt.Sprintf("thumbnails/%s.jpg", input.VideoID)
			thumbnailURL := fmt.Sprintf("https://storage.example.com/%s", thumbnailKey)

			return GenerateThumbnailOutput{
				ThumbnailURL: thumbnailURL,
				Width:        width,
				Height:       height,
			}, nil
		},
	)
}
