package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/config"
)

func main() {
	ctx := context.Background()

	ctx, err := config.New(ctx)
	if err != nil {
		log.Fatal(err)
	}

	ctx = config.InitGenAI(ctx)

	g := config.GetGenkitCtx(ctx)

	// マルチモーダルインプット
	// image, err := os.ReadFile("photo.png")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// resp, err := genkit.Generate(ctx, g,
	// 	ai.WithSystem("日本語で回答を生成してください。"),
	// 	ai.WithMessages(
	// 		ai.NewUserMessage(
	// 			ai.NewMediaPart("image/jpeg", "data:image/jpeg;base64,"+base64.StdEncoding.EncodeToString(image)),
	// 			ai.NewTextPart("Compose a poem about this image."),
	// 		),
	// 	),
	// )

	// マルチモーダルアウトプット(画像)
	// resp, err := genkit.Generate(ctx, g,
	// 	ai.WithModel(googlegenai.VertexAIModel(g, "imagen-3.0-fast-generate-001")),
	// 	ai.WithPrompt("Generate an image of a sunset over a mountain range."),
	// )

	// マルチモーダルアウトプット(動画)
	resp, err := genkit.Generate(ctx, g,
		ai.WithModelName("vertexai/veo-3.1-fast-generate-001"),
		ai.WithPrompt("A majestic dragon soaring over a mystical forest at dawn."),
	)

	if err != nil {
		log.Fatal(err)
	}

	for i, content := range resp.Message.Content {
		fmt.Println("Response:")
		fmt.Println("contentType:", content.ContentType)
		if content.IsImage() {
			err := SaveMediaPartToFile(content, fmt.Sprintf("output_image_%d.jpg", i))
			if err != nil {
				log.Fatal(err)
			}
		}
		if content.IsVideo() {
			fmt.Println("vodeo URL", content.Text)
		}
	}
}

// SaveMediaPartToFile saves a media part to a file.
// base64形式のデータをデコードして保存する。
func SaveMediaPartToFile(content *ai.Part, filename string) error {
	// content.Textからbase64データを取得
	// 形式: "data:image/jpeg;base64,<base64データ>" または直接base64データ
	base64Data := content.Text

	// "data:image/xxx;base64," のプレフィックスを削除
	if strings.Contains(base64Data, ";base64,") {
		parts := strings.Split(base64Data, ";base64,")
		if len(parts) == 2 {
			base64Data = parts[1]
		}
	}

	// base64デコード
	decodedData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return fmt.Errorf("failed to decode base64 data: %w", err)
	}

	// ファイルに書き込み
	err = os.WriteFile(filename, decodedData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("Image saved to %s\n", filename)
	return nil
}
