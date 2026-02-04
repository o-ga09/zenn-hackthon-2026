package context

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/config"
	"gorm.io/gorm"
)

type CtxUserKey string
type CtxRequestIDKey string
type CfgKey string
type DBKey string
type SkipOptimisticLockKey string

const ConfigKey CfgKey = "config"
const USERID CtxUserKey = "userID"
const REQUESTID CtxRequestIDKey = "requestId"
const DB DBKey = "db"
const SkipOptimisticLock SkipOptimisticLockKey = "skipOptimisticLock"

func GetCtxFromUser(ctx context.Context) string {
	userID, ok := ctx.Value(USERID).(string)
	if !ok {
		return ""
	}
	return userID
}

func SetCtxFromUser(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, USERID, userID)
}

func SetRequestID(ctx context.Context) context.Context {
	reqID := GetRequestID(ctx)
	if reqID != "" {
		return ctx
	}
	return context.WithValue(ctx, REQUESTID, uuid.NewString())
}

func GetRequestID(ctx context.Context) string {
	if reqID, ok := ctx.Value(REQUESTID).(string); ok {
		return reqID
	}
	return ""
}

func SetDB(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, DB, db)
}

func GetDB(ctx context.Context) *gorm.DB {
	if db, ok := ctx.Value(DB).(*gorm.DB); ok {
		return db.WithContext(ctx)
	}
	return nil
}

func SetRequestTime(ctx context.Context, time time.Time) context.Context {
	return context.WithValue(ctx, "requestTime", time)
}

func SetConfig(ctx context.Context, cfg *config.Config) context.Context {
	return context.WithValue(ctx, config.CtxEnvKey, cfg)
}

// WithSkipOptimisticLock は楽観ロックチェックをスキップするコンテキストを返す
// 一括更新など、楽観ロックチェックが不要な場合に使用する
func WithSkipOptimisticLock(ctx context.Context) context.Context {
	return context.WithValue(ctx, SkipOptimisticLock, true)
}
