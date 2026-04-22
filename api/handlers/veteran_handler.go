package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/cmsgov/emmy-api/pkg/circuitbreaker"
	"github.com/cmsgov/emmy-api/pkg/core"
	"github.com/cmsgov/emmy-api/pkg/education"
	"github.com/cmsgov/emmy-api/pkg/reporting"
	"github.com/cmsgov/emmy-api/pkg/resilience"
	"github.com/cmsgov/emmy-api/pkg/veteran"
	"github.com/gofiber/fiber/v2"
)

func VeteranDisabilityHandler(cfg *core.Config, service veteran.Service, reporter reporting.Reporter, logger *slog.Logger) fiber.Handler {
	const contextTimeout = 5 * time.Second

	if logger == nil {
		logger = slog.Default()
	}

	logger = logger.With(slog.String("handler", "VeteranDisabilityHandler"))

	return func(c *fiber.Ctx) error {
		requestStartTime := time.Now()
		clientID, ok := c.Locals("sub").(string)
		if !ok {
			clientID = ""
		}

		ctx := c.UserContext()

		var req veteran.Request
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
		}

		if !isValidVeteranRequest(&req) {
			return fiber.NewError(fiber.StatusBadRequest, "request must include first name, last name, date of birth, and either SSN or a complete address")
		}

		ctx, cancel := context.WithTimeout(ctx, contextTimeout)
		defer cancel()

		datasourceStartTime := time.Now()
		result, err := service.LookupDisabilityRating(ctx, req)
		datasourceDuration := time.Since(datasourceStartTime)
		if err != nil {
			var statusCode int
			switch {
			case errors.Is(err, veteran.ErrNotFound):
				statusCode = fiber.StatusNotFound
			case errors.Is(err, resilience.ErrCircuitOpen), errors.Is(err, circuitbreaker.ErrCircuitOpen):
				statusCode = fiber.StatusServiceUnavailable
			default:
				statusCode = fiber.StatusBadGateway
				logger.ErrorContext(ctx, "veteran disability lookup failed",
					slog.Any("error", err),
					slog.String("sub", clientID),
				)
			}

			reporter.Report(c.Context(), &reporting.ReportData{
				Endpoint:   c.Path(),
				Success:    false,
				DataSource: "VA",
				ClientID:   clientID,
				Timestamp:  time.Now(),
				StatusCode: statusCode,
			})

			if statusCode == fiber.StatusNotFound {
				return c.SendStatus(fiber.StatusNotFound)
			}
			return fiber.NewError(statusCode, http.StatusText(statusCode))
		}

		responseTimestamp := time.Now()
		rid, ok := ctx.Value(core.RequestContextKey).(string)
		if !ok || rid == "" {
			rid = "unknown"
		}
		result.Metadata = education.Metadata{
			APIVersion:               cfg.ServiceVersion,
			Environment:              cfg.Environment,
			RequestTimestamp:         requestStartTime.UTC().Format("2006-01-02T15:04:05.000Z"),
			ResponseTimestamp:        responseTimestamp.UTC().Format("2006-01-02T15:04:05.000Z"),
			DatasourceDurationMillis: datasourceDuration.Milliseconds(),
			TransactionID:            rid,
		}

		reporter.Report(c.Context(), &reporting.ReportData{
			Endpoint:   c.Path(),
			Success:    true,
			DataSource: "VA",
			ClientID:   clientID,
			Timestamp:  time.Now(),
			StatusCode: fiber.StatusOK,
		})

		return c.Status(fiber.StatusOK).JSON(result)
	}
}

func isValidVeteranRequest(req *veteran.Request) bool {
	if strings.TrimSpace(req.FirstName) == "" ||
		strings.TrimSpace(req.LastName) == "" ||
		strings.TrimSpace(req.DateOfBirth) == "" {
		return false
	}

	if strings.TrimSpace(req.SSN) != "" {
		return true
	}

	if req.Address == nil {
		return false
	}

	return strings.TrimSpace(req.Address.Street1) != "" &&
		strings.TrimSpace(req.Address.City) != "" &&
		strings.TrimSpace(req.Address.State) != "" &&
		strings.TrimSpace(req.Address.PostalCode) != "" &&
		strings.TrimSpace(req.Address.Country) != ""
}
