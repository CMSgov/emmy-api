package education

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

func (s *service) Submit(ctx context.Context, reqBody Request) (Response, error) {
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

	body, err := json.Marshal(reqBody)
	if err != nil {
		log.Error("nsc submit marshal failed", slog.Any("error", err))
		return Response{}, fmt.Errorf("marshal submit body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.cfg.SubmitURL, bytes.NewReader(body))
	if err != nil {
		log.Error("nsc submit create request failed", slog.Any("error", err))
		return Response{}, fmt.Errorf("create submit request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	log.Debug("nsc submit request prepared",
		slog.String("method", req.Method),
		slog.String("host", req.URL.Host),
		slog.String("path", req.URL.Path),
	)

	start := time.Now()
	resp, err := s.client.Do(req)
	latency := time.Since(start)

	if err != nil {
		log.Error("nsc submit request failed",
			slog.Any("error", err),
			slog.Duration("latency", latency),
		)
		return Response{}, fmt.Errorf("submit request: %w", err)
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)

	log.Info("nsc submit response received",
		slog.Int("status", resp.StatusCode),
		slog.String("content_type", resp.Header.Get("Content-Type")),
		slog.Duration("latency", latency),
	)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		snippet := string(respBytes)
		if len(snippet) > 800 {
			snippet = snippet[:800] + "..."
		}

		log.Error("nsc submit non-2xx",
			slog.Int("status", resp.StatusCode),
			slog.String("www_authenticate", resp.Header.Get("WWW-Authenticate")),
			slog.String("body_snippet", snippet),
		)

		return Response{}, fmt.Errorf("nsc submit failed: status=%d", resp.StatusCode)
	}

	var nscResp NSCResponse
	if err := json.Unmarshal(respBytes, &nscResp); err != nil {
		log.Error("nsc submit decode failed", slog.Any("error", err))
		return Response{}, fmt.Errorf("decode nsc response: %w", err)
	}

	var out Response
	if nscResp.TransactionDetails.NSCHit == "Y" {
		out.EnrollmentStatus = "FULL_TIME"
	} else {
		out.EnrollmentStatus = "UNKNOWN"
	}

	log.Debug("nsc submit decoded successfully")
	return out, nil
}
