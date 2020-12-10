-- name: CreateSubscription :exec
INSERT INTO subscriptions (created_at, user_name, artist_id)
VALUES (now(), @user_name, @artist_id)
ON CONFLICT  DO NOTHING;
