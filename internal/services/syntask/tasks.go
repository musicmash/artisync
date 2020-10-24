package syntask

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/musicmash/artisync/internal/db"
	"github.com/musicmash/artisync/internal/guard"
	"github.com/musicmash/artisync/internal/pipelines/syntask"
)

type Task struct {
	ID      uuid.UUID `json:"id"`
	State   string    `json:"state"`
	Details *string   `json:"details"`
}

type Mgr struct {
	// client *spotify.Client
	conn     *db.Conn
	pipeline syntask.Pipeline
}

func New(conn *db.Conn) *Mgr {
	return &Mgr{conn: conn, pipeline: syntask.New(conn)}
}

func (m *Mgr) GetSyncTaskState(ctx context.Context, id uuid.UUID) (*Task, error) {
	res, err := m.conn.GetOneTimeSyncTaskState(ctx, id)
	if err == nil {
		// active task is found
		task := Task{ID: res.ID, State: string(res.State)}
		// TODO (m.kalinin): add details
		// if res.Details.Valid {
		// 	task.Details = &res.Details.String
		// }

		return &task, nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return nil, guard.NewClientError(ErrTaskNotFound)
	}

	return nil, guard.NewInternalError(err)
}

func (m *Mgr) GetOrCreateOneTimeSyncTaskForUser(ctx context.Context, userName string, code string) (*Task, error) {
	return m.getOrCreateSyncTaskForUser(ctx, false, userName, code)
}

func (m *Mgr) GetOrCreateDailySyncTaskForUser(ctx context.Context, userName string, code string) (*Task, error) {
	return m.getOrCreateSyncTaskForUser(ctx, true, userName, code)
}

//nolint:lll
func (m *Mgr) getOrCreateSyncTaskForUser(ctx context.Context, scheduleDaily bool, userName string, code string) (*Task, error) {
	// if any sync sync task already spawned, but not finished
	res, err := m.conn.IsAnySyncTaskProcessingForUser(ctx, userName)
	if err == nil {
		task := Task{ID: res.ID, State: string(res.State)}
		if res.Details.Valid {
			task.Details = &res.Details.String
		}

		return &task, nil
	}

	// task not found
	if errors.Is(err, sql.ErrNoRows) {
		task, err := m.pipeline.Run(ctx, &syntask.PipelineOpts{
			UserName:          userName,
			SpotifyAuthCode:   code,
			ScheduleDailySync: scheduleDaily,
		})
		if err != nil {
			// TODO (m.kalinin): check if spotify code is broken?
			return nil, guard.NewInternalError(err)
		}

		result := Task(*task)
		return &result, nil //nolint:nlreturn
	}

	// another error here
	return nil, guard.NewInternalError(fmt.Errorf("can't get current tasks for %v: %w", userName, err))
}
