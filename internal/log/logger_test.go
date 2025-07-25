package log

import (
	"bytes"
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestLoggerInitialization(t *testing.T) {
	tests := []struct {
		name   string
		config Config
	}{
		{
			name: "file only",
			config: Config{
				Level:       slog.LevelInfo,
				LogFilePath: "/tmp/test.log",
				Verbose:     false,
			},
		},
		{
			name: "verbose only",
			config: Config{
				Level:   slog.LevelDebug,
				Verbose: true,
			},
		},
		{
			name: "both file and verbose",
			config: Config{
				Level:       slog.LevelInfo,
				LogFilePath: "/tmp/test2.log",
				Verbose:     true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// reset global logger for each test
			resetOnce()
			globalLogger = nil

			Initialize(tt.config)
			logger := Global()

			if logger == nil {
				t.Fatal("logger should not be nil")
			}

			if logger.Logger == nil {
				t.Fatal("embedded slog.Logger should not be nil")
			}

			// cleanup
			if tt.config.LogFilePath != "" {
				_ = os.Remove(tt.config.LogFilePath)
			}
		})
	}
}

func TestDualHandlerFormatting(t *testing.T) {
	var buf bytes.Buffer
	handler := NewDualHandler(&buf, false, slog.LevelInfo)

	record := slog.NewRecord(time.Now(), slog.LevelInfo, "test message", 0)
	record.Add("key", "value")

	err := handler.Handle(context.Background(), record)
	if err != nil {
		t.Fatalf("Handle failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Error("output should contain the message")
	}
	if !strings.Contains(output, "key=value") {
		t.Error("output should contain attributes")
	}
	if !strings.Contains(output, "[INFO]") {
		t.Error("output should contain log level")
	}
}

func TestLogLevels(t *testing.T) {
	var buf bytes.Buffer
	handler := NewDualHandler(&buf, false, slog.LevelWarn)

	// should not log (below threshold)
	if handler.Enabled(context.Background(), slog.LevelInfo) {
		t.Error("should not be enabled for info level when threshold is warn")
	}

	// should log (at threshold)
	warnRecord := slog.NewRecord(time.Now(), slog.LevelWarn, "warn message", 0)
	if !handler.Enabled(context.Background(), slog.LevelWarn) {
		t.Error("should be enabled for warn level")
	}

	_ = handler.Handle(context.Background(), warnRecord)
	output := buf.String()
	if !strings.Contains(output, "warn message") {
		t.Error("should contain warn message")
	}
}

func TestMultiHandler(t *testing.T) {
	var buf1, buf2 bytes.Buffer

	handler1 := NewDualHandler(&buf1, false, slog.LevelInfo)
	handler2 := NewDualHandler(&buf2, false, slog.LevelInfo)

	multiHandler := NewMultiHandler(handler1, handler2)

	record := slog.NewRecord(time.Now(), slog.LevelInfo, "multi test", 0)
	err := multiHandler.Handle(context.Background(), record)
	if err != nil {
		t.Fatalf("MultiHandler.Handle failed: %v", err)
	}

	output1 := buf1.String()
	output2 := buf2.String()

	if !strings.Contains(output1, "multi test") {
		t.Error("first handler should receive the message")
	}
	if !strings.Contains(output2, "multi test") {
		t.Error("second handler should receive the message")
	}
}

func TestRequestIDFunctions(t *testing.T) {
	// test request ID generation
	id1 := GenerateRequestID()
	id2 := GenerateRequestID()

	if id1 == id2 {
		t.Error("generated IDs should be unique")
	}
	if !strings.HasPrefix(id1, "req_") {
		t.Error("request ID should have req_ prefix")
	}

	// test context functions
	ctx := context.Background()
	testID := "test123"

	ctxWithID := ContextWithRequestID(ctx, testID)
	retrievedID := RequestIDFromContext(ctxWithID)

	if retrievedID != testID {
		t.Errorf("expected %s, got %s", testID, retrievedID)
	}

	// test with context without request ID
	emptyID := RequestIDFromContext(ctx)
	if emptyID != "" {
		t.Error("should return empty string for context without request ID")
	}
}

func TestGlobalLoggerFunctions(t *testing.T) {
	// reset global logger
	resetOnce()
	globalLogger = nil

	tempFile := filepath.Join(os.TempDir(), "test_global.log")
	defer func() { _ = os.Remove(tempFile) }()

	Initialize(Config{
		Level:       slog.LevelDebug,
		LogFilePath: tempFile,
		Verbose:     false,
	})

	// test all log level functions
	Debug("debug message", "key", "debug")
	Info("info message", "key", "info")
	Warn("warn message", "key", "warn")
	Error("error message", "key", "error")

	// read the log file
	content, err := os.ReadFile(tempFile)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	logContent := string(content)
	if !strings.Contains(logContent, "debug message") {
		t.Error("should contain debug message")
	}
	if !strings.Contains(logContent, "info message") {
		t.Error("should contain info message")
	}
	if !strings.Contains(logContent, "warn message") {
		t.Error("should contain warn message")
	}
	if !strings.Contains(logContent, "error message") {
		t.Error("should contain error message")
	}
}

func TestLoggerClose(t *testing.T) {
	tempFile := filepath.Join(os.TempDir(), "test_close.log")
	defer func() { _ = os.Remove(tempFile) }()

	// reset global logger
	resetOnce()
	globalLogger = nil

	Initialize(Config{
		LogFilePath: tempFile,
		Level:       slog.LevelInfo,
	})

	logger := Global()
	err := logger.Close()
	if err != nil {
		t.Errorf("Close should not return error: %v", err)
	}
}

func TestWithRequestID(t *testing.T) {
	// reset global logger
	resetOnce()
	globalLogger = nil

	Initialize(Config{Level: slog.LevelInfo})
	logger := Global()

	requestLogger := logger.WithRequestID("req_123")
	if requestLogger == nil {
		t.Error("WithRequestID should return a logger")
	}
}

// resetOnce resets the sync.Once for testing
func resetOnce() {
	once = sync.Once{}
}
