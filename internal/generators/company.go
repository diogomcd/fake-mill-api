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

	// Company name
	companyNouns := []string{"Solutions", "Systems", "Tech", "Digital", "Consulting", "Group", "Enterprises", "Mill", "Cooperativa", "Empresa", "Comércio", "Soluções"}
	companySuffixes := []string{"LTDA", "S.A.", "ME", "EIRELI"}
	companyName := fmt.Sprintf("%s %s %s", ds.GetRandomLastName(), companyNouns[rand.Intn(len(companyNouns))], companySuffixes[rand.Intn(len(companySuffixes))])
	tradeName := strings.Split(companyName, " ")[0] + " " + companyNouns[rand.Intn(len(companyNouns))]

	// CNPJ number
	cnpj := g.GenerateCNPJ(true, true)

	// State registration (simplified) - maximum 14 characters
	stateRegistration := fmt.Sprintf("%02d.%03d.%03d.%03d", rand.Intn(100), rand.Intn(1000), rand.Intn(1000), rand.Intn(1000))

	// Email, phone, address
	email, _, _ := g.GenerateEmail("")
	phone, _, _, _ := g.GeneratePhone("", "landline")
	address := g.GenerateAddress("", "")

	// Foundation date (1 to 20 years ago)
	foundedAt := time.Now().AddDate(-(1 + rand.Intn(20)), rand.Intn(12)+1, rand.Intn(28)+1).Format("2006-01-02")

	return &models.CompanyResponse{
		Name:              companyName,
		TradeName:         tradeName,
		CNPJ:              cnpj,
		StateRegistration: stateRegistration,
		Email:             email,
		Phone:             phone,
		FoundedAt:         foundedAt,
		Address:           *address,
	}
}
