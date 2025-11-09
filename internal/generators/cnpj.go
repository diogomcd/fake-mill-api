package generators

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

// GenerateCNPJ generates a valid or invalid CNPJ
func (g *Generator) GenerateCNPJ(formatted bool, valid bool) string {
	if !valid {
		// Generate invalid CNPJ - only random numbers
		cnpjStr := ""
		for i := 0; i < 14; i++ {
			cnpjStr += strconv.Itoa(rand.Intn(10))
		}
		if formatted {
			return FormatCNPJ(cnpjStr)
		}
		return cnpjStr
	}

	// Generate 12 random digits
	cnpj := make([]int, 12)
	for i := 0; i < 8; i++ {
		cnpj[i] = rand.Intn(10)
	}
	// Define the matrix as 0001
	cnpj[8] = 0
	cnpj[9] = 0
	cnpj[10] = 0
	cnpj[11] = 1

	// Calculate first check digit
	firstCheck := calculateCNPJCheckDigit(cnpj, []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2})
	cnpj = append(cnpj, firstCheck)

	// Calculate second check digit
	secondCheck := calculateCNPJCheckDigit(cnpj, []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2})
	cnpj = append(cnpj, secondCheck)

	cnpjStr := ""
	for _, digit := range cnpj {
		cnpjStr += strconv.Itoa(digit)
	}

	if formatted {
		return FormatCNPJ(cnpjStr)
	}
	return cnpjStr
}

// calculateCNPJCheckDigit calculates the check digit of the CNPJ
func calculateCNPJCheckDigit(cnpj []int, weights []int) int {
	sum := 0
	for i, digit := range cnpj {
		sum += digit * weights[i]
	}

	remainder := sum % 11
	if remainder < 2 {
		return 0
	}
	return 11 - remainder
}

// FormatCNPJ formats the CNPJ in the XX.XXX.XXX/XXXX-XX format
func FormatCNPJ(cnpj string) string {
	if len(cnpj) != 14 {
		return cnpj
	}
	return fmt.Sprintf("%s.%s.%s/%s-%s", cnpj[0:2], cnpj[2:5], cnpj[5:8], cnpj[8:12], cnpj[12:14])
}

// ValidateCNPJ validates a CNPJ
func ValidateCNPJ(cnpj string) bool {
	cleanCNPJ := strings.ReplaceAll(cnpj, ".", "")
	cleanCNPJ = strings.ReplaceAll(cleanCNPJ, "/", "")
	cleanCNPJ = strings.ReplaceAll(cleanCNPJ, "-", "")

	if len(cleanCNPJ) != 14 {
		return false
	}

	digits := make([]int, 14)
	for i, c := range cleanCNPJ {
		digit, err := strconv.Atoi(string(c))
		if err != nil {
			return false
		}
		digits[i] = digit
	}

	// Check if all digits are equal (invalid CNPJ)
	allEqual := true
	for i := 1; i < 14; i++ {
		if digits[i] != digits[0] {
			allEqual = false
			break
		}
	}
	if allEqual {
		return false
	}

	// Validate first digit
	firstCheck := calculateCNPJCheckDigit(digits[:12], []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2})
	if digits[12] != firstCheck {
		return false
	}

	// Validate second digit
	secondCheck := calculateCNPJCheckDigit(digits[:13], []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2})
	if digits[13] != secondCheck {
		return false
	}

	return true
}
