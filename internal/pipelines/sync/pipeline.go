package sync

import (
	"context"

	"github.com/musicmash/artisync/internal/db"
)

type Pipeline interface {
	Run(ctx context.Context) error
}

type PipelineData struct {
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

func (t *TaskPipeline) Run(ctx context.Context) error {
	data := &PipelineData{}
	for i := range t.steps {
		if err := t.steps[i](ctx, data); err != nil {
			return err
		}
	}
	return nil
}
