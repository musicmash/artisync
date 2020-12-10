// Code generated by sqlc. DO NOT EDIT.
// source: subscriptions.sql

package models

import (
	"context"
)

const createSubscription = `-- name: CreateSubscription :exec
INSERT INTO subscriptions (created_at, user_name, artist_id)
VALUES (now(), $1, $2)
ON CONFLICT  DO NOTHING
`

type CreateSubscriptionParams struct {
	UserName string `json:"user_name"`
	ArtistID int64  `json:"artist_id"`
}

func (q *Queries) CreateSubscription(ctx context.Context, arg CreateSubscriptionParams) error {
	_, err := q.db.ExecContext(ctx, createSubscription, arg.UserName, arg.ArtistID)
	return err
}
