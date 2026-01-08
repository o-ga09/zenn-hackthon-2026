package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/o-ga09/zenn-hackthon-2026/internal/domain"
	"github.com/o-ga09/zenn-hackthon-2026/internal/handler/request"
	"github.com/o-ga09/zenn-hackthon-2026/internal/handler/response"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/config"
	Ctx "github.com/o-ga09/zenn-hackthon-2026/pkg/context"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/errors"
	"gorm.io/gorm"
)

type IAuthServer interface {
	SignUp(c echo.Context) error  // セッション作成
	SignOut(c echo.Context) error // セッション削除
	GetUser(c echo.Context) error // ログインユーザー情報取得
}

type AuthServer struct {
	repo domain.IUserRepository
}

func NewAuthServer(repo domain.IUserRepository) IAuthServer {
	return &AuthServer{
		repo: repo,
	}
}

// SignUp セッション作成
func (s *AuthServer) SignUp(c echo.Context) error {
	ctx := c.Request().Context()

	// リクエストボディのバインドとバリデーション
	var req request.CreateSessionRequest
	if err := c.Bind(&req); err != nil {
		return errors.Wrap(ctx, err)
	}
	if err := c.Validate(&req); err != nil {
		return errors.Wrap(ctx, err)
	}

	// デフォルトの有効期限（14日）
	expiresIn := int64(60 * 60 * 24 * 14) // 14日
	if req.ExpiresIn != nil {
		expiresIn = *req.ExpiresIn
	}

	// Firebaseアプリケーションの取得
	app, err := config.GetFirebaseApp(ctx)
	if err != nil {
		errors.MakeSystemError(ctx, "認証サービスの初期化に失敗しました")
		return errors.Wrap(ctx, err)
	}

	client, err := app.Auth(ctx)
	if err != nil {
		errors.MakeSystemError(ctx, "認証クライアントの初期化に失敗しました")
		return errors.Wrap(ctx, err)
	}

	// IDトークンの検証
	token, err := client.VerifyIDToken(ctx, req.IDToken)
	if err != nil {
		errors.MakeAuthorizationError(ctx, "無効なIDトークンです")
		return errors.Wrap(ctx, err)
	}

	// セッションCookieの作成
	cookie, err := client.SessionCookie(ctx, req.IDToken, time.Duration(expiresIn)*time.Second)
	if err != nil {
		errors.MakeSystemError(ctx, "セッションCookieの作成に失敗しました")
		return errors.Wrap(ctx, err)
	}

	// Cookieの設定
	c.SetCookie(&http.Cookie{
		Name:     "__session",
		Value:    cookie,
		MaxAge:   int(expiresIn),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})

	// ユーザーが存在するか確認
	_, err = s.repo.FindByUID(ctx, token.UID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// ユーザーが存在しない場合は何もしない（フロントエンドでユーザー作成APIを呼ぶ）
		} else {
			return errors.Wrap(ctx, err)
		}
	}

	return c.JSON(http.StatusOK, response.SessionResponse{
		Message:   "セッションを作成しました",
		ExpiresIn: expiresIn,
	})
}

// SignOut セッション削除
func (s *AuthServer) SignOut(c echo.Context) error {
	// Cookieの削除
	c.SetCookie(&http.Cookie{
		Name:     "__session",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})

	return c.JSON(http.StatusOK, response.DeleteSessionResponse{
		Message: "セッションを削除しました",
	})
}

// GetUser ログインユーザー情報取得
func (s *AuthServer) GetUser(c echo.Context) error {
	ctx := c.Request().Context()

	// AuthMiddlewareでセットされたUIDを取得
	uid := Ctx.GetCtxFromUser(ctx)
	if uid == "" {
		errors.MakeAuthorizationError(ctx, "認証されていません")
		return errors.MakeBusinessError(ctx, "認証されていません")
	}

	// UIDでユーザーを取得
	user, err := s.repo.FindByUID(ctx, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.MakeNotFoundError(ctx, "ユーザーが見つかりません")
		}
		return errors.Wrap(ctx, err)
	}

	return c.JSON(http.StatusOK, response.ToResponse(user))
}
