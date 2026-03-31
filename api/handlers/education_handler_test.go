package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cmsgov/emmy-api/pkg/core"
	"github.com/cmsgov/emmy-api/pkg/education"
	"github.com/cmsgov/emmy-api/pkg/resilience"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

type fakeEducationService struct {
	response education.Response
	err      error
	lastReq  education.Request
}

func (s *fakeEducationService) Submit(_ context.Context, req education.Request) (education.Response, error) {
	s.lastReq = req
	return s.response, s.err
}

func TestEducationHandler_SuccessWithBody(t *testing.T) {
	app := fiber.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	cfg := &core.Config{
		NSC: core.NSCConfig{
			AccountID: "10053523",
		},
	}

	svc := &fakeEducationService{
		response: education.Response{
			EnrollmentStatus: "FULL_TIME",
		},
	}
	app.Post("/edu", EducationHandler(cfg, svc, logger))

	reqBody := education.Request{
		FirstName:   "John",
		LastName:    "Doe",
		DateOfBirth: "1990-01-01",
		Terms:       "y",
		EndClient:   "CMS",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/edu", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, fiber.StatusOK, resp.StatusCode)
	require.Equal(t, "John", svc.lastReq.FirstName)
	require.Equal(t, "Doe", svc.lastReq.LastName)
	require.Equal(t, "10053523", svc.lastReq.AccountID) // Should be defaulted from config
}

func TestEducationHandler_SuccessWithQueryParams(t *testing.T) {
	app := fiber.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	cfg := &core.Config{
		NSC: core.NSCConfig{
			AccountID: "10053523",
		},
	}

	svc := &fakeEducationService{
		response: education.Response{
			EnrollmentStatus: "FULL_TIME",
		},
	}
	app.Get("/edu", EducationHandler(cfg, svc, logger))

	req := httptest.NewRequest(http.MethodGet, "/edu?firstName=Jane&lastName=Smith&dateOfBirth=1995-05-05", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, fiber.StatusOK, resp.StatusCode)
	require.Equal(t, "Jane", svc.lastReq.FirstName)
	require.Equal(t, "Smith", svc.lastReq.LastName)
	require.Equal(t, "10053523", svc.lastReq.AccountID)
}

func TestEducationHandler_UnknownEnrollmentStatus(t *testing.T) {
	app := fiber.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	cfg := &core.Config{}
	svc := &fakeEducationService{
		response: education.Response{
			EnrollmentStatus: "UNKNOWN",
		},
	}
	app.Post("/edu", EducationHandler(cfg, svc, logger))

	req := httptest.NewRequest(http.MethodPost, "/edu", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestEducationHandler_InvalidRequest(t *testing.T) {
	app := fiber.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	cfg := &core.Config{
		NSC: core.NSCConfig{},
	}

	app.Post("/edu", EducationHandler(cfg, &fakeEducationService{}, logger))

	req := httptest.NewRequest(http.MethodPost, "/edu", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Fiber might return 400 or still try query parser.
	// In our implementation, if BodyParser fails, it tries QueryParser.
	// "invalid json" body might cause BodyParser to fail, and QueryParser might "succeed" with empty values.
	// However, Fiber's BodyParser for JSON with "invalid json" content will return an error.
	require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
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
