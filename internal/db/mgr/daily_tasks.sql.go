// Code generated by sqlc. DO NOT EDIT.
// source: daily_tasks.sql

package db

import (
	"context"
)

const createDailySyncTask = `-- name: CreateDailySyncTask :exec
INSERT INTO artist_daily_sync_tasks (user_name)
VALUES ($1)
`

func (q *Queries) CreateDailySyncTask(ctx context.Context, userName string) error {
	_, err := q.db.ExecContext(ctx, createDailySyncTask, userName)
	return err
}

const getUserDailySyncTask = `-- name: GetUserDailySyncTask :one
SELECT id, created_at, updated_at, user_name FROM artist_daily_sync_tasks
WHERE user_name = $1
LIMIT 1
`

func (q *Queries) GetUserDailySyncTask(ctx context.Context, userName string) (ArtistDailySyncTask, error) {
	row := q.db.QueryRowContext(ctx, getUserDailySyncTask, userName)
	var i ArtistDailySyncTask
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserName,
	)
	return i, err
}
