package handlers

import (
	"context"
	"time"

	redisLocal "github.com/cmsgov/emmy-api/pkg/redis"
	"github.com/gofiber/fiber/v2"
	goredis "github.com/redis/go-redis/v9"
)

const (
	contextTimeout = 2
)

func GetRDBStatus(rdb *goredis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.Context(), contextTimeout*time.Second)
		defer cancel()

		err := redisLocal.Ping(ctx, rdb)
		if err != nil {
			return err
		}

		return c.SendStatus(fiber.StatusOK)
	}
}
