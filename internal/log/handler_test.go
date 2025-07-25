package log

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"
	"time"
)

func TestDualHandlerColoring(t *testing.T) {
	tests := []struct {
		level    slog.Level
		expected string
	}{
		{slog.LevelDebug, "\033[36m"},
		{slog.LevelInfo, "\033[32m"},
		{slog.LevelWarn, "\033[33m"},
		{slog.LevelError, "\033[31m"},
	}

	for _, tt := range tests {
		t.Run(tt.level.String(), func(t *testing.T) {
			var buf bytes.Buffer
			handler := NewDualHandler(&buf, false, slog.LevelDebug)

			record := slog.NewRecord(time.Now(), tt.level, "test", 0)
			handler.Handle(context.Background(), record)

			output := buf.String()
			if !strings.Contains(output, tt.expected) {
				t.Errorf("expected color code %s in output, got: %s", tt.expected, output)
			}
		})
	}
}

func TestDualHandlerWithAttrs(t *testing.T) {
	var buf bytes.Buffer
	handler := NewDualHandler(&buf, false, slog.LevelInfo)
	
	// add some attributes to the handler
	handlerWithAttrs := handler.WithAttrs([]slog.Attr{
		slog.String("service", "test"),
		slog.Int("version", 1),
	})

	record := slog.NewRecord(time.Now(), slog.LevelInfo, "test message", 0)
	record.Add("request_id", "123")

	handlerWithAttrs.Handle(context.Background(), record)

	output := buf.String()
	if !strings.Contains(output, "service=test") {
		t.Error("should contain handler attribute service=test")
	}
	if !strings.Contains(output, "version=1") {
		t.Error("should contain handler attribute version=1")
	}
	if !strings.Contains(output, "request_id=123") {
		t.Error("should contain record attribute request_id=123")
	}
}

func TestDualHandlerWithGroup(t *testing.T) {
	var buf bytes.Buffer
	handler := NewDualHandler(&buf, false, slog.LevelInfo)
	
	groupHandler := handler.WithGroup("http")
	if groupHandler == nil {
		t.Error("WithGroup should return a handler")
	}

	// test that it creates a new handler instance
	if groupHandler == handler {
		t.Error("WithGroup should return a new handler instance")
	}
}

func TestMultiHandlerEnabled(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	
	// one handler with INFO level, another with ERROR level
	handler1 := NewDualHandler(&buf1, false, slog.LevelInfo)
	handler2 := NewDualHandler(&buf2, false, slog.LevelError)
	
	multiHandler := NewMultiHandler(handler1, handler2)

	// should be enabled for INFO (handler1 accepts it)
	if !multiHandler.Enabled(context.Background(), slog.LevelInfo) {
		t.Error("should be enabled for INFO level")
	}

	// should be enabled for ERROR (both handlers accept it)
	if !multiHandler.Enabled(context.Background(), slog.LevelError) {
		t.Error("should be enabled for ERROR level")
	}

	// should not be enabled for DEBUG (neither handler accepts it)
	if multiHandler.Enabled(context.Background(), slog.LevelDebug) {
		t.Error("should not be enabled for DEBUG level")
	}
}

func TestMultiHandlerWithAttrs(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	
	handler1 := NewDualHandler(&buf1, false, slog.LevelInfo)
	handler2 := NewDualHandler(&buf2, false, slog.LevelInfo)
	
	multiHandler := NewMultiHandler(handler1, handler2)
	
	// add attributes
	handlerWithAttrs := multiHandler.WithAttrs([]slog.Attr{
		slog.String("component", "test"),
	})

	record := slog.NewRecord(time.Now(), slog.LevelInfo, "test", 0)
	handlerWithAttrs.Handle(context.Background(), record)

	output1 := buf1.String()
	output2 := buf2.String()

	if !strings.Contains(output1, "component=test") {
		t.Error("first handler should contain the attribute")
	}
	if !strings.Contains(output2, "component=test") {
		t.Error("second handler should contain the attribute")
	}
}

func TestMultiHandlerWithGroup(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	
	handler1 := NewDualHandler(&buf1, false, slog.LevelInfo)
	handler2 := NewDualHandler(&buf2, false, slog.LevelInfo)
	
	multiHandler := NewMultiHandler(handler1, handler2)
	groupHandler := multiHandler.WithGroup("database")

	if groupHandler == nil {
		t.Error("WithGroup should return a handler")
	}
	if groupHandler == multiHandler {
		t.Error("WithGroup should return a new instance")
	}
}

func TestHandlerTimestampFormat(t *testing.T) {
	var buf bytes.Buffer
	handler := NewDualHandler(&buf, false, slog.LevelInfo)

	now := time.Date(2025, 1, 15, 14, 30, 45, 0, time.UTC)
	record := slog.NewRecord(now, slog.LevelInfo, "timestamp test", 0)

	handler.Handle(context.Background(), record)
	
	output := buf.String()
	expected := "2025-01-15 14:30:45"
	if !strings.Contains(output, expected) {
		t.Errorf("expected timestamp format %s in output: %s", expected, output)
	}
}