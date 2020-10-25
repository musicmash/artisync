-- name: CreateArtist :one
INSERT INTO spotify_artists (name, poster)
VALUES (@name, @poster)
ON CONFLICT DO NOTHING
RETURNING *;
