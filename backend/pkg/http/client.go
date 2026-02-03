package http

import (
	"fmt"
	"io"
	"net/http"
)

// fetchMediaData はURLからメディアデータを取得する
func FetchMediaData(url string, fallbackContentType string) ([]byte, string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, "", fmt.Errorf("failed to fetch media from URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("failed to fetch media: status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read media data: %w", err)
	}

	// Content-Typeを取得
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = fallbackContentType
	}
	if contentType == "" {
		contentType = http.DetectContentType(data)
	}

	return data, contentType, nil
}
