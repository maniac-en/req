package main

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/maniac-en/req/internal/database"
	"github.com/maniac-en/req/internal/tui"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

// Embed migration files into the binary
//
//go:embed db/migrations/*.sql
var migrationsFS embed.FS

var (
	USERHOMEDIR string
	DBDIR       string
	DBPATH      string
	DB          *sql.DB
)

type Config struct {
	DB *database.Queries
}

func init() {
	// setup DB path based on user's home path
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("error reading user's home path")
	}
	USERHOMEDIR = userHomeDir
	DBDIR = filepath.Join(USERHOMEDIR, ".cache", "req")
	if err := os.MkdirAll(DBDIR, 0o755); err != nil {
		log.Fatal("error creating a database directory")
	}
	DBPATH = filepath.Join(DBDIR, "app.db")

	// Run migrations
	if err := runMigrations(); err != nil {
		log.Fatalf("error running migrations: %v", err)
	}
}

func runMigrations() error {
	// Connect to database
	db, err := sql.Open("sqlite3", DBPATH)
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}

	// Create a sub-filesystem for just the migrations directory
	migrationSubFS, err := fs.Sub(migrationsFS, filepath.Join("db", "migrations"))
	if err != nil {
		return fmt.Errorf("error creating sub-filesystem: %w", err)
	}

	// Create goose provider with the embedded filesystem
	gooseProvider, err := goose.NewProvider("sqlite3", db, migrationSubFS)
	if err != nil {
		return fmt.Errorf("error creating goose provider: %w", err)
	}

	// Run migrations
	_, err = gooseProvider.Up(context.Background())
	if err != nil {
		return fmt.Errorf("error running migrations: %w", err)
	}
	DB = db
	return nil
}

func main() {
	db := database.New(DB)
	cfg := Config{
		DB: db,
	}
	_, err := cfg.DB.CreateCollection(context.Background(), "testing")
	if err != nil {
		log.Fatal(err)
	}

	// create tabs and model
	tabs := tui.InitTabs()
	model, err := tui.InitModel(tabs)
	// its really hard for this to throw an error rn, but i want this here
	// in case our programme is able to create some creative errors on init
	// in the future
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
