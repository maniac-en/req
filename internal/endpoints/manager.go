package endpoints

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/maniac-en/req/internal/crud"
	"github.com/maniac-en/req/internal/database"
	"github.com/maniac-en/req/internal/log"
)

func NewEndpointsManager(db *database.Queries) *EndpointsManager {
	return &EndpointsManager{DB: db}
}

func (e *EndpointsManager) Create(ctx context.Context, name string) (EndpointEntity, error) {
	return EndpointEntity{}, fmt.Errorf("use CreateEndpoint to create an endpoint with full data")
}

func (e *EndpointsManager) Read(ctx context.Context, id int64) (EndpointEntity, error) {
	if err := crud.ValidateID(id); err != nil {
		log.Debug("endpoint read failed validation", "id", id)
		return EndpointEntity{}, crud.ErrInvalidInput
	}

	log.Debug("reading endpoint", "id", id)
	endpoint, err := e.DB.GetEndpoint(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Debug("endpoint not found", "id", id)
			return EndpointEntity{}, crud.ErrNotFound
		}
		log.Error("failed to read endpoint", "id", id, "error", err)
		return EndpointEntity{}, err
	}

	return EndpointEntity{Endpoint: endpoint}, nil
}

func (e *EndpointsManager) Update(ctx context.Context, id int64, name string) (EndpointEntity, error) {
	return EndpointEntity{}, fmt.Errorf("use UpdateEndpoint to update an endpoint with full data")
}

func (e *EndpointsManager) Delete(ctx context.Context, id int64) error {
	if err := crud.ValidateID(id); err != nil {
		log.Debug("endpoint delete failed validation", "id", id)
		return crud.ErrInvalidInput
	}

	log.Debug("deleting endpoint", "id", id)
	err := e.DB.DeleteEndpoint(ctx, id)
	if err != nil {
		log.Error("failed to delete endpoint", "id", id, "error", err)
		return err
	}

	log.Info("deleted endpoint", "id", id)
	return nil
}

func (e *EndpointsManager) List(ctx context.Context) ([]EndpointEntity, error) {
	return nil, fmt.Errorf("use ListByCollection to list endpoints for a specific collection")
}

func (e *EndpointsManager) ListByCollection(ctx context.Context, collectionID int64, limit, offset int) (*PaginatedEndpoints, error) {
	if err := crud.ValidateID(collectionID); err != nil {
		log.Debug("endpoint list failed collection validation", "collection_id", collectionID)
		return nil, crud.ErrInvalidInput
	}

	log.Debug("listing paginated endpoints", "collection_id", collectionID, "limit", limit, "offset", offset)

	total, err := e.DB.CountEndpointsByCollection(ctx, collectionID)
	if err != nil {
		log.Error("failed to count endpoints", "collection_id", collectionID, "error", err)
		return nil, err
	}

	endpoints, err := e.DB.ListEndpointsPaginated(ctx, database.ListEndpointsPaginatedParams{
		CollectionID: collectionID,
		Limit:        int64(limit),
		Offset:       int64(offset),
	})
	if err != nil {
		log.Error("failed to get paginated endpoints", "collection_id", collectionID, "limit", limit, "offset", offset, "error", err)
		return nil, err
	}

	entities := make([]EndpointEntity, len(endpoints))
	for i, endpoint := range endpoints {
		entities[i] = EndpointEntity{Endpoint: endpoint}
	}

	// Calculate pagination metadata
	totalPages := int((total + int64(limit) - 1) / int64(limit)) // Ceiling division
	currentPage := (offset / limit) + 1
	hasNext := offset+len(endpoints) < int(total)
	hasPrev := offset > 0

	result := &PaginatedEndpoints{
		Endpoints:   entities,
		Total:       total,
		Offset:      int(offset),
		Limit:       int(limit),
		HasNext:     hasNext,
		HasPrev:     hasPrev,
		TotalPages:  totalPages,
		CurrentPage: currentPage,
	}
	log.Info("retrieved endpoints", "collection_id", collectionID, "count", len(entities), "total", total, "page", currentPage, "total_pages", totalPages)
	return result, nil
}

func (e *EndpointsManager) CreateEndpoint(ctx context.Context, data EndpointData) (EndpointEntity, error) {
	if err := crud.ValidateID(data.CollectionID); err != nil {
		log.Debug("endpoint creation failed collection validation", "collection_id", data.CollectionID)
		return EndpointEntity{}, crud.ErrInvalidInput
	}
	if err := crud.ValidateName(data.Name); err != nil {
		log.Debug("endpoint creation failed name validation", "name", data.Name)
		return EndpointEntity{}, crud.ErrInvalidInput
	}
	if data.Method == "" || data.URL == "" {
		log.Debug("endpoint creation failed - method and URL required", "method", data.Method, "url", data.URL)
		return EndpointEntity{}, crud.ErrInvalidInput
	}

	headersJSON := data.Headers
	if headersJSON == "" {
		headersJSON = "{}"
	}

	queryParamsJSON := "{}"
	if len(data.QueryParams) > 0 {
		qpBytes, err := json.Marshal(data.QueryParams)
		if err != nil {
			log.Error("failed to marshal query params", "error", err)
			return EndpointEntity{}, err
		}
		queryParamsJSON = string(qpBytes)
	}

	log.Debug("creating endpoint", "collection_id", data.CollectionID, "name", data.Name, "method", data.Method, "url", data.URL)
	endpoint, err := e.DB.CreateEndpoint(ctx, database.CreateEndpointParams{
		CollectionID: data.CollectionID,
		Name:         data.Name,
		Method:       data.Method,
		Url:          data.URL,
		Headers:      headersJSON,
		QueryParams:  queryParamsJSON,
		RequestBody:  data.RequestBody,
	})
	if err != nil {
		log.Error("failed to create endpoint", "collection_id", data.CollectionID, "name", data.Name, "error", err)
		return EndpointEntity{}, err
	}

	log.Info("created endpoint", "id", endpoint.ID, "name", endpoint.Name, "collection_id", endpoint.CollectionID)
	return EndpointEntity{Endpoint: endpoint}, nil
}

func (e *EndpointsManager) UpdateEndpoint(ctx context.Context, id int64, data EndpointData) (EndpointEntity, error) {
	if err := crud.ValidateID(id); err != nil {
		log.Debug("endpoint update failed ID validation", "id", id)
		return EndpointEntity{}, crud.ErrInvalidInput
	}
	if err := crud.ValidateName(data.Name); err != nil {
		log.Debug("endpoint update failed name validation", "name", data.Name)
		return EndpointEntity{}, crud.ErrInvalidInput
	}
	if data.Method == "" || data.URL == "" {
		log.Debug("endpoint update failed - method and URL required", "method", data.Method, "url", data.URL)
		return EndpointEntity{}, crud.ErrInvalidInput
	}

	headersJSON := data.Headers
	if headersJSON == "" {
		headersJSON = "{}"
	}

	queryParamsJSON := "{}"
	if len(data.QueryParams) > 0 {
		qpBytes, err := json.Marshal(data.QueryParams)
		if err != nil {
			log.Error("failed to marshal query params", "error", err)
			return EndpointEntity{}, err
		}
		queryParamsJSON = string(qpBytes)
	}

	log.Debug("updating endpoint", "id", id, "name", data.Name, "method", data.Method, "url", data.URL)
	endpoint, err := e.DB.UpdateEndpoint(ctx, database.UpdateEndpointParams{
		Name:        data.Name,
		Method:      data.Method,
		Url:         data.URL,
		Headers:     headersJSON,
		QueryParams: queryParamsJSON,
		RequestBody: data.RequestBody,
		ID:          id,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			log.Debug("endpoint not found for update", "id", id)
			return EndpointEntity{}, crud.ErrNotFound
		}
		log.Error("failed to update endpoint", "id", id, "name", data.Name, "error", err)
		return EndpointEntity{}, err
	}

	log.Info("updated endpoint", "id", endpoint.ID, "name", endpoint.Name)
	return EndpointEntity{Endpoint: endpoint}, nil
}
