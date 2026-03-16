package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadEnvFile(t *testing.T) {
	err := loadEnvFile(".env.example")

	require.NoErrorf(t, err, `There was an error loading ".env.example": %v`, err)
}

func TestLoadEnv(t *testing.T) {
	err := LoadEnv()

	require.NoErrorf(t, err, "LoadEnv should not return an error. Got %v", err)
}

func TestGetEnv_KeyValue(t *testing.T) {
	t.Setenv("xyz", "abc")

	result := getEnv("xyz", "development")

	expected := "abc"

	assert.Equalf(t, expected, result, `getEnv("xyz", "development) = %q; expected: %q`, result, expected)
}

func TestGetEnv_FallbackValue(t *testing.T) {
	t.Setenv("xyz", "")

	result := getEnv("xyz", "development")

	expected := "development"

	assert.Equalf(t, expected, result, `getEnv("xyz", "development") = %q; expected: %q`, result, expected)
}
