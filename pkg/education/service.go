package education

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/cmsgov/emmy-api/pkg/core"
)

type EducationService interface {
	Submit(ctx context.Context, req Request) (Response, error)
}

type HTTPTransport interface {
	Do(req *http.Request) (*http.Response, error)
}

type Options struct {
	// Override for testing the HTTP client
	HTTPClient HTTPTransport
	// Structured logger using slog package
	Logger *slog.Logger
	// Context timeout
	Timeout time.Duration
}

type service struct {
	cfg    *core.NSCConfig
	client HTTPTransport
	logger *slog.Logger
	opts   Options
}

func New(cfg *core.NSCConfig, opts Options) EducationService {
	logger := opts.Logger
	if logger == nil {
		logger = slog.Default()
	}

	logger = logger.With(
		slog.String("component", "education"),
		slog.String("vendor", "nsc"),
	)

	client := opts.HTTPClient

	if client == nil {
		client = nscHTTPClient(context.Background(), cfg)
	}

	return &service{
		cfg:    cfg,
		client: client,
		logger: logger,
		opts:   opts,
	}
}
