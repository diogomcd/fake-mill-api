package generators

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests for CNPJ generation focusing on:
// 1. Varied generation without duplicates (reliability)
// 2. Validation against reserved cases (data security)
// 3. Robust formatting (UX)
// 4. Edge case coverage (maintainability)

func TestGenerateCNPJ_Valid_Formatted(t *testing.T) {
	gen := NewGenerator(nil)
	cnpj := gen.GenerateCNPJ(true, true)

	assert.Len(t, cnpj, 18, "Formatted CNPJ should have 18 characters")
	assert.Contains(t, cnpj, ".", "Formatted CNPJ should contain dots")
	assert.Contains(t, cnpj, "/", "Formatted CNPJ should contain slash")
	assert.Contains(t, cnpj, "-", "Formatted CNPJ should contain dash")
	assert.True(t, ValidateCNPJ(cnpj), "Generated CNPJ should be valid")
}

func TestGenerateCNPJ_Valid_Unformatted(t *testing.T) {
	gen := NewGenerator(nil)
	cnpj := gen.GenerateCNPJ(false, true)

	assert.Len(t, cnpj, 14, "Unformatted CNPJ should have 14 characters")
	assert.False(t, strings.Contains(cnpj, "."), "Unformatted CNPJ should not contain dots")
	assert.False(t, strings.Contains(cnpj, "/"), "Unformatted CNPJ should not contain slash")
	assert.False(t, strings.Contains(cnpj, "-"), "Unformatted CNPJ should not contain dash")
	assert.True(t, ValidateCNPJ(cnpj), "Generated CNPJ should be valid")
}

func TestGenerateCNPJ_Invalid(t *testing.T) {
	gen := NewGenerator(nil)
	cnpj := gen.GenerateCNPJ(false, false)

	assert.Len(t, cnpj, 14, "Invalid CNPJ should still have 14 characters")
	assert.False(t, ValidateCNPJ(cnpj), "Invalid CNPJ should not pass validation")
}

func TestFormatCNPJ(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Valid 14-digit CNPJ",
			input:    "12345678000190",
			expected: "12.345.678/0001-90",
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
			result := FormatCNPJ(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGenerateCNPJ_Uniqueness(t *testing.T) {
	gen := NewGenerator(nil)
	cnpjs := make(map[string]bool)
	const count = 100

	for i := 0; i < count; i++ {
		cnpj := gen.GenerateCNPJ(false, true)
		assert.False(t, cnpjs[cnpj], "CNPJ %s should be unique (attempt %d)", cnpj, i+1)
		cnpjs[cnpj] = true
	}

	assert.Equal(t, count, len(cnpjs), "Should generate %d unique CNPJs", count)
}

func TestValidateCNPJ_Valid_WithMask(t *testing.T) {
	gen := NewGenerator(nil)
	validCNPJs := make([]string, 5)
	for i := 0; i < 5; i++ {
		validCNPJs[i] = gen.GenerateCNPJ(true, true)
	}

	for _, cnpj := range validCNPJs {
		t.Run(cnpj, func(t *testing.T) {
			assert.True(t, ValidateCNPJ(cnpj), "CNPJ %s should be valid", cnpj)
		})
	}
}

func TestValidateCNPJ_Valid_WithoutMask(t *testing.T) {
	gen := NewGenerator(nil)
	validCNPJs := make([]string, 5)
	for i := 0; i < 5; i++ {
		validCNPJs[i] = gen.GenerateCNPJ(false, true)
	}

	for _, cnpj := range validCNPJs {
		t.Run(cnpj, func(t *testing.T) {
			assert.True(t, ValidateCNPJ(cnpj), "CNPJ %s should be valid", cnpj)
		})
	}
}

func TestValidateCNPJ_Invalid(t *testing.T) {
	tests := []struct {
		name        string
		cnpj        string
		description string
	}{
		{
			name:        "Reserved CNPJ - all zeros",
			cnpj:         "00000000000000",
			description: "CNPJ with all zeros should be invalid",
		},
		{
			name:        "Reserved CNPJ - all ones",
			cnpj:         "11111111111111",
			description: "CNPJ with all ones should be invalid",
		},
		{
			name:        "Reserved CNPJ - all twos",
			cnpj:         "22222222222222",
			description: "CNPJ with all twos should be invalid",
		},
		{
			name:        "Reserved CNPJ - all nines",
			cnpj:         "99999999999999",
			description: "CNPJ with all nines should be invalid",
		},
		{
			name:        "Incorrect formatting - wrong separators",
			cnpj:         "12-345-678/0001.90",
			description: "CNPJ with wrong separator positions should be invalid",
		},
		{
			name:        "Incorrect formatting - missing separator",
			cnpj:         "12.345.6780001-90",
			description: "CNPJ with missing slash should be invalid",
		},
		{
			name:        "Invalid characters - letters",
			cnpj:         "AB.CDE.FGH/IJKL-MN",
			description: "CNPJ with letters should be invalid",
		},
		{
			name:        "Invalid characters - special chars",
			cnpj:         "12.345.678/0001-9@",
			description: "CNPJ with special characters should be invalid",
		},
		{
			name:        "Invalid characters - mixed",
			cnpj:         "12a.345.678/0001-90",
			description: "CNPJ with mixed characters should be invalid",
		},
		{
			name:        "Too short",
			cnpj:         "12345",
			description: "CNPJ with less than 14 digits should be invalid",
		},
		{
			name:        "Too long",
			cnpj:         "12345678000190123",
			description: "CNPJ with more than 14 digits should be invalid",
		},
		{
			name:        "Incorrect check digits",
			cnpj:         "12.345.678/0001-00",
			description: "CNPJ with incorrect check digits should be invalid",
		},
		{
			name:        "Empty string",
			cnpj:         "",
			description: "Empty CNPJ should be invalid",
		},
		{
			name:        "Whitespace only",
			cnpj:         " ",
			description: "CNPJ with only whitespace should be invalid",
		},
		{
			name:        "Multiple whitespaces",
			cnpj:         "   ",
			description: "CNPJ with multiple whitespaces should be invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateCNPJ(tt.cnpj)
			assert.False(t, result, "%s: %s", tt.description, tt.cnpj)
		})
	}
}

