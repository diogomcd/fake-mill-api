package generators

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// GenerateRG generates a valid or invalid RG
func (g *Generator) GenerateRG(stateCode string, formatted bool, valid bool) (rg, state, issuer, issueDate, expirationDate string) {
	ds := g.dataStore

	if !valid {
		// Generate invalid RG - only random numbers
		rgNumber := ""
		for i := 0; i < 8; i++ {
			rgNumber += fmt.Sprintf("%d", rand.Intn(10))
		}
		checkDigit := fmt.Sprintf("%d", rand.Intn(10)) // Random digit for invalid RG

		rg = rgNumber + checkDigit
		if formatted {
			rg = FormatRG(rg)
		}

		issuer = "SSP"
		return
	}

	var selectedState *StateData
	if stateCode != "" {
		selectedState = ds.GetStateByCode(stateCode)
	}

	if selectedState == nil {
		selectedState = ds.GetRandomState()
	}

	state = selectedState.Code

	// Issue date: between 1 month ago and 10 years ago
	now := time.Now()
	minDaysAgo := 30   // minimum 1 month
	maxDaysAgo := 3650 // maximum 10 years
	daysAgo := minDaysAgo + rand.Intn(maxDaysAgo-minDaysAgo+1)
	issueTime := now.AddDate(0, 0, -daysAgo)
	expirationTime := issueTime.AddDate(10, 0, 0) // 10 years of validity

	issueDate = issueTime.Format("2006-01-02")
	expirationDate = expirationTime.Format("2006-01-02")

	// Generate 8 random digits
	rgNumber := ""
	for i := 0; i < 8; i++ {
		rgNumber += fmt.Sprintf("%d", rand.Intn(10))
	}

	checkDigit := calculateRGCheckDigit(rgNumber)
	rg = rgNumber + checkDigit

	if formatted {
		rg = FormatRG(rg)
	}

	issuer = "SSP" // Secretaria de Segurança Pública

	return
}

// FormatRG formats the RG to the standard XX.XXX.XXX-X
func FormatRG(rg string) string {
	cleanRG := CleanRG(rg)
	if len(cleanRG) != 9 {
		return rg
	}
	return fmt.Sprintf("%s.%s.%s-%s",
		cleanRG[0:2],
		cleanRG[2:5],
		cleanRG[5:8],
		cleanRG[8:9],
	)
}

// CleanRG removes the formatting from the RG
func CleanRG(rg string) string {
	cleanRG := strings.ReplaceAll(rg, ".", "")
	cleanRG = strings.ReplaceAll(cleanRG, "-", "")
	cleanRG = strings.ToUpper(cleanRG)
	return cleanRG
}

// calculateRGCheckDigit calculates the check digit of the RG
// Returns the check digit as a string (can be "0" to "9" or "X")
func calculateRGCheckDigit(rg string) string {
	weights := []int{2, 3, 4, 5, 6, 7, 8, 9}
	sum := 0

	// Multiply the 8 digits by the weights from 2 to 9
	for i, c := range rg {
		digit := int(c - '0')
		sum += digit * weights[i]
	}

	// Divide by 11 and get the remainder
	remainder := sum % 11

	// Subtract the remainder from 11
	result := 11 - remainder

	// Special cases:
	// If result is 10, check digit is X
	// If result is 11, check digit is 0
	if result == 10 {
		return "X"
	}
	if result == 11 {
		return "0"
	}

	return fmt.Sprintf("%d", result)
}

// ValidateRG validates an RG by checking the check digit
func ValidateRG(rg string) bool {
	cleanRG := CleanRG(rg)

	// Check length (8 digits + 1 check digit)
	if len(cleanRG) != 9 {
		return false
	}

	// Separate the first 8 digits and the check digit
	rgNumber := cleanRG[0:8]
	checkDigit := cleanRG[8:9]

	// Verify if the first 8 characters are digits
	for _, c := range rgNumber {
		if c < '0' || c > '9' {
			return false
		}
	}

	// Verify if the check digit is a digit or X
	if checkDigit != "X" && (checkDigit[0] < '0' || checkDigit[0] > '9') {
		return false
	}

	// Calculate the correct check digit
	expectedCheckDigit := calculateRGCheckDigit(rgNumber)

	return expectedCheckDigit == checkDigit
}
