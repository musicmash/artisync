-- name: IsAnySyncTaskProcessingForUser :one
SELECT * FROM artist_one_time_sync_tasks
WHERE user_name = @user_name AND state NOT IN ('done', 'error')
LIMIT 1;

-- name: CreateOneTimeSyncTask :one
INSERT INTO artist_one_time_sync_tasks (user_name)
VALUES (@user_name)
RETURNING *;

-- name: GetOneTimeSyncTaskState :one
SELECT id, state, details FROM artist_one_time_sync_tasks
WHERE id = @id
LIMIT 1;

-- name: UpdateOneTimeSyncTaskState :exec
UPDATE artist_one_time_sync_tasks
SET state = @state, details = @details
WHERE id = @id;

-- name: GetNextOneTimeSyncTasks :many
SELECT
    artist_one_time_sync_tasks.id,
    artist_one_time_sync_tasks.user_name,
    artist_sync_refresh_tokens.value as refresh_token
FROM artist_one_time_sync_tasks
LEFT JOIN artist_sync_refresh_tokens ON (
    artist_sync_refresh_tokens.user_name=artist_one_time_sync_tasks.user_name
)
WHERE state='created' AND artist_sync_refresh_tokens.expired_at > now()
LIMIT $1
FOR UPDATE;
