// Package testutils provides shared testing utilities for backend packages.
package testutils

import (
	"context"
	"database/sql"
	"testing"

	"github.com/maniac-en/req/internal/backend/database"
	_ "github.com/mattn/go-sqlite3"
)

// SetupTestDB creates an in-memory SQLite database with specified tables
func SetupTestDB(t *testing.T, tables ...string) *database.Queries {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Create tables based on the requested tables
	for _, table := range tables {
		schema := getTableSchema(table)
		if schema == "" {
			t.Fatalf("Unknown table: %s", table)
		}

		if _, err := db.Exec(schema); err != nil {
			t.Fatalf("Failed to create %s table: %v", table, err)
		}
	}

	return database.New(db)
}

// getTableSchema returns the SQL schema for creating the specified table
func getTableSchema(table string) string {
	schemas := map[string]string{
		"collections": `
			CREATE TABLE collections (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				name TEXT NOT NULL,
				created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
				updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
			);`,
		"endpoints": `
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
			);`,
		"history": `
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
			);`,
	}

	return schemas[table]
}

// CreateTestCollection creates a test collection and returns its ID
func CreateTestCollection(t *testing.T, db *database.Queries, name string) int64 {
	collection, err := db.CreateCollection(context.Background(), name)
	if err != nil {
		t.Fatalf("Failed to create test collection: %v", err)
	}
	return collection.ID
}
