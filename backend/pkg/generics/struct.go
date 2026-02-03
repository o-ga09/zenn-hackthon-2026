package generics

import (
	"encoding/json"
	"fmt"
)

// ConvertToStruct はinterface{}（通常はmap[string]interface{}）を指定した構造体に変換する
func ConvertToStruct[T any](raw interface{}) (T, error) {
	var result T

	// すでに目的の型の場合はそのまま返す
	if typed, ok := raw.(T); ok {
		return typed, nil
	}

	// JSONを経由して変換
	jsonBytes, err := json.Marshal(raw)
	if err != nil {
		return result, fmt.Errorf("failed to marshal to JSON: %w", err)
	}

	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		return result, fmt.Errorf("failed to unmarshal from JSON: %w", err)
	}

	return result, nil
}
