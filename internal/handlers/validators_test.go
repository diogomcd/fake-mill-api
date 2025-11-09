package handlers

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/diogomcd/fake-mill-api/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestValidatePhone(t *testing.T) {
	app := fiber.New()
	app.Get("/validate/phone", ValidatePhone)

	tests := []struct {
		description         string
		url                 string
		expectedStatusCode  int
		expectedValid       bool
		expectedPhoneNumber string
		expectedCountryCode string
		expectedDDI         int32
		expectError         bool
		expectedErrorMsg    string
	}{
		{
			description:         "Valid BR mobile number (default)",
			url:                 "/validate/phone?phone_number=11999999999",
			expectedStatusCode:  200,
			expectedValid:       true,
			expectedPhoneNumber: "11999999999",
			expectedCountryCode: "BR",
			expectedDDI:         55,
			expectError:         false,
		},
		{
			description:         "Valid BR mobile number with country_code",
			url:                 "/validate/phone?phone_number=11987654321&country_code=BR",
			expectedStatusCode:  200,
			expectedValid:       true,
			expectedPhoneNumber: "11987654321",
			expectedCountryCode: "BR",
			expectedDDI:         55,
			expectError:         false,
		},
		{
			description:         "Valid BR number with +55 prefix",
			url:                 "/validate/phone?phone_number=%2B5511998765432",
			expectedStatusCode:  200,
			expectedValid:       true,
			expectedPhoneNumber: "+5511998765432",
			expectedCountryCode: "BR",
			expectedDDI:         55,
			expectError:         false,
		},
		{
			description:         "Valid US number with country_code",
			url:                 "/validate/phone?phone_number=6502530000&country_code=US",
			expectedStatusCode:  200,
			expectedValid:       true,
			expectedPhoneNumber: "6502530000",
			expectedCountryCode: "US",
			expectedDDI:         1,
			expectError:         false,
		},
		{
			description:         "Invalid number - too short",
			url:                 "/validate/phone?phone_number=12345",
			expectedStatusCode:  200,
			expectedValid:       false,
			expectedPhoneNumber: "12345",
			expectError:         false,
		},
		{
			description:         "Invalid number - bad format",
			url:                 "/validate/phone?phone_number=abcdef",
			expectedStatusCode:  200,
			expectedValid:       false,
			expectedPhoneNumber: "abcdef",
			expectError:         false,
		},
		{
			description:        "Missing phone_number query param",
			url:                "/validate/phone",
			expectedStatusCode: 400,
			expectError:        true,
			expectedErrorMsg:   "phone_number is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.url, nil)
			resp, err := app.Test(req, -1)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatusCode, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			if tt.expectError {
				var errorResponse map[string]string
				err = json.Unmarshal(body, &errorResponse)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedErrorMsg, errorResponse["error"])
			} else {
				var result models.PhoneValidationResponse
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedValid, result.Valid)
				assert.Equal(t, tt.expectedPhoneNumber, result.PhoneNumber)

				if tt.expectedValid {
					assert.Equal(t, tt.expectedCountryCode, result.CountryCode)
					assert.Equal(t, tt.expectedDDI, result.DDI)
				}
			}
		})
	}
}
