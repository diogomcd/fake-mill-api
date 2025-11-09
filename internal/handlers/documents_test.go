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

func setupDocumentsApp() *fiber.App {
	ds, err := generators.NewDataStore()
	if err != nil {
		panic(fmt.Sprintf("Failed to create DataStore: %v", err))
	}
	gen := generators.NewGenerator(ds)

	app := fiber.New()
	app.Use(middleware.InjectGenerator(gen))
	v1 := app.Group("/api/v1")
	v1.Get("/cpf", CPFHandler)
	v1.Get("/cnpj", CNPJHandler)
	v1.Get("/rg", RGHandler)
	v1.Get("/validate/cpf/:cpf", ValidateCPFHandler)
	v1.Get("/validate/cpf", ValidateCPFHandler)
	v1.Get("/validate/cnpj/:cnpj", ValidateCNPJHandler)
	v1.Get("/validate/rg/:rg", ValidateRGHandler)

	return app
}

func TestCPFHandler_Success(t *testing.T) {
	app := setupDocumentsApp()

	req := httptest.NewRequest("GET", "/api/v1/cpf", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	testutils.AssertJSONHeaders(t, resp)

	body, _ := io.ReadAll(resp.Body)
	var cpfResp models.CPFResponse
	err = json.Unmarshal(body, &cpfResp)
	assert.NoError(t, err)
	assert.NotEmpty(t, cpfResp.CPF)
	assert.True(t, cpfResp.Valid)
}

func TestCNPJHandler_Success(t *testing.T) {
	app := setupDocumentsApp()

	req := httptest.NewRequest("GET", "/api/v1/cnpj", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var cnpjResp models.CNPJResponse
	err = json.Unmarshal(body, &cnpjResp)
	assert.NoError(t, err)
	assert.NotEmpty(t, cnpjResp.CNPJ)
	assert.True(t, cnpjResp.Valid)
}

func TestRGHandler_Success(t *testing.T) {
	app := setupDocumentsApp()

	req := httptest.NewRequest("GET", "/api/v1/rg?state=SP", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var rgResp models.RGResponse
	err = json.Unmarshal(body, &rgResp)
	assert.NoError(t, err)
	assert.NotEmpty(t, rgResp.RG)
	assert.Equal(t, "SP", rgResp.State)
}

func TestValidateCPFHandler_Valid(t *testing.T) {
	app := setupDocumentsApp()

	gen := generators.NewGenerator(nil)
	validCPF := gen.GenerateCPF(false, true)

	req := httptest.NewRequest("GET", "/api/v1/validate/cpf/"+validCPF, nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var validation models.CPFValidationResponse
	err = json.Unmarshal(body, &validation)
	assert.NoError(t, err)
	assert.True(t, validation.Valid)
}

func TestValidateCPFHandler_Invalid(t *testing.T) {
	app := setupDocumentsApp()

	req := httptest.NewRequest("GET", "/api/v1/validate/cpf/00000000000", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var validation models.CPFValidationResponse
	err = json.Unmarshal(body, &validation)
	assert.NoError(t, err)
	assert.False(t, validation.Valid)
}

func TestValidateCPFHandler_MissingCPF(t *testing.T) {
	app := setupDocumentsApp()

	req := httptest.NewRequest("GET", "/api/v1/validate/cpf", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	testutils.AssertHTTPError(t, resp, 400, "cpf parameter is required")
}

func TestValidateCNPJHandler_Valid(t *testing.T) {
	app := setupDocumentsApp()

	gen := generators.NewGenerator(nil)
	validCNPJ := gen.GenerateCNPJ(false, true)

	req := httptest.NewRequest("GET", "/api/v1/validate/cnpj/"+validCNPJ, nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var validation models.CNPJValidationResponse
	err = json.Unmarshal(body, &validation)
	assert.NoError(t, err)
	assert.True(t, validation.Valid)
}

func TestValidateRGHandler_Valid(t *testing.T) {
	ds, _ := generators.NewDataStore()
	gen := generators.NewGenerator(ds)
	rg, _, _, _, _ := gen.GenerateRG("SP", false, true)

	app := setupDocumentsApp()
	req := httptest.NewRequest("GET", "/api/v1/validate/rg/"+rg, nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var validation models.RGValidationResponse
	err = json.Unmarshal(body, &validation)
	assert.NoError(t, err)
	assert.True(t, validation.Valid)
}
