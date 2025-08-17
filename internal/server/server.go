package server

import (
	"time"

	"github.com/dimuthuh28/go-api-gateway/internal/loadbalancer"
	"github.com/dimuthuh28/go-api-gateway/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func NewServer(mw *middleware.RedisMiddleware, lb *loadbalancer.RoundRobin) *fiber.App {
	app := fiber.New()

	// Use middlewares
	app.Use(mw.RateLimit(100))          // 100 requests limit
	app.Use(mw.Cache(60 * time.Second)) // cache responses for 60 seconds

	// Example route
	app.Get("/api/service1", func(c *fiber.Ctx) error {
		backend := lb.Next()
		return c.SendString("Service 1 Response via " + backend)
	})

	return app
}
