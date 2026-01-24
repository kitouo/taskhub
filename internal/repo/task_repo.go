package repo

import (
	"context"

	"github.com/kitouo/taskhub/internal/model"
)

type TaskRepo interface {
	Create(ctx context.Context, t model.Task) (model.Task, error)
	List(ctx context.Context) ([]model.Task, error)
	Get(ctx context.Context, id string) (model.Task, bool, error)
	MarkDone(ctx context.Context, id string, done bool) (model.Task, bool, error)
}
