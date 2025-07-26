package endpoints

import (
	"context"

	"github.com/maniac-en/req/internal/database"
)

func NewEndpointsManager(db *database.Queries) *EndpointManager {
	return &EndpointManager{DB: db}
}

func (e *EndpointManager) Create(ctx context.Context, name string) (EndpointEntity, error) {
	return EndpointEntity{}, nil
}

func (e *EndpointManager) Read(ctx context.Context, id int64) (EndpointEntity, error) {
	return EndpointEntity{}, nil
}

func (e *EndpointManager) Update(ctx context.Context, id int64, name string) (EndpointEntity, error) {
	return EndpointEntity{}, nil
}

func (e *EndpointManager) Delete(ctx context.Context, id int64) error {
	return nil
}

func (e *EndpointManager) List(ctx context.Context) ([]EndpointEntity, error) {
	return nil, nil
}

func (e *EndpointManager) ListPaginated(ctx context.Context, limit, offset int) (*PaginatedEndpoints, error) {
	return nil, nil
}
