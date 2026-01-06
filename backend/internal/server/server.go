package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/o-ga09/zenn-hackthon-2026/internal/router"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/logger"
)

type Server struct {
	Port   string
	engine *echo.Echo
}

func New() *Server {
	return &Server{
		Port:   "8080",
		engine: echo.New(),
	}
}

func (s *Server) Run(ctx context.Context) error {
	// ミドルウェアの設定
	s.engine.Use(middleware.Recover())
	s.engine.Use(AddID(ctx))
	s.engine.Use(AddTime())
	s.engine.Use(RequestLogger())
	s.engine.Use(SetDB())
	s.engine.Use(WithTimeout())
	s.engine.Use(CORS())
	s.engine.Use(middleware.BodyLimit("10M"))
	s.engine.Use(middleware.Gzip())

	// ルーティングの設定
	apiRoot := s.engine.Group("/api")
	router.SetupApplicationRoute(apiRoot)
	router.SetupSystemRoute(apiRoot)
	// サーバーの起動
	port := fmt.Sprintf(":%s", s.Port)
	srv := &http.Server{
		Addr:    port,
		Handler: s.engine,
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
