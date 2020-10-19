-- name: CreateSubscription :exec
INSERT INTO subscriptions (user_name, artist_id)
SELECT @user_name, artist_id FROM spotify_artist_associations
WHERE store_name = @store_name AND store_id = @store_id
ON CONFLICT DO NOTHING;
