// Package errors はアプリケーション全体で使用するエラー定義を提供する
package errors

import "errors"

// Genkit関連のエラー定義
var (
	ErrGenkitNotInitialized  = errors.New("genkit is not initialized")
	ErrStorageNotInitialized = errors.New("storage is not initialized")
	ErrFlowContextNotFound   = errors.New("flow context not found in context")
	ErrInvalidInput          = errors.New("invalid input")
	ErrMediaAnalysisFailed   = errors.New("media analysis failed")
	ErrVideoGenerationFailed = errors.New("video generation failed")
	ErrMaxMediaItemsExceeded = errors.New("max media items exceeded")
	ErrNoMediaItems          = errors.New("no media items provided")
	ErrToolExecutionFailed   = errors.New("tool execution failed")
)
