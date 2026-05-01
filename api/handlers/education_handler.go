package handlers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/cmsgov/emmy-api/pkg/core"
	"github.com/cmsgov/emmy-api/pkg/education"
	"github.com/cmsgov/emmy-api/pkg/reporting"
	"github.com/cmsgov/emmy-api/pkg/resilience"
	"github.com/gofiber/fiber/v2"
)

var errMissingRequiredField = errors.New("missing required field")

func EducationHandler(cfg *core.Config, edu education.Service, reporter reporting.Reporter, logger *slog.Logger) fiber.Handler {
	const contextTimeout time.Duration = 30 * time.Second

	if logger == nil {
		logger = slog.Default()
	}
	logger = logger.With(slog.String("handler", "EducationHandler"))

	return func(c *fiber.Ctx) error {
		requestStartTime := time.Now()
		clientID, ok := c.Locals("sub").(string)
		if !ok {
			clientID = ""
		}

		ctx := c.UserContext()

		var reqBody education.Request
		if err := c.BodyParser(&reqBody); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
		}

		if missing := missingEducationIdentityField(&reqBody); missing != "" {
			err := fmt.Errorf("%w: %s", errMissingRequiredField, missing)
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		ctx, cancel := context.WithTimeout(ctx, contextTimeout)
		defer cancel()

		datasourceStartTime := time.Now()
		result, err := edu.LookupEnrollmentStatus(ctx, reqBody)
		datasourceDuration := time.Since(datasourceStartTime)
		if err != nil {
			var statusCode int
			switch {
			case errors.Is(err, education.ErrNotFound):
				statusCode = fiber.StatusNotFound
			case errors.Is(err, resilience.ErrCircuitOpen):
				statusCode = fiber.StatusServiceUnavailable
			default:
				statusCode = fiber.StatusBadGateway
			}

			reporter.Report(c.Context(), &reporting.ReportData{
				Endpoint:   c.Path(),
				Success:    false,
				DataSource: "NSC",
				ClientID:   clientID,
				Timestamp:  time.Now(),
				StatusCode: statusCode,
			})

			logger.ErrorContext(ctx, "education verification failed", slog.Any("error", err))

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
			DataSource: "NSC",
			ClientID:   clientID,
			Timestamp:  time.Now(),
			StatusCode: fiber.StatusOK,
		})

		return c.Status(fiber.StatusOK).JSON(result)
	}
}

func BatchEducationHandler(_ *core.Config, edu education.Service, reporter reporting.Reporter, logger *slog.Logger) fiber.Handler {
	if logger == nil {
		logger = slog.Default()
	}
	logger = logger.With(slog.String("handler", "BatchEducationHandler"))

	return func(c *fiber.Ctx) error {
		clientID, ok := c.Locals("sub").(string)
		if !ok {
			clientID = ""
		}

		var reqBody education.BatchRequest
		if err := c.BodyParser(&reqBody); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
		}

		if reqBody.BatchID == "" {
			return fiber.NewError(fiber.StatusBadRequest, "missing required field batchId")
		}

		err := edu.RegisterBatch(c.UserContext(), reqBody)
		if err != nil {
			logger.ErrorContext(c.UserContext(), "Failed to initiate batch. Please reach out to emmy support with the transaction id", slog.Any("error", err))
			reporter.Report(c.Context(), &reporting.ReportData{
				Endpoint:   c.Path(),
				Success:    false,
				DataSource: "NSC",
				ClientID:   clientID,
				Timestamp:  time.Now(),
				StatusCode: fiber.StatusInternalServerError,
			})
			return fiber.NewError(fiber.StatusInternalServerError, "failed to record batch")
		}

		reporter.Report(c.Context(), &reporting.ReportData{
			Endpoint:   c.Path(),
			Success:    true,
			DataSource: "NSC",
			ClientID:   clientID,
			Timestamp:  time.Now(),
			StatusCode: fiber.StatusAccepted,
		})

		return c.SendStatus(fiber.StatusAccepted)
	}
}

func GetBatchStatusHandler(edu education.Service, reporter reporting.Reporter, logger *slog.Logger) fiber.Handler {
	if logger == nil {
		logger = slog.Default()
	}
	logger = logger.With(slog.String("handler", "GetBatchStatusHandler"))

	return func(c *fiber.Ctx) error {
		clientID, ok := c.Locals("sub").(string)
		if !ok {
			clientID = ""
		}

		batchJobID := c.Params("batchJobId")
		if batchJobID == "" {
			return fiber.NewError(fiber.StatusBadRequest, "missing required parameter batchJobId")
		}

		result, err := edu.GetBatchStatus(c.UserContext(), batchJobID)
		if err != nil {
			if errors.Is(err, education.ErrNotFound) {
				return c.SendStatus(fiber.StatusNotFound)
			}

			logger.ErrorContext(c.UserContext(), "failed to get batch status", slog.String("batchJobId", batchJobID), slog.Any("error", err))
			reporter.Report(c.Context(), &reporting.ReportData{
				Endpoint:   c.Path(),
				Success:    false,
				DataSource: "Postgres",
				ClientID:   clientID,
				Timestamp:  time.Now(),
				StatusCode: fiber.StatusInternalServerError,
			})
			return fiber.NewError(fiber.StatusInternalServerError, "failed to get batch status")
		}

		reporter.Report(c.Context(), &reporting.ReportData{
			Endpoint:   c.Path(),
			Success:    true,
			DataSource: "Postgres",
			ClientID:   clientID,
			Timestamp:  time.Now(),
			StatusCode: fiber.StatusOK,
		})

		return c.Status(fiber.StatusOK).JSON(result)
	}
}

func missingEducationIdentityField(req *education.Request) string {
	if req == nil {
		return "firstName"
	}

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
