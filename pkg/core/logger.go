package core

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

type contextHandler struct {
	slog.Handler
}

//nolint:gocritic // slog.Handler requires slog.Record by value.
func (h *contextHandler) Handle(ctx context.Context, r slog.Record) error {
	if rid, ok := ctx.Value(RequestContextKey).(string); ok {
		r.AddAttrs(slog.String("request_id", rid))
	}
	if err := h.Handler.Handle(ctx, r); err != nil {
		return fmt.Errorf("handle slog record: %w", err)
	}
	return nil
}

func (h *contextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &contextHandler{Handler: h.Handler.WithAttrs(attrs)}
}

func (h *contextHandler) WithGroup(name string) slog.Handler {
	return &contextHandler{Handler: h.Handler.WithGroup(name)}
}

func newStdoutHandler(cfg *Config) slog.Handler {
	var handler slog.Handler
	if cfg.IsProd() {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{})
	}
	return &contextHandler{Handler: handler}
}

func NewLogger(cfg *Config) *slog.Logger {
	stdoutHandler := newStdoutHandler(cfg)
	return slog.New(stdoutHandler)
}
