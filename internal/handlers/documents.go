package handlers

import (
	"strconv"

	"github.com/diogomcd/fake-mill-api/internal/generators"
	"github.com/diogomcd/fake-mill-api/internal/middleware"
	"github.com/diogomcd/fake-mill-api/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// CPFHandler handles requests to the /api/v1/cpf endpoint
// @Summary Gera CPF válido ou inválido
// @Description Gera um ou mais números de CPF válidos ou inválidos, com opção de formatação.
// @Tags Documentos
// @Accept json
// @Produce json
// @Param quantity query int false "Quantidade de CPFs (1-200)" minimum(1) maximum(200) default(1)
// @Param formatted query bool false "Retorna formatado (XXX.XXX.XXX-XX)" default(true)
// @Param valid query bool false "Gera CPF válido com dígitos verificadores corretos" default(true)
// @Success 200 {object} models.CPFResponse
// @Success 200 {array} models.CPFResponse
// @Router /cpf [get]
func CPFHandler(c *fiber.Ctx) error {
	gen := middleware.GetGenerator(c)
	formatted, _ := strconv.ParseBool(c.Query("formatted", "true"))
	valid, _ := strconv.ParseBool(c.Query("valid", "true"))

	log.Debug().
		Str("handler", "CPFHandler").
		Bool("formatted", formatted).
		Bool("valid", valid).
		Msg("CPF generation requested")

	return generateMultiple(c, func() models.CPFResponse {
		return *generateCPFResponse(gen, formatted, valid)
	})
}

// CNPJHandler handles requests to the /api/v1/cnpj endpoint
// @Summary Gera CNPJ válido ou inválido
// @Description Gera um ou mais números de CNPJ válidos ou inválidos, com opção de formatação.
// @Tags Documentos
// @Accept json
// @Produce json
// @Param quantity query int false "Quantidade de CNPJs (1-200)" minimum(1) maximum(200) default(1)
// @Param formatted query bool false "Retorna formatado (XX.XXX.XXX/XXXX-XX)" default(true)
// @Param valid query bool false "Gera CNPJ válido com dígitos verificadores corretos" default(true)
// @Success 200 {object} models.CNPJResponse
// @Success 200 {array} models.CNPJResponse
// @Router /cnpj [get]
func CNPJHandler(c *fiber.Ctx) error {
	gen := middleware.GetGenerator(c)
	formatted, _ := strconv.ParseBool(c.Query("formatted", "true"))
	valid, _ := strconv.ParseBool(c.Query("valid", "true"))

	log.Debug().
		Str("handler", "CNPJHandler").
		Bool("formatted", formatted).
		Bool("valid", valid).
		Msg("CNPJ generation requested")

	return generateMultiple(c, func() models.CNPJResponse {
		return *generateCNPJResponse(gen, formatted, valid)
	})
}

// RGHandler handles requests to the /api/v1/rg endpoint
// @Summary Gera RG válido ou inválido
// @Description Gera um ou mais números de RG válidos ou inválidos.
// @Tags Documentos
// @Accept json
// @Produce json
// @Param quantity query int false "Quantidade de RGs (1-200)" minimum(1) maximum(200) default(1)
// @Param formatted query bool false "Retorna formatado (XX.XXX.XXX-X)" default(true)
// @Param state query string false "UF do estado (ex: SP, RJ)"
// @Param valid query bool false "Gera RG válido com dígito verificador correto" default(true)
// @Success 200 {object} models.RGResponse
// @Success 200 {array} models.RGResponse
// @Router /rg [get]
func RGHandler(c *fiber.Ctx) error {
	gen := middleware.GetGenerator(c)
	state := c.Query("state", "")
	formatted, _ := strconv.ParseBool(c.Query("formatted", "true"))
	valid, _ := strconv.ParseBool(c.Query("valid", "true"))

	originalState := state
	state = gen.GetDataStore().ValidateAndSanitizeState(state)

	if originalState != "" && state == "" {
		log.Warn().
			Str("handler", "RGHandler").
			Str("requested_state", originalState).
			Str("error_type", "invalid_state_code").
			Msg("Invalid state code provided, using random state")
	}

	log.Debug().
		Str("handler", "RGHandler").
		Str("state", state).
		Bool("formatted", formatted).
		Bool("valid", valid).
		Msg("RG generation requested")

	return generateMultiple(c, func() models.RGResponse {
		rg, rgState, issuer, issueDate, expirationDate := gen.GenerateRG(state, formatted, valid)
		return models.RGResponse{
			RG:             rg,
			State:          rgState,
			Issuer:         issuer,
			IssueDate:      issueDate,
			ExpirationDate: expirationDate,
		}
	})
}

// generateCPFResponse creates the response structure for CPF
func generateCPFResponse(gen generators.DocumentGenerator, formatted bool, valid bool) *models.CPFResponse {
	cpf := gen.GenerateCPF(formatted, valid)

	return &models.CPFResponse{
		CPF:   cpf,
		Valid: valid,
	}
}

// generateCNPJResponse creates the response structure for CNPJ
func generateCNPJResponse(gen generators.DocumentGenerator, formatted bool, valid bool) *models.CNPJResponse {
	cnpj := gen.GenerateCNPJ(formatted, valid)

	return &models.CNPJResponse{
		CNPJ:  cnpj,
		Valid: valid,
	}
}
