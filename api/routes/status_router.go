package routes

import (
	"log/slog"

	"github.com/DSACMS/verification-service-api/api/handlers"
	"github.com/DSACMS/verification-service-api/api/middleware"
	"github.com/DSACMS/verification-service-api/pkg/circuitbreaker"
	"github.com/DSACMS/verification-service-api/pkg/core"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func StatusRouter(app fiber.Router, cfg core.Config, rdb *redis.Client, logger *slog.Logger) {
	if logger == nil {
		logger = slog.Default()
	}

	withBreaker := middleware.WithCircuitBreaker(func(name string) *circuitbreaker.RedisBreaker {
		return circuitbreaker.NewRedisBreaker(rdb, name, circuitbreaker.DefaultOptions(), logger)
	})

	app.Get("/health", withBreaker(handlers.GetRDBStatus(rdb)))
}
