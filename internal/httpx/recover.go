package httpx

import (
	"net/http"
	"runtime/debug"

	"github.com/kitouo/taskhub/internal/logx"
)

/*
把错误（panic）从“直接崩溃进程/断开连接”变成“可控的500响应+可排障日志”，从而提升服务稳定性与可运维性
*/
func Recover(logger logx.Logger, next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				rid := RequestIDFromContext(r.Context())

				logger.Error("panic",
					"request_id="+rid,
					"recover=", rec,
					"stack", string(debug.Stack()),
				)

				WriteError(w,
					http.StatusInternalServerError, "INTERNAL",
					"internal server error", rid)

			}
		}()

		next.ServeHTTP(w, r)

	})

}
