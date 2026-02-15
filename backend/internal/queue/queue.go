package queue

import (
	"context"
)

type IQueue interface {
	Enqueue(ctx context.Context, task *Task) error
}

type Task struct {
	ID      string      `json:"id"`
	Version int         `json:"version"`
	Type    string      `json:"type"`
	Data    interface{} `json:"data"` // VlogInput または map[string]interface{}
	Status  string      `json:"status"`
}
