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

type Config struct {
	NSC            NSCConfig
	Environment    string
	Reporting      ReportingConfig
	ServiceVersion string
	VA             VAConfig
	Redis          RedisConfig
	Port           int
	SkipAuth       bool
}

type ctxKey int

const (
	RequestContextKey ctxKey = iota
)
