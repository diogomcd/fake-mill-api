package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// Logger returns structured logging middleware
func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()
		duration := time.Since(start)

		// Request data
		method := c.Method()
		path := c.Path()
		statusCode := c.Response().StatusCode()
		ip := getClientIP(c)
		userAgent := c.Get("User-Agent")
		queryParams := c.Query("")
		contentLength := c.Request().Header.ContentLength()

		// Base log event with common fields
		baseLog := log.Info().
			Str("method", method).
			Str("path", path).
			Int("status", statusCode).
			Dur("duration_ms", duration).
			Str("ip", ip).
			Str("user_agent", userAgent).
			Int64("content_length", int64(contentLength))

		if queryParams != "" {
			baseLog = baseLog.Str("query_params", queryParams)
		}

		// If there's an error, logs as error
		if err != nil {
			log.Error().
				Err(err).
				Str("method", method).
				Str("path", path).
				Int("status", statusCode).
				Dur("duration_ms", duration).
				Str("ip", ip).
				Str("error_type", "handler_error").
				Msg("request handler failed")
			return err
		}

		// Log response based on status code
		if statusCode >= 500 {
			baseLog.Str("severity", "error").Msg("server error occurred")
		} else if statusCode >= 400 {
			baseLog.Str("severity", "warning").Msg("client error received")
		} else if statusCode >= 300 {
			baseLog.Str("severity", "info").Msg("redirect response sent")
		} else {
			baseLog.Str("severity", "debug").Msg("request completed successfully")
		}

		return err
	}
}
