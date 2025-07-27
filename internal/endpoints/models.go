package endpoints

import (
	"time"

	"github.com/maniac-en/req/internal/crud"
	"github.com/maniac-en/req/internal/database"
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
	return crud.ParseTimestamp(c.CreatedAt)
}

func (c EndpointEntity) GetUpdatedAt() time.Time {
	return crud.ParseTimestamp(c.UpdatedAt)
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
