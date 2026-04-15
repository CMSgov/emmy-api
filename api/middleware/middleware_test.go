package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type skipAuthPayload struct {
	Sub      string   `json:"sub"`
	Username string   `json:"username"`
	Scope    string   `json:"scope"`
	Groups   []string `json:"groups"`
}

func setupSkipAuthApp() *fiber.App {
	app := fiber.New()
	app.Use(SkipAuthMiddleware())
	app.Get("/whoami", func(c *fiber.Ctx) error {
		return c.JSON(skipAuthPayload{
			Sub:      c.Locals("sub").(string),
			Username: c.Locals("username").(string),
			Scope:    c.Locals("scope").(string),
			Groups:   c.Locals("groups").([]string),
		})
	})

	return app
}

func TestSubjectMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		setupHeader    func(req *http.Request)
		expectedSub    string
		expectedStatus int
	}{
		{
			name:           "No headers",
			setupHeader:    func(req *http.Request) {},
			expectedSub:    "unknown-subject",
			expectedStatus: http.StatusOK,
		},
		{
			name: "X-Sub header",
			setupHeader: func(req *http.Request) {
				req.Header.Set("X-Sub", "user-123")
			},
			expectedSub:    "user-123",
			expectedStatus: http.StatusOK,
		},
		{
			name: "Bearer Authorization header with valid JWT",
			setupHeader: func(req *http.Request) {
				tkn, _ := jwt.NewBuilder().Subject("jwt-sub-123").Build()
				signed, _ := jwt.Sign(tkn, jwt.WithKey(jwa.HS256, []byte("secret")))
				req.Header.Set("Authorization", "Bearer "+string(signed))
			},
			expectedSub:    "jwt-sub-123",
			expectedStatus: http.StatusOK,
		},
		{
			name: "Bearer Authorization header with malformed JWT",
			setupHeader: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer not-a-jwt")
			},
			expectedSub:    "", // Not checked if status is not OK
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			app.Use(SubjectMiddleware(nil))
			app.Get("/test", func(c *fiber.Ctx) error {
				return c.SendString(c.Locals("sub").(string)) //nolint:forcetypeassert // test setup guarantees type
			})

			req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
			tt.setupHeader(req)

			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.expectedStatus == http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				assert.Equal(t, tt.expectedSub, string(body))
			}
		})
	}
}

func TestMiddlewareChain(t *testing.T) {
	app := fiber.New()
	// SubjectMiddleware should not overwrite if sub is already set (e.g. by SkipAuth)
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("sub", "pre-set")
		return c.Next()
	})
	app.Use(SubjectMiddleware(nil))
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString(c.Locals("sub").(string)) //nolint:forcetypeassert // test setup guarantees type
	})

	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "pre-set", string(body))
}

func TestSkipAuthMiddleware_DefaultIdentity(t *testing.T) {
	app := setupSkipAuthApp()

	req := httptest.NewRequest(http.MethodGet, "/whoami", http.NoBody)

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, fiber.StatusOK, resp.StatusCode)

	var payload skipAuthPayload
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&payload))

	assert.Equal(t, defaultSkipAuthSub, payload.Sub)
	assert.Equal(t, defaultSkipAuthSub, payload.Username)
	assert.Equal(t, defaultSkipAuthScope, payload.Scope)
	assert.Equal(t, []string{defaultSkipAuthGroup}, payload.Groups)
}

func TestSkipAuthMiddleware_HeaderOverrides(t *testing.T) {
	app := setupSkipAuthApp()

	req := httptest.NewRequest(http.MethodGet, "/whoami", http.NoBody)
	req.Header.Set(skipAuthHeaderSub, "test-sub")
	req.Header.Set(skipAuthHeaderUsername, "test-user")
	req.Header.Set(skipAuthHeaderScope, "read:edu")
	req.Header.Set(skipAuthHeaderGroups, "admins,qa,reporting")

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, fiber.StatusOK, resp.StatusCode)

	var payload skipAuthPayload
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&payload))

	assert.Equal(t, "test-sub", payload.Sub)
	assert.Equal(t, "test-user", payload.Username)
	assert.Equal(t, "read:edu", payload.Scope)
	assert.Equal(t, []string{"admins", "qa", "reporting"}, payload.Groups)
}

func TestSubjectMiddleware_Logging(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	app := fiber.New()
	app.Use(SubjectMiddleware(logger))
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString(c.Locals("sub").(string)) //nolint:forcetypeassert // test setup guarantees type
	})

	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	req.Header.Set("Authorization", "Bearer invalid-token-for-logging-test")

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	assert.Contains(t, buf.String(), "failed to parse bearer token")
	assert.Contains(t, buf.String(), "invalid-token-for-logging-test")
}
