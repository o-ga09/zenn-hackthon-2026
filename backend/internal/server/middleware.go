package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/o-ga09/zenn-hackthon-2026/internal/infra/database/mysql"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/config"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/constant"
	Ctx "github.com/o-ga09/zenn-hackthon-2026/pkg/context"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/errors"
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

func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var err error
			ctx := c.Request().Context()

			cookie, err := c.Cookie("__session")
			if err != nil {
				return errors.MakeAuthorizationError(ctx, "セッションCookieが見つかりません")
			}

			// Firebaseアプリケーションの取得
			app, err := config.GetFirebaseApp(ctx)
			if err != nil {
				return errors.MakeAuthorizationError(ctx, "認証サービスの初期化に失敗しました")
			}

			client, err := app.Auth(ctx)
			if err != nil {
				return errors.MakeAuthorizationError(ctx, "認証クライアントの初期化に失敗しました")
			}

			// セッションCookieの検証
			sessionToken, err := client.VerifySessionCookie(ctx, cookie.Value)
			if err != nil {
				return errors.MakeAuthorizationError(ctx, "無効なセッションCookieです")
			}

			ctx = c.Request().Context()
			ctx = Ctx.SetCtxFromUser(ctx, sessionToken.UID)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

func CORS() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:3000",
			"https://tavinikkiy.com",
		},
		AllowMethods: []string{
			echo.POST,
			echo.GET,
			echo.OPTIONS,
			echo.PUT,
			echo.DELETE,
		},
		AllowHeaders:     []string{"Content-Type", "Authorization", "X-Requested-With"},
		AllowCredentials: true,
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

func ErrorHandler() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err == nil {
				return nil
			}
			code := errors.GetCode(err)
			statusCode := toHTTPStatusCode(code)
			message := errors.GetMessage(err)
			return c.JSON(statusCode, map[string]string{"error": message})
		}
	}
}

func toHTTPStatusCode(code errors.ErrCode) int {
	switch code {
	case errors.ErrCodeCritical:
		return http.StatusInternalServerError
	case errors.ErrCodeBussiness:
		return http.StatusBadRequest
	case errors.ErrCodeNotFound:
		return http.StatusNotFound
	case errors.ErrCodeConflict:
		return http.StatusConflict
	case errors.ErrCodeUnAuthorized:
		return http.StatusUnauthorized
	case errors.ErrCodeForbidden:
		return http.StatusForbidden
	case errors.ErrCodeInValidArgument:
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}

// CustomValidator is a wrapper for the validator library
type CustomValidator struct {
	validator *validator.Validate
}

// Validate implements the echo.Validator interface
func (cv *CustomValidator) Validate(i interface{}) error {
	err := cv.validator.Struct(i)
	if err != nil {
		fmt.Println("⏰")
		return errors.MakeInvalidArgumentError(context.Background(), err.Error())
	}
	return nil
}
func NewValidator() *CustomValidator {
	return &CustomValidator{validator: validator.New()}
}
