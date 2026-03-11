package api

import (
	"errors"
	"fmt"
	"log/slog"
	"runtime/debug"

	"github.com/DSACMS/verification-service-api/api/middleware"
	"github.com/DSACMS/verification-service-api/api/routes"
	"github.com/DSACMS/verification-service-api/pkg/core"

	"go.opentelemetry.io/otel/codes"

	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	slogfiber "github.com/samber/slog-fiber"

	"github.com/redis/go-redis/v9"
)

func errorHandler(logger *slog.Logger, otel core.OtelService) fiber.ErrorHandler {
	handleFiberError := func(ctx *fiber.Ctx, err *fiber.Error) error {
		span := otel.SpanFromContext(ctx.Context())
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Message)

		logger.ErrorContext(
			ctx.Context(),
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
			c.Context(),
			"panic!",
			"stack", stack,
			"err", e,
		)
	}
}

type Config struct {
	Otel   core.OtelService
	Logger *slog.Logger
	Core   core.Config

	Redis *redis.Client
}

func New(cfg *Config) (*fiber.App, error) {
	if cfg.Logger == nil {
		cfg.Logger = core.NewLogger(&cfg.Core)
	}

	logger := cfg.Logger.With(slog.String("component", "api"))

	fiberConfig := fiber.Config{
		ErrorHandler: errorHandler(cfg.Logger, cfg.Otel),
	}

	app := fiber.New(fiberConfig)

	app.Use(recover.New(recover.Config{
		EnableStackTrace:  true,
		StackTraceHandler: stackTraceHandler(cfg.Logger),
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
		AllowMethods: "*",
	}))

	app.Use(otelfiber.Middleware())

	app.Use(slogfiber.NewWithConfig(
		cfg.Logger,
		slogfiber.Config{
			WithRequestID: true,
			WithSpanID:    true,
			WithTraceID:   true,
		},
	))

	if cfg.Core.SkipAuth {
		app.Use(middleware.SkipAuthMiddleware())
	} else {
		verifier, err := middleware.NewCognitoVerifier(middleware.CognitoConfig{
			Region:     cfg.Core.Cognito.Region,
			UserPoolID: cfg.Core.Cognito.UserPoolID,
			ClientID:   cfg.Core.Cognito.AppClientID,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to initialize cognito middleware: %w", err)
		}
		app.Use(verifier.FiberMiddleware())
	}

	routes.StatusRouter(app, cfg.Core, cfg.Redis, logger)

	return app, nil
}
