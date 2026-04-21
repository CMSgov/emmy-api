package routes

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
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

	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	StatusRouter(app, cfg, rdb, logger)

	req := httptest.NewRequest(http.MethodGet, "/health", http.NoBody)

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
