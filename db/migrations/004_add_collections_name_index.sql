-- +goose Up
CREATE INDEX idx_collections_name ON collections(name);

-- +goose Down
DROP INDEX IF EXISTS idx_collections_name;