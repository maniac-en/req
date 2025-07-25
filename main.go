package main

import (
	"context"
	"database/sql"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/maniac-en/req/internal/collections"
	"github.com/maniac-en/req/internal/database"
	"github.com/maniac-en/req/internal/log"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

// Embed migration files into the binary
//
//go:embed db/migrations/*.sql
var migrationsFS embed.FS

var (
	USERHOMEDIR string
	APPDIR      string
	DBPATH      string
	LOGPATH     string
	DB          *sql.DB
)

type Config struct {
	DB          *database.Queries
	Collections *collections.CollectionsManager
}

func initPaths() error {
	// setup paths based on user's home directory
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error reading user's home path: %w", err)
	}
	USERHOMEDIR = userHomeDir
	APPDIR = filepath.Join(USERHOMEDIR, ".cache", "req")
	if err := os.MkdirAll(APPDIR, 0o755); err != nil {
		return fmt.Errorf("error creating app directory: %w", err)
	}
	DBPATH = filepath.Join(APPDIR, "app.db")
	LOGPATH = filepath.Join(APPDIR, "req.log")
	return nil
}

func runMigrations() error {
	// connect to database
	db, err := sql.Open("sqlite3", DBPATH)
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}

	// create sub-filesystem for migrations
	migrationSubFS, err := fs.Sub(migrationsFS, filepath.Join("db", "migrations"))
	if err != nil {
		return fmt.Errorf("error creating sub-filesystem: %w", err)
	}

	// create goose provider
	gooseProvider, err := goose.NewProvider("sqlite3", db, migrationSubFS)
	if err != nil {
		return fmt.Errorf("error creating goose provider: %w", err)
	}

	// run migrations
	_, err = gooseProvider.Up(context.Background())
	if err != nil {
		return fmt.Errorf("error running migrations: %w", err)
	}
	DB = db
	return nil
}

func main() {
	verbose := flag.Bool("verbose", false, "enable verbose logging to terminal")
	flag.Parse()

	// initialize paths first
	if err := initPaths(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize: %v\n", err)
		os.Exit(1)
	}

	// initialize logging
	logLevel := slog.LevelInfo
	if *verbose {
		logLevel = slog.LevelDebug
	}
	log.Initialize(log.Config{
		Level:       logLevel,
		LogFilePath: LOGPATH,
		Verbose:     *verbose,
	})
	defer func() {
		if err := log.Global().Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to close logger: %v\n", err)
		}
	}()

	log.Info("starting req application")

	// run database migrations
	if err := runMigrations(); err != nil {
		log.Fatal("failed to run migrations", "error", err)
	}

	// create database client and collections manager
	db := database.New(DB)
	collectionsManager := collections.NewCollectionsManager(db)
	
	_ = &Config{
		DB:          db,
		Collections: collectionsManager,
	}

	log.Info("application initialized", "components", []string{"database", "collections", "logging"})
	log.Info("application started successfully")
}
