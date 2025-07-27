package collections

import (
	"time"

	"github.com/maniac-en/req/internal/crud"
	"github.com/maniac-en/req/internal/database"
	"github.com/maniac-en/req/internal/log"
)

type CollectionEntity struct {
	database.Collection
}

func (c CollectionEntity) GetID() int64 {
	return c.ID
}

func (c CollectionEntity) GetName() string {
	return c.Name
}

func (c CollectionEntity) GetCreatedAt() time.Time {
	parsed, err := time.Parse(time.RFC3339, c.CreatedAt)
	if err != nil {
		log.Debug("failed to parse created_at timestamp", "timestamp", c.CreatedAt, "error", err)
		return time.Time{}
	}
	return parsed
}

func (c CollectionEntity) GetUpdatedAt() time.Time {
	parsed, err := time.Parse(time.RFC3339, c.UpdatedAt)
	if err != nil {
		log.Debug("failed to parse updated_at timestamp", "timestamp", c.UpdatedAt, "error", err)
		return time.Time{}
	}
	return parsed
}

type CollectionsManager struct {
	DB *database.Queries
}

type PaginatedCollections struct {
	Collections []CollectionEntity `json:"collections"`
	crud.PaginationMetadata
}
