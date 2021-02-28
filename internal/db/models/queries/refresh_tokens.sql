-- name: CreateRefreshToken :exec
INSERT INTO artist_sync_refresh_tokens (user_name, expired_at, value)
VALUES (@user_name, @expired_at, @value)
ON CONFLICT (user_name)
DO UPDATE SET value = @value, expired_at = @expired_at;
