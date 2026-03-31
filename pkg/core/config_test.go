package core

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfigFromEnv_LoadsVAConfig(t *testing.T) {
	t.Setenv("VA_BASE_URL", "https://sandbox-api.va.gov/services/veteran_verification/v2")
	t.Setenv("VA_TOKEN_URL", "https://sandbox-api.va.gov/oauth2/token")
	t.Setenv("VA_CLIENT_ID", "client-id")
	t.Setenv("VA_AUD", "https://example.okta.com/oauth2/default/v1/token")
	t.Setenv("VA_PRIVATE_KEY_PATH", "/tmp/va-private-key.pem")
	t.Setenv("VA_TIMEOUT_SECONDS", "9")

	cfg, err := NewConfigFromEnv()
	require.NoError(t, err)

	require.Equal(t, "https://sandbox-api.va.gov/services/veteran_verification/v2", cfg.VA.BaseURL)
	require.Equal(t, "https://sandbox-api.va.gov/oauth2/token", cfg.VA.TokenURL)
	require.Equal(t, "client-id", cfg.VA.ClientID)
	require.Equal(t, "https://example.okta.com/oauth2/default/v1/token", cfg.VA.TokenAudience)
	require.Equal(t, "/tmp/va-private-key.pem", cfg.VA.PrivateKeyPath)
	require.Equal(t, 9, cfg.VA.TimeoutSeconds)
}
