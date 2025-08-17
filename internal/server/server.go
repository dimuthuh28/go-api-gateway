package server

import (
	"api-gateway-go/internal/loadbalancer"
	"api-gateway-go/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func NewServer(mw *middleware.RedisMiddleware, lb *loadbalancer.RoundRobin) *fiber.App {
	app := fiber.New()

	// Use middlewares
	app.Use(mw.RateLimit(100))
	app.Use(mw.Cache())

	// Example route
	app.Get("/api/service1", func(c *fiber.Ctx) error {
		backend := lb.Next()
		return c.SendString("Service 1 Response via " + backend)
	})

	return app
}
