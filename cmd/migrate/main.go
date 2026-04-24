package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"net/url"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
	"github.com/cmsgov/emmy-api/pkg/core"
)

func main() {
	var migrationPath string
	flag.StringVar(&migrationPath, "path", "./migrations", "path to migration files")
	flag.Parse()

	command := flag.Arg(0)
	if command == "" {
		fmt.Fprintln(os.Stderr, "Usage: migrate [options] <command>")           //nolint:errcheck // CLI usage info
		fmt.Fprintln(os.Stderr, "Commands: up, down, force <version>, version") //nolint:errcheck // CLI usage info
		os.Exit(1)
	}

	cfg, err := core.NewConfigFromEnv()
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	logger := core.NewLogger(&cfg)

	password := cfg.Database.Password
	if cfg.Database.IAMAuth {
		ctx := context.Background()
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


	m, err := migrate.New(
		"file://"+migrationPath,
		dsn,
	)
	if err != nil {
		logger.Error("failed to initialize migration", "error", err, "dsn", dsn)
		os.Exit(1)
	}

	switch command {
	case "up":
		err = m.Up()
	case "down":
		err = m.Down()
	case "version":
		currVersion, dirty, versionErr := m.Version()
		if versionErr != nil {
			logger.Error("failed to get version", "error", versionErr)
			os.Exit(1)
		}
		//nolint:forbidigo,errcheck // CLI output
		fmt.Printf("Version: %d, Dirty: %v\n", currVersion, dirty)
		return
	default:
		//nolint:forbidigo,errcheck // CLI output
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.Error("migration failed", "command", command, "error", err)
		os.Exit(1)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		logger.Info("no changes to apply")
	} else {
		logger.Info("migration completed successfully", "command", command)
	}
}
