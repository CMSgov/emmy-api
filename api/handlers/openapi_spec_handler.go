package handlers

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

const bundledOpenAPISpecPath = "api-spec/dist/openapi.bundled.json"

func OpenAPISpecHandler() fiber.Handler {
	return OpenAPISpecHandlerForPath(bundledOpenAPISpecPath)
}

func OpenAPISpecHandlerForPath(path string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		resolvedPath, err := resolveOpenAPISpecPath(path)
		if err != nil {
			return fiber.NewError(
				fiber.StatusInternalServerError,
				fmt.Sprintf("openapi spec unavailable: %v", err),
			)
		}

		body, err := os.ReadFile(resolvedPath)
		if err != nil {
			return fiber.NewError(
				fiber.StatusInternalServerError,
				fmt.Sprintf("openapi spec unavailable: %v", err),
			)
		}

		c.Type("json")
		return c.Status(fiber.StatusOK).Send(body)
	}
}

func resolveOpenAPISpecPath(path string) (string, error) {
	if filepath.IsAbs(path) {
		if _, err := os.Stat(path); err != nil {
			return "", err
		}
		return path, nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for dir := wd; ; dir = filepath.Dir(dir) {
		candidate := filepath.Join(dir, path)
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
	}

	return "", os.ErrNotExist
}
