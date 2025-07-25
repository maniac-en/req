-- +goose Up
CREATE TRIGGER update_collections_updated_at
AFTER UPDATE ON collections
FOR EACH ROW
BEGIN
    UPDATE collections 
    SET updated_at = CURRENT_TIMESTAMP 
    WHERE id = NEW.id;
END;

-- +goose Down
DROP TRIGGER IF EXISTS update_collections_updated_at;