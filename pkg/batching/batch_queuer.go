package batching

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/cmsgov/emmy-api/pkg/core"
)

type BatchQueuer interface {
	Batch(ctx context.Context, data *BatchData)
}

type batcher struct {
	logger    *slog.Logger
	sqsClient *sqs.Client
	cfg       *core.NSCConfig
}

func NewBatchQueuer(ctx context.Context, cfg *core.NSCConfig, logger *slog.Logger) BatchQueuer {
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

	return &batcher{
		cfg:       cfg,
		logger:    logger,
		sqsClient: sqsClient,
	}
}

func (r *batcher) Batch(ctx context.Context, data *BatchData) {
	if r.cfg.SQSQueueURL != "" && r.sqsClient != nil {
		body, err := json.Marshal(data)
		if err != nil {
			r.logger.ErrorContext(ctx, "failed to marshal batch data", "error", err)
			return
		}

		_, err = r.sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
			QueueUrl:    aws.String(r.cfg.SQSQueueURL),
			MessageBody: aws.String(string(body)),
		})
		if err != nil {
			r.logger.ErrorContext(ctx, "failed to send batch to SQS", "error", err)
			return
		}
		r.logger.DebugContext(ctx, "batch sent to SQS", "queue_url", r.cfg.SQSQueueURL)
	} else {
		r.logger.InfoContext(ctx, "batch queued",
			slog.String("batch", data.BatchID),
		)
	}
}
