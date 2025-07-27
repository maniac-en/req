// Package log provides structured logging with terminal and file output.
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

var (
	globalLogger *Logger
	once         sync.Once
)

// LoggerFactory creates Logger instances for dependency injection
type LoggerFactory struct{}

// NewLoggerFactory creates a new LoggerFactory
func NewLoggerFactory() *LoggerFactory {
	return &LoggerFactory{}
}

// CreateLogger creates a new Logger with the given config
func (f *LoggerFactory) CreateLogger(config Config) *Logger {
	var handlers []slog.Handler

	var fileLogger *lumberjack.Logger
	if config.LogFilePath != "" {
		fileLogger = &lumberjack.Logger{
			Filename:   config.LogFilePath,
			MaxSize:    10,
			MaxBackups: 2,
			MaxAge:     7,
			Compress:   true,
		}
		fileHandler := slog.NewJSONHandler(fileLogger, &slog.HandlerOptions{
			Level: config.Level,
		})
		handlers = append(handlers, fileHandler)
	}

	if config.Verbose {
		terminalHandler := NewDualHandler(os.Stderr, false, config.Level)
		handlers = append(handlers, terminalHandler)
	}

	var handler slog.Handler
	if len(handlers) == 1 {
		handler = handlers[0]
	} else if len(handlers) > 1 {
		handler = NewMultiHandler(handlers...)
	} else {
		handler = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: config.Level,
		})
	}

	return &Logger{
		Logger:     slog.New(handler),
		fileLogger: fileLogger,
	}
}

type Config struct {
	Level       slog.Level
	LogFilePath string
	Verbose     bool
}

func Initialize(config Config) {
	once.Do(func() {
		factory := NewLoggerFactory()
		globalLogger = factory.CreateLogger(config)
	})
}

func Global() *Logger {
	if globalLogger == nil {
		Initialize(Config{
			Level:   slog.LevelInfo,
			Verbose: false,
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

func GenerateRequestID() string {
	return fmt.Sprintf("req_%d", time.Now().UnixNano())
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
