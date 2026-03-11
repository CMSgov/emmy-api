package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
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
