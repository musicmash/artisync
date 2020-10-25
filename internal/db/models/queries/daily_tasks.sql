-- name: CreateDailySyncTask :exec
INSERT INTO artist_daily_sync_tasks (user_name)
VALUES (@user_name)
ON CONFLICT (user_name)
DO UPDATE SET updated_at = now();

-- name: GetUserDailySyncTask :one
SELECT * FROM artist_daily_sync_tasks
WHERE user_name = @user_name
LIMIT 1;
