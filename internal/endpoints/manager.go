package endpoints

import (
	"context"

	"github.com/maniac-en/req/internal/database"
)

func NewEndpointsManager(db *database.Queries) *EndpointsManager {
	return &EndpointsManager{DB: db}
}

func (e *EndpointsManager) Create(ctx context.Context, name string) (EndpointEntity, error) {
	return EndpointEntity{}, nil
}

func (e *EndpointsManager) Read(ctx context.Context, id int64) (EndpointEntity, error) {
	return EndpointEntity{}, nil
}

func (e *EndpointsManager) Update(ctx context.Context, id int64, name string) (EndpointEntity, error) {
	return EndpointEntity{}, nil
}

func (e *EndpointsManager) Delete(ctx context.Context, id int64) error {
	return nil
}

func (e *EndpointsManager) List(ctx context.Context) ([]EndpointEntity, error) {
	return nil, nil
}

func (e *EndpointsManager) ListPaginated(ctx context.Context, limit, offset int) (*PaginatedEndpoints, error) {
	return nil, nil
}
