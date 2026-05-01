package routes

import (
	"database/sql"
	"log/slog"

	"github.com/cmsgov/emmy-api/api/handlers"
	"github.com/cmsgov/emmy-api/api/middleware"
	"github.com/cmsgov/emmy-api/pkg/circuitbreaker"
	"github.com/cmsgov/emmy-api/pkg/core"
	"github.com/cmsgov/emmy-api/pkg/education"
	"github.com/cmsgov/emmy-api/pkg/encryption"
	"github.com/cmsgov/emmy-api/pkg/reporting"
	"github.com/cmsgov/emmy-api/pkg/veteran"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type RouterParams struct {
	CFG        *core.Config
	RDB        *redis.Client
	Reporter   reporting.Reporter
	Logger     *slog.Logger
	Encryption encryption.Service
	DB         *sql.DB
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

	edu := education.New(&params.CFG.NSC, education.Options{
		Logger:     params.Logger,
		DB:         params.DB,
		Encryption: params.Encryption,
	})
	veteranService := veteran.New(&params.CFG.VA, veteran.Options{
		Logger: params.Logger,
	})

	// One breaker per endpoint
	withCB := middleware.WithCircuitBreaker(func(name string) circuitbreaker.Breaker {
		return circuitbreaker.NewRedisBreaker(
			params.RDB,
			name,
			circuitbreaker.DefaultOptions(),
			params.Logger,
		)
	})

	api.Post("/v0/education-enrollments", withCB(handlers.EducationHandler(cfg, edu, reporter, logger)))
	api.Post("/v0/batch-education-enrollments", withCB(handlers.BatchEducationHandler(cfg, edu, reporter, logger)))
	api.Get("/v0/batch-education-enrollments/:batchJobId", withCB(handlers.GetBatchStatusHandler(edu, reporter, logger)))
	api.Get("/v0/batch-education-enrollments/:batchJobId/details", withCB(handlers.GetBatchDetailsHandler(edu, reporter, logger)))
	api.Post("/v0/veteran-disability-ratings", withCB(handlers.VeteranDisabilityHandler(cfg, veteranService, reporter, logger)))
}
