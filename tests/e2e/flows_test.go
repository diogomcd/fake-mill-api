package e2e

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/diogomcd/fake-mill-api/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestFlow_PersonCPFValidation(t *testing.T) {
	if testing.Short() {
		t.Skip("Pulando teste E2E no modo short")
	}

	// Pequeno delay para evitar rate limiting de testes anteriores
	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get(baseURL + "/api/v1/person")
	if err != nil {
		t.Skipf("Servidor não está rodando em %s: %v", baseURL, err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == 429 {
		t.Skip("Rate limited, skipping test")
		return
	}
	assert.Equal(t, 200, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var person models.Person
	err = json.Unmarshal(body, &person)
	assert.NoError(t, err)

	// Encode the CPF for use in the URL
	encodedCPF := url.PathEscape(person.CPF.Unmasked)
	validationResp, err := http.Get(baseURL + "/api/v1/validate/cpf/" + encodedCPF)
	assert.NoError(t, err)
	defer validationResp.Body.Close()
	assert.Equal(t, 200, validationResp.StatusCode)

	validationBody, err := io.ReadAll(validationResp.Body)
	assert.NoError(t, err)

	var validation models.CPFValidationResponse
	err = json.Unmarshal(validationBody, &validation)
	assert.NoError(t, err)
	assert.True(t, validation.Valid, "Generated CPF should be valid")
}

func TestFlow_CompanyCNPJValidation(t *testing.T) {
	if testing.Short() {
		t.Skip("Pulando teste E2E no modo short")
	}

	// Small delay to avoid rate limiting from previous tests
	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get(baseURL + "/api/v1/company")
	if err != nil {
		t.Skipf("Servidor não está rodando em %s: %v", baseURL, err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == 429 {
		t.Skip("Rate limited, skipping test")
		return
	}
	assert.Equal(t, 200, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var company models.CompanyResponse
	err = json.Unmarshal(body, &company)
	assert.NoError(t, err)

	// Encode the CNPJ for use in the URL (may contain dots and bars)
	encodedCNPJ := url.PathEscape(company.CNPJ)
	validationResp, err := http.Get(baseURL + "/api/v1/validate/cnpj/" + encodedCNPJ)
	assert.NoError(t, err)
	defer validationResp.Body.Close()
	assert.Equal(t, 200, validationResp.StatusCode)

	validationBody, err := io.ReadAll(validationResp.Body)
	assert.NoError(t, err)

	var validation models.CNPJValidationResponse
	err = json.Unmarshal(validationBody, &validation)
	assert.NoError(t, err)
	assert.True(t, validation.Valid, "Generated CNPJ should be valid")
}

func TestFlow_MultipleDocumentsUniqueness(t *testing.T) {
	if testing.Short() {
		t.Skip("Pulando teste E2E no modo short")
	}

	// Small delay to avoid rate limiting from previous tests
	time.Sleep(200 * time.Millisecond)

	resp, err := http.Get(baseURL + "/api/v1/person?quantity=20")
	if err != nil {
		t.Skipf("Servidor não está rodando em %s: %v", baseURL, err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == 429 {
		t.Skip("Rate limited, skipping test")
		return
	}
	assert.Equal(t, 200, resp.StatusCode, "Expected 200 but got %d", resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var people []models.Person
	err = json.Unmarshal(body, &people)
	if err != nil {
		t.Logf("Failed to unmarshal as array, body: %s", string(body))
	}
	assert.NoError(t, err)
	assert.Len(t, people, 20)

	cpfs := make(map[string]bool)
	for _, person := range people {
		cpf := person.CPF.Unmasked
		assert.False(t, cpfs[cpf], "CPF %s should be unique", cpf)
		cpfs[cpf] = true

		// Encode the CPF for use in the URL
		encodedCPF := url.PathEscape(cpf)
		validationResp, err := http.Get(baseURL + "/api/v1/validate/cpf/" + encodedCPF)
		assert.NoError(t, err)
		if validationResp.StatusCode == 429 {
			validationResp.Body.Close()
			t.Logf("Rate limited during validation, skipping remaining validations")
			break
		}
		assert.Equal(t, 200, validationResp.StatusCode)

		validationBody, err := io.ReadAll(validationResp.Body)
		assert.NoError(t, err)
		validationResp.Body.Close()

		var validation models.CPFValidationResponse
		err = json.Unmarshal(validationBody, &validation)
		assert.NoError(t, err)
		assert.True(t, validation.Valid, "CPF %s should be valid", cpf)
	}

	assert.Equal(t, 20, len(cpfs), "All 20 CPFs should be unique")
}
