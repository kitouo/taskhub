package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kitouo/taskhub/internal/config"
	"github.com/kitouo/taskhub/internal/httpx"
	"github.com/kitouo/taskhub/internal/logx"
)

func newMux() http.Handler {

	// 创建路由
	mux := http.NewServeMux()

	// 路由注册
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "taskhub/dev\n")
	})

	mux.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		rid := httpx.RequestIDFromContext(r.Context())
		httpx.WriteError(w, http.StatusBadRequest,
			"INVALID_ARGUMENT", "bad request example", rid)
	})

	// 流量接口
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	})

	return mux
}

func main() {

	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	logger := logx.New("server=api", logx.ParseLevel(cfg.LogLevel))
	logger.Info("starting", "config="+cfg.SafeString())

	h := newMux()
	h = httpx.AccessLogger(logger, h)
	h = httpx.Recover(logger, h)
	h = httpx.WithRequestID(h) // 最外层

	srv := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      h,
		ReadTimeout:  time.Duration(cfg.ReadTimeoutSec) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeoutSec) * time.Second,
		IdleTimeout:  time.Duration(cfg.IdleTimeoutSec) * time.Second,
	}

	/*
		加入signal监听与shutdown流程

	*/
	// 启动 server 放到 goroutine
	go func() {
		logger.Info("http server starting", "addr="+srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("http server crashed", "err", err)
		}
	}()

	/*
		等待信号
		让主goroutine阻塞等待“退出信号”，一旦收到信号就继续往下执行优雅退出逻辑
	*/
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	logger.Info("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(cfg.ShutdownTimeoutSec)*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {

		logger.Warn("http server shutdown failed", "err=", err)

	} else {
		logger.Info("http server shutdown gracefully")
	}

}
