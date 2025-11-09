package generators

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/diogomcd/fake-mill-api/internal/models"
)

// GeneratePerson generates complete fake person data
func (g *Generator) GeneratePerson(gender, stateCode string) *models.Person {
	ds := g.dataStore

	if gender == "" || (gender != "male" && gender != "female") {
		genders := []string{"male", "female"}
		gender = genders[rand.Intn(len(genders))]
	}

	// Choose consistent state for RG, address and phone
	var selectedState *StateData
	if stateCode != "" {
		selectedState = ds.GetStateByCode(stateCode)
	}
	if selectedState == nil {
		selectedState = ds.GetRandomState()
	}
	actualStateCode := selectedState.Code

	personName := g.generatePersonName(gender)
	filiation := g.generateFiliation()
	cpf := g.generatePersonCPF()
	rg := g.generatePersonRG(actualStateCode)
	birthdate := generateBirthdate()
	age := calculateAge(birthdate)
	height := generateHeight(gender)
	weight := generateWeight(gender)
	bmi := calculateBMI(weight.Kilograms, height.Meters)
	zodiacSign := ds.GetZodiacSign(birthdate)
	favoriteColor := ds.GetRandomColor()
	bloodType := ds.GetRandomBloodType()
	email := g.generatePersonEmail(personName.FirstName, personName.LastName)
	phone := g.generatePersonPhone(actualStateCode)
	address := g.GenerateAddress(actualStateCode, "")
	profession := g.generatePersonProfession()
	company := g.generatePersonCompany()
	education := ds.GetRandomEducationLevel()
	maritalStatus := ds.GetRandomMaritalStatus()
	birthCity := ds.GetRandomCity(actualStateCode)

	return &models.Person{
		Name:          personName,
		CPF:           cpf,
		RG:            rg,
		Birthdate:     birthdate,
		Age:           age,
		Gender:        gender,
		Height:        height,
		Weight:        weight,
		BMI:           bmi,
		ZodiacSign:    zodiacSign,
		FavoriteColor: favoriteColor,
		BloodType:     bloodType,
		Filiation:     filiation,
		Email:         email,
		Phone:         phone,
		Address:       *address,
		Profession:    profession,
		Company:       company,
		Education:     education,
		MaritalStatus: maritalStatus,
		BirthCity:     birthCity,
	}
}

// generatePersonProfession generates professional information for the person
func (g *Generator) generatePersonProfession() models.PersonProfession {
	profession := g.dataStore.GetRandomProfession()

	return models.PersonProfession{
		Title: profession.Title,
		Area:  profession.Area,
	}
}

// generatePersonCompany generates company information for the person
func (g *Generator) generatePersonCompany() models.PersonCompany {
	company := g.GenerateCompany()

	return models.PersonCompany{
		Name: company.Name,
		CNPJ: company.CNPJ,
	}
}

// generatePersonName generates full name of the person
func (g *Generator) generatePersonName(gender string) models.PersonName {
	ds := g.dataStore

	// Generate first names (1 or 2 with 20% chance)
	var firstNames []string
	firstNameCount := 1
	if rand.Intn(100) < 20 { // 20% chance of having 2 first names
		firstNameCount = 2
	}

	for i := 0; i < firstNameCount; i++ {
		if gender == "male" {
			firstNames = append(firstNames, ds.GetRandomMaleFirstName())
		} else {
			firstNames = append(firstNames, ds.GetRandomFemaleFirstName())
		}
	}

	// Generate last names (2 or 3)
	lastNameCount := 2 + rand.Intn(2) // 2 or 3 last names
	var lastNames []string
	for i := 0; i < lastNameCount; i++ {
		lastNames = append(lastNames, ds.GetRandomLastName())
	}

	// Combine first names and last names
	firstName := strings.Join(firstNames, " ")
	lastName := strings.Join(lastNames, " ")
	fullName := firstName + " " + lastName

	return models.PersonName{
		FirstName: firstName,
		LastName:  lastName,
		FullName:  fullName,
	}
}

// generateFiliation generates family information (father and mother)
func (g *Generator) generateFiliation() models.Filiation {
	ds := g.dataStore

	// Father: male name + 2 last names
	fatherFirstName := ds.GetRandomMaleFirstName()
	fatherLastNames := []string{ds.GetRandomLastName(), ds.GetRandomLastName()}
	father := fatherFirstName + " " + strings.Join(fatherLastNames, " ")

	// Mother: female name + 2 last names
	motherFirstName := ds.GetRandomFemaleFirstName()
	motherLastNames := []string{ds.GetRandomLastName(), ds.GetRandomLastName()}
	mother := motherFirstName + " " + strings.Join(motherLastNames, " ")

	return models.Filiation{
		Father: father,
		Mother: mother,
	}
}

// generatePersonCPF generates CPF as object
func (g *Generator) generatePersonCPF() models.PersonCPF {
	cpfUnmasked := g.GenerateCPF(false, true)
	cpfMasked := FormatCPF(cpfUnmasked)

	return models.PersonCPF{
		Masked:   cpfMasked,
		Unmasked: cpfUnmasked,
	}
}

// generatePersonRG generates RG as object
func (g *Generator) generatePersonRG(stateCode string) models.PersonRG {
	rgUnmasked, state, issuer, issueDate, expirationDate := g.GenerateRG(stateCode, false, true)
	rgMasked := FormatRG(rgUnmasked)

	return models.PersonRG{
		Masked:         rgMasked,
		Unmasked:       rgUnmasked,
		State:          state,
		Issuer:         issuer,
		IssueDate:      issueDate,
		ExpirationDate: expirationDate,
	}
}

// generateHeight generates height in different metrics
func generateHeight(gender string) models.Height {
	var heightCm float64
	if gender == "male" {
		// Men: average 175cm, deviation of 7cm
		heightCm = 175 + (rand.NormFloat64() * 7)
	} else {
		// Women: average 162cm, deviation of 6cm
		heightCm = 162 + (rand.NormFloat64() * 6)
	}

	// Limit reasonable values
	if heightCm < 140 {
		heightCm = 140
	} else if heightCm > 200 {
		heightCm = 200
	}

	heightM := heightCm / 100
	heightInches := heightCm * 0.393701
	heightFeet := heightInches / 12

	return models.Height{
		Centimeters: math.Round(heightCm*100) / 100,
		Meters:      math.Round(heightM*100) / 100,
		Inches:      math.Round(heightInches*100) / 100,
		Feet:        math.Round(heightFeet*100) / 100,
	}
}

// generateWeight generates weight based on gender with realistic variation
func generateWeight(gender string) models.Weight {
	var minWeight, maxWeight float64

	if gender == "male" {
		minWeight = 60
		maxWeight = 100
	} else {
		minWeight = 50
		maxWeight = 85
	}

	// Generate weight with more concentrated distribution in the middle
	weightKg := minWeight + (maxWeight-minWeight)*rand.Float64()
	// Add small Gaussian variation for more realism
	weightKg += rand.NormFloat64() * 5

	weightPounds := weightKg * 2.20462
	weightGrams := weightKg * 1000

	return models.Weight{
		Kilograms: math.Round(weightKg*100) / 100,
		Pounds:    math.Round(weightPounds*100) / 100,
		Grams:     math.Round(weightGrams*100) / 100,
	}
}

// generatePersonEmail generates email related to name with password
func (g *Generator) generatePersonEmail(firstName, lastName string) models.PersonEmail {
	// Generate email using the centralized logic that now uses dynamic domains
	email, _, _ := g.GenerateEmail("")

	// Generate random password
	password := generatePassword()

	return models.PersonEmail{
		Address:  email,
		Password: password,
	}
}

// calculateBMI calculates the Body Mass Index
func calculateBMI(weightKg, heightM float64) float64 {
	if heightM <= 0 {
		return 0
	}
	bmi := weightKg / (heightM * heightM)
	return math.Round(bmi*100) / 100
}

// generatePassword generates a random password
func generatePassword() string {
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	length := 8 + rand.Intn(8) // 8-16 characters
	var password strings.Builder
	for i := 0; i < length; i++ {
		password.WriteByte(chars[rand.Intn(len(chars))])
	}
	return password.String()
}

// generatePersonPhone generates phone as object
func (g *Generator) generatePersonPhone(stateCode string) models.PersonPhone {
	formatted, _, _, _ := g.GeneratePhone(stateCode, "")

	// For temporary purposes, let's assume Brazil
	countryCode := "BR"
	ddi := 55

	unformatted := strings.ReplaceAll(formatted, "(", "")
	unformatted = strings.ReplaceAll(unformatted, ")", "")
	unformatted = strings.ReplaceAll(unformatted, " ", "")
	unformatted = strings.ReplaceAll(unformatted, "-", "")

	// International formats
	internationalFormat := fmt.Sprintf("+%d %s", ddi, formatted)
	nationalFormat := formatted
	e164Format := fmt.Sprintf("+%d%s", ddi, unformatted)

	return models.PersonPhone{
		InternationalFormat: internationalFormat,
		NationalFormat:      nationalFormat,
		CountryCode:         countryCode,
		DDI:                 ddi,
		E164Format:          e164Format,
	}
}

// generateBirthdate generates a random birthdate
func generateBirthdate() string {
	minAge := 18
	maxAge := 80
	age := minAge + rand.Intn(maxAge-minAge+1)

	now := time.Now()
	birthYear := now.Year() - age
	birthMonth := 1 + rand.Intn(12)
	birthDay := 1 + rand.Intn(28) // Simplified, uses 28 to avoid problems with days of the month

	return time.Date(birthYear, time.Month(birthMonth), birthDay, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
}

// calculateAge calculates the age based on the birthdate
func calculateAge(birthdate string) int {
	t, err := time.Parse("2006-01-02", birthdate)
	if err != nil {
		return 0
	}

	now := time.Now()
	age := now.Year() - t.Year()

	if now.Month() < t.Month() || (now.Month() == t.Month() && now.Day() < t.Day()) {
		age--
	}

	return age
}
