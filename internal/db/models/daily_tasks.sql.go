// Code generated by sqlc. DO NOT EDIT.
// source: daily_tasks.sql

package models

import (
	"context"
	"time"
)

const createDailySyncTask = `-- name: CreateDailySyncTask :exec
INSERT INTO artist_daily_sync_tasks (user_name)
VALUES ($1)
ON CONFLICT (user_name)
DO UPDATE SET updated_at = now()
`

func (q *Queries) CreateDailySyncTask(ctx context.Context, userName string) error {
	_, err := q.db.ExecContext(ctx, createDailySyncTask, userName)
	return err
}

const disableDailySyncTask = `-- name: DisableDailySyncTask :exec
DELETE FROM artist_daily_sync_tasks
WHERE user_name = $1
`

func (q *Queries) DisableDailySyncTask(ctx context.Context, userName string) error {
	_, err := q.db.ExecContext(ctx, disableDailySyncTask, userName)
	return err
}

const getDailyLock = `-- name: GetDailyLock :exec
SELECT id, created_at, updated_at, user_name from artist_daily_sync_tasks
FOR UPDATE
`

func (q *Queries) GetDailyLock(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, getDailyLock)
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

const resetDailyTasks = `-- name: ResetDailyTasks :execrows
UPDATE artist_daily_sync_tasks
SET updated_at=now()
WHERE
    updated_at < $1 AND
    user_name in (
        SELECT user_name FROM artist_sync_refresh_tokens
        WHERE expired_at > now()
    )
`

func (q *Queries) ResetDailyTasks(ctx context.Context, today time.Time) (int64, error) {
	result, err := q.db.ExecContext(ctx, resetDailyTasks, today)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const scheduleDailyTasks = `-- name: ScheduleDailyTasks :execrows
INSERT INTO artist_one_time_sync_tasks (user_name, state)
SELECT
    daily.user_name, 'created'
FROM
    artist_daily_sync_tasks AS daily
    LEFT JOIN artist_one_time_sync_tasks AS one ON daily.user_name = one.user_name
        AND one.created_at >= $1
    LEFT JOIN artist_sync_refresh_tokens AS token ON daily.user_name = token.user_name
        AND token.expired_at >= $1
WHERE one.created_at IS NULL AND token.value != ''
`

func (q *Queries) ScheduleDailyTasks(ctx context.Context, today time.Time) (int64, error) {
	result, err := q.db.ExecContext(ctx, scheduleDailyTasks, today)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
