package main

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"runtime/debug"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/maniac-en/req/internal/backend/collections"
	"github.com/maniac-en/req/internal/backend/database"
	"github.com/maniac-en/req/internal/backend/demo"
	"github.com/maniac-en/req/internal/backend/endpoints"
	"github.com/maniac-en/req/internal/backend/history"
	"github.com/maniac-en/req/internal/backend/http"
	"github.com/maniac-en/req/internal/log"
	"github.com/maniac-en/req/internal/tui/app"
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

var Version = "dev"

func getVersion() string {
	// Try to get version from build info (works with go install)
	if info, ok := debug.ReadBuildInfo(); ok {
		if info.Main.Version != "" && info.Main.Version != "(devel)" {
			return info.Main.Version
		}
	}

	// Fall back to injected version (release builds)
	return Version
}

func initPaths() error {
	// setup paths using OS-appropriate cache directory
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error reading user's home path: %w", err)
	}
	USERHOMEDIR = userHomeDir

	// use OS-appropriate cache directory
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		return fmt.Errorf("error reading user's cache path: %w", err)
	}
	APPDIR = filepath.Join(userCacheDir, "req")
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
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Error("failed to close migration database connection", "error", closeErr)
		}
	}()

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

	// open a new connection for the global DB
	globalDB, err := sql.Open("sqlite3", DBPATH)
	if err != nil {
		return fmt.Errorf("error opening global database connection: %w", err)
	}
	DB = globalDB
	return nil
}

func main() {
	// initialize paths first
	if err := initPaths(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize: %v\n", err)
		os.Exit(1)
	}

	// initialize logging
	logLevel := slog.LevelInfo
	if os.Getenv("REQ_DEBUG") == "1" || os.Getenv("REQ_LOG_LEVEL") == "debug" {
		logLevel = slog.LevelDebug
	}
	log.Initialize(log.Config{
		Level:       logLevel,
		LogFilePath: LOGPATH,
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

	// create database client and managers
	db := database.New(DB)
	collectionsManager := collections.NewCollectionsManager(db)
	endpointsManager := endpoints.NewEndpointsManager(db)
	httpManager := http.NewHTTPManager()
	historyManager := history.NewHistoryManager(db)

	// create clean context for dependency injection
	appContext := app.NewContext(
		collectionsManager,
		endpointsManager,
		httpManager,
		historyManager,
		getVersion(),
	)

	// populate dummy data for demo
	demoGenerator := demo.NewDemoGenerator(collectionsManager, endpointsManager)
	dummyDataCreated, err := demoGenerator.PopulateDummyData(context.Background())
	if err != nil {
		log.Error("failed to populate dummy data", "error", err)
	} else if dummyDataCreated {
		// appContext.SetDummyDataCreated(true)
	}

	log.Info("application initialized", "components", []string{"database", "collections", "endpoints", "http", "history", "logging", "demo"})
	log.Debug("configuration loaded", "collections_manager", collectionsManager != nil, "endpoints", endpointsManager != nil, "database", db != nil, "http_manager", httpManager != nil, "history_manager", historyManager != nil)
	log.Info("application started successfully")

	// Entry point for UI
	program := tea.NewProgram(app.NewAppModel(appContext), tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		log.Fatal("Fatal error:", err)
	}
}
