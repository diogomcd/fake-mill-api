package generators

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

// GenerateCPF generates a valid or invalid CPF
func (g *Generator) GenerateCPF(formatted bool, valid bool) string {
	if !valid {
		// Generate invalid CPF - only random numbers
		cpfStr := ""
		for i := 0; i < 11; i++ {
			cpfStr += strconv.Itoa(rand.Intn(10))
		}
		if formatted {
			return FormatCPF(cpfStr)
		}
		return cpfStr
	}

	// Generate 9 random digits
	cpf := make([]int, 9)
	for i := 0; i < 9; i++ {
		cpf[i] = rand.Intn(10)
	}

	// Calculate first check digit
	firstCheck := calculateCPFCheckDigit(cpf)
	cpf = append(cpf, firstCheck)

	// Calculate second check digit
	secondCheck := calculateCPFCheckDigit(cpf[:10])
	cpf = append(cpf, secondCheck)

	cpfStr := ""
	for _, digit := range cpf {
		cpfStr += strconv.Itoa(digit)
	}

	if formatted {
		return FormatCPF(cpfStr)
	}
	return cpfStr
}

// calculateCPFCheckDigit calculates the check digit of the CPF
func calculateCPFCheckDigit(cpf []int) int {
	sum := 0
	multiplier := len(cpf) + 1

	for _, digit := range cpf {
		sum += digit * multiplier
		multiplier--
	}

	remainder := sum % 11
	if remainder < 2 {
		return 0
	}
	return 11 - remainder
}

// FormatCPF formats the CPF in the XXX.XXX.XXX-XX format
func FormatCPF(cpf string) string {
	if len(cpf) != 11 {
		return cpf
	}
	return fmt.Sprintf("%s.%s.%s-%s", cpf[0:3], cpf[3:6], cpf[6:9], cpf[9:11])
}

// ValidateCPF validates a CPF format (with or without mask)
func ValidateCPF(cpf string) bool {
	cleanCPF := strings.ReplaceAll(cpf, ".", "")
	cleanCPF = strings.ReplaceAll(cleanCPF, "-", "")

	if len(cleanCPF) != 11 {
		return false
	}

	digits := make([]int, 11)
	for i, c := range cleanCPF {
		digit, err := strconv.Atoi(string(c))
		if err != nil {
			return false
		}
		digits[i] = digit
	}

	// Check if all digits are equal (invalid CPF)
	allEqual := true
	for i := 1; i < 11; i++ {
		if digits[i] != digits[0] {
			allEqual = false
			break
		}
	}
	if allEqual {
		return false
	}

	// Validate first digit
	firstCheck := calculateCPFCheckDigit(digits[:9])
	if digits[9] != firstCheck {
		return false
	}

	// Validate second digit
	secondCheck := calculateCPFCheckDigit(digits[:10])
	if digits[10] != secondCheck {
		return false
	}

	return true
}
