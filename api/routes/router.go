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

func RegisterRoutes(app fiber.Router, cfg *core.Config, rdb *redis.Client, db *sql.DB, encrypt encryption.Service, reporter reporting.Reporter, logger *slog.Logger) {
	if logger == nil {
		logger = slog.Default()
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Backend running!")
	})
	app.Get("/api-spec/v1/verify", handlers.OpenAPISpecHandler())

	api := app.Group("/api")

	edu := education.New(&cfg.NSC, education.Options{
		Logger:     logger,
		DB:         db,
		Encryption: encrypt,
	})
	veteranService := veteran.New(&cfg.VA, veteran.Options{
		Logger: logger,
	})

	// One breaker per endpoint
	withCB := middleware.WithCircuitBreaker(func(name string) circuitbreaker.Breaker {
		return circuitbreaker.NewRedisBreaker(
			rdb,
			name,
			circuitbreaker.DefaultOptions(),
			logger,
		)
	})

	api.Post("/v0/education-enrollments", withCB(handlers.EducationHandler(cfg, edu, reporter, logger)))
	api.Post("/v0/batch-education-enrollments", withCB(handlers.BatchEducationHandler(cfg, edu, reporter, logger)))
	api.Get("/v0/batch-education-enrollments/:batchJobId", handlers.GetBatchStatusHandler(edu, reporter, logger))
	api.Post("/v0/veteran-disability-ratings", withCB(handlers.VeteranDisabilityHandler(cfg, veteranService, reporter, logger)))
}
