package sync

import (
	"context"
	"fmt"
	"time"

	"github.com/musicmash/artisync/internal/log"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

func PrepareSpotifyClient(ctx context.Context, data *PipelineData) error {
	ts := data.auth.TokenSource(context.Background(), &oauth2.Token{
		AccessToken:  "fake-access-token",
		Expiry:       time.Now().UTC().AddDate(-1, 0, 0),
		RefreshToken: data.refreshToken,
	})

	token, err := ts.Token()
	if err != nil {
		return fmt.Errorf("can't reissue access-token: %w", err)
	}

	log.Info("access-token was successfully reissued")
	data.client = spotify.NewClient(data.auth.Client(ctx, token))
	return nil
}

func GetUserArtists(ctx context.Context, data *PipelineData) error {
	// get user's top artists
	// get artists followed by user
	// make unique slice of artist id's
	return nil
}

func EnsureUserArtistsExists(ctx context.Context, data *PipelineData) error {
	// iterate over artists:
	//   get association with artist_id
	//   check if association not exists
	//     create artist, association
	return nil
}

func Subscribe(ctx context.Context, data *PipelineData) error {
	// subscribe user on artists
	return nil
}
