package generators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests for bank account generation focusing on:
// 1. Valid bank codes
// 2. Agency and account format
// 3. Account type validation

func TestGenerateBankAccount_ValidBank(t *testing.T) {
	gen := NewGenerator(nil)

	bank, agency, account, accountType := gen.GenerateBankAccount("001")

	assert.Equal(t, "001", bank.Code, "Bank code should match")
	assert.Equal(t, "Banco do Brasil", bank.Name, "Bank name should match")
	assert.NotEmpty(t, agency, "Agency should not be empty")
	assert.NotEmpty(t, account, "Account should not be empty")
	assert.Contains(t, []string{"checking", "savings"}, accountType, "Account type should be checking or savings")
}

func TestGenerateBankAccount_RandomBank(t *testing.T) {
	gen := NewGenerator(nil)

	bank, agency, account, accountType := gen.GenerateBankAccount("")

	assert.NotEmpty(t, bank.Code, "Bank code should not be empty")
	assert.NotEmpty(t, bank.Name, "Bank name should not be empty")
	assert.NotEmpty(t, agency, "Agency should not be empty")
	assert.NotEmpty(t, account, "Account should not be empty")
	assert.Contains(t, []string{"checking", "savings"}, accountType, "Account type should be checking or savings")
}

func TestGenerateBankAccount_AgencyFormat(t *testing.T) {
	gen := NewGenerator(nil)

	_, agency, _, _ := gen.GenerateBankAccount("")

	assert.Contains(t, agency, "-", "Agency should contain dash")
	assert.GreaterOrEqual(t, len(agency), 5, "Agency should have at least 5 characters")
}

func TestGenerateBankAccount_AccountFormat(t *testing.T) {
	gen := NewGenerator(nil)

	_, _, account, _ := gen.GenerateBankAccount("")

	assert.Contains(t, account, "-", "Account should contain dash")
	assert.GreaterOrEqual(t, len(account), 9, "Account should have at least 9 characters")
}

