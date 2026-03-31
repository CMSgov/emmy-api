package core

import (
	"log/slog"
	"os"

	slogmulti "github.com/samber/slog-multi"
	"go.opentelemetry.io/contrib/bridges/otelslog"
)

func newStdoutHandler(cfg *Config) slog.Handler {
	if cfg.IsProd() {
		return slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})
	}
	return slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{})
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
