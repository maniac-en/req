package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
)

type DualHandler struct {
	output io.Writer
	level  slog.Level
	attrs  []slog.Attr
	groups []string
}

func NewDualHandler(output io.Writer, _ bool, level slog.Level) *DualHandler {
	return &DualHandler{
		output: output,
		level:  level,
	}
}

func (h *DualHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *DualHandler) Handle(ctx context.Context, record slog.Record) error {
	// format the log message
	output := h.formatRecord(record)

	// write to output
	_, err := fmt.Fprintln(h.output, output)
	return err
}

func (h *DualHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]slog.Attr, len(h.attrs)+len(attrs))
	copy(newAttrs, h.attrs)
	copy(newAttrs[len(h.attrs):], attrs)

	return &DualHandler{
		output: h.output,
		level:  h.level,
		attrs:  newAttrs,
		groups: h.groups,
	}
}

func (h *DualHandler) WithGroup(name string) slog.Handler {
	newGroups := make([]string, len(h.groups)+1)
	copy(newGroups, h.groups)
	newGroups[len(h.groups)] = name

	return &DualHandler{
		output: h.output,
		level:  h.level,
		attrs:  h.attrs,
		groups: newGroups,
	}
}

func (h *DualHandler) formatRecord(record slog.Record) string {
	return h.formatHuman(record)
}

func (h *DualHandler) formatHuman(record slog.Record) string {
	var color string
	switch record.Level {
	case slog.LevelDebug:
		color = "\033[36m"
	case slog.LevelInfo:
		color = "\033[32m"
	case slog.LevelWarn:
		color = "\033[33m"
	case slog.LevelError:
		color = "\033[31m"
	default:
		color = "\033[0m"
	}

	timestamp := record.Time.Format("2006-01-02 15:04:05")
	output := fmt.Sprintf("%s[%s]%s %s %s",
		color, record.Level.String(), "\033[0m", timestamp, record.Message)

	var attrs []string
	record.Attrs(func(attr slog.Attr) bool {
		attrs = append(attrs, fmt.Sprintf("%s=%v", attr.Key, attr.Value))
		return true
	})
	for _, attr := range h.attrs {
		attrs = append(attrs, fmt.Sprintf("%s=%v", attr.Key, attr.Value))
	}

	if len(attrs) > 0 {
		output += fmt.Sprintf(" {%s}", strings.Join(attrs, ", "))
	}

	return output
}

type MultiHandler struct {
	handlers []slog.Handler
}

func NewMultiHandler(handlers ...slog.Handler) *MultiHandler {
	return &MultiHandler{handlers: handlers}
}

func (h *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (h *MultiHandler) Handle(ctx context.Context, record slog.Record) error {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, record.Level) {
			if err := handler.Handle(ctx, record); err != nil {
				return err
			}
		}
	}
	return nil
}

func (h *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		newHandlers[i] = handler.WithAttrs(attrs)
	}
	return &MultiHandler{handlers: newHandlers}
}

func (h *MultiHandler) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		newHandlers[i] = handler.WithGroup(name)
	}
	return &MultiHandler{handlers: newHandlers}
}
