package handlers

import (
	"github.com/diogomcd/fake-mill-api/internal/middleware"
	"github.com/diogomcd/fake-mill-api/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// BankAccountHandler handles requests to the /api/v1/bank-account endpoint
// @Summary Gera dados bancários fictícios
// @Description Gera dados de uma conta bancária fictícia, incluindo banco, agência e número da conta.
// @Tags Financeiro
// @Accept json
// @Produce json
// @Param quantity query int false "Quantidade de contas (1-200)" minimum(1) maximum(200) default(1)
// @Param bank query string false "Código do banco (ex: 001, 237)"
// @Success 200 {object} models.BankAccountResponse
// @Success 200 {array} models.BankAccountResponse
// @Router /bank-account [get]
func BankAccountHandler(c *fiber.Ctx) error {
	gen := middleware.GetGenerator(c)
	bankCode := c.Query("bank", "")

	log.Debug().
		Str("handler", "BankAccountHandler").
		Str("bank_code", bankCode).
		Msg("Bank account generation requested")

	return generateMultiple(c, func() models.BankAccountResponse {
		bank, agency, account, accountType := gen.GenerateBankAccount(bankCode)
		return models.BankAccountResponse{
			Bank: models.Bank{
				Code: bank.Code,
				Name: bank.Name,
			},
			Agency:      agency,
			Account:     account,
			AccountType: accountType,
		}
	})
}

// CreditCardHandler handles requests to the /api/v1/credit-card endpoint
// @Summary Gera dados de cartão de crédito fictício
// @Description Gera dados de um cartão de crédito fictício, incluindo número, bandeira, CVV e data de validade.
// @Tags Financeiro
// @Accept json
// @Produce json
// @Param quantity query int false "Quantidade de cartões (1-200)" minimum(1) maximum(200) default(1)
// @Param brand query string false "Bandeira do cartão" Enums(visa, mastercard, elo, amex)
// @Success 200 {object} models.CreditCardResponse
// @Success 200 {array} models.CreditCardResponse
// @Router /credit-card [get]
func CreditCardHandler(c *fiber.Ctx) error {
	gen := middleware.GetGenerator(c)
	brand := c.Query("brand", "")

	log.Debug().
		Str("handler", "CreditCardHandler").
		Str("brand", brand).
		Msg("Credit card generation requested")

	return generateMultiple(c, func() models.CreditCardResponse {
		number, cardBrand, cvv, expirationDate, holderName := gen.GenerateCreditCard(brand)
		return models.CreditCardResponse{
			Number:         number,
			Brand:          cardBrand,
			CVV:            cvv,
			ExpirationDate: expirationDate,
			HolderName:     holderName,
		}
	})
}
