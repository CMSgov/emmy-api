package core

import (
	"context"
	"log/slog"
	"os"

	slogmulti "github.com/samber/slog-multi"
	"go.opentelemetry.io/contrib/bridges/otelslog"
)

type contextHandler struct {
	slog.Handler
}

func (h *contextHandler) Handle(ctx context.Context, r slog.Record) error {
	if rid, ok := ctx.Value(RequestContextKey).(string); ok {
		r.AddAttrs(slog.String("request_id", rid))
	}
	return h.Handler.Handle(ctx, r)
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

func NewLoggerWithOtel(cfg *Config, otel OtelService) *slog.Logger {
	stdoutHandler := newStdoutHandler(cfg)
	otelHandler := otelslog.NewHandler(
		"emmy-api",
		otelslog.WithLoggerProvider(otel.LoggerProvider()),
	)

	return slog.New(
		slogmulti.Fanout(
			stdoutHandler,
			otelHandler,
		),
	)
}
