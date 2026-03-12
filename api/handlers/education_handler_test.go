package handlers

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DSACMS/verification-service-api/pkg/core"
	"github.com/DSACMS/verification-service-api/pkg/education"
	"github.com/DSACMS/verification-service-api/pkg/resilience"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

type fakeEducationService struct {
	response education.Response
	err      error
}

func (s *fakeEducationService) Submit(_ context.Context, _ education.Request) (education.Response, error) {
	return s.response, s.err
}

func TestEducationHandler_CircuitOpenReturnsServiceUnavailable(t *testing.T) {
	app := fiber.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	cfg := &core.Config{
		NSC: core.NSCConfig{
			AccountID: "10053523",
		},
	}

	app.Get("/edu", EducationHandler(cfg, &fakeEducationService{
		err: resilience.ErrCircuitOpen,
	}, logger))

	req := httptest.NewRequest(http.MethodGet, "/edu", http.NoBody)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, fiber.StatusServiceUnavailable, resp.StatusCode)
}

func TestEducationHandler_VendorErrorReturnsBadGateway(t *testing.T) {
	app := fiber.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	cfg := &core.Config{
		NSC: core.NSCConfig{
			AccountID: "10053523",
		},
	}

	app.Get("/edu", EducationHandler(cfg, &fakeEducationService{
		err: errors.New("vendor request failed"),
	}, logger))

	req := httptest.NewRequest(http.MethodGet, "/edu", http.NoBody)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, fiber.StatusBadGateway, resp.StatusCode)
}

func TestEducationHandler_EmitsVerificationSpans(t *testing.T) {
	recorder := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))
	previous := otel.GetTracerProvider()
	otel.SetTracerProvider(tp)
	t.Cleanup(func() {
		otel.SetTracerProvider(previous)
		_ = tp.Shutdown(context.Background())
	})

	app := fiber.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	cfg := &core.Config{
		NSC: core.NSCConfig{
			AccountID: "10053523",
		},
	}

	app.Get("/edu", EducationHandler(cfg, &fakeEducationService{
		response: education.Response{},
	}, logger))

	req := httptest.NewRequest(http.MethodGet, "/edu", http.NoBody)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, fiber.StatusOK, resp.StatusCode)

	spans := recorder.Ended()
	names := make([]string, 0, len(spans))
	for _, sp := range spans {
		names = append(names, sp.Name())
	}

	require.Contains(t, names, "verification.request")
	require.Contains(t, names, "decision.engine")
}
