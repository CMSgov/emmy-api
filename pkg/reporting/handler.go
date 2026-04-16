package reporting

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
)

// LambdaHandler handles SQS events and logs ReportData
type LambdaHandler struct {
	logger *slog.Logger
}

// NewLambdaHandler creates a new LambdaHandler
func NewLambdaHandler(logger *slog.Logger) *LambdaHandler {
	if logger == nil {
		logger = slog.Default()
	}
	return &LambdaHandler{
		logger: logger.With(slog.String("component", "reporting-lambda")),
	}
}

// HandleRequest processes SQS messages
func (h *LambdaHandler) HandleRequest(ctx context.Context, sqsEvent events.SQSEvent) error {
	for i := range sqsEvent.Records {
		message := &sqsEvent.Records[i]
		var data ReportData
		err := json.Unmarshal([]byte(message.Body), &data)
		if err != nil {
			h.logger.ErrorContext(ctx, "failed to unmarshal SQS message body",
				slog.String("message_id", message.MessageId),
				slog.String("error", err.Error()),
			)
			// Continue to next message instead of failing the whole batch
			continue
		}

		h.logger.InfoContext(ctx, "received report data",
			slog.String("message_id", message.MessageId),
			slog.String("endpoint", data.Endpoint),
			slog.String("data_source", data.DataSource),
			slog.String("client_id", data.ClientID),
			slog.Bool("success", data.Success),
			slog.Time("timestamp", data.Timestamp),
			slog.Int("status_code", data.StatusCode),
		)
	}

	return nil
}
