package history

import (
	"time"

	"github.com/maniac-en/req/internal/database"
)

type HistoryManager struct {
	DB *database.Queries
}

type HistoryEntity struct {
	database.History
}

func (h HistoryEntity) GetID() int64 {
	return h.ID
}

func (h HistoryEntity) GetName() string {
	return h.Method + " " + h.Url
}

func (h HistoryEntity) GetCreatedAt() time.Time {
	t, _ := time.Parse(time.RFC3339, h.ExecutedAt)
	return t
}

func (h HistoryEntity) GetUpdatedAt() time.Time {
	return h.GetCreatedAt()
}