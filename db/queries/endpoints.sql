-- name: CreateEndpoint :one
INSERT INTO endpoints (
    collection_id,
    name,
    method,
    url,
    headers,
    query_params,
    request_body
) VALUES (
    ?, ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: GetEndpoint :one
SELECT * FROM endpoints
WHERE id = ? LIMIT 1;

-- name: ListEndpoints :many
SELECT * FROM endpoints
WHERE collection_id = ?
ORDER BY name;

-- name: UpdateEndpoint :one
UPDATE endpoints
SET
    name = ?,
    method = ?,
    url = ?,
    headers = ?,
    query_params = ?,
    request_body = ?
WHERE
    id = ?
RETURNING *;

-- name: DeleteEndpoint :exec
DELETE FROM endpoints
WHERE id = ?;
