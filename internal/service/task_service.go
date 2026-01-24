package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/kitouo/taskhub/internal/model"
	"github.com/kitouo/taskhub/internal/repo"
)

var ErrInvalidTitle = errors.New("invalid title")

type TaskService struct {
	repo repo.TaskRepo
}

func NewTaskService(repo repo.TaskRepo) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) Create(ctx context.Context, title string) (model.Task, error) {
	title = strings.TrimSpace(title)
	if title == "" || len(title) > 200 {
		return model.Task{}, ErrInvalidTitle
	}
	t := model.Task{
		ID:       NewID(),
		Title:    title,
		Done:     false,
		CreateAt: time.Now().UTC(),
	}
	return s.repo.Create(ctx, t)
}

func (s *TaskService) List(ctx context.Context) ([]model.Task, error) {
	return s.repo.List(ctx)
}

func (s *TaskService) Get(ctx context.Context, id string) (model.Task, bool, error) {
	return s.repo.Get(ctx, id)
}

func (s *TaskService) MarkDone(ctx context.Context, id string, done bool) (model.Task, bool, error) {
	return s.repo.MarkDone(ctx, id, done)
}

func NewID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
