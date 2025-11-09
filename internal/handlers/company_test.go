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

func setupCompanyApp() *fiber.App {
	ds, err := generators.NewDataStore()
	if err != nil {
		panic(fmt.Sprintf("Failed to create DataStore: %v", err))
	}
	gen := generators.NewGenerator(ds)

	app := fiber.New()
	app.Use(middleware.InjectGenerator(gen))
	v1 := app.Group("/api/v1")
	v1.Get("/company", CompanyHandler)

	return app
}

func TestCompanyHandler_Success(t *testing.T) {
	app := setupCompanyApp()

	req := httptest.NewRequest("GET", "/api/v1/company", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	testutils.AssertJSONHeaders(t, resp)

	body, _ := io.ReadAll(resp.Body)
	var company models.CompanyResponse
	err = json.Unmarshal(body, &company)
	assert.NoError(t, err)
	assert.NotEmpty(t, company.Name)
	assert.NotEmpty(t, company.TradeName)
	assert.NotEmpty(t, company.CNPJ)
	assert.NotEmpty(t, company.Email)
	assert.NotEmpty(t, company.Phone)
	assert.NotEmpty(t, company.Address.Street)
}
