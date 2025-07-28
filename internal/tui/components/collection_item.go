package components

import (
	"fmt"
	"strconv"

	"github.com/maniac-en/req/internal/backend/collections"
)

type CollectionItem struct {
	collection collections.CollectionEntity
}

func NewCollectionItem(collection collections.CollectionEntity) CollectionItem {
	return CollectionItem{collection: collection}
}

func (i CollectionItem) FilterValue() string {
	return i.collection.Name
}

func (i CollectionItem) GetID() string {
	return strconv.FormatInt(i.collection.ID, 10)
}

func (i CollectionItem) GetTitle() string {
	return i.collection.Name
}

func (i CollectionItem) GetDescription() string {
	return fmt.Sprintf("ID: %d", i.collection.ID)
}

func (i CollectionItem) GetCollection() collections.CollectionEntity {
	return i.collection
}
