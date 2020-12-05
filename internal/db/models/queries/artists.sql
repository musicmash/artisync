-- name: CreateArtist :one
INSERT INTO artists (name, poster)
VALUES (@name, @poster)
ON CONFLICT DO NOTHING
RETURNING *;
