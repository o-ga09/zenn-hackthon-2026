package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/o-ga09/zenn-hackthon-2026/internal/handler"
	"github.com/o-ga09/zenn-hackthon-2026/internal/infra/database/mysql"
	"github.com/o-ga09/zenn-hackthon-2026/internal/infra/storage"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/config"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/logger"
)

type Server struct {
	Port   string
	Engine *echo.Echo
	User   handler.IUserServer
	Auth   handler.IAuthServer
	Image  handler.IImageServer
	VLog   handler.IVLogServer
	Agent  handler.IAgentServer
}

func New(ctx context.Context) *Server {
	env := config.GetCtxEnv(ctx)
	// ハンドラーの初期化
	r2Storage, err := storage.NewCloudflareR2Storage(ctx, env.CLOUDFLARE_R2_ACCOUNT_ID, env.CLOUDFLARE_R2_ACCESSKEY, env.CLOUDFLARE_R2_SECRETKEY, env.CLOUDFLARE_R2_BUCKET_NAME)
	if err != nil {
		log.Fatal(err)
	}
	userHandler := handler.NewUserServer(&mysql.UserRepository{}, r2Storage)
	authHandler := handler.NewAuthServer(&mysql.UserRepository{}, r2Storage)
	imageHandler := handler.NewImageServer(&mysql.MediaRepository{}, r2Storage)
	vlogHandler := handler.NewVLogServer(&mysql.VLogRepository{})
	agentHandler := handler.NewAgentServer(ctx, r2Storage)

	// Echoインスタンス作成
	e := echo.New()
	e.Validator = NewValidator()
	e.Binder = NewCustomBinder()

	return &Server{
		Port:   "8080",
		Engine: e,
		User:   userHandler,
		Auth:   authHandler,
		Image:  imageHandler,
		VLog:   vlogHandler,
		Agent:  agentHandler,
	}
}

func (s *Server) Run(ctx context.Context) error {
	// ミドルウェアの設定
	s.Engine.Use(middleware.Recover())
	s.Engine.Use(AddID(ctx))
	s.Engine.Use(AddTime())
	s.Engine.Use(RequestLogger())
	s.Engine.Use(SetDB())
	s.Engine.Use(WithTimeout())
	s.Engine.Use(CORS())
	s.Engine.Use(middleware.BodyLimit("10M"))
	s.Engine.Use(middleware.Gzip())
	s.Engine.Use(ErrorHandler())

	// ルーティングの設定
	s.SetupApplicationRoute()

	// サーバーの起動
	port := fmt.Sprintf(":%s", s.Port)
	srv := &http.Server{
		Addr:    port,
		Handler: s.Engine,
	}

	// サーバーの起動
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(ctx, fmt.Sprintf("Failed to listen and serve: %v", err))
		}
	}()

	logger.Info(ctx, fmt.Sprintf("Server is running on %s", s.Port))
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info(ctx, "graceful shutdown")

	// サーバーのタイムアウト設定
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	// サーバーのシャットダウン
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error(ctx, fmt.Sprintf("failed to shutdown server: %v", err))
		return err
	}
	return nil
}
