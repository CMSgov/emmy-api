package education

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/cmsgov/emmy-api/pkg/core"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func nscHTTPClient(ctx context.Context, cfg *core.NSCConfig, logger *slog.Logger) *http.Client {
	if logger == nil {
		logger = slog.Default()
	}

	base := &http.Client{
		Transport: &nscLoggingTransport{
			base:     http.DefaultTransport,
			logger:   logger,
			tokenURL: cfg.TokenURL,
		},
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			if len(via) > 0 {
				r.Header = via[0].Header.Clone()
			}
			return nil
		},
	}

	ctx = context.WithValue(ctx, oauth2.HTTPClient, base)

	cc := &clientcredentials.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		TokenURL:     cfg.TokenURL,
		Scopes:       []string{"vs.api.insights"},
	}

	return oauth2.NewClient(ctx, &nscLoggingTokenSource{
		base:   cc.TokenSource(ctx),
		logger: logger,
	})
}

type nscLoggingTokenSource struct {
	base   oauth2.TokenSource
	logger *slog.Logger
}

func (s *nscLoggingTokenSource) Token() (*oauth2.Token, error) {
	start := time.Now()
	s.logger.Info("nsc oauth token request started")

	token, err := s.base.Token()
	latency := time.Since(start)
	if err != nil {
		s.logger.Error("nsc oauth token request failed",
			slog.Any("error", err),
			slog.Duration("latency", latency),
		)
		return nil, err
	}

	s.logger.Info("nsc oauth token request succeeded",
		slog.Duration("latency", latency),
		slog.Bool("access_token_present", token != nil && token.AccessToken != ""),
		slog.String("token_type", token.TokenType),
		slog.Bool("expiry_set", token != nil && !token.Expiry.IsZero()),
		slog.Int64("expiry_unix", tokenExpiryUnix(token)),
	)

	return token, nil
}

type nscLoggingTransport struct {
	base     http.RoundTripper
	logger   *slog.Logger
	tokenURL string
}

func (t *nscLoggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	transport := t.base
	if transport == nil {
		transport = http.DefaultTransport
	}

	start := time.Now()
	resp, err := transport.RoundTrip(req)
	latency := time.Since(start)

	if req != nil && req.URL.String() == t.tokenURL {
		if err != nil {
			t.logger.Error("nsc oauth token endpoint request failed",
				slog.Any("error", err),
				slog.Duration("latency", latency),
			)
			return nil, err
		}

		t.logger.Info("nsc oauth token endpoint response received",
			slog.Int("status", resp.StatusCode),
			slog.String("content_type", resp.Header.Get("Content-Type")),
			slog.Duration("latency", latency),
		)
		return resp, nil
	}

	if req != nil && req.URL != nil {
		t.logger.Debug("nsc outbound request dispatched",
			slog.String("host", req.URL.Host),
			slog.String("path", req.URL.Path),
			slog.Bool("authorization_header_present", req.Header.Get("Authorization") != ""),
			slog.Duration("latency", latency),
		)
	}

	return resp, err
}

func tokenExpiryUnix(token *oauth2.Token) int64 {
	if token == nil || token.Expiry.IsZero() {
		return 0
	}

	return token.Expiry.Unix()
}
