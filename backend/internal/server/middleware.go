package server

import (
	"context"
	"log/slog"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/o-ga09/zenn-hackthon-2026/internal/infra/database/mysql"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/constant"
	Ctx "github.com/o-ga09/zenn-hackthon-2026/pkg/context"
)

type RequestInfo struct {
	status                                            int
	contents_length                                   int64
	method, path, sourceIP, query, user_agent, errors string
	elapsed                                           time.Duration
}

func AddID(ctx context.Context) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := Ctx.SetRequestID(ctx)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

func WithTimeout() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
			defer cancel()
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

func RequestLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			req := c.Request()
			slog.Log(req.Context(), constant.SeverityInfo, "処理開始", "request Id", Ctx.GetRequestID(req.Context()))

			err := next(c)

			res := c.Response()
			r := &RequestInfo{
				status:          res.Status,
				contents_length: req.ContentLength,
				method:          req.Method,
				path:            req.URL.Path,
				sourceIP:        c.RealIP(),
				query:           req.URL.RawQuery,
				user_agent:      req.UserAgent(),
				errors:          "",
				elapsed:         time.Since(start),
			}
			if err != nil {
				r.errors = err.Error()
			}
			slog.Log(req.Context(), constant.SeverityInfo, "処理終了", "Request", r.LogValue(), "requestId", Ctx.GetRequestID(req.Context()))
			return err
		}
	}
}

func (r *RequestInfo) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Int("status", r.status),
		slog.Int64("Content-length", r.contents_length),
		slog.String("method", r.method),
		slog.String("path", r.path),
		slog.String("sourceIP", r.sourceIP),
		slog.String("query", r.query),
		slog.String("user_agent", r.user_agent),
		slog.String("errors", r.errors),
		slog.String("elapsed", r.elapsed.String()),
	)
}

func CORS() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			echo.POST,
			echo.GET,
			echo.OPTIONS,
		},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: false,
		MaxAge:           86400, // 24 hours in seconds
	})
}

func SetDB() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			db, err := mysql.Connect(ctx)
			if err != nil {
				slog.Log(ctx, constant.SeverityError, "DB接続に失敗しました", "error", err.Error())
				return err
			}
			ctx = Ctx.SetDB(ctx, db)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

func AddTime() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := Ctx.SetRequestTime(c.Request().Context(), time.Now())
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}
