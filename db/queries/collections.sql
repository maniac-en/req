-- name: CreateCollection :one
INSERT INTO collections (name) VALUES (?) RETURNING *;

-- name: GetCollectionsPaginated :many
SELECT * FROM collections
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: GetCollections :many
SELECT
    c.id,
    c.name,
    c.created_at,
    c.updated_at,
    COUNT(e.id) AS endpoint_count
FROM collections c
LEFT JOIN endpoints e ON e.collection_id = c.id
GROUP BY c.id, c.name, c.created_at, c.updated_at
ORDER BY c.created_at DESC;

-- name: CountCollections :one
SELECT COUNT(*) FROM collections;

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
