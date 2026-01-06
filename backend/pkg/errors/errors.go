package errors

import (
	"context"
	"errors"
	"fmt"

	"github.com/morikuni/failure/v2"
	Ctx "github.com/o-ga09/zenn-hackthon-2026/pkg/context"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/logger"
)

const (
	callStack = "callStack"
)

type ErrType string

type ErrCode string

var (
	ErrCodeUnAuthorized    ErrCode = "unauthorized"     // 401
	ErrCodeForbidden       ErrCode = "forbidden"        // 403
	ErrCodeInValidArgument ErrCode = "invalid argument" // 400
	ErrCodeBussiness       ErrCode = "business error"   // 400
	ErrCodeConflict        ErrCode = "conflict"         // 409
	ErrCodeNotFound        ErrCode = "not found"        // 404
	ErrCodeCritical        ErrCode = "critical error"   // 500
)

var (
	ErrTypeUnAuthorized ErrType = "unauthorized"
	ErrTypeForbidden    ErrType = "forbidden"
	ErrTypeBussiness    ErrType = "business error"
	ErrTypeConflict     ErrType = "conflict"
	ErrTypeNotFound     ErrType = "not found"
	ErrTypeCritical     ErrType = "critical error"
)

var (
	// ドメインエラー
	ErrInvalidFirebaseID  = errors.New("不正なFirebaseIDです。")
	ErrInvalidUserID      = errors.New("不正なUserIDです。")
	ErrInvalidName        = errors.New("不正なユーザー名です。")
	ErrInvalidDisplayName = errors.New("不正な表示名です。")
	ErrInvalidGroupID     = errors.New("不正なグループIDです。")
	ErrInvalidRelationID  = errors.New("不正なリレーションIDです。")
	ErrInvalidTwitterID   = errors.New("不正なTwitterIDです。")
	ErrInvalidGender      = errors.New("性別の値の範囲が不正です。")
	ErrInvalidDateTime    = errors.New("日付のフォーマットが不正です。")
	ErrInvalidProfileURL  = errors.New("不正なプロフィールURLです。")
	ErrInvalidUserType    = errors.New("無効なユーザータイプフォーマットです。")
	ErrFollowed           = errors.New("すでにフォロー済みです。")
	ErrFollowSelf         = errors.New("自分自身をフォローすることはできません。")
	ErrRequestNotNil      = errors.New("リクエストが正しくありません。")

	// ulidエラー
	ErrEmptyULID   = errors.New("empty ulid")
	ErrInvalidULID = errors.New("invalid ulid")

	// データベースエラー
	ErrRecordNotFound         = errors.New("record not found")
	ErrConflict               = errors.New("conflict")
	ErrOptimisticLockConflict = errors.New("optimistic lock conflict")
	ErrForeignKeyConstraint   = errors.New("foreign key constraint error")
	ErrUniqueConstraint       = errors.New("unique constraint error")

	// 画像エラー
	ErrInvalidImageType  = errors.New("ファイルの種類が不正です。")
	ErrFailedImageName   = errors.New("ファイル名の生成に失敗しました。")
	ErrFailedDecodeImage = errors.New("画像のデコードに失敗しました。")
	ErrNotFoundImage     = errors.New("画像が見つかりません。")

	// リクエストエラー
	ErrRequestBodyNil = errors.New("リクエストボディが空です。")

	// その他エラー
	ErrSystem           = errors.New("システムエラーが発生しました。")
	ErrAuthorized       = errors.New("認証に失敗しました。")
	ErrUnauthorized     = errors.New("認可に失敗しました。")
	ErrInvalidArgument  = errors.New("バリデーションエラーが発生しました。")
	ErrInvalidOperation = errors.New("無効な操作です。")
	ErrNotFound         = errors.New("指定されたデータが見つかりません。")
)

// ginのcontextに認証エラーをセットして、ログ出力する
func MakeAuthorizationError(ctx context.Context, msg string) {
	var wrapped error
	if msg == "" {
		wrapped = failure.Translate(ErrAuthorized, ErrTypeUnAuthorized)
	} else {
		wrapped = failure.Translate(errors.New(msg), ErrTypeUnAuthorized)
	}
	c := Ctx.GetCtxGinCtx(ctx)
	if c != nil {
		c.Error(wrapped)
	}
	stack := getCallstack(wrapped)
	errMessage := GetMessage(wrapped)
	logger.Warn(ctx, errMessage, callStack, stack)
	c.Abort()
}

// ginのcontextに認可エラーをセットして、ログ出力する
func MakeAuthorizedError(ctx context.Context, msg string) {
	var wrapped error
	if msg == "" {
		wrapped = failure.Translate(ErrUnauthorized, ErrTypeForbidden)
	} else {
		wrapped = failure.Translate(errors.New(msg), ErrTypeForbidden)
	}
	c := Ctx.GetCtxGinCtx(ctx)
	if c != nil {
		c.Error(wrapped)
	}
	stack := getCallstack(wrapped)
	errMessage := GetMessage(wrapped)
	logger.Warn(ctx, errMessage, callStack, stack)
	c.Abort()
}

// ginのcontextにシステムエラーをセットして、ログ出力する
func MakeSystemError(ctx context.Context, msg string) {
	var wrapped error
	if msg == "" {
		wrapped = failure.Translate(ErrSystem, ErrTypeCritical)
	} else {
		wrapped = failure.Translate(errors.New(msg), ErrTypeCritical)
	}
	c := Ctx.GetCtxGinCtx(ctx)
	if c != nil {
		c.Error(wrapped)
	}
	stack := getCallstack(wrapped)
	errMessage := GetMessage(wrapped)
	logger.Error(ctx, errMessage, callStack, stack)
	c.Abort()
}

func MakeBusinessError(ctx context.Context, msg string) {
	var wrapped error
	if msg == "" {
		wrapped = failure.Translate(ErrSystem, ErrTypeBussiness)
	} else {
		wrapped = failure.Translate(errors.New(msg), ErrTypeBussiness)
	}
	c := Ctx.GetCtxGinCtx(ctx)
	if c != nil {
		c.Error(wrapped)
	}
	stack := getCallstack(wrapped)
	errMessage := GetMessage(wrapped)
	logger.Warn(ctx, errMessage, callStack, stack)
	c.Abort()
}

func MakeConflictError(ctx context.Context, msg string) {
	var wrapped error
	if msg == "" {
		wrapped = failure.Translate(ErrConflict, ErrTypeConflict)
	} else {
		wrapped = failure.Translate(errors.New(msg), ErrTypeConflict)
	}
	c := Ctx.GetCtxGinCtx(ctx)
	if c != nil {
		c.Error(wrapped)
	}
	stack := getCallstack(wrapped)
	errMessage := GetMessage(wrapped)
	logger.Warn(ctx, errMessage, callStack, stack)
	c.Abort()
}

func MakeNotFoundError(ctx context.Context, msg string) {
	var wrapped error
	if msg == "" {
		wrapped = failure.Translate(ErrNotFound, ErrTypeNotFound)
	} else {
		wrapped = failure.Translate(errors.New(msg), ErrTypeNotFound)
	}
	c := Ctx.GetCtxGinCtx(ctx)
	if c != nil {
		c.Error(wrapped)
	}
	stack := getCallstack(wrapped)
	errMessage := GetMessage(wrapped)
	logger.Warn(ctx, errMessage, callStack, stack)
	c.Abort()
}

// エラーをラップして、返す
// もしエラーがラップされていない場合は、システムエラーでラップして返す
func Wrap(ctx context.Context, err error) error {
	if !IsWrapped(err) {
		wrapped := failure.Translate(err, ErrTypeCritical)
		return wrapped
	}
	return failure.Wrap(err)
}

func New(ctx context.Context, err string) error {
	if err == "" {
		return failure.New(ctx)
	}
	return failure.New(ctx, failure.Field(failure.Message(err)))
}

func IsWrapped(err error) bool {
	return failure.Is(err, ErrTypeUnAuthorized, ErrTypeForbidden, ErrTypeBussiness, ErrTypeCritical, ErrTypeNotFound, ErrTypeConflict)
}

func Is(err error, target error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, target) {
		return true
	}
	return false
}

func GetMessage(err error) string {
	if err == nil {
		return ""
	}
	if IsWrapped(err) {
		return failure.ForceUnwrap(err).Error()
	}
	return err.Error()
}

func GetCode(err error) ErrCode {
	if err == nil {
		return ""
	}

	types := failure.CodeOf(err)
	if code, ok := types.(ErrType); ok {
		switch code {
		case ErrTypeUnAuthorized:
			return ErrCodeUnAuthorized
		case ErrTypeForbidden:
			return ErrCodeForbidden
		case ErrTypeBussiness:
			return ErrCodeInValidArgument
		case ErrTypeCritical:
			return ErrCodeCritical
		case ErrTypeNotFound:
			return ErrCodeNotFound
		case ErrTypeConflict:
			return ErrCodeConflict
		}
	}
	return ""
}

func getCallstack(err error) string {
	if err == nil {
		return ""
	}
	callstack := failure.CallStackOf(err)
	msg := ""
	for _, frame := range callstack.Frames() {
		msg += fmt.Sprintf("%+v\n", frame)
	}
	return msg
}
