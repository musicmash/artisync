-- name: CreateDailySyncTask :exec
INSERT INTO artist_daily_sync_tasks (user_name)
VALUES (@user_name)
ON CONFLICT (user_name)
DO UPDATE SET updated_at = now();

-- name: GetUserDailySyncTask :one
SELECT * FROM artist_daily_sync_tasks
WHERE user_name = @user_name
LIMIT 1;

-- name: ScheduleDailyTasks :execrows
INSERT INTO artist_one_time_sync_tasks (user_name, state)
SELECT
    daily.user_name, 'created'
FROM
    artist_daily_sync_tasks AS daily
LEFT JOIN artist_one_time_sync_tasks AS one
    ON daily.user_name = one.user_name
AND one.created_at >= @yesterday
LEFT JOIN artist_sync_refresh_tokens AS token
    ON daily.user_name = token.user_name AND token.expired_at >= now()
WHERE
    daily.updated_at < @today
    AND one.created_at IS NULL
    AND token.value != '';

-- name: ResetDailyTasks :execrows
UPDATE artist_daily_sync_tasks
SET updated_at=now()
WHERE
    updated_at < @today AND
    user_name in (
        SELECT user_name FROM artist_sync_refresh_tokens
        WHERE expired_at > now()
    );
