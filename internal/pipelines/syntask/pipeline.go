package syntask

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/musicmash/artisync/internal/db"
	"github.com/musicmash/artisync/internal/log"
)

type Task struct {
	ID      uuid.UUID       `json:"id"`
	State   string          `json:"state"`
	Details json.RawMessage `json:"details"`
}

type PipelineOpts struct {
	UserName          string
	SpotifyAuthCode   string
	ScheduleDailySync bool
}

type Pipeline interface {
	Run(ctx context.Context, opts *PipelineOpts) (*Task, error)
}

type PipelineData struct {
	task         *Task
	conn         *db.Conn
	refreshToken string
}

type Step func(ctx context.Context, opts *PipelineOpts, data *PipelineData) error

type TaskPipeline struct {
	// client *spotify.Client
	conn  *db.Conn
	steps []Step
}

func New(mgr *db.Conn) Pipeline {
	return &TaskPipeline{
		conn: mgr,
		steps: []Step{
			GetRefreshTokenStep,
			ScheduleSyncTaskStep,
		},
	}
}

func (t *TaskPipeline) Run(ctx context.Context, opts *PipelineOpts) (*Task, error) {
	data := &PipelineData{conn: t.conn}
	for i := range t.steps {
		if err := t.steps[i](ctx, opts, data); err != nil {
			return nil, err
		}
	}

	if data.task == nil {
		log.Error(ErrInternalEmptyTask.Error())
		return nil, ErrInternalEmptyTask
	}

	return data.task, nil
}
