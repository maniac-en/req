package collections

import (
	"time"

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
	Total       int64              `json:"total"`
	HasNext     bool               `json:"has_next"`
	HasPrev     bool               `json:"has_prev"`
	Limit       int                `json:"limit"`
	Offset      int                `json:"offset"`
	TotalPages  int                `json:"total_pages"`
	CurrentPage int                `json:"current_page"`
}
