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

type PaginatedHistory struct {
	Items       []HistoryEntity
	Total       int64
	HasNext     bool
	HasPrev     bool
	Limit       int
	Offset      int
	TotalPages  int
	CurrentPage int
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
