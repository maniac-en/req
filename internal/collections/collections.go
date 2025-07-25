package collections

import (
	"context"

	"github.com/maniac-en/req/internal/database"
)

func NewCollectionsManager(q *database.Queries) *CollectionsManager {
	collectionsManager := CollectionsManager{
		DB: q,
	}
	return &collectionsManager
}

func (c *CollectionsManager) GetAllCollections() ([]database.Collection, error) {
	dbCollections, err := c.DB.GetAllCollections(context.Background())
	return dbCollections, err
}

func (c *CollectionsManager) CreateCollection(name string) (database.Collection, error) {
	return c.DB.CreateCollection(context.Background(), name)
}

func (c *CollectionsManager) UpdateCollectionName(name string, collectionId int) (database.Collection, error) {
	id := int64(collectionId)
	return c.DB.UpdateCollectionName(context.Background(), database.UpdateCollectionNameParams{
		Name: name,
		ID:   id,
	})
}

func (c *CollectionsManager) DeleteCollection(id int) error {
	err := c.DB.DeleteCollection(context.Background(), int64(id))
	if err != nil {
		return err
	}
	return nil
}
