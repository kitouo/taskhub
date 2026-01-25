package main

import (
	"context"

	"github.com/kitouo/taskhub/internal/app"
	"github.com/kitouo/taskhub/internal/config"
	"github.com/kitouo/taskhub/internal/logx"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	logger := logx.New("server=api", logx.ParseLevel(cfg.LogLevel))
	logger.Info("starting server", "config="+cfg.SafeString())

	application, err := app.New(cfg, logger)
	if err != nil {
		logger.Error("app init failed", "err=", err)
		panic(err)
	}
	_ = application.Run(context.Background())
}
