package core

func WithRedisAddr(addr string) func(*Config) {
	return func(c *Config) {
		c.Redis.Addr = addr
	}
}

func WithRedisPassword(pw string) func(*Config) {
	return func(c *Config) {
		c.Redis.Password = pw
	}
}

func WithRedisDB(db int) func(*Config) {
	return func(c *Config) {
		c.Redis.DB = db
	}
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

func WithVABaseURL(baseURL string) func(*Config) {
	return func(c *Config) {
		c.VA.BaseURL = baseURL
	}
}

func WithVATokenURL(tokenURL string) func(*Config) {
	return func(c *Config) {
		c.VA.TokenURL = tokenURL
	}
}

func WithVAClientID(clientID string) func(*Config) {
	return func(c *Config) {
		c.VA.ClientID = clientID
	}
}

func WithVATokenAudience(aud string) func(*Config) {
	return func(c *Config) {
		c.VA.TokenAudience = aud
	}
}

func WithVAPrivateKeyPath(path string) func(*Config) {
	return func(c *Config) {
		c.VA.PrivateKeyPath = path
	}
}

func WithVATimeoutSeconds(timeout int) func(*Config) {
	return func(c *Config) {
		c.VA.TimeoutSeconds = timeout
	}
}
