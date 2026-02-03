package server

import (
	"mime/multipart"
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

	// クエリパラメータのバインド
	if err := cb.bindQueryParams(i, c); err != nil {
		return err
	}

	// multipart/form-dataのバインド
	if err := cb.bindFormParams(i, c); err != nil {
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

// bindQueryParams クエリパラメータを構造体フィールドにバインド
func (cb *CustomBinder) bindQueryParams(i interface{}, c echo.Context) error {
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

		// queryタグを取得
		queryTag := field.Tag.Get("query")
		if queryTag == "" {
			continue
		}

		// タグ名を解析（例: "name" or "name,required"）
		queryName := strings.Split(queryTag, ",")[0]
		if queryName == "" {
			continue
		}

		// クエリパラメータから値を取得
		queryValue := c.QueryParam(queryName)
		if queryValue == "" {
			continue
		}

		// フィールドに値を設定
		if fieldValue.CanSet() && fieldValue.Kind() == reflect.String {
			fieldValue.SetString(queryValue)
		}
	}

	return nil
}

func (cb *CustomBinder) bindFormParams(i interface{}, c echo.Context) error {
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

		// formタグを取得
		formTag := field.Tag.Get("form")
		if formTag == "" {
			continue
		}

		// タグ名を解析（例: "name" or "name,required"）
		formName := strings.Split(formTag, ",")[0]
		if formName == "" {
			continue
		}

		// フィールドの型に応じて処理を分岐
		switch fieldValue.Kind() {
		case reflect.Ptr:
			// *multipart.FileHeader の場合
			if fieldValue.Type() == reflect.TypeOf((*multipart.FileHeader)(nil)) {
				file, err := c.FormFile(formName)
				if err != nil {
					// ファイルが存在しない場合はスキップ
					continue
				}
				if fieldValue.CanSet() {
					fieldValue.Set(reflect.ValueOf(file))
				}
			}
		case reflect.Slice:
			// []*multipart.FileHeader の場合
			if fieldValue.Type() == reflect.TypeOf([]*multipart.FileHeader{}) {
				form, err := c.MultipartForm()
				if err != nil {
					continue
				}
				if files, exists := form.File[formName]; exists && len(files) > 0 {
					if fieldValue.CanSet() {
						fieldValue.Set(reflect.ValueOf(files))
					}
				}
			}
		case reflect.String:
			// 文字列フィールドの場合（従来の処理）
			formValue := c.FormValue(formName)
			if formValue == "" {
				continue
			}
			if fieldValue.CanSet() {
				fieldValue.SetString(formValue)
			}
		}
	}

	return nil
}
