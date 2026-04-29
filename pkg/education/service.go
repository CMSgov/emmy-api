package education

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/cmsgov/emmy-api/pkg/core"
	"github.com/cmsgov/emmy-api/pkg/encryption"
)

var ErrNotFound = errors.New("education enrollment not found")
var ErrDatabaseConnectionNotAvailable = errors.New("database connection not available")

type Service interface {
	LookupEnrollmentStatus(ctx context.Context, req Request) (Response, error)
	RegisterBatch(ctx context.Context, req BatchRequest) error
}

type HTTPTransport interface {
	Do(req *http.Request) (*http.Response, error)
}

type Options struct {
	// Override for testing the HTTP client
	HTTPClient HTTPTransport
	// Structured logger using slog package
	Logger *slog.Logger
	// Context timeout
	Timeout time.Duration
	// Database connection for batch enrollment
	DB *sql.DB
	// Encryption service for SSN
	Encryption encryption.Service
}

type service struct {
	cfg        *core.NSCConfig
	client     HTTPTransport
	logger     *slog.Logger
	db         *sql.DB
	encryption encryption.Service
	opts       Options
}

func New(cfg *core.NSCConfig, opts Options) Service {
	logger := opts.Logger
	if logger == nil {
		logger = slog.Default()
	}

	logger = logger.With(
		slog.String("component", "education"),
		slog.String("vendor", "nsc"),
	)

	client := opts.HTTPClient

	if client == nil {
		client = nscHTTPClient(context.Background(), cfg, logger)
	}

	return &service{
		cfg:        cfg,
		client:     client,
		logger:     logger,
		db:         opts.DB,
		encryption: opts.Encryption,
		opts:       opts,
	}
}

func (s *service) RegisterBatch(ctx context.Context, req BatchRequest) error {
	logger := slog.Default()

	if s.db == nil {
		return fmt.Errorf("RegisterBatch: s.db is nil: %w", ErrDatabaseConnectionNotAvailable)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("RegisterBatch: BeginTx failed: %w", err)
	}

	defer func() {
		if rollbackErr := tx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			logger.Error("RegisterBatch rollback failed", "error", rollbackErr)
		}
	}()

	var batchDBID string
	err = tx.QueryRowContext(
		ctx,
		"INSERT INTO enrollment_batches (batch_id, submitted_by, callback_url, status) VALUES ($1, $2, $3, $4) RETURNING id",
		req.BatchID, req.SubmittedBy, req.CallbackURL, "QUEUED",
	).Scan(&batchDBID)
	if err != nil {
		return fmt.Errorf("RegisterBatch: insert enrollment_batches failed: batch_id=%q submitted_by=%q: %w",
			req.BatchID, req.SubmittedBy, err)
	}

	for i, student := range req.Students {
		ssn := student.SSN
		if s.encryption != nil && ssn != "" {
			ssn, err = s.encryption.Encrypt(ctx, ssn)
			if err != nil {
				return fmt.Errorf("RegisterBatch: encrypt SSN failed at student index=%d record_id=%q: %w",
					i, student.RecordID, err)
			}
		}

		_, err = tx.ExecContext(
			ctx,
			"INSERT INTO batch_students (batch_db_id, record_id, first_name, last_name, date_of_birth, ssn, status) VALUES ($1, $2, $3, $4, $5, $6, $7)",
			batchDBID, student.RecordID, student.FirstName, student.LastName, student.DateOfBirth, ssn, "QUEUED",
		)
		if err != nil {
			return fmt.Errorf("RegisterBatch: insert batch_students failed at student index=%d record_id=%q batch_db_id=%q: %w",
				i, student.RecordID, batchDBID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("RegisterBatch: Commit failed: batch_id=%q batch_db_id=%q: %w",
			req.BatchID, batchDBID, err)
	}

	return nil
}
