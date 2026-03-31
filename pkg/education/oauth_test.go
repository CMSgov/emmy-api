package education

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cmsgov/emmy-api/pkg/core"
	"github.com/stretchr/testify/require"
)

func TestNSCHTTPClient_PreservesAuthOnRedirect(t *testing.T) {
	var seenAuth string

	// Test NSC server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {

		// OAuth token endpoint
		case "/token":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"access_token": "test-token",
				"token_type":   "Bearer",
				"expires_in":   3600,
			})

		// first submit then redirect
		case "/submit":
			http.Redirect(w, r, "/final", http.StatusFound)

		// redirect location
		case "/final":
			seenAuth = r.Header.Get("Authorization")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`ok`))

		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	cfg := &core.NSCConfig{
		ClientID:     "abc",
		ClientSecret: "secret",
		TokenURL:     ts.URL + "/token",
	}

	client := nscHTTPClient(context.Background(), cfg)

	req, err := http.NewRequest("GET", ts.URL+"/submit", nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	require.Equal(t, "Bearer test-token", seenAuth)
}
