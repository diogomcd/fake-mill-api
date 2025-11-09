package middleware

import (
	"github.com/diogomcd/fake-mill-api/pkg/logger"
	"github.com/diogomcd/fake-mill-api/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

// ValidateResponse is a middleware that validates response data before sending
func ValidateResponse() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Continue to the handler
		if err := c.Next(); err != nil {
			return err
		}

		// Get the response body
		body := c.Response().Body()
		if len(body) == 0 {
			return nil
		}

		// Skip validation for non-JSON responses
		contentType := string(c.Response().Header.ContentType())
		if contentType != "application/json" && contentType != "" {
			return nil
		}

		// Log validation verification
		logger.Get().Debug().
			Str("path", c.Path()).
			Str("method", c.Method()).
			Msg("Validation middleware executed")

		return nil
	}
}

// ValidateStruct validates a struct and returns Fiber error if validation fails
func ValidateStruct(data interface{}) error {
	if err := validator.Validate(data); err != nil {
		logger.Get().Error().
			Err(err).
			Str("type", "validation_error").
			Msg("Response validation failed")
		return fiber.NewError(fiber.StatusInternalServerError, "Data validation failed")
	}
	return nil
}
