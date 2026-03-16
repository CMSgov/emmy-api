package handlers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/DSACMS/verification-service-api/pkg/core"
	"github.com/DSACMS/verification-service-api/pkg/education"
	"github.com/DSACMS/verification-service-api/pkg/resilience"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

var verificationTracer = otel.Tracer("verification-service-api/verification")

func EducationHandler(cfg *core.Config, edu education.EducationService, logger *slog.Logger) fiber.Handler {
	const contextTimeout time.Duration = 5 * time.Second

	if logger == nil {
		logger = slog.Default()
	}
	logger = logger.With(slog.String("handler", "TestEducationHandler"))

	return func(c *fiber.Ctx) error {
		ctx, verificationSpan := verificationTracer.Start(
			c.UserContext(),
			"verification.request",
		)
		defer verificationSpan.End()

		verificationSpan.SetAttributes(
			attribute.String("vendor.name", "nsc"),
			attribute.String("http.route", c.Path()),
			attribute.String("http.method", c.Method()),
		)

		ctx, cancel := context.WithTimeout(ctx, contextTimeout)
		defer cancel()

		ctx, decisionSpan := verificationTracer.Start(ctx, "decision.engine")

		reqBody := education.Request{
			AccountID:        cfg.NSC.AccountID,
			OrganizationName: "Lynette",
			DateOfBirth:      "1988-10-24",
			LastName:         "Oyola",
			FirstName:        "Lynette",
			Terms:            "y",
			EndClient:        "CMS",
		}

		result, err := edu.Submit(ctx, reqBody)
		if err != nil {
			decisionSpan.RecordError(err)
			decisionSpan.SetStatus(codes.Error, "verification failed")
			decisionSpan.End()

			logger.ErrorContext(ctx, "education verification failed", slog.Any("error", err))

			verificationSpan.RecordError(err)

			status := fiber.StatusBadGateway
			if errors.Is(err, resilience.ErrCircuitOpen) {
				status = fiber.StatusServiceUnavailable
			}
			verificationSpan.SetStatus(codes.Error, http.StatusText(status))

			return fiber.NewError(
				status,
				fmt.Sprintf("education verification failed: %v", err),
			)
		}

		decisionSpan.SetStatus(codes.Ok, "decision completed")
		decisionSpan.End()

		verificationSpan.SetStatus(codes.Ok, "verification completed")
		return c.Status(fiber.StatusOK).JSON(result)
	}
}
