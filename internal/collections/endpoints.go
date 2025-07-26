package collections

import (
	"context"
	"fmt"

	"github.com/maniac-en/req/internal/database"
	"github.com/maniac-en/req/internal/log"
)

func (c *CollectionsManager) CreateEndpoint(
	ctx context.Context,
	name, method, url string,
	collectionId int64,
) (database.Endpoint, error) {
	_, err := c.DB.GetCollection(ctx, collectionId)
	log.Debug("Getting the collection with collectionId:", collectionId)
	if err != nil {
		log.Error("collection doesn't exists", "id", collectionId, "error", err)
		return database.Endpoint{}, fmt.Errorf("collection doesn't exists!")
	}

	if err := validateName(name); err != nil {
		log.Error("invalid endpoint name", "name", name, "error", err)
		return database.Endpoint{}, err
	}

	log.Info("creating endpoint", "name", name)
	endpoint, err := c.DB.CreateEndpoint(ctx, database.CreateEndpointParams{
		Name:         name,
		CollectionID: collectionId,
		Method:       method,
		Url:          url,
	})
	if err != nil {
		log.Error("failed to create endpoint", "name", name, "error", err)
		return database.Endpoint{}, err
	}
	log.Info("created collection", "id", endpoint.ID, "name", endpoint.Name)
	return endpoint, nil
}

func (c *CollectionsManager) GetEndpoints(ctx context.Context, collectionId int64) ([]database.Endpoint, error) {
	log.Debug("fetching all endpoints with collectionId", collectionId)
	dbEndpoints, err := c.DB.ListEndpoints(ctx, collectionId)
	if err != nil {
		log.Error("failed to fetch endpoints", "id", collectionId, "error", err)
		return nil, err
	}
	log.Info("endpoints fetched", "count", len(dbEndpoints))
	return dbEndpoints, nil
}

func (c *CollectionsManager) GetEndpoint(ctx context.Context, id int64) (database.Endpoint, error) {
	log.Debug("fetching endpoint", "id", id)
	endpoint, err := c.DB.GetEndpoint(ctx, id)
	if err != nil {
		log.Error("failed to fetch endpoint", "id", id, "err", err)
		return database.Endpoint{}, fmt.Errorf("failed to fetch endpoint")
	}
	log.Info("fetched endpoint", "id", id)
	return endpoint, nil
}

func (c *CollectionsManager) UpdateEndpoint(
	ctx context.Context,
	id int64,
	name string,
	method string,
	url string,
	headers string,
	queryParams string,
	requestBody string,
) (database.Endpoint, error) {
	log.Debug("fetching endpoint to update", "id", id)

	endpoint, err := c.DB.GetEndpoint(ctx, id)
	if err != nil {
		log.Error("failed to fetch endpoint", "id", id, "err", err)
		return database.Endpoint{}, err
	}

	patchEndpoint(&endpoint, name, method, url, headers, queryParams, requestBody)
	log.Info("updating endpoint", "id", id)
	updated, err := c.DB.UpdateEndpoint(ctx, database.UpdateEndpointParams{
		ID:          endpoint.ID,
		Name:        endpoint.Name,
		Method:      endpoint.Method,
		Url:         endpoint.Url,
		Headers:     endpoint.Headers,
		QueryParams: endpoint.QueryParams,
		RequestBody: endpoint.RequestBody,
	})
	if err != nil {
		log.Error("failed to update endpoint", "id", id, "err", err)
		return database.Endpoint{}, err
	}

	log.Info("updated endpoint", "id", updated.ID, "name", updated.Name, "method", updated.Method, "url", updated.Url)
	return updated, nil
}

func patchEndpoint(
	endpoint *database.Endpoint,
	name string,
	method string,
	url string,
	headers string,
	queryParams string,
	requestBody string,
) {
	if name != "" {
		endpoint.Name = name
	}
	if method != "" {
		endpoint.Method = method
	}
	if url != "" {
		endpoint.Url = url
	}
	if headers != "" {
		endpoint.Headers = headers
	}
	if queryParams != "" {
		endpoint.QueryParams = queryParams
	}
	if requestBody != "" {
		endpoint.RequestBody = requestBody
	}
}
