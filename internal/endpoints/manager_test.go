package endpoints

import (
	"context"
	"database/sql"
	"testing"

	"github.com/maniac-en/req/internal/crud"
	"github.com/maniac-en/req/internal/database"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *database.Queries {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	schema := `
	CREATE TABLE collections (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE endpoints (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		collection_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		method TEXT NOT NULL,
		url TEXT NOT NULL,
		headers TEXT DEFAULT '{}' NOT NULL,
		query_params TEXT DEFAULT '{}' NOT NULL,
		request_body TEXT DEFAULT '' NOT NULL,
		created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (collection_id) REFERENCES collections(id) ON DELETE CASCADE
	);`

	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("Failed to create test schema: %v", err)
	}

	return database.New(db)
}

func createTestCollection(t *testing.T, db *database.Queries) int64 {
	collection, err := db.CreateCollection(context.Background(), "Test Collection")
	if err != nil {
		t.Fatalf("Failed to create test collection: %v", err)
	}
	return collection.ID
}

func TestEndpointsManagerCRUD(t *testing.T) {
	db := setupTestDB(t)
	manager := NewEndpointsManager(db)
	ctx := context.Background()
	collectionID := createTestCollection(t, db)

	t.Run("Create returns error", func(t *testing.T) {
		_, err := manager.Create(ctx, "Test Endpoint")
		if err == nil {
			t.Error("Expected Create to return error")
		}
	})

	t.Run("Read", func(t *testing.T) {
		// Create endpoint using CreateEndpoint
		data := EndpointData{
			CollectionID: collectionID,
			Name:         "Read Test Endpoint",
			Method:       "GET",
			URL:          "https://api.example.com/test",
			Headers:      `{"Content-Type": "application/json"}`,
			QueryParams:  map[string]string{"param": "value"},
			RequestBody:  "",
		}
		created, err := manager.CreateEndpoint(ctx, data)
		if err != nil {
			t.Fatalf("CreateEndpoint failed: %v", err)
		}

		endpoint, err := manager.Read(ctx, created.GetID())
		if err != nil {
			t.Fatalf("Read failed: %v", err)
		}
		if endpoint.GetName() != "Read Test Endpoint" {
			t.Errorf("Expected name 'Read Test Endpoint', got %s", endpoint.GetName())
		}
		if endpoint.Method != "GET" {
			t.Errorf("Expected method 'GET', got %s", endpoint.Method)
		}
	})

	t.Run("Update returns error", func(t *testing.T) {
		_, err := manager.Update(ctx, 1, "Updated Name")
		if err == nil {
			t.Error("Expected Update to return error")
		}
	})

	t.Run("Delete", func(t *testing.T) {
		// Create endpoint to delete
		data := EndpointData{
			CollectionID: collectionID,
			Name:         "Delete Test Endpoint",
			Method:       "DELETE",
			URL:          "https://api.example.com/delete",
			Headers:      "{}",
			QueryParams:  map[string]string{},
			RequestBody:  "",
		}
		created, err := manager.CreateEndpoint(ctx, data)
		if err != nil {
			t.Fatalf("CreateEndpoint failed: %v", err)
		}

		err = manager.Delete(ctx, created.GetID())
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		_, err = manager.Read(ctx, created.GetID())
		if err != crud.ErrNotFound {
			t.Errorf("Expected ErrNotFound after delete, got %v", err)
		}
	})

	t.Run("List returns error", func(t *testing.T) {
		_, err := manager.List(ctx)
		if err == nil {
			t.Error("Expected List to return error")
		}
	})
}

func TestCreateEndpoint(t *testing.T) {
	db := setupTestDB(t)
	manager := NewEndpointsManager(db)
	ctx := context.Background()
	collectionID := createTestCollection(t, db)

	t.Run("Valid endpoint creation", func(t *testing.T) {
		data := EndpointData{
			CollectionID: collectionID,
			Name:         "Test API Endpoint",
			Method:       "POST",
			URL:          "https://api.example.com/users",
			Headers:      `{"Content-Type": "application/json", "Authorization": "Bearer token"}`,
			QueryParams:  map[string]string{"format": "json", "version": "v1"},
			RequestBody:  `{"name": "John", "email": "john@example.com"}`,
		}

		endpoint, err := manager.CreateEndpoint(ctx, data)
		if err != nil {
			t.Fatalf("CreateEndpoint failed: %v", err)
		}

		if endpoint.GetName() != "Test API Endpoint" {
			t.Errorf("Expected name 'Test API Endpoint', got %s", endpoint.GetName())
		}
		if endpoint.Method != "POST" {
			t.Errorf("Expected method 'POST', got %s", endpoint.Method)
		}
		if endpoint.Url != "https://api.example.com/users" {
			t.Errorf("Expected URL 'https://api.example.com/users', got %s", endpoint.Url)
		}
		if endpoint.CollectionID != collectionID {
			t.Errorf("Expected collection_id %d, got %d", collectionID, endpoint.CollectionID)
		}
	})

	t.Run("Invalid collection ID", func(t *testing.T) {
		data := EndpointData{
			CollectionID: -1,
			Name:         "Test Endpoint",
			Method:       "GET",
			URL:          "https://api.example.com",
		}

		_, err := manager.CreateEndpoint(ctx, data)
		if err != crud.ErrInvalidInput {
			t.Errorf("Expected ErrInvalidInput, got %v", err)
		}
	})

	t.Run("Empty name", func(t *testing.T) {
		data := EndpointData{
			CollectionID: collectionID,
			Name:         "",
			Method:       "GET",
			URL:          "https://api.example.com",
		}

		_, err := manager.CreateEndpoint(ctx, data)
		if err != crud.ErrInvalidInput {
			t.Errorf("Expected ErrInvalidInput, got %v", err)
		}
	})

	t.Run("Empty method", func(t *testing.T) {
		data := EndpointData{
			CollectionID: collectionID,
			Name:         "Test Endpoint",
			Method:       "",
			URL:          "https://api.example.com",
		}

		_, err := manager.CreateEndpoint(ctx, data)
		if err != crud.ErrInvalidInput {
			t.Errorf("Expected ErrInvalidInput, got %v", err)
		}
	})

	t.Run("Empty URL", func(t *testing.T) {
		data := EndpointData{
			CollectionID: collectionID,
			Name:         "Test Endpoint",
			Method:       "GET",
			URL:          "",
		}

		_, err := manager.CreateEndpoint(ctx, data)
		if err != crud.ErrInvalidInput {
			t.Errorf("Expected ErrInvalidInput, got %v", err)
		}
	})
}

func TestUpdateEndpoint(t *testing.T) {
	db := setupTestDB(t)
	manager := NewEndpointsManager(db)
	ctx := context.Background()
	collectionID := createTestCollection(t, db)

	t.Run("Valid endpoint update", func(t *testing.T) {
		// Create endpoint first
		data := EndpointData{
			CollectionID: collectionID,
			Name:         "Original Endpoint",
			Method:       "GET",
			URL:          "https://api.example.com/original",
			Headers:      "{}",
			QueryParams:  map[string]string{},
			RequestBody:  "",
		}
		created, err := manager.CreateEndpoint(ctx, data)
		if err != nil {
			t.Fatalf("CreateEndpoint failed: %v", err)
		}

		// Update the endpoint
		updateData := EndpointData{
			Name:        "Updated Endpoint",
			Method:      "PUT",
			URL:         "https://api.example.com/updated",
			Headers:     `{"Content-Type": "application/json"}`,
			QueryParams: map[string]string{"updated": "true"},
			RequestBody: `{"updated": true}`,
		}

		updated, err := manager.UpdateEndpoint(ctx, created.GetID(), updateData)
		if err != nil {
			t.Fatalf("UpdateEndpoint failed: %v", err)
		}

		if updated.GetName() != "Updated Endpoint" {
			t.Errorf("Expected name 'Updated Endpoint', got %s", updated.GetName())
		}
		if updated.Method != "PUT" {
			t.Errorf("Expected method 'PUT', got %s", updated.Method)
		}
		if updated.Url != "https://api.example.com/updated" {
			t.Errorf("Expected URL 'https://api.example.com/updated', got %s", updated.Url)
		}
	})

	t.Run("Update non-existent endpoint", func(t *testing.T) {
		data := EndpointData{
			Name:   "Test",
			Method: "GET",
			URL:    "https://api.example.com",
		}

		_, err := manager.UpdateEndpoint(ctx, 99999, data)
		if err != crud.ErrNotFound {
			t.Errorf("Expected ErrNotFound, got %v", err)
		}
	})

	t.Run("Update with invalid ID", func(t *testing.T) {
		data := EndpointData{
			Name:   "Test",
			Method: "GET",
			URL:    "https://api.example.com",
		}

		_, err := manager.UpdateEndpoint(ctx, -1, data)
		if err != crud.ErrInvalidInput {
			t.Errorf("Expected ErrInvalidInput, got %v", err)
		}
	})
}

func TestListByCollection(t *testing.T) {
	db := setupTestDB(t)
	manager := NewEndpointsManager(db)
	ctx := context.Background()
	collectionID := createTestCollection(t, db)

	// Create test endpoints
	endpoints := []EndpointData{
		{CollectionID: collectionID, Name: "Endpoint 1", Method: "GET", URL: "https://api.example.com/1"},
		{CollectionID: collectionID, Name: "Endpoint 2", Method: "POST", URL: "https://api.example.com/2"},
		{CollectionID: collectionID, Name: "Endpoint 3", Method: "PUT", URL: "https://api.example.com/3"},
		{CollectionID: collectionID, Name: "Endpoint 4", Method: "DELETE", URL: "https://api.example.com/4"},
		{CollectionID: collectionID, Name: "Endpoint 5", Method: "PATCH", URL: "https://api.example.com/5"},
	}

	for _, data := range endpoints {
		_, err := manager.CreateEndpoint(ctx, data)
		if err != nil {
			t.Fatalf("Failed to create endpoint %s: %v", data.Name, err)
		}
	}

	t.Run("Valid pagination", func(t *testing.T) {
		result, err := manager.ListByCollection(ctx, collectionID, 2, 0)
		if err != nil {
			t.Fatalf("ListByCollection failed: %v", err)
		}

		if len(result.Endpoints) != 2 {
			t.Errorf("Expected 2 endpoints, got %d", len(result.Endpoints))
		}
		if result.Total != 5 {
			t.Errorf("Expected total 5, got %d", result.Total)
		}
		if !result.HasNext {
			t.Error("Expected HasNext to be true")
		}
		if result.HasPrev {
			t.Error("Expected HasPrev to be false for offset 0")
		}
		if result.TotalPages != 3 {
			t.Errorf("Expected 3 total pages, got %d", result.TotalPages)
		}
		if result.CurrentPage != 1 {
			t.Errorf("Expected current page 1, got %d", result.CurrentPage)
		}
	})

	t.Run("Second page", func(t *testing.T) {
		result, err := manager.ListByCollection(ctx, collectionID, 2, 2)
		if err != nil {
			t.Fatalf("ListByCollection failed: %v", err)
		}

		if len(result.Endpoints) != 2 {
			t.Errorf("Expected 2 endpoints, got %d", len(result.Endpoints))
		}
		if !result.HasNext {
			t.Error("Expected HasNext to be true")
		}
		if !result.HasPrev {
			t.Error("Expected HasPrev to be true for offset > 0")
		}
		if result.CurrentPage != 2 {
			t.Errorf("Expected current page 2, got %d", result.CurrentPage)
		}
	})

	t.Run("Last page", func(t *testing.T) {
		result, err := manager.ListByCollection(ctx, collectionID, 2, 4)
		if err != nil {
			t.Fatalf("ListByCollection failed: %v", err)
		}

		if len(result.Endpoints) != 1 {
			t.Errorf("Expected 1 endpoint, got %d", len(result.Endpoints))
		}
		if result.HasNext {
			t.Error("Expected HasNext to be false for last page")
		}
		if !result.HasPrev {
			t.Error("Expected HasPrev to be true")
		}
		if result.CurrentPage != 3 {
			t.Errorf("Expected current page 3, got %d", result.CurrentPage)
		}
	})

	t.Run("Invalid collection ID", func(t *testing.T) {
		_, err := manager.ListByCollection(ctx, -1, 10, 0)
		if err != crud.ErrInvalidInput {
			t.Errorf("Expected ErrInvalidInput, got %v", err)
		}
	})

	t.Run("Empty collection", func(t *testing.T) {
		emptyCollectionID := createTestCollection(t, db)
		result, err := manager.ListByCollection(ctx, emptyCollectionID, 10, 0)
		if err != nil {
			t.Fatalf("ListByCollection failed: %v", err)
		}

		if len(result.Endpoints) != 0 {
			t.Errorf("Expected 0 endpoints, got %d", len(result.Endpoints))
		}
		if result.Total != 0 {
			t.Errorf("Expected total 0, got %d", result.Total)
		}
		if result.HasNext {
			t.Error("Expected HasNext to be false for empty collection")
		}
		if result.HasPrev {
			t.Error("Expected HasPrev to be false for empty collection")
		}
	})
}

func TestEndpointsManagerValidation(t *testing.T) {
	db := setupTestDB(t)
	manager := NewEndpointsManager(db)
	ctx := context.Background()

	t.Run("Read with invalid ID", func(t *testing.T) {
		_, err := manager.Read(ctx, -1)
		if err != crud.ErrInvalidInput {
			t.Errorf("Expected ErrInvalidInput, got %v", err)
		}
	})

	t.Run("Read non-existent", func(t *testing.T) {
		_, err := manager.Read(ctx, 99999)
		if err != crud.ErrNotFound {
			t.Errorf("Expected ErrNotFound, got %v", err)
		}
	})

	t.Run("Delete with invalid ID", func(t *testing.T) {
		err := manager.Delete(ctx, -1)
		if err != crud.ErrInvalidInput {
			t.Errorf("Expected ErrInvalidInput, got %v", err)
		}
	})
}