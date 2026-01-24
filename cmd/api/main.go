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

	application := app.New(cfg, logger)
	_ = application.Run(context.Background())
}
