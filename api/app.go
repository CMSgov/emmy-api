package api

import (
	"context"
	"errors"
	"log/slog"
	"runtime/debug"

	"github.com/cmsgov/emmy-api/api/middleware"
	"github.com/cmsgov/emmy-api/api/routes"
	"github.com/cmsgov/emmy-api/pkg/core"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	slogfiber "github.com/samber/slog-fiber"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func RequestIDToUserContext() fiber.Handler {
	return func(c *fiber.Ctx) error {
		rid := c.Get(fiber.HeaderXRequestID)
		if rid == "" {
			rid = uuid.NewString()
		}

		ctx := context.WithValue(c.UserContext(), core.RequestContextKey, rid)
		c.SetUserContext(ctx)

		return c.Next()
	}
}

func errorHandler(logger *slog.Logger) fiber.ErrorHandler {
	handleFiberError := func(ctx *fiber.Ctx, err *fiber.Error) error {
		logger.ErrorContext(
			ctx.UserContext(),
			"fiber error",
			"code", err.Code,
			"message", err.Message,
		)

		return ctx.
			Status(err.Code).
			SendString(err.Message)
	}

	return func(ctx *fiber.Ctx, err error) error {
		var e *fiber.Error
		if !errors.As(err, &e) {
			e = fiber.ErrInternalServerError
		}
		return handleFiberError(ctx, e)
	}
}

func stackTraceHandler(logger *slog.Logger) func(*fiber.Ctx, any) {
	return func(c *fiber.Ctx, e any) {
		stack := debug.Stack()
		logger.ErrorContext(
			c.UserContext(),
			"panic!",
			"stack", stack,
			"err", e,
		)
	}
}

type Config struct {
	Logger *slog.Logger
	Redis  *redis.Client
	Core   core.Config
}

func New(cfg *Config) (*fiber.App, error) {
	if cfg.Logger == nil {
		cfg.Logger = core.NewLogger(&cfg.Core)
	}

	logger := cfg.Logger.With(slog.String("component", "api"))

	fiberConfig := fiber.Config{
		ErrorHandler: errorHandler(cfg.Logger),
	}

	app := fiber.New(fiberConfig)

	app.Use(RequestIDToUserContext())

	app.Use(slogfiber.NewWithConfig(
		cfg.Logger,
		slogfiber.Config{
			WithRequestID: true,
			WithSpanID:    true,
			WithTraceID:   true,
		},
	))

	app.Use(recover.New(recover.Config{
		EnableStackTrace:  true,
		StackTraceHandler: stackTraceHandler(cfg.Logger),
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
		AllowMethods: "*",
	}))

	app.Use(middleware.SubjectMiddleware(cfg.Logger))

	routes.StatusRouter(app, cfg.Redis, logger)

	if cfg.Core.SkipAuth {
		app.Use(middleware.SkipAuthMiddleware())
	}

	return app, nil
}
