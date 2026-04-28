package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/cmsgov/emmy-api/api"
	"github.com/cmsgov/emmy-api/pkg/core"
	"github.com/cmsgov/emmy-api/pkg/encryption"
	"github.com/cmsgov/emmy-api/pkg/redis"
	"github.com/cmsgov/emmy-api/pkg/reporting"
	_ "github.com/lib/pq"

	"github.com/gofiber/fiber/v2"
)

var (
	ErrRunFailed      = errors.New("application failed to run")
	ErrPortOutOfRange = errors.New("port out of range")
)

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

	rdb := redis.NewClient(redis.Config{
		Addr:               cfg.Redis.Addr,
		Password:           cfg.Redis.Password,
		DB:                 cfg.Redis.DB,
		UseTLS:             cfg.Redis.UseTLS,
		InsecureSkipVerify: cfg.Redis.InsecureSkipVerify,
	}, logger)

	err = redis.Ping(ctx, rdb)
	if err != nil {
		logger.ErrorContext(ctx, "redis ping failed", "err", err)
		if closeErr := rdb.Close(); closeErr != nil {
			logger.Warn("redis close failed after ping error", "err", closeErr)
		}
		return ErrRunFailed
	}

	defer func() {
		if closeErr := rdb.Close(); closeErr != nil {
			logger.Warn("redis close failed", "err", closeErr)
		}
	}()

	db, err := initDatabase(ctx, &cfg, logger)
	if err != nil {
		logger.ErrorContext(ctx, "database initialization failed", "err", err)
		return ErrRunFailed
	}
	defer func() {
		if closeError := db.Close(); closeError != nil {
			logger.Error("failed to close db", "error", closeError)
		}
	}()

	awsCfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "failed to load AWS config", "err", err)
		return ErrRunFailed
	}

	var encryptionService encryption.Service
	if cfg.KMS.KeyID != "" {
		kmsClient := kms.NewFromConfig(awsCfg)
		encryptionService = encryption.NewKMSService(kmsClient, cfg.KMS.KeyID)
		logger.Info("KMS encryption service initialized", "key_id", cfg.KMS.KeyID)
	} else {
		if cfg.Environment != "development" {
			logger.Error("KMS_KEY_ID not set", "env", cfg.Environment)
			return ErrRunFailed
		}

		encryptionService = encryption.NewMockEncryptionService()
		logger.Warn("KMS_KEY_ID not set, using mock encryption service")
	}

	reporter := reporting.NewReporter(ctx, cfg.Reporting, logger)
	app, err := api.New(&api.Config{
		Core:       cfg,
		Logger:     initLogger,
		Redis:      rdb,
		DB:         db,
		Encryption: encryptionService,
		Reporter:   reporter,
	})
	if err != nil {
		logger.ErrorContext(
			ctx,
			"Error building app",
			"err", err,
		)
		return ErrRunFailed
	}

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

func initDatabase(ctx context.Context, cfg *core.Config, logger *slog.Logger) (*sql.DB, error) {
	password := cfg.Database.Password
	if cfg.Database.IAMAuth {
		var awsConfigErr error
		awsConfig, awsConfigErr := config.LoadDefaultConfig(ctx)
		if awsConfigErr != nil {
			logger.Error("failed to load AWS config for IAM auth", "error", awsConfigErr)
		} else {
			endpoint := net.JoinHostPort(cfg.Database.Host, cfg.Database.Port)
			token, tokenErr := auth.BuildAuthToken(ctx, endpoint, awsConfig.Region, cfg.Database.User, awsConfig.Credentials)
			if tokenErr != nil {
				logger.Error("failed to build auth token", "error", tokenErr)
			} else {
				password = url.QueryEscape(token)
			}
		}
	}

	dbHostPort := net.JoinHostPort(cfg.Database.Host, cfg.Database.Port)
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		cfg.Database.User, password, dbHostPort, cfg.Database.Name, cfg.Database.SSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		if closeErr := db.Close(); closeErr != nil {
			logger.Error("failed to close db", "error", closeErr)
		}
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func listenAddr(port int) (string, error) {
	if port < 1 || port > 65535 {
		return "", fmt.Errorf("%w: %d", ErrPortOutOfRange, port)
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
