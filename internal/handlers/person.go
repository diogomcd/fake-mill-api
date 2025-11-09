package handlers

import (
	"github.com/diogomcd/fake-mill-api/internal/middleware"
	"github.com/diogomcd/fake-mill-api/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// PersonHandler handles requests to the /api/v1/person endpoint
// @Summary Gera dados de uma pessoa
// @Description Gera um ou mais registros de pessoas fictícias com nome, CPF, RG, data de nascimento, etc.
// @Tags Pessoa
// @Accept json
// @Produce json
// @Param quantity query int false "Quantidade de registros (1-200)" minimum(1) maximum(200) default(1)
// @Param gender query string false "Gênero da pessoa" Enums(male, female, random)
// @Param state query string false "UF do estado (ex: SP, RJ)"
// @Success 200 {object} models.Person
// @Success 200 {array} models.Person
// @Router /person [get]
func PersonHandler(c *fiber.Ctx) error {
	gen := middleware.GetGenerator(c)
	gender := c.Query("gender", "")
	state := c.Query("state", "")

	if gender != "" && gender != "male" && gender != "female" && gender != "random" {
		log.Warn().
			Str("handler", "PersonHandler").
			Str("requested_gender", gender).
			Str("error_type", "invalid_gender_value").
			Msg("Invalid gender parameter received, ignoring")
		gender = ""
	}

	if gender == "random" {
		gender = ""
	}

	originalState := state
	state = gen.GetDataStore().ValidateAndSanitizeState(state)

	if originalState != "" && state == "" {
		log.Warn().
			Str("handler", "PersonHandler").
			Str("requested_state", originalState).
			Str("error_type", "invalid_state_code").
			Msg("Invalid state code provided, using random state")
	}

	log.Debug().
		Str("handler", "PersonHandler").
		Str("gender", gender).
		Str("state", state).
		Msg("Person generation requested")

	return generateMultiple(c, func() models.Person {
		return *gen.GeneratePerson(gender, state)
	})
}
