package handlers

import (
	"github.com/diogomcd/fake-mill-api/internal/generators"
	"github.com/diogomcd/fake-mill-api/internal/middleware"
	"github.com/diogomcd/fake-mill-api/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// EmailHandler handles requests to the /api/v1/email endpoint
// @Summary Gera email fictício
// @Description Gera um ou mais emails fictícios, com opção de domínio customizado.
// @Tags Contato
// @Accept json
// @Produce json
// @Param quantity query int false "Quantidade de emails (1-200)" minimum(1) maximum(200) default(1)
// @Param domain query string false "Domínio customizado (ex: minhaempresa.com)"
// @Success 200 {object} models.EmailResponse
// @Success 200 {array} models.EmailResponse
// @Router /email [get]
func EmailHandler(c *fiber.Ctx) error {
	gen := middleware.GetGenerator(c)
	domain := c.Query("domain", "")

	log.Debug().
		Str("handler", "EmailHandler").
		Str("domain", domain).
		Msg("Email generation requested")

	return generateMultiple(c, func() models.EmailResponse {
		email, username, emailDomain := gen.GenerateEmail(domain)
		return models.EmailResponse{
			Email:    email,
			Username: username,
			Domain:   emailDomain,
		}
	})
}

// PhoneHandler handles requests to the /api/v1/phone endpoint
// @Summary Gera telefone brasileiro
// @Description Gera um ou mais números de telefone brasileiros, com DDD e tipo.
// @Tags Contato
// @Accept json
// @Produce json
// @Param quantity query int false "Quantidade de telefones (1-200)" minimum(1) maximum(200) default(1)
// @Param state query string false "UF do estado (ex: SP, RJ)"
// @Param type query string false "Tipo de telefone" Enums(mobile, landline, random) default(random)
// @Success 200 {object} models.PhoneResponse
// @Success 200 {array} models.PhoneResponse
// @Router /phone [get]
func PhoneHandler(c *fiber.Ctx) error {
	gen := middleware.GetGenerator(c)
	state := c.Query("state", "")
	phoneType := c.Query("type", "random")

	originalState := state
	state = gen.GetDataStore().ValidateAndSanitizeState(state)

	if originalState != "" && state == "" {
		log.Warn().
			Str("handler", "PhoneHandler").
			Str("requested_state", originalState).
			Str("error_type", "invalid_state_code").
			Msg("Invalid state code provided, using random state")
	}

	log.Debug().
		Str("handler", "PhoneHandler").
		Str("state", state).
		Str("type", phoneType).
		Msg("Phone generation requested")

	return generateMultiple(c, func() models.PhoneResponse {
		phone, ddd, phoneState, actualPhoneType := gen.GeneratePhone(state, phoneType)
		return models.PhoneResponse{
			Phone:       phone,
			Formatted:   phone,
			Unformatted: generators.UnformatPhone(phone),
			Type:        actualPhoneType,
			DDD:         ddd,
			State:       phoneState,
		}
	})
}
