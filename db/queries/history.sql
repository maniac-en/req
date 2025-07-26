-- name: CreateHistoryEntry :one
INSERT INTO history (
    collection_id, collection_name, endpoint_name,
    method, url, status_code, duration, response_size,
    request_headers, query_params, request_body,
    response_body, response_headers, executed_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetHistoryById :one
SELECT * FROM history
WHERE id = ?;

-- name: GetHistoryByCollection :many
SELECT id, endpoint_name, status_code, executed_at, url, method FROM history
WHERE collection_id = ?
ORDER BY executed_at DESC
LIMIT ? OFFSET ?;

-- name: CountHistoryByCollection :one
SELECT COUNT(*) FROM history
WHERE collection_id = ?;

-- name: DeleteHistoryEntry :exec
DELETE FROM history
WHERE id = ?;

-- name: DeleteOldHistory :exec
DELETE FROM history
WHERE executed_at < datetime('now', '-30 days');