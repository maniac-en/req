// Package log provides structured file-based logging with rotation.
package log

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	*slog.Logger
	fileLogger *lumberjack.Logger
}

type Config struct {
	Level       slog.Level
	LogFilePath string
}

var (
	globalLogger *Logger
	once         sync.Once
)

func Initialize(config Config) {
	once.Do(func() {
		globalLogger = createLogger(config)
	})
}

func createLogger(config Config) *Logger {
	var fileLogger *lumberjack.Logger
	var handler slog.Handler

	if config.LogFilePath != "" {
		fileLogger = &lumberjack.Logger{
			Filename:   config.LogFilePath,
			MaxSize:    10, // MB
			MaxBackups: 2,
			MaxAge:     7, // days
			Compress:   true,
		}
		handler = slog.NewJSONHandler(fileLogger, &slog.HandlerOptions{
			Level: config.Level,
		})
	} else {
		// Fallback to stderr if no file path provided
		handler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level: config.Level,
		})
	}

	return &Logger{
		Logger:     slog.New(handler),
		fileLogger: fileLogger,
	}
}

func Global() *Logger {
	if globalLogger == nil {
		Initialize(Config{
			Level: slog.LevelInfo,
		})
	}
	return globalLogger
}

func (l *Logger) Close() error {
	if l.fileLogger != nil {
		return l.fileLogger.Close()
	}
	return nil
}

func (l *Logger) WithRequestID(requestID string) *slog.Logger {
	return l.With("request_id", requestID)
}

// Global convenience functions
func Debug(msg string, args ...any) {
	Global().Debug(msg, args...)
}

func Info(msg string, args ...any) {
	Global().Info(msg, args...)
}

func Warn(msg string, args ...any) {
	Global().Warn(msg, args...)
}

func Error(msg string, args ...any) {
	Global().Error(msg, args...)
}

func Fatal(msg string, args ...any) {
	Global().Error(msg, args...)
	os.Exit(1)
}

type contextKey string

const requestIDKey contextKey = "request_id"

func ContextWithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

func RequestIDFromContext(ctx context.Context) string {
	if reqID := ctx.Value(requestIDKey); reqID != nil {
		if id, ok := reqID.(string); ok {
			return id
		}
	}
	return ""
}
