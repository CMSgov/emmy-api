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
	"github.com/cmsgov/emmy-api/pkg/education"
)

const disabilityRatingPath = "/restricted/disability_rating"

var ErrNotFound = errors.New("veteran not found")

var (
	errVADisabilityRatingFailed = errors.New("va disability rating failed")
	errVAServiceConfigRequired  = errors.New("va config is required")
	errVABaseURLRequired        = errors.New("va base url is required")
	errVATokenURLRequired       = errors.New("va token url is required")
	errVAClientIDRequired       = errors.New("va client id is required")
	errVATokenAudienceRequired  = errors.New("va token audience is required")
	errVAPrivateKeyRequired     = errors.New("va private key path is required")
	errVAHTTPClientRequired     = errors.New("va http client is required")
)

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
	Address     *Address `json:"address,omitempty"`
	FirstName   string   `json:"firstName"`
	MiddleName  string   `json:"middleName,omitempty"`
	LastName    string   `json:"lastName"`
	DateOfBirth string   `json:"dateOfBirth"`
	SSN         string   `json:"ssn,omitempty"`
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
	RawData                  any                `json:"rawData"`
	LegalEffectiveDate       *string            `json:"legalEffectiveDate"`
	CombinedEffectiveDate    *string            `json:"combinedEffectiveDate"`
	EarliestRatingEndDate    *string            `json:"earliestRatingEndDate"`
	DataSource               core.DataSource    `json:"dataSource"`
	Metadata                 education.Metadata `json:"metadata"`
	CombinedDisabilityRating int                `json:"combinedDisabilityRating"`
}

type service struct {
	cfg    *core.VAConfig
	client HTTPTransport
	logger *slog.Logger
	opts   Options
}

type disabilityRatingRequest struct {
	SSN                string `json:"ssn,omitempty"`
	FirstName          string `json:"first_name"`                     //nolint:tagliatelle // VA API request uses snake_case.
	MiddleName         string `json:"middle_name,omitempty"`          //nolint:tagliatelle // VA API request uses snake_case.
	LastName           string `json:"last_name"`                      //nolint:tagliatelle // VA API request uses snake_case.
	BirthDate          string `json:"birth_date"`                     //nolint:tagliatelle // VA API request uses snake_case.
	StreetAddressLine1 string `json:"street_address_line1,omitempty"` //nolint:tagliatelle // VA API request uses snake_case.
	StreetAddressLine2 string `json:"street_address_line2,omitempty"` //nolint:tagliatelle // VA API request uses snake_case.
	StreetAddressLine3 string `json:"street_address_line3,omitempty"` //nolint:tagliatelle // VA API request uses snake_case.
	City               string `json:"city,omitempty"`
	State              string `json:"state,omitempty"`
	Zipcode            string `json:"zipcode,omitempty"`
	Country            string `json:"country,omitempty"`
}

type disabilityRatingResponse struct {
	Data struct {
		Attributes struct {
			CombinedEffectiveDate    *string            `json:"combined_effective_date"`    //nolint:tagliatelle // VA API response uses snake_case.
			LegalEffectiveDate       *string            `json:"legal_effective_date"`       //nolint:tagliatelle // VA API response uses snake_case.
			IndividualRatings        []individualRating `json:"individual_ratings"`         //nolint:tagliatelle // VA API response uses snake_case.
			CombinedDisabilityRating int                `json:"combined_disability_rating"` //nolint:tagliatelle // VA API response uses snake_case.
		} `json:"attributes"`
	} `json:"data"`
}

type individualRating struct {
	EffectiveDate      *string `json:"effective_date"`  //nolint:tagliatelle // VA API uses snake_case
	RatingEndDate      *string `json:"rating_end_date"` //nolint:tagliatelle // VA API uses snake_case
	Decision           string  `json:"decision"`
	DiagnosticText     string  `json:"diagnostic_text"`      //nolint:tagliatelle // VA API uses snake_case
	DiagnosticTypeCode string  `json:"diagnostic_type_code"` //nolint:tagliatelle // VA API uses snake_case
	DiagnosticTypeName string  `json:"diagnostic_type_name"` //nolint:tagliatelle // VA API uses snake_case
	DisabilityRatingID string  `json:"disability_rating_id"` //nolint:tagliatelle // VA API uses snake_case
	RatingPercentage   int     `json:"rating_percentage"`    //nolint:tagliatelle // VA API uses snake_case
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

//nolint:gocritic // Request is treated as an immutable API payload value.
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

	respBytes, err := io.ReadAll(resp.Body)
	if closeErr := resp.Body.Close(); closeErr != nil {
		return Response{}, fmt.Errorf("close disability rating response body: %w", closeErr)
	}
	if err != nil {
		return Response{}, fmt.Errorf("read disability rating response body: %w", err)
	}

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
		return Response{}, fmt.Errorf("%w: status=%d", errVADisabilityRatingFailed, resp.StatusCode)
	}

	var out disabilityRatingResponse
	if err := json.Unmarshal(respBytes, &out); err != nil {
		return Response{}, fmt.Errorf("decode disability rating response: %w", err)
	}

	var earliestEndDate *string
	for _, r := range out.Data.Attributes.IndividualRatings {
		if r.RatingEndDate != nil && *r.RatingEndDate != "" {
			if earliestEndDate == nil || *r.RatingEndDate < *earliestEndDate {
				earliestEndDate = r.RatingEndDate
			}
		}
	}

	var rawBody any
	if err := json.Unmarshal(respBytes, &rawBody); err != nil {
		return Response{}, fmt.Errorf("decode raw disability rating body: %w", err)
	}

	return Response{
		CombinedDisabilityRating: out.Data.Attributes.CombinedDisabilityRating,
		CombinedEffectiveDate:    out.Data.Attributes.CombinedEffectiveDate,
		LegalEffectiveDate:       out.Data.Attributes.LegalEffectiveDate,
		EarliestRatingEndDate:    earliestEndDate,
		RawData:                  rawBody,
		DataSource:               core.DataSourceVA,
	}, nil
}

//nolint:gocritic // Request is treated as an immutable API payload value.
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
		return errVAServiceConfigRequired
	}

	switch {
	case s.cfg.BaseURL == "":
		return errVABaseURLRequired
	case s.cfg.TokenURL == "":
		return errVATokenURLRequired
	case s.cfg.ClientID == "":
		return errVAClientIDRequired
	case s.cfg.TokenAudience == "":
		return errVATokenAudienceRequired
	case s.cfg.PrivateKeyPath == "":
		return errVAPrivateKeyRequired
	case s.client == nil:
		return errVAHTTPClientRequired
	}

	return nil
}
