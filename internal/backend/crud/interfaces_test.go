package crud

import (
	"context"
	"testing"
	"time"
)

type TestEntity struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (e TestEntity) GetID() int64 {
	return e.ID
}

func (e TestEntity) GetName() string {
	return e.Name
}

func (e TestEntity) GetCreatedAt() time.Time {
	return e.CreatedAt
}

func (e TestEntity) GetUpdatedAt() time.Time {
	return e.UpdatedAt
}

func TestCRUDInterface(t *testing.T) {
	ctx := context.Background()

	t.Run("Create", func(t *testing.T) {
		manager := &testCRUDManager{}

		entity, err := manager.Create(ctx, "test-entity")
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		if entity.GetName() != "test-entity" {
			t.Errorf("expected name 'test-entity', got %s", entity.GetName())
		}

		if entity.GetID() == 0 {
			t.Error("expected ID to be set")
		}

		if entity.GetCreatedAt().IsZero() {
			t.Error("expected CreatedAt to be set")
		}
	})

	t.Run("Read", func(t *testing.T) {
		manager := &testCRUDManager{}

		// Create entity first
		created, err := manager.Create(ctx, "test-read")
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		// Read it back
		entity, err := manager.Read(ctx, created.GetID())
		if err != nil {
			t.Fatalf("Read failed: %v", err)
		}

		if entity.GetID() != created.GetID() {
			t.Errorf("expected ID %d, got %d", created.GetID(), entity.GetID())
		}

		if entity.GetName() != "test-read" {
			t.Errorf("expected name 'test-read', got %s", entity.GetName())
		}
	})

	t.Run("Update", func(t *testing.T) {
		manager := &testCRUDManager{}

		// Create entity first
		created, err := manager.Create(ctx, "original-name")
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		// Update it
		updated, err := manager.Update(ctx, created.GetID(), "updated-name")
		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		if updated.GetName() != "updated-name" {
			t.Errorf("expected name 'updated-name', got %s", updated.GetName())
		}

		if updated.GetUpdatedAt().Before(created.GetUpdatedAt()) {
			t.Error("expected UpdatedAt to be updated")
		}
	})

	t.Run("Delete", func(t *testing.T) {
		manager := &testCRUDManager{}

		// Create entity first
		created, err := manager.Create(ctx, "to-delete")
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		// Delete it
		err = manager.Delete(ctx, created.GetID())
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		// Verify it's gone
		_, err = manager.Read(ctx, created.GetID())
		if err == nil {
			t.Error("expected Read to fail after Delete")
		}
	})

	t.Run("List", func(t *testing.T) {
		manager := &testCRUDManager{}

		// Create some entities
		names := []string{"entity1", "entity2", "entity3"}
		for _, name := range names {
			_, err := manager.Create(ctx, name)
			if err != nil {
				t.Fatalf("Create failed: %v", err)
			}
		}

		// List them
		entities, err := manager.List(ctx)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}

		if len(entities) < len(names) {
			t.Errorf("expected at least %d entities, got %d", len(names), len(entities))
		}
	})
}

func TestCRUDValidation(t *testing.T) {
	ctx := context.Background()
	manager := &testCRUDManager{}

	t.Run("Create with empty name", func(t *testing.T) {
		_, err := manager.Create(ctx, "")
		if err == nil {
			t.Error("expected Create with empty name to fail")
		}
	})

	t.Run("Read non-existent entity", func(t *testing.T) {
		_, err := manager.Read(ctx, 99999)
		if err == nil {
			t.Error("expected Read of non-existent entity to fail")
		}
	})

	t.Run("Update non-existent entity", func(t *testing.T) {
		_, err := manager.Update(ctx, 99999, "new-name")
		if err == nil {
			t.Error("expected Update of non-existent entity to fail")
		}
	})

	t.Run("Delete non-existent entity", func(t *testing.T) {
		err := manager.Delete(ctx, 99999)
		if err == nil {
			t.Error("expected Delete of non-existent entity to fail")
		}
	})
}

type testCRUDManager struct {
	entities map[int64]TestEntity
	nextID   int64
}

func (m *testCRUDManager) Create(ctx context.Context, name string) (TestEntity, error) {
	if m.entities == nil {
		m.entities = make(map[int64]TestEntity)
		m.nextID = 1
	}

	if name == "" {
		return TestEntity{}, ErrInvalidInput
	}

	now := time.Now()
	entity := TestEntity{
		ID:        m.nextID,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	m.entities[m.nextID] = entity
	m.nextID++

	return entity, nil
}

func (m *testCRUDManager) Read(ctx context.Context, id int64) (TestEntity, error) {
	if m.entities == nil {
		return TestEntity{}, ErrNotFound
	}

	entity, exists := m.entities[id]
	if !exists {
		return TestEntity{}, ErrNotFound
	}

	return entity, nil
}

func (m *testCRUDManager) Update(ctx context.Context, id int64, name string) (TestEntity, error) {
	if m.entities == nil {
		return TestEntity{}, ErrNotFound
	}

	entity, exists := m.entities[id]
	if !exists {
		return TestEntity{}, ErrNotFound
	}

	if name == "" {
		return TestEntity{}, ErrInvalidInput
	}

	entity.Name = name
	entity.UpdatedAt = time.Now()
	m.entities[id] = entity

	return entity, nil
}

func (m *testCRUDManager) Delete(ctx context.Context, id int64) error {
	if m.entities == nil {
		return ErrNotFound
	}

	_, exists := m.entities[id]
	if !exists {
		return ErrNotFound
	}

	delete(m.entities, id)
	return nil
}

func (m *testCRUDManager) List(ctx context.Context) ([]TestEntity, error) {
	if m.entities == nil {
		return []TestEntity{}, nil
	}

	entities := make([]TestEntity, 0, len(m.entities))
	for _, entity := range m.entities {
		entities = append(entities, entity)
	}

	return entities, nil
}
