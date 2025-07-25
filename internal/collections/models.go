package collections

import (
	"time"

	"github.com/maniac-en/req/internal/database"
)

type CollectionsManager struct {
	DB *database.Queries
}

type Collection struct {
	Name      string
	Id        int
	CreatedAt time.Time
	UpdatedAt time.Time
}
