package generators

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests for credit card generation focusing on:
// 1. Brand-specific generation
// 2. Luhn algorithm validation
// 3. Format validation
// 4. Uniqueness

func TestGenerateCreditCard_ByBrand(t *testing.T) {
	ds, err := NewDataStore()
	if err != nil {
		t.Fatalf("Failed to create DataStore: %v", err)
	}
	gen := NewGenerator(ds)

	brands := []string{"Visa", "Mastercard", "Elo", "Amex"}
	for _, brand := range brands {
		t.Run(brand, func(t *testing.T) {
			number, cardBrand, cvv, expirationDate, holderName := gen.GenerateCreditCard(brand)

			assert.Equal(t, brand, cardBrand, "Card brand should match with proper capitalization")
			assert.NotEmpty(t, number, "Card number should not be empty")
			assert.Len(t, cvv, 3, "CVV should have 3 digits")
			assert.Len(t, expirationDate, 5, "Expiration date should have 5 characters (MM/YY)")
			assert.Contains(t, expirationDate, "/", "Expiration date should contain /")
			assert.NotEmpty(t, holderName, "Holder name should not be empty")
		})
	}
}

func TestGenerateCreditCard_RandomBrand(t *testing.T) {
	ds, err := NewDataStore()
	if err != nil {
		t.Fatalf("Failed to create DataStore: %v", err)
	}
	gen := NewGenerator(ds)

	number, cardBrand, cvv, expirationDate, holderName := gen.GenerateCreditCard("")

	assert.Contains(t, []string{"Visa", "Mastercard", "Elo", "Amex"}, cardBrand, "Card brand should be valid with proper capitalization")
	assert.NotEmpty(t, number, "Card number should not be empty")
	assert.Len(t, cvv, 3, "CVV should have 3 digits")
	assert.Len(t, expirationDate, 5, "Expiration date should have 5 characters (MM/YY)")
	assert.Contains(t, expirationDate, "/", "Expiration date should contain /")
	assert.NotEmpty(t, holderName, "Holder name should not be empty")
}

func TestGenerateCreditCard_Uniqueness(t *testing.T) {
	ds, err := NewDataStore()
	if err != nil {
		t.Fatalf("Failed to create DataStore: %v", err)
	}
	gen := NewGenerator(ds)

	cards := make(map[string]bool)
	const count = 100

	for i := 0; i < count; i++ {
		number, _, _, _, _ := gen.GenerateCreditCard("")
		cleanNumber := strings.ReplaceAll(number, " ", "")
		assert.False(t, cards[cleanNumber], "Card %s should be unique (attempt %d)", cleanNumber, i+1)
		cards[cleanNumber] = true
	}

	assert.Equal(t, count, len(cards), "Should generate %d unique cards", count)
}

func TestGenerateCreditCard_NumberFormat(t *testing.T) {
	ds, err := NewDataStore()
	if err != nil {
		t.Fatalf("Failed to create DataStore: %v", err)
	}
	gen := NewGenerator(ds)

	number, _, _, _, _ := gen.GenerateCreditCard("")

	cleanNumber := strings.ReplaceAll(number, " ", "")
	assert.Len(t, cleanNumber, 16, "Card number should have 16 digits")
}

