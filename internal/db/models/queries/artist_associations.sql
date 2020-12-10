-- name: GetArtistAssociation :one
SELECT * FROM artist_associations
WHERE store_name = @store_name and store_id = @store_id
LIMIT 1;

-- name: CreateArtistAssociation :one
INSERT INTO artist_associations (artist_id, store_name, store_id)
VALUES (@artist_id, @store_name, @store_id)
ON CONFLICT DO NOTHING
RETURNING *;
