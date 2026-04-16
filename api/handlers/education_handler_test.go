package handlers

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/cmsgov/emmy-api/pkg/core"
	"github.com/cmsgov/emmy-api/pkg/education"
	"github.com/cmsgov/emmy-api/pkg/reporting"
	"github.com/cmsgov/emmy-api/pkg/resilience"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

type fakeReporter struct {
	calls []*reporting.ReportData
}

func (r *fakeReporter) Report(_ context.Context, data *reporting.ReportData) {
	r.calls = append(r.calls, data)
}

type fakeEducationService struct {
	response education.Response
	err      error
	calls    int
	lastReq  education.Request
}

func (s *fakeEducationService) LookupEnrollmentStatus(_ context.Context, req education.Request) (education.Response, error) {
	s.calls++
	s.lastReq = req
	return s.response, s.err
}

func TestEducationHandler_Success(t *testing.T) {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		rid := c.Get(fiber.HeaderXRequestID)
		ctx := context.WithValue(c.UserContext(), core.RequestContextKey, rid)
		c.SetUserContext(ctx)
		return c.Next()
	})
	cfg := &core.Config{Environment: "test", ServiceVersion: "1.3.0"}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	service := &fakeEducationService{
		response: education.Response{EnrollmentStatus: education.EnrollmentStatusUnknown},
	}
	reporter := &fakeReporter{}

	app.Post("/api/v0/education-enrollments", EducationHandler(cfg, service, reporter, logger))

	req := httptest.NewRequest(http.MethodPost, "/api/v0/education-enrollments", strings.NewReader(`{
		"firstName":"Lynette",
		"middleName":"Marie",
		"lastName":"Oyola",
		"dateOfBirth":"1988-10-24",
		"ssn":"123-45-6789"
	}`))
	req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)
	req.Header.Set(fiber.HeaderXRequestID, "test-request-id")

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	require.Equal(t, fiber.StatusOK, resp.StatusCode)
	require.Contains(t, string(body), `"enrollmentStatus":"ENROLLMENT_STATUS_UNKNOWN"`)
	require.Contains(t, string(body), `"metadata"`)
	require.Contains(t, string(body), `"apiVersion":"1.3.0"`)
	require.Contains(t, string(body), `"environment":"test"`)
	require.Contains(t, string(body), `"transaction-id":"test-request-id"`)
	require.Equal(t, 1, service.calls)
	require.Equal(t, "Lynette", service.lastReq.FirstName)
	require.Equal(t, "Marie", service.lastReq.MiddleName)
	require.Equal(t, "123-45-6789", service.lastReq.SSN)

	require.Len(t, reporter.calls, 1)
	require.True(t, reporter.calls[0].Success)
	require.Equal(t, "NSC", reporter.calls[0].DataSource)
	require.Equal(t, fiber.StatusOK, reporter.calls[0].StatusCode)
	require.WithinDuration(t, time.Now(), reporter.calls[0].Timestamp, 2*time.Second)
}

func TestEducationHandler_InvalidJSONReturnsBadRequest(t *testing.T) {
	app := fiber.New()
	cfg := &core.Config{}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	reporter := &fakeReporter{}

	app.Post("/api/v0/education-enrollments", EducationHandler(cfg, &fakeEducationService{}, reporter, logger))

	req := httptest.NewRequest(http.MethodPost, "/api/v0/education-enrollments", strings.NewReader(`{`))
	req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestEducationHandler_MissingRequiredFieldReturnsBadRequest(t *testing.T) {
	app := fiber.New()
	cfg := &core.Config{}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	reporter := &fakeReporter{}

	app.Post("/api/v0/education-enrollments", EducationHandler(cfg, &fakeEducationService{}, reporter, logger))

	req := httptest.NewRequest(http.MethodPost, "/api/v0/education-enrollments", strings.NewReader(`{
		"lastName":"Oyola",
		"dateOfBirth":"1988-10-24",
		"ssn":"123-45-6789"
	}`))
	req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestEducationHandler_NotFoundReturnsNotFound(t *testing.T) {
	app := fiber.New()
	cfg := &core.Config{}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	reporter := &fakeReporter{}

	app.Post("/api/v0/education-enrollments", EducationHandler(cfg, &fakeEducationService{
		err: education.ErrNotFound,
	}, reporter, logger))

	req := httptest.NewRequest(http.MethodPost, "/api/v0/education-enrollments", strings.NewReader(`{
		"firstName":"Lynette",
		"lastName":"Oyola",
		"dateOfBirth":"1988-10-24",
		"ssn":"123-45-6789"
	}`))
	req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	require.Len(t, reporter.calls, 1)
	require.False(t, reporter.calls[0].Success)
}

func TestEducationHandler_CircuitOpenReturnsServiceUnavailable(t *testing.T) {
	app := fiber.New()
	cfg := &core.Config{}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	app.Post("/api/v0/education-enrollments", EducationHandler(cfg, &fakeEducationService{
		err: resilience.ErrCircuitOpen,
	}, &fakeReporter{}, logger))

	req := httptest.NewRequest(http.MethodPost, "/api/v0/education-enrollments", strings.NewReader(`{
		"firstName":"Lynette",
		"lastName":"Oyola",
		"dateOfBirth":"1988-10-24",
		"ssn":"123-45-6789"
	}`))
	req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, fiber.StatusServiceUnavailable, resp.StatusCode)
}

func TestEducationHandler_VendorErrorReturnsBadGateway(t *testing.T) {
	app := fiber.New()
	cfg := &core.Config{}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	app.Post("/api/v0/education-enrollments", EducationHandler(cfg, &fakeEducationService{
		err: errors.New("vendor request failed"),
	}, &fakeReporter{}, logger))

	req := httptest.NewRequest(http.MethodPost, "/api/v0/education-enrollments", strings.NewReader(`{
		"firstName":"Lynette",
		"lastName":"Oyola",
		"dateOfBirth":"1988-10-24",
		"ssn":"123-45-6789"
	}`))
	req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)

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
	cfg := &core.Config{}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	app.Post("/api/v0/education-enrollments", EducationHandler(cfg, &fakeEducationService{
		response: education.Response{EnrollmentStatus: education.EnrollmentStatusFullTime},
	}, &fakeReporter{}, logger))

	req := httptest.NewRequest(http.MethodPost, "/api/v0/education-enrollments", strings.NewReader(`{
		"firstName":"Lynette",
		"lastName":"Oyola",
		"dateOfBirth":"1988-10-24",
		"ssn":"123-45-6789"
	}`))
	req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)

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
