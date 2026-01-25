package api

import (
	"context"
	"net/http"
	"time"

	"github.com/kitouo/taskhub/internal/service"
)

type Router struct {
	task       *TaskHandler
	readyCheck func(context.Context) error
}

func NewRouter(svc *service.TaskService, readyCheck func(context.Context) error) http.Handler {

	r := &Router{
		task:       NewTaskHandler(svc),
		readyCheck: readyCheck,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", r.healthz)
	mux.HandleFunc("/readyz", r.readyz)

	// tasks
	mux.HandleFunc("/tasks", r.task.HandleTasks)     // GET/POST
	mux.HandleFunc("/tasks/", r.task.HandleTaskByID) // GET/PATCH

	return mux
}

func (r *Router) healthz(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func (r *Router) readyz(w http.ResponseWriter, req *http.Request) {
	/*
		如果注入了readyCheck，说明存在外部以来
		这里设置一个很短的超时，避免/readyz被卡死
	*/
	if r.readyCheck != nil {
		ctx, cancel := context.WithTimeout(req.Context(), 1*time.Second)
		defer cancel()

		if err := r.readyCheck(ctx); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte("not ready"))
			return
		}
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
