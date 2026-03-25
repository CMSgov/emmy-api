package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestOpenAPISpecHandler_ReturnsBundledJSON(t *testing.T) {
	app := fiber.New()
	app.Get("/api-spec/v1/verify", OpenAPISpecHandler())

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
	require.NotEmpty(t, body)
	require.JSONEq(t, string(expectedBody), string(body))
}

func TestOpenAPISpecHandler_ReturnsInternalServerErrorWhenArtifactMissing(t *testing.T) {
	app := fiber.New()
	missingPath := filepath.Join(t.TempDir(), "missing-openapi.json")
	app.Get("/api-spec/v1/verify", OpenAPISpecHandlerForPath(missingPath))

	req := httptest.NewRequest(http.MethodGet, "/api-spec/v1/verify", http.NoBody)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}
