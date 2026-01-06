package main

import (
	"context"
	"log"

	"github.com/o-ga09/zenn-hackthon-2026/internal/server"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/config"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/logger"
)

func main() {
	ctx := context.Background()

	ctx, err := config.New(ctx)
	if err != nil {
		log.Fatal(err)
	}

	logger.Logger(ctx)

	srv := server.New()
	if err := srv.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
