package main

import (
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
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
	handler := reporting.NewLambdaHandler(logger)

	lambda.Start(handler.HandleRequest)
}
