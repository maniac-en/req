package endpoints

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/maniac-en/req/internal/backend/crud"
	"github.com/maniac-en/req/internal/backend/database"
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
		log.Warn("endpoint read failed validation", "id", id)
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
		log.Warn("endpoint delete failed validation", "id", id)
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

func (e *EndpointsManager) ListByCollectionByPage(ctx context.Context, collectionID int64, limit, offset int) (*PaginatedEndpoints, error) {
	if err := crud.ValidateID(collectionID); err != nil {
		log.Warn("endpoint list failed collection validation", "collection_id", collectionID)
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

	pagination := crud.CalculatePagination(total, limit, offset)

	result := &PaginatedEndpoints{
		Endpoints:          entities,
		PaginationMetadata: pagination,
	}
	log.Info("retrieved endpoints", "collection_id", collectionID, "count", len(entities), "total", pagination.Total, "page", pagination.CurrentPage, "total_pages", pagination.TotalPages)
	return result, nil
}

func (e *EndpointsManager) ListByCollection(ctx context.Context, collectionID int64) ([]EndpointEntity, error) {
	if err := crud.ValidateID(collectionID); err != nil {
		log.Warn("endpoint list failed collection validation", "collection_id", collectionID)
		return nil, crud.ErrInvalidInput
	}

	log.Debug("listing endpoints", "collection_id", collectionID, "limit")

	endpoints, err := e.DB.ListEndpointsByCollection(context.Background(), collectionID)
	if err != nil && err != sql.ErrNoRows {
		log.Warn("error occured while fetching endpoints", "collection_id", collectionID)
		return nil, err
	}

	entities := make([]EndpointEntity, len(endpoints))
	for i, endpoint := range endpoints {
		entities[i] = EndpointEntity{Endpoint: endpoint}
	}

	log.Info("retrieved endpoints", "collection_id", collectionID, "count")
	return entities, nil
}

func (e *EndpointsManager) CreateEndpoint(ctx context.Context, data EndpointData) (EndpointEntity, error) {
	if err := crud.ValidateID(data.CollectionID); err != nil {
		log.Warn("endpoint creation failed collection validation", "collection_id", data.CollectionID)
		return EndpointEntity{}, crud.ErrInvalidInput
	}
	if err := crud.ValidateName(data.Name); err != nil {
		log.Warn("endpoint creation failed name validation", "name", data.Name)
		return EndpointEntity{}, crud.ErrInvalidInput
	}
	if data.Method == "" {
		log.Warn("endpoint creation failed - method required", "method", data.Method)
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

func (e *EndpointsManager) UpdateEndpointName(ctx context.Context, id int64, name string) (EndpointEntity, error) {
	if err := crud.ValidateID(id); err != nil {
		log.Warn("endpoint update failed ID validation", "id", id)
		return EndpointEntity{}, crud.ErrInvalidInput
	}

	if err := crud.ValidateName(name); err != nil {
		log.Warn("endpoint update failed name validation", "name", name)
		return EndpointEntity{}, crud.ErrInvalidInput
	}

	log.Debug("updating endpoint name", "id", id, "name", name)

	endpoint, err := e.DB.UpdateEndpointName(ctx, database.UpdateEndpointNameParams{
		Name: name,
		ID:   id,
	})

	if err != nil {
		if err == sql.ErrNoRows {
			log.Debug("endpoint not found for update", "id", id)
			return EndpointEntity{}, crud.ErrNotFound
		}
		log.Error("failed to update endpoint", "id", id, "name", name, "error", err)
		return EndpointEntity{}, err
	}

	log.Info("updated endpoint", "id", endpoint.ID, "name", endpoint.Name)
	return EndpointEntity{Endpoint: endpoint}, nil
}

func (e *EndpointsManager) UpdateEndpoint(ctx context.Context, id int64, data EndpointData) (EndpointEntity, error) {
	if err := crud.ValidateID(id); err != nil {
		log.Warn("endpoint update failed ID validation", "id", id)
		return EndpointEntity{}, crud.ErrInvalidInput
	}
	if err := crud.ValidateName(data.Name); err != nil {
		log.Warn("endpoint update failed name validation", "name", data.Name)
		return EndpointEntity{}, crud.ErrInvalidInput
	}
	if data.Method == "" {
		log.Warn("endpoint update failed - method required", "method", data.Method)
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
