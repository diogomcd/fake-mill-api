package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/diogomcd/fake-mill-api/internal/generators"
	"github.com/diogomcd/fake-mill-api/internal/handlers/testutils"
	"github.com/diogomcd/fake-mill-api/internal/middleware"
	"github.com/diogomcd/fake-mill-api/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupFinancialApp() *fiber.App {
	ds, err := generators.NewDataStore()
	if err != nil {
		panic(fmt.Sprintf("Failed to create DataStore: %v", err))
	}
	gen := generators.NewGenerator(ds)

	app := fiber.New()
	app.Use(middleware.InjectGenerator(gen))
	v1 := app.Group("/api/v1")
	v1.Get("/bank-account", BankAccountHandler)
	v1.Get("/credit-card", CreditCardHandler)

	return app
}

func TestBankAccountHandler_Success(t *testing.T) {
	app := setupFinancialApp()

	req := httptest.NewRequest("GET", "/api/v1/bank-account", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	testutils.AssertJSONHeaders(t, resp)

	body, _ := io.ReadAll(resp.Body)
	var bankResp models.BankAccountResponse
	err = json.Unmarshal(body, &bankResp)
	assert.NoError(t, err)
	assert.NotEmpty(t, bankResp.Bank.Code)
	assert.NotEmpty(t, bankResp.Bank.Name)
	assert.NotEmpty(t, bankResp.Agency)
	assert.NotEmpty(t, bankResp.Account)
	assert.Contains(t, []string{"checking", "savings"}, bankResp.AccountType)
}

func TestCreditCardHandler_Success(t *testing.T) {
	app := setupFinancialApp()

	req := httptest.NewRequest("GET", "/api/v1/credit-card", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var cardResp models.CreditCardResponse
	err = json.Unmarshal(body, &cardResp)
	assert.NoError(t, err)
	assert.NotEmpty(t, cardResp.Number)
	assert.NotEmpty(t, cardResp.Brand)
	assert.Contains(t, []string{"Visa", "Mastercard", "Elo", "Amex"}, cardResp.Brand)
	assert.Len(t, cardResp.CVV, 3)
	assert.Len(t, cardResp.ExpirationDate, 5)
	assert.NotEmpty(t, cardResp.HolderName)
}

func TestCreditCardHandler_BrandVisa(t *testing.T) {
	app := setupFinancialApp()

	req := httptest.NewRequest("GET", "/api/v1/credit-card?brand=Visa", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var cardResp models.CreditCardResponse
	err = json.Unmarshal(body, &cardResp)
	assert.NoError(t, err)
	assert.Equal(t, "Visa", cardResp.Brand)
}
