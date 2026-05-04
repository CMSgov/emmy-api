package routes

import (
	"log/slog"

	"github.com/cmsgov/emmy-api/api/handlers"
	"github.com/cmsgov/emmy-api/pkg/core"
	"github.com/cmsgov/emmy-api/pkg/education"
	"github.com/cmsgov/emmy-api/pkg/reporting"
	"github.com/cmsgov/emmy-api/pkg/veteran"
	"github.com/gofiber/fiber/v2"
)

type RouterParams struct {
	CFG      *core.Config
	Reporter reporting.Reporter
	Logger   *slog.Logger
	EDU      education.Service
	VA       veteran.Service
	WithCB   func(fiber.Handler) fiber.Handler
}

func RegisterRoutes(app fiber.Router, params RouterParams) {
	if params.Logger == nil {
		params.Logger = slog.Default()
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Backend running!")
	})
	app.Get("/api-spec/v1/verify", handlers.OpenAPISpecHandler())

	api := app.Group("/api")

	api.Post("/v0/education-enrollments", params.WithCB(handlers.EducationHandler(params.CFG, params.EDU, params.Reporter, params.Logger)))

	api.Post("/v0/batch-education-enrollments", params.WithCB(handlers.BatchEducationHandler(params.CFG, params.EDU, params.Reporter, params.Logger)))

	api.Get("/v0/batch-education-enrollments/:batchJobId", handlers.GetBatchStatusHandler(params.EDU, params.Reporter, params.Logger))

	api.Get("/v0/batch-education-enrollments/:batchJobId/details", handlers.GetBatchDetailsHandler(params.EDU, params.Reporter, params.Logger))

	api.Post("/v0/veteran-disability-ratings", params.WithCB(handlers.VeteranDisabilityHandler(params.CFG, params.VA, params.Reporter, params.Logger)))
}
