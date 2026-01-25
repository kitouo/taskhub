package mysqlrepo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/kitouo/taskhub/internal/model"
)

type TaskRepo struct {
	db *sql.DB
}

func NewTaskRepo(db *sql.DB) *TaskRepo {
	return &TaskRepo{db: db}
}

func (r *TaskRepo) Create(ctx context.Context, t model.Task) (model.Task, error) {
	// MySQL使用TINYINT(1)表示布尔
	doneInt := 0
	if t.Done {
		doneInt = 1
	}

	createdAt := t.CreatedAt.UTC()

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO tasks(id, title, done, created_at) VALUES (?, ?, ?, ?)`,
		t.ID, t.Title, doneInt, createdAt,
	)
	if err != nil {
		return model.Task{}, fmt.Errorf("insert task: %w", err)
	}
	return t, nil
}

func (r *TaskRepo) List(ctx context.Context) ([]model.Task, error) {
	// ORDER BY created_at 使列表稳定
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, title, done, created_at FROM tasks ORDER BY created_at ASC`,
	)

	if err != nil {
		return nil, fmt.Errorf("query tasks: %w", err)
	}
	defer rows.Close()

	out := make([]model.Task, 0)
	for rows.Next() {
		var (
			t       model.Task
			doneInt int
			ct      time.Time
		)
		if err := rows.Scan(&t.ID, &t.Title, &doneInt, &ct); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		t.Done = doneInt == 1
		t.CreatedAt = ct.UTC()
		out = append(out, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err: %w", err)
	}
	return out, nil
}

func (r *TaskRepo) Get(ctx context.Context, id string) (model.Task, bool, error) {

	row := r.db.QueryRowContext(ctx,
		`SELECT id, title, done, created_at FROM tasks WHERE id = ?`,
		id,
	)

	var (
		t       model.Task
		doneInt int
		ct      time.Time
	)
	if err := row.Scan(&t.ID, &t.Title, &doneInt, &ct); err != nil {
		if err == sql.ErrNoRows {
			return model.Task{}, false, nil
		}
		return model.Task{}, false, fmt.Errorf("get task: %w", err)
	}

	t.Done = doneInt == 1
	t.CreatedAt = ct.UTC()
	return t, true, nil
}

func (r *TaskRepo) MarkDone(ctx context.Context, id string, done bool) (model.Task, bool, error) {
	doneInt := 0
	if done {
		doneInt = 1
	}

	res, err := r.db.ExecContext(ctx,
		`UPDATE tasks SET done = ? WHERE id = ?`,
		doneInt, id,
	)
	if err != nil {
		return model.Task{}, false, fmt.Errorf("update task: %w", err)
	}

	aff, err := res.RowsAffected()
	if err != nil {
		return model.Task{}, false, fmt.Errorf("rows affected: %w", err)
	}
	if aff == 0 {
		return model.Task{}, false, nil
	}

	// 更新后再读一次，保持返回结构与Get一致
	return r.Get(ctx, id)
}
