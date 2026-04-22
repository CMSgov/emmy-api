package routes

import (
	"log/slog"

	"github.com/cmsgov/emmy-api/api/handlers"
	"github.com/cmsgov/emmy-api/api/middleware"
	"github.com/cmsgov/emmy-api/pkg/circuitbreaker"
	"github.com/cmsgov/emmy-api/pkg/core"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func StatusRouter(app fiber.Router, cfg core.Config, rdb *redis.Client, logger *slog.Logger) {
	if logger == nil {
		logger = slog.Default()
	}

	withBreaker := middleware.WithCircuitBreaker(func(name string) circuitbreaker.Breaker {
		return circuitbreaker.NewRedisBreaker(rdb, name, circuitbreaker.DefaultOptions(), logger)
	})

	app.Get("/health", withBreaker(handlers.GetRDBStatus(rdb)))
}
