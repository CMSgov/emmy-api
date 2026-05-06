package routes

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/cmsgov/emmy-api/pkg/core"
	"github.com/cmsgov/emmy-api/pkg/encryption"
	"github.com/cmsgov/emmy-api/pkg/reporting"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func getRedisAddr() string {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		return "localhost:6379"
	}
	return addr
}

func TestRegisterRoutes_RegistersOpenAPISpecEndpoint(t *testing.T) {
	app := fiber.New()
	logger := slog.New(slog.DiscardHandler)
	rdb := redis.NewClient(&redis.Options{
		Addr: getRedisAddr(),
	})
	reporter := reporting.NewMockReporter()
	encrypt := encryption.NewMockEncryptionService()

	RegisterRoutes(app, RouterParams{
		CFG:        &core.Config{},
		RDB:        rdb,
		Reporter:   reporter,
		Logger:     logger,
		Encryption: encrypt,
	})

	req := httptest.NewRequest(http.MethodGet, "/api-spec/v1/verify", http.NoBody)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	expectedBody, err := os.ReadFile(filepath.Join("..", "..", "api-spec", "v0", "dist", "openapi.bundled.json"))
	require.NoError(t, err)

	require.Equal(t, fiber.StatusOK, resp.StatusCode)
	require.Contains(t, resp.Header.Get(fiber.HeaderContentType), fiber.MIMEApplicationJSON)
	require.JSONEq(t, string(expectedBody), string(body))
}

func TestRegisterRoutes_RegistersVeteranDisabilityEndpoint(t *testing.T) {
	app := fiber.New()
	logger := slog.New(slog.DiscardHandler)
	rdb := redis.NewClient(&redis.Options{
		Addr: getRedisAddr(),
	})
	reporter := reporting.NewMockReporter()
	encrypt := encryption.NewMockEncryptionService()

	RegisterRoutes(app, RouterParams{
		CFG:        &core.Config{},
		RDB:        rdb,
		Reporter:   reporter,
		Logger:     logger,
		Encryption: encrypt,
	})

	routes := app.GetRoutes(true)
	for _, route := range routes {
		if route.Method == http.MethodPost && route.Path == "/api/v0/veteran-disability-ratings" {
			return
		}
	}

	t.Fatal("expected POST /api/v0/veteran-disability-ratings to be registered")
}

func TestRegisterRoutes_RegistersEducationEnrollmentsEndpoint(t *testing.T) {
	app := fiber.New()
	logger := slog.New(slog.DiscardHandler)
	rdb := redis.NewClient(&redis.Options{
		Addr: getRedisAddr(),
	})
	reporter := reporting.NewMockReporter()
	encrypt := encryption.NewMockEncryptionService()

	RegisterRoutes(app, RouterParams{
		CFG:        &core.Config{},
		RDB:        rdb,
		Reporter:   reporter,
		Logger:     logger,
		Encryption: encrypt,
	})

	routes := app.GetRoutes(true)
	for _, route := range routes {
		if route.Method == http.MethodPost && route.Path == "/api/v0/education-enrollments" {
			return
		}
	}

	t.Fatal("expected POST /api/v0/education-enrollments to be registered")
}
