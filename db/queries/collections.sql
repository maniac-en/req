-- name: CreateCollection :one
INSERT INTO collections (name) VALUES (?) RETURNING *;

-- name: GetAllCollections :many
SELECT * FROM collections;

-- name: UpdateCollectionName :one
UPDATE collections
SET name = ?
WHERE id = ?
RETURNING *;

-- name: DeleteCollection :exec
DELETE FROM collections
WHERE id = ?;

-- name: GetCollection :one
SELECT * FROM collections
WHERE id = ?;
