-- +goose Up
CREATE TABLE history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    collection_id INTEGER,
    collection_name TEXT,
    endpoint_name TEXT,
    method TEXT NOT NULL,
    url TEXT NOT NULL,
    status_code INTEGER NOT NULL,
    duration INTEGER NOT NULL, -- duration in milliseconds
    response_size INTEGER DEFAULT 0,
    request_headers TEXT DEFAULT '{}',
    query_params TEXT DEFAULT '{}',
    request_body TEXT DEFAULT '',
    response_body TEXT DEFAULT '',
    response_headers TEXT DEFAULT '{}',
    executed_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_history_collection_id ON history(collection_id);
CREATE INDEX idx_history_executed_at ON history(executed_at);
CREATE INDEX idx_history_status_code ON history(status_code);

-- +goose Down
DROP INDEX IF EXISTS idx_history_status_code;
DROP INDEX IF EXISTS idx_history_executed_at;
DROP INDEX IF EXISTS idx_history_collection_id;
DROP TABLE IF EXISTS history;