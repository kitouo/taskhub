package memory

import (
	"context"
	"sync"

	"github.com/kitouo/taskhub/internal/model"
)

type TaskRepo struct {
	mu    sync.RWMutex
	byID  map[string]model.Task
	order []string
}

func NewTaskRepo() *TaskRepo {
	return &TaskRepo{
		byID: make(map[string]model.Task),
	}
}

func (r *TaskRepo) Create(ctx context.Context, task model.Task) (model.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.byID[task.ID] = task
	r.order = append(r.order, task.ID)
	return task, nil
}

func (r *TaskRepo) List(ctx context.Context) ([]model.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]model.Task, 0, len(r.order))
	for _, id := range r.order {
		out = append(out, r.byID[id])
	}

	return out, nil
}

func (r *TaskRepo) Get(ctx context.Context, id string) (model.Task, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	task, ok := r.byID[id]
	return task, ok, nil
}

func (r *TaskRepo) MarkDone(ctx context.Context, id string, done bool) (model.Task, bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	task, ok := r.byID[id]

	if !ok {
		return model.Task{}, false, nil
	}

	task.Done = done
	r.byID[id] = task
	return task, true, nil
}
