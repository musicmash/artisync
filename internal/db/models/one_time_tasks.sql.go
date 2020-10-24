// Code generated by sqlc. DO NOT EDIT.
// source: one_time_tasks.sql

package models

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createOneTimeSyncTask = `-- name: CreateOneTimeSyncTask :one
INSERT INTO "artist_one_time_sync_tasks" (user_name)
VALUES ($1)
RETURNING id, created_at, updated_at, user_name, state, details
`

func (q *Queries) CreateOneTimeSyncTask(ctx context.Context, userName string) (ArtistOneTimeSyncTask, error) {
	row := q.db.QueryRowContext(ctx, createOneTimeSyncTask, userName)
	var i ArtistOneTimeSyncTask
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserName,
		&i.State,
		&i.Details,
	)
	return i, err
}

const getNextScheduledTask = `-- name: GetNextScheduledTask :one
SELECT
    artist_one_time_sync_tasks.id,
    artist_one_time_sync_tasks.user_name,
    artist_sync_refresh_tokens.value as refresh_token
FROM artist_one_time_sync_tasks
LEFT JOIN artist_sync_refresh_tokens ON (
    artist_sync_refresh_tokens.user_name=artist_one_time_sync_tasks.user_name
)
WHERE state = 'scheduled' AND artist_sync_refresh_tokens.expired_at > now()
LIMIT 1
`

type GetNextScheduledTaskRow struct {
	ID           uuid.UUID `json:"id"`
	UserName     string    `json:"user_name"`
	RefreshToken string    `json:"refresh_token"`
}

func (q *Queries) GetNextScheduledTask(ctx context.Context) (GetNextScheduledTaskRow, error) {
	row := q.db.QueryRowContext(ctx, getNextScheduledTask)
	var i GetNextScheduledTaskRow
	err := row.Scan(&i.ID, &i.UserName, &i.RefreshToken)
	return i, err
}

const getOneTimeSyncTaskState = `-- name: GetOneTimeSyncTaskState :one
SELECT id, state FROM "artist_one_time_sync_tasks"
WHERE id = $1
LIMIT 1
`

type GetOneTimeSyncTaskStateRow struct {
	ID    uuid.UUID `json:"id"`
	State TaskState `json:"state"`
}

func (q *Queries) GetOneTimeSyncTaskState(ctx context.Context, id uuid.UUID) (GetOneTimeSyncTaskStateRow, error) {
	row := q.db.QueryRowContext(ctx, getOneTimeSyncTaskState, id)
	var i GetOneTimeSyncTaskStateRow
	err := row.Scan(&i.ID, &i.State)
	return i, err
}

const isAnySyncTaskProcessingForUser = `-- name: IsAnySyncTaskProcessingForUser :one
SELECT id, created_at, updated_at, user_name, state, details FROM "artist_one_time_sync_tasks"
WHERE user_name = $1 AND state NOT IN ('done', 'error')
LIMIT 1
`

func (q *Queries) IsAnySyncTaskProcessingForUser(ctx context.Context, userName string) (ArtistOneTimeSyncTask, error) {
	row := q.db.QueryRowContext(ctx, isAnySyncTaskProcessingForUser, userName)
	var i ArtistOneTimeSyncTask
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserName,
		&i.State,
		&i.Details,
	)
	return i, err
}

const updateOneTimeSyncTaskState = `-- name: UpdateOneTimeSyncTaskState :exec
UPDATE artist_one_time_sync_tasks
SET state = $1, details = $2
WHERE id = $3
`

type UpdateOneTimeSyncTaskStateParams struct {
	State   TaskState      `json:"state"`
	Details sql.NullString `json:"details"`
	ID      uuid.UUID      `json:"id"`
}

func (q *Queries) UpdateOneTimeSyncTaskState(ctx context.Context, arg UpdateOneTimeSyncTaskStateParams) error {
	_, err := q.db.ExecContext(ctx, updateOneTimeSyncTaskState, arg.State, arg.Details, arg.ID)
	return err
}