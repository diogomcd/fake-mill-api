package middleware

import (
	"github.com/diogomcd/fake-mill-api/internal/config"
	"github.com/gofiber/fiber/v2"
)

// CORS returns the CORS middleware with default configuration
func CORS() fiber.Handler {
	return CORSWithConfig(config.CORSConfig{
		AllowOrigins: "*",
		AllowMethods: "GET, OPTIONS, POST",
		AllowHeaders: "Content-Type, Authorization",
		MaxAge:       "86400",
	})
}

// CORSWithConfig returns the CORS middleware with custom configuration
func CORSWithConfig(cfg config.CORSConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", cfg.AllowOrigins)
		c.Set("Access-Control-Allow-Methods", cfg.AllowMethods)
		c.Set("Access-Control-Allow-Headers", cfg.AllowHeaders)
		c.Set("Access-Control-Max-Age", cfg.MaxAge)

		if c.Method() == fiber.MethodOptions {
			return c.SendStatus(fiber.StatusNoContent)
		}

		return c.Next()
	}
}
