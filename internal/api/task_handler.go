package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/kitouo/taskhub/internal/httpx"
	"github.com/kitouo/taskhub/internal/service"
)

type TaskHandler struct {
	svc *service.TaskService
}

func NewTaskHandler(svc *service.TaskService) *TaskHandler {
	return &TaskHandler{svc: svc}
}

type createTaskRequest struct {
	Title string `json:"title"`
}

type markDoneRequest struct {
	Done bool `json:"done"`
}

/*
HandleTasks /tasks: GET list, POST create
*/
func (h *TaskHandler) HandleTasks(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		tasks, err := h.svc.List(r.Context())
		if err != nil {
			h.writeInternal(w, r)
			return
		}
		httpx.WriteJson(w, http.StatusOK, tasks)
		return
	case http.MethodPost:
		var req createTaskRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.writeBadRequest(w, r, "INVALID_JSON", "invalid json body")
			return
		}
		t, err := h.svc.Create(r.Context(), req.Title)
		if err == service.ErrInvalidTitle {
			h.writeBadRequest(w, r, "INVALID_ARGUMENT", "title is required (<= 200)")
			return
		}
		if err != nil {
			h.writeInternal(w, r)
			return
		}
		httpx.WriteJson(w, http.StatusCreated, t)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

/*
HandleTaskByID tasks/{id} or /tasks/{id}/done : GET, PATCH
*/
func (h *TaskHandler) HandleTaskByID(w http.ResponseWriter, r *http.Request) {
	// 从请求URL里把/tasks/这一段去掉，然后把剩下的部分两端的/去掉，得到“纯净的资源标识（id）或子路径”。
	path := strings.TrimPrefix(r.URL.Path, "/tasks/")
	path = strings.Trim(path, "/")
	if path == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	parts := strings.Split(path, "/")
	fmt.Println(parts)
	id := parts[0]

	// GET /tasks/{id}
	if r.Method == http.MethodGet && len(parts) == 1 {
		t, ok, err := h.svc.Get(r.Context(), id)
		if err != nil {
			h.writeInternal(w, r)
			return
		}
		if !ok {
			h.writeNotFound(w, r, "NOT_FOUND", "task not found")
			return
		}
		httpx.WriteJson(w, http.StatusOK, t)
		return
	}

	// PATCH /tasks/{id}  (body: {"done":true})
	if r.Method == http.MethodPatch && len(parts) == 1 {
		var req markDoneRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.writeBadRequest(w, r, "INVALID_JSON", "invalid json body")
			return
		}

		t, ok, err := h.svc.MarkDone(r.Context(), id, req.Done)
		if err != nil {
			h.writeInternal(w, r)
			return
		}
		if !ok {
			h.writeNotFound(w, r, "NOT_FOUND", "task not found")
			return
		}
		httpx.WriteJson(w, http.StatusOK, t)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *TaskHandler) writeInternal(w http.ResponseWriter, r *http.Request) {
	rid := httpx.RequestIDFromContext(r.Context())
	httpx.WriteError(w, http.StatusInternalServerError, "INTERNAL", "internal server error", rid)
}

func (h *TaskHandler) writeBadRequest(w http.ResponseWriter, r *http.Request, code, msg string) {
	rid := httpx.RequestIDFromContext(r.Context())
	httpx.WriteError(w, http.StatusBadRequest, code, msg, rid)
}
func (h *TaskHandler) writeNotFound(w http.ResponseWriter, r *http.Request, code, msg string) {
	rid := httpx.RequestIDFromContext(r.Context())
	httpx.WriteError(w, http.StatusNotFound, code, msg, rid)
}
