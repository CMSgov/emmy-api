package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
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
		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Name)
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
