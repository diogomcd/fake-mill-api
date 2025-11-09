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

func setupTestApp() (*fiber.App, *generators.Generator) {
	ds, err := generators.NewDataStore()
	if err != nil {
		panic(fmt.Sprintf("Failed to create DataStore: %v", err))
	}
	gen := generators.NewGenerator(ds)

	app := fiber.New()
	app.Use(middleware.InjectGenerator(gen))
	v1 := app.Group("/api/v1")
	v1.Get("/person", PersonHandler)

	return app, gen
}

func TestPersonHandler_Success_NoParams(t *testing.T) {
	app, _ := setupTestApp()

	req := httptest.NewRequest("GET", "/api/v1/person", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	testutils.AssertJSONHeaders(t, resp)

	body, _ := io.ReadAll(resp.Body)
	var person models.Person
	err = json.Unmarshal(body, &person)
	assert.NoError(t, err)
	assert.NotEmpty(t, person.Name.FullName)
}

func TestPersonHandler_Success_Quantity5(t *testing.T) {
	app, _ := setupTestApp()

	req := httptest.NewRequest("GET", "/api/v1/person?quantity=5", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var people []models.Person
	err = json.Unmarshal(body, &people)
	assert.NoError(t, err)
	assert.Len(t, people, 5)
}

func TestPersonHandler_Success_Quantity1(t *testing.T) {
	app, _ := setupTestApp()

	req := httptest.NewRequest("GET", "/api/v1/person?quantity=1", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var person models.Person
	err = json.Unmarshal(body, &person)
	assert.NoError(t, err)
	assert.NotEmpty(t, person.Name.FullName)
}

func TestPersonHandler_Success_GenderMale(t *testing.T) {
	app, _ := setupTestApp()

	req := httptest.NewRequest("GET", "/api/v1/person?gender=male", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var person models.Person
	err = json.Unmarshal(body, &person)
	assert.NoError(t, err)
	assert.Equal(t, "male", person.Gender)
}

func TestPersonHandler_Success_GenderFemale(t *testing.T) {
	app, _ := setupTestApp()

	req := httptest.NewRequest("GET", "/api/v1/person?gender=female", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var person models.Person
	err = json.Unmarshal(body, &person)
	assert.NoError(t, err)
	assert.Equal(t, "female", person.Gender)
}

func TestPersonHandler_Success_StateSP(t *testing.T) {
	app, _ := setupTestApp()

	req := httptest.NewRequest("GET", "/api/v1/person?state=SP", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var person models.Person
	err = json.Unmarshal(body, &person)
	assert.NoError(t, err)
	assert.Equal(t, "SP", person.Address.State)
}

func TestPersonHandler_Error_QuantityTooHigh(t *testing.T) {
	app, _ := setupTestApp()

	req := httptest.NewRequest("GET", "/api/v1/person?quantity=201", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var person models.Person
	err = json.Unmarshal(body, &person)
	assert.NoError(t, err)
}

func TestPersonHandler_Error_InvalidGender(t *testing.T) {
	app, _ := setupTestApp()

	req := httptest.NewRequest("GET", "/api/v1/person?gender=invalid", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var person models.Person
	err = json.Unmarshal(body, &person)
	assert.NoError(t, err)
}

func TestPersonHandler_Error_InvalidState(t *testing.T) {
	app, _ := setupTestApp()

	req := httptest.NewRequest("GET", "/api/v1/person?state=XX", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var person models.Person
	err = json.Unmarshal(body, &person)
	assert.NoError(t, err)
}
