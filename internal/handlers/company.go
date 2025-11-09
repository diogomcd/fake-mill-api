package handlers

import (
	"github.com/diogomcd/fake-mill-api/internal/middleware"
	"github.com/diogomcd/fake-mill-api/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// CompanyHandler handles requests to the /api/v1/company endpoint
// @Summary Gera dados completos de empresa fictícia
// @Description Gera dados completos de uma ou mais empresas fictícias, incluindo nome, CNPJ, endereço, etc.
// @Tags Empresa
// @Accept json
// @Produce json
// @Param quantity query int false "Quantidade de empresas (1-200)" minimum(1) maximum(200) default(1)
// @Success 200 {object} models.CompanyResponse
// @Success 200 {array} models.CompanyResponse
// @Router /company [get]
func CompanyHandler(c *fiber.Ctx) error {
	gen := middleware.GetGenerator(c)

	log.Debug().
		Str("handler", "CompanyHandler").
		Msg("Company generation requested")

	return generateMultiple(c, func() models.CompanyResponse {
		return *gen.GenerateCompany()
	})
}
