package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/DSACMS/verification-service-api/pkg/circuitbreaker"
	"github.com/gofiber/fiber/v2"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type CognitoConfig struct {
	Region     string
	UserPoolID string
	ClientID   string
}

type CognitoVerifier struct {
	issuer  string
	jwksURL string
	cache   *jwk.Cache
	client  *http.Client
	cfg     CognitoConfig
}

const (
	skipAuthHeaderSub      = "x-skip-auth-sub"
	skipAuthHeaderUsername = "x-skip-auth-username"
	skipAuthHeaderScope    = "x-skip-auth-scope"
	skipAuthHeaderGroups   = "x-skip-auth-groups"

	defaultSkipAuthSub   = "skip-auth-user"
	defaultSkipAuthScope = "local"
	defaultSkipAuthGroup = "local-dev"
)

func NewCognitoVerifier(cfg CognitoConfig) (*CognitoVerifier, error) {
	if cfg.Region == "" {
		return nil, errors.New("Region is required")
	}

	if cfg.UserPoolID == "" {
		return nil, errors.New("UserPoolID is required")
	}

	if cfg.ClientID == "" {
		return nil, errors.New("ClientID is required")
	}

	issuer := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", cfg.Region, cfg.UserPoolID)
	jwksURL := issuer + "/.well-known/jwks.json"

	cache := jwk.NewCache(context.Background())
	// register the JWKS URL with a refresh window
	cache.Register(jwksURL)

	return &CognitoVerifier{
		issuer:  issuer,
		jwksURL: jwksURL,
		cache:   cache,
		client:  &http.Client{Timeout: 5 * time.Second},
		cfg:     cfg,
	}, nil
}

// copy for testing
func NewCognitoVerifierWithURLs(cfg CognitoConfig, issuer, jwksURL string) (*CognitoVerifier, error) {
	if cfg.ClientID == "" {
		return nil, errors.New("ClientID is required")
	}

	if issuer == "" {
		return nil, errors.New("issuer is required")
	}

	if jwksURL == "" {
		return nil, errors.New("jwksURL is required")
	}

	cache := jwk.NewCache(context.Background())
	cache.Register(jwksURL)

	return &CognitoVerifier{
		issuer:  issuer,
		jwksURL: jwksURL,
		cache:   cache,
		client:  &http.Client{Timeout: 5 * time.Second},
		cfg:     cfg,
	}, nil
}

func (v *CognitoVerifier) FiberMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		raw := c.Get("x-amzn-oidc-accesstoken")
		if raw == "" {
			return fiber.ErrUnauthorized
		}

		ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
		defer cancel()

		keyset, err := v.cache.Get(ctx, v.jwksURL)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "unable to load jwks")
		}

		tok, err := jwt.Parse(
			[]byte(raw),
			jwt.WithKeySet(keyset),
			jwt.WithValidate(true),

			jwt.WithIssuer(v.issuer),

			jwt.WithClaimValue("token_use", "access"),
		)
		if err != nil {
			return fiber.ErrUnauthorized
		}

		// access tokens commonly carry client id in "client_id"
		// wrote return values + call in if statement since these are one-offs.
		// when written like this, the return values aren't accessible outside the scope of this if statement
		if cid, ok := tok.Get("client_id"); !ok || cid != v.cfg.ClientID {
			return fiber.ErrUnauthorized
		}

		// put useful info on context
		// save to c.locals to store temporary variables in the request's scope. They are only available to routes matching the request and go away when the request is handled
		if sub, ok := tok.Get("sub"); ok {
			c.Locals("sub", sub)
		}
		if username, ok := tok.Get("username"); ok {
			c.Locals("username", username)
		}
		if scope, ok := tok.Get("scope"); ok {
			c.Locals("scope", scope)
		}
		if groups, ok := tok.Get("cognito:groups"); ok {
			c.Locals("groups", groups)
		}

		return c.Next()
	}
}

// SkipAuthMiddleware injects a stable local identity when auth is disabled.
// Optional override headers:
//   - x-skip-auth-sub
//   - x-skip-auth-username
//   - x-skip-auth-scope
//   - x-skip-auth-groups (comma-separated)
func SkipAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		sub := c.Get(skipAuthHeaderSub)
		if sub == "" {
			sub = defaultSkipAuthSub
		}

		username := c.Get(skipAuthHeaderUsername)
		if username == "" {
			username = sub
		}

		scope := c.Get(skipAuthHeaderScope)
		if scope == "" {
			scope = defaultSkipAuthScope
		}

		groupsHeader := c.Get(skipAuthHeaderGroups)
		groups := []string{defaultSkipAuthGroup}
		if groupsHeader != "" {
			parsed := parseGroups(groupsHeader)
			if len(parsed) > 0 {
				groups = parsed
			}
		}

		c.Locals("sub", sub)
		c.Locals("username", username)
		c.Locals("scope", scope)
		c.Locals("groups", groups)

		return c.Next()
	}
}

func parseGroups(value string) []string {
	parts := strings.Split(value, ",")
	groups := make([]string, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		groups = append(groups, trimmed)
	}

	return groups
}

// wraps another handler and runs before c.Next
// It can block the entire request
// figures out what breaker to use since we're intentionally reusing one breaker per endpoint
// calls breaker.Allow() to either block the request of let it continue
func WithCircuitBreaker(newBreaker func(name string) *circuitbreaker.RedisBreaker) func(fiber.Handler) fiber.Handler {
	var mu sync.RWMutex

	// the lookup table is usable by multiple requests
	// We want multiple requests to be able to hit this middleware without failing or creating inconsistent state such as multiple requests creating multiple breakers for the same endpoint
	breakers := make(map[string]*circuitbreaker.RedisBreaker)

	getBreaker := func(name string) *circuitbreaker.RedisBreaker {
		// allows concurrent read requests but not write requests
		mu.RLock()
		// read breaker in map and assign to variable
		b := breakers[name]

		// release the read lock/unlock
		mu.RUnlock()
		// if the breaker exists, return it
		if b != nil {
			return b
		}

		// blocks goroutine requests
		mu.Lock()
		// release the mutex lock immediately after this function runs
		defer mu.Unlock()
		// now that mutext is locked, check if breaker exists
		b = breakers[name]
		// if breaker exists, return it
		if b != nil {
			return b
		}

		// if the breaker does not exist, use the newBreaker higher order function to create a breaker
		b = newBreaker(name)
		// now that a new breaker has been created, add it to the breakers table
		breakers[name] = b
		// return the breaker
		return b
	}

	return func(next fiber.Handler) fiber.Handler {
		return func(c *fiber.Ctx) error {
			name := breakerName(c)
			breaker := getBreaker(name)

			err := breaker.Allow(c.Context())
			if err != nil {
				if errors.Is(err, circuitbreaker.ErrCircuitOpen) {
					return c.SendStatus(fiber.StatusServiceUnavailable)
				}

				return c.SendStatus(fiber.StatusServiceUnavailable)
			}

			return next(c)
		}
	}
}

func breakerName(c *fiber.Ctx) string {
	var path string
	r := c.Route()
	if r != nil && r.Path != "" {
		path = r.Path
	} else {
		path = c.Path()
	}

	return c.Method() + " " + path

}
