package riverqueue

import (
	"context"
	"log/slog"

	"github.com/rs/zerolog"
)

// zerologHandler adapts zerolog.Logger to slog.Handler interface
type zerologHandler struct {
	logger zerolog.Logger
}

func NewZerologAdapter(logger zerolog.Logger) *slog.Logger {
	return slog.New(&zerologHandler{logger: logger})
}

func (h *zerologHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

func (h *zerologHandler) Handle(ctx context.Context, record slog.Record) error {
	var event *zerolog.Event

	switch record.Level {
	case slog.LevelDebug:
		event = h.logger.Debug()
	case slog.LevelInfo:
		event = h.logger.Info()
	case slog.LevelWarn:
		event = h.logger.Warn()
	case slog.LevelError:
		event = h.logger.Error()
	default:
		event = h.logger.Info()
	}

	record.Attrs(func(attr slog.Attr) bool {
		event = event.Interface(attr.Key, attr.Value.Any())
		return true
	})

	event.Msg(record.Message)
	return nil
}

func (h *zerologHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	ctx := h.logger.With()
	for _, attr := range attrs {
		ctx = ctx.Interface(attr.Key, attr.Value.Any())
	}
	return &zerologHandler{logger: ctx.Logger()}
}

func (h *zerologHandler) WithGroup(name string) slog.Handler {
	return &zerologHandler{logger: h.logger.With().Str("group", name).Logger()}
}
