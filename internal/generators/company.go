package generators

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/diogomcd/fake-mill-api/internal/models"
)

// GenerateCompany generates complete fake company data
func (g *Generator) GenerateCompany() *models.CompanyResponse {
	ds := g.dataStore

	// Select a random company type (area)
	companyType := ds.GetRandomCompanyType()

	// Generate company name based on type
	prefix := companyType.Prefixes[rand.Intn(len(companyType.Prefixes))]
	suffix := companyType.Suffixes[rand.Intn(len(companyType.Suffixes))]
	companyWord := ds.GetRandomCompanyWord()
	legalSuffix := ds.GetRandomLegalSuffix()

	// Build company name with legal suffix
	companyName := fmt.Sprintf("%s %s %s %s", prefix, companyWord, suffix, legalSuffix)

	// Generate trade name (nome fantasia) based on pattern
	pattern := companyType.TradeNamePatterns[rand.Intn(len(companyType.TradeNamePatterns))]
	tradeName := strings.ReplaceAll(pattern, "{word}", companyWord)

	// CNPJ number
	cnpj := g.GenerateCNPJ(true, true)

	// Address (must be generated first to get state for phone DDD)
	address := g.GenerateAddress("", "")

	// Email, phone, shareCapital
	email := g.generateCompanyEmail(tradeName)
	phone := g.generateCompanyPhone(address.State)

	// Company size (porte)
	companySize := ds.GetRandomCompanySize()
	shareCapital := g.generateCompanyShareCapital(companySize)

	// Opening date (data de abertura) - 1 to 30 years ago
	yearsAgo := 1 + rand.Intn(30)
	openingDate := time.Now().AddDate(-yearsAgo, -rand.Intn(12), -rand.Intn(28)).Format("2006-01-02")

	// Foundation date (same as opening date)
	foundedAt := openingDate

	return &models.CompanyResponse{
		Name:         companyName,
		TradeName:    tradeName,
		CNPJ:         cnpj,
		Email:        email,
		Phone:        phone,
		Area:         companyType.Area,
		Size:         companySize.Name,
		OpeningDate:  openingDate,
		ShareCapital: shareCapital,
		FoundedAt:    foundedAt,
		Address:      *address,
	}
}

// generateCompanyEmail generates email for the company as object
func (g *Generator) generateCompanyEmail(tradeName string) models.CompanyEmail {
	ds := g.dataStore

	// Clean trade name (remove special characters and spaces)
	cleanTradeName := strings.ToLower(tradeName)
	cleanTradeName = strings.ReplaceAll(cleanTradeName, " ", "")
	cleanTradeName = strings.ReplaceAll(cleanTradeName, "á", "a")
	cleanTradeName = strings.ReplaceAll(cleanTradeName, "é", "e")
	cleanTradeName = strings.ReplaceAll(cleanTradeName, "í", "i")
	cleanTradeName = strings.ReplaceAll(cleanTradeName, "ó", "o")
	cleanTradeName = strings.ReplaceAll(cleanTradeName, "ú", "u")
	cleanTradeName = strings.ReplaceAll(cleanTradeName, "â", "a")
	cleanTradeName = strings.ReplaceAll(cleanTradeName, "ê", "e")
	cleanTradeName = strings.ReplaceAll(cleanTradeName, "ô", "o")
	cleanTradeName = strings.ReplaceAll(cleanTradeName, "ã", "a")
	cleanTradeName = strings.ReplaceAll(cleanTradeName, "õ", "o")
	cleanTradeName = strings.ReplaceAll(cleanTradeName, "ç", "c")

	// Get random domain extension
	domainExtension := ds.GetRandomEmailExtension()
	domain := fmt.Sprintf("%s.%s", cleanTradeName, domainExtension)
	emailAddress := fmt.Sprintf("contato@%s", domain)

	return models.CompanyEmail{
		Address: emailAddress,
		Domain:  domain,
	}
}

// generateCompanyPhone generates phone for the company as object
func (g *Generator) generateCompanyPhone(stateCode string) models.CompanyPhone {
	formatted, _, _, _ := g.GeneratePhone(stateCode, "landline")

	// For Brazil
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

	return models.CompanyPhone{
		InternationalFormat: internationalFormat,
		NationalFormat:      nationalFormat,
		CountryCode:         countryCode,
		DDI:                 ddi,
		E164Format:          e164Format,
	}
}

// generateCompanyShareCapital generates share capital for the company as object
func (g *Generator) generateCompanyShareCapital(companySize CompanySizeData) models.CompanyShareCapital {
	// Share capital (capital social) based on company size
	capitalMin := companySize.MinCapital
	capitalMax := companySize.MaxCapital
	capitalValue := capitalMin + rand.Float64()*(capitalMax-capitalMin)

	// Format with thousands separator and decimal places (Brazilian format)
	// Unformatted: "1234.56"
	unformatted := fmt.Sprintf("%.2f", capitalValue)

	// FormattedWithoutR: "1.234,56"
	formattedWithoutR := formatBRLCurrency(capitalValue, false)

	// Formatted: "R$ 1.234,56"
	formatted := formatBRLCurrency(capitalValue, true)

	return models.CompanyShareCapital{
		Formatted:         formatted,
		Unformatted:       unformatted,
		FormattedWithoutR: formattedWithoutR,
		Value:             capitalValue,
	}
}

// formatBRLCurrency formats a float value as Brazilian currency
func formatBRLCurrency(value float64, includeSymbol bool) string {
	// Convert to string with 2 decimal places
	str := fmt.Sprintf("%.2f", value)

	// Split integer and decimal parts
	parts := strings.Split(str, ".")
	integerPart := parts[0]
	decimalPart := parts[1]

	// Add thousands separator
	var result strings.Builder
	for i, digit := range integerPart {
		if i > 0 && (len(integerPart)-i)%3 == 0 {
			result.WriteRune('.')
		}
		result.WriteRune(digit)
	}

	// Build final string
	formatted := fmt.Sprintf("%s,%s", result.String(), decimalPart)

	if includeSymbol {
		return fmt.Sprintf("R$ %s", formatted)
	}
	return formatted
}
