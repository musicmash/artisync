package sync

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/musicmash/artisync/internal/db/models"
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
	var (
		timeRange = "medium"
		limit     = 50
	)
	opts := spotify.Options{Timerange: &timeRange, Limit: &limit}

	artists := []spotify.FullArtist{}
	results, err := data.client.CurrentUsersTopArtistsOpt(&opts)
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
	var (
		after = ""
		limit = 50
	)

	artists := []spotify.FullArtist{}
	results, err := data.client.CurrentUsersFollowedArtistsOpt(limit, after)
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

		after = results.Artists[len(results.Artists)-1].ID.String()
		results, err = data.client.CurrentUsersFollowedArtistsOpt(limit, after)
		if err != nil {
			return fmt.Errorf("can't get next page with artists that user follows: %w", err)
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

func getPoster(images []spotify.Image) string {
	if len(images) == 0 {
		return ""
	}

	return images[0].URL
}

//nolint:lll
func createArtist(ctx context.Context, db *models.Queries, fullArtist *spotify.FullArtist) (models.SpotifyArtist, error) {
	poster := getPoster(fullArtist.Images)
	return db.CreateArtist(ctx, models.CreateArtistParams{
		Name: fullArtist.Name,
		Poster: sql.NullString{
			String: poster,
			Valid:  len(poster) > 0,
		},
	})
}

//nolint:lll
func createAssociation(ctx context.Context, db *models.Queries, artistID int64, storeID string) (models.SpotifyArtistAssociation, error) {
	return db.CreateArtistAssociation(ctx, models.CreateArtistAssociationParams{
		ArtistID:  artistID,
		StoreID:   storeID,
		StoreName: "spotify",
	})
}

func SubscribeUserOnArtists(ctx context.Context, data *PipelineData) error {
	err := data.mashDB.ExecTx(ctx, func(db *models.Queries) error {
		for i := range data.userArtists {
			artistStoreID := data.userArtists[i].ID.String()
			association, err := db.GetArtistAssociation(ctx, models.GetArtistAssociationParams{
				StoreID:   artistStoreID,
				StoreName: "spotify",
			})
			switch {
			case errors.Is(err, sql.ErrNoRows):
				artist, err := createArtist(ctx, db, &data.userArtists[i])
				if err != nil {
					return fmt.Errorf("can't create artist with spotify_id %s: %w", artistStoreID, err)
				}

				association, err = createAssociation(ctx, db, artist.ID, artistStoreID)
				if err != nil {
					return fmt.Errorf("can't create association for artist with spotify_id %s: %w", artistStoreID, err)
				}
			case err != nil:
				return fmt.Errorf("cant get association for artist with spotify_id: %s: %w", artistStoreID, err)
			}

			err = db.CreateSubscription(ctx, models.CreateSubscriptionParams{
				UserName: data.userName,
				ArtistID: association.ArtistID,
			})
			if err != nil {
				return fmt.Errorf("cant subscribe user on artist with spotify_id: %s: %w", artistStoreID, err)
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("can't commit tx: %w", err)
	}

	return nil
}
