package collections_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/maniac-en/req/internal/collections"
	"github.com/maniac-en/req/internal/database"
	_ "github.com/mattn/go-sqlite3"
)

const schema = `
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
);
`

func setupTestDB(t *testing.T) (*sql.DB, *database.Queries, func()) {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	// schema, err := os.ReadFile("testdata/schema.sql")
	// if err != nil {
	// 	t.Fatalf("failed to read schema: %v", err)
	// }
	_, err = db.Exec(string(schema))
	if err != nil {
		t.Fatalf("failed to execute schema: %v", err)
	}

	queries := database.New(db)

	cleanup := func() {
		db.Close()
	}
	return db, queries, cleanup
}

func TestCreateCollection(t *testing.T) {
	_, queries, cleanup := setupTestDB(t)
	defer cleanup()

	manager := collections.NewCollectionsManager(queries)
	ctx := context.Background()

	collection, err := manager.CreateCollection(ctx, "Test Collection")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if collection.Name != "Test Collection" {
		t.Errorf("expected name to be 'Test Collection', got %s", collection.Name)
	}
}

func TestCreateCollection_Validation(t *testing.T) {
	_, queries, cleanup := setupTestDB(t)
	defer cleanup()
	manager := collections.NewCollectionsManager(queries)
	ctx := context.Background()

	_, err := manager.CreateCollection(ctx, "")
	if err == nil {
		t.Errorf("expected error for empty name")
	}

	longName := make([]byte, 101)
	for i := range longName {
		longName[i] = 'a'
	}
	_, err = manager.CreateCollection(ctx, string(longName))
	if err == nil {
		t.Errorf("expected error for long name")
	}

	_, err = manager.CreateCollection(ctx, "Invalid/Name")
	if err == nil {
		t.Errorf("expected error for invalid characters")
	}
}

func TestUpdateCollectionName(t *testing.T) {
	_, queries, cleanup := setupTestDB(t)
	defer cleanup()
	manager := collections.NewCollectionsManager(queries)
	ctx := context.Background()

	collection, _ := manager.CreateCollection(ctx, "Old Name")
	updated, err := manager.UpdateCollectionName(ctx, "New Name", int(collection.ID))
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if updated.Name != "New Name" {
		t.Errorf("expected name to be 'New Name', got %s", updated.Name)
	}
}

func TestDeleteCollection(t *testing.T) {
	_, queries, cleanup := setupTestDB(t)
	defer cleanup()
	manager := collections.NewCollectionsManager(queries)
	ctx := context.Background()

	collection, _ := manager.CreateCollection(ctx, "ToDelete")
	err := manager.DeleteCollection(ctx, int(collection.ID))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	collections, _ := manager.GetAllCollections(ctx)
	if len(collections) != 0 {
		t.Errorf("expected 0 collections after delete")
	}
}

func TestGetAllCollections(t *testing.T) {
	_, queries, cleanup := setupTestDB(t)
	defer cleanup()
	manager := collections.NewCollectionsManager(queries)
	ctx := context.Background()

	manager.CreateCollection(ctx, "One")
	manager.CreateCollection(ctx, "Two")

	all, err := manager.GetAllCollections(ctx)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(all) != 2 {
		t.Errorf("expected 2 collections, got %d", len(all))
	}
}

func TestCreateEndpoint(t *testing.T) {
	_, queries, cleanup := setupTestDB(t)
	defer cleanup()
	manager := collections.NewCollectionsManager(queries)
	ctx := context.Background()

	collection, _ := manager.CreateCollection(ctx, "API Group")
	endpoint, err := manager.CreateEndpoint(ctx, "Get Users", "GET", "/users", collection.ID)
	if err != nil {
		t.Fatalf("endpoint creation failed: %v", err)
	}
	if endpoint.Name != "Get Users" {
		t.Errorf("expected endpoint name 'Get Users', got %s", endpoint.Name)
	}
}

func TestUpdateEndpoint(t *testing.T) {
	_, queries, cleanup := setupTestDB(t)
	defer cleanup()
	manager := collections.NewCollectionsManager(queries)
	ctx := context.Background()

	collection, _ := manager.CreateCollection(ctx, "Test")
	ep, _ := manager.CreateEndpoint(ctx, "Original", "GET", "/", collection.ID)

	updated, err := manager.UpdateEndpoint(ctx, ep.ID, "Updated", "POST", "/new", "hdr", "q", "body")
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if updated.Method != "POST" || updated.Url != "/new" {
		t.Errorf("update did not apply correctly")
	}
}

func TestGetEndpoint(t *testing.T) {
	_, queries, cleanup := setupTestDB(t)
	defer cleanup()
	manager := collections.NewCollectionsManager(queries)
	ctx := context.Background()

	collection, _ := manager.CreateCollection(ctx, "Test")
	ep, _ := manager.CreateEndpoint(ctx, "TestEp", "GET", "/test", collection.ID)
	fetched, err := manager.GetEndpoint(ctx, ep.ID)
	if err != nil {
		t.Fatalf("failed to fetch: %v", err)
	}
	if fetched.ID != ep.ID {
		t.Errorf("expected ID %d, got %d", ep.ID, fetched.ID)
	}
}

func TestGetEndpoints(t *testing.T) {
	_, queries, cleanup := setupTestDB(t)
	defer cleanup()
	manager := collections.NewCollectionsManager(queries)
	ctx := context.Background()

	collection, _ := manager.CreateCollection(ctx, "Batch")
	manager.CreateEndpoint(ctx, "One", "GET", "/1", collection.ID)
	manager.CreateEndpoint(ctx, "Two", "GET", "/2", collection.ID)

	eps, err := manager.GetEndpoints(ctx, collection.ID)
	if err != nil {
		t.Fatalf("expected no error: %v", err)
	}
	if len(eps) != 2 {
		t.Errorf("expected 2 endpoints, got %d", len(eps))
	}
}

func TestCreateEndpoint_InvalidCollectionID(t *testing.T) {
	_, queries, cleanup := setupTestDB(t)
	defer cleanup()
	manager := collections.NewCollectionsManager(queries)
	ctx := context.Background()

	_, err := manager.CreateEndpoint(ctx, "Name", "GET", "/url", 999)
	if err == nil {
		t.Errorf("expected error for invalid collection ID")
	}
}

func TestUpdateEndpoint_NonExistent(t *testing.T) {
	_, queries, cleanup := setupTestDB(t)
	defer cleanup()
	manager := collections.NewCollectionsManager(queries)
	ctx := context.Background()

	_, err := manager.UpdateEndpoint(ctx, 999, "Updated", "POST", "/new", "hdr", "q", "body")
	if err == nil {
		t.Errorf("expected error for non-existent endpoint")
	}
}

func TestGetEndpoint_NotFound(t *testing.T) {
	_, queries, cleanup := setupTestDB(t)
	defer cleanup()
	manager := collections.NewCollectionsManager(queries)
	ctx := context.Background()

	_, err := manager.GetEndpoint(ctx, 999)
	if err == nil {
		t.Errorf("expected error for non-existent endpoint")
	}
}

func newBrokenTestDB(t *testing.T) *database.Queries {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open sqlite memory db: %v", err)
	}
	db.Close()
	return database.New(db)
}

func TestGetEndpoints_DBError(t *testing.T) {
	db := newBrokenTestDB(t)

	manager := collections.NewCollectionsManager(db)
	ctx := context.Background()

	_, err := manager.GetEndpoints(ctx, 999)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	t.Logf("got expected error: %v", err)
}

func TestDeleteCollection_InvalidID(t *testing.T) {
	_, queries, cleanup := setupTestDB(t)
	defer cleanup()
	manager := collections.NewCollectionsManager(queries)
	ctx := context.Background()

	err := manager.DeleteCollection(ctx, 999)
	if err != nil {
		t.Logf("received expected delete error: %v", err)
	}
}
