-- +goose Up
CREATE TABLE endpoints (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    collection_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    method TEXT NOT NULL,
    url TEXT NOT NULL,
    headers TEXT DEFAULT '{}' NOT NULL,
    query_params TEXT DEFAULT '{}' NOT NULL,
    request_body TEXT DEFAULT '' NOT NULL,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (collection_id) REFERENCES collections(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE endpoints;
<<<<<<< Updated upstream

=======
>>>>>>> Stashed changes
