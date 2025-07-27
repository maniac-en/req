package collections

import (
	"context"
	"fmt"
	"testing"

	"github.com/maniac-en/req/internal/crud"
	"github.com/maniac-en/req/internal/testutils"
)


func TestCollectionsManagerCRUD(t *testing.T) {
	db := testutils.SetupTestDB(t, "collections")
	manager := NewCollectionsManager(db)
	ctx := context.Background()

	t.Run("Create", func(t *testing.T) {
		collection, err := manager.Create(ctx, "Test Collection")
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
		if collection.GetName() != "Test Collection" {
			t.Errorf("Expected name 'Test Collection', got %s", collection.GetName())
		}
		if collection.GetID() <= 0 {
			t.Errorf("Expected positive ID, got %d", collection.GetID())
		}
	})

	t.Run("Read", func(t *testing.T) {
		created, _ := manager.Create(ctx, "Read Test")
		collection, err := manager.Read(ctx, created.GetID())
		if err != nil {
			t.Fatalf("Read failed: %v", err)
		}
		if collection.GetName() != "Read Test" {
			t.Errorf("Expected name 'Read Test', got %s", collection.GetName())
		}
	})

	t.Run("Update", func(t *testing.T) {
		created, _ := manager.Create(ctx, "Update Test")
		updated, err := manager.Update(ctx, created.GetID(), "Updated Name")
		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}
		if updated.GetName() != "Updated Name" {
			t.Errorf("Expected name 'Updated Name', got %s", updated.GetName())
		}
	})

	t.Run("Delete", func(t *testing.T) {
		created, _ := manager.Create(ctx, "Delete Test")
		err := manager.Delete(ctx, created.GetID())
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}
		_, err = manager.Read(ctx, created.GetID())
		if err != crud.ErrNotFound {
			t.Errorf("Expected ErrNotFound after delete, got %v", err)
		}
	})

	t.Run("List", func(t *testing.T) {
		manager.Create(ctx, "List Test 1")
		manager.Create(ctx, "List Test 2")
		collections, err := manager.List(ctx)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}
		if len(collections) < 2 {
			t.Errorf("Expected at least 2 collections, got %d", len(collections))
		}
	})

	t.Run("ListPaginated", func(t *testing.T) {
		for i := 1; i <= 5; i++ {
			manager.Create(ctx, fmt.Sprintf("Pagination Test %d", i))
		}

		paginated, err := manager.ListPaginated(ctx, 2, 0)
		if err != nil {
			t.Fatalf("ListPaginated failed: %v", err)
		}

		if len(paginated.Collections) != 2 {
			t.Errorf("Expected 2 collections, got %d", len(paginated.Collections))
		}
		if paginated.Total < 5 {
			t.Errorf("Expected total >= 5, got %d", paginated.Total)
		}
		if !paginated.HasNext {
			t.Error("Expected HasNext to be true")
		}
		if paginated.HasPrev {
			t.Error("Expected HasPrev to be false for offset 0")
		}

		// Test second page
		paginated2, err := manager.ListPaginated(ctx, 2, 2)
		if err != nil {
			t.Fatalf("ListPaginated page 2 failed: %v", err)
		}
		if !paginated2.HasPrev {
			t.Error("Expected HasPrev to be true for offset > 0")
		}
	})
}

func TestCollectionsManagerValidation(t *testing.T) {
	db := testutils.SetupTestDB(t, "collections")
	manager := NewCollectionsManager(db)
	ctx := context.Background()

	t.Run("Create with empty name", func(t *testing.T) {
		_, err := manager.Create(ctx, "")
		if err != crud.ErrInvalidInput {
			t.Errorf("Expected ErrInvalidInput, got %v", err)
		}
	})

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
}
