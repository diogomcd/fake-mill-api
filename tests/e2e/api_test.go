package e2e

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/diogomcd/fake-mill-api/internal/models"
	"github.com/stretchr/testify/assert"
)

const baseURL = "http://127.0.0.1:8080"

func TestHealthCheck(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	resp, err := http.Get(baseURL + "/api/health")
	if err != nil {
		t.Skipf("Server is not running on %s: %v", baseURL, err)
		return
	}
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var health map[string]interface{}
	err = json.Unmarshal(body, &health)
	assert.NoError(t, err)
	assert.Equal(t, "ok", health["status"])
}

func TestPersonGeneration_Complete(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	resp, err := http.Get(baseURL + "/api/v1/person")
	if err != nil {
		t.Skipf("Server is not running on %s: %v", baseURL, err)
		return
	}
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var person models.Person
	err = json.Unmarshal(body, &person)
	assert.NoError(t, err)

	assert.NotEmpty(t, person.Name.FullName)
	assert.NotEmpty(t, person.CPF.Masked)
	assert.NotEmpty(t, person.RG.Masked)
	assert.NotEmpty(t, person.Email.Address)
	assert.NotEmpty(t, person.Phone.InternationalFormat)
	assert.NotEmpty(t, person.Address.Street)
}

func TestCPFValidation_GeneratedCPF(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	resp, err := http.Get(baseURL + "/api/v1/cpf")
	if err != nil {
		t.Skipf("Server is not running on %s: %v", baseURL, err)
		return
	}
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var cpfResp models.CPFResponse
	err = json.Unmarshal(body, &cpfResp)
	assert.NoError(t, err)
	assert.True(t, cpfResp.Valid)

	// Encode the CPF for use in the URL (may contain dots and dashes)
	encodedCPF := url.PathEscape(cpfResp.CPF)
	validationResp, err := http.Get(baseURL + "/api/v1/validate/cpf/" + encodedCPF)
	assert.NoError(t, err)
	defer validationResp.Body.Close()
	assert.Equal(t, 200, validationResp.StatusCode)

	validationBody, err := io.ReadAll(validationResp.Body)
	assert.NoError(t, err)

	var validation models.CPFValidationResponse
	err = json.Unmarshal(validationBody, &validation)
	assert.NoError(t, err)
	assert.True(t, validation.Valid)
}

func TestCNPJValidation_GeneratedCNPJ(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	resp, err := http.Get(baseURL + "/api/v1/cnpj")
	if err != nil {
		t.Skipf("Server is not running on %s: %v", baseURL, err)
		return
	}
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var cnpjResp models.CNPJResponse
	err = json.Unmarshal(body, &cnpjResp)
	assert.NoError(t, err)
	assert.True(t, cnpjResp.Valid)

	// Remove formatting of CNPJ for use in the URL (to avoid problems with / in path)
	// Remove dots, slashes and dashes
	cnpjUnmasked := strings.ReplaceAll(cnpjResp.CNPJ, ".", "")
	cnpjUnmasked = strings.ReplaceAll(cnpjUnmasked, "/", "")
	cnpjUnmasked = strings.ReplaceAll(cnpjUnmasked, "-", "")

	// Encode the CNPJ without formatting for use in the URL
	encodedCNPJ := url.PathEscape(cnpjUnmasked)
	validationResp, err := http.Get(baseURL + "/api/v1/validate/cnpj/" + encodedCNPJ)
	assert.NoError(t, err)
	defer validationResp.Body.Close()
	assert.Equal(t, 200, validationResp.StatusCode)

	validationBody, err := io.ReadAll(validationResp.Body)
	assert.NoError(t, err)

	var validation models.CNPJValidationResponse
	err = json.Unmarshal(validationBody, &validation)
	assert.NoError(t, err)
	assert.True(t, validation.Valid)
}

func TestRateLimiting_Real(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	// Check if the server is running
	resp, err := http.Get(baseURL + "/api/health")
	if err != nil {
		t.Skipf("Server is not running on %s: %v", baseURL, err)
		return
	}
	resp.Body.Close()

	client := &http.Client{}

	// Make 59 requests (leaving space to not exceed the limit of 60)
	successCount := 0
	for i := 0; i < 59; i++ {
		req, _ := http.NewRequest("GET", baseURL+"/api/v1/person", nil)
		resp, err := client.Do(req)
		assert.NoError(t, err)
		if resp != nil && resp.StatusCode == 200 {
			successCount++
			resp.Body.Close()
		}
		// Small delay to avoid too fast requests
		time.Sleep(10 * time.Millisecond)
	}

	// The 60th request should pass
	req, _ := http.NewRequest("GET", baseURL+"/api/v1/person", nil)
	resp, err = client.Do(req)
	assert.NoError(t, err)
	if resp != nil && resp.StatusCode == 200 {
		successCount++
		resp.Body.Close()
	}

	// The 61st request should be blocked
	time.Sleep(100 * time.Millisecond)
	req, _ = http.NewRequest("GET", baseURL+"/api/v1/person", nil)
	resp, err = client.Do(req)
	assert.NoError(t, err)
	if resp != nil {
		assert.Equal(t, 429, resp.StatusCode, "Request 61 should be rate limited")
		resp.Body.Close()
	}
}

func TestConcurrentRequests(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	// Check if the server is running
	resp, err := http.Get(baseURL + "/api/health")
	if err != nil {
		t.Skipf("Server is not running on %s: %v", baseURL, err)
		return
	}
	resp.Body.Close()

	const numGoroutines = 50
	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := http.Get(baseURL + "/api/v1/person")
			if err != nil {
				errors <- err
				return
			}
			if resp.StatusCode != 200 && resp.StatusCode != 429 {
				errors <- assert.AnError
			}
			resp.Body.Close()
		}()
	}

	wg.Wait()
	close(errors)

	errorCount := 0
	for err := range errors {
		if err != nil {
			errorCount++
		}
	}

	assert.Less(t, errorCount, numGoroutines/2, "Most requests should succeed")
}
