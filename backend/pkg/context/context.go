package context

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/config"
	"gorm.io/gorm"
)

type CtxUserKey string
type CtxGinKey string
type CtxRequestIDKey string
type CfgKey string
type DBKey string

const ConfigKey CfgKey = "config"
const GinCtx CtxGinKey = "ginCtx"
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

func SetCtxGinCtx(ctx context.Context, c *gin.Context) context.Context {
	return context.WithValue(ctx, GinCtx, c)
}

func GetCtxGinCtx(ctx context.Context) *gin.Context {
	if c, ok := ctx.Value(GinCtx).(*gin.Context); ok {
		return c
	}
	return nil
}

func GetCfgFromCtx(ctx context.Context) *config.Config {
	return ctx.Value(config.CtxEnvKey).(*config.Config)
}

func SetDB(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, DB, db)
}

func SetRequestTime(ctx context.Context, time time.Time) context.Context {
	return context.WithValue(ctx, "requestTime", time)
}
