package log

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"
	"time"
)

func TestDualHandlerColors(t *testing.T) {
	tests := []struct {
		level    slog.Level
		expected string
	}{
		{slog.LevelInfo, "\033[32m"},
		{slog.LevelError, "\033[31m"},
	}

	for _, tt := range tests {
		var buf bytes.Buffer
		handler := NewDualHandler(&buf, false, slog.LevelDebug)

		record := slog.NewRecord(time.Now(), tt.level, "test", 0)
		handler.Handle(context.Background(), record)

		output := buf.String()
		if !strings.Contains(output, tt.expected) {
			t.Errorf("expected color %s for %s level", tt.expected, tt.level)
		}
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
		t.Fatalf("multihandler failed: %v", err)
	}

	if !strings.Contains(buf1.String(), "multi test") {
		t.Error("first handler should receive message")
	}
	if !strings.Contains(buf2.String(), "multi test") {
		t.Error("second handler should receive message")
	}
}
