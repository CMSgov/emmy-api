package core

type OtlpConfig struct {
	Endpoint string
	Insecure bool
}

type OtelConfig struct {
	OtlpExporter OtlpConfig
	Disable      bool
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type NSCConfig struct {
	SubmitURL    string
	TokenURL     string
	ClientSecret string
	ClientID     string
	AccountID    string
}

type Config struct {
	Environment string
	Otel        OtelConfig
	Port        int
	SkipAuth    bool
	Redis       RedisConfig
	NSC         NSCConfig
}
