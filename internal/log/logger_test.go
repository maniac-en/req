package log

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"testing"
)

func TestInitialize(t *testing.T) {
	t.Run("creates working logger with file", func(t *testing.T) {
		tempFile, err := os.CreateTemp("", "test*.log")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		defer os.Remove(tempFile.Name())
		tempFile.Close()

		Initialize(Config{
			Level:       slog.LevelInfo,
			LogFilePath: tempFile.Name(),
		})

		logger := Global()
		if logger == nil {
			t.Fatal("Initialize should create global logger")
		}

		// Test that we can log
		logger.Info("test message")

		err = logger.Close()
		if err != nil {
			t.Errorf("Close should not error: %v", err)
		}
	})

	t.Run("creates logger without file path", func(t *testing.T) {
		Initialize(Config{Level: slog.LevelInfo})

		logger := Global()
		if logger == nil {
			t.Fatal("Initialize should create global logger even without file path")
		}
	})
}

func TestRequestIDFunctions(t *testing.T) {
	t.Run("generates unique request IDs", func(t *testing.T) {
		id1 := GenerateRequestID()
		id2 := GenerateRequestID()

		if id1 == id2 {
			t.Error("IDs should be unique")
		}
		if !strings.HasPrefix(id1, "req_") {
			t.Error("ID should have req_ prefix")
		}
	})

	t.Run("context request ID functions", func(t *testing.T) {
		ctx := ContextWithRequestID(context.Background(), "test123")
		retrieved := RequestIDFromContext(ctx)

		if retrieved != "test123" {
			t.Errorf("expected test123, got %s", retrieved)
		}
	})

	t.Run("request ID from empty context", func(t *testing.T) {
		retrieved := RequestIDFromContext(context.Background())
		if retrieved != "" {
			t.Errorf("expected empty string, got %s", retrieved)
		}
	})
}

func TestGlobalLoggerFunctions(t *testing.T) {
	t.Run("global functions work", func(t *testing.T) {
		// These should not panic
		Info("test info message")
		Debug("test debug message")
		Warn("test warn message")
		Error("test error message")

		logger := Global()
		if logger == nil {
			t.Error("Global() should return logger")
		}
	})

	t.Run("with request ID", func(t *testing.T) {
		logger := Global()
		loggerWithID := logger.WithRequestID("test-req-123")

		if loggerWithID == nil {
			t.Error("WithRequestID should return logger")
		}
	})
}