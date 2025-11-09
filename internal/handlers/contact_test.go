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

func setupContactApp() *fiber.App {
	ds, err := generators.NewDataStore()
	if err != nil {
		panic(fmt.Sprintf("Failed to create DataStore: %v", err))
	}
	gen := generators.NewGenerator(ds)

	app := fiber.New()
	app.Use(middleware.InjectGenerator(gen))
	v1 := app.Group("/api/v1")
	v1.Get("/email", EmailHandler)
	v1.Get("/phone", PhoneHandler)

	return app
}

func TestEmailHandler_Success(t *testing.T) {
	app := setupContactApp()

	req := httptest.NewRequest("GET", "/api/v1/email", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	testutils.AssertJSONHeaders(t, resp)

	body, _ := io.ReadAll(resp.Body)
	var emailResp models.EmailResponse
	err = json.Unmarshal(body, &emailResp)
	assert.NoError(t, err)
	assert.NotEmpty(t, emailResp.Email)
	assert.NotEmpty(t, emailResp.Username)
	assert.NotEmpty(t, emailResp.Domain)
}

func TestEmailHandler_CustomDomain(t *testing.T) {
	app := setupContactApp()

	req := httptest.NewRequest("GET", "/api/v1/email?domain=example.com", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var emailResp models.EmailResponse
	err = json.Unmarshal(body, &emailResp)
	assert.NoError(t, err)
	assert.Equal(t, "example.com", emailResp.Domain)
}

func TestPhoneHandler_Success(t *testing.T) {
	app := setupContactApp()

	req := httptest.NewRequest("GET", "/api/v1/phone", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var phoneResp models.PhoneResponse
	err = json.Unmarshal(body, &phoneResp)
	assert.NoError(t, err)
	assert.NotEmpty(t, phoneResp.Phone)
	assert.NotEmpty(t, phoneResp.DDD)
	assert.NotEmpty(t, phoneResp.State)
}

func TestPhoneHandler_StateSP(t *testing.T) {
	app := setupContactApp()

	req := httptest.NewRequest("GET", "/api/v1/phone?state=SP", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var phoneResp models.PhoneResponse
	err = json.Unmarshal(body, &phoneResp)
	assert.NoError(t, err)
	assert.Equal(t, "SP", phoneResp.State)
}

func TestPhoneHandler_InvalidState(t *testing.T) {
	app := setupContactApp()

	req := httptest.NewRequest("GET", "/api/v1/phone?state=XX", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var phoneResp models.PhoneResponse
	err = json.Unmarshal(body, &phoneResp)
	assert.NoError(t, err)
	assert.NotEqual(t, "XX", phoneResp.State)
}
