package tasks

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
)

type Task struct {
	ID      uuid.UUID       `json:"id"`
	State   string          `json:"state"`
	Details json.RawMessage `json:"details"`
}

type Repository interface {
	GetTask(ctx context.Context, id uuid.UUID) (*Task, error)
}
