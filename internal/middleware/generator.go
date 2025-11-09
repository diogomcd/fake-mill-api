package middleware

import (
	"github.com/diogomcd/fake-mill-api/internal/generators"
	"github.com/gofiber/fiber/v2"
)

const generatorKey = "generator"

// InjectGenerator creates a middleware to inject the Generator into the context
// Accepts IGenerator to facilitate mocking in tests
func InjectGenerator(gen generators.IGenerator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals(generatorKey, gen)
		return c.Next()
	}
}

// GetGenerator retrieves the Generator from the context as IGenerator
func GetGenerator(c *fiber.Ctx) generators.IGenerator {
	gen, ok := c.Locals(generatorKey).(generators.IGenerator)
	if !ok {
		panic("generator not found in context")
	}
	return gen
}
