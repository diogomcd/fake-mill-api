package handlers

import (
	"github.com/diogomcd/fake-mill-api/internal/middleware"
	"github.com/diogomcd/fake-mill-api/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// AddressHandler handles requests to the /api/v1/address endpoint
// @Summary Gera endereço brasileiro fictício
// @Description Gera um ou mais endereços brasileiros fictícios, com opção de estado.
// @Tags Endereço
// @Accept json
// @Produce json
// @Param quantity query int false "Quantidade de endereços (1-200)" minimum(1) maximum(200) default(1)
// @Param state query string false "UF do estado (ex: SP, RJ)"
// @Success 200 {object} models.Address
// @Success 200 {array} models.Address
// @Router /address [get]
func AddressHandler(c *fiber.Ctx) error {
	gen := middleware.GetGenerator(c)
	state := c.Query("state", "")
	city := c.Query("city", "")

	originalState := state
	state = gen.GetDataStore().ValidateAndSanitizeState(state)

	if originalState != "" && state == "" {
		log.Warn().
			Str("handler", "AddressHandler").
			Str("requested_state", originalState).
			Str("error_type", "invalid_state_code").
			Msg("Invalid state code provided, using random state")
	}

	log.Debug().
		Str("handler", "AddressHandler").
		Str("state", state).
		Str("city", city).
		Msg("Address generation requested")

	return generateMultiple(c, func() models.Address {
		return *gen.GenerateAddress(state, city)
	})
}

// ZipcodeHandler handles requests to the /api/v1/zipcode endpoint
// @Summary Gera CEP válido
// @Description Gera um ou mais CEPs válidos, com opção de estado.
// @Tags Endereço
// @Accept json
// @Produce json
// @Param quantity query int false "Quantidade de CEPs (1-200)" minimum(1) maximum(200) default(1)
// @Param state query string false "UF do estado (ex: SP, RJ)"
// @Success 200 {object} models.ZipcodeResponse
// @Success 200 {array} models.ZipcodeResponse
// @Router /zipcode [get]
func ZipcodeHandler(c *fiber.Ctx) error {
	gen := middleware.GetGenerator(c)
	state := c.Query("state", "")

	originalState := state
	state = gen.GetDataStore().ValidateAndSanitizeState(state)

	if originalState != "" && state == "" {
		log.Warn().
			Str("handler", "ZipcodeHandler").
			Str("requested_state", originalState).
			Str("error_type", "invalid_state_code").
			Msg("Invalid state code provided, using random state")
	}

	log.Debug().
		Str("handler", "ZipcodeHandler").
		Str("state", state).
		Msg("Zipcode generation requested")

	return generateMultiple(c, func() models.ZipcodeResponse {
		formatted, unformatted, state, city := gen.GenerateZipcodeDetails(state)
		return models.ZipcodeResponse{
			Zipcode:     formatted,
			Formatted:   formatted,
			Unformatted: unformatted,
			State:       state,
			City:        city,
		}
	})
}
