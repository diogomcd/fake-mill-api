package generators

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests for RG generation focusing on:
// 1. Generation by state with correct issuing agencies
// 2. Variability across multiple generations
// 3. Format validation
// 4. Edge case coverage

func TestGenerateRG_ByState(t *testing.T) {
	ds, err := NewDataStore()
	if err != nil {
		t.Fatalf("Failed to create DataStore: %v", err)
	}
	gen := NewGenerator(ds)

	states := []string{"SP", "RJ", "MG"}
	for _, state := range states {
		t.Run(state, func(t *testing.T) {
			rg, returnedState, issuer, issueDate, expirationDate := gen.GenerateRG(state, false, true)

			assert.NotEmpty(t, rg, "RG should not be empty")
			assert.Equal(t, state, returnedState, "Returned state should match requested state")
			assert.Equal(t, "SSP", issuer, "Issuer should be SSP")
			assert.NotEmpty(t, issueDate, "Issue date should not be empty")
			assert.NotEmpty(t, expirationDate, "Expiration date should not be empty")
			assert.True(t, ValidateRG(rg), "Generated RG should be valid")
		})
	}
}

func TestGenerateRG_Randomness(t *testing.T) {
	ds, err := NewDataStore()
	if err != nil {
		t.Fatalf("Failed to create DataStore: %v", err)
	}
	gen := NewGenerator(ds)

	states := []string{"SP", "RJ", "MG"}
	const countPerState = 50

	for _, state := range states {
		t.Run(state, func(t *testing.T) {
			rgs := make(map[string]bool)

			for i := 0; i < countPerState; i++ {
				rg, _, _, _, _ := gen.GenerateRG(state, false, true)
				assert.False(t, rgs[rg], "RG %s should be unique (attempt %d)", rg, i+1)
				rgs[rg] = true
			}

			assert.Equal(t, countPerState, len(rgs), "Should generate %d unique RGs for state %s", countPerState, state)
		})
	}
}

func TestGenerateRG_Formatted(t *testing.T) {
	ds, err := NewDataStore()
	if err != nil {
		t.Fatalf("Failed to create DataStore: %v", err)
	}
	gen := NewGenerator(ds)

	rg, _, _, _, _ := gen.GenerateRG("SP", true, true)

	assert.Contains(t, rg, ".", "Formatted RG should contain dots")
	assert.Contains(t, rg, "-", "Formatted RG should contain dash")
	assert.True(t, ValidateRG(rg), "Formatted RG should be valid")
}

func TestGenerateRG_Unformatted(t *testing.T) {
	ds, err := NewDataStore()
	if err != nil {
		t.Fatalf("Failed to create DataStore: %v", err)
	}
	gen := NewGenerator(ds)

	rg, _, _, _, _ := gen.GenerateRG("SP", false, true)

	assert.False(t, strings.Contains(rg, "."), "Unformatted RG should not contain dots")
	assert.False(t, strings.Contains(rg, "-"), "Unformatted RG should not contain dash")
	assert.Len(t, rg, 9, "Unformatted RG should have 9 characters")
	assert.True(t, ValidateRG(rg), "Unformatted RG should be valid")
}

func TestGenerateRG_Invalid(t *testing.T) {
	ds, err := NewDataStore()
	if err != nil {
		t.Fatalf("Failed to create DataStore: %v", err)
	}
	gen := NewGenerator(ds)

	rg, _, _, _, _ := gen.GenerateRG("SP", false, false)

	assert.Len(t, rg, 9, "Invalid RG should still have 9 characters")
	assert.False(t, ValidateRG(rg), "Invalid RG should not pass validation")
}

func TestFormatRG(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Valid 9-digit RG",
			input:    "123456789",
			expected: "12.345.678-9",
		},
		{
			name:     "RG with X check digit",
			input:    "12345678X",
			expected: "12.345.678-X",
		},
		{
			name:     "Invalid length",
			input:    "12345",
			expected: "12345",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatRG(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateRG_Invalid(t *testing.T) {
	tests := []struct {
		name        string
		rg          string
		description string
	}{
		{
			name:        "Too short",
			rg:          "12345",
			description: "RG with less than 9 characters should be invalid",
		},
		{
			name:        "Too long",
			rg:          "1234567890",
			description: "RG with more than 9 characters should be invalid",
		},
		{
			name:        "Invalid check digit",
			rg:          "123456780",
			description: "RG with incorrect check digit should be invalid",
		},
		{
			name:        "Empty string",
			rg:          "",
			description: "Empty RG should be invalid",
		},
		{
			name:        "Invalid characters",
			rg:          "12345678A",
			description: "RG with invalid characters should be invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateRG(tt.rg)
			assert.False(t, result, "%s: %s", tt.description, tt.rg)
		})
	}
}

