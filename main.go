package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DSACMS/verification-service-api/api"
	"github.com/DSACMS/verification-service-api/api/routes"
	"github.com/DSACMS/verification-service-api/pkg/core"
	"github.com/DSACMS/verification-service-api/pkg/redis"

	"github.com/gofiber/fiber/v2"
)

var ErrRunFailed = errors.New("application failed to run")

func main() {
	err := run()
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}

func run() error {
	logger := core.NewLogger(nil)

	err := core.LoadEnv()
	if err != nil {
		logger.Error(
			"Failed to load environment",
			"err", err,
		)
		return ErrRunFailed
	}

	cfg, err := core.NewConfigFromEnv()
	if err != nil {
		logger.Error(
			"Failed to get configuration",
			"err", err,
		)
		return ErrRunFailed
	}

	logger.Info("raw abc123 env", "SKIP_AUTH", os.Getenv("SKIP_AUTH"))

	logger.Info("config loaded",
		"env", cfg.Environment,
		"port", cfg.Port,
		"skip_auth", cfg.SkipAuth,
	)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	initLogger := core.NewLogger(&cfg)

	otel, err := core.NewOtelService(ctx, &cfg)
	if err != nil {
		initLogger.ErrorContext(
			ctx,
			"Otel error",
			"err", err,
		)
		return ErrRunFailed
	}
	defer otel.Shutdown(ctx, initLogger)

	logger = core.NewLoggerWithOtel(&cfg, otel)

	app, err := api.New(&api.Config{
		Core:   cfg,
		Logger: logger,
		Otel:   otel,
	})
	if err != nil {
		logger.ErrorContext(
			ctx,
			"Error building app",
			"err", err,
		)
		return ErrRunFailed
	}

	rdb := redis.NewClient(redis.Config{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	}, logger)

	err = redis.Ping(ctx, rdb)
	if err != nil {
		logger.ErrorContext(ctx, "redis ping failed", "err", err)
		_ = rdb.Close()
		return ErrRunFailed
	}

	defer func() {
		if err := rdb.Close(); err != nil {
			logger.Warn("redis close failed", "err", err)
		}
	}()

	routes.RegisterRoutes(app, &cfg, rdb, logger)

	addr, err := listenAddr(cfg.Port)
	if err != nil {
		logger.ErrorContext(
			ctx,
			"Invalid server port",
			"port", cfg.Port,
			"err", err,
		)
		return ErrRunFailed
	}

	err = runServer(ctx, app, addr)
	if err != nil {
		logger.ErrorContext(
			ctx,
			"Server error",
			"err", err,
		)
		return ErrRunFailed
	}

	return nil
}

func listenAddr(port int) (string, error) {
	if port < 1 || port > 65535 {
		return "", fmt.Errorf("port out of range: %d", port)
	}

	return fmt.Sprintf(":%d", port), nil
}

func runServer(ctx context.Context, app *fiber.App, addr string) error {
	srvErr := make(chan error, 1)

	go func() {
		srvErr <- app.Listen(addr)
	}()

	select {
	case err := <-srvErr:
		return err
	case <-ctx.Done():
	}

	err := app.ShutdownWithTimeout(5 * time.Second)
	if err != nil {
		return fmt.Errorf("error during shutdown: %w", err)
	}

	return nil
}
