package sync

import (
	"context"

	"github.com/musicmash/artisync/internal/db"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

type PipelineOpts struct {
	UserName     string
	RefreshToken string
}

type Pipeline interface {
	Run(ctx context.Context, opts *PipelineOpts) error
}

type PipelineData struct {
	userName     string
	userArtists  []spotify.FullArtist
	refreshToken string
	client       spotify.Client
	auth         oauth2.Config
}

type Step func(ctx context.Context, data *PipelineData) error

type TaskPipeline struct {
	conn  *db.Conn
	auth  oauth2.Config
	steps []Step
}

func New(mgr *db.Conn, auth oauth2.Config) Pipeline {
	return &TaskPipeline{
		conn: mgr,
		auth: auth,
		steps: []Step{
			PrepareSpotifyClient,
			GetUserTopArtists,
			GetArtistsThatUserFollows,
			UniqueArtists,
			SubscribeUserOnArtists,
		},
	}
}

func (t *TaskPipeline) Run(ctx context.Context, opts *PipelineOpts) error {
	data := &PipelineData{
		userName:     opts.UserName,
		refreshToken: opts.RefreshToken,
		auth:         t.auth,
	}
	for i := range t.steps {
		if err := t.steps[i](ctx, data); err != nil {
			return err
		}
	}
	return nil
}
