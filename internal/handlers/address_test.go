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

func setupAddressApp() *fiber.App {
	ds, err := generators.NewDataStore()
	if err != nil {
		panic(fmt.Sprintf("Failed to create DataStore: %v", err))
	}
	gen := generators.NewGenerator(ds)

	app := fiber.New()
	app.Use(middleware.InjectGenerator(gen))
	v1 := app.Group("/api/v1")
	v1.Get("/address", AddressHandler)
	v1.Get("/zipcode", ZipcodeHandler)

	return app
}

func TestAddressHandler_Success(t *testing.T) {
	app := setupAddressApp()

	req := httptest.NewRequest("GET", "/api/v1/address", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	testutils.AssertJSONHeaders(t, resp)

	body, _ := io.ReadAll(resp.Body)
	var address models.Address
	err = json.Unmarshal(body, &address)
	assert.NoError(t, err)
	assert.NotEmpty(t, address.Street)
	assert.NotEmpty(t, address.State)
}

func TestAddressHandler_StateSP(t *testing.T) {
	app := setupAddressApp()

	req := httptest.NewRequest("GET", "/api/v1/address?state=SP", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var address models.Address
	err = json.Unmarshal(body, &address)
	assert.NoError(t, err)
	assert.Equal(t, "SP", address.State)
}

func TestZipcodeHandler_Success(t *testing.T) {
	app := setupAddressApp()

	req := httptest.NewRequest("GET", "/api/v1/zipcode", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var zipcodeResp models.ZipcodeResponse
	err = json.Unmarshal(body, &zipcodeResp)
	assert.NoError(t, err)
	assert.NotEmpty(t, zipcodeResp.Zipcode)
	assert.NotEmpty(t, zipcodeResp.State)
	assert.NotEmpty(t, zipcodeResp.City)
}
