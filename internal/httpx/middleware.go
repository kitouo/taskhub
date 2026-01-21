package httpx

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/kitouo/taskhub/internal/logx"
)

// 避免字符串撞名
type ctxKey string

const requestIDKey ctxKey = "request_id"

// RequestIDFromContext 从context.Context里安全地取出request_id(请求链路id)
func RequestIDFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(requestIDKey).(string); ok {
		return v
	}
	return ""
}

// WithRequestID 给每个请求“绑定”一个request_id，并且让它贯穿整条请求链路：响应头、context、日志都能拿到同一个ID
func WithRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// 从请求头里尝试读取X-Request-Id
		rid := r.Header.Get("X-Request-ID")

		// 如果请求头没带request_id，就自己生成一个，保证每个请求必定有rid
		if rid == "" {
			rid = newRequestID()
		}

		/*
			把request_id写到响应头
			调用方（浏览器、前端、其他服务）能在响应里看到这个id
			用户报错时能把request_id给你，你就能在日志中精准定位
		*/
		w.Header().Set("X-Request-ID", rid)

		// 把request_id写进context，让后续所有处理代码都能通过r.Context()拿到同一个request_id
		ctx := context.WithValue(r.Context(), requestIDKey, rid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AccessLogger(logger logx.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := &wrapWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		next.ServeHTTP(ww, r)

		// 耗时
		lat := time.Since(start).Milliseconds()

		// AccessLog从响应头读取request_id
		rid := ww.Header().Get("X-Request-Id")
		if rid == "" {
			rid = RequestIDFromContext(r.Context())
		}

		logger.Info("request",
			"request_id="+rid,
			"method="+r.Method,
			"path="+r.URL.Path,
			"status=", ww.status,
			"latency_ms=", lat,
		)

	})
}

type wrapWriter struct {
	http.ResponseWriter
	status int
}

func (w *wrapWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// 生成新的request_id
func newRequestID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
