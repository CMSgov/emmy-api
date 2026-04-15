package handlers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/cmsgov/emmy-api/pkg/circuitbreaker"
	"github.com/cmsgov/emmy-api/pkg/reporting"
	"github.com/cmsgov/emmy-api/pkg/resilience"
	"github.com/cmsgov/emmy-api/pkg/veteran"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func VeteranDisabilityHandler(service veteran.Service, reporter reporting.Reporter, logger *slog.Logger) fiber.Handler {
	const contextTimeout = 5 * time.Second

	if logger == nil {
		logger = slog.Default()
	}

	logger = logger.With(slog.String("handler", "VeteranDisabilityHandler"))

	return func(c *fiber.Ctx) error {
		clientID, ok := c.Locals("sub").(string)
		if !ok {
			clientID = ""
		}

		ctx, verificationSpan := verificationTracer.Start(
			c.UserContext(),
			"verification.request",
		)
		defer verificationSpan.End()

		verificationSpan.SetAttributes(
			attribute.String("vendor.name", "va"),
			attribute.String("http.route", c.Path()),
			attribute.String("http.method", c.Method()),
			attribute.String("user.id", clientID),
		)

		var req veteran.Request
		if err := c.BodyParser(&req); err != nil {
			verificationSpan.RecordError(err)
			verificationSpan.SetStatus(codes.Error, http.StatusText(fiber.StatusBadRequest))
			return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
		}

		if missing := missingVeteranIdentityField(req); missing != "" {
			err := fmt.Errorf("missing required field %s", missing)
			verificationSpan.RecordError(err)
			verificationSpan.SetStatus(codes.Error, http.StatusText(fiber.StatusBadRequest))
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if !hasVeteranLookupInput(req) {
			verificationSpan.SetStatus(codes.Error, http.StatusText(fiber.StatusNotFound))
			return c.SendStatus(fiber.StatusNotFound)
		}

		ctx, cancel := context.WithTimeout(ctx, contextTimeout)
		defer cancel()

		ctx, decisionSpan := verificationTracer.Start(ctx, "decision.engine")

		result, err := service.LookupDisabilityRating(ctx, req)
		if err != nil {
			decisionSpan.RecordError(err)
			decisionSpan.SetStatus(codes.Error, "verification failed")
			decisionSpan.End()

			var statusCode int
			switch {
			case errors.Is(err, veteran.ErrNotFound):
				statusCode = fiber.StatusNotFound
				verificationSpan.SetStatus(codes.Error, http.StatusText(fiber.StatusNotFound))
			case errors.Is(err, resilience.ErrCircuitOpen), errors.Is(err, circuitbreaker.ErrCircuitOpen):
				statusCode = fiber.StatusServiceUnavailable
				verificationSpan.SetStatus(codes.Error, http.StatusText(fiber.StatusServiceUnavailable))
			default:
				statusCode = fiber.StatusBadGateway
				logger.ErrorContext(ctx, "veteran disability lookup failed",
					slog.Any("error", err),
					slog.String("sub", clientID),
				)
				verificationSpan.SetStatus(codes.Error, http.StatusText(fiber.StatusBadGateway))
			}

			reporter.Report(c.Context(), reporting.ReportData{
				Endpoint:   c.Path(),
				Success:    false,
				DataSource: "VA",
				ClientID:   clientID,
				Timestamp:  time.Now(),
				StatusCode: statusCode,
			})

			verificationSpan.RecordError(err)

			if statusCode == fiber.StatusNotFound {
				return c.SendStatus(fiber.StatusNotFound)
			}
			return fiber.NewError(statusCode, http.StatusText(statusCode))
		}

		decisionSpan.SetStatus(codes.Ok, "decision completed")
		decisionSpan.End()

		reporter.Report(c.Context(), reporting.ReportData{
			Endpoint:   c.Path(),
			Success:    true,
			DataSource: "VA",
			ClientID:   clientID,
			Timestamp:  time.Now(),
			StatusCode: fiber.StatusOK,
		})

		verificationSpan.SetStatus(codes.Ok, "verification completed")
		return c.Status(fiber.StatusOK).JSON(result)
	}
}

func missingVeteranIdentityField(req veteran.Request) string {
	switch {
	case strings.TrimSpace(req.FirstName) == "":
		return "firstName"
	case strings.TrimSpace(req.LastName) == "":
		return "lastName"
	case strings.TrimSpace(req.DateOfBirth) == "":
		return "dateOfBirth"
	default:
		return ""
	}
}

func hasVeteranLookupInput(req veteran.Request) bool {
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
