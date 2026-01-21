package main

import (
	"fmt"
	"net/http"

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
		_, _ = w.Write([]byte("OK"))
	})

	mux.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "taskhub/dev\n")
	})

	return mux
}

func main() {

	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	logger := logx.New("service=api")
	logger.Info("starting", "config="+cfg.SafeString())

	h := newMux()
	h = httpx.AccessLogger(logger, h)
	h = httpx.WithRequestID(h) // 最外层

	addr := ":" + cfg.HTTPPort
	if err := http.ListenAndServe(addr, h); err != nil {
		logger.Error("server stopped", "err=", err)
	}

}
