package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kitouo/taskhub/internal/api"
	"github.com/kitouo/taskhub/internal/config"
	"github.com/kitouo/taskhub/internal/httpx"
	"github.com/kitouo/taskhub/internal/logx"
	"github.com/kitouo/taskhub/internal/repo/memory"
	"github.com/kitouo/taskhub/internal/service"
)

type App struct {
	cfg    config.Config
	logger logx.Logger
	srv    *http.Server
}

func New(cfg config.Config, logger logx.Logger) *App {

	// wire dependencies 线路依赖
	taskRepo := memory.NewTaskRepo()
	taskService := service.NewTaskService(taskRepo)
	handler := api.NewRouter(taskService)

	// middleware chain
	h := handler
	h = httpx.AccessLogger(logger, h)
	h = httpx.Recover(logger, h)
	h = httpx.WithRequestID(h)

	srv := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      h,
		ReadTimeout:  time.Duration(cfg.ReadTimeoutSec) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeoutSec) * time.Second,
		IdleTimeout:  time.Duration(cfg.IdleTimeoutSec) * time.Second,
	}

	return &App{cfg: cfg, logger: logger, srv: srv}
}

func (a *App) Run(ctx context.Context) error {
	// start server
	go func() {
		a.logger.Info("http server starting", "addr="+a.srv.Addr)
		if err := a.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Error("http server crashed", "err=", err)
		}
	}()

	// wait for signal or ctx done 监听退出信号
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(stop)

	// 等待“外部退出条件”
	select {
	case <-ctx.Done():
		a.logger.Info("shutdown by context")
	case <-stop:
		a.logger.Info("shutdown signal received")
	}

	sdCtx, cancel := context.WithTimeout(context.Background(), time.Duration(a.cfg.ShutdownTimeoutSec)*time.Second)
	defer cancel()

	if err := a.srv.Shutdown(sdCtx); err != nil {
		a.logger.Error("http server shutdown failed", "err=", err)
		return err
	}
	a.logger.Info("http server shutdown gracefully")
	return nil
}
