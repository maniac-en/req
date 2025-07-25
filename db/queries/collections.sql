-- name: CreateCollection :one
INSERT INTO collections (name) VALUES (?) RETURNING *;
