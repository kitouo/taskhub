package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kitouo/taskhub/internal/api"
	"github.com/kitouo/taskhub/internal/config"
	"github.com/kitouo/taskhub/internal/db"
	"github.com/kitouo/taskhub/internal/httpx"
	"github.com/kitouo/taskhub/internal/logx"
	"github.com/kitouo/taskhub/internal/repo"
	"github.com/kitouo/taskhub/internal/repo/memory"
	mysqlrepo "github.com/kitouo/taskhub/internal/repo/mysql"
	"github.com/kitouo/taskhub/internal/service"
)

type App struct {
	cfg    config.Config
	logger logx.Logger
	srv    *http.Server
	/*
		closeFunc用于释放外部资源（例如MySQL连接池）
		memory模式下可以是nil或空函数
	*/
	closeFunc func() error
}

func New(cfg config.Config, logger logx.Logger) (*App, error) {

	// wire dependencies 线路依赖
	var taskRepo repo.TaskRepo
	/*
		readyCheck：注入到 router，用于 /readyz
			- memory：nil（默认 ok）
			- mysql ：Ping MySQL
	*/
	var readyCheck func(context.Context) error
	var closeFunc func() error

	switch cfg.RepoMode {
	case "memory":
		// 内存模式：无外部依赖，启动永远 ready
		taskRepo = memory.NewTaskRepo()
		readyCheck = nil
		closeFunc = nil
	case "mysql":
		// MySQL模式：需要DB_DSN，否则db.Open会返回error
		dbConn, err := db.Open(cfg.DBDriver, cfg.DBDNS)
		if err != nil {
			return nil, err
		}

		// 如果后续任何步骤失败，记得及时Close，避免资源泄露
		if err := db.MigrateMySQL(dbConn); err != nil {
			_ = dbConn.Close()
			return nil, err
		}

		//注入readiness：/readyz会在1s超时内ping一次DB
		readyCheck = func(ctx context.Context) error {
			return db.Ping(ctx, dbConn)
		}

		// 退出时关闭连接池
		closeFunc = dbConn.Close

		//使用MySQL repo实现
		taskRepo = mysqlrepo.NewTaskRepo(dbConn)
	default:
		return nil, fmt.Errorf("unsupported REPO_MODE: %s", cfg.RepoMode)
	}

	taskSvc := service.NewTaskService(taskRepo)

	handler := api.NewRouter(taskSvc, readyCheck)

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

	return &App{
		cfg:       cfg,
		logger:    logger,
		srv:       srv,
		closeFunc: closeFunc,
	}, nil
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

	// 释放外部资源
	if a.closeFunc != nil {
		if err := a.closeFunc(); err != nil {
			a.logger.Error("close resources failed", "err=", err)
		}
	}
	return nil
}
