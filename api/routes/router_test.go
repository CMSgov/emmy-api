package routes

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/DSACMS/verification-service-api/pkg/core"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestRegisterRoutes_RegistersOpenAPISpecEndpoint(t *testing.T) {
	app := fiber.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	RegisterRoutes(app, &core.Config{}, rdb, logger)

	req := httptest.NewRequest(http.MethodGet, "/api-spec/v1/verify", http.NoBody)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	expectedBody, err := os.ReadFile(filepath.Join("..", "..", "api-spec", "dist", "openapi.bundled.json"))
	require.NoError(t, err)

	require.Equal(t, fiber.StatusOK, resp.StatusCode)
	require.Contains(t, resp.Header.Get(fiber.HeaderContentType), fiber.MIMEApplicationJSON)
	require.JSONEq(t, string(expectedBody), string(body))
}
