package collections

import (
	"context"

	"github.com/maniac-en/req/internal/database"
	"github.com/maniac-en/req/internal/log"
)

func NewCollectionsManager(q *database.Queries) *CollectionsManager {
	collectionsManager := CollectionsManager{
		DB: q,
	}
	return &collectionsManager
}

func (c *CollectionsManager) GetAllCollections() ([]database.Collection, error) {
	log.Debug("fetching all collections")
	dbCollections, err := c.DB.GetAllCollections(context.Background())
	if err != nil {
		log.Error("failed to fetch collections", "error", err)
		return nil, err
	}
	log.Info("fetched collections", "count", len(dbCollections))
	return dbCollections, nil
}

func (c *CollectionsManager) CreateCollection(name string) (database.Collection, error) {
	log.Info("creating collection", "name", name)
	collection, err := c.DB.CreateCollection(context.Background(), name)
	if err != nil {
		log.Error("failed to create collection", "name", name, "error", err)
		return database.Collection{}, err
	}
	log.Info("created collection", "id", collection.ID, "name", collection.Name)
	return collection, nil
}

func (c *CollectionsManager) UpdateCollectionName(name string, collectionId int) (database.Collection, error) {
	id := int64(collectionId)
	log.Info("updating collection name", "id", id, "new_name", name)
	collection, err := c.DB.UpdateCollectionName(context.Background(), database.UpdateCollectionNameParams{
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

func (c *CollectionsManager) DeleteCollection(id int) error {
	log.Info("deleting collection", "id", id)
	err := c.DB.DeleteCollection(context.Background(), int64(id))
	if err != nil {
		log.Error("failed to delete collection", "id", id, "error", err)
		return err
	}
	log.Info("deleted collection", "id", id)
	return nil
}
