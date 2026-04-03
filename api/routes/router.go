package routes

import (
	"log/slog"

	"github.com/cmsgov/emmy-api/api/handlers"
	"github.com/cmsgov/emmy-api/api/middleware"
	"github.com/cmsgov/emmy-api/pkg/circuitbreaker"
	"github.com/cmsgov/emmy-api/pkg/core"
	"github.com/cmsgov/emmy-api/pkg/education"
	"github.com/cmsgov/emmy-api/pkg/veteran"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func RegisterRoutes(app fiber.Router, cfg *core.Config, rdb *redis.Client, logger *slog.Logger) {
	if logger == nil {
		logger = slog.Default()
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Backend running!")
	})
	app.Get("/api-spec/v1/verify", handlers.OpenAPISpecHandler())

	api := app.Group("/api")

	edu := education.New(&cfg.NSC, education.Options{
		Logger: logger,
	})
	veteranService := veteran.New(&cfg.VA, veteran.Options{
		Logger: logger,
	})

	// One breaker per endpoint
	withCB := middleware.WithCircuitBreaker(func(name string) *circuitbreaker.RedisBreaker {
		return circuitbreaker.NewRedisBreaker(
			rdb,
			name,
			circuitbreaker.DefaultOptions(),
			logger,
		)
	})

	api.Post("/v0/education-enrollments", withCB(handlers.EducationHandler(edu, logger)))
	api.Post("/v0/veteran-disability-ratings", withCB(handlers.VeteranDisabilityHandler(veteranService, logger)))
}
