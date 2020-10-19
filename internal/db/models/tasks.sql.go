// Code generated by sqlc. DO NOT EDIT.
// source: tasks.sql

package models

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createSyncTask = `-- name: CreateSyncTask :exec
INSERT INTO "artist_once_sync_tasks" (user_name)
VALUES ($1)
`

func (q *Queries) CreateSyncTask(ctx context.Context, userName string) error {
	_, err := q.db.ExecContext(ctx, createSyncTask, userName)
	return err
}

const getNextScheduledTask = `-- name: GetNextScheduledTask :one
SELECT
    artist_once_sync_tasks.id,
    artist_once_sync_tasks.user_name,
    artist_sync_refresh_tokens.value as refresh_token
FROM artist_once_sync_tasks
LEFT JOIN artist_sync_refresh_tokens ON (
    artist_sync_refresh_tokens.user_name=artist_once_sync_tasks.user_name
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

const getSyncTaskState = `-- name: GetSyncTaskState :one
SELECT id, state FROM "artist_once_sync_tasks"
WHERE id = $1
LIMIT 1
`

type GetSyncTaskStateRow struct {
	ID    uuid.UUID `json:"id"`
	State TaskState `json:"state"`
}

func (q *Queries) GetSyncTaskState(ctx context.Context, id uuid.UUID) (GetSyncTaskStateRow, error) {
	row := q.db.QueryRowContext(ctx, getSyncTaskState, id)
	var i GetSyncTaskStateRow
	err := row.Scan(&i.ID, &i.State)
	return i, err
}

const isAnySyncTaskProcessingForUser = `-- name: IsAnySyncTaskProcessingForUser :one
SELECT id, created_at, updated_at, user_name, state, details FROM "artist_once_sync_tasks"
WHERE user_name = $1 AND state NOT IN ('done', 'error')
LIMIT 1
`

func (q *Queries) IsAnySyncTaskProcessingForUser(ctx context.Context, userName string) (ArtistOnceSyncTask, error) {
	row := q.db.QueryRowContext(ctx, isAnySyncTaskProcessingForUser, userName)
	var i ArtistOnceSyncTask
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

const updateSyncTaskState = `-- name: UpdateSyncTaskState :exec
UPDATE artist_once_sync_tasks
SET state = $1, details = $2
WHERE id = $3
`

type UpdateSyncTaskStateParams struct {
	State   TaskState      `json:"state"`
	Details sql.NullString `json:"details"`
	ID      uuid.UUID      `json:"id"`
}

func (q *Queries) UpdateSyncTaskState(ctx context.Context, arg UpdateSyncTaskStateParams) error {
	_, err := q.db.ExecContext(ctx, updateSyncTaskState, arg.State, arg.Details, arg.ID)
	return err
}
