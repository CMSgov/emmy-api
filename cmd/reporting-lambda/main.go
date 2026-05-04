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
	"github.com/cmsgov/emmy-api/pkg/core"
	"github.com/cmsgov/emmy-api/pkg/reporting"
	_ "github.com/lib/pq"
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
				}
			}
		}

		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.Database.Host, cfg.Database.Port, cfg.Database.User, password, cfg.Database.Name, cfg.Database.SSLMode)
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			logger.Error("failed to open database connection", "error", err)
		} else {
			err = db.Ping()
			if err != nil {
				logger.Error("failed to ping database", "error", err)
				err = db.Close()
				if err != nil {
					logger.Error("failed to close database", "error", err)
				}
				db = nil
			}
		}
	}

	handler := reporting.NewLambdaHandler(logger, db)

	lambda.Start(handler.HandleRequest)
}
