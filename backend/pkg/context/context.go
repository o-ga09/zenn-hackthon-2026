package context

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CtxUserKey string
type CtxRequestIDKey string
type CfgKey string
type DBKey string

const ConfigKey CfgKey = "config"
const USERID CtxUserKey = "userID"
const REQUESTID CtxRequestIDKey = "requestId"
const DB DBKey = "db"

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
