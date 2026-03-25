package routes

import (
	"log/slog"

	"github.com/DSACMS/verification-service-api/api/handlers"
	"github.com/DSACMS/verification-service-api/api/middleware"
	"github.com/DSACMS/verification-service-api/pkg/circuitbreaker"
	"github.com/DSACMS/verification-service-api/pkg/core"
	"github.com/DSACMS/verification-service-api/pkg/education"
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

	// One breaker per endpoint
	withCB := middleware.WithCircuitBreaker(func(name string) *circuitbreaker.RedisBreaker {
		return circuitbreaker.NewRedisBreaker(
			rdb,
			name,
			circuitbreaker.DefaultOptions(),
			logger,
		)
	})

	api.Get("/edu", withCB(handlers.EducationHandler(cfg, edu, logger)))
}
