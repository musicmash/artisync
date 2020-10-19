-- name: CreateRefreshToken :exec
INSERT INTO "artist_sync_refresh_tokens" (user_name, expired_at, value)
VALUES (@user_name, @expired_at, @value);

-- name: GetUserSyncTask :one
SELECT * FROM "artist_once_sync_tasks"
WHERE user_name = @user_name;
