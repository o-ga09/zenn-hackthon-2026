package server

import (
	"reflect"
	"strings"

	"github.com/labstack/echo"
)

// CustomBinder カスタムバインダー
type CustomBinder struct {
	defaultBinder echo.Binder
}

// NewCustomBinder カスタムバインダーの生成
func NewCustomBinder() *CustomBinder {
	return &CustomBinder{
		defaultBinder: &echo.DefaultBinder{},
	}
}

// Bind リクエストデータを構造体にバインド
func (cb *CustomBinder) Bind(i interface{}, c echo.Context) error {
	// デフォルトのバインド処理（クエリパラメータ、JSONボディなど）
	if err := cb.defaultBinder.Bind(i, c); err != nil && err != echo.ErrUnsupportedMediaType {
		return err
	}

	// パスパラメータのバインド
	if err := cb.bindPathParams(i, c); err != nil {
		return err
	}

	return nil
}

// bindPathParams パスパラメータを構造体フィールドにバインド
func (cb *CustomBinder) bindPathParams(i interface{}, c echo.Context) error {
	typ := reflect.TypeOf(i)
	val := reflect.ValueOf(i)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	// 構造体でない場合は何もしない
	if typ.Kind() != reflect.Struct {
		return nil
	}

	// 各フィールドをチェック
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		// paramタグを取得
		paramTag := field.Tag.Get("param")
		if paramTag == "" {
			continue
		}

		// タグ名を解析（例: "id" or "id,required"）
		paramName := strings.Split(paramTag, ",")[0]
		if paramName == "" {
			continue
		}

		// パスパラメータから値を取得
		paramValue := c.Param(paramName)
		if paramValue == "" {
			continue
		}

		// フィールドに値を設定
		if fieldValue.CanSet() && fieldValue.Kind() == reflect.String {
			fieldValue.SetString(paramValue)
		}
	}

	return nil
}
