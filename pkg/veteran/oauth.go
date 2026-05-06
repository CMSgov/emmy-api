package veteran

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/cmsgov/emmy-api/pkg/core"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

const (
	disabilityRatingScope   = "disability_rating_restricted.read"
	clientAssertionType     = "urn:ietf:params:oauth:client-assertion-type:jwt-bearer"
	clientAssertionLifetime = 60 * time.Second
)

var (
	errVAConfigRequired            = errors.New("va config is required")
	errVATokenClientRequired       = errors.New("va token client is required")
	errVATokenRequestFailed        = errors.New("token request failed")
	errVATokenResponseMissingToken = errors.New("token response missing access_token")
	errVAPrivateKeyPathRequired    = errors.New("va private key path is required")
	errVAPrivateKeyInvalidPEM      = errors.New("private key is not valid PEM")
	errVAPrivateKeyNotRSA          = errors.New("private key is not RSA")
)

type tokenResponse struct {
	AccessToken string `json:"access_token"` //nolint:tagliatelle // OAuth response uses snake_case.
	TokenType   string `json:"token_type"`   //nolint:tagliatelle // OAuth response uses snake_case.
	ExpiresIn   int64  `json:"expires_in"`   //nolint:tagliatelle // OAuth response uses snake_case.
}

type vaTokenSource struct {
	cfg    *core.VAConfig
	client *http.Client
	now    func() time.Time
	logger *slog.Logger
}

func vaHTTPClient(ctx context.Context, cfg *core.VAConfig) *http.Client {
	return vaHTTPClientWithBase(ctx, cfg, newVABaseHTTPClient())
}

func vaHTTPClientWithBase(_ context.Context, cfg *core.VAConfig, base *http.Client) *http.Client {
	if base == nil {
		base = newVABaseHTTPClient()
	}

	baseTransport := base.Transport
	if baseTransport == nil {
		baseTransport = http.DefaultTransport
	}

	source := oauth2.ReuseTokenSource(nil, &vaTokenSource{
		cfg:    cfg,
		client: base,
		now:    time.Now,
	})

	return &http.Client{
		CheckRedirect: base.CheckRedirect,
		Transport: &oauth2.Transport{
			Source: source,
			Base:   baseTransport,
		},
	}
}

func newVABaseHTTPClient() *http.Client {
	return &http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			if len(via) > 0 {
				r.Header = via[0].Header.Clone()
			}
			return nil
		},
	}
}

func (s *vaTokenSource) Token() (*oauth2.Token, error) {
	if s == nil || s.cfg == nil {
		return nil, errVAConfigRequired
	}
	if s.client == nil {
		return nil, errVATokenClientRequired
	}
	if s.now == nil {
		s.now = time.Now
	}

	if s.logger == nil {
		s.logger = slog.Default()
	}

	assertion, err := signedClientAssertion(s.cfg, s.now())
	if err != nil {
		return nil, fmt.Errorf("build client assertion: %w", err)
	}

	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_assertion_type", clientAssertionType)
	form.Set("client_assertion", assertion)
	form.Set("scope", disabilityRatingScope)

	req, err := http.NewRequest(http.MethodPost, s.cfg.TokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request token: %w", err)
	}

	respBytes, err := io.ReadAll(resp.Body)
	if closeErr := resp.Body.Close(); closeErr != nil {
		return nil, fmt.Errorf("close token response body: %w", closeErr)
	}
	if err != nil {
		return nil, fmt.Errorf("read token response body: %w", err)
	}
	s.logger.Error("va oauth token request failed",
		slog.Any("error", resp.Status),
	)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("%w: status=%d", errVATokenRequestFailed, resp.StatusCode)
	}

	var tokenResp tokenResponse
	if err := json.Unmarshal(respBytes, &tokenResp); err != nil {
		return nil, fmt.Errorf("decode token response: %w", err)
	}
	if tokenResp.AccessToken == "" {
		return nil, errVATokenResponseMissingToken
	}

	tokenType := tokenResp.TokenType
	if tokenType == "" {
		tokenType = "Bearer"
	}

	token := &oauth2.Token{
		AccessToken: tokenResp.AccessToken,
		TokenType:   tokenType,
	}
	if tokenResp.ExpiresIn > 0 {
		token.Expiry = s.now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	}

	return token, nil
}

func signedClientAssertion(cfg *core.VAConfig, now time.Time) (string, error) {
	key, err := loadRSAPrivateKey(cfg.PrivateKeyPath)
	if err != nil {
		return "", err
	}

	headerJSON, err := json.Marshal(map[string]string{
		"alg": "RS256",
		"typ": "JWT",
	})
	if err != nil {
		return "", fmt.Errorf("marshal jwt header: %w", err)
	}

	payloadJSON, err := json.Marshal(map[string]any{
		"iss": cfg.ClientID,
		"sub": cfg.ClientID,
		"aud": cfg.TokenAudience,
		"jti": strings.ToLower(uuid.NewString()),
		"iat": now.Unix(),
		"exp": now.Add(clientAssertionLifetime).Unix(),
	})
	if err != nil {
		return "", fmt.Errorf("marshal jwt claims: %w", err)
	}

	encodedHeader := base64.RawURLEncoding.EncodeToString(headerJSON)
	encodedPayload := base64.RawURLEncoding.EncodeToString(payloadJSON)
	signingInput := encodedHeader + "." + encodedPayload

	hash := sha256.Sum256([]byte(signingInput))
	signature, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hash[:])
	if err != nil {
		return "", fmt.Errorf("sign jwt: %w", err)
	}

	encodedSignature := base64.RawURLEncoding.EncodeToString(signature)
	return signingInput + "." + encodedSignature, nil
}

func loadRSAPrivateKey(path string) (*rsa.PrivateKey, error) {
	if strings.TrimSpace(path) == "" {
		return nil, errVAPrivateKeyPathRequired
	}

	//nolint:gosec // Path is sourced from validated runtime configuration.
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read private key: %w", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errVAPrivateKeyInvalidPEM
	}

	if key, parseErr := x509.ParsePKCS1PrivateKey(block.Bytes); parseErr == nil {
		return key, nil
	}

	pkcs8Key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}

	key, ok := pkcs8Key.(*rsa.PrivateKey)
	if !ok {
		return nil, errVAPrivateKeyNotRSA
	}

	return key, nil
}
