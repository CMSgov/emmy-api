package core

import (
	"errors"
	"fmt"
)

const (
	defaultConfigEnvironment    string = "development"
	defaultConfigPort           int    = 8000
	defaultSkipAuth             bool   = false
	defaultOtelDisable          bool   = false
	defaultOTLPExporterEndpoint string = "localhost:4317"
	defaultOTLPInsecure         bool   = false
	defaultCognitoRegion        string = "us-east-1"
	defaultCognitoUserPoolID    string = "UNSET"
	defaultCognitoAppClientID   string = "UNSET"
	defaultRedisAddr            string = "localhost:6379"
	defaultRedisPassword        string = ""
	defaultRedisDB              int    = 0

	keyNSCSubmitURL string = "NSC_SUBMIT_URL"
	keyTokenURL     string = "NSC_TOKEN_URL"
	keyClientSecret string = "NSC_CLIENT_SECRET"
	keyClientID     string = "NSC_CLIENT_ID"
	keyAccountID    string = "NSC_ACCOUNT_ID"
)

func DefaultConfig() Config {
	return Config{
		Environment: defaultConfigEnvironment,
		Port:        defaultConfigPort,
		SkipAuth:    defaultSkipAuth,

		Otel: OtelConfig{
			Disable: defaultOtelDisable,
			OtlpExporter: OtlpConfig{
				Endpoint: defaultOTLPExporterEndpoint,
				Insecure: defaultOTLPInsecure,
			},
		},

		Cognito: CognitoConfig{
			Region:      defaultCognitoRegion,
			UserPoolID:  defaultCognitoUserPoolID,
			AppClientID: defaultCognitoAppClientID,
		},

		Redis: RedisConfig{
			Addr:     defaultRedisAddr,
			Password: defaultRedisPassword,
			DB:       defaultRedisDB,
		},

		NSC: NSCConfig{
			SubmitURL:    getEnv(keyNSCSubmitURL, ""),
			TokenURL:     getEnv(keyTokenURL, ""),
			ClientSecret: getEnv(keyClientSecret, ""),
			ClientID:     getEnv(keyClientID, ""),
			AccountID:    getEnv(keyAccountID, ""),
		},
	}
}

func NewConfig(options ...func(*Config)) Config {
	cfg := DefaultConfig()
	for _, opt := range options {
		opt(&cfg)
	}
	return cfg
}

// NewConfigFromEnv builds Config from DefaultConfig and selected environment variables.
// Supported keys:
//
// - ENVIRONMENT, PORT, SKIP_AUTH
//
// - OTEL_DISABLE, OTEL_OTLP_EXPORTER_ENDPOINT, OTEL_OTLP_EXPORTER_INSECURE
//
// - COGNITO_REGION, COGNITO_USER_POOL_ID, COGNITO_APP_CLIENT_ID
//
// - REDIS_ADDR, REDIS_PASSWORD, REDIS_DB
//
// - NSC_SUBMIT_URL, NSC_TOKEN_URL, NSC_CLIENT_SECRET, NSC_CLIENT_ID, NSC_ACCOUNT_ID
//
// Provided options are applied after env loading and override both defaults and env file values.
func NewConfigFromEnv(options ...func(*Config)) (Config, error) {
	cfg := DefaultConfig()

	err := errors.Join(
		setFromEnv(&cfg.Environment, "ENVIRONMENT"),
		setFromEnv(&cfg.Port, "PORT"),
		setFromEnv(&cfg.SkipAuth, "SKIP_AUTH"),

		setFromEnv(&cfg.Otel.Disable, "OTEL_DISABLE"),
		setFromEnv(&cfg.Otel.OtlpExporter.Endpoint, "OTEL_OTLP_EXPORTER_ENDPOINT"),
		setFromEnv(&cfg.Otel.OtlpExporter.Insecure, "OTEL_OTLP_EXPORTER_INSECURE"),

		setFromEnv(&cfg.Cognito.Region, "COGNITO_REGION"),
		setFromEnv(&cfg.Cognito.UserPoolID, "COGNITO_USER_POOL_ID"),
		setFromEnv(&cfg.Cognito.AppClientID, "COGNITO_APP_CLIENT_ID"),

		setFromEnv(&cfg.Redis.Addr, "REDIS_ADDR"),
		setFromEnv(&cfg.Redis.Password, "REDIS_PASSWORD"),
		setFromEnv(&cfg.Redis.DB, "REDIS_DB"),

		setFromEnv(&cfg.NSC.SubmitURL, "NSC_SUBMIT_URL"),
		setFromEnv(&cfg.NSC.TokenURL, "NSC_TOKEN_URL"),
		setFromEnv(&cfg.NSC.ClientSecret, "NSC_CLIENT_SECRET"),
		setFromEnv(&cfg.NSC.ClientID, "NSC_CLIENT_ID"),
		setFromEnv(&cfg.NSC.AccountID, "NSC_ACCOUNT_ID"),
	)

	for _, opt := range options {
		opt(&cfg)
	}

	return cfg, err
}

func LoadEnv(environment ...string) error {
	filenames := []string{
		".env.local",
		".env",
	}

	env := getEnv("ENVIRONMENT", DefaultConfig().Environment)
	if len(environment) > 0 {
		env = environment[0]
	}

	if env != "" {
		file := ".env." + env + ".local"
		filenames = append([]string{file}, filenames...)
	}

	var errs error

	for _, filename := range filenames {
		if err := loadEnvFile(filename); err != nil {
			errs = errors.Join(
				errs,
				fmt.Errorf("error loading %s: %w", filename, err),
			)
		}
	}

	return errs
}
