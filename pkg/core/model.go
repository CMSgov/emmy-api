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

type VAConfig struct {
	BaseURL        string
	TokenURL       string
	ClientID       string
	TokenAudience  string
	PrivateKeyPath string
	TimeoutSeconds int
}

type Config struct {
	Environment string
	Otel        OtelConfig
	Port        int
	SkipAuth    bool
	Redis       RedisConfig
	NSC         NSCConfig
	VA          VAConfig
}

type ctxKey int

const (
	RequestContextKey ctxKey = iota
)
