-- name: IsAnySyncTaskProcessingForUser :one
SELECT * FROM "artist_once_sync_tasks"
WHERE user_name = @user_name AND state NOT IN ('done', 'error')
LIMIT 1;

-- name: CreateSyncTask :one
INSERT INTO "artist_once_sync_tasks" (user_name)
VALUES (@user_name)
RETURNING *;

-- name: GetSyncTaskState :one
SELECT id, state FROM "artist_once_sync_tasks"
WHERE id = @id
LIMIT 1;

-- name: UpdateSyncTaskState :exec
UPDATE artist_once_sync_tasks
SET state = @state, details = @details
WHERE id = @id;

-- name: GetNextScheduledTask :one
SELECT
    artist_once_sync_tasks.id,
    artist_once_sync_tasks.user_name,
    artist_sync_refresh_tokens.value as refresh_token
FROM artist_once_sync_tasks
LEFT JOIN artist_sync_refresh_tokens ON (
    artist_sync_refresh_tokens.user_name=artist_once_sync_tasks.user_name
)
WHERE state = 'scheduled' AND artist_sync_refresh_tokens.expired_at > now()
LIMIT 1;
