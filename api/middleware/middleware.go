package middleware

import (
	"errors"
	"strings"
	"sync"

	"github.com/DSACMS/verification-service-api/pkg/circuitbreaker"
	"github.com/gofiber/fiber/v2"
)

const (
	skipAuthHeaderSub      = "x-skip-auth-sub"
	skipAuthHeaderUsername = "x-skip-auth-username"
	skipAuthHeaderScope    = "x-skip-auth-scope"
	skipAuthHeaderGroups   = "x-skip-auth-groups"

	defaultSkipAuthSub   = "skip-auth-user"
	defaultSkipAuthScope = "local"
	defaultSkipAuthGroup = "local-dev"
)

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
