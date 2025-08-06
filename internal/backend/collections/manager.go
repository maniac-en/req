package collections

import (
	"context"
	"database/sql"

	"github.com/maniac-en/req/internal/backend/crud"
	"github.com/maniac-en/req/internal/backend/database"
	"github.com/maniac-en/req/internal/log"
)

func NewCollectionsManager(db *database.Queries) *CollectionsManager {
	return &CollectionsManager{DB: db}
}

func (c *CollectionsManager) Create(ctx context.Context, name string) (CollectionEntity, error) {
	if err := crud.ValidateName(name); err != nil {
		log.Warn("collection creation failed validation", "name", name)
		return CollectionEntity{}, crud.ErrInvalidInput
	}

	log.Debug("creating collection", "name", name)
	collection, err := c.DB.CreateCollection(ctx, name)
	if err != nil {
		log.Error("failed to create collection", "name", name, "error", err)
		return CollectionEntity{}, err
	}

	log.Info("created collection", "id", collection.ID, "name", collection.Name)
	return CollectionEntity{Collection: collection}, nil
}

func (c *CollectionsManager) Read(ctx context.Context, id int64) (CollectionEntity, error) {
	if err := crud.ValidateID(id); err != nil {
		log.Warn("collection read failed validation", "id", id)
		return CollectionEntity{}, crud.ErrInvalidInput
	}

	log.Debug("reading collection", "id", id)
	collection, err := c.DB.GetCollection(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Debug("collection not found", "id", id)
			return CollectionEntity{}, crud.ErrNotFound
		}
		log.Error("failed to read collection", "id", id, "error", err)
		return CollectionEntity{}, err
	}

	return CollectionEntity{Collection: collection}, nil
}

func (c *CollectionsManager) Update(ctx context.Context, id int64, name string) (CollectionEntity, error) {
	if err := crud.ValidateID(id); err != nil {
		log.Warn("collection update failed ID validation", "id", id)
		return CollectionEntity{}, crud.ErrInvalidInput
	}
	if err := crud.ValidateName(name); err != nil {
		log.Warn("collection update failed name validation", "name", name)
		return CollectionEntity{}, crud.ErrInvalidInput
	}

	log.Debug("updating collection", "id", id, "name", name)
	collection, err := c.DB.UpdateCollectionName(ctx, database.UpdateCollectionNameParams{
		Name: name,
		ID:   id,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			log.Debug("collection not found for update", "id", id)
			return CollectionEntity{}, crud.ErrNotFound
		}
		log.Error("failed to update collection", "id", id, "name", name, "error", err)
		return CollectionEntity{}, err
	}

	log.Info("updated collection", "id", collection.ID, "name", collection.Name)
	return CollectionEntity{Collection: collection}, nil
}

func (c *CollectionsManager) Delete(ctx context.Context, id int64) error {
	if err := crud.ValidateID(id); err != nil {
		log.Warn("collection delete failed validation", "id", id)
		return crud.ErrInvalidInput
	}

	log.Debug("deleting collection", "id", id)
	err := c.DB.DeleteCollection(ctx, id)
	if err != nil {
		log.Error("failed to delete collection", "id", id, "error", err)
		return err
	}

	log.Info("deleted collection", "id", id)
	return nil
}

func (c *CollectionsManager) List(ctx context.Context) ([]CollectionEntity, error) {
	log.Debug("listing all collections without pagination")
	collections, err := c.DB.GetCollections(ctx)
	collectionsEntity := []CollectionEntity{}
	for _, collection := range collections {
		collectionsEntity = append(collectionsEntity, CollectionEntity{Collection: collection})
	}
	if err != nil {
		return nil, err
	}
	return collectionsEntity, nil
}

func (c *CollectionsManager) ListPaginated(ctx context.Context, limit, offset int) (*PaginatedCollections, error) {
	// Warn about unusual pagination parameters
	if limit > 1000 {
		log.Warn("large pagination limit requested", "limit", limit)
	}
	if offset < 0 {
		log.Warn("negative pagination offset", "offset", offset)
	}

	log.Debug("listing paginated collections", "limit", limit, "offset", offset)

	total, err := c.DB.CountCollections(ctx)
	if err != nil {
		log.Error("failed to count collections", "error", err)
		return nil, err
	}

	collections, err := c.DB.GetCollectionsPaginated(ctx, database.GetCollectionsPaginatedParams{
		Limit:  int64(limit),
		Offset: int64(offset),
	})
	if err != nil {
		log.Error("failed to get paginated collections", "limit", limit, "offset", offset, "error", err)
		return nil, err
	}

	entities := make([]CollectionEntity, len(collections))
	for i, collection := range collections {
		entities[i] = CollectionEntity{Collection: collection}
	}

	pagination := crud.CalculatePagination(total, limit, offset)

	result := &PaginatedCollections{
		Collections:        entities,
		PaginationMetadata: pagination,
	}
	log.Info("retrieved collections", "count", len(entities), "total", pagination.Total, "page", pagination.CurrentPage, "total_pages", pagination.TotalPages)
	return result, nil
}
