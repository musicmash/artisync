package sync

import "context"

func PrepareSpotifyClient(ctx context.Context, data *PipelineData) error {
	// exchange access token and init client
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
