package middleware

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

type RedisMiddleware struct {
	Rdb *redis.Client
}

func (rm *RedisMiddleware) RateLimit(limit int64) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ip := c.IP()
		ctx := context.Background()
		count, _ := rm.Rdb.Incr(ctx, ip).Result()

		if count > limit {
			return c.Status(fiber.StatusTooManyRequests).SendString("Too many requests")
		}
		return c.Next()
	}
}

func (rm *RedisMiddleware) Cache(ttl time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := context.Background()
		cached, err := rm.Rdb.Get(ctx, c.Path()).Result()
		if err == nil {
			return c.SendString(cached)
		}
		resp := "Response from backend"
		rm.Rdb.Set(ctx, c.Path(), resp, ttl)
		return c.SendString(resp)
	}
}
