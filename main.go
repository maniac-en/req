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

	tea "github.com/charmbracelet/bubbletea"
	"github.com/maniac-en/req/global"
	"github.com/maniac-en/req/internal/app"
	"github.com/maniac-en/req/internal/collections"
	"github.com/maniac-en/req/internal/database"
	"github.com/maniac-en/req/internal/endpoints"
	"github.com/maniac-en/req/internal/history"
	"github.com/maniac-en/req/internal/http"
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
	Endpoints   *endpoints.EndpointsManager
	HTTP        *http.HTTPManager
	History     *history.HistoryManager
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

	config := &Config{
		DB:          db,
		Collections: collectionsManager,
		Endpoints:   endpointsManager,
		HTTP:        httpManager,
		History:     historyManager,
	}
	appContext := &global.AppContext{
		Collections: collectionsManager,
		Endpoints:   endpointsManager,
		HTTP:        httpManager,
		History:     historyManager,
	}
	global.SetAppContext(appContext)

	log.Info("application initialized", "components", []string{"database", "collections", "endpoints", "http", "history", "logging"})
	log.Debug("configuration loaded", "collections_manager", config.Collections != nil, "endpoints", config.Endpoints != nil, "database", config.DB != nil, "http_manager", config.HTTP != nil, "history_manager", config.History != nil)
	log.Info("application started successfully")

	// Entry point for UI
	program := tea.NewProgram(app.InitialModel(), tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		log.Fatal("Fatal error:", err)
	}
}
