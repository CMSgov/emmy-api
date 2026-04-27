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
	Environment    string
	Port           int
	SkipAuth       bool
	Redis          RedisConfig
	NSC            NSCConfig
	VA             VAConfig
	Reporting      ReportingConfig
	ServiceVersion string
}

type ctxKey int

const (
	RequestContextKey ctxKey = iota
)
