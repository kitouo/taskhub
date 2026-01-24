package api

import (
	"net/http"

	"github.com/kitouo/taskhub/internal/service"
)

type Router struct {
	task *TaskHandler
}

func NewRouter(svc *service.TaskService) http.Handler {

	r := &Router{
		task: NewTaskHandler(svc),
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
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
