package core

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

type RedisConfig struct {
	Addr               string
	Password           string
	DB                 int
	UseTLS             bool
	InsecureSkipVerify bool
}

type NSCConfig struct {
	SubmitURL    string
	TokenURL     string
	ClientSecret string
	ClientID     string
	AccountID    string
}

type Config struct {
	Cognito     CognitoConfig
	Environment string
	Otel        OtelConfig
	Port        int
	SkipAuth    bool
	Redis       RedisConfig
	NSC         NSCConfig
}
