package testutil

import (
	"context"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	Ctx "github.com/o-ga09/zenn-hackthon-2026/pkg/context"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/errors"
)

func SetUpTestRouter(t *testing.T) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	auth := AuthMiddleware()

	// withReqId := AddID(t.Context())
	errhandler := ErrorHandler()

	// 共通ミドルウェア設定
	// router.Use(withReqId)
	router.Use(auth)
	router.Use(errhandler)

	return router
}

func AddID(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request = c.Request.WithContext(ctx)
		ctx := Ctx.SetRequestID(c.Request.Context())
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := Ctx.SetCtxGinCtx(c.Request.Context(), c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()

		ginCtx := Ctx.GetCtxGinCtx(c.Request.Context())
		if ginCtx == nil {
			ginCtx = c
		}
		ginErr := ginCtx.Errors.Last()
		if ginErr != nil {
			// エラーメッセージの取得
			err := ginErr.Err
			message := errors.GetMessage(err)
			code := errors.GetCode(err)
			// レスポンスの設定
			c.JSON(ToHTTPCode(code), gin.H{
				"error": message,
				"code":  code,
			})
		}
	}
}

func ToHTTPCode(code errors.ErrCode) int {
	switch code {
	case errors.ErrCodeUnAuthorized:
		return http.StatusUnauthorized
	case errors.ErrCodeForbidden:
		return http.StatusForbidden
	case errors.ErrCodeInValidArgument:
		return http.StatusBadRequest
	case errors.ErrCodeCritical:
		return http.StatusInternalServerError
	case errors.ErrCodeNotFound:
		return http.StatusNotFound
	case errors.ErrCodeConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		userID := c.Request.Header.Get("x-tavinikkiy-user")

		ctx = Ctx.SetCtxFromUser(ctx, userID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
