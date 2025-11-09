package generators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests for company generation focusing on:
// 1. Complete company generation
// 2. Valid CNPJ
// 3. Name uniqueness

func TestGenerateCompany_Complete(t *testing.T) {
	ds, err := NewDataStore()
	if err != nil {
		t.Fatalf("Failed to create DataStore: %v", err)
	}
	gen := NewGenerator(ds)

	company := gen.GenerateCompany()

	assert.NotEmpty(t, company.Name, "Company name should not be empty")
	assert.NotEmpty(t, company.TradeName, "Trade name should not be empty")
	assert.NotEmpty(t, company.CNPJ, "CNPJ should not be empty")
	assert.NotEmpty(t, company.StateRegistration, "State registration should not be empty")
	assert.NotEmpty(t, company.Email, "Email should not be empty")
	assert.NotEmpty(t, company.Phone, "Phone should not be empty")
	assert.NotEmpty(t, company.FoundedAt, "Founded date should not be empty")
	assert.NotEmpty(t, company.Address.Street, "Address should not be empty")
}

func TestGenerateCompany_ValidCNPJ(t *testing.T) {
	ds, err := NewDataStore()
	if err != nil {
		t.Fatalf("Failed to create DataStore: %v", err)
	}
	gen := NewGenerator(ds)

	company := gen.GenerateCompany()

	assert.True(t, ValidateCNPJ(company.CNPJ), "Generated CNPJ should be valid")
}

func TestGenerateCompany_NameUniqueness(t *testing.T) {
	ds, err := NewDataStore()
	if err != nil {
		t.Fatalf("Failed to create DataStore: %v", err)
	}
	gen := NewGenerator(ds)

	names := make(map[string]bool)
	const count = 50

	for i := 0; i < count; i++ {
		company := gen.GenerateCompany()
		names[company.Name] = true
	}

	assert.GreaterOrEqual(t, len(names), 45, "Should generate at least 45 unique company names out of 50 (found %d unique)", len(names))
}
