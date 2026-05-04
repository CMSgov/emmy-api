package routes

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/cmsgov/emmy-api/pkg/core"
	"github.com/cmsgov/emmy-api/pkg/education"
	"github.com/cmsgov/emmy-api/pkg/reporting"
	"github.com/cmsgov/emmy-api/pkg/veteran"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

type stubEducationService struct{}

func (s *stubEducationService) LookupEnrollmentStatus(_ context.Context, _ education.Request) (education.Response, error) {
	return education.Response{}, nil
}

func (s *stubEducationService) RegisterBatch(_ context.Context, _ education.BatchRequest) error {
	return nil
}

func (s *stubEducationService) GetBatchStatus(_ context.Context, _ string) (education.BatchJobStatusResponse, error) {
	return education.BatchJobStatusResponse{}, nil
}

func (s *stubEducationService) GetBatchDetails(_ context.Context, _ string) (education.BatchJobDetailsResponse, error) {
	return education.BatchJobDetailsResponse{}, nil
}

type stubVeteranService struct{}

func (s *stubVeteranService) LookupDisabilityRating(_ context.Context, _ veteran.Request) (veteran.Response, error) {
	return veteran.Response{}, nil
}

func TestRegisterRoutes_RegistersOpenAPISpecEndpoint(t *testing.T) {
	app := fiber.New()
	logger := slog.New(slog.DiscardHandler)
	reporter := reporting.NewMockReporter()

	RegisterRoutes(app, RouterParams{
		CFG:      &core.Config{},
		Reporter: reporter,
		Logger:   logger,
		EDU:      &stubEducationService{},
		VA:       &stubVeteranService{},
		WithCB: func(next fiber.Handler) fiber.Handler {
			return next
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api-spec/v1/verify", http.NoBody)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	expectedBody, err := os.ReadFile(filepath.Join("..", "..", "api-spec", "v0", "dist", "openapi.bundled.json"))
	require.NoError(t, err)

	require.Equal(t, fiber.StatusOK, resp.StatusCode)
	require.Contains(t, resp.Header.Get(fiber.HeaderContentType), fiber.MIMEApplicationJSON)
	require.JSONEq(t, string(expectedBody), string(body))
}

func TestRegisterRoutes_RegistersVeteranDisabilityEndpoint(t *testing.T) {
	app := fiber.New()
	logger := slog.New(slog.DiscardHandler)
	reporter := reporting.NewMockReporter()

	RegisterRoutes(app, RouterParams{
		CFG:      &core.Config{},
		Reporter: reporter,
		Logger:   logger,
		EDU:      &stubEducationService{},
		VA:       &stubVeteranService{},
		WithCB: func(next fiber.Handler) fiber.Handler {
			return next
		},
	})

	routes := app.GetRoutes(true)
	for _, route := range routes {
		if route.Method == http.MethodPost && route.Path == "/api/v0/veteran-disability-ratings" {
			return
		}
	}

	t.Fatal("expected POST /api/v0/veteran-disability-ratings to be registered")
}

func TestRegisterRoutes_RegistersEducationEnrollmentsEndpoint(t *testing.T) {
	app := fiber.New()
	logger := slog.New(slog.DiscardHandler)
	reporter := reporting.NewMockReporter()

	RegisterRoutes(app, RouterParams{
		CFG:      &core.Config{},
		Reporter: reporter,
		Logger:   logger,
		EDU:      &stubEducationService{},
		VA:       &stubVeteranService{},
		WithCB: func(next fiber.Handler) fiber.Handler {
			return next
		},
	})

	routes := app.GetRoutes(true)
	for _, route := range routes {
		if route.Method == http.MethodPost && route.Path == "/api/v0/education-enrollments" {
			return
		}
	}

	t.Fatal("expected POST /api/v0/education-enrollments to be registered")
}
