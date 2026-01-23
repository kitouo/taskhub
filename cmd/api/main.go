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

	addr := ":" + cfg.HTTPPort
	if err := http.ListenAndServe(addr, h); err != nil {
		logger.Error("server stopped", "err=", err)
	}

}
