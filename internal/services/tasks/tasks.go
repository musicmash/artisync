package tasks

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/musicmash/artisync/internal/db"
	"github.com/musicmash/artisync/internal/db/models"
	"github.com/musicmash/artisync/internal/guard"
	"github.com/musicmash/artisync/internal/repository/tasks"
)

type Service struct {
	conn *db.Conn
}

func New(conn *db.Conn) *Service {
	return &Service{conn: conn}
}

func convertResultToRepoTask(res models.GetOneTimeSyncTaskStateRow) *tasks.Task {
	task := tasks.Task{
		ID:      res.ID,
		State:   string(res.State),
		Details: res.Details,
	}

	return &task
}

func (s *Service) GetTask(ctx context.Context, id uuid.UUID) (*tasks.Task, error) {
	res, err := s.conn.GetOneTimeSyncTaskState(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, guard.NewClientError(ErrTaskNotFound)
	}
	if err != nil {
		return nil, guard.NewInternalError(err)
	}

	return convertResultToRepoTask(res), nil
}
