package sync

import (
	"context"
	"errors"
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

func GetUserTopArtists(ctx context.Context, data *PipelineData) error {
	artists := []spotify.FullArtist{}
	results, err := data.client.CurrentUsersTopArtists()
	if err != nil {
		return fmt.Errorf("can't get user top artists: %w", err)
	}

	for page := 1; ; page++ {
		artists = append(artists, results.Artists...)
		if results.Total == len(artists) {
			break
		}

		err = data.client.NextPage(results)
		if errors.Is(err, spotify.ErrNoMorePages) {
			break
		}
		if err != nil {
			return fmt.Errorf("can't get next page with user top artists: %w", err)
		}
	}

	log.Infof("got %d top artists for user", len(artists))
	data.userArtists = artists
	return nil
}

func GetArtistsThatUserFollows(ctx context.Context, data *PipelineData) error {
	artists := []spotify.FullArtist{}
	results, err := data.client.CurrentUsersFollowedArtists()
	if err != nil {
		return fmt.Errorf("can't get artists that user follows: %w", err)
	}

	for {
		if len(results.Artists) == 0 {
			break
		}

		artists = append(artists, results.Artists...)
		if results.Total == len(artists) {
			break
		}

		after := results.Artists[len(results.Artists)-1].ID.String()
		results, err = data.client.CurrentUsersFollowedArtistsOpt(50, after)
		if err != nil {
			return fmt.Errorf("cen't get next page of followed artists: %w", err)
		}
	}

	log.Infof("got %d artists that user follows", len(artists))
	data.userArtists = append(data.userArtists, artists...)
	return nil
}

func UniqueArtists(ctx context.Context, data *PipelineData) error {
	uniqueArtists := []spotify.FullArtist{}
	uniqueIDs := make(map[string]struct{}, len(data.userArtists))
	for i := range data.userArtists {
		artistID := data.userArtists[i].ID.String()
		if _, exists := uniqueIDs[artistID]; exists {
			continue
		}

		uniqueIDs[artistID] = struct{}{}
		uniqueArtists = append(uniqueArtists, data.userArtists[i])
	}

	log.Infof("got %d unique artists from %d artists", len(uniqueArtists), len(data.userArtists))
	data.userArtists = uniqueArtists
	return nil
}

func EnsureUserArtistsExists(ctx context.Context, data *PipelineData) error {
	// iterate over artists:
	//   get association with artist_id
	//   check if association not exists
	//     create artist, association
	return nil
}

func SubscribeUserOnArtists(ctx context.Context, data *PipelineData) error {
	// subscribe user on artists
	return nil
}
