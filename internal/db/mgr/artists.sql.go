// Code generated by sqlc. DO NOT EDIT.
// source: artists.sql

package db

import (
	"context"
)

const createArtist = `-- name: CreateArtist :one
INSERT INTO "spotify_artists" (name, poster)
VALUES ($1, $2)
ON CONFLICT DO NOTHING
RETURNING id, created_at, name, poster
`

type CreateArtistParams struct {
	Name   string `json:"name"`
	Poster string `json:"poster"`
}

func (q *Queries) CreateArtist(ctx context.Context, arg CreateArtistParams) (SpotifyArtist, error) {
	row := q.db.QueryRowContext(ctx, createArtist, arg.Name, arg.Poster)
	var i SpotifyArtist
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Name,
		&i.Poster,
	)
	return i, err
}
