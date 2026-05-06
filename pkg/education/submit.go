package education

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

var (
	errLegacyEnrollmentStatusRequired = errors.New("enrollmentStatus is required")
	errUnsupportedLegacyNSCStatusCode = errors.New("unsupported legacy nsc status code")
)

func (s *service) LookupEnrollmentStatus(ctx context.Context, reqBody Request) (Response, error) {
	if s.opts.Timeout > 0 {
		if _, hasDeadline := ctx.Deadline(); !hasDeadline {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, s.opts.Timeout)
			defer cancel()
		}
	}

	log := s.logger.With(
		slog.String("nsc_submit_url", s.cfg.SubmitURL),
	)

	nscReqBody := toNSCRequest(s.cfg, reqBody)

	body, err := json.Marshal(nscReqBody)
	if err != nil {
		log.ErrorContext(ctx, "nsc submit marshal failed", slog.Any("error", err))
		return Response{}, fmt.Errorf("marshal submit body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.cfg.SubmitURL, bytes.NewReader(body))
	if err != nil {
		log.ErrorContext(ctx, "nsc submit create request failed", slog.Any("error", err))
		return Response{}, fmt.Errorf("create submit request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	log.InfoContext(ctx, "nsc submit request prepared",
		slog.String("method", req.Method),
		slog.String("host", req.URL.Host),
		slog.String("path", req.URL.Path),
		slog.String("content_type", req.Header.Get("Content-Type")),
		slog.String("accept", req.Header.Get("Accept")),
		slog.Int("payload_bytes", len(body)),
		slog.Bool("has_ssn", strings.TrimSpace(reqBody.SSN) != ""),
		slog.Bool("has_date_of_birth", strings.TrimSpace(reqBody.DateOfBirth) != ""),
		slog.Bool("has_middle_name", strings.TrimSpace(reqBody.MiddleName) != ""),
		slog.Bool("has_address", nscReqBody.Address1 != ""),
		slog.Bool("has_context_deadline", hasContextDeadline(ctx)),
		slog.Int64("deadline_remaining_ms", deadlineRemainingMillis(ctx)),
	)

	start := time.Now()
	resp, err := s.client.Do(req)
	latency := time.Since(start)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(ctx.Err(), context.DeadlineExceeded) {
			log.ErrorContext(ctx, "nsc submit request timed out",
				slog.Any("error", err),
				slog.Duration("latency", latency),
				slog.Int64("deadline_remaining_ms", deadlineRemainingMillis(ctx)),
			)
		} else {
			log.ErrorContext(ctx, "nsc submit request failed",
				slog.Any("error", err),
				slog.Duration("latency", latency),
				slog.Bool("context_canceled", errors.Is(ctx.Err(), context.Canceled)),
				slog.Bool("context_deadline_exceeded", errors.Is(ctx.Err(), context.DeadlineExceeded)),
			)
		}
		return Response{}, fmt.Errorf("submit request: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.WarnContext(ctx, errCloseNSCSubmitResponseBody.Error(), slog.Any("error", closeErr))
		}
	}()

	respBytes, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return Response{}, fmt.Errorf("%w: %w", errReadNSCSubmitResponseBody, readErr)
	}
	snippet := bodySnippet(respBytes)

	log.InfoContext(ctx, "nsc submit response received",
		slog.Int("status", resp.StatusCode),
		slog.String("status_text", http.StatusText(resp.StatusCode)),
		slog.String("content_type", resp.Header.Get("Content-Type")),
		slog.String("server", resp.Header.Get("Server")),
		slog.String("x_cache", resp.Header.Get("X-Cache")),
		slog.String("via", resp.Header.Get("Via")),
		slog.Int64("content_length", resp.ContentLength),
		slog.Duration("latency", latency),
	)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.ErrorContext(ctx, "nsc submit non-2xx",
			slog.Int("status", resp.StatusCode),
			slog.String("status_text", http.StatusText(resp.StatusCode)),
			slog.String("content_type", resp.Header.Get("Content-Type")),
			slog.String("www_authenticate", resp.Header.Get("WWW-Authenticate")),
			slog.String("server", resp.Header.Get("Server")),
			slog.String("x_cache", resp.Header.Get("X-Cache")),
			slog.String("via", resp.Header.Get("Via")),
			slog.Duration("latency", latency),
			slog.String("body_snippet", snippet),
		)

		return Response{}, fmt.Errorf("%w: status=%d", errNSCSubmitFailed, resp.StatusCode)
	}

	var out nscResponse
	if err := json.Unmarshal(respBytes, &out); err != nil {
		log.ErrorContext(ctx, "nsc submit decode failed",
			slog.Any("error", err),
			slog.String("content_type", resp.Header.Get("Content-Type")),
			slog.Duration("latency", latency),
			slog.String("body_snippet", snippet),
		)
		return Response{}, fmt.Errorf("decode nsc response: %w", err)
	}

	var rawBody any
	if err := json.Unmarshal(respBytes, &rawBody); err != nil {
		log.ErrorContext(ctx, "nsc submit raw decode failed",
			slog.Any("error", err),
			slog.Duration("latency", latency),
			slog.String("body_snippet", snippet),
		)
		return Response{}, fmt.Errorf("decode raw nsc response: %w", err)
	}

	log.InfoContext(ctx, "nsc submit decoded successfully",
		slog.String("status_code", strings.TrimSpace(out.Status.Code)),
		slog.String("current_enrollment_status", currentEnrollmentStatus(out)),
		slog.Int("enrollment_records", enrollmentRecordCount(out)),
	)

	return translateNSCResponse(out, rawBody)
}

//nolint:gocritic // Request is treated as an immutable API payload value.
func toNSCRequest(cfg *core.NSCConfig, reqBody Request) nscRequest {
	accountID := ""
	if cfg != nil {
		accountID = cfg.AccountID
	}

	out := nscRequest{
		AccountID:   accountID,
		DateOfBirth: reqBody.DateOfBirth,
		LastName:    reqBody.LastName,
		FirstName:   reqBody.FirstName,
		MiddleName:  reqBody.MiddleName,
		SSN:         reqBody.SSN,
		EndClient:   "CMS",
		Terms:       "y",
	}

	if reqBody.Address != nil {
		out.Address1 = reqBody.Address.Street1
		out.Address2 = reqBody.Address.Street2
		out.City = reqBody.Address.City
		out.State = reqBody.Address.State
		out.ZipCode = reqBody.Address.PostalCode
	}

	return out
}

//nolint:gocritic // Keeping value semantics is acceptable for this internal translation helper.
func translateNSCResponse(resp nscResponse) (Response, error) {
	if isNSCNoHit(resp) || isNSCNotCurrentlyEnrolled(resp) {
		return Response{}, ErrNotFound
	}

	status, ok := resolveEnrollmentStatus(resp)
	if !ok {
		return Response{}, errNSCMissingEnrollmentStatus
	}

	details := []EnrollmentDetail{}
	for _, d := range resp.EnrollmentDetails {
		for _, ed := range d.EnrollmentData {
			s, ok := normalizeEnrollmentStatus(ed.EnrollmentStatus)
			if !ok {
				continue
			}
			details = append(details, EnrollmentDetail{
				SchoolName:       d.OfficialSchoolName,
				TermBeginDate:    ed.TermBeginDate,
				TermEndDate:      ed.TermEndDate,
				EnrollmentStatus: s,
			})
		}
	}

	return Response{
		EnrollmentStatus:  status,
		EnrollmentDetails: details,
		RawData:           rawBody,
		DataSource:        core.DataSourceNSC,
	}, nil
}

//nolint:gocritic // Response is treated as an immutable API payload value.
func isNSCNoHit(resp nscResponse) bool {
	switch strings.ToUpper(strings.TrimSpace(resp.TransactionDetails.NSCHit)) {
	case "N", "NO", "FALSE", "0":
		return true
	}

	return false
}

//nolint:gocritic // Response is treated as an immutable API payload value.
func isNSCNotCurrentlyEnrolled(resp nscResponse) bool {
	if len(resp.EnrollmentDetails) == 0 {
		return false
	}

	for _, detail := range resp.EnrollmentDetails {
		if strings.EqualFold(strings.TrimSpace(detail.CurrentEnrollmentStatus), "CN") {
			return true
		}
	}

	return false
}

//nolint:gocritic // Response is treated as an immutable API payload value.
func resolveEnrollmentStatus(resp nscResponse) (EnrollmentStatus, bool) {
	var best EnrollmentStatus

	rank := func(s EnrollmentStatus) int {
		switch s {
		case EnrollmentStatusFullTime:
			return 5
		case EnrollmentStatusThreeQuartersTime:
			return 4
		case EnrollmentStatusHalfTime:
			return 3
		case EnrollmentStatusLessThanHalfTime:
			return 2
		case EnrollmentStatusUnknown:
			return 1
		default:
			return 0
		}
	}

	for _, detail := range resp.EnrollmentDetails {
		for _, item := range detail.EnrollmentData {
			if status, ok := normalizeEnrollmentStatus(item.EnrollmentStatus); ok {
				if rank(status) > rank(best) {
					best = status
				}
			}
		}

		if status, ok := normalizeCurrentEnrollmentStatus(detail.CurrentEnrollmentStatus); ok {
			if rank(status) > rank(best) {
				best = status
			}
		}
	}

	if best != "" {
		return best, true
	}

	if isNSCPositiveHit(resp) {
		return EnrollmentStatusUnknown, true
	}

	return "", false
}

//nolint:gocritic // Response is treated as an immutable API payload value.
func isNSCPositiveHit(resp nscResponse) bool {
	switch strings.ToUpper(strings.TrimSpace(resp.TransactionDetails.NSCHit)) {
	case "Y", "YES", "TRUE", "1":
		return true
	}

	return strings.EqualFold(strings.TrimSpace(resp.Status.Code), "0")
}

func normalizeEnrollmentStatus(value string) (EnrollmentStatus, bool) {
	if value == "" {
		return "", false
	}

	normalized := strings.ToUpper(strings.TrimSpace(value))
	normalized = strings.NewReplacer("-", "_", " ", "_").Replace(normalized)

	switch normalized {
	case string(EnrollmentStatusFullTime), "F":
		return EnrollmentStatusFullTime, true
	case string(EnrollmentStatusThreeQuartersTime), "Q", "THREE_QUARTER_TIME":
		return EnrollmentStatusThreeQuartersTime, true
	case string(EnrollmentStatusHalfTime), "H":
		return EnrollmentStatusHalfTime, true
	case string(EnrollmentStatusLessThanHalfTime), "L":
		return EnrollmentStatusLessThanHalfTime, true
	case string(EnrollmentStatusUnknown), "Y":
		return EnrollmentStatusUnknown, true
	default:
		return "", false
	}
}

func normalizeCurrentEnrollmentStatus(value string) (EnrollmentStatus, bool) {
	normalized := strings.ToUpper(strings.TrimSpace(value))

	switch normalized {
	case "CC":
		return EnrollmentStatusUnknown, true
	case "CN":
		return "", false
	default:
		return normalizeEnrollmentStatus(normalized)
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}

	return ""
}

func bodySnippet(body []byte) string {
	snippet := string(body)
	if len(snippet) > 800 {
		snippet = snippet[:800] + "..."
	}

	return snippet
}

func hasContextDeadline(ctx context.Context) bool {
	_, ok := ctx.Deadline()
	return ok
}

func deadlineRemainingMillis(ctx context.Context) int64 {
	deadline, ok := ctx.Deadline()
	if !ok {
		return -1
	}

	return time.Until(deadline).Milliseconds()
}

//nolint:gocritic // Response is treated as an immutable API payload value.
func currentEnrollmentStatus(resp nscResponse) string {
	if len(resp.EnrollmentDetails) == 0 {
		return ""
	}

	return strings.TrimSpace(resp.EnrollmentDetails[0].CurrentEnrollmentStatus)
}

//nolint:gocritic // Response is treated as an immutable API payload value.
func enrollmentRecordCount(resp nscResponse) int {
	count := 0
	for _, detail := range resp.EnrollmentDetails {
		count += len(detail.EnrollmentData)
	}

	return count
}

type legacySubmitResponse struct {
	Status legacySubmitStatus `json:"status"`
}

type legacySubmitStatus struct {
	Code string `json:"code"`
}

func mapLegacyEnrollmentStatus(respBytes []byte) (EnrollmentStatus, error) {
	var legacy legacySubmitResponse
	if err := json.Unmarshal(respBytes, &legacy); err != nil {
		return "", fmt.Errorf("decode legacy nsc response: %w", err)
	}

	switch legacy.Status.Code {
	case "0":
		return EnrollmentStatusEnrolled, nil
	case "":
		return "", errLegacyEnrollmentStatusRequired
	default:
		return "", fmt.Errorf("%w %q", errUnsupportedLegacyNSCStatusCode, legacy.Status.Code)
	}
}
