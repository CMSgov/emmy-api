package core

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

type ReportingConfig struct {
	SQSQueueURL string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
	IAMAuth  bool
}

type Config struct {
	Environment string
	Port        int
	SkipAuth    bool
	ServiceVersion string
	Redis       RedisConfig
	Database    DatabaseConfig
	NSC         NSCConfig
	VA          VAConfig
	Reporting   ReportingConfig
}

type ctxKey int

const (
	RequestContextKey ctxKey = iota
)
