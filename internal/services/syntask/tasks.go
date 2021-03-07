package syntask

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/musicmash/artisync/internal/db"
	"github.com/musicmash/artisync/internal/guard"
	"github.com/musicmash/artisync/internal/pipelines/syntask"
)

type Task struct {
	ID      uuid.UUID       `json:"id"`
	State   string          `json:"state"`
	Details json.RawMessage `json:"details"`
}

type Mgr struct {
	conn     *db.Conn
	pipeline syntask.Pipeline
}

func New(conn *db.Conn, pipeline syntask.Pipeline) *Mgr {
	return &Mgr{conn: conn, pipeline: pipeline}
}

func (m *Mgr) GetSyncTaskState(ctx context.Context, id uuid.UUID) (*Task, error) {
	res, err := m.conn.GetOneTimeSyncTaskState(ctx, id)
	if err == nil {
		// active task is found
		task := Task{ID: res.ID, State: string(res.State), Details: res.Details}
		return &task, nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return nil, guard.NewClientError(ErrTaskNotFound)
	}

	return nil, guard.NewInternalError(err)
}

func (m *Mgr) GetOrCreateOneTimeSyncTaskForUser(ctx context.Context, userName string, code string) (*Task, error) {
	const scheduledDaily = false

	return m.getOrCreateSyncTaskForUser(ctx, scheduledDaily, userName, code)
}

func (m *Mgr) GetOrCreateDailySyncTaskForUser(ctx context.Context, userName string, code string) (*Task, error) {
	const scheduledDaily = true

	return m.getOrCreateSyncTaskForUser(ctx, scheduledDaily, userName, code)
}

//nolint:lll
func (m *Mgr) getOrCreateSyncTaskForUser(ctx context.Context, scheduleDaily bool, userName string, code string) (*Task, error) {
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
	return &result, nil
}
