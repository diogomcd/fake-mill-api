package testutils

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// AssertHTTPError validates HTTP error response
func AssertHTTPError(t *testing.T, resp *http.Response, expectedStatus int, expectedErrorMsg string) {
	t.Helper()
	assert.Equal(t, expectedStatus, resp.StatusCode, "Expected status code %d, got %d", expectedStatus, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err, "Failed to read response body")

	var errorResponse map[string]interface{}
	err = json.Unmarshal(body, &errorResponse)
	assert.NoError(t, err, "Failed to unmarshal error response")

	if expectedErrorMsg != "" {
		errorMsg, ok := errorResponse["error"].(string)
		assert.True(t, ok, "Error field should be a string")
		assert.Contains(t, errorMsg, expectedErrorMsg, "Error message should contain: %s", expectedErrorMsg)
	}
}

// AssertJSONHeaders validates JSON headers
func AssertJSONHeaders(t *testing.T, resp *http.Response) {
	t.Helper()
	contentType := resp.Header.Get("Content-Type")
	assert.Contains(t, contentType, "application/json", "Content-Type should be application/json")
}

// AssertValidationError validates validation error response
func AssertValidationError(t *testing.T, resp *http.Response, field string) {
	t.Helper()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Expected status code 400")

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err, "Failed to read response body")

	var errorResponse map[string]interface{}
	err = json.Unmarshal(body, &errorResponse)
	assert.NoError(t, err, "Failed to unmarshal error response")

	if field != "" {
		errorMsg, ok := errorResponse["error"].(string)
		if ok {
			assert.Contains(t, errorMsg, field, "Error message should mention field: %s", field)
		}
	}
}
