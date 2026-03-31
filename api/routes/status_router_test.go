package routes

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cmsgov/emmy-api/pkg/core"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatusEndpoint(t *testing.T) {
	app := fiber.New()

	cfg := core.Config{}

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	StatusRouter(app, cfg, rdb, logger)

	req := httptest.NewRequest(http.MethodGet, "/status", http.NoBody)

	expected := fiber.StatusOK

	result, err := app.Test(req)

	require.NoErrorf(t, err, "app.Test(req) returned error: %v", err)
	defer result.Body.Close()

	assert.Equalf(
		t,
		expected,
		result.StatusCode,
		"app.Test(req) returned status %v; expected: %v",
		result.StatusCode,
		expected,
	)
}
