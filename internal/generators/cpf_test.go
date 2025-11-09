package generators

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests for CPF generation focusing on:
// 1. Varied generation without duplicates (reliability)
// 2. Validation against reserved cases (data security)
// 3. Robust formatting (UX)
// 4. Edge case coverage (maintainability)

func TestGenerateCPF_Valid_Formatted(t *testing.T) {
	gen := NewGenerator(nil)
	cpf := gen.GenerateCPF(true, true)

	assert.Len(t, cpf, 14, "Formatted CPF should have 14 characters")
	assert.Contains(t, cpf, ".", "Formatted CPF should contain dots")
	assert.Contains(t, cpf, "-", "Formatted CPF should contain dash")
	assert.True(t, ValidateCPF(cpf), "Generated CPF should be valid")
}

func TestGenerateCPF_Valid_Unformatted(t *testing.T) {
	gen := NewGenerator(nil)
	cpf := gen.GenerateCPF(false, true)

	assert.Len(t, cpf, 11, "Unformatted CPF should have 11 characters")
	assert.False(t, strings.Contains(cpf, "."), "Unformatted CPF should not contain dots")
	assert.False(t, strings.Contains(cpf, "-"), "Unformatted CPF should not contain dash")
	assert.True(t, ValidateCPF(cpf), "Generated CPF should be valid")
}

func TestGenerateCPF_Invalid(t *testing.T) {
	gen := NewGenerator(nil)
	cpf := gen.GenerateCPF(false, false)

	assert.Len(t, cpf, 11, "Invalid CPF should still have 11 characters")
	assert.False(t, ValidateCPF(cpf), "Invalid CPF should not pass validation")
}

func TestFormatCPF(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Valid 11-digit CPF",
			input:    "12345678910",
			expected: "123.456.789-10",
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
			result := FormatCPF(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGenerateCPF_Uniqueness(t *testing.T) {
	gen := NewGenerator(nil)
	cpfs := make(map[string]bool)
	const count = 100

	for i := 0; i < count; i++ {
		cpf := gen.GenerateCPF(false, true)
		assert.False(t, cpfs[cpf], "CPF %s should be unique (attempt %d)", cpf, i+1)
		cpfs[cpf] = true
	}

	assert.Equal(t, count, len(cpfs), "Should generate %d unique CPFs", count)
}

func TestGenerateCPF_Distribution(t *testing.T) {
	gen := NewGenerator(nil)
	digitCounts := make(map[int]int)
	const count = 100

	for i := 0; i < count; i++ {
		cpf := gen.GenerateCPF(false, true)
		for _, char := range cpf {
			digit := int(char - '0')
			digitCounts[digit]++
		}
	}

	for digit := 0; digit <= 9; digit++ {
		assert.Greater(t, digitCounts[digit], 0, "Digit %d should appear at least once", digit)
	}
}

func TestValidateCPF_Valid_WithMask(t *testing.T) {
	gen := NewGenerator(nil)
	validCPFs := make([]string, 5)
	for i := 0; i < 5; i++ {
		validCPFs[i] = gen.GenerateCPF(true, true)
	}

	for _, cpf := range validCPFs {
		t.Run(cpf, func(t *testing.T) {
			assert.True(t, ValidateCPF(cpf), "CPF %s should be valid", cpf)
		})
	}
}

func TestValidateCPF_Valid_WithoutMask(t *testing.T) {
	gen := NewGenerator(nil)
	validCPFs := make([]string, 5)
	for i := 0; i < 5; i++ {
		validCPFs[i] = gen.GenerateCPF(false, true)
	}

	for _, cpf := range validCPFs {
		t.Run(cpf, func(t *testing.T) {
			assert.True(t, ValidateCPF(cpf), "CPF %s should be valid", cpf)
		})
	}
}

func TestValidateCPF_Invalid(t *testing.T) {
	tests := []struct {
		name        string
		cpf         string
		description string
	}{
		{
			name:        "Reserved CPF - all zeros",
			cpf:         "00000000000",
			description: "CPF with all zeros should be invalid",
		},
		{
			name:        "Reserved CPF - all ones",
			cpf:         "11111111111",
			description: "CPF with all ones should be invalid",
		},
		{
			name:        "Reserved CPF - all twos",
			cpf:         "22222222222",
			description: "CPF with all twos should be invalid",
		},
		{
			name:        "Reserved CPF - all threes",
			cpf:         "33333333333",
			description: "CPF with all threes should be invalid",
		},
		{
			name:        "Reserved CPF - all fours",
			cpf:         "44444444444",
			description: "CPF with all fours should be invalid",
		},
		{
			name:        "Reserved CPF - all fives",
			cpf:         "55555555555",
			description: "CPF with all fives should be invalid",
		},
		{
			name:        "Reserved CPF - all sixes",
			cpf:         "66666666666",
			description: "CPF with all sixes should be invalid",
		},
		{
			name:        "Reserved CPF - all sevens",
			cpf:         "77777777777",
			description: "CPF with all sevens should be invalid",
		},
		{
			name:        "Reserved CPF - all eights",
			cpf:         "88888888888",
			description: "CPF with all eights should be invalid",
		},
		{
			name:        "Reserved CPF - all nines",
			cpf:         "99999999999",
			description: "CPF with all nines should be invalid",
		},
		{
			name:        "Incorrect formatting - wrong separators",
			cpf:         "123-456-789.10",
			description: "CPF with wrong separator positions should be invalid",
		},
		{
			name:        "Incorrect formatting - missing separator",
			cpf:         "123.456.78910",
			description: "CPF with missing dash should be invalid",
		},
		{
			name:        "Invalid characters - letters",
			cpf:         "ABC.DEF.GHI-JK",
			description: "CPF with letters should be invalid",
		},
		{
			name:        "Invalid characters - special chars",
			cpf:         "123.456.789-0@",
			description: "CPF with special characters should be invalid",
		},
		{
			name:        "Invalid characters - mixed",
			cpf:         "12a.456.789-10",
			description: "CPF with mixed characters should be invalid",
		},
		{
			name:        "Too short",
			cpf:         "12345",
			description: "CPF with less than 11 digits should be invalid",
		},
		{
			name:        "Too long",
			cpf:         "123456789101112",
			description: "CPF with more than 11 digits should be invalid",
		},
		{
			name:        "Incorrect check digits",
			cpf:         "123.456.789-00",
			description: "CPF with incorrect check digits should be invalid",
		},
		{
			name:        "Empty string",
			cpf:         "",
			description: "Empty CPF should be invalid",
		},
		{
			name:        "Whitespace only",
			cpf:         " ",
			description: "CPF with only whitespace should be invalid",
		},
		{
			name:        "Multiple whitespaces",
			cpf:         "   ",
			description: "CPF with multiple whitespaces should be invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateCPF(tt.cpf)
			assert.False(t, result, "%s: %s", tt.description, tt.cpf)
		})
	}
}

