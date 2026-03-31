package veteran

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/cmsgov/emmy-api/pkg/core"
	"github.com/stretchr/testify/require"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func TestVAHTTPClient_AddsBearerTokenFromClientAssertionFlow(t *testing.T) {
	keyPath := writeRSAPrivateKeyPEM(t)
	var tokenRequest *http.Request
	var protectedAuthHeader string

	base := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			switch req.URL.String() {
			case "https://example.test/token":
				tokenRequest = req.Clone(req.Context())
				body, err := io.ReadAll(req.Body)
				require.NoError(t, err)
				req.Body = io.NopCloser(strings.NewReader(string(body)))

				form, err := url.ParseQuery(string(body))
				require.NoError(t, err)
				require.Equal(t, "client_credentials", form.Get("grant_type"))
				require.Equal(t, clientAssertionType, form.Get("client_assertion_type"))
				require.Equal(t, disabilityRatingScope, form.Get("scope"))
				require.NotEmpty(t, form.Get("client_assertion"))

				claims := decodeJWTClaims(t, form.Get("client_assertion"))
				require.Equal(t, "client-id", claims["iss"])
				require.Equal(t, "client-id", claims["sub"])
				require.Equal(t, "https://example.okta.com/oauth2/default/v1/token", claims["aud"])
				require.NotEmpty(t, claims["jti"])

				return jsonResponse(http.StatusOK, `{
					"access_token":"test-token",
					"token_type":"Bearer",
					"expires_in":300
				}`), nil
			case "https://example.test/restricted/disability_rating":
				protectedAuthHeader = req.Header.Get("Authorization")
				return jsonResponse(http.StatusOK, `{"data":{"attributes":{"combined_disability_rating":70}}}`), nil
			default:
				t.Fatalf("unexpected request to %s", req.URL.String())
				return nil, nil
			}
		}),
	}

	client := vaHTTPClientWithBase(context.Background(), &core.VAConfig{
		BaseURL:        "https://example.test",
		TokenURL:       "https://example.test/token",
		ClientID:       "client-id",
		TokenAudience:  "https://example.okta.com/oauth2/default/v1/token",
		PrivateKeyPath: keyPath,
	}, base)

	req, err := http.NewRequest(http.MethodPost, "https://example.test/restricted/disability_rating", nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.NotNil(t, tokenRequest)
	require.Equal(t, "application/x-www-form-urlencoded", tokenRequest.Header.Get("Content-Type"))
	require.Equal(t, "Bearer test-token", protectedAuthHeader)
}

func TestSignedClientAssertion_UsesConfiguredAudienceAndTTL(t *testing.T) {
	keyPath := writeRSAPrivateKeyPEM(t)
	now := time.Unix(1_700_000_000, 0).UTC()

	assertion, err := signedClientAssertion(&core.VAConfig{
		ClientID:       "client-id",
		TokenAudience:  "https://example.okta.com/oauth2/default/v1/token",
		PrivateKeyPath: keyPath,
	}, now)
	require.NoError(t, err)

	claims := decodeJWTClaims(t, assertion)
	require.Equal(t, "client-id", claims["iss"])
	require.Equal(t, "client-id", claims["sub"])
	require.Equal(t, "https://example.okta.com/oauth2/default/v1/token", claims["aud"])
	require.EqualValues(t, now.Unix(), claims["iat"])
	require.EqualValues(t, now.Add(clientAssertionLifetime).Unix(), claims["exp"])
	require.NotEmpty(t, claims["jti"])
}

func jsonResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func writeRSAPrivateKeyPEM(t *testing.T) string {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	path := t.TempDir() + "/va-private-key.pem"
	pemBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})

	err = os.WriteFile(path, pemBytes, 0o600)
	require.NoError(t, err)

	return path
}

func decodeJWTClaims(t *testing.T, signed string) map[string]any {
	t.Helper()

	parts := strings.Split(signed, ".")
	require.Len(t, parts, 3)

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	require.NoError(t, err)

	var claims map[string]any
	err = json.Unmarshal(payload, &claims)
	require.NoError(t, err)

	if aud, ok := claims["aud"].([]any); ok && len(aud) == 1 {
		claims["aud"] = aud[0]
	}

	return claims
}
