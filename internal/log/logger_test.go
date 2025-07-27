package log

import (
	"bytes"
	"context"
	"log/slog"
	"os"
	"strings"
	"testing"
	"time"
)

func TestLoggerFactory(t *testing.T) {
	factory := NewLoggerFactory()
	
	t.Run("creates working logger", func(t *testing.T) {
		config := Config{Level: slog.LevelInfo, Verbose: true}
		logger := factory.CreateLogger(config)
		
		if logger == nil {
			t.Fatal("factory should create logger")
		}
	})
	
	t.Run("creates independent instances", func(t *testing.T) {
		config := Config{Level: slog.LevelInfo}
		logger1 := factory.CreateLogger(config)
		logger2 := factory.CreateLogger(config)
		
		if logger1 == logger2 {
			t.Error("factory should create independent instances")
		}
	})
}

func TestDualHandler(t *testing.T) {
	var buf bytes.Buffer
	handler := NewDualHandler(&buf, false, slog.LevelInfo)
	
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "test message", 0)
	record.Add("key", "value")
	
	err := handler.Handle(context.Background(), record)
	if err != nil {
		t.Fatalf("handler failed: %v", err)
	}
	
	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Error("output should contain message")
	}
	if !strings.Contains(output, "key=value") {
		t.Error("output should contain attributes")
	}
}

func TestRequestIDFunctions(t *testing.T) {
	id1 := GenerateRequestID()
	id2 := GenerateRequestID()
	
	if id1 == id2 {
		t.Error("IDs should be unique")
	}
	if !strings.HasPrefix(id1, "req_") {
		t.Error("ID should have req_ prefix")
	}
	
	ctx := ContextWithRequestID(context.Background(), "test123")
	retrieved := RequestIDFromContext(ctx)
	
	if retrieved != "test123" {
		t.Errorf("expected test123, got %s", retrieved)
	}
}

func TestGlobalLogger(t *testing.T) {
	// Test that global functions work
	Info("test info message")
	Debug("test debug message") 
	Warn("test warn message")
	Error("test error message")
	
	logger := Global()
	if logger == nil {
		t.Error("Global() should return logger")
	}
}

func TestLoggerWithFile(t *testing.T) {
	tempFile, err := os.CreateTemp("", "test*.log")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()
	
	factory := NewLoggerFactory()
	logger := factory.CreateLogger(Config{
		Level:       slog.LevelInfo,
		LogFilePath: tempFile.Name(),
	})
	
	logger.Info("test file logging")
	
	err = logger.Close()
	if err != nil {
		t.Errorf("Close should not error: %v", err)
	}
}