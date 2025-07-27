package collections

import (
	"time"

	"github.com/maniac-en/req/internal/crud"
	"github.com/maniac-en/req/internal/database"
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
	return crud.ParseTimestamp(c.CreatedAt)
}

func (c CollectionEntity) GetUpdatedAt() time.Time {
	return crud.ParseTimestamp(c.UpdatedAt)
}

type CollectionsManager struct {
	DB *database.Queries
}

type PaginatedCollections struct {
	Collections []CollectionEntity `json:"collections"`
	crud.PaginationMetadata
}
