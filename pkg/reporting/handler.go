package reporting

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
)

type eventStore interface {
	LookupAgencyName(ctx context.Context, clientID string) (string, error)
	InsertAPIEvent(ctx context.Context, data *ReportData, agencyName *string) error
}

type sqlEventStore struct {
	db *sql.DB
}

func (s *sqlEventStore) LookupAgencyName(ctx context.Context, clientID string) (string, error) {
	const query = `SELECT agency_name FROM client_agencies WHERE client_id = $1`

	var agencyName string

	row := s.db.QueryRowContext(ctx, query, clientID)

	if err := row.Scan(&agencyName); err != nil {
		return "", err
	}

	return agencyName, nil
}

func (s *sqlEventStore) InsertAPIEvent(ctx context.Context, data *ReportData, agencyName *string) error {
	const query = `INSERT INTO api_events (data_source, endpoint, client_id, agency_name, success, timestamp, status_code)
				  VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := s.db.ExecContext(ctx, query,
		data.DataSource,
		data.Endpoint,
		data.ClientID,
		agencyName,
		data.Success,
		data.Timestamp,
		data.StatusCode,
	)

	return err
}

// LambdaHandler handles SQS events and logs ReportData
type LambdaHandler struct {
	logger *slog.Logger
	store  eventStore
}

// NewLambdaHandler creates a new LambdaHandler
func NewLambdaHandler(logger *slog.Logger, db *sql.DB) *LambdaHandler {
	if logger == nil {
		logger = slog.Default()
	}

	var store eventStore

	if db != nil {
		store = &sqlEventStore{db: db}
	}

	return &LambdaHandler{
		logger: logger.With(slog.String("component", "reporting-lambda")),
		store:  store,
	}
}

// HandleRequest processes SQS messages
func (h *LambdaHandler) HandleRequest(ctx context.Context, sqsEvent events.SQSEvent) error {
	for i := range sqsEvent.Records {
		message := &sqsEvent.Records[i]

		var data ReportData

		if err := json.Unmarshal([]byte(message.Body), &data); err != nil {
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

		if h.store != nil {
			var agencyName *string

			name, err := h.store.LookupAgencyName(ctx, data.ClientID)

			switch {
			case err == nil:
				agencyName = &name
			case errors.Is(err, sql.ErrNoRows):
				h.logger.WarnContext(ctx, "no agency mapping found for client id",
					slog.String("message_id", message.MessageId),
					slog.String("client_id", data.ClientID),
				)
			default:
				h.logger.ErrorContext(ctx, "failed to look up agency name for client id",
					slog.String("message_id", message.MessageId),
					slog.String("client_id", data.ClientID),
					slog.String("error", err.Error()),
				)
			}

			if err = h.store.InsertAPIEvent(ctx, &data, agencyName); err != nil {
				h.logger.ErrorContext(ctx, "failed to insert report data into database",
					slog.String("message_id", message.MessageId),
					slog.String("error", err.Error()),
				)
				// We continue even if DB insert fails, as the data is already logged
			}
		}
	}

	return nil
}
