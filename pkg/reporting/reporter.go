package reporting

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/cmsgov/emmy-api/pkg/core"
)

type ReportData struct {
	Endpoint   string    `json:"endpoint"`
	DataSource string    `json:"data_source"`
	ClientID   string    `json:"client_id"`
	Success    bool      `json:"success"`
	Timestamp  time.Time `json:"timestamp"`
	StatusCode int       `json:"status_code"`
}

type Reporter interface {
	Report(ctx context.Context, data *ReportData)
}

type reporter struct {
	logger    *slog.Logger
	sqsClient *sqs.Client
	cfg       core.ReportingConfig
}

func NewReporter(ctx context.Context, cfg core.ReportingConfig, logger *slog.Logger) Reporter {
	if logger == nil {
		logger = slog.Default()
	}
	logger = logger.With(slog.String("component", "reporting"))

	var sqsClient *sqs.Client
	if cfg.SQSQueueURL != "" {
		awsCfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			logger.Error("failed to load AWS config", "error", err)
		} else {
			sqsClient = sqs.NewFromConfig(awsCfg)
		}
	}

	return &reporter{
		cfg:       cfg,
		logger:    logger,
		sqsClient: sqsClient,
	}
}

func (r *reporter) Report(ctx context.Context, data *ReportData) {
	if r.cfg.SQSQueueURL != "" && r.sqsClient != nil {
		body, err := json.Marshal(data)
		if err != nil {
			r.logger.ErrorContext(ctx, "failed to marshal report data", "error", err)
			return
		}

		_, err = r.sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
			QueueUrl:    aws.String(r.cfg.SQSQueueURL),
			MessageBody: aws.String(string(body)),
		})
		if err != nil {
			r.logger.ErrorContext(ctx, "failed to send report to SQS", "error", err)
			return
		}
		r.logger.DebugContext(ctx, "report sent to SQS", "queue_url", r.cfg.SQSQueueURL)
	} else {
		r.logger.InfoContext(ctx, "api call report",
			slog.String("endpoint", data.Endpoint),
			slog.Bool("success", data.Success),
			slog.String("data_source", data.DataSource),
			slog.String("client_id", data.ClientID),
			slog.Time("timestamp", data.Timestamp),
			slog.Int("status_code", data.StatusCode),
		)
	}
}
