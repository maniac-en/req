package collections

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/maniac-en/req/internal/database"
)

func (c *CollectionsManager) CreateEndpoint(ctx context.Context, name, method, url string, collectionId int64) (database.Endpoint, error) {
	return c.DB.CreateEndpoint(ctx, database.CreateEndpointParams{
		Name:         name,
		CollectionID: collectionId,
		Method:       method,
		Url:          url,
	})
}

func (c *CollectionsManager) GetEndpoints(ctx context.Context, collectionId int64) ([]database.Endpoint, error) {
	return c.DB.ListEndpoints(ctx, collectionId)
}

func (c *CollectionsManager) GetEndpoint(ctx context.Context, id int64) (database.Endpoint, error) {
	return c.DB.GetEndpoint(ctx, id)
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
	existingEndpoint, err := c.DB.GetEndpoint(ctx, id)
	if err != nil {
		return database.Endpoint{}, err
	}

	if name != "" {
		existingEndpoint.Name = name
	}
	if method != "" {
		existingEndpoint.Method = method
	}
	if url != "" {
		existingEndpoint.Url = url
	}
	if headers != "" {
		if !json.Valid([]byte(headers)) {
			return database.Endpoint{}, fmt.Errorf("invalid JSON for header parameters: %s", headers)
		}
		existingEndpoint.Headers = headers
	}
	if queryParams != "" {
		if !json.Valid([]byte(queryParams)) {
			return database.Endpoint{}, fmt.Errorf("invalid JSON for query parameters: %s", queryParams)
		}
		existingEndpoint.QueryParams = queryParams
	}
	if requestBody != "" {
		existingEndpoint.RequestBody = requestBody
	}

	return c.DB.UpdateEndpoint(ctx, database.UpdateEndpointParams{
		ID:          existingEndpoint.ID,
		Name:        existingEndpoint.Name,
		Method:      existingEndpoint.Method,
		Url:         existingEndpoint.Url,
		Headers:     existingEndpoint.Headers,
		QueryParams: existingEndpoint.QueryParams,
		RequestBody: existingEndpoint.RequestBody,
	})
}
