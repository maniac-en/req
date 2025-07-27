// Package history tracks HTTP request execution records.
package history

import (
	"time"

	"github.com/maniac-en/req/internal/crud"
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
	return crud.ParseTimestamp(h.ExecutedAt)
}

// GetUpdatedAt returns creation time since history entries are immutable
func (h HistoryEntity) GetUpdatedAt() time.Time {
	return h.GetCreatedAt()
}

type PaginatedHistory struct {
	Items []HistoryEntity `json:"items"`
	crud.PaginationMetadata
}

type ExecutionData struct {
	CollectionID    int64
	CollectionName  string
	EndpointName    string
	Method          string
	URL             string
	Headers         map[string]string
	QueryParams     map[string]string
	RequestBody     string
	StatusCode      int
	ResponseBody    string
	ResponseHeaders map[string][]string
	Duration        time.Duration
	ResponseSize    int64
}
