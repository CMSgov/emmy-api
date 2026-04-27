package core

import (
	"errors"
	"fmt"
)

const (
	defaultConfigEnvironment string = "development"
	defaultConfigPort        int    = 3000
	defaultSkipAuth          bool   = false
	defaultVATimeoutSeconds  int    = 5
	defaultRedisAddr         string = "localhost:6379"
	defaultRedisPassword     string = ""
	defaultRedisDB           int    = 0
	defaultRedisUseTLS       bool   = true
	defaultRedisInsecureSkip bool   = false

	keyNSCSubmitURL   string = "NSC_SUBMIT_URL"
	keyTokenURL       string = "NSC_TOKEN_URL"     //nolint:gosec // Environment variable key name, not a credential.
	keyClientSecret   string = "NSC_CLIENT_SECRET" //nolint:gosec // Environment variable key name, not a credential.
	keyClientID       string = "NSC_CLIENT_ID"
	keyAccountID      string = "NSC_ACCOUNT_ID"
	keyVABaseURL      string = "VA_BASE_URL"
	keyVATokenURL     string = "VA_TOKEN_URL" //nolint:gosec // Environment variable key name, not a credential.
	keyVAClientID     string = "VA_CLIENT_ID"
	keyVAAudience     string = "VA_AUD"
	keyVAKeyPath      string = "VA_PRIVATE_KEY_PATH"
	keyVATimeout      string = "VA_TIMEOUT_SECONDS"
	keySQSQueueURL    string = "SQS_QUEUE_URL"
	keyServiceVersion string = "SERVICE_VERSION"
)

func DefaultConfig() Config {
	return Config{
		Environment: defaultConfigEnvironment,
		Port:        defaultConfigPort,
		SkipAuth:    defaultSkipAuth,

		Redis: RedisConfig{
			Addr:               defaultRedisAddr,
			Password:           defaultRedisPassword,
			DB:                 defaultRedisDB,
			UseTLS:             defaultRedisUseTLS,
			InsecureSkipVerify: defaultRedisInsecureSkip,
		},

		NSC: NSCConfig{
			SubmitURL:    getEnv(keyNSCSubmitURL, ""),
			TokenURL:     getEnv(keyTokenURL, ""),
			ClientSecret: getEnv(keyClientSecret, ""),
			ClientID:     getEnv(keyClientID, ""),
			AccountID:    getEnv(keyAccountID, ""),
		},

		VA: VAConfig{
			BaseURL:        getEnv(keyVABaseURL, ""),
			TokenURL:       getEnv(keyVATokenURL, ""),
			ClientID:       getEnv(keyVAClientID, ""),
			TokenAudience:  getEnv(keyVAAudience, ""),
			PrivateKeyPath: getEnv(keyVAKeyPath, ""),
			TimeoutSeconds: defaultVATimeoutSeconds,
		},
		Reporting: ReportingConfig{
			SQSQueueURL: getEnv(keySQSQueueURL, ""),
		},
		ServiceVersion: getEnv(keyServiceVersion, "1.3.0"),
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
// - REDIS_ADDR, REDIS_PASSWORD, REDIS_DB
//
// - NSC_SUBMIT_URL, NSC_TOKEN_URL, NSC_CLIENT_SECRET, NSC_CLIENT_ID, NSC_ACCOUNT_ID
//
// - VA_BASE_URL, VA_TOKEN_URL, VA_CLIENT_ID, VA_AUD, VA_PRIVATE_KEY_PATH, VA_TIMEOUT_SECONDS
//
// - SQS_QUEUE_URL, SERVICE_VERSION
//
// Provided options are applied after env loading and override both defaults and env file values.
func NewConfigFromEnv(options ...func(*Config)) (Config, error) {
	cfg := DefaultConfig()

	err := errors.Join(
		setFromEnv(&cfg.Environment, "ENVIRONMENT"),
		setFromEnv(&cfg.Port, "PORT"),
		setFromEnv(&cfg.SkipAuth, "SKIP_AUTH"),

		setFromEnv(&cfg.Redis.Addr, "REDIS_ADDR"),
		setFromEnv(&cfg.Redis.Password, "REDIS_PASSWORD"),
		setFromEnv(&cfg.Redis.DB, "REDIS_DB"),
		setFromEnv(&cfg.Redis.UseTLS, "REDIS_USE_TLS"),
		setFromEnv(&cfg.Redis.InsecureSkipVerify, "REDIS_INSECURE_SKIP_VERIFY"),

		setFromEnv(&cfg.NSC.SubmitURL, "NSC_SUBMIT_URL"),
		setFromEnv(&cfg.NSC.TokenURL, "NSC_TOKEN_URL"),
		setFromEnv(&cfg.NSC.ClientSecret, "NSC_CLIENT_SECRET"),
		setFromEnv(&cfg.NSC.ClientID, "NSC_CLIENT_ID"),
		setFromEnv(&cfg.NSC.AccountID, "NSC_ACCOUNT_ID"),

		setFromEnv(&cfg.VA.BaseURL, "VA_BASE_URL"),
		setFromEnv(&cfg.VA.TokenURL, "VA_TOKEN_URL"),
		setFromEnv(&cfg.VA.ClientID, "VA_CLIENT_ID"),
		setFromEnv(&cfg.VA.TokenAudience, "VA_AUD"),
		setFromEnv(&cfg.VA.PrivateKeyPath, "VA_PRIVATE_KEY_PATH"),
		setFromEnv(&cfg.VA.TimeoutSeconds, "VA_TIMEOUT_SECONDS"),

		setFromEnv(&cfg.Reporting.SQSQueueURL, "SQS_QUEUE_URL"),
		setFromEnv(&cfg.ServiceVersion, "SERVICE_VERSION"),
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
