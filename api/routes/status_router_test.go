package routes

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatusEndpoint(t *testing.T) {
	app := fiber.New()

	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	logger := slog.New(slog.DiscardHandler)

	StatusRouter(app, rdb, logger)

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
