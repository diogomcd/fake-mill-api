package generators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests for address generation focusing on:
// 1. Complete address generation
// 2. State-specific generation
// 3. Street variability
// 4. Coordinate validation

func TestGenerateAddress_Complete(t *testing.T) {
	ds, err := NewDataStore()
	if err != nil {
		t.Fatalf("Failed to create DataStore: %v", err)
	}
	gen := NewGenerator(ds)

	address := gen.GenerateAddress("", "")

	assert.NotEmpty(t, address.Street, "Street should not be empty")
	assert.NotEmpty(t, address.Number, "Number should not be empty")
	assert.NotEmpty(t, address.Neighborhood, "Neighborhood should not be empty")
	assert.NotEmpty(t, address.City, "City should not be empty")
	assert.NotEmpty(t, address.State, "State should not be empty")
	assert.NotEmpty(t, address.Zipcode, "Zipcode should not be empty")
	assert.NotNil(t, address.Coordinates, "Coordinates should not be nil")
	assert.GreaterOrEqual(t, address.Coordinates.Lat, -90.0, "Latitude should be >= -90")
	assert.LessOrEqual(t, address.Coordinates.Lat, 90.0, "Latitude should be <= 90")
	assert.GreaterOrEqual(t, address.Coordinates.Lng, -180.0, "Longitude should be >= -180")
	assert.LessOrEqual(t, address.Coordinates.Lng, 180.0, "Longitude should be <= 180")
}

func TestGenerateAddress_ByState(t *testing.T) {
	ds, err := NewDataStore()
	if err != nil {
		t.Fatalf("Failed to create DataStore: %v", err)
	}
	gen := NewGenerator(ds)

	states := []string{"SP", "RJ", "MG"}
	for _, state := range states {
		t.Run(state, func(t *testing.T) {
			address := gen.GenerateAddress(state, "")
			assert.Equal(t, state, address.State, "Address state should match requested state")
		})
	}
}

func TestGenerateAddress_StreetVariability(t *testing.T) {
	ds, err := NewDataStore()
	if err != nil {
		t.Fatalf("Failed to create DataStore: %v", err)
	}
	gen := NewGenerator(ds)

	streets := make(map[string]bool)
	const count = 50

	for i := 0; i < count; i++ {
		address := gen.GenerateAddress("SP", "")
		streets[address.Street] = true
	}

	assert.Greater(t, len(streets), 1, "Should generate varied street names")
}

func TestGenerateAddress_ZipcodeFormat(t *testing.T) {
	ds, err := NewDataStore()
	if err != nil {
		t.Fatalf("Failed to create DataStore: %v", err)
	}
	gen := NewGenerator(ds)

	address := gen.GenerateAddress("SP", "")

	assert.Len(t, address.Zipcode, 9, "Zipcode should have 9 characters (XXXXX-XXX)")
	assert.Contains(t, address.Zipcode, "-", "Zipcode should contain dash")
}

