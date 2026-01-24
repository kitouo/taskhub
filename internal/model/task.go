package model

import "time"

type Task struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	Done     bool      `json:"done"`
	CreateAt time.Time `json:"create_at"`
}
