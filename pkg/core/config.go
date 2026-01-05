package core

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type OtlpConfig struct {
	Endpoint string
	Insecure bool
}

type OtelConfig struct {
	OtlpExporter OtlpConfig
	Disable      bool
}

type CognitoConfig struct {
	Region      string
	UserPoolID  string
	AppClientID string
}

// Values needed to connect to the national student clearing house
// see:
//   - https://docs.studentclearinghouse.org/vs/insights-json/verification-services-request
//
// - https://docs.studentclearinghouse.org/vs/verification-services/oauth
type NSCConfig struct {
	// Endpoint for the insights API
	Endpoint string
	// Client ID for NSC API
	ClientID string
	// Client Secret for NSC API
	ClientSecret string
	// Access Token URL for NSC API
	AccessTokenURL string
}

type Config struct {
	Cognito     CognitoConfig
	Environment string
	Otel        OtelConfig
	Port        int
	SkipAuth    bool
	// Configuration for connecting to NSC
	NSCConfig NSCConfig
}

func WithEnvironment(environment string) func(*Config) {
	return func(c *Config) {
		c.Environment = environment
	}
}

func WithPort(port int) func(*Config) {
	return func(c *Config) {
		c.Port = port
	}
}

func WithSkipAuth(value ...bool) func(*Config) {
	val := true
	if len(value) > 0 {
		val = value[0]
	}

	return func(c *Config) {
		c.SkipAuth = val
	}
}

func WithOtlpEndpoint(endpoint string) func(*Config) {
	return func(c *Config) {
		c.Otel.OtlpExporter.Endpoint = endpoint
	}
}

func WithOtlpInsecure(insecure bool) func(*Config) {
	return func(c *Config) {
		c.Otel.OtlpExporter.Insecure = insecure
	}
}

func WithOtelDisable(value ...bool) func(*Config) {
	val := true
	if len(value) > 0 {
		val = value[0]
	}

	return func(c *Config) {
		c.Otel.Disable = val
	}
}

func WithCognitoRegion(region string) func(*Config) {
	return func(c *Config) {
		c.Cognito.Region = region
	}
}

func WithCognitoUserPoolID(userPoolID string) func(*Config) {
	return func(c *Config) {
		c.Cognito.UserPoolID = userPoolID
	}
}

func WithCognitoAppClientID(appClientID string) func(*Config) {
	return func(c *Config) {
		c.Cognito.AppClientID = appClientID
	}
}

func WithNSCEndpoint(endpoint string) func(*Config) {
	return func(c *Config) {
		c.NSCConfig.Endpoint = endpoint
	}
}

func WithNSCClientID(clientID string) func(*Config) {
	return func(c *Config) {
		c.NSCConfig.ClientID = clientID
	}
}

func WithNSCClientSecret(clientSecret string) func(*Config) {
	return func(c *Config) {
		c.NSCConfig.ClientID = clientSecret
	}
}

func WithNSCAccessTokenURL(accessTokenURL string) func(*Config) {
	return func(c *Config) {
		c.NSCConfig.ClientID = accessTokenURL
	}
}

func DefaultConfig() Config {
	return Config{
		Environment: "development",
		Port:        8000,
		SkipAuth:    false,
		Otel: OtelConfig{
			Disable: false,
			OtlpExporter: OtlpConfig{
				Endpoint: "localhost:4317",
				Insecure: false,
			},
		},
		Cognito: CognitoConfig{
			Region:      "us-east-1",
			UserPoolID:  "UNSET",
			AppClientID: "UNSET",
		},
		NSCConfig: NSCConfig{
			Endpoint:       "UNSET",
			ClientID:       "UNSET",
			ClientSecret:   "UNSET",
			AccessTokenURL: "UNSET",
		},
	}
}

func NewConfig(options ...func(*Config)) Config {
	config := DefaultConfig()
	for _, opt := range options {
		opt(&config)
	}
	return config
}

func setFromEnv(loc any, key string) error {
	strValue := os.Getenv(key)
	if strValue == "" {
		return nil
	}

	switch v := loc.(type) {
	case *string:
		(*v) = strValue
	case *bool:
		val, err := strconv.ParseBool(strValue)
		if err != nil {
			return fmt.Errorf("failed to parse %s as a bool: %w", strValue, err)
		}
		(*v) = val
	case *int:
		val, err := strconv.ParseInt(strValue, 10, strconv.IntSize)
		if err != nil {
			return fmt.Errorf("failed to parse %s as an int: %w", strValue, err)
		}
		(*v) = int(val)
	}
	return nil
}

func NewConfigFromEnv(options ...func(*Config)) (Config, error) {
	config := DefaultConfig()
	err := errors.Join(
		setFromEnv(&config.Environment, "ENVIRONMENT"),
		setFromEnv(&config.Port, "PORT"),
		setFromEnv(&config.SkipAuth, "SKIP_AUTH"),
		setFromEnv(&config.Otel.Disable, "OTEL_DISABLE"),
		setFromEnv(&config.Otel.OtlpExporter.Endpoint, "OTEL_OTLP_EXPORTER_ENDPOINT"),
		setFromEnv(&config.Otel.OtlpExporter.Insecure, "OTEL_OTLP_EXPORTER_INSECURE"),
		setFromEnv(&config.Cognito.Region, "COGNITO_REGION"),
		setFromEnv(&config.Cognito.UserPoolID, "COGNITO_USER_POOL_ID"),
		setFromEnv(&config.Cognito.AppClientID, "COGNITO_APP_CLIENT_ID"),
	)

	for _, opt := range options {
		opt(&config)
	}

	return config, err
}

func loadEnvFile(filename string) error {
	err := godotenv.Load(filename)
	if err == nil || errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return fmt.Errorf("error loading file %s: %w", filename, err)
}

func LoadEnv(environment ...string) error {
	var err error
	var errs error

	filenames := []string{
		".env.local",
		".env",
	}

	env := getEnv("ENVIRONMENT", DefaultConfig().Environment)
	if len(environment) > 0 {
		env = environment[0]
	}
	if env != "" {
		file := ".env" + env + ".local"
		filenames = append([]string{file}, filenames...)
	}

	for _, filename := range filenames {
		err = loadEnvFile(filename)
		if err != nil {
			errs = errors.Join(
				errs,
				fmt.Errorf("error loading %s: %w", filename, err),
			)
		}
	}

	return errs
}

func getEnv(key, fallback string) string {
	// returns value of associated env key
	value := os.Getenv(key)

	if value != "" {
		return value
	}

	return fallback
}

func (c *Config) IsProd() bool {
	return c.Environment == "production"
}
