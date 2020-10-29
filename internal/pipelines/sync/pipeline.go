package sync

import (
	"context"

	"github.com/musicmash/artisync/internal/db"
)

type PipelineOpts struct {
	UserName     string
	RefreshToken string
}

type Pipeline interface {
	Run(ctx context.Context, opts *PipelineOpts) error
}

type PipelineData struct {
	UserName    string
	UserArtists []string
}

type Step func(ctx context.Context, data *PipelineData) error

type TaskPipeline struct {
	conn  *db.Conn
	steps []Step
}

func New(mgr *db.Conn) Pipeline {
	return &TaskPipeline{
		conn: mgr,
		steps: []Step{
			GetUserArtists,
			EnsureUserArtistsExists,
			Subscribe,
		},
	}
}

func (t *TaskPipeline) Run(ctx context.Context, opts *PipelineOpts) error {
	data := &PipelineData{
		UserName: opts.UserName,
	}
	for i := range t.steps {
		if err := t.steps[i](ctx, data); err != nil {
			return err
		}
	}
	return nil
}
