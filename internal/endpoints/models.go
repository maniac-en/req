package endpoints

import (
	"time"

	"github.com/maniac-en/req/internal/crud"
	"github.com/maniac-en/req/internal/database"
	"github.com/maniac-en/req/internal/log"
)

type EndpointEntity struct {
	database.Endpoint
}

func (c EndpointEntity) GetID() int64 {
	return c.ID
}

func (c EndpointEntity) GetName() string {
	return c.Name
}

func (c EndpointEntity) GetCreatedAt() time.Time {
	parsed, err := time.Parse(time.RFC3339, c.CreatedAt)
	if err != nil {
		log.Debug("failed to parse created_at timestamp", "timestamp", c.CreatedAt, "error", err)
		return time.Time{}
	}
	return parsed
}

func (c EndpointEntity) GetUpdatedAt() time.Time {
	parsed, err := time.Parse(time.RFC3339, c.UpdatedAt)
	if err != nil {
		log.Debug("failed to parse updated_at timestamp", "timestamp", c.UpdatedAt, "error", err)
		return time.Time{}
	}
	return parsed
}

type EndpointsManager struct {
	DB *database.Queries
}

type PaginatedEndpoints struct {
	Endpoints []EndpointEntity `json:"endpoints"`
	crud.PaginationMetadata
}

type EndpointData struct {
	CollectionID int64
	Name         string
	Method       string
	URL          string
	Headers      string
	QueryParams  map[string]string
	RequestBody  string
}
