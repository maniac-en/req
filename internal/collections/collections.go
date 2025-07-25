package collections

import (
	"context"
	"fmt"
	"strings"

	"github.com/maniac-en/req/internal/database"
	"github.com/maniac-en/req/internal/log"
)

func NewCollectionsManager(q *database.Queries) *CollectionsManager {
	collectionsManager := CollectionsManager{
		DB: q,
	}
	return &collectionsManager
}

func (c *CollectionsManager) GetAllCollections(ctx context.Context) ([]database.Collection, error) {
	log.Debug("fetching all collections")
	dbCollections, err := c.DB.GetAllCollections(ctx)
	if err != nil {
		log.Error("failed to fetch collections", "error", err)
		return nil, err
	}
	log.Info("fetched collections", "count", len(dbCollections))
	return dbCollections, nil
}

func (c *CollectionsManager) CreateCollection(ctx context.Context, name string) (database.Collection, error) {
	if err := validateCollectionName(name); err != nil {
		log.Error("invalid collection name", "name", name, "error", err)
		return database.Collection{}, err
	}

	log.Info("creating collection", "name", name)
	collection, err := c.DB.CreateCollection(ctx, name)
	if err != nil {
		log.Error("failed to create collection", "name", name, "error", err)
		return database.Collection{}, err
	}
	log.Info("created collection", "id", collection.ID, "name", collection.Name)
	return collection, nil
}

func (c *CollectionsManager) UpdateCollectionName(ctx context.Context, name string, collectionId int) (database.Collection, error) {
	if err := validateCollectionName(name); err != nil {
		log.Error("invalid collection name", "name", name, "error", err)
		return database.Collection{}, err
	}

	id := int64(collectionId)
	log.Info("updating collection name", "id", id, "new_name", name)
	collection, err := c.DB.UpdateCollectionName(ctx, database.UpdateCollectionNameParams{
		Name: name,
		ID:   id,
	})
	if err != nil {
		log.Error("failed to update collection", "id", id, "name", name, "error", err)
		return database.Collection{}, err
	}
	log.Info("updated collection", "id", collection.ID, "name", collection.Name)
	return collection, nil
}

func (c *CollectionsManager) DeleteCollection(ctx context.Context, id int) error {
	log.Info("deleting collection", "id", id)
	err := c.DB.DeleteCollection(ctx, int64(id))
	if err != nil {
		log.Error("failed to delete collection", "id", id, "error", err)
		return err
	}
	log.Info("deleted collection", "id", id)
	return nil
}

func validateCollectionName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("collection name cannot be empty")
	}
	if len(name) > 100 {
		return fmt.Errorf("collection name cannot exceed 100 characters")
	}
	if strings.ContainsAny(name, "/\\:*?\"<>|") {
		return fmt.Errorf("collection name contains invalid characters")
	}
	return nil
}
