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

	// State registration (simplified) - maximum 14 characters
	stateRegistration := fmt.Sprintf("%02d.%03d.%03d.%03d", rand.Intn(100), rand.Intn(1000), rand.Intn(1000), rand.Intn(1000))

	// Email, phone, address
	email, _, _ := g.GenerateEmail("")
	phone, _, _, _ := g.GeneratePhone("", "landline")
	address := g.GenerateAddress("", "")

	// Company size (porte)
	companySize := ds.GetRandomCompanySize()

	// Share capital (capital social) based on company size
	capitalMin := companySize.MinCapital
	capitalMax := companySize.MaxCapital
	capitalValue := capitalMin + rand.Float64()*(capitalMax-capitalMin)
	shareCapital := fmt.Sprintf("R$ %.2f", capitalValue)

	// Opening date (data de abertura) - 1 to 30 years ago
	yearsAgo := 1 + rand.Intn(30)
	openingDate := time.Now().AddDate(-yearsAgo, -rand.Intn(12), -rand.Intn(28)).Format("2006-01-02")

	// Foundation date (same as opening date)
	foundedAt := openingDate

	return &models.CompanyResponse{
		Name:              companyName,
		TradeName:         tradeName,
		CNPJ:              cnpj,
		StateRegistration: stateRegistration,
		Email:             email,
		Phone:             phone,
		Area:              companyType.Area,
		Size:              companySize.Name,
		OpeningDate:       openingDate,
		ShareCapital:      shareCapital,
		FoundedAt:         foundedAt,
		Address:           *address,
	}
}
