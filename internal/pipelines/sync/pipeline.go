package sync

import (
	"context"

	"github.com/musicmash/artisync/internal/db"
	"github.com/zmb3/spotify"
)

type PipelineOpts struct {
	UserName     string
	RefreshToken string
}

type Pipeline interface {
	Run(ctx context.Context, opts *PipelineOpts) error
}

type PipelineData struct {
	userName    string
	userArtists []string
	client      *spotify.Client
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
			PrepareSpotifyClient,
			GetUserArtists,
			EnsureUserArtistsExists,
			Subscribe,
		},
	}
}

func (t *TaskPipeline) Run(ctx context.Context, opts *PipelineOpts) error {
	data := &PipelineData{
		userName: opts.UserName,
	}
	for i := range t.steps {
		if err := t.steps[i](ctx, data); err != nil {
			return err
		}
	}
	return nil
}
