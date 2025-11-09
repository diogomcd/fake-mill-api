package generators

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Tests for person generation focusing on:
// 1. Complete person generation
// 2. Gender filtering
// 3. State filtering
// 4. Age/birthdate consistency
// 5. BMI/height/weight consistency
// 6. Name diversity

func TestGeneratePerson_Complete(t *testing.T) {
	ds, err := NewDataStore()
	if err != nil {
		t.Fatalf("Failed to create DataStore: %v", err)
	}
	gen := NewGenerator(ds)

	person := gen.GeneratePerson("", "")

	assert.NotEmpty(t, person.Name.FullName, "Full name should not be empty")
	assert.NotEmpty(t, person.CPF.Masked, "CPF should not be empty")
	assert.NotEmpty(t, person.RG.Masked, "RG should not be empty")
	assert.NotEmpty(t, person.Email.Address, "Email should not be empty")
	assert.NotEmpty(t, person.Phone.InternationalFormat, "Phone should not be empty")
	assert.NotEmpty(t, person.Address.Street, "Address should not be empty")
}

func TestGeneratePerson_GenderFilter(t *testing.T) {
	ds, err := NewDataStore()
	if err != nil {
		t.Fatalf("Failed to create DataStore: %v", err)
	}
	gen := NewGenerator(ds)

	malePerson := gen.GeneratePerson("male", "")
	assert.Equal(t, "male", malePerson.Gender, "Gender should be male")

	femalePerson := gen.GeneratePerson("female", "")
	assert.Equal(t, "female", femalePerson.Gender, "Gender should be female")
}

func TestGeneratePerson_StateFilter(t *testing.T) {
	ds, err := NewDataStore()
	if err != nil {
		t.Fatalf("Failed to create DataStore: %v", err)
	}
	gen := NewGenerator(ds)

	person := gen.GeneratePerson("", "SP")
	assert.Equal(t, "SP", person.Address.State, "Address state should match requested state")
	assert.Equal(t, "SP", person.RG.State, "RG state should match requested state")
}

func TestGeneratePerson_AgeBirthdateConsistency(t *testing.T) {
	ds, err := NewDataStore()
	if err != nil {
		t.Fatalf("Failed to create DataStore: %v", err)
	}
	gen := NewGenerator(ds)

	person := gen.GeneratePerson("", "")

	birthdate, err := time.Parse("2006-01-02", person.Birthdate)
	assert.NoError(t, err, "Birthdate should be valid")

	expectedAge := calculateAgeFromBirthdate(birthdate)
	assert.Equal(t, person.Age, expectedAge, "Age should match calculated age from birthdate")
}

func TestGeneratePerson_BMIConsistency(t *testing.T) {
	ds, err := NewDataStore()
	if err != nil {
		t.Fatalf("Failed to create DataStore: %v", err)
	}
	gen := NewGenerator(ds)

	person := gen.GeneratePerson("", "")

	expectedBMI := person.Weight.Kilograms / (person.Height.Meters * person.Height.Meters)
	assert.InDelta(t, expectedBMI, person.BMI, 0.1, "BMI should match calculated BMI")
}

func TestGeneratePerson_NameDiversity(t *testing.T) {
	ds, err := NewDataStore()
	if err != nil {
		t.Fatalf("Failed to create DataStore: %v", err)
	}
	gen := NewGenerator(ds)

	names := make(map[string]bool)
	const count = 50

	for i := 0; i < count; i++ {
		person := gen.GeneratePerson("", "")
		names[person.Name.FullName] = true
	}

	assert.Greater(t, len(names), 40, "Should generate at least 40 unique names")
}

func calculateAgeFromBirthdate(birthdate time.Time) int {
	now := time.Now()
	age := now.Year() - birthdate.Year()
	if now.YearDay() < birthdate.YearDay() {
		age--
	}
	return age
}

