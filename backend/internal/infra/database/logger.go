package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
	Ctx "github.com/o-ga09/zenn-hackthon-2026/pkg/context"
	CtxLogger "github.com/o-ga09/zenn-hackthon-2026/pkg/logger"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type SentryLogger struct {
	slowThreshold time.Duration
	logLevel      logger.LogLevel
}

func NewSentryLogger() *SentryLogger {
	return &SentryLogger{
		slowThreshold: 200 * time.Millisecond, // スロークエリの閾値
		logLevel:      logger.Info,            // ログレベル
	}
}

func (l *SentryLogger) LogMode(level logger.LogLevel) logger.Interface {
	newlogger := *l
	newlogger.logLevel = level
	return &newlogger
}

func (l *SentryLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Info {
		CtxLogger.Info(ctx, fmt.Sprintf(msg, data...))
	}
}

func (l *SentryLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Warn {
		CtxLogger.Warn(ctx, fmt.Sprintf(msg, data...))
	}
}

func (l *SentryLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Error {
		CtxLogger.Error(ctx, fmt.Sprintf(msg, data...))
	}
}

func (l *SentryLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.logLevel <= 0 {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	// クエリの実行時間を計測
	isSlowQuery := elapsed > l.slowThreshold && l.slowThreshold != 0

	span := sentry.StartSpan(ctx, "gorm.query")
	span.Description = "Database Query"
	span.SetData("sql", sql)
	span.SetData("rows", rows)
	span.SetData("elapsed", elapsed.String())
	span.SetData("is_slow_query", isSlowQuery)

	// データベース接続情報の追加
	if db := ctx.Value(Ctx.DB); db != nil {
		if gormDB, ok := db.(*gorm.DB); ok {
			if sqlDB, err := gormDB.DB(); err == nil {
				stats := sqlDB.Stats()
				span.SetData("db.max_open_connections", stats.MaxOpenConnections)
				span.SetData("db.open_connections", stats.OpenConnections)
				span.SetData("db.in_use", stats.InUse)
				span.SetData("db.idle", stats.Idle)
			}
		}
	}

	if err != nil {
		span.Status = sentry.SpanStatusInternalError
		span.SetData("error", err.Error())
	} else {
		span.Status = sentry.SpanStatusOK
	}

	defer span.Finish()
	switch {
	case err != nil && !errors.Is(err, gorm.ErrRecordNotFound) && l.logLevel >= logger.Error:
		if rows == -1 {
			CtxLogger.Error(ctx, fmt.Sprintf("[%.3fms] [rows:%s]", float64(elapsed.Nanoseconds())/1e6, "-"), "sql", sql, "error", err.Error())
		} else {
			CtxLogger.Error(ctx, fmt.Sprintf("[%.3fms] [rows:%d]", float64(elapsed.Nanoseconds())/1e6, rows), "sql", sql, "error", err.Error())
		}
	case elapsed > l.slowThreshold && l.slowThreshold != 0 && l.logLevel >= logger.Warn:
		if rows == -1 {
			CtxLogger.Warn(ctx, fmt.Sprintf("[%.3fms] [rows:%s]", float64(elapsed.Nanoseconds())/1e6, "-"), "sql", sql)
		} else {
			CtxLogger.Warn(ctx, fmt.Sprintf("[%.3fms] [rows:%d]", float64(elapsed.Nanoseconds())/1e6, rows), "sql", sql)
		}
	case l.logLevel == logger.Info:
		if rows == -1 {
			CtxLogger.Info(ctx, fmt.Sprintf("[%.3fms] [rows:%s]", float64(elapsed.Nanoseconds())/1e6, "-"), "sql", sql)
		} else {
			CtxLogger.Info(ctx, fmt.Sprintf("[%.3fms] [rows:%d]", float64(elapsed.Nanoseconds())/1e6, rows), "sql", sql)
		}
	}
}
