package generators

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests for email generation focusing on:
// 1. Varied generation without duplicates
// 2. Custom domain support
// 3. Format validation
// 4. Edge case coverage

func TestGenerateEmail_Random(t *testing.T) {
	ds, err := NewDataStore()
	if err != nil {
		t.Fatalf("Failed to create DataStore: %v", err)
	}
	gen := NewGenerator(ds)

	email, username, domain := gen.GenerateEmail("")

	assert.Contains(t, email, "@", "Email should contain @")
	assert.Contains(t, email, ".", "Email should contain dot")
	assert.Equal(t, username+"@"+domain, email, "Email should be properly formatted")
	assert.NotEmpty(t, username, "Username should not be empty")
	assert.NotEmpty(t, domain, "Domain should not be empty")
}

func TestGenerateEmail_CustomDomain(t *testing.T) {
	ds, err := NewDataStore()
	if err != nil {
		t.Fatalf("Failed to create DataStore: %v", err)
	}
	gen := NewGenerator(ds)

	customDomain := "example.com"
	email, username, domain := gen.GenerateEmail(customDomain)

	assert.Equal(t, customDomain, domain, "Domain should match custom domain")
	assert.Contains(t, email, "@"+customDomain, "Email should contain custom domain")
	assert.NotEmpty(t, username, "Username should not be empty")
}

func TestGenerateEmail_Uniqueness(t *testing.T) {
	ds, err := NewDataStore()
	if err != nil {
		t.Fatalf("Failed to create DataStore: %v", err)
	}
	gen := NewGenerator(ds)

	emails := make(map[string]bool)
	const count = 100

	for i := 0; i < count; i++ {
		email, _, _ := gen.GenerateEmail("")
		assert.False(t, emails[email], "Email %s should be unique (attempt %d)", email, i+1)
		emails[email] = true
	}

	assert.Equal(t, count, len(emails), "Should generate %d unique emails", count)
}

func TestGenerateEmail_FormatValidation(t *testing.T) {
	ds, err := NewDataStore()
	if err != nil {
		t.Fatalf("Failed to create DataStore: %v", err)
	}
	gen := NewGenerator(ds)

	email, _, _ := gen.GenerateEmail("")

	parts := strings.Split(email, "@")
	assert.Len(t, parts, 2, "Email should have exactly one @")
	assert.NotEmpty(t, parts[0], "Username part should not be empty")
	assert.NotEmpty(t, parts[1], "Domain part should not be empty")

	domainParts := strings.Split(parts[1], ".")
	assert.GreaterOrEqual(t, len(domainParts), 2, "Domain should have at least one dot")
}

func TestGenerateEmail_InvalidFormats(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		description string
	}{
		{
			name:        "Without @",
			email:       "testexample.com",
			description: "Email without @ should be invalid",
		},
		{
			name:        "Without domain",
			email:       "test@",
			description: "Email without domain should be invalid",
		},
		{
			name:        "Without user",
			email:       "@example.com",
			description: "Email without user should be invalid",
		},
		{
			name:        "Multiple @",
			email:       "test@test@example.com",
			description: "Email with multiple @ should be invalid",
		},
		{
			name:        "Empty string",
			email:       "",
			description: "Empty email should be invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parts := strings.Split(tt.email, "@")
			isValid := len(parts) == 2 && parts[0] != "" && parts[1] != "" && strings.Contains(parts[1], ".")
			assert.False(t, isValid, "%s: %s", tt.description, tt.email)
		})
	}
}

