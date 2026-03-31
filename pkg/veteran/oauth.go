package veteran

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/cmsgov/emmy-api/pkg/core"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"golang.org/x/oauth2"
)

const (
	disabilityRatingScope   = "disability_rating_restricted.read"
	clientAssertionType     = "urn:ietf:params:oauth:client-assertion-type:jwt-bearer"
	clientAssertionLifetime = 60 * time.Second
)

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

type vaTokenSource struct {
	cfg    *core.VAConfig
	client *http.Client
	now    func() time.Time
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
		return nil, errors.New("va config is required")
	}
	if s.client == nil {
		return nil, errors.New("va token client is required")
	}
	if s.now == nil {
		s.now = time.Now
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
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("token request failed: status=%d", resp.StatusCode)
	}

	var tokenResp tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("decode token response: %w", err)
	}
	if tokenResp.AccessToken == "" {
		return nil, errors.New("token response missing access_token")
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

	token, err := jwt.NewBuilder().
		Audience([]string{cfg.TokenAudience}).
		Issuer(cfg.ClientID).
		Subject(cfg.ClientID).
		JwtID(strings.ToLower(uuid.NewString())).
		IssuedAt(now).
		Expiration(now.Add(clientAssertionLifetime)).
		Build()
	if err != nil {
		return "", fmt.Errorf("build jwt claims: %w", err)
	}

	signed, err := jwt.Sign(token, jwt.WithKey(jwa.RS256, key))
	if err != nil {
		return "", fmt.Errorf("sign jwt: %w", err)
	}

	return string(signed), nil
}

func loadRSAPrivateKey(path string) (*rsa.PrivateKey, error) {
	if strings.TrimSpace(path) == "" {
		return nil, errors.New("va private key path is required")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read private key: %w", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("private key is not valid PEM")
	}

	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}

	pkcs8Key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}

	key, ok := pkcs8Key.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("private key is not RSA")
	}

	return key, nil
}
