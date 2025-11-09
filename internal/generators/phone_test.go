package generators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests for phone generation focusing on:
// 1. Generation by state with correct DDD
// 2. Mobile vs landline types
// 3. Format validation
// 4. Edge case coverage

func TestGeneratePhone_ByState(t *testing.T) {
	ds, err := NewDataStore()
	if err != nil {
		t.Fatalf("Failed to create DataStore: %v", err)
	}
	gen := NewGenerator(ds)

	states := []string{"SP", "RJ", "MG"}
	for _, state := range states {
		t.Run(state, func(t *testing.T) {
			phone, ddd, returnedState, phoneType := gen.GeneratePhone(state, "")

			assert.NotEmpty(t, phone, "Phone should not be empty")
			assert.NotEmpty(t, ddd, "DDD should not be empty")
			assert.Equal(t, state, returnedState, "Returned state should match requested state")
			assert.Contains(t, []string{"mobile", "landline"}, phoneType, "Phone type should be mobile or landline")
			assert.Contains(t, phone, "("+ddd+")", "Phone should contain DDD")
		})
	}
}

func TestGeneratePhone_Mobile(t *testing.T) {
	ds, err := NewDataStore()
	if err != nil {
		t.Fatalf("Failed to create DataStore: %v", err)
	}
	gen := NewGenerator(ds)

	phone, _, _, phoneType := gen.GeneratePhone("SP", "mobile")

	assert.Equal(t, "mobile", phoneType, "Phone type should be mobile")
	assert.Contains(t, phone, "9", "Mobile phone should start with 9")
}

func TestGeneratePhone_Landline(t *testing.T) {
	ds, err := NewDataStore()
	if err != nil {
		t.Fatalf("Failed to create DataStore: %v", err)
	}
	gen := NewGenerator(ds)

	phone, _, _, phoneType := gen.GeneratePhone("SP", "landline")

	assert.Equal(t, "landline", phoneType, "Phone type should be landline")
	unformatted := UnformatPhone(phone)
	firstDigit := int(unformatted[2] - '0')
	assert.True(t, firstDigit >= 3 && firstDigit <= 7, "Landline first digit should be between 3 and 7, got %d in phone %s", firstDigit, phone)
}

func TestGeneratePhone_Randomness(t *testing.T) {
	ds, err := NewDataStore()
	if err != nil {
		t.Fatalf("Failed to create DataStore: %v", err)
	}
	gen := NewGenerator(ds)

	phones := make(map[string]bool)
	const count = 100

	for i := 0; i < count; i++ {
		phone, _, _, _ := gen.GeneratePhone("", "")
		assert.False(t, phones[phone], "Phone %s should be unique (attempt %d)", phone, i+1)
		phones[phone] = true
	}

	assert.Equal(t, count, len(phones), "Should generate %d unique phones", count)
}

func TestFormatPhone(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Valid 10-digit phone",
			input:    "1198765432",
			expected: "(11) 9876-5432",
		},
		{
			name:     "Valid 11-digit phone",
			input:    "11987654321",
			expected: "(11) 98765-4321",
		},
		{
			name:     "Too short",
			input:    "12345",
			expected: "12345",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatPhone(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUnformatPhone(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Formatted phone",
			input:    "(11) 98765-4321",
			expected: "11987654321",
		},
		{
			name:     "Already unformatted",
			input:    "11987654321",
			expected: "11987654321",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UnformatPhone(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGeneratePhone_InvalidDDD(t *testing.T) {
	ds, err := NewDataStore()
	if err != nil {
		t.Fatalf("Failed to create DataStore: %v", err)
	}
	gen := NewGenerator(ds)

	phone, ddd, _, _ := gen.GeneratePhone("XX", "")

	assert.NotEmpty(t, phone, "Phone should still be generated")
	assert.NotEmpty(t, ddd, "DDD should still be generated")
	assert.NotEqual(t, "00", ddd, "DDD should not be 00")
	assert.NotEqual(t, "99", ddd, "DDD should not be 99")
}

func TestGeneratePhone_InvalidFormats(t *testing.T) {
	tests := []struct {
		name        string
		phone       string
		description string
	}{
		{
			name:        "Too short",
			phone:       "123",
			description: "Phone with less than 10 digits should be invalid",
		},
		{
			name:        "Too long",
			phone:       "119876543210",
			description: "Phone with more than 11 digits should be invalid",
		},
		{
			name:        "Empty string",
			phone:       "",
			description: "Empty phone should be invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			unformatted := UnformatPhone(tt.phone)
			isValid := len(unformatted) >= 10 && len(unformatted) <= 11
			assert.False(t, isValid, "%s: %s", tt.description, tt.phone)
		})
	}
}

