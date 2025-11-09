package handlers

import (
	"strconv"

	"github.com/diogomcd/fake-mill-api/internal/middleware"
	"github.com/diogomcd/fake-mill-api/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

const (
	minQuantity = 1
	maxQuantity = 200
)

func parseQuantity(quantityStr string) int {
	quantity, err := strconv.Atoi(quantityStr)
	if err != nil {
		logger.Get().Warn().
			Err(err).
			Str("input", quantityStr).
			Str("reason", "invalid_integer_format").
			Msg("Failed to parse quantity parameter, using default value")
		return minQuantity
	}

	if quantity < minQuantity || quantity > maxQuantity {
		logger.Get().Warn().
			Int("requested", quantity).
			Int("min", minQuantity).
			Int("max", maxQuantity).
			Str("reason", "quantity_out_of_range").
			Msg("Quantity parameter out of range, using default value")
		return minQuantity
	}

	return quantity
}

func generateMultiple[T any](c *fiber.Ctx, generator func() T) error {
	quantity := parseQuantity(c.Query("quantity", "1"))

	if quantity == 1 {
		data := generator()

		if err := middleware.ValidateStruct(data); err != nil {
			log.Error().
				Err(err).
				Str("handler", "generateMultiple").
				Str("path", c.Path()).
				Str("error_type", "single_response_validation_failed").
				Msg("Single response validation failed")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
				"code":  "validation_failed",
			})
		}

		return c.JSON(data)
	}

	results := make([]T, quantity)
	for i := 0; i < quantity; i++ {
		data := generator()

		if err := middleware.ValidateStruct(data); err != nil {
			log.Error().
				Err(err).
				Int("index", i).
				Int("quantity", quantity).
				Str("handler", "generateMultiple").
				Str("path", c.Path()).
				Str("error_type", "batch_item_validation_failed").
				Msg("Batch item validation failed")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
				"code":  "batch_validation_failed",
				"index": i,
			})
		}

		results[i] = data
	}
	return c.JSON(results)
}
