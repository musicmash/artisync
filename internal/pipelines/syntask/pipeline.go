package syntask

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/musicmash/artisync/internal/db"
	"golang.org/x/oauth2"
)

type Task struct {
	ID      uuid.UUID       `json:"id"`
	State   string          `json:"state"`
	Details json.RawMessage `json:"details"`
}

type PipelineOpts struct {
	auth              *oauth2.Config
	UserName          string
	SpotifyAuthCode   string
	ScheduleDailySync bool
}

type Pipeline interface {
	Run(ctx context.Context, opts *PipelineOpts) (*Task, error)
}

type PipelineData struct {
	auth  *oauth2.Config
	task  *Task
	conn  *db.Conn
	token *oauth2.Token
}

type Step func(ctx context.Context, opts *PipelineOpts, data *PipelineData) error

type TaskPipeline struct {
	auth  *oauth2.Config
	conn  *db.Conn
	steps []Step
}

func New(auth *oauth2.Config, mgr *db.Conn) Pipeline {
	return &TaskPipeline{
		auth: auth,
		conn: mgr,
		steps: []Step{
			GetRefreshTokenStep,
			ScheduleSyncTaskStep,
		},
	}
}

func (t *TaskPipeline) Run(ctx context.Context, opts *PipelineOpts) (*Task, error) {
	data := &PipelineData{auth: t.auth, conn: t.conn}
	for i := range t.steps {
		if err := t.steps[i](ctx, opts, data); err != nil {
			return nil, err
		}
	}

	if data.task == nil {
		return nil, ErrInternalEmptyTask
	}

	return data.task, nil
}
