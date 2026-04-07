package veteran

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/cmsgov/emmy-api/pkg/core"
)

const disabilityRatingPath = "/restricted/disability_rating"

var ErrNotFound = errors.New("veteran not found")

type Service interface {
	LookupDisabilityRating(ctx context.Context, req Request) (Response, error)
}

type HTTPTransport interface {
	Do(req *http.Request) (*http.Response, error)
}

type Options struct {
	HTTPClient HTTPTransport
	Logger     *slog.Logger
	Timeout    time.Duration
}

type Request struct {
	FirstName   string   `json:"firstName"`
	MiddleName  string   `json:"middleName,omitempty"`
	LastName    string   `json:"lastName"`
	DateOfBirth string   `json:"dateOfBirth"`
	SSN         string   `json:"ssn,omitempty"`
	Address     *Address `json:"address,omitempty"`
}

type Address struct {
	Street1    string `json:"street1,omitempty"`
	Street2    string `json:"street2,omitempty"`
	Street3    string `json:"street3,omitempty"`
	City       string `json:"city,omitempty"`
	State      string `json:"state,omitempty"`
	PostalCode string `json:"postalCode,omitempty"`
	Country    string `json:"country,omitempty"`
}

type Response struct {
	CombinedDisabilityRating int             `json:"combinedDisabilityRating"`
	RawData                  any             `json:"rawData"`
	DataSource               core.DataSource `json:"dataSource"`
}

type service struct {
	cfg    *core.VAConfig
	client HTTPTransport
	logger *slog.Logger
	opts   Options
}

type disabilityRatingRequest struct {
	SSN                string `json:"ssn,omitempty"`
	FirstName          string `json:"first_name"`
	MiddleName         string `json:"middle_name,omitempty"`
	LastName           string `json:"last_name"`
	BirthDate          string `json:"birth_date"`
	StreetAddressLine1 string `json:"street_address_line1,omitempty"`
	StreetAddressLine2 string `json:"street_address_line2,omitempty"`
	StreetAddressLine3 string `json:"street_address_line3,omitempty"`
	City               string `json:"city,omitempty"`
	State              string `json:"state,omitempty"`
	Zipcode            string `json:"zipcode,omitempty"`
	Country            string `json:"country,omitempty"`
}

type disabilityRatingResponse struct {
	Data struct {
		Attributes struct {
			CombinedDisabilityRating int `json:"combined_disability_rating"`
		} `json:"attributes"`
	} `json:"data"`
}

func New(cfg *core.VAConfig, opts Options) Service {
	logger := opts.Logger
	if logger == nil {
		logger = slog.Default()
	}

	logger = logger.With(
		slog.String("component", "veteran"),
		slog.String("vendor", "va"),
	)

	if opts.Timeout == 0 && cfg != nil && cfg.TimeoutSeconds > 0 {
		opts.Timeout = time.Duration(cfg.TimeoutSeconds) * time.Second
	}

	client := opts.HTTPClient
	if client == nil && cfg != nil {
		client = vaHTTPClient(context.Background(), cfg)
	}

	return &service{
		cfg:    cfg,
		client: client,
		logger: logger,
		opts:   opts,
	}
}

func (s *service) LookupDisabilityRating(ctx context.Context, reqBody Request) (Response, error) {
	if err := s.validateConfig(); err != nil {
		return Response{}, err
	}

	if s.opts.Timeout > 0 {
		if _, hasDeadline := ctx.Deadline(); !hasDeadline {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, s.opts.Timeout)
			defer cancel()
		}
	}

	vaReq := toDisabilityRatingRequest(reqBody)
	body, err := json.Marshal(vaReq)
	if err != nil {
		return Response{}, fmt.Errorf("marshal disability rating body: %w", err)
	}

	url := strings.TrimRight(s.cfg.BaseURL, "/") + disabilityRatingPath
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return Response{}, fmt.Errorf("create disability rating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	start := time.Now()
	resp, err := s.client.Do(req)
	latency := time.Since(start)
	if err != nil {
		s.logger.ErrorContext(ctx, "va disability rating request failed",
			slog.Any("error", err),
			slog.Duration("latency", latency),
		)
		return Response{}, fmt.Errorf("do disability rating request: %w", err)
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusNotFound {
		return Response{}, ErrNotFound
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		snippet := string(respBytes)
		if len(snippet) > 800 {
			snippet = snippet[:800] + "..."
		}

		s.logger.ErrorContext(ctx, "va disability rating non-2xx",
			slog.Int("status", resp.StatusCode),
			slog.String("body_snippet", snippet),
			slog.Duration("latency", latency),
		)
		return Response{}, fmt.Errorf("va disability rating failed: status=%d", resp.StatusCode)
	}

	var out disabilityRatingResponse
	if err := json.Unmarshal(respBytes, &out); err != nil {
		return Response{}, fmt.Errorf("decode disability rating response: %w", err)
	}

	var rawBody any
	if err := json.Unmarshal(respBytes, &rawBody); err != nil {
		return Response{}, fmt.Errorf("decode raw disability rating body: %w", err)
	}

	return Response{
		CombinedDisabilityRating: out.Data.Attributes.CombinedDisabilityRating,
		RawData:                  rawBody,
		DataSource:               core.DataSourceVA,
	}, nil
}

func toDisabilityRatingRequest(reqBody Request) disabilityRatingRequest {
	out := disabilityRatingRequest{
		SSN:        reqBody.SSN,
		FirstName:  reqBody.FirstName,
		MiddleName: reqBody.MiddleName,
		LastName:   reqBody.LastName,
		BirthDate:  reqBody.DateOfBirth,
	}

	if reqBody.Address != nil {
		out.StreetAddressLine1 = reqBody.Address.Street1
		out.StreetAddressLine2 = reqBody.Address.Street2
		out.StreetAddressLine3 = reqBody.Address.Street3
		out.City = reqBody.Address.City
		out.State = reqBody.Address.State
		out.Zipcode = reqBody.Address.PostalCode
		out.Country = reqBody.Address.Country
	}

	return out
}

func (s *service) validateConfig() error {
	if s.cfg == nil {
		return errors.New("va config is required")
	}

	switch {
	case s.cfg.BaseURL == "":
		return errors.New("va base url is required")
	case s.cfg.TokenURL == "":
		return errors.New("va token url is required")
	case s.cfg.ClientID == "":
		return errors.New("va client id is required")
	case s.cfg.TokenAudience == "":
		return errors.New("va token audience is required")
	case s.cfg.PrivateKeyPath == "":
		return errors.New("va private key path is required")
	case s.client == nil:
		return errors.New("va http client is required")
	}

	return nil
}
