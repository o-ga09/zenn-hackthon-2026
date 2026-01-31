package queue

import (
	"context"

	"github.com/o-ga09/zenn-hackthon-2026/internal/agent"
)

type IQueue interface {
	Enqueue(ctx context.Context, task *Task) error
	Dequeue(ctx context.Context) (*Task, error)
}

type Task struct {
	ID     string           `json:"id"`
	Type   string           `json:"type"`
	Data   *agent.VlogInput `json:"data"`
	Status string           `json:"status"`
}
