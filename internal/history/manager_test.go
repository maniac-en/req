package history

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/maniac-en/req/internal/database"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *database.Queries {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory database: %v", err)
	}

	// Use same schema as migration
	schema := `
	CREATE TABLE history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		collection_id INTEGER,
		collection_name TEXT,
		endpoint_name TEXT,
		method TEXT NOT NULL,
		url TEXT NOT NULL,
		status_code INTEGER NOT NULL,
		duration INTEGER NOT NULL,
		response_size INTEGER DEFAULT 0,
		request_headers TEXT DEFAULT '{}',
		query_params TEXT DEFAULT '{}',
		request_body TEXT DEFAULT '',
		response_body TEXT DEFAULT '',
		response_headers TEXT DEFAULT '{}',
		executed_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}

	return database.New(db)
}


func TestHistoryManagerCRUD(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)
	manager := NewHistoryManager(db)
	
	t.Run("Create returns error", func(t *testing.T) {
		_, err := manager.Create(ctx, "test")
		if err == nil {
			t.Error("expected Create to return error")
		}
	})
	
	t.Run("Update returns error", func(t *testing.T) {
		_, err := manager.Update(ctx, 1, "test")
		if err == nil {
			t.Error("expected Update to return error")
		}
	})
	
	t.Run("List returns error", func(t *testing.T) {
		_, err := manager.List(ctx)
		if err == nil {
			t.Error("expected List to return error")
		}
	})
	
	t.Run("Read with invalid ID", func(t *testing.T) {
		_, err := manager.Read(ctx, 0)
		if err == nil {
			t.Error("expected Read with invalid ID to return error")
		}
	})
	
	t.Run("Delete with invalid ID", func(t *testing.T) {
		err := manager.Delete(ctx, 0)
		if err == nil {
			t.Error("expected Delete with invalid ID to return error")
		}
	})
}

func TestRecordExecution(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)
	manager := NewHistoryManager(db)
	
	t.Run("valid execution", func(t *testing.T) {
		data := ExecutionData{
			CollectionID:   1,
			CollectionName: "Test Collection",
			EndpointName:   "Get Users",
			Method:         "GET",
			URL:            "https://api.example.com/users",
			Headers:        map[string]string{"Authorization": "Bearer token"},
			QueryParams:    map[string]string{"limit": "10"},
			RequestBody:    "",
			StatusCode:     200,
			ResponseBody:   `{"users": []}`,
			ResponseHeaders: map[string][]string{"Content-Type": {"application/json"}},
			Duration:       150 * time.Millisecond,
			ResponseSize:   100,
		}
		
		entity, err := manager.RecordExecution(ctx, data)
		if err != nil {
			t.Fatalf("RecordExecution failed: %v", err)
		}
		
		if entity.GetID() == 0 {
			t.Error("expected entity to have ID")
		}
		
		if entity.Method != "GET" {
			t.Errorf("expected method GET, got %s", entity.Method)
		}
		
		// Test Read functionality
		readEntity, err := manager.Read(ctx, entity.GetID())
		if err != nil {
			t.Fatalf("Read failed: %v", err)
		}
		
		if readEntity.GetID() != entity.GetID() {
			t.Errorf("expected ID %d, got %d", entity.GetID(), readEntity.GetID())
		}
	})
	
	t.Run("invalid execution data", func(t *testing.T) {
		tests := []struct {
			name string
			data ExecutionData
		}{
			{
				name: "empty method",
				data: ExecutionData{Method: "", URL: "https://example.com", StatusCode: 200},
			},
			{
				name: "empty URL",
				data: ExecutionData{Method: "GET", URL: "", StatusCode: 200},
			},
			{
				name: "invalid status code",
				data: ExecutionData{Method: "GET", URL: "https://example.com", StatusCode: 999},
			},
		}
		
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := manager.RecordExecution(ctx, tt.data)
				if err == nil {
					t.Error("expected RecordExecution to fail")
				}
			})
		}
	})
}

func TestListByCollection(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)
	manager := NewHistoryManager(db)
	
	// Create test data
	testData := []ExecutionData{
		{CollectionID: 1, Method: "GET", URL: "https://api.example.com/users", StatusCode: 200},
		{CollectionID: 1, Method: "POST", URL: "https://api.example.com/users", StatusCode: 201},
		{CollectionID: 1, Method: "GET", URL: "https://api.example.com/posts", StatusCode: 200},
		{CollectionID: 2, Method: "GET", URL: "https://api.other.com/data", StatusCode: 200},
	}
	
	for _, data := range testData {
		_, err := manager.RecordExecution(ctx, data)
		if err != nil {
			t.Fatalf("failed to create test data: %v", err)
		}
	}
	
	t.Run("valid pagination", func(t *testing.T) {
		result, err := manager.ListByCollection(ctx, 1, 2, 0)
		if err != nil {
			t.Fatalf("ListByCollection failed: %v", err)
		}
		
		if result.Total != 3 {
			t.Errorf("expected total 3, got %d", result.Total)
		}
		
		if result.TotalPages != 2 {
			t.Errorf("expected 2 total pages, got %d", result.TotalPages)
		}
		
		if result.CurrentPage != 1 {
			t.Errorf("expected current page 1, got %d", result.CurrentPage)
		}
		
		if !result.HasNext {
			t.Error("expected HasNext to be true")
		}
		
		if result.HasPrev {
			t.Error("expected HasPrev to be false")
		}
		
		if len(result.Items) != 2 {
			t.Errorf("expected 2 items, got %d", len(result.Items))
		}
	})
	
	t.Run("invalid collection ID", func(t *testing.T) {
		_, err := manager.ListByCollection(ctx, 0, 10, 0)
		if err == nil {
			t.Error("expected ListByCollection with invalid ID to fail")
		}
	})
	
	t.Run("second page", func(t *testing.T) {
		result, err := manager.ListByCollection(ctx, 1, 2, 2)
		if err != nil {
			t.Fatalf("ListByCollection failed: %v", err)
		}
		
		if result.CurrentPage != 2 {
			t.Errorf("expected current page 2, got %d", result.CurrentPage)
		}
		
		if !result.HasPrev {
			t.Error("expected HasPrev to be true")
		}
		
		if result.HasNext {
			t.Error("expected HasNext to be false")
		}
		
		if len(result.Items) != 1 {
			t.Errorf("expected 1 item, got %d", len(result.Items))
		}
	})
	
	t.Run("collection filtering", func(t *testing.T) {
		result, err := manager.ListByCollection(ctx, 2, 10, 0)
		if err != nil {
			t.Fatalf("ListByCollection failed: %v", err)
		}
		
		if result.Total != 1 {
			t.Errorf("expected total 1, got %d", result.Total)
		}
		
		if len(result.Items) != 1 {
			t.Errorf("expected 1 item, got %d", len(result.Items))
		}
	})
}

func TestDelete(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)
	manager := NewHistoryManager(db)
	
	// Create test entry
	data := ExecutionData{
		CollectionID: 1,
		Method:       "GET",
		URL:          "https://example.com",
		StatusCode:   200,
	}
	
	entity, err := manager.RecordExecution(ctx, data)
	if err != nil {
		t.Fatalf("failed to create test data: %v", err)
	}
	
	// Delete it
	err = manager.Delete(ctx, entity.GetID())
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	
	// Verify it's gone
	_, err = manager.Read(ctx, entity.GetID())
	if err == nil {
		t.Error("expected Read to fail after Delete")
	}
}

