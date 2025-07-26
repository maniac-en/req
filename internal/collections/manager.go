package collections

import (
	"context"
	"database/sql"

	"github.com/maniac-en/req/internal/crud"
	"github.com/maniac-en/req/internal/database"
	"github.com/maniac-en/req/internal/log"
)

func NewCollectionsManager(db *database.Queries) *CollectionsManager {
	return &CollectionsManager{DB: db}
}

func (c *CollectionsManager) Create(ctx context.Context, name string) (CollectionEntity, error) {
	if err := crud.ValidateName(name); err != nil {
		log.Debug("collection creation failed validation", "name", name)
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
		log.Debug("collection read failed validation", "id", id)
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
		log.Debug("collection update failed ID validation", "id", id)
		return CollectionEntity{}, crud.ErrInvalidInput
	}
	if err := crud.ValidateName(name); err != nil {
		log.Debug("collection update failed name validation", "name", name)
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
		log.Debug("collection delete failed validation", "id", id)
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
	log.Debug("listing all collections")
	collections, err := c.DB.GetAllCollections(ctx)
	if err != nil {
		log.Error("failed to list collections", "error", err)
		return nil, err
	}

	entities := make([]CollectionEntity, len(collections))
	for i, collection := range collections {
		entities[i] = CollectionEntity{Collection: collection}
	}

	log.Debug("listed collections", "count", len(entities))
	return entities, nil
}
