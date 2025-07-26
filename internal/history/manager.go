package history

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/maniac-en/req/internal/crud"
	"github.com/maniac-en/req/internal/database"
	"github.com/maniac-en/req/internal/log"
)

func NewHistoryManager(db *database.Queries) *HistoryManager {
	return &HistoryManager{DB: db}
}

func (h *HistoryManager) Create(ctx context.Context, name string) (HistoryEntity, error) {
	return HistoryEntity{}, fmt.Errorf("use RecordExecution to create history entries")
}

func (h *HistoryManager) Read(ctx context.Context, id int64) (HistoryEntity, error) {
	if err := crud.ValidateID(id); err != nil {
		return HistoryEntity{}, err
	}

	history, err := h.DB.GetHistoryById(ctx, id)
	if err != nil {
		log.Error("failed to read history entry", "id", id, "error", err)
		return HistoryEntity{}, err
	}

	log.Debug("read history entry", "id", id)
	return HistoryEntity{History: history}, nil
}

func (h *HistoryManager) Update(ctx context.Context, id int64, name string) (HistoryEntity, error) {
	return HistoryEntity{}, fmt.Errorf("history entries are immutable")
}

func (h *HistoryManager) Delete(ctx context.Context, id int64) error {
	if err := crud.ValidateID(id); err != nil {
		return err
	}

	err := h.DB.DeleteHistoryEntry(ctx, id)
	if err != nil {
		log.Error("failed to delete history entry", "id", id, "error", err)
		return err
	}

	log.Info("deleted history entry", "id", id)
	return nil
}

func (h *HistoryManager) List(ctx context.Context) ([]HistoryEntity, error) {
	return nil, fmt.Errorf("use ListByCollection to list history entries")
}

func (h *HistoryManager) ListByCollection(ctx context.Context, collectionID int64, limit, offset int) (PaginatedHistory, error) {
	if err := crud.ValidateID(collectionID); err != nil {
		return PaginatedHistory{}, err
	}

	// Get total count
	total, err := h.DB.CountHistoryByCollection(ctx, sql.NullInt64{Int64: collectionID, Valid: true})
	if err != nil {
		log.Error("failed to count history by collection", "collection_id", collectionID, "error", err)
		return PaginatedHistory{}, err
	}

	// Get paginated results
	summaries, err := h.DB.GetHistoryByCollection(ctx, database.GetHistoryByCollectionParams{
		CollectionID: sql.NullInt64{Int64: collectionID, Valid: true},
		Limit:        int64(limit),
		Offset:       int64(offset),
	})
	if err != nil {
		log.Error("failed to list history by collection", "collection_id", collectionID, "error", err)
		return PaginatedHistory{}, err
	}

	// Convert to entities
	entities := make([]HistoryEntity, len(summaries))
	for i, summary := range summaries {
		entities[i] = HistoryEntity{History: database.History{
			ID:           summary.ID,
			Method:       summary.Method,
			Url:          summary.Url,
			StatusCode:   summary.StatusCode,
			ExecutedAt:   summary.ExecutedAt,
			EndpointName: summary.EndpointName,
		}}
	}

	// Calculate pagination metadata
	totalPages := int((total + int64(limit) - 1) / int64(limit)) // Ceiling division
	currentPage := (offset / limit) + 1
	hasNext := (offset + limit) < int(total)
	hasPrev := offset > 0

	result := PaginatedHistory{
		Items:       entities,
		Total:       total,
		HasNext:     hasNext,
		HasPrev:     hasPrev,
		Limit:       limit,
		Offset:      offset,
		TotalPages:  totalPages,
		CurrentPage: currentPage,
	}

	log.Info("listed history by collection", "collection_id", collectionID, "count", len(entities), "total", total, "page", currentPage, "total_pages", totalPages)
	return result, nil
}

func (h *HistoryManager) RecordExecution(ctx context.Context, data ExecutionData) (HistoryEntity, error) {
	if err := validateExecutionData(data); err != nil {
		log.Error("invalid execution data", "error", err)
		return HistoryEntity{}, err
	}

	log.Debug("recording execution", "method", data.Method, "url", data.URL, "status", data.StatusCode)

	// Marshal to JSON for database storage
	requestHeaders, err := json.Marshal(data.Headers)
	if err != nil {
		return HistoryEntity{}, fmt.Errorf("failed to marshal request headers: %w", err)
	}

	queryParams, err := json.Marshal(data.QueryParams)
	if err != nil {
		return HistoryEntity{}, fmt.Errorf("failed to marshal query params: %w", err)
	}

	responseHeaders, err := json.Marshal(data.ResponseHeaders)
	if err != nil {
		return HistoryEntity{}, fmt.Errorf("failed to marshal response headers: %w", err)
	}

	params := database.CreateHistoryEntryParams{
		CollectionID:    sql.NullInt64{Int64: data.CollectionID, Valid: data.CollectionID > 0},
		CollectionName:  sql.NullString{String: data.CollectionName, Valid: data.CollectionName != ""},
		EndpointName:    sql.NullString{String: data.EndpointName, Valid: data.EndpointName != ""},
		Method:          data.Method,
		Url:             data.URL,
		StatusCode:      int64(data.StatusCode),
		Duration:        data.Duration.Milliseconds(),
		ResponseSize:    sql.NullInt64{Int64: data.ResponseSize, Valid: data.ResponseSize > 0},
		RequestHeaders:  sql.NullString{String: string(requestHeaders), Valid: true},
		QueryParams:     sql.NullString{String: string(queryParams), Valid: true},
		RequestBody:     sql.NullString{String: data.RequestBody, Valid: data.RequestBody != ""},
		ResponseBody:    sql.NullString{String: data.ResponseBody, Valid: data.ResponseBody != ""},
		ResponseHeaders: sql.NullString{String: string(responseHeaders), Valid: true},
		ExecutedAt:      time.Now().Format(time.RFC3339),
	}

	history, err := h.DB.CreateHistoryEntry(ctx, params)
	if err != nil {
		log.Error("failed to record execution", "error", err)
		return HistoryEntity{}, err
	}

	log.Info("recorded execution", "id", history.ID, "collection_id", data.CollectionID, "status", data.StatusCode)
	return HistoryEntity{History: history}, nil
}

func validateExecutionData(data ExecutionData) error {
	if err := crud.ValidateName(data.Method); err != nil {
		log.Debug("execution validation failed: invalid method", "method", data.Method)
		return fmt.Errorf("invalid method: %w", err)
	}

	if err := crud.ValidateName(data.URL); err != nil {
		log.Debug("execution validation failed: invalid URL", "url", data.URL)
		return fmt.Errorf("invalid URL: %w", err)
	}

	if data.StatusCode < 100 || data.StatusCode > 599 {
		log.Debug("execution validation failed: invalid status code", "status_code", data.StatusCode)
		return fmt.Errorf("invalid status code: %d", data.StatusCode)
	}

	return nil
}
