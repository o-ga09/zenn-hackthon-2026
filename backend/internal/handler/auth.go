package handler

import (
	"net/http"
	"strings"
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
	repo    domain.IUserRepository
	storage domain.IUserStorage
}

func NewAuthServer(repo domain.IUserRepository, storage domain.IUserStorage) IAuthServer {
	return &AuthServer{
		repo:    repo,
		storage: storage,
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

	cfg := config.GetCtxEnv(ctx)
	domainName := cfg.COOKIE_DOMAIN
	isSecure := true
	sameSite := http.SameSiteNoneMode

	// ローカル環境の設定
	if cfg.Env == "local" || cfg.Env == "dev" {
		domainName = "localhost"
		isSecure = false
		sameSite = http.SameSiteLaxMode
	}

	// Cookieの設定
	c.SetCookie(&http.Cookie{
		Name:     "__session",
		Value:    cookie,
		MaxAge:   int(expiresIn),
		HttpOnly: true,
		Secure:   isSecure,
		Domain:   domainName,
		Path:     "/",
		SameSite: sameSite,
	})

	// ユーザーが存在するか確認
	_, err = s.repo.FindByUID(ctx, &domain.User{UID: token.UID})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.Wrap(ctx, err)
	}

	return c.JSON(http.StatusOK, response.SessionResponse{
		Message:   "セッションを作成しました",
		ExpiresIn: expiresIn,
	})
}

// SignOut セッション削除
func (s *AuthServer) SignOut(c echo.Context) error {
	ctx := c.Request().Context()
	cfg := config.GetCtxEnv(ctx)

	domain := cfg.COOKIE_DOMAIN
	isSecure := true
	sameSite := http.SameSiteNoneMode

	// ローカル環境の設定
	if cfg.Env == "local" || cfg.Env == "dev" {
		domain = "" // ローカルではDomainを空にする
		isSecure = false
		sameSite = http.SameSiteLaxMode
	}

	// Cookieの削除
	c.SetCookie(&http.Cookie{
		Name:     "__session",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   isSecure,
		Domain:   domain,
		Path:     "/",
		SameSite: sameSite,
	})

	return c.JSON(http.StatusOK, response.DeleteSessionResponse{
		Message: "セッションを削除しました",
	})
}

// GetUser ログインユーザー情報取得
func (s *AuthServer) GetUser(c echo.Context) error {
	ctx := c.Request().Context()

	// AuthMiddlewareでセットされたUIDを取得
	id := Ctx.GetCtxFromUser(ctx)
	if id == "" {
		return errors.MakeBusinessError(ctx, "認証されていません")
	}

	// IDでユーザーを取得
	user, err := s.repo.FindByUID(ctx, &domain.User{BaseModel: domain.BaseModel{ID: id}})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.MakeNotFoundError(ctx, "ユーザーが見つかりません")
		}
		return errors.Wrap(ctx, err)
	}

	if user.ProfileImage.Valid && !strings.HasPrefix(user.ProfileImage.String, "https://") {
		user.ProfileImage.String, err = s.storage.Get(ctx, user.ProfileImage.String)
		if err != nil {
			return errors.Wrap(ctx, err)
		}
	}

	return c.JSON(http.StatusOK, response.ToResponse(user))
}
