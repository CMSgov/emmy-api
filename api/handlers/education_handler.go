package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/DSACMS/verification-service-api/pkg/core"
	"github.com/DSACMS/verification-service-api/pkg/education"
	"github.com/gofiber/fiber/v2"
)

func EducationHandler(cfg *core.Config, edu education.EducationService, logger *slog.Logger) fiber.Handler {
	const contextTimeout time.Duration = 5 * time.Second

	if logger == nil {
		logger = slog.Default()
	}
	logger = logger.With(slog.String("handler", "TestEducationHandler"))

	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.Context(), contextTimeout)
		defer cancel()

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
			logger.ErrorContext(ctx, "education verification failed", slog.Any("error", err))
			return fiber.NewError(
				fiber.StatusBadGateway,
				fmt.Sprintf("education verification failed: %v", err),
			)
		}

		return c.Status(fiber.StatusOK).JSON(result)
	}
}
