package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
	_ "github.com/lib/pq"
	"github.com/cmsgov/emmy-api/pkg/core"
	"github.com/cmsgov/emmy-api/pkg/reporting"
)

func main() {
	// Simple initialization for Lambda environment
	cfg, err := core.NewConfigFromEnv()
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	logger := core.NewLogger(&cfg)

	var db *sql.DB
	if cfg.Database.Host != "" {
		password := cfg.Database.Password
		iamAuthEnabled := cfg.Database.IAMAuth
		logger.Info("database connection info",
			"host", cfg.Database.Host,
			"port", cfg.Database.Port,
			"user", cfg.Database.User,
			"iam_auth_enabled", iamAuthEnabled,
			"has_password", password != "")

		if iamAuthEnabled {
			ctx := context.Background()
			var awsConfigErr error
			awsConfig, awsConfigErr := config.LoadDefaultConfig(ctx)
			if awsConfigErr != nil {
				logger.Error("failed to load AWS config for IAM auth", "error", awsConfigErr)
			} else {
				endpoint := fmt.Sprintf("%s:%s", cfg.Database.Host, cfg.Database.Port)
				token, tokenErr := auth.BuildAuthToken(ctx, endpoint, awsConfig.Region, cfg.Database.User, awsConfig.Credentials)
				if tokenErr != nil {
					logger.Error("failed to build auth token", "error", tokenErr)
				} else {
					password = token
					logger.Info("generated IAM auth token",
						"user", cfg.Database.User,
						"endpoint", endpoint,
						"region", awsConfig.Region)
				}
			}
		}

		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.Database.Host, cfg.Database.Port, cfg.Database.User, password, cfg.Database.Name, cfg.Database.SSLMode)
		db, err = sql.Open("postgres", dsn)
		logger.Info("connecting to database", "host", cfg.Database.Host, "dbname", cfg.Database.Name)
		if err != nil {
			logger.Error("failed to open database connection", "error", err)
		} else {
			err = db.Ping()
			if err != nil {
				logger.Error("failed to ping database", "error", err)
				db.Close()
				db = nil
			} else {
				logger.Info("connected to database", "host", cfg.Database.Host, "dbname", cfg.Database.Name)
				defer db.Close()
			}
		}
	}

	handler := reporting.NewLambdaHandler(logger, db)

	lambda.Start(handler.HandleRequest)
}
