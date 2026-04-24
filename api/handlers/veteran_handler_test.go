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

	"github.com/cmsgov/emmy-api/pkg/core"
	"github.com/cmsgov/emmy-api/pkg/resilience"
	"github.com/cmsgov/emmy-api/pkg/veteran"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

type fakeVeteranService struct {
	response veteran.Response
	err      error
	calls    int
	lastReq  veteran.Request
}

func (s *fakeVeteranService) LookupDisabilityRating(_ context.Context, req veteran.Request) (veteran.Response, error) {
	s.calls++
	s.lastReq = req
	return s.response, s.err
}

func TestVeteranDisabilityHandler_Success(t *testing.T) {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		rid := c.Get(fiber.HeaderXRequestID)
		ctx := context.WithValue(c.UserContext(), core.RequestContextKey, rid)
		c.SetUserContext(ctx)
		return c.Next()
	})
	cfg := &core.Config{Environment: "test", ServiceVersion: "1.3.0"}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	service := &fakeVeteranService{
		response: veteran.Response{CombinedDisabilityRating: 70},
	}
	reporter := &fakeReporter{}

	app.Post("/api/v0/veteran-disability-ratings", VeteranDisabilityHandler(cfg, service, reporter, logger))

	req := httptest.NewRequest(http.MethodPost, "/api/v0/veteran-disability-ratings", strings.NewReader(`{
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
	require.Contains(t, string(body), `"combinedDisabilityRating":70`)
	require.Contains(t, string(body), `"metadata"`)
	require.Contains(t, string(body), `"apiVersion":"1.3.0"`)
	require.Contains(t, string(body), `"environment":"test"`)
	require.Contains(t, string(body), `"transaction-id":"test-request-id"`)
	require.Equal(t, 1, service.calls)
	require.Equal(t, "123-45-6789", service.lastReq.SSN)

	require.Len(t, reporter.calls, 1)
	require.True(t, reporter.calls[0].Success)
	require.Equal(t, "VA", reporter.calls[0].DataSource)
	require.Equal(t, fiber.StatusOK, reporter.calls[0].StatusCode)
}

func TestVeteranDisabilityHandler_AddressOnlySuccess(t *testing.T) {
	app := fiber.New()
	cfg := &core.Config{}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	service := &fakeVeteranService{
		response: veteran.Response{CombinedDisabilityRating: 70},
	}
	reporter := &fakeReporter{}

	app.Post("/api/v0/veteran-disability-ratings", VeteranDisabilityHandler(cfg, service, reporter, logger))

	req := httptest.NewRequest(http.MethodPost, "/api/v0/veteran-disability-ratings", strings.NewReader(`{
		"firstName":"Lynette",
		"lastName":"Oyola",
		"dateOfBirth":"1988-10-24",
		"address":{
			"street1":"17020 Tortoise St",
			"city":"Round Rock",
			"state":"TX",
			"postalCode":"78664",
			"country":"USA"
		}
	}`))
	req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, fiber.StatusOK, resp.StatusCode)
	require.Equal(t, 1, service.calls)
	require.NotNil(t, service.lastReq.Address)
	require.Equal(t, "17020 Tortoise St", service.lastReq.Address.Street1)
}

func TestVeteranDisabilityHandler_InvalidJSONReturnsBadRequest(t *testing.T) {
	app := fiber.New()
	cfg := &core.Config{}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	app.Post("/api/v0/veteran-disability-ratings", VeteranDisabilityHandler(cfg, &fakeVeteranService{}, &fakeReporter{}, logger))

	req := httptest.NewRequest(http.MethodPost, "/api/v0/veteran-disability-ratings", strings.NewReader(`{`))
	req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestVeteranDisabilityHandler_MissingRequiredFieldReturnsBadRequest(t *testing.T) {
	app := fiber.New()
	cfg := &core.Config{}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	app.Post("/api/v0/veteran-disability-ratings", VeteranDisabilityHandler(cfg, &fakeVeteranService{}, &fakeReporter{}, logger))

	req := httptest.NewRequest(http.MethodPost, "/api/v0/veteran-disability-ratings", strings.NewReader(`{
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

func TestVeteranDisabilityHandler_MissingSSNAndAddressReturnsBadRequestWithoutCallingProvider(t *testing.T) {
	app := fiber.New()
	cfg := &core.Config{}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	service := &fakeVeteranService{}

	app.Post("/api/v0/veteran-disability-ratings", VeteranDisabilityHandler(cfg, service, &fakeReporter{}, logger))

	req := httptest.NewRequest(http.MethodPost, "/api/v0/veteran-disability-ratings", strings.NewReader(`{
		"firstName":"Lynette",
		"lastName":"Oyola",
		"dateOfBirth":"1988-10-24"
	}`))
	req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	require.Equal(t, 0, service.calls)
}

func TestVeteranDisabilityHandler_IncompleteAddressWithoutSSNReturnsBadRequestWithoutCallingProvider(t *testing.T) {
	app := fiber.New()
	cfg := &core.Config{}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	service := &fakeVeteranService{}

	app.Post("/api/v0/veteran-disability-ratings", VeteranDisabilityHandler(cfg, service, &fakeReporter{}, logger))

	req := httptest.NewRequest(http.MethodPost, "/api/v0/veteran-disability-ratings", strings.NewReader(`{
		"firstName":"Lynette",
		"lastName":"Oyola",
		"dateOfBirth":"1988-10-24",
		"address":{
			"street1":"17020 Tortoise St",
			"city":"Round Rock"
		}
	}`))
	req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	require.Equal(t, 0, service.calls)
}

func TestVeteranDisabilityHandler_UpstreamNotFoundReturnsNotFound(t *testing.T) {
	app := fiber.New()
	cfg := &core.Config{}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	app.Post("/api/v0/veteran-disability-ratings", VeteranDisabilityHandler(cfg, &fakeVeteranService{
		err: veteran.ErrNotFound,
	}, &fakeReporter{}, logger))

	req := httptest.NewRequest(http.MethodPost, "/api/v0/veteran-disability-ratings", strings.NewReader(`{
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
}

func TestVeteranDisabilityHandler_UpstreamErrorReturnsBadGateway(t *testing.T) {
	app := fiber.New()
	cfg := &core.Config{}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	app.Post("/api/v0/veteran-disability-ratings", VeteranDisabilityHandler(cfg, &fakeVeteranService{
		err: errors.New("provider failed"),
	}, &fakeReporter{}, logger))

	req := httptest.NewRequest(http.MethodPost, "/api/v0/veteran-disability-ratings", strings.NewReader(`{
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

func TestVeteranDisabilityHandler_CircuitOpenReturnsServiceUnavailable(t *testing.T) {
	app := fiber.New()
	cfg := &core.Config{}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	app.Post("/api/v0/veteran-disability-ratings", VeteranDisabilityHandler(cfg, &fakeVeteranService{
		err: resilience.ErrCircuitOpen,
	}, &fakeReporter{}, logger))

	req := httptest.NewRequest(http.MethodPost, "/api/v0/veteran-disability-ratings", strings.NewReader(`{
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
